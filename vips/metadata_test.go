package vips

import (
	"testing"

	"github.com/davidbyttow/gimage"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	buf, err := gimage.ReadFile("../fixtures/canyon.jpg")
	require.NoError(t, err)

	image, err := LoadBuffer(buf)
	require.NoError(t, err)
	require.Equal(t, image.Type(), ImageTypeJpeg)

	// metadata, err := LoadMetadata(image)
	// require.NoError(t, err)
	// assert.Equal(t, 2560, metadata.Size.Width)
	// assert.Equal(t, 1600, metadata.Size.Height)

	// image, err = image.Shrink(0.5, 0.5)
	// require.NoError(t, err)

	// metadata, err = LoadMetadata(image)
	// require.NoError(t, err)
	// assert.Equal(t, 1280, metadata.Size.Width)
	// assert.Equal(t, 800, metadata.Size.Height)
}
