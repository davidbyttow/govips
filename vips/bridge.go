package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"fmt"
	"log"
	"math"
	"runtime"
	dbg "runtime/debug"
	"unsafe"
)

const (
	defaultQuality     = 90
	defaultCompression = 6
	maxScaleFactor     = 10
)

type vipsLabelOptions struct {
	Text      *C.char
	Font      *C.char
	Width     C.int
	Height    C.int
	OffsetX   C.int
	OffsetY   C.int
	Alignment C.VipsAlign
	DPI       C.int
	Margin    C.int
	Opacity   C.float
	Color     [3]C.double
}

type vipsLoadOptions struct {
	cOpts C.ImageLoadOptions
}

func vipsCall(name string, options []*Option) error {
	operation := vipsOperationNew(name)

	return vipsCallOperation(operation, options)
}

func vipsOperationNew(name string) *C.VipsOperation {
	cName := C.CString(name)
	defer freeCString(cName)

	return C.vips_operation_new(cName)
}

func vipsCallOperation(operation *C.VipsOperation, options []*Option) error {
	// todo: replace with https://jcupitt.github.io/libvips/API/current/VipsOperation.html#vips-cache-operation-build

	for _, option := range options {
		if option.Output() {
			continue
		}
		defer option.Close()

		cName := C.CString(option.Name)
		defer freeCString(cName)

		C.gobject_set_property((*C.VipsObject)(unsafe.Pointer(operation)), cName, option.GValue())
	}

	if ret := C.vips_cache_operation_buildp(&operation); ret != 0 {
		return handleVipsError(nil)
	}

	defer unrefPointer(unsafe.Pointer(operation))

	for _, option := range options {
		if !option.Output() {
			continue
		}
		defer option.Close()

		cName := C.CString(option.Name)
		defer freeCString(cName)

		C.g_object_get_property((*C.GObject)(unsafe.Pointer(operation)), (*C.gchar)(cName), option.GValue())
	}

	return nil
}

func vipsHasProfile(in *C.VipsImage) bool {
	return int(C.has_profile_embed(in)) > 0
}

func vipsLoadFromBuffer(buf []byte, opts ...LoadOption) (*C.VipsImage, ImageType, error) {
	// Reference buf here so it's not garbage collected during image initialization.
	defer runtime.KeepAlive(buf)

	var out *C.VipsImage

	imageType := vipsDetermineImageType(buf)
	if imageType == ImageTypeUnknown {
		if len(buf) > 2 {
			log.Printf("Failed to understand image format size=%d %x %x %x", len(buf), buf[0], buf[1], buf[2])
		} else {
			log.Printf("Failed to understand image format size=%d", len(buf))
		}
		return nil, ImageTypeUnknown, ErrUnsupportedImageFormat
	}

	bufLength := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	var loadOpts vipsLoadOptions
	for _, opt := range opts {
		opt(&loadOpts)
	}

	err := C.init_image(imageBuf, bufLength, C.int(imageType), &loadOpts.cOpts, &out)
	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError(out)
	}

	return out, imageType, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-copy
func vipsCopyImage(in *C.VipsImage) (*C.VipsImage, error) {
	var out *C.VipsImage

	err := C.copy_image(in, &out)
	if int(err) != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsExportBuffer(image *C.VipsImage, params *ExportParams) ([]byte, ImageType, error) {
	tmpImage, err := vipsPrepareForExport(image, params)
	if err != nil {
		return nil, ImageTypeUnknown, err
	}

	// If these are equal, then we don't want to deref the original image as
	// the original will be returned if the target colorspace is not supported
	if tmpImage != image {
		defer unrefImage(tmpImage)
	}

	cLen := C.size_t(0)
	var cErr C.int
	interlaced := C.int(boolToInt(params.Interlaced))
	quality := C.int(params.Quality)
	lossless := C.int(boolToInt(params.Lossless))
	stripMetadata := C.int(boolToInt(params.StripMetadata))
	format := params.Format

	if format != ImageTypeUnknown && !IsTypeSupported(format) {
		return nil, ImageTypeUnknown, fmt.Errorf("cannot save to %#v", ImageTypes[format])
	}

	if params.BackgroundColor != nil {
		tmpImage, err = vipsFlattenBackground(tmpImage, *params.BackgroundColor)
		if err != nil {
			return nil, ImageTypeUnknown, err
		}
	}

	var ptr unsafe.Pointer

	switch format {
	case ImageTypeWEBP:
		incOpCounter("save_webp_buffer")
		cErr = C.save_webp_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality, lossless)
	case ImageTypePNG:
		incOpCounter("save_png_buffer")
		cErr = C.save_png_buffer(tmpImage, &ptr, &cLen, stripMetadata, C.int(params.Compression), quality, interlaced)
	case ImageTypeTIFF:
		incOpCounter("save_tiff_buffer")
		cErr = C.save_tiff_buffer(tmpImage, &ptr, &cLen)
	case ImageTypeHEIF:
		incOpCounter("save_heif_buffer")
		cErr = C.save_heif_buffer(tmpImage, &ptr, &cLen, quality, lossless)
	default:
		incOpCounter("save_jpeg_buffer")
		format = ImageTypeJPEG
		cErr = C.save_jpeg_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality, interlaced)
	}

	if int(cErr) != 0 {
		return nil, ImageTypeUnknown, handleVipsError(nil)
	}

	buf := C.GoBytes(ptr, C.int(cLen))
	gFreePointer(ptr)

	return buf, format, nil
}

