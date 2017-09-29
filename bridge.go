package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	defaultQuality = 80
)

var stringBuffer4096 = fixedString(4096)

func vipsForeignFindLoad(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_load(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindLoadBuffer(bytes []byte) (string, error) {
	cOperationName := C.vips_foreign_find_load_buffer(
		byteArrayPointer(bytes),
		C.size_t(len(bytes)))
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSave(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsForeignFindSaveBuffer(filename string) (string, error) {
	cFilename := C.CString(filename)
	defer freeCString(cFilename)

	cOperationName := C.vips_foreign_find_save_buffer(cFilename)
	if cOperationName == nil {
		return "", ErrUnsupportedImageFormat
	}
	return C.GoString(cOperationName), nil
}

func vipsInterpolateNew(name string) (*C.VipsInterpolate, error) {
	cName := C.CString(name)
	defer freeCString(cName)

	interp := C.vips_interpolate_new(cName)
	if interp == nil {
		return nil, ErrInvalidInterpolator
	}
	return interp, nil
}

func vipsOperationNew(name string) *C.VipsOperation {
	cName := C.CString(name)
	defer freeCString(cName)
	return C.vips_operation_new(cName)
}

func vipsFilenameSplit8(file string) (string, string) {
	cFile := C.CString(file)
	defer freeCString(cFile)

	cFilename := C.CString(stringBuffer4096)
	defer freeCString(cFilename)

	cOptionString := C.CString(stringBuffer4096)
	defer freeCString(cOptionString)

	C.vips__filename_split8(cFile, cFilename, cOptionString)

	fileName := C.GoString(cFilename)
	optionString := C.GoString(cOptionString)
	return fileName, optionString
}

func vipsCallString(name string, options *Options, optionString string) error {
	operation := vipsOperationNew(name)
	//defer C.g_object_unref(C.gpointer(operation))

	if optionString != "" {
		cOptionString := C.CString(optionString)
		defer freeCString(cOptionString)

		if C.vips_object_set_from_string(
			(*C.VipsObject)(unsafe.Pointer(operation)),
			cOptionString) != 0 {
			return handleVipsError()
		}
	}
	return vipsCallOperation(operation, options)
}

func vipsCall(name string, options *Options) error {
	operation := vipsOperationNew(name)
	//defer C.g_object_unref(C.gpointer(operation))

	return vipsCallOperation(operation, options)
}

func vipsCallOperation(operation *C.VipsOperation, options *Options) error {
	// TODO(d): Unref the outputs
	if options != nil {
		for _, option := range options.Options {
			if option.IsOutput {
				continue
			}

			cName := C.CString(option.Name)
			defer freeCString(cName)

			C.gobject_set_property(
				(*C.VipsObject)(unsafe.Pointer(operation)),
				cName,
				&option.GValue)
		}
	}

	if ret := C.vips_cache_operation_buildp(&operation); ret != 0 {
		C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))
		return handleVipsError()
	}

	// We defer this here because the pointer may have changed.
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(operation)))

	if options != nil {
		for _, option := range options.Options {
			if !option.IsOutput {
				continue
			}

			cName := C.CString(option.Name)
			defer freeCString(cName)

			C.g_object_get_property(
				(*C.GObject)(unsafe.Pointer(operation)),
				(*C.gchar)(cName),
				&option.GValue)
			option.Deserialize()
		}
	}

	return nil
}

func vipsPrepareForExport(image *C.VipsImage, options *ExportOptions) (*C.VipsImage, error) {
	var outImage *C.VipsImage

	if options.StripProfile {
		C.remove_icc_profile(image)
	}

	if options.Quality == 0 {
		options.Quality = defaultQuality
	}

	// Use a default interpretation and cast it to C type
	if options.Interpretation == 0 {
		options.Interpretation = InterpretationSrgb
	}

	interpretation := C.VipsInterpretation(options.Interpretation)

	// Apply the proper colour space
	if int(C.is_colorspace_supported(image)) == 1 {
		err := C.to_colorspace(image, &outImage, interpretation)
		if int(err) != 0 {
			return nil, handleVipsError()
		}
		image = outImage
	}

	return image, nil
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

	return image, imageType, nil
}

func vipsExportBuffer(image *C.VipsImage, options *ExportOptions) ([]byte, error) {
	tmpImage, err := vipsPrepareForExport(image, options)
	if err != nil {
		return nil, err
	}

	// If these are equal, then we don't want to deref the original image as
	// the original will be returned if the target colorspace is not supported
	if tmpImage != image {
		defer C.g_object_unref(C.gpointer(tmpImage))
	}

	cLen := C.size_t(0)
	cErr := C.int(0)
	interlaced := C.int(boolToInt(options.Interlaced))
	quality := C.int(options.Quality)
	stripMetadata := C.int(boolToInt(options.StripMetadata))

	if options.Type != ImageTypeUnknown && !IsTypeSupported(options.Type) {
		return nil, fmt.Errorf("cannot save to %#v", imageTypes[options.Type])
	}

	var ptr unsafe.Pointer
	switch options.Type {
	case ImageTypeWEBP:
		cErr = C.save_webp_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality)
	case ImageTypePNG:
		cErr = C.save_png_buffer(tmpImage, &ptr, &cLen, stripMetadata, C.int(options.Compression), quality, interlaced)
	case ImageTypeTIFF:
		cErr = C.save_tiff_buffer(tmpImage, &ptr, &cLen)
	default:
		cErr = C.save_jpeg_buffer(tmpImage, &ptr, &cLen, stripMetadata, quality, interlaced)
	}

	if int(cErr) != 0 {
		return nil, handleVipsError()
	}

	buf := C.GoBytes(ptr, C.int(cLen))

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
