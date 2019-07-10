package vips_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wix-playground/govips/pkg/vips"
)

func TestTransform(t *testing.T) {
	if testing.Short() {
		return
	}
	buf, format, err := vips.NewTransform().
		LoadFile("../../assets/fixtures/canyon.jpg").
		Scale(0.25).
		Apply()

	require.NoError(t, err)
	require.True(t, len(buf) > 0)
	assert.Equal(t, format, vips.ImageTypeJPEG)

	image, err := vips.NewImageFromBuffer(buf)
	require.NoError(t, err)

	assert.Equal(t, 640, image.Width())
	assert.Equal(t, 400, image.Height())

	image.Close()
	vips.ShutdownThread()
	vips.Shutdown()

	vips.PrintObjectReport("Final")
}
