// https://libvips.github.io/libvips/API/current/VipsForeignSave.html

#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>


enum types {
	UNKNOWN = 0,
	JPEG,
	WEBP,
	PNG,
	TIFF,
	GIF,
	PDF,
	SVG,
	MAGICK,
	HEIF,
	BMP
};

int load_jpeg_buffer(void *buf, size_t len, VipsImage **out, int shrink, int fail, int autorotate);
int load_png_buffer(void *buf, size_t len, VipsImage **out);
int load_webp_buffer(void *buf, size_t len, VipsImage **out, int shrink);
int load_tiff_buffer(void *buf, size_t len, VipsImage **out, int page, int n, int autorotate, int subifd);
int load_gif_buffer(void *buf, size_t len, VipsImage **out, int page, int n);
int load_pdf_buffer(void *buf, size_t len, VipsImage **out, int page, int n, double dpi, double scale);
int load_svg_buffer(void *buf, size_t len, VipsImage **out, double dpi, double scale, int unlimited);
int load_heif_buffer(void *buf, size_t len, VipsImage **out, int page, int n, int thumbnail);
int load_magick_buffer(void *buf, size_t len, VipsImage **out, int page, int n, char *density);

int save_jpeg_buffer(VipsImage* image, void **buf, size_t *len, int strip, int quality, int interlace);
int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int interlace);
int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless, int effort);
int save_heif_buffer(VipsImage *in, void **buf, size_t *len, int quality, int lossless);
int save_tiff_buffer(VipsImage *in, void **buf, size_t *len);
