// https://libvips.github.io/libvips/API/current/libvips-colour.html

#include <stdlib.h>
#include <vips/vips.h>

int is_colorspace_supported(VipsImage *in);
int to_colorspace(VipsImage *in, VipsImage **out, VipsInterpretation space);

int optimize_icc_profile(VipsImage *in, VipsImage **out, int isCmyk, char *srgb_profile_path, char *gray_profile_path);
