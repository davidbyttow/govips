package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/jhford/govips/pkg/vips"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
  tests := map[string]vips.ImageType{
    "koala.jpg": vips.ImageTypeJPEG,
    "koala.webp": vips.ImageTypeWEBP,
    "koala.png": vips.ImageTypePNG,
    "koala.tiff": vips.ImageTypeTIFF,
    "koala.gif": vips.ImageTypeGIF,
    //"koala.bmp": vips.ImageTypeBMP,
    //"koala.bmp2": vips.ImageTypeBMP,
    //"koala.bmp3": vips.ImageTypeBMP,
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
