#include "arithmetic.h"

int add(VipsImage *left, VipsImage *right, VipsImage **out) {
	return vips_add(left, right, out, NULL);
}

int multiply(VipsImage *left, VipsImage *right, VipsImage **out) {
	return vips_multiply(left, right, out, NULL);
}

int linear(VipsImage *in, VipsImage **out, double *a, double *b, int n) {
	return vips_linear(in, out, a, b, n, NULL);
}

int linear1(VipsImage *in, VipsImage **out, double a, double b) {
	return vips_linear1(in, out, a, b, NULL);
}

int invert_image(VipsImage *in, VipsImage **out) {
	return vips_invert(in, out, NULL);
}
