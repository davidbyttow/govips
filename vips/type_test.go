package vips

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJPEG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(assets + "canyon.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJPEG, imageType)
}

func TestHEIF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(assets + "citron.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}
