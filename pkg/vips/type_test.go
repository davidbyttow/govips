package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wix-playground/govips/pkg/vips"
)

func TestJPEG(t *testing.T) {
	buf, _ := ioutil.ReadFile("../../assets/fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := vips.DetermineImageType(buf)
	assert.Equal(t, vips.ImageTypeJPEG, imageType)
}
