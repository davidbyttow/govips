package vips

// #cgo pkg-config: vips
// #include "create.h"
import "C"

// XYZ loads an image buffer and creates a new Image
func XYZ(width int, height int) (*ImageRef, error) {
	startupIfNeeded()

	image, err := vipsXYZ(width, height)
	if err != nil {
		return nil, err
	}

	return newImageRef(image, ImageTypeBMP, nil), nil
}

func vipsXYZ(width int, height int) (*C.VipsImage, error) {
	var out *C.VipsImage

	if err := C.xyz(&out, C.int(width), C.int(height)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
