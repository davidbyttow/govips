package vips

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform(t *testing.T) {
	Startup(nil)

	in, err := NewImageFromFile(resources + "canyon.jpg")
	require.NoError(t, err)

	buf, metadata, err := NewTransform().
		Scale(0.25).
		ApplyAndExport(in)

	require.NoError(t, err)
	require.True(t, len(buf) > 0)
	assert.Equal(t, metadata.Format, ImageTypeJPEG)

	out, err := NewImageFromBuffer(buf)
	require.NoError(t, err)

	assert.Equal(t, 640, out.Width())
	assert.Equal(t, 400, out.Height())

	in.Close()
	out.Close()

	PrintObjectReport("Final")
}

func TestResize(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(512, 256)
	})
}

func TestFlatten(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.BackgroundColor(&Color{R: 255, G: 192, B: 203}).StripProfile()
	})
}

func TestResizeWithICC(t *testing.T) {
	goldenTest(t, resources+"icc.jpg", func(tx *Transform) {
		tx.StripMetadata()
		tx.ResizeWidth(300)
	})
}

func TestResizeAndStripICC(t *testing.T) {
	goldenTest(t, resources+"icc.jpg", func(tx *Transform) {
		tx.StripMetadata().ResizeWidth(300).StripProfile()
	})
}

func TestResizeCrop(t *testing.T) {
	goldenTest(t, resources+"colors.png", func(tx *Transform) {
		tx.Resize(100, 300).
			ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestResizeShapes(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256)
	})
}

func TestRelativeResizeShapes(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.ScaleHeight(0.5)
	})
}

func TestCenterCrop(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestBottomRightCrop(t *testing.T) {
	goldenTest(t, resources+"shapes.png", func(tx *Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(ResizeStrategyCrop).
			Anchor(AnchorBottomRight)
	})
}

func TestOffsetCrop(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(500, 720).
			CropOffsetX(120).
			ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestOffsetCropBounds(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(100, 100).
			CropOffsetX(120).
			ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestRelativeOffsetCrop(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Resize(500, 720).
			CropRelativeOffsetX(0.1066).
			ResizeStrategy(ResizeStrategyCrop)
	})
}

func TestRotate(t *testing.T) {
	goldenTest(t, resources+"canyon.jpg", func(tx *Transform) {
		tx.Rotate(Angle90)
	})
}

func TestAutoRotate(t *testing.T) {
	goldenTest(t, resources+"canyon.jpg", func(tx *Transform) {
		tx.AutoRotate()
	})
}

func TestScale3x(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.Scale(3.0)
	})
}

func TestMaxScale(t *testing.T) {
	goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
		tx.MaxScale(1.0).ResizeWidth(100000)
	})
}

func TestOverlay(t *testing.T) {
	if testing.Short() {
		return
	}
	var tomatoesData, cloverData []byte
	t.Run("tomatoes", func(t *testing.T) {
		tomatoesData = goldenTest(t, resources+"tomatoes.png", func(tx *Transform) {
			tx.ResizeWidth(320)
		})
	})
	t.Run("clover", func(t *testing.T) {
		cloverData = goldenTest(t, resources+"clover.png", func(tx *Transform) {
			tx.ResizeWidth(64)
		})
	})
	tomatoes, err := NewImageFromBuffer(tomatoesData)
	require.NoError(t, err)
	defer tomatoes.Close()

	clover, err := NewImageFromBuffer(cloverData)
	require.NoError(t, err)
	defer clover.Close()

	err = tomatoes.Composite(clover, BlendModeOver, 0, 0)
	require.NoError(t, err)

	buf, _, err := tomatoes.Export(nil)
	require.NoError(t, err)
	assertGoldenMatch(t, resources+"tomatoes.png", buf)
}

func TestBandJoin(t *testing.T) {
	image1, err := NewImageFromFile(resources + "tomatoes.png")
	require.NoError(t, err)
	defer image1.Close()

	image2, err := NewImageFromFile(resources + "clover.png")
	require.NoError(t, err)
	defer image2.Close()

	err = image1.BandJoin(image2)
	require.NoError(t, err)

	buf, _, err := image1.Export(nil)
	require.NoError(t, err)
	assertGoldenMatch(t, resources+"tomatoes.png", buf)
}

func goldenTest(t *testing.T, file string, fn func(t *Transform)) []byte {
	if testing.Short() {
		return nil
	}

	Startup(nil)

	i, err := NewImageFromFile(file)
	require.NoError(t, err)
	defer i.Close()

	tx := NewTransform()

	fn(tx)

	buf, _, err := tx.ApplyAndExport(i)
	require.NoError(t, err)
	assertGoldenMatch(t, file, buf)

	return buf
}

func assertGoldenMatch(t *testing.T, file string, buf []byte) {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		panic("bad filename")
	}
	name := strings.Replace(t.Name(), "/", "_", -1)
	prefix := file[:i] + "." + name
	ext := file[i:]
	goldenFile := prefix + ".golden" + ext
	golden, _ := ioutil.ReadFile(goldenFile)
	if golden != nil {
		if !assert.Equal(t, golden, buf) {
			failed := prefix + ".failed" + ext
			err := ioutil.WriteFile(failed, buf, 0666)
			if err != nil {
				panic(err)
			}
		}
		return
	}
	t.Log("Writing golden file: " + goldenFile)
	err := ioutil.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
}
