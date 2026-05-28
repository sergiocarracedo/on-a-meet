package config

type Config struct {
	Camera       string `mapstructure:"camera"`
	Interval     string `mapstructure:"interval"`
	OnCommand    string `mapstructure:"on-command"`
	OffCommand   string `mapstructure:"off-command"`
	DetectMethod string `mapstructure:"detect-method"`
	Debounce     int    `mapstructure:"debounce"`
}

func Defaults() Config {
	return Config{
		Interval:     "1s",
		DetectMethod: "v4l2",
		Debounce:     3,
	}
}
