package vips

// #cgo pkg-config: vips
// #include "foreign.h"
import "C"
import (
	"bytes"
	"log"
	"math"
	"runtime"
	"unsafe"
)

// ImageType represents an image type
type ImageType int

// ImageType enum
const (
	ImageTypeUnknown ImageType = C.UNKNOWN
	ImageTypeGIF     ImageType = C.GIF
	ImageTypeJPEG    ImageType = C.JPEG
	ImageTypeMagick  ImageType = C.MAGICK
	ImageTypePDF     ImageType = C.PDF
	ImageTypePNG     ImageType = C.PNG
	ImageTypeSVG     ImageType = C.SVG
	ImageTypeTIFF    ImageType = C.TIFF
	ImageTypeWEBP    ImageType = C.WEBP
	ImageTypeHEIF    ImageType = C.HEIF
)

var imageTypeExtensionMap = map[ImageType]string{
	ImageTypeGIF:    ".gif",
	ImageTypeJPEG:   ".jpeg",
	ImageTypeMagick: ".magick",
	ImageTypePDF:    ".pdf",
	ImageTypePNG:    ".png",
	ImageTypeSVG:    ".svg",
	ImageTypeTIFF:   ".tiff",
	ImageTypeWEBP:   ".webp",
	ImageTypeHEIF:   ".heic",
}

var ImageTypes = map[ImageType]string{
	ImageTypeGIF:    "gif",
	ImageTypeJPEG:   "jpeg",
	ImageTypeMagick: "magick",
	ImageTypePDF:    "pdf",
	ImageTypePNG:    "png",
	ImageTypeSVG:    "svg",
	ImageTypeTIFF:   "tiff",
	ImageTypeWEBP:   "webp",
	ImageTypeHEIF:   "heif",
}

// FileExt returns the canonical extension for the ImageType
func (i ImageType) FileExt() string {
	if ext, ok := imageTypeExtensionMap[i]; ok {
		return ext
	}
	return ""
}

func IsTypeSupported(imageType ImageType) bool {
	startupIfNeeded()

	return supportedImageTypes[imageType]
}

// DetermineImageType attempts to determine the image type of the given buffer
func DetermineImageType(buf []byte) ImageType {
	if len(buf) < 12 {
		return ImageTypeUnknown
	}
	if isJPEG(buf) {
		return ImageTypeJPEG
	}
	if isPNG(buf) {
		return ImageTypePNG
	}
	if IsTypeSupported(ImageTypeGIF) && isGIF(buf) {
		return ImageTypeGIF
	}
	if IsTypeSupported(ImageTypeTIFF) && isTIFF(buf) {
		return ImageTypeTIFF
	}
	if IsTypeSupported(ImageTypeWEBP) && isWEBP(buf) {
		return ImageTypeWEBP
	}
	if IsTypeSupported(ImageTypeHEIF) && isHEIF(buf) {
		return ImageTypeHEIF
	}
	if IsTypeSupported(ImageTypeSVG) && isSVG(buf) {
		return ImageTypeSVG
	}
	if IsTypeSupported(ImageTypePDF) && isPDF(buf) {
		return ImageTypePDF
	}

	return ImageTypeUnknown
}

var jpeg = []byte("\xFF\xD8\xFF")

func isJPEG(buf []byte) bool {
	return bytes.HasPrefix(buf, jpeg)
}

var gif = []byte("\x47\x49\x46")

func isGIF(buf []byte) bool {
	return bytes.HasPrefix(buf, gif)
}

var png = []byte("\x89\x50\x4E\x47")

func isPNG(buf []byte) bool {
	return bytes.HasPrefix(buf, png)
}

var tifII = []byte("\x49\x49\x2A\x00")
var tifMM = []byte("\x4D\x4D\x00\x2A")

func isTIFF(buf []byte) bool {
	return bytes.HasPrefix(buf, tifII) || bytes.HasPrefix(buf, tifMM)
}

var webp = []byte("\x57\x45\x42\x50")

func isWEBP(buf []byte) bool {
	return bytes.Equal(buf[8:12], webp)
}

// https://github.com/strukturag/libheif/blob/master/libheif/heif.cc
var ftyp = []byte("ftyp")
var heic = []byte("heic")

func isHEIF(buf []byte) bool {
	return bytes.Equal(buf[4:8], ftyp) && bytes.Equal(buf[8:12], heic)
}

