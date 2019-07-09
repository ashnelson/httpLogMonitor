package log

import (
	"fmt"
	"io"
	"os"
)

const (
	infoFmt    = "\x1b[32;1mINFO: %s\x1b[0m\n"
	errorFmt   = "\x1b[31;1mERROR: %s\x1b[0m\n"
	warningFmt = "\x1b[33;1mWARN: %s\x1b[0m\n"
	alertFmt   = "\x1b[35;1mALERT: %s\x1b[0m\n"
)

var logWtr io.Writer = os.Stdout

// SetLogOutput sets the logged output to write to the specified writer
func SetLogOutput(wtr io.Writer) {
	logWtr = wtr
}

// Log logs the message with no tag, coloring, or formatting
func Log(msg string, args ...interface{}) {
	fmt.Fprintf(logWtr, msg, args...)
}

// LogError logs teh message with an "ERROR:" tag with red text
func LogError(msg string, args ...interface{}) {
	fmt.Fprintf(logWtr, errorFmt, fmt.Sprintf(msg, args...))
}

// LogWarning logs teh message with a "WARN:" tag with yellow text
func LogWarning(msg string, args ...interface{}) {
	fmt.Fprintf(logWtr, warningFmt, fmt.Sprintf(msg, args...))
}

// LogInfo logs teh message with an "INFO:" tag with green text
func LogInfo(msg string, args ...interface{}) {
	fmt.Fprintf(logWtr, infoFmt, fmt.Sprintf(msg, args...))
}

// LogAlert logs teh message with an "ALERT:" tag with magenta text
func LogAlert(msg string, args ...interface{}) {
	fmt.Fprintf(logWtr, alertFmt, fmt.Sprintf(msg, args...))
}
