// Package vips provides go bindings for libvips, a fast image processing library.
package vips

// #cgo pkg-config: vips
// #include <vips/vips.h>
// #include "govips.h"
// #include "icc_profiles.h"
import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

// Version is the full libvips version string (x.y.z)
const Version = string(C.VIPS_VERSION)

// MajorVersion is the libvips major component of the version string (x in x.y.z)
const MajorVersion = int(C.VIPS_MAJOR_VERSION)

// MinorVersion is the libvips minor component of the version string (y in x.y.z)
const MinorVersion = int(C.VIPS_MINOR_VERSION)

// MicroVersion is the libvips micro component of the version string (z in x.y.z)
// Also known as patch version
const MicroVersion = int(C.VIPS_MICRO_VERSION)

const (
	defaultConcurrencyLevel = 1
	defaultMaxCacheMem      = 50 * 1024 * 1024
	defaultMaxCacheSize     = 100
	defaultMaxCacheFiles    = 0
)

var (
	running             = false
	initLock            sync.Mutex
	statCollectorDone   chan struct{}
	once                sync.Once
	typeLoaders         = make(map[string]ImageType)
	supportedImageTypes = make(map[ImageType]bool)
)

// Config allows fine-tuning of libvips library
type Config struct {
	ConcurrencyLevel int
	MaxCacheFiles    int
	MaxCacheMem      int
	MaxCacheSize     int
	ReportLeaks      bool
	CacheTrace       bool
	CollectStats     bool
}

// Startup sets up the libvips support and ensures the versions are correct. Pass in nil for
// default configuration.
func Startup(config *Config) {
	initLock.Lock()
	runtime.LockOSThread()
	defer initLock.Unlock()
	defer runtime.UnlockOSThread()

	if running {
		govipsLog("govips", LogLevelInfo, "warning libvips already started")
		return
	}

	if C.VIPS_MAJOR_VERSION < 8 {
		panic("govips requires libvips version 8.10+")
	}

	if C.VIPS_MAJOR_VERSION == 8 && C.VIPS_MINOR_VERSION < 10 {
		panic("govips requires libvips version 8.10+")
	}

	cName := C.CString("govips")
	defer freeCString(cName)

	// Initialize govips logging handler and verbosity filter to historical default
	if !currentLoggingOverriden {
		LoggingSettings(nil, LogLevelInfo)
		currentLoggingOverriden = false
	}

	// Override default glib logging handler to intercept logging messages
	C.vips_set_logging_handler()

	err := C.vips_init(cName)
	if err != 0 {
		panic(fmt.Sprintf("Failed to start vips code=%v", err))
	}

	err = C.icc_profiles_init()
	if err != 0 {
		panic(fmt.Sprintf("Failed to initialize icc profiles=%v", err))
	}

	running = true

	if config != nil {
		if config.CollectStats {
			statCollectorDone = collectStats()
		}

		C.vips_leak_set(toGboolean(config.ReportLeaks))

		if config.ConcurrencyLevel >= 0 {
			C.vips_concurrency_set(C.int(config.ConcurrencyLevel))
		} else {
			C.vips_concurrency_set(defaultConcurrencyLevel)
		}

		if config.MaxCacheFiles >= 0 {
			C.vips_cache_set_max_files(C.int(config.MaxCacheFiles))
		} else {
			C.vips_cache_set_max_files(defaultMaxCacheFiles)
		}

		if config.MaxCacheMem >= 0 {
			C.vips_cache_set_max_mem(C.size_t(config.MaxCacheMem))
		} else {
			C.vips_cache_set_max_mem(defaultMaxCacheMem)
		}

		if config.MaxCacheSize >= 0 {
			C.vips_cache_set_max(C.int(config.MaxCacheSize))
		} else {
			C.vips_cache_set_max(defaultMaxCacheSize)
		}

		if config.CacheTrace {
			C.vips_cache_set_trace(toGboolean(true))
		}
	} else {
		C.vips_concurrency_set(defaultConcurrencyLevel)
		C.vips_cache_set_max(defaultMaxCacheSize)
		C.vips_cache_set_max_mem(defaultMaxCacheMem)
		C.vips_cache_set_max_files(defaultMaxCacheFiles)
	}

	govipsLog("govips", LogLevelInfo, fmt.Sprintf("vips %s started with concurrency=%d cache_max_files=%d cache_max_mem=%d cache_max=%d",
		Version,
		int(C.vips_concurrency_get()),
		int(C.vips_cache_get_max_files()),
		int(C.vips_cache_get_max_mem()),
		int(C.vips_cache_get_max())))

	initTypes()
}

// Shutdown libvips
func Shutdown() {
	if statCollectorDone != nil {
		statCollectorDone <- struct{}{}
	}

	initLock.Lock()
	runtime.LockOSThread()
	defer initLock.Unlock()
	defer runtime.UnlockOSThread()

	if !running {
		govipsLog("govips", LogLevelInfo, "warning libvips not started")
		return
	}

	C.vips_shutdown()
	running = false
}

// ShutdownThread clears the cache for for the given thread.
func ShutdownThread() {
	C.vips_thread_shutdown()
}

// ClearCache drops the whole operation cache, handy for leak tracking.
func ClearCache() {
	C.vips_cache_drop_all()
}

// PrintCache prints the whole operation cache to stdout for debugging purposes.
func PrintCache() {
	C.vips_cache_print()
}

// PrintObjectReport outputs all of the current internal objects in libvips
func PrintObjectReport(label string) {
	govipsLog("govips", LogLevelInfo, fmt.Sprintf("\n=======================================\nvips live objects: %s...\n", label))
	C.vips_object_print_all()
	govipsLog("govips", LogLevelInfo, "=======================================\n\n")
}

// MemoryStats is a data structure that houses various memory statistics from ReadVipsMemStats()
type MemoryStats struct {
	Mem     int64
	MemHigh int64
	Files   int64
	Allocs  int64
}

// ReadVipsMemStats returns various memory statistics such as allocated memory and open files.
func ReadVipsMemStats(stats *MemoryStats) {
	stats.Mem = int64(C.vips_tracked_get_mem())
	stats.MemHigh = int64(C.vips_tracked_get_mem_highwater())
	stats.Allocs = int64(C.vips_tracked_get_allocs())
	stats.Files = int64(C.vips_tracked_get_files())
}

func startupIfNeeded() {
	if !running {
		govipsLog("govips", LogLevelInfo, "libvips was forcibly started automatically, consider calling Startup/Shutdown yourself")
		Startup(nil)
	}
}

// InitTypes initializes caches and figures out which image types are supported
func initTypes() {
	once.Do(func() {
		cType := C.CString("VipsOperation")
		defer freeCString(cType)

		for k, v := range ImageTypes {
			name := strings.ToLower("VipsForeignLoad" + v)
			typeLoaders[name] = k
			typeLoaders[name+"buffer"] = k

			cFunc := C.CString(v + "load")
			//noinspection GoDeferInLoop
			defer freeCString(cFunc)

			ret := C.vips_type_find(cType, cFunc)

			supportedImageTypes[k] = int(ret) != 0

			govipsLog("govips", LogLevelInfo, fmt.Sprintf("registered image type loader type=%s", v))
		}
	})
}
