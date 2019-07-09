package stats

import (
	"sync/atomic"
	"time"

	"github.com/ashnelson/httpLogMonitor/config"
	"github.com/ashnelson/httpLogMonitor/log"
)

var trafficCount int32

// initializeMonitors starts the stats and alerts monitors each in their own
// goroutine
func initializeMonitors(cfg config.StatsConfig) {
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
			for name, curSection := range recordedTraffic.sections {
				atomic.AddInt32(&trafficCount, int32(curSection.totalCount))
				printSection(name, curSection)
			}

			// Reset the traffic for the next interval
			recordedTraffic.sections = make(map[string]*section, len(recordedTraffic.sections))
			recordedTraffic.lock.Unlock()
		}
	}()
}

// monitorAlerts monitors the total traffic for the specified interval and
// prints an alert if the traffic exceeds the specified threshold. An alert will
// be printed for each interval until the traffic trends back below the
// specified threshold in which case a recovery message will be printed.
func monitorAlerts(intervalSeconds, alertThreshold int) {
	var (
		totalHitsPerInterval int32
		alertTriggered       bool
	)

	// Monitor the rate per second
	go func() {
		var cnt int32

		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			cnt = atomic.LoadInt32(&trafficCount)

			if int(cnt) > alertThreshold {
				alertTriggered = true
				atomic.AddInt32(&totalHitsPerInterval, cnt)
			} else if alertTriggered {
				alertTriggered = false
				atomic.StoreInt32(&totalHitsPerInterval, int32(0))
			}

			atomic.StoreInt32(&trafficCount, int32(0))
		}
	}()

	// Monitor the rate over the interval time period
	go func() {
		const alertTimeFmt = "01-02-2006 15:04:05"
		var (
			curTime       string
			needToRecover bool
		)

		ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
		for range ticker.C {
			curTime = time.Now().Format(alertTimeFmt)

			if alertTriggered {
				needToRecover = true
				log.LogAlert("High traffic generated an alert; nbrHits: %d, triggeredAt: %s", totalHitsPerInterval, curTime)
			} else if needToRecover {
				needToRecover = false
				log.LogAlert("Alert recovered at %s", curTime)
			}
		}
	}()
}

// printSection prints the section details as formatted string
func printSection(sectionName string, curSection *section) {
	log.Log(`%s
	TotalHits:    %d
	TotalBytes:   %d
	Methods:      %v
	Users:        %v
	RemoteHosts:  %v
	StatusCodes:  %v
`, sectionName, curSection.totalCount, curSection.totalBytes, curSection.methods,
		curSection.users, curSection.remoteHosts, curSection.statusCode)
}
