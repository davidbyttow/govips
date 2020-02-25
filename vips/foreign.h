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

int load_image_buffer(void *buf, size_t len, int imageType, VipsImage **out);

int find_image_type_loader(int t);
int find_image_type_saver(int t);

int save_jpeg_buffer(VipsImage* image, void **buf, size_t *len, int strip, int quality, int interlace);
int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int interlace);
int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless);
int save_heif_buffer(VipsImage *in, void **buf, size_t *len, int quality, int lossless);
int save_tiff_buffer(VipsImage *in, void **buf, size_t *len);
