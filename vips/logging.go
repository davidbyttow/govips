package vips

// #cgo pkg-config: vips
// #include <glib.h>
import "C"
import "log"

type logLevel int

const (
	logLevelError    logLevel = C.G_LOG_LEVEL_ERROR
	logLevelCritical logLevel = C.G_LOG_LEVEL_CRITICAL
	logLevelWarning  logLevel = C.G_LOG_LEVEL_WARNING
	logLevelMessage  logLevel = C.G_LOG_LEVEL_MESSAGE
	logLevelInfo     logLevel = C.G_LOG_LEVEL_INFO
	logLevelDebug    logLevel = C.G_LOG_LEVEL_DEBUG
)

//export govipsLoggingHandler
func govipsLoggingHandler(messageDomain *C.char, messageLevel C.int, message *C.char) {
	govipsLog(C.GoString(messageDomain), logLevel(messageLevel), C.GoString(message))
}

func govipsLog(messageDomain string, messageLevel logLevel, message string) {
	var messageLevelDescription string
	switch logLevel(messageLevel) {
	case logLevelError:
		messageLevelDescription = "error"
	case logLevelCritical:
		messageLevelDescription = "critical"
	case logLevelWarning:
		messageLevelDescription = "warning"
	case logLevelMessage:
		messageLevelDescription = "message"
	case logLevelInfo:
		messageLevelDescription = "info"
	case logLevelDebug:
		messageLevelDescription = "debug"
	}

	log.Println(messageDomain, "[", messageLevelDescription, "]", message)
}
