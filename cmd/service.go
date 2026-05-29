package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Manage the on-a-meet system service",
	Long:    `Install, uninstall, start, stop, and restart the systemd (Linux) or launchd (macOS) service.`,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
