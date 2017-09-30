package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/spf13/cast"
)

// OptionType represents the data type of an option
type OptionType string

// OptionType enum
const (
	OptionTypeBool         OptionType = "bool"
	OptionTypeInt          OptionType = "int"
	OptionTypeDouble       OptionType = "double"
	OptionTypeString       OptionType = "string"
	OptionTypeVipsImage    OptionType = "VipsImage"
	OptionTypeBlob         OptionType = "Blob"
	OptionTypeInterpolator OptionType = "Interpolator"
)

var optionSerializers = map[OptionType]OptionTypeSerializer{
	OptionTypeBool:         boolSerializer{},
	OptionTypeInt:          intSerializer{},
	OptionTypeDouble:       doubleSerializer{},
	OptionTypeString:       stringSerializer{},
	OptionTypeVipsImage:    vipsImageSerializer{},
	OptionTypeBlob:         blobSerializer{},
	OptionTypeInterpolator: interpolatorSerializer{},
}

type OptionTypeSerializer interface {
	Serialize(*C.GValue, interface{})
	Deserialize(interface{}, *C.GValue)
	String(interface{}) string
}

type boolSerializer struct{}

func (t boolSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boolean(dst, toGboolean(src.(bool)))
}

func (t boolSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*bool) = fromGboolean(C.g_value_get_boolean(src))
}

func (t boolSerializer) String(i interface{}) string {
	return cast.ToString(i)
}

type intSerializer struct{}

func (t intSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_int(dst, C.gint(src.(int)))
}

func (t intSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*int) = int(C.g_value_get_int(src))
}

func (t intSerializer) String(i interface{}) string {
	return cast.ToString(i)
}

type doubleSerializer struct{}

func (t doubleSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_double(dst, C.gdouble(src.(float64)))
}

func (t doubleSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*float64) = float64(C.g_value_get_double(src))
}

func (t doubleSerializer) String(i interface{}) string {
	return cast.ToString(i)
}

type stringSerializer struct{}

func (t stringSerializer) Serialize(dst *C.GValue, src interface{}) {
	cStr := C.CString(src.(string))
	defer freeCString(cStr)
	C.g_value_set_string(dst, (*C.gchar)(cStr))
}

func (t stringSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(*string) = C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(src))))
}

func (t stringSerializer) String(i interface{}) string {
	return cast.ToString(i)
}

type vipsImageSerializer struct{}

func (t vipsImageSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, C.gpointer(src.(*C.VipsImage)))
}

func (t vipsImageSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(**C.VipsImage) = (*C.VipsImage)(C.g_value_get_object(src))
}

func (t vipsImageSerializer) String(i interface{}) string {
	image := i.(*C.VipsImage)
	return cast.ToString(image)
}

type blobSerializer struct{}

func (t blobSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_boxed(dst, C.gconstpointer(src.(*Blob).cBlob))
}

func (t blobSerializer) Deserialize(dst interface{}, src *C.GValue) {
	*dst.(**Blob) = newBlob((*C.VipsBlob)(C.g_value_dup_boxed(src)))
}

func (t blobSerializer) String(i interface{}) string {
	blob := i.(*Blob)
	if blob == nil {
		return "nil"
	}
	return fmt.Sprintf("(%v,%d)", blob, blob.Length())
}

type interpolatorSerializer struct{}

func (t interpolatorSerializer) Serialize(dst *C.GValue, src interface{}) {
	C.g_value_set_object(dst, C.gpointer(src.(*Interpolator).interp))
}

func (t interpolatorSerializer) Deserialize(dst interface{}, src *C.GValue) {
	panic("Interpolator output not implemented")
}

func (t interpolatorSerializer) String(i interface{}) string {
	interp := i.(*Interpolator)
	return fmt.Sprintf("%v", interp)
}

type Option struct {
	Name       string
	Value      interface{}
	OptionType OptionType
	GValue     C.GValue
	IsOutput   bool
}

func newOption(name string, value interface{}, optionType OptionType, gType C.GType, isOutput bool) *Option {
	o := &Option{
		Name:       name,
		Value:      value,
		OptionType: optionType,
		IsOutput:   isOutput,
	}
	C.g_value_init(&o.GValue, gType)
	return o
}

