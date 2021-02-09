package vips

// #include "foreign.h"
import "C"
import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/png"
	"math"
	"runtime"
	"unsafe"

	"golang.org/x/image/bmp"
	"golang.org/x/net/html/charset"
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
	ImageTypeBMP     ImageType = C.BMP
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
	ImageTypeBMP:    ".bmp",
}

// ImageTypes defines the various image types supported by govips
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
	ImageTypeBMP:    "bmp",
}

// TiffCompression represents method for compressing a tiff at export
type TiffCompression int

// TiffCompression enum
const (
	TiffCompressionNone     TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_NONE
	TiffCompressionJpeg     TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_JPEG
	TiffCompressionDeflate  TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_DEFLATE
	TiffCompressionPackbits TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_PACKBITS
	TiffCompressionFax4     TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_CCITTFAX4
	TiffCompressionLzw      TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_LZW
	TiffCompressionWebp     TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_WEBP
	TiffCompressionZstd     TiffCompression = C.VIPS_FOREIGN_TIFF_COMPRESSION_ZSTD
)

// TiffPredictor represents method for compressing a tiff at export
type TiffPredictor int

// TiffPredictor enum
const (
	TiffPredictorNone       TiffPredictor = C.VIPS_FOREIGN_TIFF_PREDICTOR_NONE
	TiffPredictorHorizontal TiffPredictor = C.VIPS_FOREIGN_TIFF_PREDICTOR_HORIZONTAL
	TiffPredictorFloat      TiffPredictor = C.VIPS_FOREIGN_TIFF_PREDICTOR_FLOAT
)

// FileExt returns the canonical extension for the ImageType
func (i ImageType) FileExt() string {
	if ext, ok := imageTypeExtensionMap[i]; ok {
		return ext
	}
	return ""
}

// IsTypeSupported checks whether given image type is supported by govips
func IsTypeSupported(imageType ImageType) bool {
	startupIfNeeded()

	return supportedImageTypes[imageType]
}

// DetermineImageType attempts to determine the image type of the given buffer
func DetermineImageType(buf []byte) ImageType {
	if len(buf) < 12 {
		return ImageTypeUnknown
	} else if isJPEG(buf) {
		return ImageTypeJPEG
	} else if isPNG(buf) {
		return ImageTypePNG
	} else if isGIF(buf) {
		return ImageTypeGIF
	} else if isTIFF(buf) {
		return ImageTypeTIFF
	} else if isWEBP(buf) {
		return ImageTypeWEBP
	} else if isHEIF(buf) {
		return ImageTypeHEIF
	} else if isSVG(buf) {
		return ImageTypeSVG
	} else if isPDF(buf) {
		return ImageTypePDF
	} else if isBMP(buf) {
		return ImageTypeBMP
	} else {
		return ImageTypeUnknown
	}
}

var jpeg = []byte("\xFF\xD8\xFF")

func isJPEG(buf []byte) bool {
	return bytes.HasPrefix(buf, jpeg)
}

var gifHeader = []byte("\x47\x49\x46")

func isGIF(buf []byte) bool {
	return bytes.HasPrefix(buf, gifHeader)
}

var pngHeader = []byte("\x89\x50\x4E\x47")

func isPNG(buf []byte) bool {
	return bytes.HasPrefix(buf, pngHeader)
}

var tifII = []byte("\x49\x49\x2A\x00")
var tifMM = []byte("\x4D\x4D\x00\x2A")

func isTIFF(buf []byte) bool {
	return bytes.HasPrefix(buf, tifII) || bytes.HasPrefix(buf, tifMM)
}

var webpHeader = []byte("\x57\x45\x42\x50")

func isWEBP(buf []byte) bool {
	return bytes.Equal(buf[8:12], webpHeader)
}

// https://github.com/strukturag/libheif/blob/master/libheif/heif.cc
var ftyp = []byte("ftyp")
var heic = []byte("heic")
var mif1 = []byte("mif1")
var msf1 = []byte("msf1")
var avif = []byte("avif")

func isHEIF(buf []byte) bool {
	return bytes.Equal(buf[4:8], ftyp) && (bytes.Equal(buf[8:12], heic) ||
		bytes.Equal(buf[8:12], avif) ||
		bytes.Equal(buf[8:12], mif1) ||
		bytes.Equal(buf[8:12], msf1))
}

var svg = []byte("<svg")

func isSVG(buf []byte) bool {
	sub := buf[:int(math.Min(1024.0, float64(len(buf))))]
	if bytes.Contains(sub, svg) {
		data := &struct {
			XMLName xml.Name `xml:"svg"`
		}{}
		reader := bytes.NewReader(buf)
		decoder := xml.NewDecoder(reader)
		decoder.Strict = false
		decoder.CharsetReader = charset.NewReaderLabel

		err := decoder.Decode(data)

		return err == nil && data.XMLName.Local == "svg"
	}

	return false
}

