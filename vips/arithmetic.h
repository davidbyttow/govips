// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html

#include <stdlib.h>
#include <vips/vips.h>

int add(VipsImage *left, VipsImage *right, VipsImage **out);
int multiply(VipsImage *left, VipsImage *right, VipsImage **out);
int linear(VipsImage *in, VipsImage **out, double *a, double *b, int n);
int linear1(VipsImage *in, VipsImage **out, double a, double b);
int invert_image(VipsImage *in, VipsImage **out);