var svg = []byte("<svg ")

func isSVG(buf []byte) bool {
	sub := buf[:int(math.Min(500.0, float64(len(buf))))]
	return bytes.Contains(sub, svg)
}

var pdf = []byte("\x25\x50\x44\x46")

func isPDF(buf []byte) bool {
	return bytes.HasPrefix(buf, pdf)
}

func vipsLoadFromBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	// Reference buf here so it's not garbage collected during image initialization.
	defer runtime.KeepAlive(buf)

	var out *C.VipsImage

	imageType := DetermineImageType(buf)
	if imageType == ImageTypeUnknown {
		if len(buf) > 2 {
			log.Printf("Failed to understand image format size=%d %x %x %x", len(buf), buf[0], buf[1], buf[2])
		} else {
			log.Printf("Failed to understand image format size=%d", len(buf))
		}
		return nil, ImageTypeUnknown, ErrUnsupportedImageFormat
	}

	bufLength := C.size_t(len(buf))
	imageBuf := unsafe.Pointer(&buf[0])

	if err := C.load_image_buffer(imageBuf, bufLength, C.int(imageType), &out); err != 0 {
		return nil, ImageTypeUnknown, handleImageError(out)
	}

	return out, imageType, nil
}

func vipsSavePNGToBuffer(in *C.VipsImage, stripMetadata bool, compression, quality int, interlaced bool) ([]byte, error) {
	incOpCounter("save_png_buffer")
	var ptr unsafe.Pointer
	cLen := C.size_t(0)

	strip := C.int(boolToInt(stripMetadata))
	comp := C.int(compression)
	qual := C.int(quality)
	inter := C.int(boolToInt(interlaced))

	if err := C.save_png_buffer(in, &ptr, &cLen, strip, comp, qual, inter); err != 0 {
		return nil, handleSaveBufferError(ptr)
	}

	return toBuff(ptr, cLen), nil
}

func vipsSaveWebPToBuffer(in *C.VipsImage, stripMetadata bool, quality int, lossless bool) ([]byte, error) {
	incOpCounter("save_webp_buffer")
	var ptr unsafe.Pointer
	cLen := C.size_t(0)

	strip := C.int(boolToInt(stripMetadata))
	qual := C.int(quality)
	loss := C.int(boolToInt(lossless))

	if err := C.save_webp_buffer(in, &ptr, &cLen, strip, qual, loss); err != 0 {
		return nil, handleSaveBufferError(ptr)
	}

	return toBuff(ptr, cLen), nil
}

func vipsSaveTIFFToBuffer(in *C.VipsImage) ([]byte, error) {
	incOpCounter("save_tiff_buffer")
	var ptr unsafe.Pointer
	cLen := C.size_t(0)

	if err := C.save_tiff_buffer(in, &ptr, &cLen); err != 0 {
		return nil, handleSaveBufferError(ptr)
	}

	return toBuff(ptr, cLen), nil
}

func vipsSaveHEIFToBuffer(in *C.VipsImage, quality int, lossless bool) ([]byte, error) {
	incOpCounter("save_heif_buffer")
	var ptr unsafe.Pointer
	cLen := C.size_t(0)

	qual := C.int(quality)
	loss := C.int(boolToInt(lossless))

	if err := C.save_heif_buffer(in, &ptr, &cLen, qual, loss); err != 0 {
		return nil, handleSaveBufferError(ptr)
	}

	return toBuff(ptr, cLen), nil
}

func vipsSaveJPEGToBuffer(in *C.VipsImage, quality int, stripMetadata, interlaced bool) ([]byte, error) {
	incOpCounter("save_jpeg_buffer")
	var ptr unsafe.Pointer
	cLen := C.size_t(0)

	strip := C.int(boolToInt(stripMetadata))
	qual := C.int(quality)
	inter := C.int(boolToInt(interlaced))

	if err := C.save_jpeg_buffer(in, &ptr, &cLen, strip, qual, inter); err != 0 {
		return nil, handleSaveBufferError(ptr)
	}

	return toBuff(ptr, cLen), nil
}

func toBuff(ptr unsafe.Pointer, cLen C.size_t) []byte {
	buf := C.GoBytes(ptr, C.int(cLen))
	gFreePointer(ptr)

	return buf
}
