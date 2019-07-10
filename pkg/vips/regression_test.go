package vips_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wix-playground/govips/pkg/vips"
)

func TestEmbed(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(512, 256)
	})
}

func TestFlatten(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.BackgroundColor(vips.Color{R: 255, G: 192, B: 203}).StripProfile()
	})
}

func TestResizeWithICC(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/icc.jpg", func(tx *vips.Transform) {
		tx.StripMetadata()
		tx.ResizeWidth(300)
	})
}

func TestResizeAndStripICC(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/icc.jpg", func(tx *vips.Transform) {
		tx.StripMetadata().ResizeWidth(300).StripProfile()
	})
}

func TestResizeCrop(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/colors.png", func(tx *vips.Transform) {
		tx.Resize(100, 300).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestResizeShapes(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256)
	})
}

func TestRelativeResizeShapes(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.ScaleHeight(0.5)
	})
}

func TestCenterCrop(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestBottomRightCrop(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/shapes.png", func(tx *vips.Transform) {
		tx.Resize(341, 256).
			ResizeStrategy(vips.ResizeStrategyCrop).
			Anchor(vips.AnchorBottomRight)
	})
}

func TestOffsetCrop(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(500, 720).
			CropOffsetX(120).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestOffsetCropBounds(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(100, 100).
			CropOffsetX(120).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestRelativeOffsetCrop(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Resize(500, 720).
			CropRelativeOffsetX(0.1066).
			ResizeStrategy(vips.ResizeStrategyCrop)
	})
}

func TestRotate(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/canyon.jpg", func(tx *vips.Transform) {
		tx.Rotate(vips.Angle90)
	})
}

func TestScale3x(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.Scale(3.0)
	})
}

func TestMaxScale(t *testing.T) {
	goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
		tx.MaxScale(1.0).ResizeWidth(100000)
	})
}

func TestOverlay(t *testing.T) {
	if testing.Short() {
		return
	}
	var tomatoesData, cloverData []byte
	t.Run("tomatoes", func(t *testing.T) {
		tomatoesData = goldenTest(t, "../../assets/fixtures/tomatoes.png", func(tx *vips.Transform) {
			tx.ResizeWidth(320)
		})
	})
	t.Run("clover", func(t *testing.T) {
		cloverData = goldenTest(t, "../../assets/fixtures/clover.png", func(tx *vips.Transform) {
			tx.ResizeWidth(64)
		})
	})
	tomatoes, err := vips.NewImageFromBuffer(tomatoesData)
	require.NoError(t, err)
	clover, err := vips.NewImageFromBuffer(cloverData)
	require.NoError(t, err)

	err = tomatoes.Composite(clover, vips.BlendModeOver)
	require.NoError(t, err)
	buf, _, err := vips.NewTransform().Image(tomatoes).Apply()
	require.NoError(t, err)
	assertGoldenMatch(t, "../../assets/fixtures/tomatoes.png", buf)
}

func goldenTest(t *testing.T, file string, fn func(t *vips.Transform)) []byte {
	if testing.Short() {
		return nil
	}
	tx := vips.NewTransform().LoadFile(file)
	fn(tx)
	buf, _, err := tx.Apply()
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
