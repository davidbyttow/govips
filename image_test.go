package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteToBytes(t *testing.T) {
	buf, err := ioutil.ReadFile("fixtures/canyon.jpg")
	require.Nil(t, err)

	image, err := vips.NewImageFromBuffer(buf)
	require.Nil(t, err)

	image = image.Resize(0.25)

	buf, err = image.Export(vips.ExportOptions{})
	require.Nil(t, err)
	assert.True(t, len(buf) > 0)

	image, err = vips.NewImageFromBuffer(buf)
	require.Nil(t, err)
	assert.Equal(t, 640, image.Width())
	assert.Equal(t, 400, image.Height())
}

func TestLoadFromMemory(t *testing.T) {
	size := 200

	bytes := make([]byte, size*size*3)
	for i := 0; i < size*size; i++ {
		bytes[i*3] = 0xFF
		bytes[i*3+1] = 0
		bytes[i*3+2] = 0
	}

	_, err := vips.NewImageFromMemory(bytes, size, size, 3, vips.BandFormatUchar)
	require.Nil(t, err)
}
