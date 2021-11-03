package vips

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ICCProfileLengths(t *testing.T) {
	assert.Equal(t, len(sRGBV2MicroICCProfile), sRGBV2MicroICCProfileLength)
	assert.Equal(t, len(sGrayV2MicroICCProfile), sGrayV2MicroICCProfileLength)
	assert.Equal(t, len(sRGBIEC6196621ICCProfile), sRGBIEC6196621ICCProfileLength)
	assert.Equal(t, len(genericGrayGamma22ICCProfile), genericGrayGamma22ICCProfileLength)
}

func Test_ICCProfileInitialisation(t *testing.T) {
	initializeICCProfiles()

	srgbProfile, err := ioutil.ReadFile(filepath.Join(temporaryDirectory, sRGBV2MicroICCProfilePath))
	assert.NoError(t, err)
	assert.Equal(t, sRGBV2MicroICCProfile, srgbProfile)

	grayProfile, err := ioutil.ReadFile(filepath.Join(temporaryDirectory, sGrayV2MicroICCProfilePath))
	assert.NoError(t, err)
	assert.Equal(t, sGrayV2MicroICCProfile, grayProfile)

	srgbProfile2, err := ioutil.ReadFile(filepath.Join(temporaryDirectory, sRGBIEC6196621ICCProfilePath))
	assert.NoError(t, err)
	assert.Equal(t, sRGBIEC6196621ICCProfile, srgbProfile2)

	grayProfile2, err := ioutil.ReadFile(filepath.Join(temporaryDirectory, genericGrayGamma22ICCProfilePath))
	assert.NoError(t, err)
	assert.Equal(t, genericGrayGamma22ICCProfile, grayProfile2)
}
