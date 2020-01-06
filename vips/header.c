#include "header.h"
#include "icc_profiles.h"
#include <unistd.h>

unsigned long has_icc_profile(VipsImage *in) {
    return vips_image_get_typeof(in, VIPS_META_ICC_NAME);
}


// todo: move to color.(go, h, c)
// https://libvips.github.io/libvips/API/8.6/libvips-colour.html#vips-icc-transform
int icc_transform(VipsImage *in, VipsImage **out, int isCmyk) {
    int channels = vips_image_get_bands(in);

    int result;

    char *srgb_profile_path = SRGB_V2_MICRO_ICC_PATH;  // SRGB_IEC61966_2_1_ICC_PATH
    char *gray_profile_path = SGRAY_V2_MICRO_ICC_PATH;  // GENERIC_GRAY_GAMMA_2_2_ICC_PATH

    if (channels > 2) {
    	if (isCmyk == 1) {
    		result = vips_icc_transform(in, out, srgb_profile_path, "input_profile", "cmyk", "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    	} else {
        result = vips_icc_transform(in, out, srgb_profile_path, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    	}
    } else {
			result = vips_icc_transform(in, out, gray_profile_path, "input_profile", gray_profile_path, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    }

    return result;
}

gboolean remove_icc_profile(VipsImage *in) {
    return vips_image_remove(in, VIPS_META_ICC_NAME);
}

// won't remove the ICC profile
void remove_metadata(VipsImage *in) {
    gchar ** fields = vips_image_get_fields(in);

    for (int i=0; fields[i] != NULL; i++) {
        if (strncmp(fields[i], VIPS_META_ICC_NAME, 16)) {
            vips_image_remove(in, fields[i]);
        }
    }

    g_strfreev(fields);
}

int get_meta_orientation(VipsImage *in) {
	int orientation = 0;
	if (vips_image_get_typeof(in, VIPS_META_ORIENTATION) != 0) {
		vips_image_get_int(in, VIPS_META_ORIENTATION, &orientation);
	}

    return orientation;
}

void set_meta_orientation(VipsImage *in, int orientation) {
	vips_image_set_int(in, VIPS_META_ORIENTATION, orientation);
}
