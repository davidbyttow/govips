package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"unsafe"

	go_debug "github.com/tj/go-debug"
)

func byteArrayPointer(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&b[0])
}

func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
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
	if b != 0 {
		return true
	}
	return false
}

func fixedString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = '0'
	}
	return string(b)
}

var debugFn = go_debug.Debug("govips")

func debug(fmt string, values ...interface{}) {
	if len(values) > 0 {
		debugFn(fmt, values...)
	} else {
		debugFn(fmt)
	}
}
