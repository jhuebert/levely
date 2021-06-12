package service

import (
	"github.com/jhuebert/levely/config"
	"github.com/spf13/viper"
)

type Config struct {
	DisplayLevelTolerance float64 `json:"displayLevelTolerance"`
	DisplayUpdateRate     float64 `json:"displayUpdateRate"`
	DisplaySseEnabled     bool    `json:"displaySseEnabled"`
}

func (s *Service) GetConfig() Config {
	return Config{
		DisplayLevelTolerance: viper.GetFloat64(config.DisplayLevelTolerance),
		DisplayUpdateRate:     viper.GetFloat64(config.DisplayUpdateRate),
		DisplaySseEnabled:     viper.GetBool(config.DisplaySseEnabled),
	}
}
