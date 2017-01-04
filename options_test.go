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
