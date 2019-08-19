package vips

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DetermineImageType__JPEG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "canyon.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJPEG, imageType)
}

func Test_DetermineImageType__HEIF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "citron.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func Test_DetermineImageType__PNG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "clover.png")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePNG, imageType)
}

func Test_DetermineImageType__TIFF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "galaxy.tif")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeTIFF, imageType)
}

func Test_DetermineImageType__WEBP(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "dice.webp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeWEBP, imageType)
}

func Test_DetermineImageType__SVG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "gopher.svg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeSVG, imageType)
}

func Test_DetermineImageType__PDF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "42.pdf")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePDF, imageType)
}

func Test_DetermineImageType__BMP(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "teddy.bmp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeBMP, imageType)
}
