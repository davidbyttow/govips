#include "introspect.h"
#include <string.h>
#include <stdio.h>

// Map a GType + GParamSpec to our ArgType enum values.
// These must match the ArgType constants in types.go.
enum {
    ARG_TYPE_UNKNOWN = 0,
    ARG_TYPE_IMAGE,
    ARG_TYPE_DOUBLE,
    ARG_TYPE_INT,
    ARG_TYPE_BOOL,
    ARG_TYPE_STRING,
    ARG_TYPE_ENUM,
    ARG_TYPE_FLAGS,
    ARG_TYPE_ARRAY_DOUBLE,
    ARG_TYPE_ARRAY_INT,
    ARG_TYPE_ARRAY_IMAGE,
    ARG_TYPE_BLOB,
    ARG_TYPE_INTERPOLATE,
    ARG_TYPE_SOURCE,
    ARG_TYPE_TARGET,
};

// Map argument flags to our ArgFlags bitmask.
// These must match the ArgFlags constants in types.go.
enum {
    ARG_FLAG_INPUT    = 1 << 0,
    ARG_FLAG_OUTPUT   = 1 << 1,
    ARG_FLAG_REQUIRED = 1 << 2,
    ARG_FLAG_MODIFY   = 1 << 3,
};

static int classify_gtype(GType type) {
    if (g_type_is_a(type, VIPS_TYPE_IMAGE))
        return ARG_TYPE_IMAGE;
    if (g_type_is_a(type, VIPS_TYPE_INTERPOLATE))
        return ARG_TYPE_INTERPOLATE;
    if (g_type_is_a(type, VIPS_TYPE_SOURCE))
        return ARG_TYPE_SOURCE;
    if (g_type_is_a(type, VIPS_TYPE_TARGET))
        return ARG_TYPE_TARGET;
    if (g_type_is_a(type, VIPS_TYPE_BLOB))
        return ARG_TYPE_BLOB;
    if (g_type_is_a(type, VIPS_TYPE_ARRAY_DOUBLE))
        return ARG_TYPE_ARRAY_DOUBLE;
    if (g_type_is_a(type, VIPS_TYPE_ARRAY_INT))
        return ARG_TYPE_ARRAY_INT;
    if (g_type_is_a(type, VIPS_TYPE_ARRAY_IMAGE))
        return ARG_TYPE_ARRAY_IMAGE;
    if (G_TYPE_IS_ENUM(type))
        return ARG_TYPE_ENUM;
    if (G_TYPE_IS_FLAGS(type))
        return ARG_TYPE_FLAGS;

    GType fundamental = G_TYPE_FUNDAMENTAL(type);
    if (fundamental == G_TYPE_DOUBLE || fundamental == G_TYPE_FLOAT)
        return ARG_TYPE_DOUBLE;
    if (fundamental == G_TYPE_INT || fundamental == G_TYPE_UINT ||
        fundamental == G_TYPE_INT64 || fundamental == G_TYPE_UINT64 ||
        fundamental == G_TYPE_LONG || fundamental == G_TYPE_ULONG)
        return ARG_TYPE_INT;
    if (fundamental == G_TYPE_BOOLEAN)
        return ARG_TYPE_BOOL;
    if (fundamental == G_TYPE_STRING)
        return ARG_TYPE_STRING;

    return ARG_TYPE_UNKNOWN;
}

static int convert_flags(VipsArgumentFlags vflags) {
    int flags = 0;
    if (vflags & VIPS_ARGUMENT_INPUT)
        flags |= ARG_FLAG_INPUT;
    if (vflags & VIPS_ARGUMENT_OUTPUT)
        flags |= ARG_FLAG_OUTPUT;
    if (vflags & VIPS_ARGUMENT_REQUIRED)
        flags |= ARG_FLAG_REQUIRED;
    if (vflags & VIPS_ARGUMENT_MODIFY)
        flags |= ARG_FLAG_MODIFY;
    return flags;
}

// Callback for vips_argument_map to collect argument info.
typedef struct {
    OpInfo *op_info;
} ArgMapData;

