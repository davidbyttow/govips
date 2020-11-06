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
int optimize_icc_profile(VipsImage *in, VipsImage **out, int isCmyk) {
	// todo: check current embedded profile, and skip if already set
    int channels = vips_image_get_bands(in);
	int result;
	const char* target_profile;

    if (channels > 2) {
    	if (isCmyk == 1) {
    		result = vips_icc_transform(in, out, SRGB_PROFILE_PATH, "input_profile", "cmyk", "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    	} else {
        	result = vips_icc_transform(in, out, SRGB_PROFILE_PATH, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
        	// ignore embedded errors
        	if (result != 0) {
        		result = 0;
        		*out = in;
        	}
    	}
    	target_profile = SRGB_PROFILE_PATH;
    } else {
		result = vips_icc_transform(in, out, GRAY_PROFILE_PATH, "input_profile", GRAY_PROFILE_PATH, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
     	target_profile = GRAY_PROFILE_PATH;
    }
    vips_image_set_string(*out, "target-icc-profile", target_profile);
    return result;
}