func vipsPrepareForExport(in *C.VipsImage, params *ExportParams) (*C.VipsImage, error) {
	if params.StripProfile && vipsHasProfile(in) {
		C.remove_icc_profile(in)
	}

	if params.Quality == 0 {
		params.Quality = defaultQuality
	}

	if params.Compression == 0 {
		params.Compression = defaultCompression
	}

	// Use a default interpretation and cast it to C type
	if params.Interpretation == 0 {
		params.Interpretation = Interpretation(in.Type)
	}

	interpretation := C.VipsInterpretation(params.Interpretation)

	// Apply the proper colour space
	if int(C.is_colorspace_supported(in)) == 1 && interpretation != in.Type {
		var out *C.VipsImage

		err := C.to_colorspace(in, &out, interpretation)
		if int(err) != 0 {
			return nil, handleVipsError(out)
		}

		return out, nil
	}

	return in, nil
}

func vipsDetermineImageType(buf []byte) ImageType {
	if len(buf) < 12 {
		return ImageTypeUnknown
	}
	if buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF {
		return ImageTypeJPEG
	}
	if IsTypeSupported(ImageTypeGIF) && buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46 {
		return ImageTypeGIF
	}
	if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		return ImageTypePNG
	}
	if IsTypeSupported(ImageTypeTIFF) &&
		((buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x0) ||
			(buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x0 && buf[3] == 0x2A)) {
		return ImageTypeTIFF
	}
	if IsTypeSupported(ImageTypePDF) && buf[0] == 0x25 && buf[1] == 0x50 && buf[2] == 0x44 && buf[3] == 0x46 {
		return ImageTypePDF
	}
	if IsTypeSupported(ImageTypeWEBP) && buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50 {
		return ImageTypeWEBP
	}
	if IsTypeSupported(ImageTypeSVG) && buf[0] == 0x3C && buf[1] == 0x3F && buf[2] == 0x78 && buf[3] == 0x6D {
		return ImageTypeSVG
	}
	// https://github.com/strukturag/libheif/blob/master/libheif/heif.cc
	if IsTypeSupported(ImageTypeHEIF) && (buf[4] == 'f' && buf[5] == 't' && buf[6] == 'y' && buf[7] == 'p') &&
		(buf[8] == 'h' && buf[9] == 'e' && buf[10] == 'i' && buf[11] == 'c') {
		return ImageTypeHEIF
	}
	return ImageTypeUnknown
}

func vipsFlattenBackground(in *C.VipsImage, color Color) (*C.VipsImage, error) {
	incOpCounter("flatten")
	var out *C.VipsImage

	if int(C.has_alpha_channel(in)) > 0 {

		bg := [3]C.double{
			C.double(color.R),
			C.double(color.G),
			C.double(color.B),
		}

		err := C.flatten_image_background(in, &out, bg[0], bg[1], bg[2])
		if int(err) != 0 {
			return nil, handleVipsError(out)
		}
		unrefImage(in)

		in = out
	}

	return in, nil
}

