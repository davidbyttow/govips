package vips

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestXYZ(t *testing.T) {
	Startup(nil)

	img, err := XYZ(100, 100)
	require.NoError(t, err)
	require.NotNil(t, img)
}
