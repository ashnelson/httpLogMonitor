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
		// The total count of hits for each interval
		totalHitsPerInterval int32
		// The number of times the alert threshold was exceeded for the interval
		rateExceededCnt int32
	)

	// Monitor the rate per second
	go func() {
		var cnt int32

		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			// Load total traffic and reset for next interval
			cnt = atomic.LoadInt32(&trafficCount)
			atomic.StoreInt32(&trafficCount, int32(0))

			if int(cnt) > alertThreshold {
				atomic.AddInt32(&rateExceededCnt, int32(1))
				atomic.AddInt32(&totalHitsPerInterval, cnt)
			} else {
				atomic.StoreInt32(&totalHitsPerInterval, int32(0))
			}
		}
	}()

	// Monitor the rate over the interval time period
	go func() {
		const alertTimeFmt = "01-02-2006 15:04:05"
		var (
			curTime            string
			needToRecover      bool
			rateCnt, totalHits int32
		)

		ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
		for range ticker.C {
			curTime = time.Now().Format(alertTimeFmt)
			rateCnt = atomic.LoadInt32(&rateExceededCnt)
			totalHits = atomic.LoadInt32(&totalHitsPerInterval)

			if int(rateCnt) >= intervalSeconds {
				needToRecover = true
				log.LogAlert("High traffic generated an alert; nbrHits: %d, triggeredAt: %s", totalHits, curTime)
			} else if needToRecover {
				needToRecover = false
				log.LogAlert("Alert recovered at %s", curTime)
			}

			// Reset the count for the next interval
			atomic.StoreInt32(&rateExceededCnt, int32(0))
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
