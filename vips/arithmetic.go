package vips

// #include "arithmetic.h"
import "C"
import "unsafe"

// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-find-trim
func vipsFindTrim(in *C.VipsImage, threshold float64, backgroundColor *Color) (int, int, int, int, error) {
	incOpCounter("findTrim")
	var left, top, width, height C.int

	if err := C.find_trim(in, &left, &top, &width, &height, C.double(threshold), C.double(backgroundColor.R),
		C.double(backgroundColor.G), C.double(backgroundColor.B)); err != 0 {
		return -1, -1, -1, -1, handleVipsError()
	}

	return int(left), int(top), int(width), int(height), nil
}

// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html#vips-getpoint
func vipsGetPoint(in *C.VipsImage, n int, x int, y int) ([]float64, error) {
	incOpCounter("getpoint")
	var out *C.double
	defer gFreePointer(unsafe.Pointer(out))

	if err := C.getpoint(in, &out, C.int(n), C.int(x), C.int(y)); err != 0 {
		return nil, handleVipsError()
	}

	// maximum n is 4
	return (*[4]float64)(unsafe.Pointer(out))[:n:n], nil
}

// https://www.libvips.org/API/current/libvips-arithmetic.html#vips-min
func vipsMin(in *C.VipsImage) (float64, int, int, error) {
	incOpCounter("min")
	var out C.double
	var x, y C.int

	if err := C.minOp(in, &out, &x, &y, C.int(1)); err != 0 {
		return 0, 0, 0, handleVipsError()
	}

	return float64(out), int(x), int(y), nil
}
