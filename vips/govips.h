#include <stdlib.h>
#include <glib.h>
#include <vips/vips.h>

#if (VIPS_MAJOR_VERSION < 8)
error_requires_version_8
#endif

    extern void
    govipsLoggingHandler(char *log_domain, int log_level, char *message);

static void govips_logging_handler(
    const gchar *log_domain, GLogLevelFlags log_level,
    const gchar *message, gpointer user_data);

void vips_set_logging_handler(void);