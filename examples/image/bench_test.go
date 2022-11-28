package image

import (
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"runtime"
	"testing"
	"time"
)

func init() {
	vips.LoggingSettings(func(domain string, level vips.LogLevel, msg string) {
		fmt.Println(domain, level, msg)
	}, vips.LogLevelDebug)

	// Disable the cache so that after GC, libvips does not hold reference to any object
	vips.Startup(&vips.Config{
		ConcurrencyLevel: 0,
		MaxCacheFiles:    0,
		MaxCacheMem:      0,
		MaxCacheSize:     0,
		ReportLeaks:      false,
		CacheTrace:       false,
		CollectStats:     false,
	})
}

func BenchmarkBlack(b *testing.B) {
	blackTest := func(n int) {
		image, err := vips.Black(100+n, 100+n)
		if err != nil {
			panic(err)
		}
		if _, _, err = image.ExportNative(); err != nil {
			panic(err)
		}
	}

	for i := 0; i < b.N; i++ {
		blackTest(i)
	}
	// Forcing the GC to run to clear all the memory
	runtime.GC()
	runtime.GC()

	// Waiting for 1 second for GC to complete
	time.Sleep(1 * time.Second)

	vips.PrintObjectReport("BenchmarkBlack")
}
