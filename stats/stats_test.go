package stats

import (
	"fmt"
	"testing"
)

func TestSplitLogLine(t *testing.T) {
	const logLn = `219.22.67.80 - langworth5892 [04/07/2019:20:01:16 -0500] "POST /optimize/recontextualize HTTP/1.0" 204 25818`

	for _, ln := range splitLogLine(logLn) {
		// TODO assert output
		fmt.Println(ln)
	}
}
