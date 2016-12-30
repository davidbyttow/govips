package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type Interpolator struct {
	name   string
	interp *C.VipsInterpolate
}

func NewInterpolate(name string) (*Interpolator, error) {
	interp := C.vips_interpolate_new(C.CString(name))
	if interp == nil {
		return nil, ErrInvalidInterpolator
	}
	return &Interpolator{
		name:   name,
		interp: interp,
	}, nil
}

func (i Interpolator) Name() string {
	return i.name
}
