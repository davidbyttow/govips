// https://libvips.github.io/libvips/API/current/libvips-header.html

#include <stdlib.h>
#include <vips/vips.h>


unsigned long has_icc_profile(VipsImage *in);
int remove_icc_profile(VipsImage *in);
