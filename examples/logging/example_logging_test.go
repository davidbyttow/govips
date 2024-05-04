package package_test

import (
	"fmt"
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

/* govips example: how to modify logging by using your own log handler
 *
 * Defines a myLogger() function compatible with govips logging,
 * which uses log.Printf() to print all logs to stderr.
 * Sets myLogger as the logging handler and Info as the minimum logging
 * verbosity *before* starting govips so startup-related messages already
 * go to the logger instead of stderr.
 *
 * The program also uses the myLogger() function directly for its own
 * logging purposes, by defining the logging domain and making it easy
 * to see which program is creating the log messages.
 */

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func myLogger(messageDomain string, verbosity vips.LogLevel, message string) {
	var messageLevelDescription string
	switch verbosity {
	case vips.LogLevelError:
		messageLevelDescription = "error"
	case vips.LogLevelCritical:
		messageLevelDescription = "critical"
	case vips.LogLevelWarning:
		messageLevelDescription = "warning"
	case vips.LogLevelMessage:
		messageLevelDescription = "message"
	case vips.LogLevelInfo:
		messageLevelDescription = "info"
	case vips.LogLevelDebug:
		messageLevelDescription = "debug"
	}

	log.Printf("[%v.%v] %v", messageDomain, messageLevelDescription, message)
}

func main() {
	vips.LoggingSettings(myLogger, vips.LogLevelInfo)
	vips.Startup(nil)
	defer vips.Shutdown()

	image1, err := vips.NewImageFromFile("input.jpg")
	checkError(err)

	myLogger("loggingExample", vips.LogLevelInfo, fmt.Sprintf("Before: %v x %v", image1.Height(), image1.Width()))
	err = image1.Resize(0.9, vips.KernelAuto)
	checkError(err)
	myLogger("loggingExample", vips.LogLevelInfo, fmt.Sprintf("Before: %v x %v", image1.Height(), image1.Width()))

	image1bytes, _, err := image1.ExportJpeg(nil)
	checkError(err)
	err = os.WriteFile("output.jpg", image1bytes, 0644)
	checkError(err)
}
