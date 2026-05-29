package cmd

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the on-a-meet service",
	Long:  `Starts the systemd (Linux) or launchd (macOS) service unit if it is installed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet service start")
		}

		svc, err := service.New(&noopProgram{}, serviceConfig(""))
		if err != nil {
			return fmt.Errorf("service init failed: %w", err)
		}

		output.Info.Println("Starting service...")
		if err := svc.Start(); err != nil {
			return fmt.Errorf("service start failed: %w (try 'sudo on-a-meet service install' first)", err)
		}
		output.Success.Println("Service started")

		return nil
	},
}

func init() {
	serviceCmd.AddCommand(startCmd)
}
