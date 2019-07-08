package main

import (
	"os"
	"os/signal"
	"syscall"

	"./config"
	"./log"
	"./stats"

	"github.com/hpcloud/tail"
)

func main() {
	stats.InitializeRecorder(config.New())

	// Open log file and tail it
	tailCfg := tail.Config{
		ReOpen: true,
		Follow: true,
	}
	tailFile, err := tail.TailFile("flog.log", tailCfg)
	if err != nil {
		log.PrintError("Failed to open file; Details: %s", err)
		return
	}

	// Make sure to release file handlers if application is terminated
	gracefulShutdown(func() {
		tailFile.Cleanup()
		close(stats.LogCh)
	})

	// Read and record each log line as it gets written to the log file
	for ln := range tailFile.Lines {
		//stats.RecordLogLine(ln.Text)
		stats.LogCh <- ln.Text
	}
}

// gracefulShutdown listens for OS interrupt signals and if a signal is detected
// the shutdown function is ran and the application closed
func gracefulShutdown(fn func()) {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT)

	go func() {
		interrupt := <-interruptCh

		log.PrintWarning("%q detected; closing file and shutting down", interrupt)
		fn()
		os.Exit(0)
	}()
}
