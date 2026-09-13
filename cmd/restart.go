package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the on-a-meet service to reload config",
	Long:  `Stops and starts the systemd (Linux) or launchd (macOS) service unit to pick up config changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requirePrivileges(runtime.GOOS, os.Geteuid(), "service restart"); err != nil {
			return err
		}

		if err := applyEnvironmentFile(runtime.GOOS, viper.GetString("environment-file")); err != nil {
			output.Warning.Printfln("Failed to patch environment file path: %v", err)
		}

		svc, err := newServiceHandle()
		if err != nil {
			return err
		}

		output.Info.Println("Restarting service...")
		if err := svc.Stop(); err != nil {
			output.Warning.Printfln("Service stop failed (may not be running): %v", err)
		} else {
			output.Success.Println("Service stopped")
		}

		if err := svc.Start(); err != nil {
			return fmt.Errorf("service start failed: %w (try 'sudo on-a-meet service install' first)", err)
		}
		output.Success.Println("Service restarted")

		return nil
	},
}

func init() {
	serviceCmd.AddCommand(restartCmd)
}
