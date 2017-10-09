package vips_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			CropOffset(120, 0).Kernel(vips.KernelNearest).
			ResizeStrategy(vips.ResizeStrategyCrop)
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
		assert.Equal(t, golden, buf)
		return
	}
	t.Log("Writing golden file: " + goldenFile)
	err := ioutil.WriteFile(goldenFile, buf, 0644)
	assert.NoError(t, err)
}
