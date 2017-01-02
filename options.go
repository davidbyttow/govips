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

type ValueDeserializer func(*Option)

var valueDeserializers = map[OptionType]ValueDeserializer{
	OptionTypeBool: func(option *Option) {
		b := option.value.(*bool)
		*b = fromGboolean(C.g_value_get_boolean(&option.gvalue))
	},
	OptionTypeInt: func(option *Option) {
		v := option.value.(*int)
		*v = int(C.g_value_get_int(&option.gvalue))
	},
	OptionTypeDouble: func(option *Option) {
		v := option.value.(*float64)
		*v = float64(C.g_value_get_double(&option.gvalue))
	},
	OptionTypeString: func(option *Option) {
		s := option.value.(*string)
		*s = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(&option.gvalue))))
	},
	OptionTypeImage: func(option *Option) {
		image := option.value.(**Image)
		*image = newImage((*C.VipsImage)(C.g_value_get_object(&option.gvalue)))
	},
	OptionTypeBlob: func(option *Option) {
		blob := option.value.(**Blob)
		debug("%#v", option.gvalue)
		*blob = newBlob((*C.VipsBlob)(C.g_value_dup_boxed(&option.gvalue)))
	},
}

func deserialize(option *Option) {
	fn := valueDeserializers[option.optionType]
	fn(option)
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
	return newOption(name, value, optionType, gType, false)
}

func newOutput(name string, value interface{}, optionType OptionType, gType C.GType) *Option {
	return newOption(name, value, optionType, gType, true)
}

func (o *Option) Deserialize() {
	if !o.isOutput {
		panic("Option is not an output")
	}
	deserialize(o)
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
	o := t.addInput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	C.g_value_set_boolean(&o.gvalue, toGboolean(b))
	return t
}

func (t *Options) SetInt(name string, v int) *Options {
	o := t.addInput(name, v, OptionTypeInt, C.G_TYPE_INT)
	C.g_value_set_int(&o.gvalue, C.gint(v))
	return t
}

func (t *Options) SetDouble(name string, v float64) *Options {
	o := t.addInput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	C.g_value_set_double(&o.gvalue, C.gdouble(v))
	return t
}

func (t *Options) SetString(name string, s string) *Options {
	o := t.addInput(name, s, OptionTypeString, C.G_TYPE_STRING)
	c_s := C.CString(s)
	defer freeCString(c_s)
	C.g_value_set_string(&o.gvalue, (*C.gchar)(c_s))
	return t
}

func (t *Options) SetImage(name string, image *Image) *Options {
	o := t.addInput(name, image, OptionTypeImage, C.vips_image_get_type())
	C.g_value_set_object(&o.gvalue, image.image)
	return t
}

func (t *Options) SetBlob(name string, blob *Blob) *Options {
	o := t.addInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	C.g_value_set_boxed(&o.gvalue, blob.c_blob)
	return t
}

func (t *Options) SetInterpolator(name string, interp *Interpolator) *Options {
	o := t.addInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
	C.g_value_set_object(&o.gvalue, interp.interp)
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
