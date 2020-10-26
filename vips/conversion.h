// https://libvips.github.io/libvips/API/current/libvips-conversion.html

#include <stdlib.h>
#include <vips/vips.h>

int copy_image(VipsImage *in, VipsImage **out);

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, double r, double g, double b);

int flip_image(VipsImage *in, VipsImage **out, int direction);

int extract_image_area(VipsImage *in, VipsImage **out, int left, int top, int width, int height);

int extract_band(VipsImage *in, VipsImage **out, int band, int num);

int rot_image(VipsImage *in, VipsImage **out, VipsAngle angle);
int autorot_image(VipsImage *in, VipsImage **out);

int zoom_image(VipsImage *in, VipsImage **out, int xfac, int yfac);

int bandjoin(VipsImage **in, VipsImage **out, int n);
int flatten_image(VipsImage *in, VipsImage **out, double r, double g, double b);
int add_alpha(VipsImage *in, VipsImage **out);
int premultiply_alpha(VipsImage *in, VipsImage **out);
int unpremultiply_alpha(VipsImage *in, VipsImage **out);
int cast(VipsImage *in, VipsImage **out, int bandFormat);
double max_alpha(VipsImage *in);

int composite2_image(VipsImage *base, VipsImage *overlay, VipsImage **out, int mode, gint x, gint y);

int is_16bit(VipsInterpretation interpretation);
