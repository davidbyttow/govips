#include "foreign.h"

#include "lang.h"

int load_image_buffer(void *buf, size_t len, int imageType, gboolean unlimitedSvgSize, VipsImage **out) {
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
    code = vips_svgload_buffer(buf, len, out,  "unlimited", unlimitedSvgSize, NULL);
  } else if (imageType == HEIF) {
    // added autorotate on load as currently it addresses orientation issues
    // https://github.com/libvips/libvips/pull/1680
    code = vips_heifload_buffer(buf, len, out, "autorotate", TRUE, NULL);
  } else if (imageType == MAGICK) {
    code = vips_magickload_buffer(buf, len, out, NULL);
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

  return 0;
}

// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-pngsave-buffer
int set_png_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "strip", params->stripMetadata,
                  "compression", params->pngCompression, "interlace",
                  params->interlace, "filter", params->pngFilter, NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }

  return 0;
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/webpsave.c#L524
// https://libvips.github.io/libvips/API/current/VipsForeignSave.html#vips-webpsave-buffer
int set_webp_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(
    VIPS_OBJECT(operation),
    "strip", params->stripMetadata,
    "lossless", params->webpLossless,
    "near_lossless", params->webpNearLossless,
    "reduction_effort", params->webpReductionEffort,
    "profile", params->webpIccProfile ? params->webpIccProfile : "none",
    NULL
  );

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
  return 0;
}

// todo: support additional params
// https://github.com/libvips/libvips/blob/master/libvips/foreign/heifsave.c#L653
int set_heif_options(VipsOperation *operation, SaveParams *params) {
  vips_object_set(VIPS_OBJECT(operation), "lossless", params->heifLossless,
                  NULL);

  if (params->quality) {
    vips_object_set(VIPS_OBJECT(operation), "Q", params->quality, NULL);
  }
  return 0;
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
  return 0;
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
  return 0;
}

static SaveParams defaultSaveParams = {
  .inputImage = NULL,
  .outputBuffer = NULL,
  .outputFormat = JPEG,
  .outputLen = 0,

  .interlace = FALSE,
  .quality = 0,
  .stripMetadata = FALSE,

  .jpegOptimizeCoding = TRUE,
  .jpegSubsample = VIPS_FOREIGN_JPEG_SUBSAMPLE_OFF,

  .pngCompression = 6,
  .pngFilter = VIPS_FOREIGN_PNG_FILTER_NONE,

  .webpLossless = FALSE,
  .webpNearLossless = FALSE,
  .webpReductionEffort = 4,
  .webpIccProfile = NULL,

  .heifLossless = FALSE,

  .tiffCompression = VIPS_FOREIGN_TIFF_COMPRESSION_LZW,
  .tiffPredictor = VIPS_FOREIGN_TIFF_PREDICTOR_HORIZONTAL,
  .tiffPyramid = FALSE,
  .tiffTile = FALSE,
  .tiffTileHeight = 256,
  .tiffTileWidth = 256,
  .tiffXRes = 1.0,
  .tiffYRes = 1.0
};

SaveParams create_save_params(ImageType outputFormat) {
  SaveParams params = defaultSaveParams;
  params.outputFormat = outputFormat;
  return params;
}

void init_save_params(SaveParams *params) { *params = defaultSaveParams; }
