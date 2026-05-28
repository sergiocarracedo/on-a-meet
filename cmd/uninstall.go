package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove on-a-meet system service",
	Long:  `Stops and removes the systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("uninstall: not yet implemented (Phase 4)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
