#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

#if (VIPS_MAJOR_VERSION < 8)
  error_requires_version_8
#endif

void SetProperty(VipsObject* object, const char* name, const GValue* value);
int VipsJpegsaveBuffer(VipsImage* image, void **buf, size_t *len, int strip, int quality, int interlace);
