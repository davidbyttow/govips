package vips

import (
	"runtime"
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
		ImageTypeWEBP,
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

	run(func() {
		var wg sync.WaitGroup
		for i, tr := range transforms {
			wg.Add(1)
			go func(i int, tr transform) {
				defer wg.Done()

				in, err := NewImageFromFile(resources + "jpg-24bit-icc-iec.jpg")
				require.NoError(t, err)
				defer in.Close()

				buf, _, err := NewTransform().
					ResizeStrategy(tr.Resize).
					Resize(tr.Width, tr.Height).
					Flip(tr.Flip).
					Kernel(tr.Kernel).
					Format(tr.Format).
					GaussianBlur(tr.Blur).
					Interpolator(tr.Interp).
					Zoom(tr.Zoom, tr.Zoom).
					Quality(tr.Quality).
					Compression(tr.Compression).
					ApplyAndExport(in)
				require.NoError(t, err)

				out, err := NewImageFromBuffer(buf)
				require.NoError(t, err)
				defer out.Close()

				assert.Equal(t, tr.Format, out.Format())
			}(i, tr)
		}
		wg.Wait()
	})
}

func run(fn func()) {
	Startup(nil)

	fn()

	//ShutdownThread()
	//Shutdown()
	runtime.GC()

	PrintObjectReport("Finished")
}
