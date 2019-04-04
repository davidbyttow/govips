package vips_test

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"testing"

	"../vips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestNewImageFromNativeImage(t *testing.T) {
	tests := []struct {
		name      string
		img       draw.Image
		fillcolor color.Color
	}{
		{
			"image.Gray",
			image.NewGray(image.Rect(0, 0, 64, 64)),
			color.RGBA{0x7F, 0x7F, 0x7F, 0xFF},
		},
		{
			"image.RGBA",
			image.NewRGBA(image.Rect(0, 0, 64, 64)),
			color.RGBA{0x0, 0x0, 0xFF, 0xFF},
		},
		{
			"image.Gray16",
			image.NewGray16(image.Rect(0, 0, 64, 64)),
			color.RGBA{0x7F, 0x7F, 0x7F, 0xFF},
		},
		{
			"image.RGBA64",
			image.NewRGBA64(image.Rect(0, 0, 64, 64)),
			color.RGBA{0x0, 0x0, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill the image with a color to be able to test it
			fill := image.NewUniform(tt.fillcolor)
			draw.Draw(tt.img, tt.img.Bounds(), fill, image.ZP, draw.Src)
			vimg, err := vips.NewImageFromNativeImage(tt.img)
			if assert.NoError(t, err) {
				assert.NotNil(t, vimg)
			}
			// Encode the image to be able to examine it pixel by pixel
			buf, _, err := vips.NewTransform().
				Image(vimg).
				Format(vips.ImageTypePNG).
				OutputBytes().
				Apply()
			assert.NoError(t, err)
			// Decode the result
			decoded, err := png.Decode(bytes.NewBuffer(buf))
			if assert.NoError(t, err) {
				// Check it pixel by pixel
			PixelCheck:
				for x := decoded.Bounds().Min.X; x < decoded.Bounds().Max.X; x++ {
					for y := decoded.Bounds().Min.Y; y < decoded.Bounds().Max.Y; y++ {
						ar, ag, ab, aa := tt.fillcolor.RGBA()
						br, bg, bb, ba := decoded.At(x, y).RGBA()
						if ar != br || ag != bg || ab != bb || aa != ba {
							// Only assert if it would fail, and then break or the test hangs
							assert.Equal(t, tt.fillcolor, decoded.At(x, y))
							break PixelCheck
						}
					}
				}
			}
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
