#include <stdlib.h>
#include <vips/vips.h>

void clear_source(VipsSourceCustom **image);
void clear_target(VipsTargetCustom **image);


VipsSourceCustom * create_go_custom_source( void * source_ptr );
static gint64 go_read ( VipsSourceCustom *source_custom, gpointer buffer, gint64 length, gpointer source_ptr );
static gint64 go_seek ( VipsSourceCustom *source_custom, gint64 offset, int whence, gpointer source_ptr );

VipsTargetCustom * create_go_custom_target( void * target_ptr );
static gint64 go_write ( VipsTargetCustom *target_custom, void *data, gint64 length, void *target_ptr );
