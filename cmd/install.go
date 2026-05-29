package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

type noopProgram struct{}

func (p *noopProgram) Start(s service.Service) error { return nil }
func (p *noopProgram) Stop(s service.Service) error  { return nil }

func serviceConfig(user string) *service.Config {
	return &service.Config{
		Name:        "on-a-meet",
		DisplayName: "on-a-meet",
		Description: "Camera state monitoring service",
		Arguments: []string{
			"detect",
			"--config", "/etc/on-a-meet/config.yaml",
		},
		WorkingDirectory: "/",
		UserName:         user,
	}
}

func patchUnitEnvironmentFile(envFile string) error {
	if envFile == "" {
		return nil
	}
	unitPath := "/etc/systemd/system/on-a-meet.service"
	data, err := os.ReadFile(unitPath)
	if err != nil {
		return fmt.Errorf("reading unit file: %w", err)
	}
	oldLine := "EnvironmentFile=-/etc/sysconfig/on-a-meet"
	newLine := "EnvironmentFile=-" + envFile
	if !strings.Contains(string(data), oldLine) {
		return nil
	}
	content := strings.ReplaceAll(string(data), oldLine, newLine)
	if err := os.WriteFile(unitPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing unit file: %w", err)
	}
	return exec.Command("systemctl", "daemon-reload").Run()
}

func installService() error {
	originalUser := os.Getenv("SUDO_USER")
	svc, err := service.New(&noopProgram{}, serviceConfig(originalUser))
	if err != nil {
		return fmt.Errorf("service init failed: %w", err)
	}

	// Stop existing service if running
	_ = svc.Stop()

	output.Info.Println("Installing service...")
	if err := svc.Install(); err != nil {
		// If unit already exists, remove and retry
		_ = svc.Uninstall()
		if err := svc.Install(); err != nil {
			return fmt.Errorf("service install failed: %w", err)
		}
	}
	output.Success.Println("Service unit created")

	if err := patchUnitEnvironmentFile(viper.GetString("environment-file")); err != nil {
		output.Warning.Printfln("Failed to patch environment file path: %v", err)
	}

	output.Info.Println("Starting service...")
	if err := svc.Start(); err != nil {
		return fmt.Errorf("service start failed: %w", err)
	}
	output.Success.Println("Service started")

	return nil
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install on-a-meet as a system service",
	Long:  `Creates and enables a systemd (Linux) or launchd (macOS) service unit.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet service install")
		}

		return installService()
	},
}

func init() {
	serviceCmd.AddCommand(installCmd)
}
