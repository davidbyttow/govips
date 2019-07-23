package vips

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// todo: add missing tests...

func TestLoadImage_AccessMode_Default(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "test.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	defer img.Close()

	if assert.NoError(t, err) {
		assert.NotNil(t, img)
		// check random access by encoding twice
		_, _, err = img.Export(nil)
		assert.NoError(t, err)
		_, _, err = img.Export(nil)
		assert.NoError(t, err)
	}
}

func TestLoadImage_AccessMode_Random(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "test.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	defer img.Close()

	if assert.NoError(t, err) {
		assert.NotNil(t, img)
		// check random access by encoding twice
		_, _, err = img.Export(nil)
		assert.NoError(t, err)
		_, _, err = img.Export(nil)
		assert.NoError(t, err)
	}
}

func TestImageTypeSupport_HEIF(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "citron.heic")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	defer img.Close()

	if assert.NoError(t, err) {
		assert.NotNil(t, img)
	}

	_, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypeHEIF, metadata.Format)
}

func TestImageRef_HasAlpha__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "with_alpha.png")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasAlpha()
	assert.True(t, got)
}

func TestImageRef_HasAlpha__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "test.png")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasAlpha()
	assert.False(t, got)
}

func TestImageRef_AddAlpha(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "test.png")
	require.NoError(t, err)
	defer img.Close()

	err = img.AddAlpha()
	assert.NoError(t, err)
	assert.True(t, img.HasAlpha(), "has alpha")

	_, _, err = img.Export(nil)
	assert.NoError(t, err)
}

func TestImageRef_AddAlpha__Idempotent(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "with_alpha.png")
	require.NoError(t, err)
	defer img.Close()

	err = img.AddAlpha()
	assert.NoError(t, err)

	assert.True(t, img.HasAlpha(), "has alpha")
	_, _, err = img.Export(nil)
	assert.NoError(t, err)
}

func TestImageRef_HasProfile__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "with_icc_profile.jpg")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasProfile()
	assert.True(t, got)
}

func TestImageRef_HasProfile__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "without_icc_profile.jpg")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasProfile()
	assert.False(t, got)
}

func TestImageRef_Close(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "test.png")
	assert.NoError(t, err)

	image.Close()

	assert.Nil(t, image.image)

	// todo: how do I check that vips GC it as well?

	PrintObjectReport("Final")
}

func TestImageRef_Linear1(t *testing.T) {
	image, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Linear1(3, 4)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Sharpen(t *testing.T) {
	image, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Sharpen(3, 4, 5)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Embed(t *testing.T) {
	image, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Embed(10, 20, 100, 200, ExtendBlack)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_GetOrientation_HasEXIF(t *testing.T) {
	image, err := NewImageFromFile(resources + "rotate.jpg")
	require.NoError(t, err)
	defer image.Close()

	o := image.GetOrientation()

	assert.Equal(t, 6, o)
}

func TestImageRef_GetOrientation_NoEXIF(t *testing.T) {
	image, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image.Close()

	o := image.GetOrientation()

	assert.Equal(t, 0, o)
}
