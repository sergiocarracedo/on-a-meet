package config

import (
	"runtime"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

type Config struct {
	Camera          string `mapstructure:"camera"`
	Interval        string `mapstructure:"interval"`
	OnCommand       string `mapstructure:"on-command"`
	OffCommand      string `mapstructure:"off-command"`
	DetectMethod    string `mapstructure:"detect-method"`
	Debounce        int    `mapstructure:"debounce"`
	Timeout         string `mapstructure:"timeout"`
	Verbose         bool   `mapstructure:"verbose"`
	EnvironmentFile string `mapstructure:"environment-file"`
}

func Defaults() Config {
	return Config{
		Interval:     "1s",
		DetectMethod: detector.DefaultMethod(runtime.GOOS),
		Debounce:     3,
		Timeout:      "30s",
	}
}
