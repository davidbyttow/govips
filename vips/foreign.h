// https://libvips.github.io/libvips/API/current/VipsForeignSave.html

#include <stdlib.h>
#include <vips/vips.h>
#include <vips/foreign.h>

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
  
  // PNG
  int pngCompression;

  // WEBP
  BOOL webpLossless;
  int webpReductionEffort;

  // HEIF
  BOOL heifLossless;

  // TIFF
  VipsForeignTiffCompression tiffCompression;
} SaveParams;

int save_to_buffer(SaveParams *params);
