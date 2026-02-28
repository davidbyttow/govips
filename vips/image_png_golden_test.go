package vips

import (
	"bytes"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage_PNG_64bit_OptimizeICCProfile(t *testing.T) {
	goldenTest(t, resources+"png-alpha-64bit.png",
		func(img *ImageRef) error {
			return img.OptimizeICCProfile()
		},
		nil,
		exportPng(NewPngExportParams()))
}

func TestImage_Resize_Downscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(0.9, KernelLanczos3)
	}, nil, nil)
}

func TestImage_Resize_Upscale_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Resize(1.1, KernelLanczos3)
	}, nil, nil)
}

func TestImage_Embed_ExtendWhite_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Embed(0, 0, 1000, 500, ExtendWhite)
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{255, 255, 255, 255})
	}, nil)
}

func TestImage_EmbedBackgroundRGBA_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.EmbedBackgroundRGBA(0, 0, 1000, 500, &ColorRGBA{R: 238, G: 238, B: 238, A: 50})
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{238, 238, 238, 50})
	}, nil)
}

func TestImage_EmbedBackground_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.EmbedBackground(0, 0, 1000, 500, &Color{R: 238, G: 238, B: 238})
	}, func(img *ImageRef) {
		point, err := img.GetPoint(999, 0)
		assert.NoError(t, err)
		assert.Equal(t, point, []float64{238, 238, 238, 255})
	}, nil)
}

func TestImageRef_PngToWebp_OptimizeICCProfile_Lossless(t *testing.T) {
	exportParams := NewWebpExportParams()
	exportParams.Quality = 90
	exportParams.Lossless = true

	testWebpOptimizeIccProfile(t, exportParams)
}

func TestImage_AutoRotate_0(t *testing.T) {
	// TODO: revisit - libvips 8.15.1 returns orientation=1 instead of 0 for
	// images with no orientation tag. Behavior varies across libvips versions.
	t.Skip("orientation behavior differs across libvips versions")
	goldenTest(t, resources+"png-24bit.png",
		func(img *ImageRef) error {
			return img.AutoRotate()
		},
		func(result *ImageRef) {
			assert.Equal(t, 0, result.Orientation())
		}, nil)
}

func TestImage_Sharpen_24bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		// usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Sharpen_8bit_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		// usm_0.66_1.00_0.01
		sigma := 1 + (0.66 / 2)
		x1 := 0.01 * 100
		m2 := 1.0

		return img.Sharpen(sigma, x1, m2)
	}, nil, nil)
}

func TestImage_Sobel(t *testing.T) {
	goldenTest(t, resources+"png-8bit+alpha.png", func(img *ImageRef) error {
		return img.Sobel()
	}, nil, nil)
}

func TestImage_Modulate_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Modulate(1.1, 1.2, 0)
	}, nil, nil)
}

func TestImage_ModulateHSV_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.ModulateHSV(1.1, 1.2, 120)
	}, nil, nil)
}

func TestImageRef_Linear1(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		return img.Linear1(3, 4)
	}, nil, nil)
}

func TestImageRef_Linear_Alpha(t *testing.T) {
	goldenTest(t, resources+"png-24bit+alpha.png", func(img *ImageRef) error {
		return img.Linear([]float64{1.1, 1.2, 1.3, 1.4}, []float64{1, 2, 3, 4})
	}, nil, nil)
}

func TestImage_Rank(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		return img.Rank(15, 15, 224)
	}, nil, nil)
}

func TestImage_GetPointWhite(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 255.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPointYellow(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(400, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 0.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPointWhiteR(t *testing.T) {
	goldenTest(t, resources+"png-24bit.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 3, len(point))
		assert.Equal(t, 255.0, point[0])
		assert.Equal(t, 255.0, point[1])
		assert.Equal(t, 255.0, point[2])

		return err
	}, nil, nil)
}

func TestImage_GetPoint_WithAlpha(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		point, err := img.GetPoint(10, 10)

		assert.Equal(t, 4, len(point))
		assert.Equal(t, 0.0, point[0])
		assert.Equal(t, 0.0, point[1])
		assert.Equal(t, 0.0, point[2])
		assert.Equal(t, 255.0, point[3])

		return err
	}, nil, nil)
}

func TestImage_GetPoint_WithAlpha2(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		point, err := img.GetPoint(0, 0)

		assert.Equal(t, 4, len(point))
		assert.Equal(t, 0.0, point[0])
		assert.Equal(t, 0.0, point[1])
		assert.Equal(t, 0.0, point[2])
		assert.Equal(t, 0.0, point[3])

		return err
	}, nil, nil)
}

func TestImage_Decode_PNG(t *testing.T) {
	goldenTest(t, resources+"png-8bit.png", func(img *ImageRef) error {
		goImg, err := img.ToImage(nil)
		assert.NoError(t, err)

		buf := new(bytes.Buffer)
		err = png.Encode(buf, goImg)
		assert.Nil(t, err)

		config, format, err := image.DecodeConfig(buf)
		assert.Nil(t, err)
		assert.Equal(t, "png", format)
		assert.NotNil(t, config)
		assert.Equal(t, 150, config.Height)
		assert.Equal(t, 200, config.Width)
		return nil
	}, nil, nil)
}

func TestImage_ExtractBand(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		return img.ExtractBand(2, 1)
	}, nil, nil)
}

func TestImage_Flatten(t *testing.T) {
	goldenTest(t, resources+"with_alpha.png", func(img *ImageRef) error {
		return img.Flatten(&Color{R: 32, G: 64, B: 128})
	}, nil, nil)
}
