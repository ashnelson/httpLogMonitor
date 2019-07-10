package stats

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ashnelson/httpLogMonitor/config"
	"github.com/ashnelson/httpLogMonitor/log"
)

// The indexes for accessing the log line after it's split into a slice
const (
	remoteHostsIdx = iota
	rfc931Idx
	authUserIdx
	dateIdx
	requestIdx
	statusIdx
	bytesIdx
)

var (
	logCh           chan string
	recordedTraffic sections
)

// Sections contains a map of all of the sections with the key being the section
// name and the value the section object
type sections struct {
	sections map[string]*section
	lock     sync.Mutex
}

// section contains all of the information about a specific section of the logs
type section struct {
	totalCount  int
	totalBytes  int
	methods     map[string]int
	users       map[string]int
	remoteHosts map[string]int
	statusCode  map[string]int

	lock sync.Mutex
}

// InitializeRecorder takes a StatsConfig and initializes the log chan and
// recordedTraffic object and starts a goroutine to read log lines from the
// log chan
func InitializeRecorder(cfg config.StatsConfig) {
	logCh = make(chan string, 128)
	recordedTraffic = sections{
		sections: make(map[string]*section, 11),
	}

	go logRecorder()

	initializeMonitors(cfg)
}

// RecordLog processes and records the provided log line
func RecordLog(logLn string) {
	logLn = strings.TrimSpace(logLn)
	if logLn != "" {
		logCh <- logLn
	}
}

// Shutdown handles all shutdown functionality such as closing the log channel
func Shutdown() {
	close(logCh)
}

// logRecorder decodes the log line and records the log in the Sections map.
// If the Section doesn't already exist in the Sections map it will be created.
func logRecorder() {
	var logLn string

	for {
		select {
		case logLn = <-logCh:
			go recordLogLine(logLn)
		}
	}

	return
}

// recordLogLine gets the correct section from recordedTraffic or creates one if
// one doesn't already exist and updates the section data with the info from the
// provided log line
func recordLogLine(logLn string) {
	// Increment traffic count for alert monitoring
	atomic.AddInt32(&trafficCount, int32(1))

	logFields := splitLogLine(logLn)
	sectionName := getSectionName(logFields[requestIdx])

	// Safely get the section
	recordedTraffic.lock.Lock()
	curSection, ok := recordedTraffic.sections[sectionName]
	if !ok {
		curSection = &section{
			methods:     make(map[string]int, 3),
			users:       make(map[string]int, 11),
			remoteHosts: make(map[string]int, 11),
			statusCode:  make(map[string]int, 5),
		}
		recordedTraffic.sections[sectionName] = curSection
	}
	recordedTraffic.lock.Unlock()

	// Get the log line as a slice
	reqBytes, err := strconv.Atoi(logFields[bytesIdx])
	if err != nil {
		log.LogError("Failed to cast request bytes (%s) to integer; Details: %s", logFields[bytesIdx], err)
	}

	// Safely update the section
	curSection.lock.Lock()
	defer curSection.lock.Unlock()

	curSection.totalCount++
	curSection.totalBytes = curSection.totalBytes + reqBytes
	curSection.methods[getRequestMethod(logFields[requestIdx])]++
	curSection.users[logFields[authUserIdx]]++
	curSection.remoteHosts[logFields[remoteHostsIdx]]++
	curSection.statusCode[logFields[statusIdx]]++
}

// splitLogLine parses the log line into the separate fields and returns them
// as a slice of strings
func splitLogLine(logLn string) []string {
	const nbrFields = 7

	strSep := " "
	idxOfStrSep := 0
	logFields := make([]string, 0, nbrFields)

	for i := 0; i < 7; i++ {
		// Check if the date field was hit and change the separator to keep
		// date and request fields intact
		if logLn[0] == '[' {
			strSep = `"`
		}

		// If the remaining log line doesn't contain a double quote character
		// the only fields left are the status code and bytes
		if !strings.Contains(logLn, `"`) {
			logFields = append(logFields, strings.Split(logLn[1:], " ")...)
			break
		}

		// Get the index of the string separator
		if idxOfStrSep = strings.Index(logLn, strSep); idxOfStrSep < 0 {
			break
		}

		// Append the log field to the slice
		logFields = append(logFields, strings.TrimSpace(logLn[:idxOfStrSep]))
		logLn = logLn[idxOfStrSep+1:]
	}

	return logFields
}

// getSectionName parses the reqStr and returns the section name
//		ex: reqStr = "GET /foo/bar HTTP/1.0"
//			section name = "foo"
func getSectionName(reqStr string) string {
	splitReq := strings.Split(reqStr, " ")
	if len(splitReq) < 3 {
		return ""
	}

	splitPath := strings.Split(splitReq[1], "/")
	if len(splitPath) < 2 {
		return ""
	}

	return "/" + splitPath[1]
}

// getRequestMethod parses the reqStr and returns the HTTP method string
func getRequestMethod(reqStr string) string {
	splitReq := strings.Split(reqStr, " ")
	if len(splitReq) < 1 {
		return ""
	}

	return splitReq[0]
}