func newInput(name string, value interface{}, optionType OptionType, gType C.GType) *Option {
	o := newOption(name, value, optionType, gType, false)
	optionSerializers[o.OptionType].Serialize(&o.GValue, o.Value)
	return o
}

func newOutput(name string, value interface{}, optionType OptionType, gType C.GType) *Option {
	return newOption(name, value, optionType, gType, true)
}

func (o *Option) String() string {
	if o.IsOutput {
		return fmt.Sprintf("%s* %s", o.OptionType, o.Name)
	}
	value := optionSerializers[o.OptionType].String(o.Value)
	return fmt.Sprintf("%s %s=%s", o.OptionType, o.Name, value)
}

func (o *Option) Deserialize() {
	if !o.IsOutput {
		panic("Option is not an output")
	}
	optionSerializers[o.OptionType].Deserialize(o.Value, &o.GValue)
}

// Options specifies optional parameters for an operation
type Options struct {
	Options []*Option
}

// NewOptions returns a new option set
func NewOptions(options ...OptionFunc) *Options {
	return (&Options{}).With(options...)
}

func (t *Options) DeserializeOutputs() {
	for _, o := range t.Options {
		if o.IsOutput {
			o.Deserialize()
		}
	}
}

func (t *Options) AddInput(name string, i interface{}, optionType OptionType, gType C.GType) *Option {
	o := newInput(name, i, optionType, gType)
	t.Options = append(t.Options, o)
	return o
}

func (t *Options) AddOutput(name string, i interface{}, optionType OptionType, gType C.GType) *Option {
	o := newOutput(name, i, optionType, gType)
	t.Options = append(t.Options, o)
	return o
}

func (t *Options) Release() {
	for _, o := range t.Options {
		C.g_value_unset(&o.GValue)
	}
}

// OptionFunc is a typeref that applies an option
type OptionFunc func(t *Options)

// With applies the given options
func (t *Options) With(options ...OptionFunc) *Options {
	for _, o := range options {
		o(t)
	}
	return t
}

// BoolInput sets a boolean value for an optional parameter
func BoolInput(name string, b bool) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	}
}

// IntInput sets a integer value for an optional parameter
func IntInput(name string, v int) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, v, OptionTypeInt, C.G_TYPE_INT)
	}
}

// DoubleInput sets a double value for an optional parameter
func DoubleInput(name string, v float64) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	}
}

// StringInput sets a string value for an optional parameter
func StringInput(name string, s string) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, s, OptionTypeString, C.G_TYPE_STRING)
	}
}

// VipsImageInput sets a Image value for a parameter
func VipsImageInput(name string, image *C.VipsImage) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, image, OptionTypeVipsImage, C.vips_image_get_type())
	}
}

// BlobInput sets a Blob value for an optional parameter
func BlobInput(name string, blob *Blob) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	}
}

// InterpolatorInput sets a Interpolator value for an optional parameter
func InterpolatorInput(name string, interp *Interpolator) OptionFunc {
	return func(t *Options) {
		t.AddInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
	}
}

// BoolOutput specifies a boolean output parameter for an operation
func BoolOutput(name string, b *bool) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, b, OptionTypeBool, C.G_TYPE_BOOLEAN)
	}
}

// IntOutput specifies a integer output parameter for an operation
func IntOutput(name string, v *int) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, v, OptionTypeInt, C.G_TYPE_INT)
	}
}

// DoubleOutput specifies a boolean output parameter for an operation
func DoubleOutput(name string, v *float64) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, v, OptionTypeDouble, C.G_TYPE_DOUBLE)
	}
}

// StringOutput specifies a string output parameter for an operation
func StringOutput(name string, s *string) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, s, OptionTypeString, C.G_TYPE_STRING)
	}
}

// VipsImageOutput specifies a Image output parameter for an operation
func VipsImageOutput(name string, image **C.VipsImage) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, image, OptionTypeVipsImage, C.vips_image_get_type())
	}
}

// BlobOutput specifies a Blob output parameter for an operation
func BlobOutput(name string, blob **Blob) OptionFunc {
	return func(t *Options) {
		t.AddOutput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
	}
}
