// Package govips provides a library for transforming images that is built on lipvips. Libvips
// is an extremely fast C-library. Therefore, govips requires that libvips 8+ be installed
// and available on the target environment.
package govips

//go:generate scripts/codegen.sh

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"fmt"
	"sync"
)

// TODO(d): Tune these. Concurrency is set to a safe level but assumes
// openslide is not enabled.
const (
	concurrencyLevel = 25
	maxCacheMem      = 100 * 1024 * 1024
	maxCacheSize     = 500
)

// VipsVersion if the primary version of libvips
const VipsVersion = string(C.VIPS_VERSION)

// VipsMajorVersion is the major version of libvips
const VipsMajorVersion = int(C.VIPS_MAJOR_VERSION)

// VipsMinorVersion if the minor vesrion of libvips
const VipsMinorVersion = int(C.VIPS_MINOR_VERSION)

var (
	running  = false
	initLock sync.Mutex
)

// Startup sets up the libvips support and ensures the versions are correct
func Startup() {
	initLock.Lock()
	defer initLock.Unlock()

	if running {
		panic("libvips already running")
	}

	if C.VIPS_MAJOR_VERSION < 8 {
		panic("Requires libvips version 8+")
	}

	cName := C.CString("govips")
	defer freeCString(cName)

	err := C.vips_init(cName)
	if err != 0 {
		panic(fmt.Sprintf("Failed to start vips code=%d", err))
	}

	running = true

	C.vips_leak_set(toGboolean(true))
	C.vips_cache_set_max_mem(maxCacheMem)
	C.vips_cache_set_max(maxCacheSize)
	C.vips_concurrency_set(concurrencyLevel)

	initTypes()
}

func startupIfNeeded() {
	if !running {
		debug("libvips was forcibly started automatically, consider calling Startup/Shutdown yourself")
		Startup()
	}
}

// Shutdown libvips
func Shutdown() {
	initLock.Lock()
	defer initLock.Unlock()

	if !running {
		panic("libvips not running")
	}

	C.vips_shutdown()
	running = false
}

// ShutdownThread clears the cache for for the given thread
func ShutdownThread() {
	C.vips_thread_shutdown()
}
