package vips_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wix-playground/govips/pkg/vips"
)

func TestLoadImage_AccessMode(t *testing.T) {
	srcBytes, err := ioutil.ReadFile("testdata/test.png")
	require.NoError(t, err)

	// defaults to random access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src)
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src, vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src, vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestNewImageFromFile_AccessMode(t *testing.T) {
	// defaults to random access
	{
		img, err := vips.NewImageFromFile("testdata/test.png")
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		img, err := vips.NewImageFromFile("testdata/test.png", vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		img, err := vips.NewImageFromFile("testdata/test.png", vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestNewImageFromBuffer_AccessMode(t *testing.T) {
	src, err := ioutil.ReadFile("testdata/test.png")
	require.NoError(t, err)

	// defaults to random access
	{
		img, err := vips.NewImageFromBuffer(src)
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// random access
	{
		img, err := vips.NewImageFromBuffer(src, vips.WithAccessMode(vips.AccessRandom))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check random access by encoding twice
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
		}
	}

	// sequential access
	{
		img, err := vips.NewImageFromBuffer(src, vips.WithAccessMode(vips.AccessSequential))
		if assert.NoError(t, err) {
			assert.NotNil(t, img)
			// check sequential access by encoding twice where the second fails
			_, _, err = img.Export(vips.ExportParams{})
			assert.NoError(t, err)
			_, _, err = img.Export(vips.ExportParams{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "out of order")
		}
	}
}

func TestImageRef_HasAlpha(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"image without alpha layer",
			"testdata/test.png",
			false,
		},
		{
			"image with alpha layer",
			"testdata/with_alpha.png",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := vips.NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.HasAlpha()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestImageRef_HasProfile(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			"image with profile",
			"testdata/with_icc_profile.jpg",
			true,
		},
		{
			"image without profile",
			"testdata/without_icc_profile.jpg",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ref, err := vips.NewImageFromFile(tt.path)
			require.NoError(t, err)
			got := ref.HasProfile()
			assert.Equal(t, tt.want, got)
		})
	}
}
