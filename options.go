package gimage

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"

type OptionType int

const (
	OptionTypeBool OptionType = iota
	OptionTypeInt
	OptionTypeDouble
	OptionTypeString
	OptionTypeImage
)

type ValueDeserializer func(*Option)

var valueDeserializers = map[OptionType]ValueDeserializer{
	OptionTypeImage: func(option *Option) {
		image := option.value.(*Image)
		image.image = (*C.VipsImage)(C.g_value_get_object(&option.gvalue))
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

func newOption(name string, value interface{}, optionType OptionType, isOutput bool) *Option {
	return &Option{
		name:       name,
		value:      value,
		optionType: optionType,
		isOutput:   isOutput,
	}
}

func newInput(name string, value interface{}, optionType OptionType) *Option {
	return newOption(name, value, optionType, false)
}

func newOutput(name string, value interface{}, optionType OptionType) *Option {
	return newOption(name, value, optionType, true)
}

func (o *Option) Deserialize() {
	deserialize(o)
}

type Options struct {
	options []*Option
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetBool(name string, b bool) *Options {
	option := newInput(name, b, OptionTypeBool)
	C.g_value_init(&option.gvalue, C.G_TYPE_BOOLEAN)
	C.g_value_set_boolean(&option.gvalue, toGboolean(b))
	o.options = append(o.options, option)
	return o
}

func (o *Options) SetInt(name string, v int) *Options {
	option := newInput(name, v, OptionTypeInt)
	C.g_value_init(&option.gvalue, C.G_TYPE_INT)
	C.g_value_set_int(&option.gvalue, C.gint(v))
	o.options = append(o.options, option)
	return o
}

func (o *Options) SetDouble(name string, v float64) *Options {
	option := newInput(name, v, OptionTypeDouble)
	C.g_value_init(&option.gvalue, C.G_TYPE_DOUBLE)
	C.g_value_set_double(&option.gvalue, C.gdouble(v))
	o.options = append(o.options, option)
	return o
}

func (o *Options) SetImage(name string, image Image) *Options {
	option := newInput(name, image, OptionTypeImage)
	C.g_value_init(&option.gvalue, C.vips_image_get_type())
	C.g_value_set_object(&option.gvalue, image.image)
	o.options = append(o.options, option)
	return o
}

func (o *Options) SetImageOut(name string, image *Image) *Options {
	option := newOutput(name, image, OptionTypeImage)
	C.g_value_init(&option.gvalue, C.vips_image_get_type())
	o.options = append(o.options, option)
	return o
}
