
#include "bridge.h"

int is_16bit(VipsInterpretation interpretation) {
	return interpretation == VIPS_INTERPRETATION_RGB16 || interpretation == VIPS_INTERPRETATION_GREY16;
}

int init_image(void *buf, size_t len, int imageType, VipsImage **out) {
	int code = 1;

	if (imageType == JPEG) {
		code = vips_jpegload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == PNG) {
		code = vips_pngload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == WEBP) {
		code = vips_webpload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == TIFF) {
		code = vips_tiffload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#if (VIPS_MAJOR_VERSION >= 8)
#if (VIPS_MINOR_VERSION >= 3)
	} else if (imageType == GIF) {
		code = vips_gifload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == PDF) {
		code = vips_pdfload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
	} else if (imageType == SVG) {
		code = vips_svgload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#endif
	} else if (imageType == MAGICK) {
		code = vips_magickload_buffer(buf, len, out, "access", VIPS_ACCESS_RANDOM, NULL);
#endif
	}

	return code;
}

int remove_icc_profile(VipsImage *in) {
  return vips_image_remove(in, VIPS_META_ICC_NAME);
}

int load_jpeg_buffer(void *buf, size_t len, VipsImage **out, int shrink) {
  if (shrink > 0) {
    return vips_jpegload_buffer(buf, len, out, "shrink", shrink, NULL);
  } else {
    return vips_jpegload_buffer(buf, len, out, NULL);
  }
}

int save_jpeg_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int interlace) {
	return vips_jpegsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"optimize_coding", TRUE,
		"interlace", INT_TO_GBOOLEAN(interlace),
		NULL
	);
}

int save_png_buffer(VipsImage *in, void **buf, size_t *len, int strip, int compression, int quality, int interlace) {
	return vips_pngsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"compression", compression,
		"interlace", INT_TO_GBOOLEAN(interlace),
		"filter", VIPS_FOREIGN_PNG_FILTER_NONE,
		NULL
	);
}

int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality) {
	return vips_webpsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		NULL
	);
}

int save_tiff_buffer(VipsImage *in, void **buf, size_t *len) {
	return vips_tiffsave_buffer(in, buf, len, NULL);
}

int is_colorspace_supported(VipsImage *in) {
	return vips_colourspace_issupported(in) ? 1 : 0;
}

int to_colorspace(VipsImage *in, VipsImage **out, VipsInterpretation space) {
	return vips_colourspace(in, out, space, NULL);
}

int flip_image(VipsImage *in, VipsImage **out, int direction) {
	return vips_flip(in, out, direction, NULL);
}

int shrink_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink) {
	return vips_shrink(in, out, xshrink, yshrink, NULL);
}

int reduce_image(VipsImage *in, VipsImage **out, double xshrink, double yshrink) {
	return vips_reduce(in, out, xshrink, yshrink, NULL);
}

int zoom_image(VipsImage *in, VipsImage **out, int xfac, int yfac) {
	return vips_zoom(in, out, xfac, yfac, NULL);
}

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, double r, double g, double b) {
	if (extend == VIPS_EXTEND_BACKGROUND) {
		double background[3] = {r, g, b};
		VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);
		return vips_embed(in, out, left, top, width, height, "extend", extend, "background", vipsBackground, NULL);
	}
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
}

int extract_image_area(VipsImage *in, VipsImage **out, int left, int top, int width, int height) {
	return vips_extract_area(in, out, left, top, width, height, NULL);
}

int flatten_image_background(VipsImage *in, VipsImage **out, double r, double g, double b) {
	if (is_16bit(in->Type)) {
		r = 65535 * r / 255;
		g = 65535 * g / 255;
		b = 65535 * b / 255;
	}

	double background[3] = {r, g, b};
	VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);

	return vips_flatten(in, out,
		"background", vipsBackground,
		"max_alpha", is_16bit(in->Type) ? 65535.0 : 255.0,
		NULL
	);
}

int transform_image(VipsImage *in, VipsImage **out, double a, double b, double c, double d, VipsInterpolate *interpolator) {
	return vips_affine(in, out, a, b, c, d, "interpolate", interpolator, NULL);
}

int find_image_loader(int t) {
  switch (t) {
    case GIF:
      return vips_type_find("VipsOperation", "gifload");
  }
	if (t == GIF) {
		return vips_type_find("VipsOperation", "gifload");
	} else if (t == PDF) {
		return vips_type_find("VipsOperation", "pdfload");
	} if (t == TIFF) {
		return vips_type_find("VipsOperation", "tiffload");
	}
	if (t == SVG) {
		return vips_type_find("VipsOperation", "svgload");
	}
	if (t == WEBP) {
		return vips_type_find("VipsOperation", "webpload");
	}
	if (t == PNG) {
		return vips_type_find("VipsOperation", "pngload");
	}
	if (t == JPEG) {
		return vips_type_find("VipsOperation", "jpegload");
	}
	if (t == MAGICK) {
		return vips_type_find("VipsOperation", "magickload");
	}
	return 0;
}

int find_image_type_saver(int t) {
	if (t == TIFF) {
		return vips_type_find("VipsOperation", "tiffsave_buffer");
	}
	if (t == WEBP) {
		return vips_type_find("VipsOperation", "webpsave_buffer");
	}
	if (t == PNG) {
		return vips_type_find("VipsOperation", "pngsave_buffer");
	}
	if (t == JPEG) {
		return vips_type_find("VipsOperation", "jpegsave_buffer");
	}
	return 0;
}

void gobject_set_property(VipsObject *object, const char *name, const GValue *value) {
  VipsObjectClass *object_class = VIPS_OBJECT_GET_CLASS( object );
  GType type = G_VALUE_TYPE( value );

  GParamSpec *pspec;
  VipsArgumentClass *argument_class;
  VipsArgumentInstance *argument_instance;

  if( vips_object_get_argument( object, name,
    &pspec, &argument_class, &argument_instance ) ) {
    vips_warn( NULL, "%s", vips_error_buffer() );
    vips_error_clear();
    return;
  }

  if( G_IS_PARAM_SPEC_ENUM( pspec ) &&
    type == G_TYPE_STRING ) {
    GType pspec_type = G_PARAM_SPEC_VALUE_TYPE( pspec );

    int enum_value;
    GValue value2 = { 0 };

    if( (enum_value = vips_enum_from_nick( object_class->nickname,
      pspec_type, g_value_get_string( value ) )) < 0 ) {
      vips_warn( NULL, "%s", vips_error_buffer() );
      vips_error_clear();
      return;
    }

    g_value_init( &value2, pspec_type );
    g_value_set_enum( &value2, enum_value );
    g_object_set_property( G_OBJECT( object ), name, &value2 );
    g_value_unset( &value2 );
  } else {
    g_object_set_property( G_OBJECT( object ), name, value );
  }
}
