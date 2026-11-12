package vips

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	running = false
	require.NoError(t, Startup(&Config{CollectStats: true, CacheTrace: true}))
	running = false
	require.NoError(t, startupIfNeeded())
}

func TestSetBlockUntrusted(t *testing.T) {
	require.NoError(t, Startup(nil))
	if !IsTypeSupported(ImageTypePDF) {
		t.Skip("PDF loader is not available")
	}

	SetBlockUntrusted(true)
	t.Cleanup(func() { SetBlockUntrusted(false) })

	img, err := NewImageFromFile(resources + "pdf.pdf")
	if img != nil {
		defer img.Close()
	}
	require.Error(t, err)
	require.Contains(t, err.Error(), "blocked")

	img, err = NewImageFromFile(resources + "png-24bit.png")
	require.NoError(t, err)
	defer img.Close()

	SetBlockUntrusted(false)
	img, err = NewImageFromFile(resources + "pdf.pdf")
	require.NoError(t, err)
	defer img.Close()
}
