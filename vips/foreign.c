#include "foreign.h"

#include "lang.h"

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
  return vips_jpegsave_buffer(
      params->inputImage, &params->outputBuffer, &params->outputLen, "strip",
      params->stripMetadata, "Q", params->quality, "optimize_coding",
      params->jpegOptimizeCoding, "interlace", params->interlace,
      "subsample_mode", params->jpegSubsample, NULL);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-pngsave-buffer
int save_png_buffer(SaveParams *params) {
  return vips_pngsave_buffer(
      params->inputImage, &params->outputBuffer, &params->outputLen, "strip",
      params->stripMetadata, "compression", params->pngCompression, "interlace",
      params->interlace, "filter", params->pngFilter, NULL);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/webpsave.c#L524
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-webpsave-buffer
int save_webp_buffer(SaveParams *params) {
  return vips_webpsave_buffer(
      params->inputImage, &params->outputBuffer, &params->outputLen, "strip",
      params->stripMetadata, "Q", params->quality, "lossless",
      params->webpLossless, "reduction_effort", params->webpReductionEffort,
      NULL);
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/heifsave.c#L653
int save_heif_buffer(SaveParams *params) {
  return vips_heifsave_buffer(params->inputImage, &params->outputBuffer,
                              &params->outputLen, "Q", params->quality,
                              "lossless", params->heifLossless, NULL);
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-tiffsave-buffer
int save_tiff_buffer(SaveParams *params) {
  return vips_tiffsave_buffer(
      params->inputImage, &params->outputBuffer, &params->outputLen, "strip",
      params->stripMetadata, "Q", params->quality, "compression",
      params->tiffCompression, "predictor", params->tiffPredictor, "pyramid",
      params->tiffPyramid, "tile_height", params->tiffTileHeight, "tile_width",
      params->tiffTileWidth, "tile", params->tiffTile, "xres", params->tiffXRes,
      "yres", params->tiffYRes, NULL);
}

int save_to_buffer(SaveParams *params) {
  switch (params->outputFormat) {
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

static SaveParams defaultSaveParams = {
  inputImage : NULL,
  outputBuffer : NULL,
  outputFormat : JPEG,
  outputLen : 0,

  interlace : FALSE,
  quality : 0,
  stripMetadata : FALSE,

  jpegOptimizeCoding : FALSE,
  jpegSubsample : VIPS_FOREIGN_JPEG_SUBSAMPLE_ON,

  pngCompression : 6,
  pngFilter : VIPS_FOREIGN_PNG_FILTER_NONE,

  webpLossless : FALSE,
  webpReductionEffort : 4,

  heifLossless : FALSE,

  tiffCompression : VIPS_FOREIGN_TIFF_COMPRESSION_LZW,
  tiffPredictor : VIPS_FOREIGN_TIFF_PREDICTOR_HORIZONTAL,
  tiffPyramid : FALSE,
  tiffTile : FALSE,
  tiffTileHeight : 256,
  tiffTileWidth : 256,
  tiffXRes : 1.0,
  tiffYRes : 1.0
};

SaveParams create_save_params(ImageType outputFormat) {
  SaveParams params = defaultSaveParams;
  params.outputFormat = outputFormat;
  return params;
}

void init_save_params(SaveParams *params) { *params = defaultSaveParams; }