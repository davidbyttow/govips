// https://libvips.github.io/libvips/API/current/libvips-arithmetic.html

#include <vips/vips.h>

int find_trim(VipsImage *in, int *left, int *top, int *width, int *height,
              double threshold, double r, double g, double b);
int getpoint(VipsImage *in, double **vector, int n, int x, int y);
int minOp(VipsImage *in, double *out, int *x, int *y, int size);
