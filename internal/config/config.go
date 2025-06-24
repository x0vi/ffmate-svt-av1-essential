package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

type ConfigDefinition struct {
	AppName    string `mapstructure:"appName"`
	AppVersion string `mapstructure:"appVersion"`

	FFMpeg string `mapstructure:"ffmpeg"`

	Port               uint   `mapstructure:"port"`
	Tray               bool   `mapstructure:"tray"`
	Database           string `mapstructure:"database"`
	Debug              string `mapstructure:"debug"`
	Loglevel           string `mapstructure:"loglevel"`
	MaxConcurrentTasks uint   `mapstructure:"maxConcurrentTasks"`
	SendTelemetry      bool   `mapstructure:"sendTelemetry"`
	NoUI               bool   `mapstructure:"noUI"`

	Mutex sync.RWMutex
}

var config ConfigDefinition

func Init() {
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("failed to unmarshal config: %s\n", err)
		os.Exit(1)
	}
	config.Mutex = sync.RWMutex{}

	if config.Debug == "" {
		config.Debug = os.Getenv("DEBUGO")
	}
}

func Config() *ConfigDefinition {
	return &config
}
