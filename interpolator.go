package govips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import "runtime"

// Interpolator represents an interpolator when resizing images
type Interpolator struct {
	interp *C.VipsInterpolate
}

// NewInterpolator creates a new Interpolator from the given name
func NewInterpolator(name string) (*Interpolator, error) {
	interp, err := vipsInterpolateNew(name)
	if err != nil {
		return nil, err
	}
	out := &Interpolator{
		interp: interp,
	}
	runtime.SetFinalizer(out, finalizeInterpolator)
	return out, nil
}

func finalizeInterpolator(i *Interpolator) {
	C.g_object_unref(C.gpointer(i.interp))
}
