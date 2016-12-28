package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

import (
	"unsafe"
)

func cPtr(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&b[0])
}

func toGboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}
