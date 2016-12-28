#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

#if (VIPS_MAJOR_VERSION < 8)
  error_requires_version_8
#endif

int init_image(void *buf, size_t len, int imageType, VipsImage **out);
void set_property(VipsObject *object, const char *name, const GValue *value);
