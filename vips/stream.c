#include "stream.h"
#include "_cgo_export.h"

void clear_source(VipsSource **ref) {
    // https://developer.gnome.org/gobject/stable/gobject-The-Base-Object-Type.html#g-clear-object
    if (G_IS_OBJECT(*ref)) g_clear_object(ref);
}

void clear_target(VipsTarget **ref) {
    // https://developer.gnome.org/gobject/stable/gobject-The-Base-Object-Type.html#g-clear-object
    if (G_IS_OBJECT(*ref)) g_clear_object(ref);
}


static gint64 go_read ( VipsSourceCustom *source_custom, gpointer buffer, gint64 length, gpointer source_ptr )
{
    return goSourceRead(source_ptr, buffer, length);
}

static gint64 go_seek ( VipsSourceCustom *source_custom, gint64 offset, int whence, gpointer source_ptr )
{
    return goSourceSeek(source_ptr, offset, whence);
}

static gint64 go_write ( VipsTargetCustom *target_custom, gpointer buffer, gint64 length, gpointer target_ptr )
{
    return goTargetWrite(target_ptr, buffer, length);
}


VipsSource * create_go_custom_source( void * source_ptr )
{
    VipsSourceCustom * source_custom = vips_source_custom_new();

    g_signal_connect( source_custom, "read", G_CALLBACK(go_read), source_ptr );
    g_signal_connect( source_custom, "seek", G_CALLBACK(go_seek), source_ptr );

    return (VipsSource *)source_custom;
}

VipsTarget * create_go_custom_target( void * target_ptr )
{
    VipsTargetCustom * target_custom = vips_target_custom_new();

    g_signal_connect( target_custom, "write", G_CALLBACK(go_write), target_ptr );

    return (VipsTarget *)target_custom;
}


