package vips

func determineImageType(buf []byte) ImageType {
	if len(buf) == 0 {
		return ImageTypeUnknown
	}
	if buf[0] == 0x89 && buf[1] == 0x50 && buf[2] == 0x4E && buf[3] == 0x47 {
		return ImageTypePng
	}
	if buf[0] == 0xFF && buf[1] == 0xD8 && buf[2] == 0xFF {
		return ImageTypeJpeg
	}
	if IsImageTypeSupported(ImageTypeWebp) && buf[8] == 0x57 && buf[9] == 0x45 && buf[10] == 0x42 && buf[11] == 0x50 {
		return ImageTypeWebp
	}
	if IsImageTypeSupported(ImageTypeTiff) &&
		((buf[0] == 0x49 && buf[1] == 0x49 && buf[2] == 0x2A && buf[3] == 0x0) ||
			(buf[0] == 0x4D && buf[1] == 0x4D && buf[2] == 0x0 && buf[3] == 0x2A)) {
		return ImageTypeTiff
	}
	if IsImageTypeSupported(ImageTypeGif) && buf[0] == 0x47 && buf[1] == 0x49 && buf[2] == 0x46 {
		return ImageTypeGif
	}
	if IsImageTypeSupported(ImageTypePdf) && buf[0] == 0x25 && buf[1] == 0x50 && buf[2] == 0x44 && buf[3] == 0x46 {
		return ImageTypePdf
	}
	// if IsImageTypeSupported(Svg) && IsSVGImage(buf) {
	// 	return Svg
	// }
	return ImageTypeUnknown
}
