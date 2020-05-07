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

func TestImageRef_WebP(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)
	defer img.Close()

	// check random access by encoding twice
	_, _, err = img.Export(nil)
	assert.NoError(t, err)
	buf, _, err := img.Export(nil)
	assert.NoError(t, err)

	assert.Equal(t, 45252, len(buf))
}

func TestImageRef_WebP__ReducedEffort(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "webp+alpha.webp")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	require.NoError(t, err)
	require.NotNil(t, img)
	defer img.Close()

	params := NewDefaultWEBPExportParams()
	params.Effort = 2
	buf, _, err := img.Export(params)
	assert.NoError(t, err)
	assert.Equal(t, 48850, len(buf))
}

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

func TestImageRef_OverSizedMetadata(t *testing.T) {
	Startup(nil)

	srcBytes, err := ioutil.ReadFile(resources + "png-bad-metadata.png")
	require.NoError(t, err)

	src := bytes.NewReader(srcBytes)
	img, err := NewImageFromReader(src)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}

func TestImageRef_Resize__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.Resize(-1, KernelLanczos3)
	require.Error(t, err)
}

func TestImageRef_ExtractArea__Error(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.ExtractArea(1, 2, 10000, 4)
	require.Error(t, err)
}

func TestImageRef_HasAlpha__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit+alpha.png")
	require.NoError(t, err)
	defer img.Close()

	assert.True(t, img.HasAlpha())
}

func TestImageRef_HasAlpha__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	assert.False(t, img.HasAlpha())
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

	assert.True(t, img.HasProfile())
}

func TestImageRef_HasIPTC__True(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	defer img.Close()

	assert.True(t, img.HasIPTC())
}

func TestImageRef_HasIPTC__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	defer img.Close()

	assert.False(t, img.HasIPTC())
}

func TestImageRef_HasProfile__False(t *testing.T) {
	Startup(nil)

	img, err := NewImageFromFile(resources + "jpg-24bit.jpg")
	require.NoError(t, err)
	defer img.Close()

	assert.False(t, img.HasProfile())
}

func TestImageRef_GetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)
	defer image.Close()

	assert.Equal(t, 6, image.GetOrientation())
}

func TestImageRef_GetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	assert.Equal(t, 0, image.GetOrientation())
}

func TestImageRef_SetOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)
	defer image.Close()

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.GetOrientation())
}

func TestImageRef_SetOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.SetOrientation(5)
	require.NoError(t, err)

	assert.Equal(t, 5, image.GetOrientation())
}

func TestImageRef_RemoveOrientation__HasEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-6.jpg")
	require.NoError(t, err)
	defer image.Close()

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.GetOrientation())
}

func TestImageRef_RemoveOrientation__NoEXIF(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer image.Close()

	err = image.RemoveOrientation()
	require.NoError(t, err)

	assert.Equal(t, 0, image.GetOrientation())
}

func TestImageRef_RemoveMetadata__RetainsProfile(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	defer image.Close()

	require.True(t, image.HasIPTC())

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.False(t, image.HasIPTC())
	assert.True(t, image.HasICCProfile())
}

func TestImageRef_RemoveMetadata__RetainsOrientation(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-orientation-5.jpg")
	require.NoError(t, err)
	defer image.Close()

	err = image.RemoveMetadata()
	require.NoError(t, err)

	assert.Equal(t, 5, image.GetOrientation())
}

func TestImageRef_RemoveICCProfile(t *testing.T) {
	Startup(nil)

	image, err := NewImageFromFile(resources + "jpg-24bit-icc-adobe-rgb.jpg")
	require.NoError(t, err)
	defer image.Close()

	require.True(t, image.HasIPTC())

	err = image.RemoveICCProfile()
	require.NoError(t, err)

	assert.False(t, image.HasICCProfile())
	assert.True(t, image.HasIPTC())
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
