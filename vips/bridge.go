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

func (i Image) Shrink(hShrink, vShrink float64) (*Image, error) {
	defer C.g_object_unref(C.gpointer(i.image))

	var out *C.VipsImage

	ret := C.shrink(
		i.image,
		&out,
		C.double(float64(hShrink)),
		C.double(float64(vShrink)))

	if ret != 0 {
		return nil, handleVipsError()
	}

	return i.SetImage(out), nil
}
