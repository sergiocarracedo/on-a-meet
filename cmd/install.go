package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

type noopProgram struct{}

func (p *noopProgram) Start(s service.Service) error { return nil }
func (p *noopProgram) Stop(s service.Service) error  { return nil }

// serviceOptions returns the platform-specific knobs handed to
// kardianos/service.
//
// Linux gets nil, which renders exactly the systemd unit this tool has always
// produced. macOS needs three things set explicitly:
//
//   - UserService, so the plist goes to ~/Library/LaunchAgents rather than
//     /Library/LaunchDaemons. A LaunchDaemon has no GUI session, so the whole
//     point of this tool — running the user's on/off command, typically a
//     notification — would silently do nothing there.
//   - RunAtLoad, which unlike systemd's Install() defaults to false, so
//     without it the agent would not come back after logout or reboot.
//   - LogDirectory, because a user service otherwise logs to ~/on-a-meet.out.log.
func serviceOptions(goos, home string) service.KeyValue {
	if goos != "darwin" {
		return nil
	}
	return service.KeyValue{
		"UserService":  true,
		"RunAtLoad":    true,
		"KeepAlive":    true,
		"LogDirectory": filepath.Join(home, "Library", "Logs"),
	}
}

func serviceConfig(goos, home, user string) *service.Config {
	cfg := &service.Config{
		Name:        appName,
		DisplayName: appName,
		Description: "Camera state monitoring service",
		Arguments: []string{
			"detect",
			"--config", serviceConfigPath(goos, home),
		},
		WorkingDirectory: "/",
		UserName:         user,
		Option:           serviceOptions(goos, home),
	}
	if goos == "darwin" {
		// A LaunchAgent already runs as the owning user, and "/" is a poor
		// working directory for the commands it launches.
		cfg.UserName = ""
		cfg.WorkingDirectory = home
	}
	return cfg
}

// rewriteSystemdEnvironmentFile points a rendered systemd unit at envFile.
// Pure and testable; changed is false when the unit has no such directive.
func rewriteSystemdEnvironmentFile(unit, envFile string) (string, bool) {
	oldLine := "EnvironmentFile=-/etc/sysconfig/" + appName
	if !strings.Contains(unit, oldLine) {
		return unit, false
	}
	return strings.ReplaceAll(unit, oldLine, "EnvironmentFile=-"+envFile), true
}

// applyEnvironmentFile wires the environment-file config key into the
// platform's service definition.
//
// On macOS this is deliberately a no-op rather than a gap. launchd has no
// EnvironmentFile equivalent, only EnvironmentVariables baked in at install
// time — but on-a-meet already implements this key itself: the executor
// re-reads the file on every command execution and expands it into the
// command. So the feature works on macOS with no launchd involvement, and
// edits take effect without reinstalling the service.
//
// Baking the values into the plist was considered and rejected: it would
// freeze them at install time, write secrets in plaintext under
// ~/Library/LaunchAgents, and diverge from the live-reload behaviour on Linux.
func applyEnvironmentFile(goos, envFile string) error {
	if envFile == "" || goos == "darwin" {
		return nil
	}
	unitPath := "/etc/systemd/system/" + appName + ".service"
	data, err := os.ReadFile(unitPath)
	if err != nil {
		return fmt.Errorf("reading unit file: %w", err)
	}
	content, changed := rewriteSystemdEnvironmentFile(string(data), envFile)
	if !changed {
		return nil
	}
	if err := os.WriteFile(unitPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("writing unit file: %w", err)
	}
	return exec.Command("systemctl", "daemon-reload").Run()
}

// newServiceHandle builds the service handle used by every service
// subcommand, so none of them can drift on how goos and home are resolved.
func newServiceHandle() (service.Service, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolving home directory: %w", err)
	}
	svc, err := service.New(&noopProgram{}, serviceConfig(runtime.GOOS, home, os.Getenv("SUDO_USER")))
	if err != nil {
		return nil, fmt.Errorf("service init failed: %w", err)
	}
	return svc, nil
}

func installService() error {
	goos := runtime.GOOS

	svc, err := newServiceHandle()
	if err != nil {
		return err
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

	if err := applyEnvironmentFile(goos, viper.GetString("environment-file")); err != nil {
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
		if err := requirePrivileges(runtime.GOOS, os.Geteuid(), "service install"); err != nil {
			return err
		}

		return installService()
	},
}

func init() {
	serviceCmd.AddCommand(installCmd)
}
