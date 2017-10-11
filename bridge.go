package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	dbg "runtime/debug"
	"unsafe"
)

const (
	defaultQuality     = 90
	defaultCompression = 6
)

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
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))

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

func vipsPrepareForExport(input *C.VipsImage, params *ExportParams) (*C.VipsImage, error) {
	if params.StripProfile {
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
		defer C.g_object_unref(C.gpointer(input))
		err := C.to_colorspace(input, &out, interpretation)
		if int(err) != 0 {
			return nil, handleVipsError()
		}
		input = out
	}

	return input, nil
}

func vipsLoadFromBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	var image *C.VipsImage
	imageType := vipsDetermineImageType(buf)

	if imageType == ImageTypeUnknown {
		return nil, ImageTypeUnknown, errors.New("Unsupported image format")
	}

	len := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	err := C.init_image(imageBuf, len, C.int(imageType), &image)
	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}

	runtime.KeepAlive(buf)

	return image, imageType, nil
}

func vipsExportBuffer(image *C.VipsImage, params *ExportParams) ([]byte, error) {
	tmpImage, err := vipsPrepareForExport(image, params)
	if err != nil {
		return nil, err
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
	stripMetadata := C.int(boolToInt(params.StripMetadata))

	if params.Format != ImageTypeUnknown && !IsTypeSupported(params.Format) {
		return nil, fmt.Errorf("cannot save to %#v", imageTypes[params.Format])
	}

	var ptr unsafe.Pointer
	switch params.Format {
	case ImageTypeWEBP:
		cErr = C.save_webp_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality)
	case ImageTypePNG:
		cErr = C.save_png_buffer(tmpImage, &ptr, &cLen, stripMetadata, C.int(params.Compression), quality, interlaced)
	case ImageTypeTIFF:
		cErr = C.save_tiff_buffer(tmpImage, &ptr, &cLen)
	default:
		cErr = C.save_jpeg_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality, interlaced)
	}

	if int(cErr) != 0 {
		return nil, handleVipsError()
	}

	buf := C.GoBytes(ptr, C.int(cLen))

	runtime.KeepAlive(ptr)

	C.g_free(C.gpointer(ptr))
	C.vips_error_clear()

	return buf, nil
}

func isTypeSupported(imageType ImageType) bool {
	return supportedImageTypes[imageType]
}

func isColorspaceIsSupportedBuffer(buf []byte) (bool, error) {
	image, _, err := vipsLoadFromBuffer(buf)
	if err != nil {
		return false, err
	}
	C.g_object_unref(C.gpointer(image))
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
	return ImageTypeUnknown
}

func vipsShrinkJPEG(buf []byte, input *C.VipsImage, shrink int) (*C.VipsImage, error) {
	var image *C.VipsImage
	var ptr = unsafe.Pointer(&buf[0])
	defer C.g_object_unref(C.gpointer(input))

	err := C.load_jpeg_buffer(ptr, C.size_t(len(buf)), &image, C.int(shrink))
	if err != 0 {
		return nil, handleVipsError()
	}

	return image, nil
}

func vipsShrink(input *C.VipsImage, shrink int) (*C.VipsImage, error) {
	var image *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))

	err := C.shrink_image(input, &image, C.double(float64(shrink)), C.double(float64(shrink)))
	if err != 0 {
		return nil, handleVipsError()
	}

	return image, nil
}

func vipsFlattenBackground(input *C.VipsImage, color Color) (*C.VipsImage, error) {
	var image *C.VipsImage
	defer C.g_object_unref(C.gpointer(input))

	bg := [3]C.double{
		C.double(color.R),
		C.double(color.G),
		C.double(color.B),
	}

	if int(C.has_alpha_channel(input)) > 0 {
		err := C.flatten_image_background(input, &image, bg[0], bg[1], bg[2])
		if int(err) != 0 {
			return nil, handleVipsError()
		}
	}

	return image, nil
}

func handleVipsError() error {
	defer C.vips_thread_shutdown()
	defer C.vips_error_clear()

	s := C.GoString(C.vips_error_buffer())
	stack := string(dbg.Stack())
	return fmt.Errorf("%s\nStack:\n%s", s, stack)
}
