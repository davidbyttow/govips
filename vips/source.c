#include "source.h"

GoSourceArguments * create_go_source_arguments( int image_id )
{
	GoSourceArguments * source_args;
	source_args = malloc(sizeof(GoSourceArguments));
	source_args->image_id = image_id;

	return source_args;
}

static gint64 go_read ( VipsSourceCustom *source_custom, void *buffer, gint64 length, GoSourceArguments * source_args )
{
    return goSourceRead(source_args->image_id, buffer, length);
}

static gint64 go_seek ( VipsSourceCustom *source_custom, gint64 offset, int whence, GoSourceArguments * source_args )
{
	return goSourceSeek(source_args->image_id, offset, whence);
}

VipsSourceCustom * create_go_custom_source( GoSourceArguments * source_args )
{
	VipsSourceCustom * source_custom = vips_source_custom_new();

	g_signal_connect( source_custom, "read", G_CALLBACK(go_read), source_args );
	g_signal_connect( source_custom, "seek", G_CALLBACK(go_seek), source_args );

	return source_custom;
}
