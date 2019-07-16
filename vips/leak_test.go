package vips

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type size struct {
	w, h int
}

var (
	resizeStrategies = []ResizeStrategy{
		ResizeStrategyCrop,
		ResizeStrategyStretch,
		ResizeStrategyEmbed,
	}
	sizes = []size{
		{100, 100},
		{500, 0},
		{0, 500},
		{1000, 1000},
	}
	formats = []ImageType{
		ImageTypeJPEG,
		ImageTypePNG,
	}
)

type transform struct {
	Resize      ResizeStrategy
	Width       int
	Height      int
	Flip        FlipDirection
	Format      ImageType
	Zoom        int
	Blur        float64
	Kernel      Kernel
	Interp      Interpolator
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
					Flip:        FlipBoth,
					Kernel:      KernelLanczos3,
					Format:      format,
					Blur:        4,
					Interp:      InterpolateBicubic,
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
		for i, tr := range transforms {
			wg.Add(1)
			go func(i int, tr transform) {
				defer wg.Done()

				buf, _, err := NewTransform().
					LoadFile(assets+"canyon.jpg").
					ResizeStrategy(tr.Resize).
					Resize(tr.Width, tr.Height).
					Flip(tr.Flip).
					Kernel(tr.Kernel).
					Format(tr.Format).
					GaussBlur(tr.Blur).
					Interpolator(tr.Interp).
					Zoom(tr.Zoom, tr.Zoom).
					Quality(tr.Quality).
					Compression(tr.Compression).
					OutputBytes().
					Apply()
				require.NoError(t, err)

				image, err := NewImageFromBuffer(buf)
				require.NoError(t, err)
				defer image.Close()

				assert.Equal(t, tr.Format, image.Format())
			}(i, tr)
		}
		wg.Wait()
	})
}

func LeakTest(fn func()) {
	Startup(nil)

	fn()

	//Shutdown()

	PrintObjectReport("Finished")
}
