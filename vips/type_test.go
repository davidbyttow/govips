package vips

import (
	"testing"

	"github.com/davidbyttow/gomore/io"
	"github.com/stretchr/testify/assert"
)

func TestJpeg(t *testing.T) {
	buf, err := io.ReadFile("../fixtures/canyon.jpg")
	if err != nil {
		t.Fail()
	}
	imageType := DetermineImageType(buf)
	assert.Equal(t, ImageTypeJpeg, imageType)
}
