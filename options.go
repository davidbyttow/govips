package govips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import "unsafe"

// OptionType represents the data type of an option
type OptionType int

// OptionType enum
const (
	OptionTypeBool OptionType = iota
	OptionTypeInt
	OptionTypeDouble
	OptionTypeString
	OptionTypeImage
	OptionTypeBlob
	OptionTypeInterpolator
)

type optionTypeSerializer interface {
	serialize(*C.GValue, interface{})
	deserialize(interface{}, *C.GValue)
}

type boolSerializer struct{}

func (t boolSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boolean(dst, toGboolean(src.(bool)))
}

func (t boolSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(*bool) = fromGboolean(C.g_value_get_boolean(src))
}

type intSerializer struct{}

func (t intSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_int(dst, C.gint(src.(int)))
}

func (t intSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(*int) = int(C.g_value_get_int(src))
}

type doubleSerializer struct{}

func (t doubleSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_double(dst, C.gdouble(src.(float64)))
}

func (t doubleSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(*float64) = float64(C.g_value_get_double(src))
}

type stringSerializer struct{}

func (t stringSerializer) serialize(dst *C.GValue, src interface{}) {
	cStr := C.CString(src.(string))
	defer freeCString(cStr)
	C.g_value_set_string(dst, (*C.gchar)(cStr))
}

func (t stringSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(*string) = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(src))))
}

type imageSerializer struct{}

func (t imageSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, C.gpointer(src.(*Image).image))
}

func (t imageSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Image) = newImage((*C.VipsImage)(C.g_value_get_object(src)))
}

type blobSerializer struct{}

func (t blobSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boxed(dst, C.gconstpointer(src.(*Blob).cBlob))
}

func (t blobSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Blob) = newBlob((*C.VipsBlob)(C.g_value_dup_boxed(src)))
}

type interpolatorSerializer struct{}

func (t interpolatorSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, C.gpointer(src.(*Interpolator).interp))
}

func (t interpolatorSerializer) deserialize(dst interface{}, src *C.GValue) {
	panic("Interpolator output not implemented")
}

var optionSerializers = map[OptionType]optionTypeSerializer{
	OptionTypeBool:         boolSerializer{},
	OptionTypeInt:          intSerializer{},
	OptionTypeDouble:       doubleSerializer{},
	OptionTypeString:       stringSerializer{},
	OptionTypeImage:        imageSerializer{},
	OptionTypeBlob:         blobSerializer{},
	OptionTypeInterpolator: interpolatorSerializer{},
}

type option struct {
	name       string
	value      interface{}
	optionType OptionType
	gvalue     C.GValue
	isOutput   bool
}

func newOption(name string, value interface{}, optionType OptionType, gType C.GType, isOutput bool) *option {
	o := &option{
		name:       name,
		value:      value,
		optionType: optionType,
		isOutput:   isOutput,
	}
	C.g_value_init(&o.gvalue, gType)
	return o
}

func newInput(name string, value interface{}, optionType OptionType, gType C.GType) *option {
	o := newOption(name, value, optionType, gType, false)
	optionSerializers[o.optionType].serialize(&o.gvalue, o.value)
	return o
}

func newOutput(name string, value interface{}, optionType OptionType, gType C.GType) *option {
	return newOption(name, value, optionType, gType, true)
}

func (o *option) Deserialize() {
	if !o.isOutput {
		panic("Option is not an output")
	}
	optionSerializers[o.optionType].deserialize(o.value, &o.gvalue)
}

// Options specifies optional parameters for an operation
type Options struct {
	options []*option
}

// NewOptions returns a new option set
func NewOptions(options ...OptionFunc) *Options {
	return (&Options{}).With(options...)
}

func (t *Options) deserializeOutputs() {
	for _, o := range t.options {
		if o.isOutput {
			o.Deserialize()
		}
	}
}

func (t *Options) addInput(name string, i interface{}, optionType OptionType, gType C.GType) *option {
	o := newInput(name, i, optionType, gType)
	t.options = append(t.options, o)
	return o
}

func (t *Options) addOutput(name string, i interface{}, optionType OptionType, gType C.GType) *option {
	o := newOutput(name, i, optionType, gType)
	t.options = append(t.options, o)
	return o
}

type OptionFunc func(t *Options)

func (t *Options) With(options ...OptionFunc) *Options {
	for _, o := range options {
		o(t)
	}
	return t
}

// BoolInput sets a boolean value for an optional parameter
func BoolInput(name string, b bool) OptionFunc {
	return func(t *Options) {
		t.addInput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	}
}

// SetInt sets a integer value for an optional parameter
func IntInput(name string, v int) OptionFunc {
	return func(t *Options) {
		t.addInput(name, v, OptionTypeInt, C.G_TYPE_INT)
	}
}

// SetDouble sets a double value for an optional parameter
func DoubleInput(name string, v float64) OptionFunc {
	return func(t *Options) {
		t.addInput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	}
}

// StringInput sets a string value for an optional parameter
func StringInput(name string, s string) OptionFunc {
	return func(t *Options) {
		t.addInput(name, s, OptionTypeString, C.G_TYPE_STRING)
	}
}

// ImageInput sets a Image value for an optional parameter
func ImageInput(name string, image *Image) OptionFunc {
	return func(t *Options) {
		t.addInput(name, image, OptionTypeImage, C.vips_image_get_type())
	}
}

// BlobInput sets a Blob value for an optional parameter
func BlobInput(name string, blob *Blob) OptionFunc {
	return func(t *Options) {
		t.addInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	}
}

// InterpolatorInput sets a Interpolator value for an optional parameter
func InterpolatorInput(name string, interp *Interpolator) OptionFunc {
	return func(t *Options) {
		t.addInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
	}
}

// BoolOutput specifies a boolean output parameter for an operation
func BoolOutput(name string, b *bool) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	}
}

// IntOutput specifies a integer output parameter for an operation
func IntOutput(name string, v *int) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, v, OptionTypeInt, C.G_TYPE_INT)
	}
}

// DoubleOutput specifies a boolean output parameter for an operation
func DoubleOutput(name string, v *float64) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	}
}

// StringOutput specifies a string output parameter for an operation
func StringOutput(name string, s *string) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, s, OptionTypeString, C.G_TYPE_STRING)
	}
}

// ImageOutput specifies a Image output parameter for an operation
func ImageOutput(name string, image **Image) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, image, OptionTypeImage, C.vips_image_get_type())
	}
}

// BlobOutput specifies a Blob output parameter for an operation
func BlobOutput(name string, blob **Blob) OptionFunc {
	return func(t *Options) {
		t.addOutput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	}
}
