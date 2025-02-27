// https://libvips.github.io/libvips/API/current/libvips-create.html

// clang-format off
// include order matters
#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>
// clang-format on
typedef struct {
  const char *Text;
  const char *Font;
  int Width;
  int Height;
  int DPI;
  gboolean RGBA;
  gboolean Justify;
  int Spacing;
  VipsAlign Align;
  VipsTextWrap Wrap;
} TextOptions;

int xyz(VipsImage **out, int width, int height);
int black(VipsImage **out, int width, int height);
int identity(VipsImage **out, int ushort);
int text(VipsImage **out, TextOptions *o);
