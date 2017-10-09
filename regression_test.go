package vips_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
)

func TestResizeCrop(t *testing.T) {
	if testing.Short() {
		return
	}
	file := "fixtures/colors.png"
	buf, err := vips.NewTransform().
		LoadFile(file).
		Resize(100, 300).
		ResizeStrategy(vips.ResizeStrategyCrop).
		Apply()
	assert.NoError(t, err)
	assertGoldenMatch(t, file, buf)
}

func TestResizeShapes(t *testing.T) {
	if testing.Short() {
		return
	}
	file := "fixtures/shapes.png"
	buf, err := vips.NewTransform().
		LoadFile(file).
		Resize(341, 256).
		Apply()
	assert.NoError(t, err)
	assertGoldenMatch(t, file, buf)
}

func TestCenterCrop(t *testing.T) {
	if testing.Short() {
		return
	}
	file := "fixtures/shapes.png"
	buf, err := vips.NewTransform().
		LoadFile(file).
		Resize(341, 256).
		ResizeStrategy(vips.ResizeStrategyCrop).
		Apply()
	assert.NoError(t, err)
	t.Name()
	assertGoldenMatch(t, file, buf)
}

func TestBottomRightCrop(t *testing.T) {
	if testing.Short() {
		return
	}
	file := "fixtures/shapes.png"
	buf, err := vips.NewTransform().
		LoadFile(file).
		Resize(341, 256).
		ResizeStrategy(vips.ResizeStrategyCrop).
		Anchor(vips.AnchorBottomRight).
		Apply()
	assert.NoError(t, err)
	t.Name()
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
