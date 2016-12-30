package gimage

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"

import "unsafe"

var STRING_BUFFER = fixedString(4096)

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

func splitFilenameAndOptions(file string) (string, string) {
	c_file := C.CString(file)
	defer freeCString(c_file)

	c_filename := C.CString(STRING_BUFFER)
	defer freeCString(c_filename)

	c_optionString := C.CString(STRING_BUFFER)
	defer freeCString(c_optionString)

	C.vips__filename_split8(c_file, c_filename, c_optionString)

	fileName := C.GoString(c_filename)
	optionString := C.GoString(c_optionString)
	return fileName, optionString
}
