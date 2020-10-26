package vips

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func Test_DetermineImageType__JPEG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "jpg-24bit-icc-iec.jpg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJPEG, imageType)
}

func Test_DetermineImageType__HEIF_HEIC(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "heic-24bit-exif.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func Test_DetermineImageType__HEIF_MIF1(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "heic-24bit.heic")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func Test_DetermineImageType__PNG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "png-24bit+alpha.png")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePNG, imageType)
}

func Test_DetermineImageType__TIFF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "tif.tif")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeTIFF, imageType)
}

func Test_DetermineImageType__WEBP(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "webp+alpha.webp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeWEBP, imageType)
}

func Test_DetermineImageType__SVG(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "svg.svg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeSVG, imageType)
}

func Test_DetermineImageType__SVG_1(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "svg_1.svg")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeSVG, imageType)
}

func Test_DetermineImageType__PDF(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "pdf.pdf")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypePDF, imageType)
}

func Test_DetermineImageType__BMP(t *testing.T) {
	Startup(&Config{})

	buf, err := ioutil.ReadFile(resources + "bmp.bmp")
	assert.NoError(t, err)
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeBMP, imageType)
}
