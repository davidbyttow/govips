package vips_test

import (
	"io/ioutil"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
)

func TestJPEG(t *testing.T) {
	buf, _ := ioutil.ReadFile("fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := vips.DetermineImageType(buf)
	assert.Equal(t, vips.ImageTypeJPEG, imageType)
}
