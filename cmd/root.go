package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var (
	version    = "dev"
	cfgFile    string
	cfgSilent  bool
	cfgVerbose bool

	rootCmd = &cobra.Command{
		Use:     "on-a-meet",
		Version: version,
		Short:   "Monitor camera on/off state and trigger commands",
		Long: `on-a-meet detects when your camera turns on or off
and executes user-defined commands on state transitions.

It polls /dev/video* devices at a configurable interval and
fires --on and --off commands with template variable substitution.`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.config/on-a-meet/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&cfgSilent, "silent", "s", false, "suppress all output")
	rootCmd.PersistentFlags().BoolVarP(&cfgVerbose, "verbose", "V", false, "enable debug output")

	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		for _, p := range configSearchPaths(runtime.GOOS, home) {
			viper.AddConfigPath(p)
		}
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ON_A_MEET")
	// Config keys are hyphenated (detect-method, environment-file), but a
	// hyphen cannot appear in an environment variable name. Without this
	// replacer none of those keys can be overridden from the environment at
	// all: ON_A_MEET_DETECT_METHOD, ON_A_MEET_ENVIRONMENT_FILE and friends
	// are simply never matched.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.SetDefault("detect-method", detector.DefaultMethod(runtime.GOOS))
	viper.SetDefault("interval", "1s")
	viper.SetDefault("debounce", 3)
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("camera", "")
	viper.SetDefault("on-command", "")
	viper.SetDefault("off-command", "")
	viper.SetDefault("silent", false)
	viper.SetDefault("verbose", false)
	viper.SetDefault("environment-file", "")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	output.Init(viper.GetBool("silent"), viper.GetBool("verbose"))
}
