package stats

import (
	"bytes"
	"io"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ashnelson/httpLogMonitor/log"
)

func TestMonitorAlerts(t *testing.T) {
	const testTrafficRate = 5

	// Set log output to something that can be read from
	logBuf := bytes.NewBuffer([]byte{})
	log.SetLogOutput(logBuf)

	// "Preload" traffic
	atomic.StoreInt32(&trafficCount, int32(testTrafficRate))

	// Monitor alerts at 1 second intervals for a 1 req/s rate
	monitorAlerts(1, 1)

	t.Run("Generate alerts without recovering", func(t *testing.T) {
		// Generate "traffic" every 100 ms for a total of 2.5 seconds
		generateTraffic(25, testTrafficRate, 100*time.Millisecond)

		// Read log before the monitor can add anything else
		logOutput := logBuf.String()

		// Read log lines and tally alerts generated
		var nbrAlertsGenerated int
		for _, msg := range strings.Split(logOutput, "\n") {
			// Count the number of alerts
			if strings.Contains(msg, "High traffic generated an alert") {
				nbrAlertsGenerated++
				continue
			}
			// Check if alert recovered unexpectedly
			if strings.Contains(msg, "Alert recovered at") {
				t.Fatalf("Monitor unexpectedly recovered from high traffic")
			}
		}

		const expectedNbrAlerts = 2
		if nbrAlertsGenerated != expectedNbrAlerts {
			t.Fatalf("Incorrect number of alerts generated\n\tExpected: %d\n\tActual:   %d", expectedNbrAlerts, nbrAlertsGenerated)
		}
	})

	t.Run("Generate alerts and recover", func(t *testing.T) {
		// Give the monitor some time to recvoer from last test
		time.Sleep(250 * time.Millisecond)

		// Clear the logs
		logBuf.Reset()

		// Generate "traffic" every 100 ms for a total of .5 seconds
		generateTraffic(5, testTrafficRate, 100*time.Millisecond)

		// Give the monitor enough time to see that no more traffic is generated
		time.Sleep(1500 * time.Millisecond)

		var (
			alertRecovered bool
			msg            string
			err            error
		)

		// Read log lines and tally alerts generated
		for {
			if msg, err = logBuf.ReadString('\n'); err == io.EOF {
				break
			}
			// Check if alert recovered unexpectedly
			if strings.Contains(msg, "Alert recovered at") {
				alertRecovered = true
				break
			}
		}

		if !alertRecovered {
			t.Fatalf("Alert failed to recover after traffic trended below alert threshold")
		}
	})
}

func generateTraffic(nbrRuns, trafficRate int, tickerInterval time.Duration) {
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	for range ticker.C {
		// Generate "traffic"
		atomic.AddInt32(&trafficCount, int32(trafficRate))

		if i > nbrRuns {
			ticker.Stop()
			break
		}
		i++
	}
}
