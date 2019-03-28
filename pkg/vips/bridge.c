
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

unsigned long has_profile_embed(VipsImage *in) {
	return vips_image_get_typeof(in, VIPS_META_ICC_NAME);
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

int save_webp_buffer(VipsImage *in, void **buf, size_t *len, int strip, int quality, int lossless) {
	return vips_webpsave_buffer(in, buf, len,
		"strip", INT_TO_GBOOLEAN(strip),
		"Q", quality,
		"lossless", INT_TO_GBOOLEAN(lossless),
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

int resize_image(VipsImage *in, VipsImage **out, double scale, double vscale, int kernel) {
	if (vscale > 0) {
		return vips_resize(in, out, scale, "vscale", vscale, "kernel", kernel, NULL);
	}
	return vips_resize(in, out, scale, "kernel", kernel, NULL);
}

int rot_image(VipsImage *in, VipsImage **out, VipsAngle angle) {
  return vips_rot(in, out, angle, NULL);
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

int text(VipsImage **out, const char *text, const char *font, int width, int height, VipsAlign align, int dpi) {
	return vips_text(out, text, "font", font, "width", width, "height", height, "align", align, "dpi", dpi, NULL);
}

int gaussian_blur(VipsImage *in, VipsImage **out, double sigma) {
	return vips_gaussblur(in, out, sigma, NULL);
}

int invert_image(VipsImage *in, VipsImage **out) {
	return vips_invert(in, out, NULL);
}

int extract_band(VipsImage *in, VipsImage **out, int band, int num) {
	if (num > 0) {
		return vips_extract_band(in, out, band, "n", num, NULL);
	}
	return vips_extract_band(in, out, band, NULL);
}

int linear1(VipsImage *in, VipsImage **out, double a, double b) {
	return vips_linear1(in, out, a, b, NULL);
}

int embed_image(VipsImage *in, VipsImage **out, int left, int top, int width, int height, int extend, double r, double g, double b) {
	if (extend == VIPS_EXTEND_BACKGROUND) {
		double background[3] = {r, g, b};
		VipsArrayDouble *vipsBackground = vips_array_double_new(background, 3);
		return vips_embed(in, out, left, top, width, height, "extend", extend, "background", vipsBackground, NULL);
	}
	return vips_embed(in, out, left, top, width, height, "extend", extend, NULL);
}

int composite(VipsImage **in, VipsImage **out, int n, int mode) {
	return vips_composite(in, out, n, &mode, NULL);
}

int add(VipsImage *left, VipsImage *right, VipsImage **out) {
	return vips_add(left, right, out, NULL);
}

int multiply(VipsImage *left, VipsImage *right, VipsImage **out) {
	return vips_multiply(left, right, out, NULL);
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
    case PDF:
      return vips_type_find("VipsOperation", "pdfload");
    case TIFF:
      return vips_type_find("VipsOperation", "tiffload");
    case SVG:
      return vips_type_find("VipsOperation", "svgload");
    case WEBP:
      return vips_type_find("VipsOperation", "webpload");
    case PNG:
      return vips_type_find("VipsOperation", "pngload");
    case JPEG:
      return vips_type_find("VipsOperation", "jpegload");
    case MAGICK:
      return vips_type_find("VipsOperation", "magickload");
  }
	return 0;
}

int find_image_type_saver(int t) {
  switch (t) {
    case TIFF:
      return vips_type_find("VipsOperation", "tiffsave_buffer");
    case WEBP:
      return vips_type_find("VipsOperation", "webpsave_buffer");
    case PNG:
      return vips_type_find("VipsOperation", "pngsave_buffer");
    case JPEG:
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
    vips_warn( NULL, "gobject warning: %s", vips_error_buffer() );
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
      vips_warn( NULL, "gobject warning: %s", vips_error_buffer() );
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

int label(VipsImage *in, VipsImage **out, LabelOptions *o) {
	double ones[3] = { 1, 1, 1 };
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);
	t[0] = in;
	if (
		vips_text(&t[1], o->Text,
			"font", o->Font,
			"width", o->Width,
			"height", o->Height,
			"align", o->Align,
			NULL) ||
		vips_linear1(t[1], &t[2], o->Opacity, 0.0, NULL) ||
		vips_cast(t[2], &t[3], VIPS_FORMAT_UCHAR, NULL) ||
		vips_embed(t[3], &t[4], o->OffsetX, o->OffsetY, t[3]->Xsize + o->OffsetX, t[3]->Ysize + o->OffsetY, NULL)
		) {
		g_object_unref(base);
		return 1;
	}
	if (
		vips_black(&t[5], 1, 1, NULL) ||
		vips_linear(t[5], &t[6], ones, o->Color, 3, NULL) ||
		vips_cast(t[6], &t[7], VIPS_FORMAT_UCHAR, NULL) ||
		vips_copy(t[7], &t[8], "interpretation", t[0]->Type, NULL) ||
		vips_embed(t[8], &t[9], 0, 0, t[0]->Xsize, t[0]->Ysize, "extend", VIPS_EXTEND_COPY, NULL)
		) {
		g_object_unref(base);
		return 1;
	}
	if (vips_ifthenelse(t[4], t[9], t[0], out, "blend", TRUE, NULL)) {
		g_object_unref(base);
		return 1;
	}
	g_object_unref(base);
	return 0;
}

/////////////////////////////////////////////////
int vips_add_band(VipsImage *in, VipsImage **out, double c) {
#if (VIPS_MAJOR_VERSION > 8 || (VIPS_MAJOR_VERSION >= 8 && VIPS_MINOR_VERSION >= 2))
	return vips_bandjoin_const1(in, out, c, NULL);
#else
	VipsImage *base = vips_image_new();
	if (
		vips_black(&base, in->Xsize, in->Ysize, NULL) ||
		vips_linear1(base, &base, 1, c, NULL)) {
			g_object_unref(base);
			return 1;
		}
	g_object_unref(base);
	return vips_bandjoin2(in, base, out, c, NULL);
#endif
}

int vips_watermark_image(VipsImage *in, VipsImage *sub, VipsImage **out, WatermarkImageOptions *o) {
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);

	// add in and sub for unreffing and later use
	t[0] = in;
	t[1] = sub;

	if (has_alpha_channel(in) == 0) {
		vips_add_band(in, &t[0], 255.0);
		// in is no longer in the array and won't be unreffed, so add it at the end
		t[8] = in;
	}

	if (has_alpha_channel(sub) == 0) {
		vips_add_band(sub, &t[1], 255.0);
		// sub is no longer in the array and won't be unreffed, so add it at the end
		t[9] = sub;
	}

	// Place watermark image in the right place and size it to the size of the
	// image that should be watermarked
	if (
		vips_embed(t[1], &t[2], o->Left, o->Top, t[0]->Xsize, t[0]->Ysize, NULL)) {
			g_object_unref(base);
		return 1;
	}

	// Create a mask image based on the alpha band from the watermark image
	// and place it in the right position
	if (
		vips_extract_band(t[1], &t[3], t[1]->Bands - 1, "n", 1, NULL) ||
		vips_linear1(t[3], &t[4], o->Opacity, 0.0, NULL) ||
		vips_cast(t[4], &t[5], VIPS_FORMAT_UCHAR, NULL) ||
		vips_copy(t[5], &t[6], "interpretation", t[0]->Type, NULL) ||
		vips_embed(t[6], &t[7], o->Left, o->Top, t[0]->Xsize, t[0]->Ysize, NULL))	{
			g_object_unref(base);
		return 1;
	}

	// Blend the mask and watermark image and write to output.
	if (vips_ifthenelse(t[7], t[2], t[0], out, "blend", TRUE, NULL)) {
		g_object_unref(base);
		return 1;
	}

	if (!t[8]) {
		t[0] = NULL;
	} else {
		t[8] = NULL;
	}

	if (!t[9]) {
		t[1] = NULL;
	} else {
		t[9] = NULL;
	}

	g_object_unref(base);
	return 0;
}


int get_text(VipsImage **text, LabelOptions *o) {
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 1);

	int ret = vips_text(&t[0], o->Text,
		"font", o->Font,
		"dpi", o->DPI,
		"width", o->Width,
		"align", o->Align,
		NULL);
	
	if (ret){
		return 1;
	}
	
	*text = t[0];
	t[0] = NULL;
	g_object_unref(base);
	return 0;
}

int get_text1(VipsImage *in, VipsImage **out, LabelOptions *o) {
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 5);

	t[1] = in;
	if (
		vips_linear1(t[1], &t[2], o->Opacity, 0.0, NULL) ||
		vips_cast(t[2], &t[3], VIPS_FORMAT_UCHAR, NULL) ||
		vips_embed(t[3], &t[4], o->OffsetX, o->OffsetY, t[3]->Xsize + o->OffsetX, t[3]->Ysize + o->OffsetY, NULL)
		) {
		t[1]= NULL;
		g_object_unref(base);
		return 1;
	}
	
	*out = t[4];
	t[1] = NULL;
	t[4] = NULL;
	g_object_unref(base);
	return 0;
}

int watermarkText(VipsImage *in, VipsImage *ti, VipsImage **out, LabelOptions *o){
	double ones[3] = { 1, 1, 1 };
	VipsImage *base = vips_image_new();
	VipsImage **t = (VipsImage **) vips_object_local_array(VIPS_OBJECT(base), 10);
	t[0] = in;
	t[4] = ti;

	if (
		vips_black(&t[5], 1, 1, NULL) ||
		vips_linear(t[5], &t[6], ones, o->Color, 3, NULL) ||
		vips_cast(t[6], &t[7], VIPS_FORMAT_UCHAR, NULL) ||
		vips_copy(t[7], &t[8], "interpretation", t[0]->Type, NULL) ||
		vips_embed(t[8], &t[9], 0, 0, t[0]->Xsize, t[0]->Ysize, "extend", VIPS_EXTEND_COPY, NULL)
		) {
		t[0] = NULL;
		t[4] = NULL;
		g_object_unref(base);
		return 1;
	}

	if (t[0]->Bands != t[9]->Bands) {
		vips_add_band(t[9], &t[1], 255.0);
	} else {
		t[1] = t[9];
	}

	if (vips_ifthenelse(t[4], t[1], t[0], out, "blend", TRUE, NULL)) {
		t[0] = NULL;
		t[4] = NULL;
		if(t[1] == t[9])
			t[1] = NULL;
		g_object_unref(base);
		return 1;
	}
	
	if(t[1] == t[9])
		t[1] = NULL;
	t[0] = NULL;
	t[4] = NULL;
	g_object_unref(base);
	return 0;
}