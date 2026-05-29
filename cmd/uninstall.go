package cmd

import (
	"fmt"
	"os"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove on-a-meet system service",
	Long:  `Stops and removes the systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet uninstall")
		}

		svc, err := service.New(&noopProgram{}, serviceConfig())
		if err != nil {
			return fmt.Errorf("service init failed: %w", err)
		}

		output.Info.Println("Stopping service...")
		if err := svc.Stop(); err != nil {
			output.Warning.Printfln("Service stop failed (may not be running): %v", err)
		} else {
			output.Success.Println("Service stopped")
		}

		output.Info.Println("Removing service...")
		if err := svc.Uninstall(); err != nil {
			return fmt.Errorf("service uninstall failed: %w", err)
		}
		output.Success.Println("Service unit removed")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
