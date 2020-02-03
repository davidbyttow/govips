package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/jhford/govips/pkg/vips"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	tests := map[string]vips.ImageType{
		"koala.jpg":  vips.ImageTypeJPEG,
		"koala.webp": vips.ImageTypeWEBP,
		"koala.png":  vips.ImageTypePNG,
		"koala.tiff": vips.ImageTypeTIFF,
		"koala.gif":  vips.ImageTypeGIF,
		"koala.bmp":  vips.ImageTypeBMP,
		"koala.bmp2": vips.ImageTypeBMP,
		"koala.bmp3": vips.ImageTypeBMP,
	}

	for input, expected := range tests {
		input := input
		expected := expected

		buf, _ := ioutil.ReadFile("../../assets/fixtures/" + input)
		assert.NotNil(t, buf)

		imageType := vips.DetermineImageType(buf)
		assert.Equal(t, expected, imageType,
			"%s <> %s", expected.OutputExt(), imageType.OutputExt())
	}
}

func TestLoading(t *testing.T) {
	tests := map[string]vips.ImageType{
		"koala.jpg":  vips.ImageTypeJPEG,
		"koala.webp": vips.ImageTypeWEBP,
		"koala.png":  vips.ImageTypePNG,
		"koala.tiff": vips.ImageTypeTIFF,
		"koala.gif":  vips.ImageTypeGIF,
		"koala.bmp":  vips.ImageTypeBMP,
		"koala.bmp2": vips.ImageTypeBMP,
		"koala.bmp3": vips.ImageTypeBMP,
	}

	for input, expected := range tests {
		input := input
		expected := expected

		image, err := vips.NewImageFromFile("../../assets/fixtures/" + input)
		assert.NoError(t, err, "loading %s", expected.OutputExt())

		assert.Equal(t, expected, image.Format())
	}
}

func TestSaving(t *testing.T) {
	tests := []vips.ImageType{
		vips.ImageTypeJPEG,
		vips.ImageTypeWEBP,
		vips.ImageTypePNG,
		vips.ImageTypeTIFF,
		vips.ImageTypeGIF,
		vips.ImageTypeBMP,
	}

	image, err := vips.NewImageFromFile("../../assets/fixtures/" + "koala.png")
	assert.NoError(t, err)
	defer func() {
		image.Close()
	}()

	for _, expected := range tests {
		expected := expected

		tx := vips.NewTransform()
		tx.Image(image)
		tx.Format(expected)
		tx.Resize(100, 100)

		buf, format, err := tx.Apply()
		assert.NoError(t, err)
		assert.Greater(t, len(buf), 0)
		assert.Equal(t, expected, format, "%s != %s", vips.ImageTypes[expected], vips.ImageTypes[format])
		actualFormat := vips.DetermineImageType(buf)
		assert.Equal(t, expected, actualFormat, "%s != %s", vips.ImageTypes[expected], vips.ImageTypes[actualFormat])
	}
}
