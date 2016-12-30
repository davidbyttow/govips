#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

#if (VIPS_MAJOR_VERSION < 8)
  error_requires_version_8
#endif


void set_property(VipsObject* object, const char* name, const GValue* value);
void filename_split8(const char* name, char *filename, char *option_string);
