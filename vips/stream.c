#include "stream.h"

void clear_source(VipsSource **ref) {
    // https://developer.gnome.org/gobject/stable/gobject-The-Base-Object-Type.html#g-clear-object
    if (G_IS_OBJECT(*ref)) g_clear_object(ref);
}

void clear_target(VipsTarget **ref) {
    // https://developer.gnome.org/gobject/stable/gobject-The-Base-Object-Type.html#g-clear-object
    if (G_IS_OBJECT(*ref)) g_clear_object(ref);
}
