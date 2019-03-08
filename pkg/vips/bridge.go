package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	dbg "runtime/debug"
	"strings"
	"unsafe"
)

const (
	defaultQuality     = 90
	defaultCompression = 6
	maxScaleFactor     = 10
	DefaultDPI         = 72
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

var stringBuffer4096 = fixedString(4096)

func vipsOperationNew(name string) *C.VipsOperation {
	cName := C.CString(name)
	defer freeCString(cName)
	return C.vips_operation_new(cName)
}

func vipsCall(name string, options []*Option) error {
	operation := vipsOperationNew(name)
	return vipsCallOperation(operation, options)
}

func vipsCallOperation(operation *C.VipsOperation, options []*Option) error {
	// Set the inputs
	for _, option := range options {
		if option.Output() {
			continue
		}
		defer option.Close()

		cName := C.CString(option.Name)
		defer freeCString(cName)
		C.gobject_set_property(
			(*C.VipsObject)(unsafe.Pointer(operation)), cName, option.GValue())
	}

	if ret := C.vips_cache_operation_buildp(&operation); ret != 0 {
		return handleVipsError()
	}

	defer C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))

	// Write back the outputs
	for _, option := range options {
		if !option.Output() {
			continue
		}
		defer option.Close()
		cName := C.CString(option.Name)
		defer freeCString(cName)

		C.g_object_get_property(
			(*C.GObject)(unsafe.Pointer(operation)), (*C.gchar)(cName), option.GValue())
	}

	return nil
}

func vipsHasProfile(input *C.VipsImage) bool {
	return int(C.has_profile_embed(input)) > 0
}

func vipsPrepareForExport(input *C.VipsImage, params *ExportParams) (*C.VipsImage, error) {
	if params.StripProfile && vipsHasProfile(input) {
		C.remove_icc_profile(input)
	}

	if params.Quality == 0 {
		params.Quality = defaultQuality
	}

	if params.Compression == 0 {
		params.Compression = defaultCompression
	}

	// Use a default interpretation and cast it to C type
	if params.Interpretation == 0 {
		params.Interpretation = Interpretation(input.Type)
	}

	interpretation := C.VipsInterpretation(params.Interpretation)

	// Apply the proper colour space
	if int(C.is_colorspace_supported(input)) == 1 && interpretation != input.Type {
		var out *C.VipsImage
		err := C.to_colorspace(input, &out, interpretation)
		if int(err) != 0 {
			return nil, handleVipsError()
		}
		input = out
	}

	return input, nil
}

func vipsLoadFromBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	// Reference buf here so it's not garbage collected during image initialization.
	defer runtime.KeepAlive(buf)

	var image *C.VipsImage
	imageType := vipsDetermineImageType(buf)

	if imageType == ImageTypeUnknown {
		if len(buf) > 2 {
			log.Printf("Failed to understand image format size=%d %x %x %x", len(buf), buf[0], buf[1], buf[2])
		} else {
			log.Printf("Failed to understand image format size=%d", len(buf))
		}
		return nil, ImageTypeUnknown, ErrUnsupportedImageFormat
	}

	len := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	err := C.init_image(imageBuf, len, C.int(imageType), &image)
	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}

	return image, imageType, nil
}

func vipsExportBuffer(image *C.VipsImage, params *ExportParams) ([]byte, ImageType, error) {
	tmpImage, err := vipsPrepareForExport(image, params)

	if err != nil {
		return nil, ImageTypeUnknown, err
	}

	// If these are equal, then we don't want to deref the original image as
	// the original will be returned if the target colorspace is not supported
	if tmpImage != image {
		defer C.g_object_unref(C.gpointer(tmpImage))
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
	default:
		incOpCounter("save_jpeg_buffer")
		format = ImageTypeJPEG
		cErr = C.save_jpeg_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality, interlaced)
	}

	if int(cErr) != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}

	buf := C.GoBytes(ptr, C.int(cLen))
	C.g_free(C.gpointer(ptr))
	return buf, format, nil
}

func isTypeSupported(imageType ImageType) bool {
	return supportedImageTypes[imageType]
}

func isColorspaceIsSupportedBuffer(buf []byte) (bool, error) {
	image, _, err := vipsLoadFromBuffer(buf)
	if err != nil {
		return false, err
	}
	defer C.g_object_unref(C.gpointer(image))
	return int(C.is_colorspace_supported(image)) == 1, nil
}

