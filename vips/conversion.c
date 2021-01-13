#include "conversion.h"

int copy_image(VipsImage *in, VipsImage **out) {
  return vips_copy(in, out, NULL);
}

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width,
                int height, int extend) {
  return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
}

int flip_image(VipsImage *in, VipsImage **out, int direction) {
  return vips_flip(in, out, direction, NULL);
}

int extract_image_area(VipsImage *in, VipsImage **out, int left, int top,
                       int width, int height) {
  return vips_extract_area(in, out, left, top, width, height, NULL);
}

int extract_band(VipsImage *in, VipsImage **out, int band, int num) {
  if (num > 0) {
    return vips_extract_band(in, out, band, "n", num, NULL);
  }
  return vips_extract_band(in, out, band, NULL);
}

int rot_image(VipsImage *in, VipsImage **out, VipsAngle angle) {
  return vips_rot(in, out, angle, NULL);
}

int autorot_image(VipsImage *in, VipsImage **out) {
  return vips_autorot(in, out, NULL);
}

int zoom_image(VipsImage *in, VipsImage **out, int xfac, int yfac) {
  return vips_zoom(in, out, xfac, yfac, NULL);
}

int bandjoin(VipsImage **in, VipsImage **out, int n) {
  return vips_bandjoin(in, out, n, NULL);
}

int bandjoin_const(VipsImage *in, VipsImage **out, double constants[], int n) {
  return vips_bandjoin_const(in, out, constants, n, NULL);
}

int similarity(VipsImage *in, VipsImage **out, double scale, double angle,
               double r, double g, double b, double a, double idx, double idy,
               double odx, double ody) {
  if (is_16bit(in->Type)) {
    r = 65535 * r / 255;
    g = 65535 * g / 255;
    b = 65535 * b / 255;
    a = 65535 * a / 255;
  }

  double background[3] = {r, g, b};
  double backgroundRGBA[4] = {r, g, b, a};

  VipsArrayDouble *vipsBackground;

  // Ignore the alpha channel if the image doesn't have one
  if (in->Bands <= 3) {
    vipsBackground = vips_array_double_new(background, 3);
  } else {
    vipsBackground = vips_array_double_new(backgroundRGBA, 4);
  }

  int code = vips_similarity(in, out, "scale", scale, "angle", angle,
                             "background", vipsBackground, "idx", idx, "idy",
                             idy, "odx", odx, "ody", ody, NULL);

  vips_area_unref(VIPS_AREA(vipsBackground));
  return code;
}

int smartcrop(VipsImage *in, VipsImage **out, int width, int height,
              int interesting) {
  return vips_smartcrop(in, out, width, height, "interesting", interesting,
                        NULL);
}

int flatten_image(VipsImage *in, VipsImage **out, double r, double g,
                  double b) {
  if (is_16bit(in->Type)) {
    r = 65535 * r / 255;
    g = 65535 * g / 255;
    b = 65535 * b / 255;
  }

  double background[3] = {r, g, b};
  VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);

  int code = vips_flatten(in, out, "background", vipsBackground, "max_alpha",
                          is_16bit(in->Type) ? 65535.0 : 255.0, NULL);

  vips_area_unref(VIPS_AREA(vipsBackground));
  return code;
}

int is_16bit(VipsInterpretation interpretation) {
  return interpretation == VIPS_INTERPRETATION_RGB16 ||
         interpretation == VIPS_INTERPRETATION_GREY16;
}

int add_alpha(VipsImage *in, VipsImage **out) {
  return vips_addalpha(in, out, NULL);
}

int premultiply_alpha(VipsImage *in, VipsImage **out) {
  return vips_premultiply(in, out, "max_alpha", max_alpha(in), NULL);
}

int unpremultiply_alpha(VipsImage *in, VipsImage **out) {
  return vips_unpremultiply(in, out, NULL);
}

int cast(VipsImage *in, VipsImage **out, int bandFormat) {
  return vips_cast(in, out, bandFormat, NULL);
}

double max_alpha(VipsImage *in) {
  switch (in->BandFmt) {
    case VIPS_FORMAT_USHORT:
      return 65535;
    case VIPS_FORMAT_FLOAT:
    case VIPS_FORMAT_DOUBLE:
      return 1.0;
    default:
      return 255;
  }
}

int composite_image(VipsImage **in, VipsImage **out, int n, int *mode, int *x,
                    int *y) {
  VipsArrayInt *xs = vips_array_int_new(x, n - 1);
  VipsArrayInt *ys = vips_array_int_new(y, n - 1);

  int code = vips_composite(in, out, n, mode, "x", xs, "y", ys, NULL);

  vips_area_unref(VIPS_AREA(xs));
  vips_area_unref(VIPS_AREA(ys));
  return code;
}

int composite2_image(VipsImage *base, VipsImage *overlay, VipsImage **out,
                     int mode, gint x, gint y) {
  return vips_composite2(base, overlay, out, mode, "x", x, "y", y, NULL);
}
