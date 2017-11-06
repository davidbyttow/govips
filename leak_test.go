package vips_test

import (
	"sync"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type size struct {
	w, h int
}

var (
	resizeStrategies = []vips.ResizeStrategy{
		vips.ResizeStrategyCrop,
		vips.ResizeStrategyStretch,
		vips.ResizeStrategyEmbed,
	}
	sizes = []size{
		{100, 100},
		{500, 0},
		{0, 500},
		{1000, 1000},
	}
	formats = []vips.ImageType{
		vips.ImageTypeJPEG,
		vips.ImageTypePNG,
	}
)

type transform struct {
	Resize      vips.ResizeStrategy
	Width       int
	Height      int
	Flip        vips.FlipDirection
	Format      vips.ImageType
	Zoom        int
	Blur        float64
	Kernel      vips.Kernel
	Interp      vips.Interpolator
	Quality     int
	Compression int
}

func TestCleanup(t *testing.T) {
	if testing.Short() {
		return
	}

	var transforms []transform
	for _, resize := range resizeStrategies {
		for _, size := range sizes {
			for _, format := range formats {
				t := transform{
					Resize:      resize,
					Width:       size.w,
					Height:      size.h,
					Flip:        vips.FlipBoth,
					Kernel:      vips.KernelLanczos3,
					Format:      format,
					Blur:        4,
					Interp:      vips.InterpolateBicubic,
					Zoom:        3,
					Quality:     80,
					Compression: 5,
				}
				transforms = append(transforms, t)
			}
		}
	}

	LeakTest(func() {
		var wg sync.WaitGroup
		for i, transform := range transforms {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				buf, err := vips.NewTransform().
					LoadFile("fixtures/canyon.jpg").
					ResizeStrategy(transform.Resize).
					Resize(transform.Width, transform.Height).
					Flip(transform.Flip).
					Kernel(transform.Kernel).
					Format(transform.Format).
					GaussBlur(transform.Blur).
					Interpolator(transform.Interp).
					Zoom(transform.Zoom, transform.Zoom).
					Quality(transform.Quality).
					Compression(transform.Compression).
					OutputBytes().
					Apply()
				require.NoError(t, err)

				image, err := vips.NewImageFromBuffer(buf)
				require.NoError(t, err)
				defer image.Close()

				assert.Equal(t, transform.Format, image.Format())
			}(i)
		}
		wg.Wait()
	})
}

func LeakTest(fn func()) {
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 1,
	})
	fn()
	vips.Shutdown()
	vips.PrintObjectReport("Finished")
}
