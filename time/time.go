package main

import (
	"fmt"
	"strings"
	"time"
)

func parseTime(timeStr string) (*time.Time, error) {
	timeFmt := "[02/Jan/2006:15:04:05 +0000]" // DataDog example log time
	if strings.Contains(timeStr, "-") {
		timeFmt = "[02/01/2006:15:04:05 -0700]" // flog log time
	}

	t, err := time.Parse(timeFmt, timeStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse time %q; Details: %s", timeStr, err)
	}

	return &t, nil
}
