package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

import (
	"fmt"
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

func fromGboolean(b C.gboolean) bool {
	if b != 0 {
		return false
	}
	return true
}

func fixedString(size int) string {
	return fmt.Sprintf("%%%d.%ds", size, size)
}
