// https://libvips.github.io/libvips/API/current/VipsForeignSave.html

// clang-format off
// include order matters
#include <stdlib.h>

#include <vips/vips.h>
#include <vips/foreign.h>
// clang-format n

#ifndef BOOL
#define BOOL int
#endif

typedef enum types {
  UNKNOWN = 0,
  JPEG,
  WEBP,
  PNG,
  TIFF,
  GIF,
  PDF,
  SVG,
  MAGICK,
  HEIF,
  BMP
} ImageType;

int load_image_buffer(void *buf, size_t len, int imageType, VipsImage **out);

typedef struct SaveParams {
  VipsImage *inputImage;
  void *outputBuffer;
  ImageType outputFormat;
  size_t outputLen;

  BOOL stripMetadata;
  int quality;
  BOOL interlace;

  // JPEG
  BOOL jpegOptimizeCoding;
  VipsForeignJpegSubsample jpegSubsample;

  // PNG
  int pngCompression;
  VipsForeignPngFilter pngFilter;

  // WEBP
  BOOL webpLossless;
  int webpReductionEffort;

  // HEIF
  BOOL heifLossless;

  // TIFF
  VipsForeignTiffCompression tiffCompression;
  VipsForeignTiffPredictor tiffPredictor;
  BOOL tiffPyramid;
  BOOL tiffTile;
  int tiffTileHeight;
  int tiffTileWidth;
  double tiffXRes;
  double tiffYRes;

} SaveParams;

SaveParams create_save_params(ImageType outputFormat);
void init_save_params(SaveParams *params);

int save_to_buffer(SaveParams *params);
