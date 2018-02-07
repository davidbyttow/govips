package vips_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbed(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(512, 256)
	})
}

func TestFlatten(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.BackgroundColor(vips.Color{R: 255, G: 192, B: 203})
	})
}

func TestResizeCrop(t *testing.T) {
	goldenTest(t, "fixtures/colors.png", func(tx *vips.Transform) {
		tx.Resize(100, 300).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestResizeShapes(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256)
	})
}

func TestRelativeResizeShapes(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.ScaleHeight(0.5)
	})
}

func TestCenterCrop(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestBottomRightCrop(t *testing.T) {
	goldenTest(t, "fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(vips.ResizeStrategyCrop).
			Anchor(vips.AnchorBottomRight)
	})
}

func TestOffsetCrop(t *testing.T) {
	goldenTest(t, "fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(500, 720).
			CropOffsetX(120).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestOffsetCropBounds(t *testing.T) {
	goldenTest(t, "fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(100, 100).
			CropOffsetX(120).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestRelativeOffsetCrop(t *testing.T) {
	goldenTest(t, "fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(500, 720).
			CropRelativeOffsetX(0.1066).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestRotate(t *testing.T) {
	goldenTest(t, "fixtures/canyon.jpg", func(tx *vips.Transform) {
		tx.Rotate(vips.Angle90)
	})
}

func goldenTest(t *testing.T, file string, fn func(t *vips.Transform)) {
	if testing.Short() {
		return
	}
	tx := vips.NewTransform().LoadFile(file)
	fn(tx)
	buf, err := tx.Apply()
	require.NoError(t, err)
	assertGoldenMatch(t, file, buf)
}

func assertGoldenMatch(t *testing.T, file string, buf []byte) {
	i := strings.LastIndex(file, ".")
	if i < 0 {
		panic("bad filename")
	}
	goldenFile := file[:i] + "." + t.Name() + ".golden" + file[i:]
	golden, _ := ioutil.ReadFile(goldenFile)
	if golden != nil {
		if !assert.Equal(t, golden, buf) {
			failed := file[:i] + "." + t.Name() + ".failed" + file[i:]
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
