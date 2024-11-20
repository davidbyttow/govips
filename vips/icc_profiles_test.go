package vips

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func Test_ICCProfileInitialisation(t *testing.T) {
	nonDefaultProfile := "non-default"
	err := ensureLoadICCPath(&nonDefaultProfile)
	require.NoError(t, err)

	err = ensureLoadICCPath(&sRGBV2MicroICCProfilePathToken)
	require.NoError(t, err)
	path, err := GetSRGBV2MicroICCProfilePath()
	require.NoError(t, err)
	assertIccProfile(t, sRGBV2MicroICCProfile, path)

	err = ensureLoadICCPath(&sGrayV2MicroICCProfilePathToken)
	require.NoError(t, err)
	path, err = GetSGrayV2MicroICCProfilePath()
	require.NoError(t, err)
	assertIccProfile(t, sGrayV2MicroICCProfile, path)

	err = ensureLoadICCPath(&sRGBIEC6196621ICCProfilePathToken)
	require.NoError(t, err)
	path, err = GetSRGBIEC6196621ICCProfilePath()
	require.NoError(t, err)
	assertIccProfile(t, sRGBIEC6196621ICCProfile, path)

	err = ensureLoadICCPath(&genericGrayGamma22ICCProfilePathToken)
	require.NoError(t, err)
	path, err = GetGenericGrayGamma22ICCProfilePath()
	require.NoError(t, err)
	assertIccProfile(t, genericGrayGamma22ICCProfile, path)
}

func assertIccProfile(t *testing.T, expectedProfile []byte, path string) {
	loadedProfile, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, expectedProfile, loadedProfile)
}
