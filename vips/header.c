#include "header.h"
#include "icc_profiles.h"
#include <unistd.h>

unsigned long has_icc_profile(VipsImage *in) {
    return vips_image_get_typeof(in, VIPS_META_ICC_NAME);
}

int icc_transform(VipsImage *in, VipsImage **out, int isCmyk) {
    int channels = vips_image_get_bands(in);

    int result;

    char *srgb_profile = srgb_v2_micro_icc_path;  // srgb_iec61966_2_1_icc_path
    char *grey_profile = generic_gray_gamma_2_2_icc_path;  // sgrey_v2_micro_icc_path

    if( channels > 2 ) {
    	if (isCmyk == 1) {
    		result = vips_icc_transform(in, out, srgb_profile, "input_profile", "cmyk", "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    	} else {
        result = vips_icc_transform(in, out, srgb_profile, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
    	}

    } else {
        result = vips_icc_transform(in, out, grey_profile, "input_profile", grey_profile, "embedded", TRUE, "intent", VIPS_INTENT_PERCEPTUAL, NULL);
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
