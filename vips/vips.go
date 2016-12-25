package vips

// #cgo pkg-config: vips
// #include "vips.h"
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
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

func loadImage(buf []byte) (*C.VipsImage, ImageType, error) {
	var image *C.VipsImage
	imageType := determineImageType(buf)

	if imageType == ImageTypeUnknown {
		return nil, ImageTypeUnknown, errors.New("Unsupported image format")
	}

	imageBuf := unsafe.Pointer(&buf[0])
	length := C.size_t(len(buf))

	err := C.init_image(imageBuf, length, C.int(imageType), &image)
	if err != 0 {
		return nil, ImageTypeUnknown, handleVipsError()
	}

	return image, imageType, nil
}

func loadJpegImage(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}

func loadWebpImage(buf []byte, shrinkFactor int) (*C.VipsImage, error) {
	return nil, nil
}

func handleVipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}

func FinalizeRequest() {
	C.finalize_request()
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
