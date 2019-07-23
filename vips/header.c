#include "header.h"

unsigned long has_icc_profile(VipsImage *in) {
	return vips_image_get_typeof(in, VIPS_META_ICC_NAME);
}

gboolean remove_icc_profile(VipsImage *in) {
  return vips_image_remove(in, VIPS_META_ICC_NAME);
}

int get_meta_orientation(VipsImage *in) {
	int orientation = 0;
	if (vips_image_get_typeof(in, VIPS_META_ORIENTATION) != 0) {
		vips_image_get_int(in, VIPS_META_ORIENTATION, &orientation);
	}

	return orientation;
}
