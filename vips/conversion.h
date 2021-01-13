// https://libvips.github.io/libvips/API/current/libvips-conversion.html

#include <stdlib.h>
#include <vips/vips.h>

int copy_image(VipsImage *in, VipsImage **out);

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width,
                int height, int extend);

int flip_image(VipsImage *in, VipsImage **out, int direction);

int extract_image_area(VipsImage *in, VipsImage **out, int left, int top,
                       int width, int height);

int extract_band(VipsImage *in, VipsImage **out, int band, int num);

int rot_image(VipsImage *in, VipsImage **out, VipsAngle angle);
int autorot_image(VipsImage *in, VipsImage **out);

int zoom_image(VipsImage *in, VipsImage **out, int xfac, int yfac);
int smartcrop(VipsImage *in, VipsImage **out, int width, int height,
              int interesting);

int bandjoin(VipsImage **in, VipsImage **out, int n);
int bandjoin_const(VipsImage *in, VipsImage **out, double constants[], int n);
int similarity(VipsImage *in, VipsImage **out, double scale, double angle,
               double r, double g, double b, double a, double idx, double idy,
               double odx, double ody);
int flatten_image(VipsImage *in, VipsImage **out, double r, double g, double b);
int add_alpha(VipsImage *in, VipsImage **out);
int premultiply_alpha(VipsImage *in, VipsImage **out);
int unpremultiply_alpha(VipsImage *in, VipsImage **out);
int cast(VipsImage *in, VipsImage **out, int bandFormat);
double max_alpha(VipsImage *in);

int composite_image(VipsImage **in, VipsImage **out, int n, int *mode, int *x,
                    int *y);
int composite2_image(VipsImage *base, VipsImage *overlay, VipsImage **out,
                     int mode, gint x, gint y);

int is_16bit(VipsInterpretation interpretation);
