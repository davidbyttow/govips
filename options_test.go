package vips_test

import (
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
)

func TestOptionPrimitives(t *testing.T) {
	var b bool
	var i int
	var d float64
	var s string
	options := vips.NewOptions(
		vips.BoolInput("b", true),
		vips.IntInput("i", 42),
		vips.DoubleInput("d", 42.2),
		vips.StringInput("s", "hi"),
		vips.BoolOutput("b", &b),
		vips.IntOutput("i", &i),
		vips.DoubleOutput("d", &d),
		vips.StringOutput("s", &s),
	)

	assert.Equal(t, true, options.Options[0].Value.(bool))
	assert.Equal(t, 42, options.Options[1].Value.(int))
	assert.Equal(t, 42.2, options.Options[2].Value.(float64))
	assert.Equal(t, "hi", options.Options[3].Value.(string))

	options.Options[4].GValue = options.Options[0].GValue
	options.Options[5].GValue = options.Options[1].GValue
	options.Options[6].GValue = options.Options[2].GValue
	options.Options[7].GValue = options.Options[3].GValue

	options.DeserializeOutputs()

	assert.Equal(t, true, b)
	assert.Equal(t, 42, i)
	assert.Equal(t, 42.2, d)
	assert.Equal(t, "hi", s)
}

func TestOptionBlob(t *testing.T) {
	bytes := []byte{42, 43, 44}
	in := vips.NewBlob(bytes)
	var out *vips.Blob
	options := vips.NewOptions(
		vips.BlobInput("in", in),
		vips.BlobOutput("out", &out),
	)

	assert.Equal(t, in, (options.Options[0].Value.(*vips.Blob)))

	options.Options[1].GValue = options.Options[0].GValue

	options.DeserializeOutputs()

	assert.Equal(t, len(bytes), out.Length())
	assert.Equal(t, bytes, out.ToBytes())
}
