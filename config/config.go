package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ashnelson/httpLogMonitor/log"
)

type Config struct {
	LogFile             string      `json:"inputLogFile"`
	ShutdownGracePeriod int         `json:"shutdownWaitSeconds"`
	StatsCfg            StatsConfig `json:"statsCfg"`
}

type StatsConfig struct {
	StatsIntervalSeconds int `json:"statsIntervalSeconds"`
	AlertIntervalSeconds int `json:"alertIntervalSeconds"`
	AlertThreshold       int `json:"alertThreshold"`
}

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
