package gimage

import (
	"testing"

	"github.com/davidbyttow/gomore/io"
	"github.com/stretchr/testify/assert"
)

func TestOperation(t *testing.T) {
	buf, _ := io.ReadFile("fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	image, _ := NewImageFromBuffer(buf, nil)
	assert.NotNil(t, image)
	assert.Equal(t, 2560, image.Width())
	assert.Equal(t, 1600, image.Height())

	out := image.Shrink(2.0, 2.0, nil)
	assert.Equal(t, 1280, out.Width())
	assert.Equal(t, 800, out.Height())
}
