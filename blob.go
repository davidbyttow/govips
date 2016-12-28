package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type Blob struct {
	blob *C.VipsBlob
}

func NewBlob(buf []byte) *Blob {
	return &Blob{
		blob: nil,
	}
}
