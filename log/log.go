package log

import (
	"fmt"
)

const (
	infoFmt    = "\x1b[32;1mINFO: %s\x1b[0m\n"
	errorFmt   = "\x1b[31;1mERROR: %s\x1b[0m\n"
	warningFmt = "\x1b[33;1mWARN: %s\x1b[0m\n"
	alertFmt   = "\x1b[35;1mWARN: %s\x1b[0m\n"
)

func PrintError(msg string, args ...interface{}) {
	fmt.Printf(errorFmt, fmt.Sprintf(msg, args...))
}

func PrintWarning(msg string, args ...interface{}) {
	fmt.Printf(warningFmt, fmt.Sprintf(msg, args...))
}

func PrintInfo(msg string, args ...interface{}) {
	fmt.Printf(infoFmt, fmt.Sprintf(msg, args...))
}

func PrintAlert(msg string, args ...interface{}) {
	fmt.Printf(alertFmt, fmt.Sprintf(msg, args...))
}
