#include "lang.h"
#include "create.h"

// https://libvips.github.io/libvips/API/current/libvips-create.html#vips-xyz
int xyz(VipsImage **out, int width, int height){
	return vips_xyz(out, width, height, NULL);
}
