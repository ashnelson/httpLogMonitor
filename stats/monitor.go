package stats

import (
	"fmt"
	"sync/atomic"
	"time"

	"../config"
	"../log"
)

var trafficCount int32

func initializeMonitors(cfg config.StatsCfg) {
	monitorStats(cfg.StatsIntervalSeconds)
	monitorAlerts(cfg.AlertIntervalSeconds, cfg.AlertThreshold)
}

// monitorStats kicks off a goroutine and prints the recordedTraffic for every
// specified interval and then clears the traffic recorder for the next interval
func monitorStats(intervalSeconds int) {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(intervalSeconds))

		for range ticker.C {
			recordedTraffic.lock.Lock()
			// Log all traffic from the current interval
			for _, sctn := range recordedTraffic.sections {
				atomic.AddInt32(&trafficCount, int32(sctn.totalCount))
				fmt.Printf("%+v\n", sctn)
			}

			// Reset the traffic for the next interval
			resetTrafficRecorder()
			recordedTraffic.lock.Unlock()
		}
	}()
}

// monitorAlerts monitors the total traffic for the specified interval and
// prints an alert if the traffic exceeds the specified threshold. An alert will
// be printed for each interval until the traffic trends back below the
// specified threshold in which case a recovery message will be printed.
func monitorAlerts(intervalSeconds, alertThreshold int) {
	go func() {
		const alertTimeFmt = "01-02-2006 15:04:05"
		var (
			cnt     int32
			curTime string
		)

		alertTriggered := false
		ticker := time.NewTicker(time.Second * time.Duration(intervalSeconds))

		for range ticker.C {
			cnt = atomic.LoadInt32(&trafficCount)
			curTime = time.Now().Format(alertTimeFmt)

			if int(cnt) > alertThreshold {
				alertTriggered = true
				log.PrintAlert("High traffic generated an alert\n\thits = {%d}, triggered at {%s}", cnt, curTime)
			} else if alertTriggered {
				log.PrintInfo("Alert recovered at %s", curTime)
				alertTriggered = false
			}

			atomic.StoreInt32(&trafficCount, int32(0))
		}
	}()
}
