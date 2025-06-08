#include <stdlib.h>
#include <vips/vips.h>

typedef struct {
  const char *Text;
  const char *Font;
  int Width;
  int Height;
  int OffsetX;
  int OffsetY;
  VipsAlign Align;
  int DPI;
  int Margin;
  float Opacity;
  double Color[3];
} LabelOptions;

int label(VipsImage *in, VipsImage **out, LabelOptions *o);
