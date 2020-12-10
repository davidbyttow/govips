#include <vips/vips.h>

typedef struct _GoSourceArguments {
	int image_id;
} GoSourceArguments;

GoSourceArguments * create_go_source_arguments( int image_id );
VipsSourceCustom * create_go_custom_source( GoSourceArguments * source_args );

static gint64 go_read ( VipsSourceCustom *source_custom, void *buffer, gint64 length, GoSourceArguments * source_args );
static gint64 go_seek ( VipsSourceCustom *source_custom, gint64 offset, int whence, GoSourceArguments * source_args );
