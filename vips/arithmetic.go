package vips

// #cgo pkg-config: vips
// #include "arithmetic.h"
import "C"

// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-add
func vipsAdd(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("add")
	var out *C.VipsImage

	if err := C.add(left, right, &out); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-multiply
func vipsMultiply(left *C.VipsImage, right *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("multiply")
	var out *C.VipsImage

	if err := C.multiply(left, right, &out); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

//  https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-linear1
func vipsLinear1(in *C.VipsImage, a, b float64) (*C.VipsImage, error) {
	incOpCounter("linear1")
	var out *C.VipsImage

	if err := C.linear1(in, &out, C.double(a), C.double(b)); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}

// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-invert
func vipsInvert(in *C.VipsImage) (*C.VipsImage, error) {
	incOpCounter("invert")
	var out *C.VipsImage

	if err := C.invert_image(in, &out); err != 0 {
		return nil, handleImageError(out)
	}

	return out, nil
}
