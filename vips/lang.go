package vips

// #cgo pkg-config: vips
// #include <vips/vips.h>
// #include <stdlib.h>
import "C"

import (
	"unsafe"
)

func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func gFreePointer(ref unsafe.Pointer) {
	C.g_free(C.gpointer(ref))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func toGboolean(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

func fromGboolean(b C.gboolean) bool {
	return b != 0
}
