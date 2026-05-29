package cmd

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the on-a-meet service to reload config",
	Long:  `Stops and starts the systemd (Linux) or launchd (macOS) service unit to pick up config changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet restart")
		}

		svc, err := service.New(&noopProgram{}, serviceConfig(""))
		if err != nil {
			return fmt.Errorf("service init failed: %w", err)
		}

		output.Info.Println("Restarting service...")
		if err := svc.Stop(); err != nil {
			output.Warning.Printfln("Service stop failed (may not be running): %v", err)
		} else {
			output.Success.Println("Service stopped")
		}

		if err := svc.Start(); err != nil {
			return fmt.Errorf("service start failed: %w", err)
		}
		output.Success.Println("Service restarted")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
