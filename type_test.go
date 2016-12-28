package gimage

import (
	"testing"

	"github.com/davidbyttow/gomore/io"
	"github.com/stretchr/testify/assert"
)

func TestJpeg(t *testing.T) {
	buf, _ := io.ReadFile("fixtures/canyon.jpg")
	assert.NotNil(t, buf)

	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJpeg, imageType)
}
