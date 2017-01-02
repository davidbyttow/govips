package govips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type Interpolator struct {
	name   string
	interp *C.VipsInterpolate
}

func NewInterpolate(name string) (*Interpolator, error) {
	interp, err := vipsInterpolateNew(name)
	if err != nil {
		return nil, err
	}
	return &Interpolator{
		name:   name,
		interp: interp,
	}, nil
}

func (i Interpolator) Name() string {
	return i.name
}
