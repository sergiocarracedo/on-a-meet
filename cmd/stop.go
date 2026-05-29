package cmd

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the on-a-meet service",
	Long:  `Stops the systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet service stop")
		}

		svc, err := service.New(&noopProgram{}, serviceConfig(""))
		if err != nil {
			return fmt.Errorf("service init failed: %w", err)
		}

		output.Info.Println("Stopping service...")
		if err := svc.Stop(); err != nil {
			return fmt.Errorf("service stop failed: %w", err)
		}
		output.Success.Println("Service stopped")

		return nil
	},
}

func init() {
	serviceCmd.AddCommand(stopCmd)
}
