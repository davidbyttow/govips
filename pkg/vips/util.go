package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import (
	"log"
	"unsafe"
)

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
  return b != 0
}

func debug(fmt string, values ...interface{}) {
	if len(values) > 0 {
		log.Printf(fmt, values...)
	} else {
		log.Print(fmt)
	}
}
