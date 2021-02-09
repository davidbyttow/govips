package vips

import (
	"testing"
)

func TestInitConfig(t *testing.T) {
	running = false
	Startup(&Config{CollectStats: true, CacheTrace: true})
	running = false
	startupIfNeeded()
}
