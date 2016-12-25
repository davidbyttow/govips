package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

// TODO(d): Tune these. Concurrency is set to a safe level but assumes
// openslide is not enabled.
const (
	concurrencyLevel = 25
	maxCacheMem      = 100 * 1024 * 1024
	maxCacheSize     = 500
)

const VipsVersion = string(C.VIPS_VERSION)
const VipsMajorVersion = int(C.VIPS_MAJOR_VERSION)
const VipsMinorVersion = int(C.VIPS_MINOR_VERSION)

func handleVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}

func FinalizeRequest() {
	C.vips_thread_shutdown()
}

var (
	lock sync.Mutex
)

func init() {
	if C.VIPS_MAJOR_VERSION < 8 {
		panic("Requires libvips version 8+")
	}

	lock.Lock()
	runtime.LockOSThread()
	defer lock.Unlock()
	defer runtime.UnlockOSThread()

	err := C.vips_init(C.CString("gimage"))
	if err != 0 {
		panic(fmt.Sprintf("Failed to start vips code=%d", err))
	}

	C.vips_cache_set_max_mem(maxCacheMem)
	C.vips_cache_set_max(maxCacheSize)
	C.vips_concurrency_set(concurrencyLevel)
}
