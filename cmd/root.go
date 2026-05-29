package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

		viper.AddConfigPath(home + "/.config/on-a-meet")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ON_A_MEET")

	viper.SetDefault("detect-method", "v4l2")
	viper.SetDefault("interval", "1s")
	viper.SetDefault("debounce", 3)
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("camera", "")
	viper.SetDefault("on-command", "")
	viper.SetDefault("off-command", "")

	output.Init(cfgSilent, cfgVerbose)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
