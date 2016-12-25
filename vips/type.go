package vips

type ImageType int

const (
	ImageTypeUnknown ImageType = iota
	ImageTypeJpeg
	ImageTypeWebp
	ImageTypePng
	ImageTypeTiff
	ImageTypeGif
	ImageTypePdf
	ImageTypeSvg
)

func IsImageTypeSupported(imageType ImageType) bool {
	return true
}
