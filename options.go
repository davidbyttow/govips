package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"unsafe"
)

type VipsFunc func(*C.GValue)

type VipsOption struct {
	Name   string
	gvalue C.GValue
	closer func(gv *C.GValue)
	output bool
}

func NewOption(name string, gtype C.GType, output bool, closer func(gv *C.GValue)) *VipsOption {
	v := &VipsOption{
		Name:   name,
		output: output,
		closer: closer,
	}
	C.g_value_init(&v.gvalue, gtype)
	return v
}

func (v *VipsOption) Output() bool {
	return v.output
}

func (v *VipsOption) Close() {
	if v.closer != nil {
		v.closer(&v.gvalue)
	}
	C.g_value_unset(&v.gvalue)
}

func (v *VipsOption) GValue() *C.GValue {
	return &v.gvalue
}

func InputBool(name string, v bool) *VipsOption {
	o := NewOption(name, C.G_TYPE_BOOLEAN, false, nil)
	C.g_value_set_boolean(&o.gvalue, toGboolean(v))
	return o
}

func OutputBool(name string, v *bool) *VipsOption {
	o := NewOption(name, C.G_TYPE_BOOLEAN, true, func(gv *C.GValue) {
		*v = fromGboolean(C.g_value_get_boolean(gv))
	})
	return o
}

func InputInt(name string, v int) *VipsOption {
	o := NewOption(name, C.G_TYPE_INT, false, nil)
	C.g_value_set_int(&o.gvalue, C.gint(v))
	return o
}

func OutputInt(name string, v *int) *VipsOption {
	o := NewOption(name, C.G_TYPE_INT, true, func(gv *C.GValue) {
		*v = int(C.g_value_get_int(gv))
	})
	return o
}

func InputDouble(name string, v float64) *VipsOption {
	o := NewOption(name, C.G_TYPE_DOUBLE, false, nil)
	C.g_value_set_double(&o.gvalue, C.gdouble(v))
	return o
}

func OutputDouble(name string, v *float64) *VipsOption {
	o := NewOption(name, C.G_TYPE_DOUBLE, true, func(gv *C.GValue) {
		*v = float64(C.g_value_get_double(gv))
	})
	return o
}

func InputString(name string, v string) *VipsOption {
	cStr := C.CString(v)
	o := NewOption(name, C.G_TYPE_STRING, false, func(gv *C.GValue) {
		freeCString(cStr)
	})
	C.g_value_set_string(&o.gvalue, (*C.gchar)(cStr))
	return o
}

func OutputString(name string, v *string) *VipsOption {
	o := NewOption(name, C.G_TYPE_STRING, true, func(gv *C.GValue) {
		*v = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(gv))))
	})
	return o
}

func InputImage(name string, v *C.VipsImage) *VipsOption {
	o := NewOption(name, C.vips_image_get_type(), false, nil)
	C.g_value_set_object(&o.gvalue, C.gpointer(v))
	return o
}

func OutputImage(name string, v **C.VipsImage) *VipsOption {
	o := NewOption(name, C.vips_image_get_type(), true, func(gv *C.GValue) {
		*v = (*C.VipsImage)(C.g_value_get_object(gv))
	})
	return o
}

func InputInterpolator(name string, interp Interpolator) *VipsOption {
	cStr := C.CString(interp.String())
	defer freeCString(cStr)
	interpolator := C.vips_interpolate_new(cStr)

	o := NewOption(name, C.vips_interpolate_get_type(), false, func(gv *C.GValue) {
		defer C.g_object_unref(C.gpointer(interpolator))
	})
	C.g_value_set_object(&o.gvalue, C.gpointer(interpolator))
	return o
}
