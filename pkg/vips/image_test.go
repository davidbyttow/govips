package vips_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/davidbyttow/govips/pkg/vips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadImage_AccessMode(t *testing.T) {
	srcBytes, err := ioutil.ReadFile("testdata/test.png")
	require.NoError(t, err)

	// random access
	{
		src := bytes.NewReader(srcBytes)
		img, err := vips.LoadImage(src, &vips.LoadOptions{Access: vips.AccessRandom})
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
		img, err := vips.LoadImage(src, &vips.LoadOptions{Access: vips.AccessSequential})
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
	// random access
	{
		img, err := vips.NewImageFromFile("testdata/test.png", &vips.LoadOptions{Access: vips.AccessRandom})
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
		img, err := vips.NewImageFromFile("testdata/test.png", &vips.LoadOptions{Access: vips.AccessSequential})
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

	// random access
	{
		img, err := vips.NewImageFromBuffer(src, &vips.LoadOptions{Access: vips.AccessRandom})
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
		img, err := vips.NewImageFromBuffer(src, &vips.LoadOptions{Access: vips.AccessSequential})
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