func isColorspaceIsSupported(image *C.VipsImage) bool {
	return int(C.is_colorspace_supported(image)) == 1
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
	if IsTypeSupported(ImageTypeSVG) && buf[0] == 0x3c && buf[1] == 0x3f && buf[2] == 0x78 && buf[3] == 0x6d {
		return ImageTypeSVG
	}
	return ImageTypeUnknown
}

func vipsFlattenBackground(input *C.VipsImage, color Color) (*C.VipsImage, error) {
	incOpCounter("flatten")
	var output *C.VipsImage

	bg := [3]C.double{
		C.double(color.R),
		C.double(color.G),
		C.double(color.B),
	}

	if int(C.has_alpha_channel(input)) > 0 {
		err := C.flatten_image_background(input, &output, bg[0], bg[1], bg[2])
		if int(err) != 0 {
			return nil, handleVipsError()
		}
		C.g_object_unref(C.gpointer(input))
		input = output
	}

	return input, nil
}

func vipsResize(input *C.VipsImage, scale, vscale float64, kernel Kernel) (*C.VipsImage, error) {
	incOpCounter("resize")
	var output *C.VipsImage

	// Let's not be insane
	scale = math.Min(scale, maxScaleFactor)
	vscale = math.Min(vscale, maxScaleFactor)

	defer C.g_object_unref(C.gpointer(input))
	if err := C.resize_image(input, &output, C.double(scale), C.double(vscale), C.int(kernel)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsRotate(input *C.VipsImage, angle Angle) (*C.VipsImage, error) {
	incOpCounter("rot")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.rot_image(input, &output, C.VipsAngle(angle)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsComposite(inputs []*C.VipsImage, mode BlendMode) (*C.VipsImage, error) {
	incOpCounter("composite")
	var output *C.VipsImage
	if err := C.composite(&inputs[0], &output, C.int(len(inputs)), C.int(mode)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsHasAlpha(image *C.VipsImage) bool {
	return int(C.has_alpha_channel(image)) > 0
}

func vipsAdd(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("add")
	defer C.g_object_unref(C.gpointer(left))
	defer C.g_object_unref(C.gpointer(right))
	var output *C.VipsImage
	if err := C.add(left, right, &output); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsMultiply(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("multiply")
	defer C.g_object_unref(C.gpointer(left))
	defer C.g_object_unref(C.gpointer(right))
	var output *C.VipsImage
	if err := C.multiply(left, right, &output); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsExtractBand(image *C.VipsImage, band, num int) (*C.VipsImage, error) {
	incOpCounter("extract")
	defer C.g_object_unref(C.gpointer(image))
	var output *C.VipsImage
	if err := C.extract_band(image, &output, C.int(band), C.int(num)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsLinear1(image *C.VipsImage, a, b float64) (*C.VipsImage, error) {
	incOpCounter("linear1")
	defer C.g_object_unref(C.gpointer(image))
	var output *C.VipsImage
	if err := C.linear1(image, &output, C.double(a), C.double(b)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsExtractArea(input *C.VipsImage, left, top, width, height int) (*C.VipsImage, error) {
	incOpCounter("extract")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.extract_image_area(input, &output, C.int(left), C.int(top), C.int(width), C.int(height)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsEmbed(input *C.VipsImage, left, top, width, height int, extend Extend) (*C.VipsImage, error) {
	incOpCounter("embed")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.embed_image(input, &output, C.int(left), C.int(top), C.int(width), C.int(height), C.int(extend), 0, 0, 0); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsFlip(input *C.VipsImage, dir Direction) (*C.VipsImage, error) {
	incOpCounter("flip")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.flip_image(input, &output, C.int(dir)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsInvert(input *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("invert")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.invert_image(input, &output); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsGaussianBlur(input *C.VipsImage, sigma float64) (*C.VipsImage, error) {
	incOpCounter("gaussblur")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.gaussian_blur(input, &output, C.double(sigma)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsZoom(input *C.VipsImage, xFactor, yFactor int) (*C.VipsImage, error) {
	incOpCounter("zoom")
	var output *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))
	if err := C.zoom_image(input, &output, C.int(xFactor), C.int(yFactor)); err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func vipsLabel(input *C.VipsImage, lp LabelParams) (*C.VipsImage, error) {
	incOpCounter("label")

	var output *C.VipsImage

	text := C.CString(lp.Text)
	font := C.CString(lp.Font)
	color := [3]C.double{C.double(lp.Color.R), C.double(lp.Color.G), C.double(lp.Color.B)}
	w := lp.Width.GetRounded(int(input.Xsize))
	h := lp.Height.GetRounded(int(input.Ysize))
	offsetX := lp.OffsetX.GetRounded(int(input.Xsize))
	offsetY := lp.OffsetY.GetRounded(int(input.Ysize))

	opts := vipsLabelOptions{
		Text:      text,
		Font:      font,
		Width:     C.int(w),
		Height:    C.int(h),
		OffsetX:   C.int(offsetX),
		OffsetY:   C.int(offsetY),
		Alignment: C.VipsAlign(lp.Alignment),
		Opacity:   C.float(lp.Opacity),
		Color:     color,
	}

	defer C.free(unsafe.Pointer(text))
	defer C.free(unsafe.Pointer(font))

	err := C.label(input, &output, (*C.LabelOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleVipsError()
	}

	return output, nil
}

func handleVipsError() error {
	defer C.vips_thread_shutdown()
	defer C.vips_error_clear()

	s := C.GoString(C.vips_error_buffer())
	stack := string(dbg.Stack())
	return fmt.Errorf("%s\nStack:\n%s", s, stack)
}

//add new feature about watermarkImage and watermarkText
type WatermarkImage struct {
	Left      int
	Top       int
	Watermark *ImageRef
	Opacity   float32
}

type vipsWatermarkImageOptions struct {
	Left    C.int
	Top     C.int
	Opacity C.float
}

type WatermarkText struct {
	Text      string
	Font      string
	Opacity   float32
	Color     Color
	Alignment string
	Margin    int
	TextInfo  TextInfo
}

type TextInfo struct {
	ImageWidth  int
	ImageHeight int
	TextWidth   int
	TextHeight  int
	OffsetX     int
	OffsetY     int
	Crop        bool
	CropX       int
	CropY       int
}

type AlignPosition int

type AlignMode struct {
	Horizon  AlignPosition
	Vertical AlignPosition
}

const (
	HorizonLeft AlignPosition = iota
	HorizonCenter
	HorizonRight
	VerticalTop
	VerticalMiddle
	VerticalBottom
)

func catchVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}

func VipsDrawWatermark(image *C.VipsImage, o WatermarkImage) (*C.VipsImage, error) {
	var out *C.VipsImage

	//	watermark, _, e := vipsRead(o.Buf)
	//	if e != nil {
	//		return nil, e
	//	}

	opts := vipsWatermarkImageOptions{C.int(o.Left), C.int(o.Top), C.float(o.Opacity)}

	err := C.vips_watermark_image(image, o.Watermark.image, &out, (*C.WatermarkImageOptions)(unsafe.Pointer(&opts)))

	if err != 0 {
		return nil, catchVipsError()
	}

	return out, nil
}

func (in *ImageRef) VipsDrawWatermark(o WatermarkImage) error {
	out, err := VipsDrawWatermark(in.image, o)
	if err != nil {
		return err
	}
	in.SetImage(out)
	return nil
}

func GetText(input *ImageRef, wt WatermarkText) (*ImageRef, error) {
	var output *C.VipsImage
	text := C.CString(wt.Text)
	font := C.CString(wt.Font)
	w := wt.TextInfo.ImageWidth

	defer C.free(unsafe.Pointer(text))
	defer C.free(unsafe.Pointer(font))
	opts := vipsLabelOptions{
		Text:  text,
		Font:  font,
		Width: C.int(w),
		DPI:   C.int(DefaultDPI),
	}

	err := C.get_text(&output, (*C.LabelOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleVipsError()
	}
	textRef := NewImageRef(output, ImageTypeJPEG)
	return textRef, nil
}

func GetText1(input *ImageRef, ti *ImageRef, wt WatermarkText) (*ImageRef, error) {
	var output *C.VipsImage

	opts := vipsLabelOptions{
		OffsetX: C.int(wt.TextInfo.OffsetX),
		OffsetY: C.int(wt.TextInfo.OffsetY),
		Opacity: C.float(wt.Opacity),
	}

	err := C.get_text1(ti.image, &output, (*C.LabelOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleVipsError()
	}
	textRef := NewImageRef(output, ImageTypeJPEG)
	return textRef, nil

}

func watermarkText(input *C.VipsImage, t *C.VipsImage, wt WatermarkText) (*C.VipsImage, error) {
	var output *C.VipsImage
	color := [3]C.double{C.double(wt.Color.R), C.double(wt.Color.G), C.double(wt.Color.B)}
	opts := vipsLabelOptions{
		Color: color,
	}

	err := C.watermarkText(input, t, &output, (*C.LabelOptions)(unsafe.Pointer(&opts)))
	if err != 0 {
		return nil, handleVipsError()
	}
	return output, nil
}

func (ti *TextInfo) getTextInfo(in *ImageRef, text *ImageRef) {
	ti.ImageWidth = in.Width()
	ti.ImageHeight = in.Height()
	ti.TextWidth = text.Width()
	ti.TextHeight = text.Height()
	ti.Crop = false
}

func (textRef *ImageRef) TextCrop(ti TextInfo, margin int) error {
	var extractWidth, extractHeight int
	if ti.CropX != 0 {
		extractWidth = ti.ImageWidth - margin
	} else {
		extractWidth = ti.TextWidth + ti.OffsetX
	}

	if ti.CropY != 0 {
		extractHeight = ti.ImageHeight - margin
	} else {
		extractHeight = ti.TextHeight + ti.OffsetY
	}

	return textRef.ExtractArea(0, 0, extractWidth, extractHeight)
}

func parseAlign(align string) AlignMode {
	position := AlignMode{
		Horizon:  HorizonRight,
		Vertical: VerticalBottom,
	}

	if align == "" {
		return position
	}
	pos := strings.Split(align, ",")
	for _, p := range pos {
		switch p {
		case "top":
			position.Vertical = VerticalTop
		case "middle":
			position.Vertical = VerticalMiddle
		case "bottom":
			position.Vertical = VerticalBottom
		case "left":
			position.Horizon = HorizonLeft
		case "center":
			position.Horizon = HorizonCenter
		case "right":
			position.Horizon = HorizonRight
		default:
			continue
		}
	}

	return position
}

func (ti *TextInfo) updateTextInfo(margin int, align string) {
	if ti.ImageWidth <= ti.TextWidth+2*margin {
		ti.OffsetX = margin
		ti.CropX = ti.TextWidth + 2*margin - ti.ImageWidth
		ti.Crop = true
	}

	if ti.ImageHeight <= ti.TextHeight+2*margin {
		ti.OffsetY = margin
		ti.CropY = ti.TextHeight + 2*margin - ti.ImageHeight
		ti.Crop = true
	}

	pos := parseAlign(align)

	if ti.CropX == 0 {
		switch pos.Horizon {
		case HorizonLeft:
			ti.OffsetX = margin
		case HorizonCenter:
			ti.OffsetX = (ti.ImageWidth - ti.TextWidth) / 2

		case HorizonRight:
			ti.OffsetX = ti.ImageWidth - ti.TextWidth - margin
		default:
			ti.OffsetX = ti.ImageWidth - ti.TextWidth - margin
		}
	}

	if ti.CropY == 0 {
		switch pos.Vertical {
		case VerticalTop:
			ti.OffsetY = margin
		case VerticalMiddle:
			ti.OffsetY = (ti.ImageHeight - ti.TextHeight) / 2
		case VerticalBottom:
			ti.OffsetY = ti.ImageHeight - ti.TextHeight - margin
		default:
			ti.OffsetY = ti.ImageHeight - ti.TextHeight - margin
		}
	}
}

func (wt *WatermarkText) updateWatermarkTextParameter(in *ImageRef, text *ImageRef) {
	var ti = &TextInfo{}
	ti.getTextInfo(in, text)
	ti.updateTextInfo(wt.Margin, wt.Alignment)
	wt.TextInfo = *ti
}

func (in *ImageRef) WatermarkText(wt WatermarkText) error {
	text, err := GetText(in, wt)
	if err != nil {
		return err
	}

	wt.updateWatermarkTextParameter(in, text)

	textRef, err := GetText1(in, text, wt)
	if err != nil {
		return err
	}

	if wt.TextInfo.Crop {
		if err = textRef.TextCrop(wt.TextInfo, wt.Margin); err != nil {
			return err
		}
	}

	out, err := watermarkText(in.image, textRef.image, wt)
	if err != nil {
		return err
	}
	in.SetImage(out)

	return nil
}
