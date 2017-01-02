package gimage

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJpeg(t *testing.T) {
	buf, _ := ioutil.ReadFile("fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJpeg, imageType)
}
