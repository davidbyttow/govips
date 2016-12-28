package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import "errors"

func loadBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	imageType := DetermineImageType(buf)

	if imageType == ImageTypeUnknown {
		return nil, ImageTypeUnknown, errors.New("Unsupported image format")
	}

	var image *C.VipsImage
	err := C.init_image(cPtr(buf),
		C.size_t(len(buf)),
		C.int(imageType),
		&image)

	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}
	return image, imageType, nil
}

func loadJpegBuffer(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}

func loadWebpBuffer(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}
