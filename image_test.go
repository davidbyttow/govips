package vips_test

import (
	"io/ioutil"
	"testing"

	vips "github.com/davidbyttow/govips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteToBytes(t *testing.T) {
	buf, err := ioutil.ReadFile("fixtures/canyon.jpg")
	require.NoError(t, err)

	vips.NewStreamFromBuffer(buf, func(stream *vips.ImageRef) error {
		err := stream.Resize(0.25, vips.InputInterpolator("interpolate", vips.InterpolateNoHalo))
		require.NoError(t, err)

		buf, err = stream.Export(vips.ExportOptions{})
		require.NoError(t, err)
		assert.True(t, len(buf) > 0)

		return err
	})

	stream, err := vips.OpenFromBuffer(buf)
	require.NoError(t, err)

	assert.Equal(t, 640, stream.Width())
	assert.Equal(t, 400, stream.Height())

	stream.Close()
	vips.ShutdownThread()
	vips.Shutdown()

	vips.PrintObjectReport("Final")
}
