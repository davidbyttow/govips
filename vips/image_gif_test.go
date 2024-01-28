package vips

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestI_GIF_Animated_Pages(t *testing.T) {
	Startup(nil)
	image, err := NewImageFromFile(resources + "gif-animated.gif")
	require.NoError(t, err)

	pages := image.Pages()
	assert.Equal(t, 8, pages)
}

func TestImage_GIF_Animated(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		nil,
		nil,
		exportGif(NewGifExportParams()))
}

func TestImage_GIF_Animated_ExportNative(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		nil,
		nil,
		nil)
}

func TestImage_GIF_Animated_to_WebP(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		nil,
		nil,
		exportWebp(NewWebpExportParams()))
}

func TestImage_GIF_Animated_to_WebP_Extra_Params(t *testing.T) {
	exportParams := NewWebpExportParams()
	exportParams.MaxKeyFrames = 100
	exportParams.MinKeyFrames = 10
	exportParams.MinSize = true
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		nil,
		nil,
		exportWebp(exportParams))
}

func TestImage_GIF_Animated_Resize(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		func(img *ImageRef) error {
			return img.Resize(2, KernelCubic)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_ResizeWithVScale(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		func(img *ImageRef) error {
			return img.ResizeWithVScale(0.5, 0.78, KernelCubic)
		},
		func(img *ImageRef) {
			assert.Equal(t, 3, img.Pages())
			assert.Equal(t, 100, img.GetPageHeight())
		},
		nil)
}

func TestImage_GIF_Animated_Rotate90(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.Rotate(Angle90)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_Rotate270(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.Rotate(Angle270)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_Embed(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.Embed(10, 20, 200, 250, ExtendWhite)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_EmbedBackground(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.EmbedBackground(10, 20, 200, 250, &Color{
				R: 255, G: 255, B: 0,
			})
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_ExtractArea(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.ExtractArea(10, 20, 80, 90)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_Crop(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			return img.Crop(10, 20, 20, 20)
		},
		func(img *ImageRef) {
			assert.Equal(t, img.GetPageHeight(), 20)
		},
		nil)
}

func TestImage_GIF_Animated_PageDelay(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		-1,
		func(img *ImageRef) error {
			delay, err := img.PageDelay()
			require.NoError(t, err)
			assert.Equal(t, delay, []int{100, 100, 100, 100, 100, 100, 100, 100})
			for i := range delay {
				delay[i] = delay[i] * 2
			}
			return img.SetPageDelay(delay)
		},
		nil,
		nil)
}

func TestImage_GIF_Animated_PageDelay2(t *testing.T) {
	goldenAnimatedTest(t, resources+"gif-animated.gif",
		3,
		func(img *ImageRef) error {
			delay, err := img.PageDelay()
			require.NoError(t, err)
			assert.Equal(t, delay, []int{100, 100, 100, 100, 100, 100, 100, 100})
			for i := range delay {
				delay[i] = delay[i] * 2
			}
			return img.SetPageDelay(delay)
		},
		nil,
		nil)
}
