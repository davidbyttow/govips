package vips

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitConfig(t *testing.T) {
	running = false
	require.NoError(t, Startup(&Config{CollectStats: true, CacheTrace: true}))
	running = false
	startupIfNeeded()
}
