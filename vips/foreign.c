#include "foreign.h"
#include "govips.h"
#include "lang.h"

int load_image_buffer_operation(VipsOperation *operation, VipsBlob *blob, VipsImage **out, ImportParams *params) {
  int result;

  // Set operation parameters
  result = vips_object_set(VIPS_OBJECT(operation), "buffer", blob, "out", out, NULL);
  if (result) {
    return result;
  }

  if (params->fail_set) {
    result = vips_object_set(VIPS_OBJECT(operation), "fail", params->fail, NULL);
    if (result) {
      govipsLoggingHandler("govips", G_LOG_LEVEL_INFO, "can't set optional parameter \"fail\"");
    }
  }

  if (params->autorotate_set) {
    result = vips_object_set(VIPS_OBJECT(operation), "autorotate", params->autorotate, NULL);
    if (result) {
      govipsLoggingHandler("govips", G_LOG_LEVEL_INFO, "can't set optional parameter \"autorotate\"");
    }
  }

  if (params->shrink_set) {
    result = vips_object_set(VIPS_OBJECT(operation), "shrink", params->shrink, NULL);
    if (result) {
      return result;
    }
  }

  // Execute operation
  result = vips_cache_operation_buildp(&operation);
  if (result) {
    return result;
  }

  // Get operation output
  g_object_get(VIPS_OBJECT(operation), "out", out, NULL);
  return 0;
}

int load_image_buffer_generic(const char *method, void *buf, size_t len, VipsImage **out, ImportParams *params) {
  VipsOperation *operation;
  VipsBlob *blob;
  int result;

  operation = vips_operation_new(method);
  if (operation == NULL) {
    return -1;
  }

  blob = vips_blob_new(NULL, buf, len);
  result = load_image_buffer_operation(operation, blob, out, params);

  vips_area_unref(VIPS_AREA(blob));
  vips_object_unref_outputs(VIPS_OBJECT(operation));
  g_object_unref(operation);
  return result;
}

int load_image_buffer(void *buf, size_t len, int imageType, VipsImage **out, ImportParams *params) {
  int code = 1;

  if (imageType == JPEG) {
    code = load_image_buffer_generic("jpegload_buffer", buf, len, out, params);
  } else if (imageType == PNG) {
    code = load_image_buffer_generic("pngload_buffer", buf, len, out, params);
  } else if (imageType == WEBP) {
    code = load_image_buffer_generic("webpload_buffer", buf, len, out, params);
  } else if (imageType == TIFF) {
    code = load_image_buffer_generic("tiffload_buffer", buf, len, out, params);
  } else if (imageType == GIF) {
    code = load_image_buffer_generic("gifload_buffer", buf, len, out, params);
  } else if (imageType == PDF) {
    code = load_image_buffer_generic("pdfload_buffer", buf, len, out, params);
  } else if (imageType == SVG) {
    code = load_image_buffer_generic("svgload_buffer", buf, len, out, params);
  } else if (imageType == HEIF) {
    // added autorotate on load as currently it addresses orientation issues
    // https://github.com/libvips/libvips/pull/1680
    params->autorotate_set = TRUE;
    params->autorotate = TRUE;
    code = load_image_buffer_generic("heifload_buffer", buf, len, out, params);
  } else if (imageType == MAGICK) {
    code = load_image_buffer_generic("magickload_buffer", buf, len, out, params);
  }

  return code;
}

typedef int (*VipsBuildOperationFn)(VipsOperation *operation,
                                    SaveParams *params);

int save_buffer(const char *operationName, SaveParams *params,
                VipsBuildOperationFn buildFn) {
  VipsBlob *blob;
  VipsOperation *operation = vips_operation_new(operationName);
  if (!operation) {
    return -1;
  }

  if (vips_object_set(VIPS_OBJECT(operation), "in", params->inputImage, NULL)) {
    return -1;
  }

  if (buildFn(operation, params)) {
    g_object_unref(operation);
    return -1;
  }

  if (vips_cache_operation_buildp(&operation)) {
    vips_object_unref_outputs(VIPS_OBJECT(operation));
    g_object_unref(operation);
    return -1;
  }

  g_object_get(VIPS_OBJECT(operation), "buffer", &blob, NULL);
  g_object_unref(operation);

  VipsArea *area = VIPS_AREA(blob);

  params->outputBuffer = (char *)(area->data);
  params->outputLen = area->length;
  area->free_fn = NULL;
  vips_area_unref(area);

  return 0;
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-jpegsave-buffer
int set_jpeg_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "strip", params->stripMetadata,
                  "optimize_coding", params->jpegOptimizeCoding, "interlace",
                  params->interlace, "subsample_mode", params->jpegSubsample,
                  NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-pngsave-buffer
int set_png_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "strip", params->stripMetadata,
                  "compression", params->pngCompression, "interlace",
                  params->interlace, "filter", params->pngFilter, NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/webpsave.c#L524
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-webpsave-buffer
int set_webp_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "strip", params->stripMetadata,
                  "lossless", params->webpLossless, "reduction_effort",
                  params->webpReductionEffort, NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/heifsave.c#L653
int set_heif_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "lossless", params->heifLossless,
                  NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-tiffsave-buffer
int set_tiff_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "strip", params->stripMetadata,
                  "compression", params->tiffCompression, "predictor",
                  params->tiffPredictor, "pyramid", params->tiffPyramid,
                  "tile_height", params->tiffTileHeight, "tile_width",
                  params->tiffTileWidth, "tile", params->tiffTile, "xres",
                  params->tiffXRes, "yres", params->tiffYRes, NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
}

int save_to_buffer(SaveParams *params) {
  switch (params->outputFormat) {
    case JPEG:
      return save_buffer("jpegsave_buffer", params, set_jpeg_options);
    case PNG:
      return save_buffer("pngsave_buffer", params, set_png_options);
    case WEBP:
      return save_buffer("webpsave_buffer", params, set_webp_options);
    case HEIF:
      return save_buffer("heifsave_buffer", params, set_heif_options);
    case TIFF:
      return save_buffer("tiffsave_buffer", params, set_tiff_options);
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