static void *collect_args(VipsObject *object, GParamSpec *pspec,
                          VipsArgumentClass *argument_class,
                          VipsArgumentInstance *argument_instance,
                          void *a, void *b) {
    ArgMapData *data = (ArgMapData *)a;
    OpInfo *op = data->op_info;

    // Skip deprecated arguments.
    if (argument_class->flags & VIPS_ARGUMENT_DEPRECATED)
        return NULL;

    // Skip non-construct arguments (internal vips bookkeeping).
    if (!(argument_class->flags & VIPS_ARGUMENT_CONSTRUCT))
        return NULL;

    if (op->n_args >= MAX_ARGS)
        return NULL;

    ArgInfo *arg = &op->args[op->n_args];
    memset(arg, 0, sizeof(ArgInfo));

    strncpy(arg->name, g_param_spec_get_name(pspec), sizeof(arg->name) - 1);

    GType type = G_PARAM_SPEC_VALUE_TYPE(pspec);
    arg->type = classify_gtype(type);
    arg->flags = convert_flags(argument_class->flags);
    arg->priority = argument_class->priority;

    // Extract enum type name.
    if (arg->type == ARG_TYPE_ENUM || arg->type == ARG_TYPE_FLAGS) {
        strncpy(arg->enum_type, g_type_name(type), sizeof(arg->enum_type) - 1);
    }

    // Extract default/min/max for numeric types.
    if (G_IS_PARAM_SPEC_DOUBLE(pspec)) {
        GParamSpecDouble *dspec = G_PARAM_SPEC_DOUBLE(pspec);
        arg->defval = dspec->default_value;
        arg->min = dspec->minimum;
        arg->max = dspec->maximum;
    } else if (G_IS_PARAM_SPEC_INT(pspec)) {
        GParamSpecInt *ispec = G_PARAM_SPEC_INT(pspec);
        arg->defval = ispec->default_value;
        arg->min = ispec->minimum;
        arg->max = ispec->maximum;
    } else if (G_IS_PARAM_SPEC_UINT(pspec)) {
        GParamSpecUInt *uspec = G_PARAM_SPEC_UINT(pspec);
        arg->defval = uspec->default_value;
        arg->min = uspec->minimum;
        arg->max = uspec->maximum;
    } else if (G_IS_PARAM_SPEC_BOOLEAN(pspec)) {
        GParamSpecBoolean *bspec = G_PARAM_SPEC_BOOLEAN(pspec);
        arg->defval = bspec->default_value ? 1.0 : 0.0;
    } else if (G_IS_PARAM_SPEC_ENUM(pspec)) {
        GParamSpecEnum *espec = G_PARAM_SPEC_ENUM(pspec);
        arg->defval = espec->default_value;
    }

    op->n_args++;
    return NULL;
}

