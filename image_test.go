package govips

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFromFile(t *testing.T) {
	image, err := NewImageFromFile("fixtures/canyon.jpg")
	require.Nil(t, err)
	assert.Equal(t, 2560, image.Width())
	assert.Equal(t, 1600, image.Height())
}

func TestWriteToFile(t *testing.T) {
	image, err := NewImageFromFile("fixtures/canyon.jpg")
	require.Nil(t, err)

	image = image.Resize(0.25)

	tempDir, err := ioutil.TempDir("", "TestWriteToFile")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	err = image.WriteToFile(tempDir + "/canyon-out.jpg")
	require.Nil(t, err)
}

func TestWriteToBytes(t *testing.T) {
	buf, err := ioutil.ReadFile("fixtures/canyon.jpg")
	require.Nil(t, err)

	image, err := NewImageFromBuffer(buf)
	require.Nil(t, err)

	image = image.Resize(0.25)

	buf, err = image.WriteToBuffer(".jpeg")
	require.Nil(t, err)
	assert.True(t, len(buf) > 0)

	image, err = NewImageFromBuffer(buf)
	require.Nil(t, err)
	assert.Equal(t, 640, image.Width())
	assert.Equal(t, 400, image.Height())
}

func TestLoadFromMemory(t *testing.T) {
	size := 200

	bytes := make([]byte, size*size*3)
	for i := 0; i < size*size; i++ {
		bytes[i*3] = 0xFF
		bytes[i*3+1] = 0
		bytes[i*3+2] = 0
	}

	image, err := NewImageFromMemory(bytes, size, size, 3, BandFormatUchar)
	require.Nil(t, err)

	tempDir, err := ioutil.TempDir("", "TestLoadFromMemory")
	require.Nil(t, err)
	defer os.RemoveAll(tempDir)

	err = image.WriteToFile(tempDir + "red-out.png")
	require.Nil(t, err)
}
