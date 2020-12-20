package vips

import (
	"bytes"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func Test_DefaultLogging(t *testing.T) {
	Startup(nil)
	if os.Getenv("CI") != "" {
		enableLogging()
	}
	LoggingSettings(nil, LogLevelInfo)

	output := captureOutput(func() {
		govipsLog("test", LogLevelInfo, "abcde12345")
	})
	dateRegexp := regexp.MustCompile(`[0-9\/\:]+ `)
	output = dateRegexp.ReplaceAllString(output, "")
	assert.Equal(t, "[test.info] abcde12345\n", output)
	if os.Getenv("CI") != "" {
		consoleLogging()
	}
}

func Test_LoggingVerbosity(t *testing.T) {
	Startup(nil)
	if os.Getenv("CI") != "" {
		enableLogging()
	}
	LoggingSettings(nil, LogLevelMessage)

	output := captureOutput(func() {
		govipsLog("test", LogLevelMessage, "abcde12345")
	})
	dateRegexp := regexp.MustCompile(`[0-9\/\:]+ `)
	output = dateRegexp.ReplaceAllString(output, "")
	assert.Equal(t, "[test.message] abcde12345\n", output)

	output2 := captureOutput(func() {
		govipsLog("test", LogLevelInfo, "fghji67890")
	})
	assert.Equal(t, "", output2)
	if os.Getenv("CI") != "" {
		consoleLogging()
	}
}

func Test_LoggingHandler(t *testing.T) {
	Startup(nil)
	if os.Getenv("CI") != "" {
		enableLogging()
	}

	var testDomain string
	var testVerbosity LogLevel
	var testMessage string
	testHandler := func(domain string, verbosity LogLevel, message string) {
		testDomain = domain
		testVerbosity = verbosity
		testMessage = message
	}
	LoggingSettings(testHandler, LogLevelInfo)

	govipsLog("domain", LogLevelCritical, "abcde12345")
	assert.Equal(t, "domain", testDomain)
	assert.Equal(t, LogLevelCritical, testVerbosity)
	assert.Equal(t, "abcde12345", testMessage)
	if os.Getenv("CI") != "" {
		consoleLogging()
	}
}
