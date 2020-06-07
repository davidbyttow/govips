#include "resample.h"

int shrink_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink) {
	return vips_shrink(in, out, xshrink, yshrink, NULL);
}

int reduce_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink) {
	return vips_reduce(in, out, xshrink, yshrink, NULL);
}

int affine_image(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator) {
	return vips_affine(in, out, a, b, c, d, "interpolate", interpolator, NULL);
}

int resize_image(VipsImage *in, VipsImage **out, double scale, gdouble vscale, int kernel) {
	if (vscale > 0) {
		return vips_resize(in, out, scale, "vscale", vscale, "kernel", kernel, NULL);
	}

	return vips_resize(in, out, scale, "kernel", kernel, NULL);
}

int alpha_resize_image(VipsImage *in, VipsImage **out, double scale, gdouble vscale) {
	if (vscale > 0) {
		return vips_alpha_resize(in, out, scale, "vscale", vscale, NULL);
	}

	return vips_alpha_resize(in, out, scale, NULL);
}
