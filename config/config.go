package config

import (
	"encoding/json"
	"io/ioutil"

	"../log"
)

type StatsCfg struct {
	StatsIntervalSeconds   int
	AlertIntervalSeconds   int
	AlertThreshold         int
	RecoverIntervalSeconds int
}

func New() StatsCfg {
	var cfg StatsCfg

	cfgFileBytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		return logErrAndGetDefaults(err)
	}

	if err := json.Unmarshal(cfgFileBytes, &cfg); err != nil {
		return logErrAndGetDefaults(err)
	}

	return cfg
}

func logErrAndGetDefaults(err error) StatsCfg {
	cfg := StatsCfg{
		StatsIntervalSeconds:   10,
		AlertIntervalSeconds:   120,
		AlertThreshold:         100,
		RecoverIntervalSeconds: 120,
	}
	log.PrintError("Failed to load configs from file\n\tDetails: %s\nUsing default values:\n\t%+v", err, cfg)
	return cfg
}
