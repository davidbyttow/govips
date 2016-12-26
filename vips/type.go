package vips

// #cgo pkg-config: vips
// #include "bridge.h"
import "C"
import (
	"strings"
	"sync"
)

type ImageType int

const (
	ImageTypeUnknown ImageType = iota
	ImageTypeGif
	ImageTypeJpeg
	ImageTypeMagick
	ImageTypePdf
	ImageTypePng
	ImageTypeSvg
	ImageTypeTiff
	ImageTypeWebp
)

var ImageTypes = map[ImageType]string{
	ImageTypeGif:    "gif",
	ImageTypeJpeg:   "jpeg",
	ImageTypeMagick: "magick",
	ImageTypePdf:    "pdf",
	ImageTypePng:    "png",
	ImageTypeSvg:    "svg",
	ImageTypeTiff:   "tiff",
	ImageTypeWebp:   "webp",
}

func DetermineImageType(buf []byte) ImageType {
	discoverSupportedImageTypes()

	size := len(buf)
	if size == 0 {
		return ImageTypeUnknown
	}

	cname := C.vips_foreign_find_load_buffer(
		cPtr(buf),
		C.size_t(size))

	if cname == nil {
		return ImageTypeUnknown
	}

	imageType := ImageTypeUnknown
	name := strings.ToLower(C.GoString(cname))
	if imageType, ok := typeLoaders[name]; ok {
		return imageType
	}

	return imageType
}

var once sync.Once

func discoverSupportedImageTypes() error {
	var err error
	once.Do(func() {
		for k, v := range ImageTypes {
			ret := C.vips_type_find(
				C.CString("VipsOperation"),
				C.CString(v+"load"))
			supportedImageTypes[k] = int(ret) != 0
		}
	})
	return err
}

var typeLoaders = make(map[string]ImageType)
var supportedImageTypes = make(map[ImageType]bool)

func initTypes() {
	for k, v := range ImageTypes {
		name := strings.ToLower("VipsForeignLoad" + v)
		typeLoaders[name] = k
		typeLoaders[name+"buffer"] = k
	}

	if err := discoverSupportedImageTypes(); err != nil {
		panic(err)
	}
}
