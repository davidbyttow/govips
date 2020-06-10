#include "lang.h"
#include "foreign.h"

int load_image_buffer(void *buf, size_t len, int imageType, VipsImage **out) {
	int code = 1;

	if (imageType == JPEG) {
		code = vips_jpegload_buffer(buf, len, out, NULL);
	} else if (imageType == PNG) {
		code = vips_pngload_buffer(buf, len, out, NULL);
	} else if (imageType == WEBP) {
		code = vips_webpload_buffer(buf, len, out, NULL);
	} else if (imageType == TIFF) {
		code = vips_tiffload_buffer(buf, len, out, NULL);
	} else if (imageType == GIF) {
		code = vips_gifload_buffer(buf, len, out, NULL);
	} else if (imageType == PDF) {
		code = vips_pdfload_buffer(buf, len, out, NULL);
	} else if (imageType == SVG) {
		code = vips_svgload_buffer(buf, len, out, NULL);
	} else if (imageType == HEIF) {
		code = vips_heifload_buffer(buf, len, out, NULL);
	} else if (imageType == MAGICK) {
		code = vips_magickload_buffer(buf, len, out, NULL);
	}

	return code;
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-jpegsave-buffer
int save_jpeg_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace) {
    return vips_jpegsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"optimize_coding", TRUE,
		"interlace", INT_TO_GBOOLEAN(interlace),
		"subsample_mode", VIPS_FOREIGN_JPEG_SUBSAMPLE_ON,
		NULL
	);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-pngsave-buffer
int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int interlace) {
	return vips_pngsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"compression", compression,
		"interlace", INT_TO_GBOOLEAN(interlace),
		"filter", VIPS_FOREIGN_PNG_FILTER_NONE,
		NULL
	);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/webpsave.c#L524
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-webpsave-buffer
int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless, int effort) {
	return vips_webpsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		"reduction_effort", effort,
		NULL
	);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/heifsave.c#L653
int save_heif_buffer(VipsImage *in, void **buf, size_t *len, int quality, int lossless) {
	return vips_heifsave_buffer(in, buf, len,
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
		NULL
	);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-tiffsave-buffer
int save_tiff_buffer(VipsImage *in, void **buf, size_t *len) {
	return vips_tiffsave_buffer(in, buf, len, NULL);
}
