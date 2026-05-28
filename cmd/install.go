package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install on-a-meet as a system service",
	Long:  `Creates and enables a systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("install: not yet implemented (Phase 4)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