// Resize executes the 'resize' operation
func vipsResize(in *C.VipsImage, scale, vscale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var out *C.VipsImage

	// Let's not be insane
	scale = math.Min(scale, maxScaleFactor)
	vscale = math.Min(vscale, maxScaleFactor)

	if err := C.resize_image(in, &out, C.double(scale), C.double(vscale), C.int(kernel)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsEmbed(in *C.VipsImage, left, top, width, height int, extend ExtendStrategy) (*C.VipsImage, error) {
	incOpCounter("embed")
	var out *C.VipsImage

	if err := C.embed_image(in, &out, C.int(left), C.int(top), C.int(width), C.int(height), C.int(extend), 0, 0, 0); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-autorot
func vipsAutoRotate(in *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("autorot")
	var out *C.VipsImage

	if err := C.autorot_image(in, &out); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-autorot
func vipsRotate(in *C.VipsImage, angle Angle) (*C.VipsImage, error) {
	incOpCounter("rot")
	var out *C.VipsImage

	if err := C.rot_image(in, &out, C.VipsAngle(angle)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsZoom(in *C.VipsImage, xFactor, yFactor int) (*C.VipsImage, error) {
	incOpCounter("zoom")
	var out *C.VipsImage

	if err := C.zoom_image(in, &out, C.int(xFactor), C.int(yFactor)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-composite
func vipsComposite(ins []*C.VipsImage, mode BlendMode) (*C.VipsImage, error) {
	incOpCounter("composite")
	var out *C.VipsImage

	if err := C.composite(&ins[0], &out, C.int(len(ins)), C.int(mode)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsBandJoin(ins []*C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("bandJoin")
	var out *C.VipsImage

	if err := C.bandjoin(&ins[0], &out, C.int(len(ins))); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsHasAlpha(in *C.VipsImage) bool {
	return int(C.has_alpha_channel(in)) > 0
}

func vipsAddAlpha(in *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("addAlpha")
	var out *C.VipsImage

	if err := C.add_alpha(in, &out); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsAdd(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("add")
	var out *C.VipsImage

	defer unrefImage(left)
	defer unrefImage(right)

	if err := C.add(left, right, &out); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsMultiply(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("multiply")
	var out *C.VipsImage

	defer unrefImage(left)
	defer unrefImage(right)

	if err := C.multiply(left, right, &out); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-conversion.html#vips-extract-band
func vipsExtractBand(in *C.VipsImage, band, num int) (*C.VipsImage, error) {
	incOpCounter("extractBand")
	var out *C.VipsImage

	if err := C.extract_band(in, &out, C.int(band), C.int(num)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

//  https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-linear1
func vipsLinear1(in *C.VipsImage, a, b float64) (*C.VipsImage, error) {
	incOpCounter("linear1")
	var out *C.VipsImage

	if err := C.linear1(in, &out, C.double(a), C.double(b)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsExtractArea(in *C.VipsImage, left, top, width, height int) (*C.VipsImage, error) {
	incOpCounter("extractArea")
	var out *C.VipsImage

	if err := C.extract_image_area(in, &out, C.int(left), C.int(top), C.int(width), C.int(height)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsFlip(in *C.VipsImage, dir Direction) (*C.VipsImage, error) {
	incOpCounter("flip")
	var out *C.VipsImage

	if err := C.flip_image(in, &out, C.int(dir)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsInvert(in *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("invert")
	var out *C.VipsImage

	if err := C.invert_image(in, &out); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-convolution.html#vips-gaussblur
func vipsGaussianBlur(in *C.VipsImage, sigma float64) (*C.VipsImage, error) {
	incOpCounter("gaussblur")
	var out *C.VipsImage

	if err := C.gaussian_blur(in, &out, C.double(sigma)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-convolution.html#vips-sharpen
func vipsSharpen(in *C.VipsImage, sigma float64, x1 float64, m2 float64) (*C.VipsImage, error) {
	incOpCounter("sharpen")
	var out *C.VipsImage

	if err := C.sharpen(in, &out, C.double(sigma), C.double(x1), C.double(m2)); err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func vipsLabel(in *C.VipsImage, params *LabelParams) (*C.VipsImage, error) {
	incOpCounter("label")
	var out *C.VipsImage

	text := C.CString(params.Text)
	defer freeCString(text)

	font := C.CString(params.Font)
	defer freeCString(font)

	color := [3]C.double{C.double(params.Color.R), C.double(params.Color.G), C.double(params.Color.B)}
	w := params.Width.GetRounded(int(in.Xsize))
	h := params.Height.GetRounded(int(in.Ysize))
	offsetX := params.OffsetX.GetRounded(int(in.Xsize))
	offsetY := params.OffsetY.GetRounded(int(in.Ysize))

	opts := vipsLabelOptions{
		Text:      text,
		Font:      font,
		Width:     C.int(w),
		Height:    C.int(h),
		OffsetX:   C.int(offsetX),
		OffsetY:   C.int(offsetY),
		Alignment: C.VipsAlign(params.Alignment),
		Opacity:   C.float(params.Opacity),
		Color:     color,
	}

	err := C.label(in, &out, (*C.LabelOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleVipsError(out)
	}

	return out, nil
}

func handleVipsError(out *C.VipsImage) error {
	if out != nil {
		unrefImage(out)
	}

	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()

	stack := string(dbg.Stack())

	return fmt.Errorf("%v\nStack:\n%s", s, stack)
}
