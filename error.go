package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import "errors"

var (
	// ErrUnsupportedImageFormat when image type is unsupported
	ErrUnsupportedImageFormat = errors.New("UnsupportedImageFormat")
)
