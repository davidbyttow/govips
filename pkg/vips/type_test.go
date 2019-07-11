package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wix-playground/govips/pkg/vips"
)

func TestJPEG(t *testing.T) {
	vips.Startup(&vips.Config{})

	buf, _ := ioutil.ReadFile("../../assets/fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := vips.DetermineImageType(buf)
	assert.Equal(t, vips.ImageTypeJPEG, imageType)
}

func TestHEIF(t *testing.T) {
	vips.Startup(&vips.Config{})

	buf, _ := ioutil.ReadFile("../../assets/fixtures/citron.heic")
	assert.NotNil(t, buf)

	imageType := vips.DetermineImageType(buf)
	assert.Equal(t, vips.ImageTypeHEIF, imageType)
}