// Callback for vips_type_map_all to collect operations.
static void *collect_ops(GType type, void *user_data) {
    IntrospectResult *result = (IntrospectResult *)user_data;

    // Only include concrete (instantiable) types.
    if (G_TYPE_IS_ABSTRACT(type))
        return NULL;

    if (result->n_ops >= MAX_OPS)
        return NULL;

    // Try to create an instance to introspect its arguments.
    const char *name = vips_nickname_find(type);
    if (!name)
        return NULL;

    VipsObject *obj = (VipsObject *)g_object_new(type, NULL);
    if (!obj)
        return NULL;

    VipsObjectClass *oclass = VIPS_OBJECT_GET_CLASS(obj);
    if (!oclass) {
        g_object_unref(obj);
        return NULL;
    }

    OpInfo *op = &result->ops[result->n_ops];
    memset(op, 0, sizeof(OpInfo));

    strncpy(op->name, name, sizeof(op->name) - 1);

    if (oclass->description)
        strncpy(op->description, oclass->description, sizeof(op->description) - 1);

    // Walk up the GType hierarchy to find the category. We look for a known
    // abstract category GType in the ancestry. Many operations are direct
    // children of VipsOperation, so we also check the direct parent name.
    {
        const char *found_category = NULL;
        GType walk = g_type_parent(type);
        while (walk && walk != VIPS_TYPE_OPERATION && walk != VIPS_TYPE_OBJECT) {
            const char *tname = g_type_name(walk);
            if (tname) {
                // Check if this is a known abstract category.
                if (strcmp(tname, "VipsArithmetic") == 0 ||
                    strcmp(tname, "VipsBinary") == 0 ||
                    strcmp(tname, "VipsUnary") == 0 ||
                    strcmp(tname, "VipsStatistic") == 0) {
                    found_category = "arithmetic";
                    // Keep walking; a more specific parent may exist.
                } else if (strcmp(tname, "VipsColour") == 0 ||
                           strcmp(tname, "VipsColourCode") == 0 ||
                           strcmp(tname, "VipsColourDifference") == 0 ||
                           strcmp(tname, "VipsColourSpace") == 0 ||
                           strcmp(tname, "VipsColourTransform") == 0) {
                    found_category = "colour";
                } else if (strcmp(tname, "VipsConversion") == 0) {
                    found_category = "conversion";
                } else if (strcmp(tname, "VipsConvolution") == 0) {
                    found_category = "convolution";
                } else if (strcmp(tname, "VipsCreate") == 0) {
                    found_category = "create";
                } else if (strcmp(tname, "VipsDraw") == 0) {
                    found_category = "draw";
                } else if (strcmp(tname, "VipsForeign") == 0 ||
                           strcmp(tname, "VipsForeignLoad") == 0 ||
                           strcmp(tname, "VipsForeignSave") == 0) {
                    found_category = "foreign";
                } else if (strcmp(tname, "VipsFreqfilt") == 0) {
                    found_category = "freqfilt";
                } else if (strcmp(tname, "VipsHistogram") == 0) {
                    found_category = "histogram";
                } else if (strcmp(tname, "VipsMorphology") == 0) {
                    found_category = "morphology";
                } else if (strcmp(tname, "VipsResample") == 0) {
                    found_category = "resample";
                }
                if (found_category)
                    break;
            }
            walk = g_type_parent(walk);
        }
        if (found_category) {
            strncpy(op->category, found_category, sizeof(op->category) - 1);
        } else {
            // Fallback: use the operation nickname as category.
            strncpy(op->category, name, sizeof(op->category) - 1);
        }
    }

    // Collect arguments.
    ArgMapData data = { .op_info = op };
    vips_argument_map(obj, collect_args, &data, NULL);

    result->n_ops++;
    g_object_unref(obj);

    return NULL;
}

int vipsgen_introspect(IntrospectResult *result) {
    memset(result, 0, sizeof(IntrospectResult));

    if (VIPS_INIT("vipsgen"))
        return -1;

    // Start from VipsOperation and recurse into all subtypes.
    vips_type_map_all(VIPS_TYPE_OPERATION, collect_ops, result);

    return 0;
}

int vipsgen_introspect_enum(const char *type_name, EnumInfo *result) {
    memset(result, 0, sizeof(EnumInfo));

    GType type = g_type_from_name(type_name);
    if (!type || !G_TYPE_IS_ENUM(type))
        return -1;

    GEnumClass *eclass = g_type_class_ref(type);
    if (!eclass)
        return -1;

    strncpy(result->c_name, type_name, sizeof(result->c_name) - 1);

    for (guint i = 0; i < eclass->n_values; i++) {
        if (result->n_values >= MAX_ENUM_VALUES)
            break;

        GEnumValue *v = &eclass->values[i];

        // Skip the "last" sentinel value that vips uses.
        if (v->value_nick && strcmp(v->value_nick, "last") == 0)
            continue;

        EnumValueInfo *ev = &result->values[result->n_values];
        if (v->value_name)
            strncpy(ev->c_name, v->value_name, sizeof(ev->c_name) - 1);
        if (v->value_nick)
            strncpy(ev->nick, v->value_nick, sizeof(ev->nick) - 1);
        ev->value = v->value;
        result->n_values++;
    }

    g_type_class_unref(eclass);
    return 0;
}

int vipsgen_introspect_enums(const char **enum_names, int n, EnumInfo *results) {
    for (int i = 0; i < n; i++) {
        if (vipsgen_introspect_enum(enum_names[i], &results[i]) != 0) {
            // Not all types are enums; just zero out the result.
            memset(&results[i], 0, sizeof(EnumInfo));
            strncpy(results[i].c_name, enum_names[i], sizeof(results[i].c_name) - 1);
        }
    }
    return 0;
}
