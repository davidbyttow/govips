#include "color.h"

#include <unistd.h>

int is_colorspace_supported(VipsImage *in) {
  return vips_colourspace_issupported(in) ? 1 : 0;
}

int to_colorspace(VipsImage *in, VipsImage **out, VipsInterpretation space) {
  return vips_colourspace(in, out, space, NULL);
}

// https://libvips.github.io/libvips/API/8.6/libvips-colour.html#vips-icc-transform
int optimize_icc_profile(VipsImage *in, VipsImage **out, int isCmyk,
                         char *srgb_profile_path, char *gray_profile_path) {
  // todo: check current embedded profile, and skip if already set

  int channels = vips_image_get_bands(in);
  int result;

  if (vips_icc_present() == 0) {
    return 1;
  }

  if (channels > 2) {
    if (isCmyk == 1) {
      result =
          vips_icc_transform(in, out, srgb_profile_path, "input_profile",
                             "cmyk", "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    } else {
      result = vips_icc_transform(in, out, srgb_profile_path, "embedded", TRUE,
                                  "intent", VIPS_INTENT_PERCEPTUAL, NULL);
      // ignore embedded errors
      if (result != 0) {
        result = 0;
        *out = in;
      }
    }
  } else {
    result = vips_icc_transform(in, out, gray_profile_path, "input_profile",
                                gray_profile_path, "embedded", TRUE, "intent",
                                VIPS_INTENT_PERCEPTUAL, NULL);
  }

  return result;
}
