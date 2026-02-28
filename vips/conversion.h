// https://libvips.github.io/libvips/API/current/libvips-conversion.html

#include <stdlib.h>
#include <vips/vips.h>

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width,
                int height, int extend);
int embed_image_background(VipsImage *in, VipsImage **out, int left, int top, int width,
                int height, double r, double g, double b, double a);
int embed_multi_page_image(VipsImage *in, VipsImage **out, int left, int top, int width,
                int height, int extend);
int embed_multi_page_image_background(VipsImage *in, VipsImage **out, int left, int top,
                int width, int height, double r, double g, double b, double a);

int flip_image(VipsImage *in, VipsImage **out, int direction);

int extract_image_area(VipsImage *in, VipsImage **out, int left, int top,
                       int width, int height);
int extract_area_multi_page(VipsImage *in, VipsImage **out, int left, int top,
                       int width, int height);

int smartcrop(VipsImage *in, VipsImage **out, int width, int height,
              int interesting);
int crop(VipsImage *in, VipsImage **out, int left, int top,
              int width, int height);

int similarity(VipsImage *in, VipsImage **out, double scale, double angle,
               double r, double g, double b, double a, double idx, double idy,
               double odx, double ody);
double max_alpha(VipsImage *in);

int composite_image(VipsImage **in, VipsImage **out, int n, int *mode, int *x,
                    int *y);

int join(VipsImage *in1, VipsImage *in2, VipsImage **out, int direction);

int is_16bit(VipsInterpretation interpretation);

int add_alpha(VipsImage *in, VipsImage **out);
