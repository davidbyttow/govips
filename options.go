package govips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import "unsafe"

type OptionType int

const (
	OptionTypeBool OptionType = iota
	OptionTypeInt
	OptionTypeDouble
	OptionTypeString
	OptionTypeImage
	OptionTypeBlob
	OptionTypeInterpolator
)

type OptionTypeSerializer interface {
	Serialize(*C.GValue, interface{})
	Deserialize(interface{}, *C.GValue)
}

type BoolSerializer struct{}

func (t BoolSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boolean(dst, toGboolean(src.(bool)))
}

func (t BoolSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*bool) = fromGboolean(C.g_value_get_boolean(src))
}

type IntSerializer struct{}

func (t IntSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_int(dst, C.gint(src.(int)))
}

func (t IntSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*int) = int(C.g_value_get_int(src))
}

type DoubleSerializer struct{}

func (t DoubleSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_double(dst, C.gdouble(src.(float64)))
}

func (t DoubleSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*float64) = float64(C.g_value_get_double(src))
}

type StringSerializer struct{}

func (t StringSerializer) Serialize(dst *C.GValue, src interface{}) {
	c_s := C.CString(src.(string))
	defer freeCString(c_s)
	C.g_value_set_string(dst, (*C.gchar)(c_s))
}

func (t StringSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*string) = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(src))))
}

type ImageSerializer struct{}

func (t ImageSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, src.(*Image).image)
}

func (t ImageSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Image) = newImage((*C.VipsImage)(C.g_value_get_object(src)))
}

type BlobSerializer struct{}

func (t BlobSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boxed(dst, src.(*Blob).c_blob)
}

func (t BlobSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Blob) = newBlob((*C.VipsBlob)(C.g_value_dup_boxed(src)))
}

type InterpolatorSerializer struct{}

func (t InterpolatorSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, src.(*Interpolator).interp)
}

func (t InterpolatorSerializer) Deserialize(dst interface{}, src *C.GValue) {
	panic("Interpolator output not implemented")
}

var optionSerializers = map[OptionType]OptionTypeSerializer{
	OptionTypeBool:         BoolSerializer{},
	OptionTypeInt:          IntSerializer{},
	OptionTypeDouble:       DoubleSerializer{},
	OptionTypeString:       StringSerializer{},
	OptionTypeImage:        ImageSerializer{},
	OptionTypeBlob:         BlobSerializer{},
	OptionTypeInterpolator: InterpolatorSerializer{},
}

type Option struct {
	name       string
	value      interface{}
	optionType OptionType
	gvalue     C.GValue
	isOutput   bool
}

func newOption(name string, value interface{}, optionType OptionType, gType C.GType, isOutput bool) *Option {
	o := &Option{
		name:       name,
		value:      value,
		optionType: optionType,
		isOutput:   isOutput,
	}
	C.g_value_init(&o.gvalue, gType)
	return o
}

func newInput(name string, value interface{}, optionType OptionType, gType C.GType) *Option {
	o := newOption(name, value, optionType, gType, false)
	optionSerializers[o.optionType].Serialize(&o.gvalue, o.value)
	return o
}

func newOutput(name string, value interface{}, optionType OptionType, gType C.GType) *Option {
	return newOption(name, value, optionType, gType, true)
}

func (o *Option) Deserialize() {
	if !o.isOutput {
		panic("Option is not an output")
	}
	optionSerializers[o.optionType].Deserialize(o.value, &o.gvalue)
}

type Options struct {
	options []*Option
}

func NewOptions() *Options {
	return &Options{}
}

func (t *Options) deserializeOutputs() {
	for _, o := range t.options {
		if o.isOutput {
			o.Deserialize()
		}
	}
}

func (t *Options) addInput(name string, i interface{}, optionType OptionType, gType C.GType) *Option {
	o := newInput(name, i, optionType, gType)
	t.options = append(t.options, o)
	return o
}

func (t *Options) addOutput(name string, i interface{}, optionType OptionType, gType C.GType) *Option {
	o := newOutput(name, i, optionType, gType)
	t.options = append(t.options, o)
	return o
}

func (t *Options) SetBool(name string, b bool) *Options {
	t.addInput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	return t
}

func (t *Options) SetInt(name string, v int) *Options {
	t.addInput(name, v, OptionTypeInt, C.G_TYPE_INT)
	return t
}

func (t *Options) SetDouble(name string, v float64) *Options {
	t.addInput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	return t
}

func (t *Options) SetString(name string, s string) *Options {
	t.addInput(name, s, OptionTypeString, C.G_TYPE_STRING)
	return t
}

func (t *Options) SetImage(name string, image *Image) *Options {
	t.addInput(name, image, OptionTypeImage, C.vips_image_get_type())
	return t
}

func (t *Options) SetBlob(name string, blob *Blob) *Options {
	t.addInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	return t
}

func (t *Options) SetInterpolator(name string, interp *Interpolator) *Options {
	t.addInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
	return t
}

func (t *Options) SetBoolOut(name string, b *bool) *Options {
	t.addOutput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	return t
}

func (t *Options) SetIntOut(name string, v *int) *Options {
	t.addOutput(name, v, OptionTypeInt, C.G_TYPE_INT)
	return t
}

func (t *Options) SetDoubleOut(name string, v *float64) *Options {
	t.addOutput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	return t
}

func (t *Options) SetStringOut(name string, s *string) *Options {
	t.addOutput(name, s, OptionTypeString, C.G_TYPE_STRING)
	return t
}

func (t *Options) SetImageOut(name string, image **Image) *Options {
	t.addOutput(name, image, OptionTypeImage, C.vips_image_get_type())
	return t
}

func (t *Options) SetBlobOut(name string, blob **Blob) *Options {
	t.addOutput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	return t
}
