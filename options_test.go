package govips

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionPrimitives(t *testing.T) {
	var b bool
	var i int
	var d float64
	var s string
	options := NewOptions().
		SetBool("b", true).
		SetInt("i", 42).
		SetDouble("d", 42.2).
		SetString("s", "hi").
		SetBoolOut("b", &b).
		SetIntOut("i", &i).
		SetDoubleOut("d", &d).
		SetStringOut("s", &s)

	assert.Equal(t, true, options.options[0].value.(bool))
	assert.Equal(t, 42, options.options[1].value.(int))
	assert.Equal(t, 42.2, options.options[2].value.(float64))
	assert.Equal(t, "hi", options.options[3].value.(string))

	options.options[4].gvalue = options.options[0].gvalue
	options.options[5].gvalue = options.options[1].gvalue
	options.options[6].gvalue = options.options[2].gvalue
	options.options[7].gvalue = options.options[3].gvalue

	options.deserializeOutputs()

	assert.Equal(t, true, b)
	assert.Equal(t, 42, i)
	assert.Equal(t, 42.2, d)
	assert.Equal(t, "hi", s)
}

func TestOptionBlob(t *testing.T) {
	bytes := []byte{42, 43, 44}
	in := NewBlob(bytes)
	var out *Blob
	options := NewOptions().
		SetBlob("in", in).
		SetBlobOut("out", &out)

	assert.Equal(t, in, (options.options[0].value.(*Blob)))

	options.options[1].gvalue = options.options[0].gvalue

	options.deserializeOutputs()

	assert.Equal(t, len(bytes), out.Length())
	assert.Equal(t, bytes, out.ToBytes())
}

// func (t *Options) SetImage(name string, image *Image) *Options {
// 	o := t.addInput(name, image, OptionTypeImage, C.vips_image_get_type())
// 	C.g_value_set_object(&o.gvalue, image.image)
// 	return t
// }

// func (t *Options) SetBlob(name string, blob *Blob) *Options {
// 	o := t.addInput(name, blob, OptionTypeBlob, C.vips_blob_get_type())
// 	C.g_value_set_boxed(&o.gvalue, blob.c_blob)
// 	return t
// }

// func (t *Options) SetInterpolator(name string, interp *Interpolator) *Options {
// 	o := t.addInput(name, interp, OptionTypeInterpolator, C.vips_interpolate_get_type())
// 	C.g_value_set_object(&o.gvalue, interp.interp)
// 	return t
// }
