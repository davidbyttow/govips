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
  BMP,
  AVIF
} ImageType;

typedef struct LoadParams {
  size_t inputLen;
  ImageType inputFormat;

  BOOL autorotate;
  BOOL fail;
  int page;
  int n;
  gdouble dpi;

  int jpegShrink;
  BOOL heifThumbnail;
  BOOL svgUnlimited;
} LoadParams;

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
  BOOL jpegTrellisQuant;
  BOOL jpegOvershootDeringing;
  BOOL jpegOptimizeScans;
  int jpegQuantTable;

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

  // AVIF
  int avifSpeed;
} SaveParams;

SaveParams create_save_params(ImageType outputFormat);
int save_to_buffer(SaveParams *params);

LoadParams create_load_params(ImageType inputFormat);
int load_image_buffer(LoadParams *params, void *buf, size_t len, VipsImage **out);
