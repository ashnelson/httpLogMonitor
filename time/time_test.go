package main

import (
	"testing"
)

func TestParseTime(t *testing.T) {
	const expectedUnixNano int64 = 1562165240000000000

	t.Run("DataDog example log time", func(t *testing.T) {
		const dataDogLogTime = "[03/Jul/2019:14:47:20 +0000]"

		parsedTime, err := parseTime(dataDogLogTime)
		if err != nil {
			t.Fatalf("Failed to parse DataDog timestamp %q\n\tDetails: %s", dataDogLogTime, err)
		}

		actualUnixNano := parsedTime.UnixNano()
		if actualUnixNano != expectedUnixNano {
			t.Fatalf("Incorrect unix nano time returned for parsed time\n\tExpected: %d\n\tActual:   %d", expectedUnixNano, actualUnixNano)
		}
	})

	t.Run("flog log time", func(t *testing.T) {
		const flogLogTime = "[03/07/2019:14:47:20 -0000]"

		parsedTime, err := parseTime(flogLogTime)
		if err != nil {
			t.Fatalf("Failed to parse DataDog timestamp %q\n\tDetails: %s", flogLogTime, err)
		}

		actualUnixNano := parsedTime.UnixNano()
		if actualUnixNano != expectedUnixNano {
			t.Fatalf("Incorrect unix nano time returned for parsed time\n\tExpected: %d\n\tActual:   %d", expectedUnixNano, actualUnixNano)
		}
	})
}
