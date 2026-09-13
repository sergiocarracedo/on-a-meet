package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the on-a-meet service",
	Long:  `Stops the systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requirePrivileges(runtime.GOOS, os.Geteuid(), "service stop"); err != nil {
			return err
		}

		svc, err := newServiceHandle()
		if err != nil {
			return err
		}

		output.Info.Println("Stopping service...")
		if err := svc.Stop(); err != nil {
			// On macOS `launchctl unload` errors when the agent simply is
			// not loaded, which is a routine state, not a failure.
			if runtime.GOOS == "darwin" {
				output.Warning.Printfln("Service stop failed (may not be running): %v", err)
				return nil
			}
			return fmt.Errorf("service stop failed: %w", err)
		}
		output.Success.Println("Service stopped")

		return nil
	},
}

func init() {
	serviceCmd.AddCommand(stopCmd)
}
