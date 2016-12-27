package vips

// #cgo pkg-config: vips
// #include "vips/vips.h"
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

var (
	once                sync.Once
	typeLoaders         = make(map[string]ImageType)
	supportedImageTypes = make(map[ImageType]bool)
)

func DetermineImageType(buf []byte) ImageType {
	InitTypes()

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

func InitTypes() {
	once.Do(func() {
		for k, v := range ImageTypes {
			name := strings.ToLower("VipsForeignLoad" + v)
			typeLoaders[name] = k
			typeLoaders[name+"buffer"] = k

			ret := C.vips_type_find(
				C.CString("VipsOperation"),
				C.CString(v+"load"))
			supportedImageTypes[k] = int(ret) != 0
		}
	})
}
