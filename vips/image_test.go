package vips

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"runtime"
	"testing"
)

// todo: add missing tests...

func TestImageRef_PNG(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "png-24bit.png")
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

func TestImageRef_HEIF(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "heic-24bit-exif.heic")
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

func TestImageRef_HEIF_MIF1(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "heic-24bit.heic")
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

func TestImageRef_BMP(t *testing.T) {
	Startup(nil)

	raw, err := ioutil.ReadFile(resources + "bmp.bmp")
	require.NoError(t, err)

	img, err := NewImageFromBuffer(raw)
	require.NoError(t, err)
	defer img.Close()

	if assert.NoError(t, err) {
		assert.NotNil(t, img)
	}

	_, metadata, err := img.Export(nil)
	assert.NoError(t, err)
	assert.Equal(t, ImageTypePNG, metadata.Format)
}

func TestImageRef_HasAlpha__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit+alpha.png")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasAlpha()
	assert.True(t, got)
}

func TestImageRef_HasAlpha__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasAlpha()
	assert.False(t, got)
}

func TestImageRef_AddAlpha(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit.png")
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

	img, err := NewImageFromFile(resources + "png-24bit+alpha.png")
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

	img, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasProfile()
	assert.True(t, got)
}

func TestImageRef_HasProfile__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	defer img.Close()

	got := img.HasProfile()
	assert.False(t, got)
}

func TestImageRef_RemoveMetadata(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	defer img.Close()

	err = img.RemoveMetadata()
	require.NoError(t, err)

	assert.True(t, img.HasProfile())
}

func TestImageRef_Linear1(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Linear1(3, 4)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Linear(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit+alpha.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Linear([]float64{1.1, 1.2, 1.3, 1.4}, []float64{1, 2, 3, 4})
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Sharpen(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Sharpen(3, 4, 5)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Modulate__Alpha(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Modulate(0.1, 0.2, 90)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Modulate(t *testing.T) {
	image, err := NewImageFromFile(resources + "jpg-24bit-icc-iec.jpg")
	require.NoError(t, err)
	defer image.Close()

	err = image.Modulate(0.1, 0.2, 90)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_Embed(t *testing.T) {
	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Embed(10, 20, 100, 200, ExtendBlack)
	require.NoError(t, err)

	_, _, err = image.Export(nil)
	require.NoError(t, err)
}

func TestImageRef_GetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)
	defer image.Close()

	o := image.GetOrientation()

	assert.Equal(t, 6, o)
}

func TestImageRef_GetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	o := image.GetOrientation()

	assert.Equal(t, 0, o)
}

func TestImageRef_ExtractArea(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.ExtractArea(1, 2, 3, 4)
	require.NoError(t, err)
}

func TestImageRef_ExtractArea__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.ExtractArea(1, 2, 10000, 4)
	require.Error(t, err)
}

func TestImageRef_Resize__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Resize(-1, KernelLanczos3)
	require.Error(t, err)
}

func TestImageRef_Close(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	image.Close()

	assert.Nil(t, image.image)

	PrintObjectReport("Final")
}

func TestImageRef_Close__AlreadyClosed(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	assert.NoError(t, err)

	go image.Close()
	go image.Close()
	go image.Close()
	go image.Close()
	defer image.Close()
	image.Close()

	assert.Nil(t, image.image)
	runtime.GC()
}

func TestImageRef_NotImage(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "txt.txt")
	require.Error(t, err)
	require.Nil(t, image)
}
