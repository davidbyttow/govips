package gimage

import (
	"testing"

	"github.com/davidbyttow/gomore/io"
	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	buf, err := io.ReadFile("fixtures/canyon.jpg")
	require.NoError(t, err)

	image, err := LoadBuffer(buf)
	require.NoError(t, err)
	require.Equal(t, image.Type(), ImageTypeJpeg)
}
