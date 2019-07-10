package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ashnelson/httpLogMonitor/log"
)

// Config contains all of the application's configuration settings
type Config struct {
	LogFile             string      `json:"inputLogFile"`
	ShutdownGracePeriod int         `json:"shutdownWaitSeconds"`
	StatsCfg            StatsConfig `json:"statsCfg"`
}

// StatsConfig contains all of the configuration settings for the stats and
// alerts monitors
type StatsConfig struct {
	StatsIntervalSeconds int `json:"statsIntervalSeconds"`
	AlertIntervalSeconds int `json:"alertIntervalSeconds"`
	AlertThreshold       int `json:"alertThreshold"`
}

// New reads the specified config file and returns it. If, for some reason, the
// config file cannot be read, a default config will be returned.
func New(cfgFile string) Config {
	var cfg Config

	cfgFileBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return logErrAndGetDefaults(err)
	}

	if err := json.Unmarshal(cfgFileBytes, &cfg); err != nil {
		return logErrAndGetDefaults(err)
	}

	return cfg
}

func logErrAndGetDefaults(err error) Config {
	cfg := Config{
		LogFile:             "/tmp/access.log",
		ShutdownGracePeriod: 3,
		StatsCfg: StatsConfig{
			StatsIntervalSeconds: 10,
			AlertIntervalSeconds: 120,
			AlertThreshold:       100,
		},
	}
	log.LogError("Failed to load configs from file; %s\n\tUsing default values:\n\t\t%+v", err, cfg)
	return cfg
}
