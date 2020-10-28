#include "lang.h"
#include "foreign.h"

int load_jpeg_buffer(void *buf, size_t len, VipsImage **out, int shrink, int fail, int autorotate) {
	return vips_jpegload_buffer(buf, len, out,
		"shrink", shrink,
		"fail", INT_TO_GBOOLEAN(fail),
		"autorotate", INT_TO_GBOOLEAN(autorotate),
		NULL);
}

int load_png_buffer(void *buf, size_t len, VipsImage **out) {
	return vips_pngload_buffer(buf, len, out, NULL);
}

int load_webp_buffer(void *buf, size_t len, VipsImage **out, int shrink) {
	return vips_webpload_buffer(buf, len, out,
		"shrink", shrink,
		NULL);
}

int load_tiff_buffer(void *buf, size_t len, VipsImage **out, int page, int n, int autorotate, int subifd) {
	return vips_tiffload_buffer(buf, len, out,
		"page", page,
		"n", n,
		"autorotate", INT_TO_GBOOLEAN(autorotate),
		"subifd", subifd,
		NULL);
}

int load_gif_buffer(void *buf, size_t len, VipsImage **out, int page, int n) {
	return vips_tiffload_buffer(buf, len, out,
		"page", page,
		"n", n,
		NULL);
}

int load_pdf_buffer(void *buf, size_t len, VipsImage **out, int page, int n, double dpi, double scale) {
	return vips_pdfload_buffer(buf, len, out,
		"page", page,
		"n", n,
		"dpi", dpi,
		"scale", scale,
		NULL);
}

int load_svg_buffer(void *buf, size_t len, VipsImage **out, double dpi, double scale, int unlimited) {
	return vips_svgload_buffer(buf, len, out,
		"dpi", dpi,
		"scale", scale,
		"unlimited", INT_TO_GBOOLEAN(unlimited),
		NULL);
}

int load_heif_buffer(void *buf, size_t len, VipsImage **out, int page, int n, int thumbnail) {
	return vips_heifload_buffer(buf, len, out,
		"page", page,
		"n", n,
		"thumbnail", INT_TO_GBOOLEAN(thumbnail),
		NULL);
}

int load_magick_buffer(void *buf, size_t len, VipsImage **out, int page, int n, char *density) {
	return vips_magickload_buffer(buf, len, out,
		"page", page,
		"n", n,
		"density", density,
		NULL);
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