var pdf = []byte("\x25\x50\x44\x46")

func isPDF(buf []byte) bool {
	return bytes.HasPrefix(buf, pdf)
}

var bmpHeader = []byte("BM")

func isBMP(buf []byte) bool {
	return bytes.HasPrefix(buf, bmpHeader)
}

func vipsLoadFromBuffer(buf []byte) (*C.VipsImage, ImageType, error) {
	src := buf
	// Reference src here so it's not garbage collected during image initialization.
	defer runtime.KeepAlive(src)

	var err error
	var out *C.VipsImage

	imageType := DetermineImageType(src)

	if imageType == ImageTypeBMP {
		src, err = bmpToPNG(src)
		if err != nil {
			return nil, ImageTypeUnknown, err
		}

		imageType = ImageTypePNG
	}

	if !IsTypeSupported(imageType) {
		govipsLog("govips", LogLevelInfo, fmt.Sprintf("failed to understand image format size=%d", len(src)))
		return nil, ImageTypeUnknown, ErrUnsupportedImageFormat
	}

	if err := C.load_image_buffer(unsafe.Pointer(&src[0]), C.size_t(len(src)), C.int(imageType), &out); err != 0 {
		return nil, ImageTypeUnknown, handleImageError(out)
	}

	return out, imageType, nil
}

func bmpToPNG(src []byte) ([]byte, error) {
	i, err := bmp.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	var w bytes.Buffer
	err = png.Encode(&w, i)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func vipsSaveJPEGToBuffer(in *C.VipsImage, params JpegExportParams) ([]byte, error) {
	incOpCounter("save_jpeg_buffer")

	p := C.create_save_params(C.JPEG)
	p.inputImage = in
	p.stripMetadata = C.int(boolToInt(params.StripMetadata))
	p.quality = C.int(params.Quality)
	p.interlace = C.int(boolToInt(params.Interlace))
	p.jpegOptimizeCoding = C.int(boolToInt(params.OptimizeCoding))
	p.jpegSubsample = C.VipsForeignJpegSubsample(params.Subsampling)

	return vipsSaveToBuffer(p)
}

func vipsSavePNGToBuffer(in *C.VipsImage, params PngExportParams) ([]byte, error) {
	incOpCounter("save_png_buffer")

	p := C.create_save_params(C.PNG)
	p.inputImage = in
	p.stripMetadata = C.int(boolToInt(params.StripMetadata))
	p.interlace = C.int(boolToInt(params.Interlace))
	p.pngCompression = C.int(params.Compression)

	return vipsSaveToBuffer(p)
}

func vipsSaveWebPToBuffer(in *C.VipsImage, params WebpExportParams) ([]byte, error) {
	incOpCounter("save_webp_buffer")

	p := C.create_save_params(C.WEBP)
	p.inputImage = in
	p.stripMetadata = C.int(boolToInt(params.StripMetadata))
	p.quality = C.int(params.Quality)
	p.webpLossless = C.int(boolToInt(params.Lossless))
	p.webpNearLossless = C.int(boolToInt(params.NearLossless))
	p.webpReductionEffort = C.int(params.ReductionEffort)

	if params.IccProfile != "" {
		p.webpIccProfile = C.CString(params.IccProfile)
		defer C.free(unsafe.Pointer(p.webpIccProfile))
	}

	return vipsSaveToBuffer(p)
}

func vipsSaveTIFFToBuffer(in *C.VipsImage, params TiffExportParams) ([]byte, error) {
	incOpCounter("save_tiff_buffer")

	p := C.create_save_params(C.TIFF)
	p.inputImage = in
	p.stripMetadata = C.int(boolToInt(params.StripMetadata))
	p.quality = C.int(params.Quality)
	p.tiffCompression = C.VipsForeignTiffCompression(params.Compression)

	return vipsSaveToBuffer(p)
}

func vipsSaveHEIFToBuffer(in *C.VipsImage, params HeifExportParams) ([]byte, error) {
	incOpCounter("save_heif_buffer")

	p := C.create_save_params(C.HEIF)
	p.inputImage = in
	p.outputFormat = C.HEIF
	p.quality = C.int(params.Quality)
	p.heifLossless = C.int(boolToInt(params.Lossless))

	return vipsSaveToBuffer(p)
}

func vipsSaveToBuffer(params C.struct_SaveParams) ([]byte, error) {
	if err := C.save_to_buffer(&params); err != 0 {
		return nil, handleSaveBufferError(params.outputBuffer)
	}

	buf := C.GoBytes(params.outputBuffer, C.int(params.outputLen))
	defer gFreePointer(params.outputBuffer)

	return buf, nil
}
