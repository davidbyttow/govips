package vips

import "unsafe"

func cPtr(b []byte) unsafe.Pointer {
	return unsafe.Pointer(&b[0])
}
