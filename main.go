package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashnelson/httpLogMonitor/config"
	"github.com/ashnelson/httpLogMonitor/log"
	"github.com/ashnelson/httpLogMonitor/stats"

	"github.com/hpcloud/tail"
)

func main() {
	// Parse flags to get config file
	var cfgFile string
	flag.StringVar(&cfgFile, "cfg", "config.json", "the config file to use for application startup")
	flag.Parse()

	// Get config values
	cfg := config.New(cfgFile)

	// Start the stats recorder
	stats.InitializeRecorder(cfg.StatsCfg)

	// Open log file and tail it
	tailCfg := tail.Config{
		ReOpen: true,
		Follow: true,
	}
	tailFile, err := tail.TailFile(cfg.LogFile, tailCfg)
	if err != nil {
		log.LogError("Failed to open file; Details: %s", err)
		return
	}

	// Release file handle and shutdown stats monitoring if application is terminated
	gracefulShutdown(cfg.ShutdownGracePeriod, func() {
		tailFile.Cleanup()
		stats.Shutdown()
	})

	// Read and record each log line as it gets written to the log file
	for ln := range tailFile.Lines {
		stats.RecordLog(ln.Text)
	}
}

// gracefulShutdown listens for OS interrupt signals and if a signal is detected
// the shutdown function is ran and the application closed with a zero exit code
func gracefulShutdown(shutdownGracePeriod int, fn func()) {
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT)

	go func() {
		interrupt := <-interruptCh

		log.LogWarning("%q detected; closing file and shutting down in %d seconds", interrupt, shutdownGracePeriod)
		time.Sleep(time.Duration(shutdownGracePeriod) * time.Second)
		fn()
		os.Exit(0)
	}()
}
