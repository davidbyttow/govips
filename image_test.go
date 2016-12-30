package gimage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO(d): Fix this, as it's resulting in an odd error when image is loaded from file twice.
// Likely has to do with operation caching and not releasing.
func TestLoadFromFile(t *testing.T) {
	// image, err := NewImageFromFile("fixtures/canyon.jpg", nil)
	// assert.Nil(t, err)
	// assert.Equal(t, 2560, image.Width())
	// assert.Equal(t, 1600, image.Height())
}

func TestWriteToFile(t *testing.T) {
	image, err := NewImageFromFile("fixtures/canyon.jpg", nil)
	assert.Nil(t, err)

	image = image.Resize(0.25, nil)

	err = image.WriteToFile("./fixtures/canyon-out.jpg", nil)
	assert.Nil(t, err)
}
