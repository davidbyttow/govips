package vips

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const assets = "../assets/fixtures/"

func TestLoadImage_AccessMode_Default(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(assets + "test.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := LoadImage(src)
	if assert.NoError(t, err) {
		assert.NotNil(t, img)
		// check random access by encoding twice
		_, _, err = img.Export(ExportParams{})
		assert.NoError(t, err)
		_, _, err = img.Export(ExportParams{})
		assert.NoError(t, err)

	}
}

func TestLoadImage_AccessMode_Random(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(assets + "test.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := LoadImage(src, WithAccessMode(AccessRandom))
	if assert.NoError(t, err) {
		assert.NotNil(t, img)
		// check random access by encoding twice
		_, _, err = img.Export(ExportParams{})
		assert.NoError(t, err)
		_, _, err = img.Export(ExportParams{})
		assert.NoError(t, err)
	}

}

func TestLoadImage_AccessMode_Sequential(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(assets + "test.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := LoadImage(src, WithAccessMode(AccessSequential))
	if assert.NoError(t, err) {
		assert.NotNil(t, img)
		// check sequential access by encoding twice where the second fails
		_, _, err = img.Export(ExportParams{})
		assert.NoError(t, err)
		_, _, err = img.Export(ExportParams{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "out of order")
	}

}

func TestImageTypeSupport_HEIF(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(assets + "citron.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	if assert.NoError(t, err) {
		assert.NotNil(t, img)
	}

	_, imageType, err := img.Export(ExportParams{})
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, imageType)
}

func TestImageRef_HasAlpha(t *testing.T) {
	Startup(nil)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"image without alpha layer",
			assets + "test.png",
			false,
		},
		{
			"image with alpha layer",
			assets + "with_alpha.png",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.HasAlpha()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestImageRef_AddAlpha(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(assets + "test.png")
	assert.NoError(t, err)

	err = image.AddAlpha()
	assert.NoError(t, err)
	assert.True(t, image.HasAlpha(), "has alpha")

	_, _, err = image.Export(ExportParams{})
	assert.NoError(t, err)
}

func TestImageRef_AddAlpha__AlreadyHasAlpha__Idempotent(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(assets + "with_alpha.png")
	assert.NoError(t, err)
	err = image.AddAlpha()
	assert.NoError(t, err)

	assert.True(t, image.HasAlpha(), "has alpha")
	_, _, err = image.Export(ExportParams{})
	assert.NoError(t, err)
}

func TestImageRef_HasProfile(t *testing.T) {
	Startup(nil)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"image with profile",
			assets + "with_icc_profile.jpg",
			true,
		},
		{
			"image without profile",
			assets + "without_icc_profile.jpg",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.HasProfile()
			assert.Equal(t, tt.want, got)
		})
	}
}
