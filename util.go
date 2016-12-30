package gimage

// #cgo pkg-config: vips
// #include <stdlib.h>
// #include "vips/vips.h"
import "C"

import "unsafe"

func cPtr(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&b[0])
}

func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
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
	b := make([]byte, size)
	for i := range b {
		b[i] = '0'
	}
	return string(b)
}
