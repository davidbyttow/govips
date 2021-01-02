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

// TODO: Pass options as discrete params objects based on types rather than long function signatures
int save_jpeg_buffer(VipsImage* image, void **buf, size_t *len, int strip, int quality, int interlace);
int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int interlace);
int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless, int effort);
int save_heif_buffer(VipsImage *in, void **buf, size_t *len, int quality, int lossless);
int save_tiff_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless);

typedef struct SaveAsParams {
  VipsImage *inputImage;
  void **outputBuffer;
  ImageType outputFormat;
  size_t *outputLen;

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
} SaveAsParams;

int save_to_buffer(SaveAsParams *params);
