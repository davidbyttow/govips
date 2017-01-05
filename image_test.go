package govips

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromFile(t *testing.T) {
	image, err := NewImageFromFile("fixtures/canyon.jpg", nil)
	require.Nil(t, err)
	assert.Equal(t, 2560, image.Width())
	assert.Equal(t, 1600, image.Height())
}

func TestWriteToFile(t *testing.T) {
	image, err := NewImageFromFile("fixtures/canyon.jpg", nil)
	require.Nil(t, err)

	image = image.Resize(0.25)

	err = image.WriteToFile("./fixtures/canyon-out.jpg", nil)
	require.Nil(t, err)
}

func TestWriteToBytes(t *testing.T) {
	buf, err := ioutil.ReadFile("fixtures/canyon.jpg")
	require.Nil(t, err)

	image, err := NewImageFromBuffer(buf, nil)
	require.Nil(t, err)

	image = image.Resize(0.25)

	buf, err = image.WriteToBuffer(".jpeg", nil)
	require.Nil(t, err)
	assert.True(t, len(buf) > 0)

	image, err = NewImageFromBuffer(buf, nil)
	require.Nil(t, err)
	assert.Equal(t, 640, image.Width())
	assert.Equal(t, 400, image.Height())
}
