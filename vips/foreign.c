#include "lang.h"
#include "foreign.h"

int load_image_buffer(void *buf, size_t len, int imageType, VipsImage **out) {
  int code = 1;

  if (imageType == JPEG) {
    code = vips_jpegload_buffer(buf, len, out, NULL);
  } else if (imageType == PNG) {
    code = vips_pngload_buffer(buf, len, out, NULL);
  } else if (imageType == WEBP) {
    code = vips_webpload_buffer(buf, len, out, NULL);
  } else if (imageType == TIFF) {
    code = vips_tiffload_buffer(buf, len, out, NULL);
  } else if (imageType == GIF) {
    code = vips_gifload_buffer(buf, len, out, NULL);
  } else if (imageType == PDF) {
    code = vips_pdfload_buffer(buf, len, out, NULL);
  } else if (imageType == SVG) {
    code = vips_svgload_buffer(buf, len, out, NULL);
  } else if (imageType == HEIF) {
    // added autorotate on load as currently it addresses orientation issues
    // https://github.com/libvips/libvips/pull/1680
    code = vips_heifload_buffer(buf, len, out, "autorotate", TRUE, NULL);
  } else if (imageType == MAGICK) {
    code = vips_magickload_buffer(buf, len, out, NULL);
  }

  return code;
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-jpegsave-buffer
int save_jpeg_buffer(SaveParams *params) {
  return vips_jpegsave_buffer(params->inputImage,
                              &params->outputBuffer,
                              &params->outputLen,
                              "strip", INT_TO_GBOOLEAN(params->stripMetadata),
                              "Q", params->quality,
                              "optimize_coding", TRUE,
                              "interlace", INT_TO_GBOOLEAN(params->interlace),
                              "subsample_mode", VIPS_FOREIGN_JPEG_SUBSAMPLE_ON,
                              NULL);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-pngsave-buffer
int save_png_buffer(SaveParams *params) {
  return vips_pngsave_buffer(params->inputImage,
                             &params->outputBuffer,
                             &params->outputLen,
                             "strip", INT_TO_GBOOLEAN(params->stripMetadata),
                             "compression", params->pngCompression,
                             "interlace", INT_TO_GBOOLEAN(params->interlace),
                             "filter", VIPS_FOREIGN_PNG_FILTER_NONE,
                             NULL);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/webpsave.c#L524
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-webpsave-buffer
int save_webp_buffer(SaveParams *params) {
  return vips_webpsave_buffer(params->inputImage,
                              &params->outputBuffer,
                              &params->outputLen,
                              "strip", INT_TO_GBOOLEAN(params->stripMetadata),
                              "Q", params->quality,
                              "lossless", INT_TO_GBOOLEAN(params->webpLossless),
                              "reduction_effort", params->webpReductionEffort,
                              NULL);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/heifsave.c#L653
int save_heif_buffer(SaveParams *params) {
  return vips_heifsave_buffer(params->inputImage,
                              &params->outputBuffer,
                              &params->outputLen,
                              "Q", params->quality,
                              "lossless", INT_TO_GBOOLEAN(params->heifLossless),
                              NULL);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-tiffsave-buffer
int save_tiff_buffer(SaveParams *params) {
  return vips_tiffsave_buffer(params->inputImage,
                              &params->outputBuffer,
                              &params->outputLen,
                              "strip", INT_TO_GBOOLEAN(params->stripMetadata),
                              "Q", params->quality,
                              "compression", params->tiffCompression,
                              "pyramid", FALSE,
                              "predictor", VIPS_FOREIGN_TIFF_PREDICTOR_HORIZONTAL,
                              "pyramid", FALSE,
                              "tile", FALSE,
                              "tile_height", 256,
                              "tile_width", 256,
                              "xres", 1.0,
                              "yres", 1.0,
                              NULL);
}

int save_to_buffer(SaveParams *params) {
  switch (params->outputFormat)
  {
  case JPEG:
    return save_jpeg_buffer(params);
  case PNG:
    return save_png_buffer(params);
  case WEBP:
    return save_webp_buffer(params);
  case HEIF:
    return save_heif_buffer(params);
  case TIFF:
    return save_tiff_buffer(params);
  default:
    g_warning("Unsupported output type given: %d", params->outputFormat);
    return -1;
  }
}

