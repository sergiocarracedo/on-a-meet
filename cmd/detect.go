package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	detectCamera   string
	detectInterval string
	detectOnCmd    string
	detectOffCmd   string
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect camera on/off state and execute commands",
	Long: `Continuously monitors camera devices and fires
user-defined commands when camera state changes.

Uses V4L2 by default to check /dev/video* device status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("detect: not yet implemented (Phase 2)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

	detectCmd.Flags().StringVarP(&detectCamera, "camera", "", "", "target specific camera (e.g., /dev/video0)")
	detectCmd.Flags().StringVarP(&detectInterval, "interval", "i", "1s", "polling interval")
	detectCmd.Flags().StringVarP(&detectOnCmd, "on", "", "", "command to run when camera turns on")
	detectCmd.Flags().StringVarP(&detectOffCmd, "off", "", "", "command to run when camera turns off")

	viper.BindPFlag("camera", detectCmd.Flags().Lookup("camera"))
	viper.BindPFlag("interval", detectCmd.Flags().Lookup("interval"))
	viper.BindPFlag("on-command", detectCmd.Flags().Lookup("on"))
	viper.BindPFlag("off-command", detectCmd.Flags().Lookup("off"))
}
