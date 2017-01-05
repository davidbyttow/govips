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
	C.g_value_set_object(dst, src.(*Image).image)
}

func (t imageSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Image) = newImage((*C.VipsImage)(C.g_value_get_object(src)))
}

type blobSerializer struct{}

func (t blobSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boxed(dst, src.(*Blob).cBlob)
}

func (t blobSerializer) deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Blob) = newBlob((*C.VipsBlob)(C.g_value_dup_boxed(src)))
}

type interpolatorSerializer struct{}

func (t interpolatorSerializer) serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, src.(*Interpolator).interp)
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

// SetBool sets a boolean value for an optional parameter
func (t *Options) SetBool(name string, b bool) *Options {
	t.addInput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	return t
}

// SetInt sets a integer value for an optional parameter
func (t *Options) SetInt(name string, v int) *Options {
	t.addInput(name, v, OptionTypeInt, C.G_TYPE_INT)
	return t
}

// SetDouble sets a double value for an optional parameter
func (t *Options) SetDouble(name string, v float64) *Options {
	t.addInput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	return t
}

// SetString sets a string value for an optional parameter
func (t *Options) SetString(name string, s string) *Options {
	t.addInput(name, s, OptionTypeString, C.G_TYPE_STRING)
	return t
}

// SetImage sets a Image value for an optional parameter
func (t *Options) SetImage(name string, image *Image) *Options {
	t.addInput(name, image, OptionTypeImage, C.vips_image_get_type())
	return t
}

// SetBlob sets a Vlob value for an optional parameter
func (t *Options) SetBlob(name string, blob *Blob) *Options {
	t.addInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	return t
}

// SetInterpolator sets a Interpolator value for an optional parameter
func (t *Options) SetInterpolator(name string, interp *Interpolator) *Options {
	t.addInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
	return t
}

// SetBoolOut specifies a boolean output parameter for an operation
func (t *Options) SetBoolOut(name string, b *bool) *Options {
	t.addOutput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	return t
}

// SetIntOut specifies a integer output parameter for an operation
func (t *Options) SetIntOut(name string, v *int) *Options {
	t.addOutput(name, v, OptionTypeInt, C.G_TYPE_INT)
	return t
}

// SetDoubleOut specifies a boolean output parameter for an operation
func (t *Options) SetDoubleOut(name string, v *float64) *Options {
	t.addOutput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	return t
}

// SetStringOut specifies a string output parameter for an operation
func (t *Options) SetStringOut(name string, s *string) *Options {
	t.addOutput(name, s, OptionTypeString, C.G_TYPE_STRING)
	return t
}

// SetImageOut specifies a Image output parameter for an operation
func (t *Options) SetImageOut(name string, image **Image) *Options {
	t.addOutput(name, image, OptionTypeImage, C.vips_image_get_type())
	return t
}

// SetBlobOut specifies a Blob output parameter for an operation
func (t *Options) SetBlobOut(name string, blob **Blob) *Options {
	t.addOutput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	return t
}
