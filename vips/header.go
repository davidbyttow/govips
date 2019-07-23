package vips

// #cgo pkg-config: vips
// #include "header.h"
import "C"

func vipsHasICCProfile(in *C.VipsImage) bool {
	return int(C.has_icc_profile(in)) > 0
}

func vipsRemoveICCProfile(in *C.VipsImage) bool {
	return fromGboolean(C.remove_icc_profile(in))
}
