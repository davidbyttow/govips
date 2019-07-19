#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

#if (VIPS_MAJOR_VERSION < 8)
  error_requires_version_8
#endif

#define INT_TO_GBOOLEAN(bool) (bool > 0 ? TRUE : FALSE)

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
	HEIF
};
typedef struct {
  const char *Text;
	const char *Font;
	int Width;
  int Height;
  int OffsetX;
  int OffsetY;
  VipsAlign Align;
	int DPI;
  int Margin;
	float Opacity;
  double Color[3];
} LabelOptions;

typedef struct {
	int access;
} ImageLoadOptions;

int init_image(void *buf, size_t len, int imageType, ImageLoadOptions *o, VipsImage **out);
int find_image_type_loader(int t);
int find_image_type_saver(int t);

//int load_jpeg_buffer(void *buf, size_t len, VipsImage **out, int shrink);

int save_jpeg_buffer(VipsImage* image, void **buf, size_t *len, int strip, int quality, int interlace);
int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int quality, int interlace);
int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless);
int save_heif_buffer(VipsImage *in, void **buf, size_t *len, int quality, int lossless);
int save_tiff_buffer(VipsImage *in, void **buf, size_t *len);

int copy_image(VipsImage *in, VipsImage **out);

int to_colorspace(VipsImage *in, VipsImage **out, VipsInterpretation space);
int is_colorspace_supported(VipsImage *in);
unsigned long has_profile_embed(VipsImage *in);
int remove_icc_profile(VipsImage *in);

int extract_band(VipsImage *in, VipsImage **out, int band, int num);
int bandjoin(VipsImage **in, VipsImage **out, int n);

int linear1(VipsImage *in, VipsImage **out, double a, double b);

// Operations
int add(VipsImage *left, VipsImage *right, VipsImage **out);
int multiply(VipsImage *left, VipsImage *right, VipsImage **out);
int resize_image(VipsImage *in, VipsImage **out, double scale, double vscale, int kernel);
int flip_image(VipsImage *in, VipsImage **out, int direction);
int shrink_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink);
int reduce_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink);
int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, double r, double g, double b);
int extract_image_area(VipsImage *in, VipsImage **out, int left, int top, int width, int height);
int flatten_image_background(VipsImage *in, VipsImage **out, double r, double g, double b);
int transform_image(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator);
int invert_image(VipsImage *in, VipsImage **out);
int composite(VipsImage **in, VipsImage **out, int n, int mode);
int label(VipsImage *in, VipsImage **out, LabelOptions *o);

int has_alpha_channel(VipsImage *image);
int add_alpha(VipsImage *in, VipsImage **out);

// conversion
// https://libvips.github.io/libvips/API/current/libvips-conversion.html
int autorot_image(VipsImage *in, VipsImage **out);
int rot_image(VipsImage *in, VipsImage **out, VipsAngle angle);
int zoom_image(VipsImage *in, VipsImage **out, int xfac, int yfac);

// Convolutions
// https://libvips.github.io/libvips/API/current/libvips-convolution.html
int gaussian_blur(VipsImage *in, VipsImage **out, double sigma);
int sharpen(VipsImage *in, VipsImage **out, double sigma, double x1, double m2);


void gobject_set_property(VipsObject* object, const char* name, const GValue* value);
