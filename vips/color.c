#include "color.h"
#include "icc_profiles.h"
#include <unistd.h>

#define SRGB_PROFILE_PATH SRGB_V2_MICRO_ICC_PATH
#define GRAY_PROFILE_PATH SGRAY_V2_MICRO_ICC_PATH

int is_colorspace_supported(VipsImage *in) {
	return vips_colourspace_issupported(in) ? 1 : 0;
}

int to_colorspace(VipsImage *in, VipsImage **out, VipsInterpretation space) {
	return vips_colourspace(in, out, space, NULL);
}

// https://libvips.github.io/libvips/API/8.6/libvips-colour.html#vips-icc-transform
int optimize_icc_profile(VipsImage *in, VipsImage **out, const char *input_profile) {
	const char* target_profile = vips_image_get_bands(in) > 2 ? SRGB_PROFILE_PATH : GRAY_PROFILE_PATH;

    if (!input_profile && !vips_image_get_typeof(in, VIPS_META_ICC_NAME)) {
    	//No input profile and no embedded ICC profile in the input image, nothing to do.
		*out = in;
		return 0;
    }

    if (vips_icc_transform(
    	in, out, target_profile,
    	"embedded", !input_profile,
    	"input_profile", input_profile ? input_profile : "none",
    	"intent", VIPS_INTENT_PERCEPTUAL,
    	NULL)) {
    	return -1;
    }

    vips_image_set_string(*out, "optimized-icc-profile", target_profile);
    return 0;
}
