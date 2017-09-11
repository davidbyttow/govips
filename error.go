package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import "errors"

var (
	// ErrUnsupportedImageFormat when image type is unsupported
	ErrUnsupportedImageFormat = errors.New("UnsupportedImageFormat")

	// ErrInvalidInterpolator when interpolator is invalid
	ErrInvalidInterpolator = errors.New("Invalid interpolator")
)

func handleVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}
