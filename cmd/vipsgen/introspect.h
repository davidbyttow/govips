#ifndef VIPSGEN_INTROSPECT_H
#define VIPSGEN_INTROSPECT_H

#include <vips/vips.h>

// Maximum number of operations and arguments we can handle.
#define MAX_OPS 1024
#define MAX_ARGS 64
#define MAX_ENUM_VALUES 128

// Argument info extracted from introspection.
typedef struct {
    char name[256];
    int type;       // maps to ArgType enum in Go
    int flags;      // combination of ArgFlags
    int priority;
    double defval;
    double min;
    double max;
    char enum_type[256]; // GType name for enum args
} ArgInfo;

// Operation info extracted from introspection.
typedef struct {
    char name[256];
    char description[1024];
    char category[256];
    int n_args;
    ArgInfo args[MAX_ARGS];
} OpInfo;

// Enum value info.
typedef struct {
    char c_name[256];
    char nick[256];
    int value;
} EnumValueInfo;

// Enum type info.
typedef struct {
    char c_name[256];
    int n_values;
    EnumValueInfo values[MAX_ENUM_VALUES];
} EnumInfo;

// Results from introspection.
typedef struct {
    int n_ops;
    OpInfo ops[MAX_OPS];
} IntrospectResult;

// Initialize vips and run introspection.
int vipsgen_introspect(IntrospectResult *result);

// Introspect a single enum type by GType name.
int vipsgen_introspect_enum(const char *type_name, EnumInfo *result);

// Introspect all enum types referenced by discovered operations.
// enum_names is an array of GType name strings, n is the count.
// results is an array of EnumInfo structs to fill.
int vipsgen_introspect_enums(const char **enum_names, int n, EnumInfo *results);

#endif
