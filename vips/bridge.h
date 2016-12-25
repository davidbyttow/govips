#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

#if (VIPS_MAJOR_VERSION < 8)
  error_requires_version_8
#endif

inline int init_image(void *buf, size_t len, int imageType, VipsImage **out) {
  int ret = 1;
  ret = vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
  return ret;
}

inline int shrink(VipsImage* in, VipsImage** out, double xshrink, double yshrink) {
  return vips_shrink(in, out, xshrink, yshrink, NULL);
}

