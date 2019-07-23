#include "header.h"

unsigned long has_icc_profile(VipsImage *in) {
	return vips_image_get_typeof(in, VIPS_META_ICC_NAME);
}

gboolean remove_icc_profile(VipsImage *in) {
  return vips_image_remove(in, VIPS_META_ICC_NAME);
}
