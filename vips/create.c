// clang-format off
// include order matters
#include "lang.h"
#include "create.h"
// clang-format on

// https://libvips.github.io/libvips/API/current/libvips-create.html#vips-text
int text(VipsImage **out, TextOptions *o) {
  return vips_text(out, o->Text, "font", o->Font, "width", o->Width, "height", o->Height, "align", o->Align,
  "dpi", o->DPI, "rgba", o->RGBA, "justify", o->Justify, "spacing", o->Spacing, "wrap", o->Wrap, NULL);
}
