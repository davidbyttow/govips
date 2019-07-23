#include "image.h"

int has_alpha_channel(VipsImage *image) {
	return vips_image_hasalpha(image);
}
