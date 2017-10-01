package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"unsafe"
)

type VipsSerializer func(*C.GValue)

type VipsOption interface {
	Name() string
	Close()
}

type VipsInput interface {
	Serialize() *C.GValue
}

type VipsOutput interface {
	GValue() *C.GValue
	Deserialize()
}

type vipsOption struct {
	name   string
	gvalue C.GValue
	fn     VipsSerializer
	closer func()
}

func (v *vipsOption) Name() string {
	return v.name
}

func (v *vipsOption) Close() {
	if v.closer != nil {
		v.closer()
	}
	C.g_value_unset(&v.gvalue)
}

type vipsInput struct {
	*vipsOption
}

func (v *vipsInput) Serialize() *C.GValue {
	gv := &v.gvalue
	v.fn(gv)
	return gv
}

func (v *vipsInput) GValue() *C.GValue {
	return &v.gvalue
}

type vipsOutput struct {
	*vipsOption
}

func (v *vipsOutput) Deserialize() {
	gv := &v.gvalue
	v.fn(gv)
}

func NewVipsInput(name string, gtype C.GType, fn VipsSerializer) *vipsInput {
	v := &vipsInput{&vipsOption{
		name: name,
		fn:   fn,
	}}
	C.g_value_init(&v.gvalue, gtype)
	return v
}

func NewVipsOutput(name string, gtype C.GType, fn VipsSerializer) *vipsOutput {
	v := &vipsOutput{&vipsOption{
		name: name,
		fn:   fn,
	}}
	C.g_value_init(&v.gvalue, gtype)
	return v
}

func InputBool(name string, v bool) VipsOption {
	i := NewVipsInput(name, C.G_TYPE_BOOLEAN, func(gv *C.GValue) {
		C.g_value_set_boolean(gv, toGboolean(v))
	})
	return i
}

func OutputBool(name string, v *bool) VipsOption {
	o := NewVipsOutput(name, C.G_TYPE_BOOLEAN, func(gv *C.GValue) {
		*v = fromGboolean(C.g_value_get_boolean(gv))
	})
	return o
}

func InputInt(name string, v int) VipsOption {
	return NewVipsInput(name, C.G_TYPE_INT, func(gv *C.GValue) {
		C.g_value_set_int(gv, C.gint(v))
	})
}

func OutputInt(name string, v *int) VipsOption {
	return NewVipsOutput(name, C.G_TYPE_INT, func(gv *C.GValue) {
		*v = int(C.g_value_get_int(gv))
	})
}

func InputDouble(name string, v float64) VipsOption {
	return NewVipsInput(name, C.G_TYPE_DOUBLE, func(gv *C.GValue) {
		C.g_value_set_double(gv, C.gdouble(v))
	})
}

func OutputDouble(name string, v *float64) VipsOption {
	return NewVipsOutput(name, C.G_TYPE_DOUBLE, func(gv *C.GValue) {
		*v = float64(C.g_value_get_double(gv))
	})
}

func InputString(name string, v string) VipsOption {
	return NewVipsInput(name, C.G_TYPE_STRING, func(gv *C.GValue) {
		cStr := C.CString(v)
		defer freeCString(cStr)
		C.g_value_set_string(gv, (*C.gchar)(cStr))
	})
}

func OutputString(name string, v *string) VipsOption {
	return NewVipsOutput(name, C.G_TYPE_STRING, func(gv *C.GValue) {
		*v = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(gv))))
	})
}

func InputImage(name string, v *C.VipsImage) VipsOption {
	return NewVipsInput(name, C.vips_image_get_type(), func(gv *C.GValue) {
		C.g_value_set_object(gv, C.gpointer(v))
	})
}

func OutputImage(name string, v **C.VipsImage) VipsOption {
	return NewVipsOutput(name, C.vips_image_get_type(), func(gv *C.GValue) {
		*v = (*C.VipsImage)(C.g_value_get_object(gv))
	})
}
