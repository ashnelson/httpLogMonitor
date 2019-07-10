# httpLogMonitor

The httpLogMonitor is a tool used to tail a common log file, parse each log line, update the stats for each URL section for a specified interval and then record the stats at the end of the interval. If the number of log lines processed exceeds the specified alert threshold, an alert will be logged. An alert message will continue to be logged for as long as the threshold is exceeded. Once the processed log lines trend back below the threshold, a recovery message will be printed.

### Config Options

The config is a standard JSON file and is read during startup.
* `inputLogFile` (default: **/tmp/access.log**) - The log file to be read by the application for monitoring
* `shutdownWaitSeconds` (default: **3**) - The amount of time to wait before the application shuts down
* `statsIntervalSeconds` (default: **10**) - The interval to accumulate stats before printing the output
* `alertIntervalSeconds` (default: **120**) - The interval to monitor the alert threshold before an alert is logged
* `alertThreshold` (default: **10**) - The number of log lines processed per second before an alert is generated

### How to run

The only dependency required is the [hpcloud/tail](https://github.com/hpcloud/tail) library and can be installed by running:

```bash
go get github.com/hpcloud/tail/...
```

If the default values need to be overridden, the `config.json` file can be modified and the application ran as follows:

```bash
go build && ./httpLogMonitor
```

The `-h` flag will output the usage text:

```bash
./httpLogMonitor -h
Usage of ./httpLogMonitor:
  -cfg string
    	the config file to use for application startup (default "config.json")
```

### Sample output

```bash
/api
	TotalHits:    3
	TotalBytes:   280
	Methods:      map[POST:2 GET:1]
	Users:        map[jill:1 mary:1 frank:1]
	RemoteHosts:  map[127.0.0.1:3]
	StatusCodes:  map[200:2 503:1]
/report
	TotalHits:    1
	TotalBytes:   123
	Methods:      map[GET:1]
	Users:        map[james:1]
	RemoteHosts:  map[127.0.0.1:1]
	StatusCodes:  map[200:1]
ALERT: High traffic generated an alert; nbrHits: 4, triggeredAt: 07-08-2019 19:52:48
ALERT: Alert recovered at 07-08-2019 19:53:49
WARN: "interrupt" detected; closing file and shutting down in 3 seconds
```

#### Third Party Libraries & Tools
* [hpcloud/tail](https://github.com/hpcloud/tail) library was used for tailing the log file
* [mingrammer/flog](https://github.com/mingrammer/flog) tool was used for generating fake logs
	* `flog -l > access.log`

#### Known Issues/Improvements
* If the file is truncated in any way, the file will be reopened and the entire contents read again
* If the log line isn't in the correct format a panic can occur
* The shutdown logic should close the file first and then wait for a period of time for the remaining logs in the channel to be processed
* The `monitor_test.go:Generate alerts without recovering` test has a race condition or bug where the alert recovers unexpectedly every once in awhile
