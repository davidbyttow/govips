package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"unsafe"
)

// Option is a type that is passed to internal libvips functions
type Option struct {
	Name   string
	gvalue C.GValue
	closer func(gv *C.GValue)
	output bool
}

// NewOption returns a new option instance
func NewOption(name string, gtype C.GType, output bool, closer func(gv *C.GValue)) *Option {
	v := &Option{
		Name:   name,
		output: output,
		closer: closer,
	}
	C.g_value_init(&v.gvalue, gtype)
	return v
}

// Output returns true if this option is an output-only option
func (v *Option) Output() bool {
	return v.output
}

// Close releases memory associated with this option
func (v *Option) Close() {
	if v.closer != nil {
		v.closer(&v.gvalue)
	}
	C.g_value_unset(&v.gvalue)
}

// GValue returns the internal gvalue type
func (v *Option) GValue() *C.GValue {
	return &v.gvalue
}

// InputBool represents a boolean input option
func InputBool(name string, v bool) *Option {
	o := NewOption(name, C.G_TYPE_BOOLEAN, false, nil)
	C.g_value_set_boolean(&o.gvalue, toGboolean(v))
	return o
}

// OutputBool represents a boolean output option
func OutputBool(name string, v *bool) *Option {
	o := NewOption(name, C.G_TYPE_BOOLEAN, true, func(gv *C.GValue) {
		*v = fromGboolean(C.g_value_get_boolean(gv))
	})
	return o
}

// InputInt represents a int input option
func InputInt(name string, v int) *Option {
	o := NewOption(name, C.G_TYPE_INT, false, nil)
	C.g_value_set_int(&o.gvalue, C.gint(v))
	return o
}

// OutputInt represents a int output option
func OutputInt(name string, v *int) *Option {
	o := NewOption(name, C.G_TYPE_INT, true, func(gv *C.GValue) {
		*v = int(C.g_value_get_int(gv))
	})
	return o
}

// InputDouble represents a float64 input option
func InputDouble(name string, v float64) *Option {
	o := NewOption(name, C.G_TYPE_DOUBLE, false, nil)
	C.g_value_set_double(&o.gvalue, C.gdouble(v))
	return o
}

// OutputDouble represents a float output option
func OutputDouble(name string, v *float64) *Option {
	o := NewOption(name, C.G_TYPE_DOUBLE, true, func(gv *C.GValue) {
		*v = float64(C.g_value_get_double(gv))
	})
	return o
}

// InputString represents a string input option
func InputString(name string, v string) *Option {
	cStr := C.CString(v)
	o := NewOption(name, C.G_TYPE_STRING, false, func(gv *C.GValue) {
		freeCString(cStr)
	})
	C.g_value_set_string(&o.gvalue, (*C.gchar)(cStr))
	return o
}

// OutputString represents a string output option
func OutputString(name string, v *string) *Option {
	o := NewOption(name, C.G_TYPE_STRING, true, func(gv *C.GValue) {
		*v = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(gv))))
	})
	return o
}

// InputImage represents a VipsImage input option
func InputImage(name string, v *C.VipsImage) *Option {
	o := NewOption(name, C.vips_image_get_type(), false, nil)
	C.g_value_set_object(&o.gvalue, C.gpointer(v))
	return o
}

// OutputImage represents a VipsImage output option
func OutputImage(name string, v **C.VipsImage) *Option {
	o := NewOption(name, C.vips_image_get_type(), true, func(gv *C.GValue) {
		*v = (*C.VipsImage)(C.g_value_get_object(gv))
	})
	return o
}

// InputInterpolator represents a Interpolator input option
func InputInterpolator(name string, interp Interpolator) *Option {
	cStr := C.CString(interp.String())
	defer freeCString(cStr)
	interpolator := C.vips_interpolate_new(cStr)

	o := NewOption(name, C.vips_interpolate_get_type(), false, func(gv *C.GValue) {
		defer C.g_object_unref(C.gpointer(interpolator))
	})
	C.g_value_set_object(&o.gvalue, C.gpointer(interpolator))
	return o
}
