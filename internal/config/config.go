package config

type Config struct {
	Camera       string `mapstructure:"camera"`
	Interval     string `mapstructure:"interval"`
	OnCommand    string `mapstructure:"on-command"`
	OffCommand   string `mapstructure:"off-command"`
	DetectMethod string `mapstructure:"detect-method"`
	Debounce     int    `mapstructure:"debounce"`
	Timeout      string `mapstructure:"timeout"`
	Verbose      bool   `mapstructure:"verbose"`
}

func Defaults() Config {
	return Config{
		Interval:     "1s",
		DetectMethod: "v4l2",
		Debounce:     3,
		Timeout:      "30s",
	}
}
