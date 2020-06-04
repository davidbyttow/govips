// https://libvips.github.io/libvips/API/current/libvips-resample.html

#include <stdlib.h>
#include <vips/vips.h>

int shrink_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink);
int reduce_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink);
int affine_image(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator);
int resize_image(VipsImage *in, VipsImage **out, double scale, gdouble vscale, int kernel);
int alpha_resize_image(VipsImage *in, VipsImage **out, double scale, gdouble vscale);
