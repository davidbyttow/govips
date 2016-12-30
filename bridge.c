
#include "bridge.h"

void set_property(VipsObject *object, const char *name, const GValue *value) {
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

void filename_split8(const char* name, char *filename, char *option_string) {
  vips__filename_split8(name, filename, option_string);
}
