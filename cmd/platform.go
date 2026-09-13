package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

// requirePrivileges checks that the process is running with the privileges the
// platform's service model needs.
//
// Linux installs a system-wide systemd unit under /etc and needs root. macOS
// installs a per-user LaunchAgent, and must NOT run under sudo: the plist path
// is resolved from the *current* user's home, so under sudo it lands in
// /var/root/Library/LaunchAgents, `launchctl load` reports success, and the
// agent then never runs for the actual user. Detection itself needs no
// elevation on macOS — neither system_profiler nor `log` requires root.
func requirePrivileges(goos string, euid int, action string) error {
	if goos == "darwin" {
		if euid == 0 {
			return fmt.Errorf(
				"do not run `%s` with sudo on macOS — it installs a per-user LaunchAgent, "+
					"and running as root would install it for the root user instead, where it would never run. "+
					"Re-run as yourself, without sudo", action)
		}
		return nil
	}
	if euid != 0 {
		return fmt.Errorf("root privileges required — please re-run with sudo: sudo on-a-meet %s", action)
	}
	return nil
}

// needsSudoReexec reports whether onboarding has to hand off to a root
// re-exec to apply its configuration. macOS writes only to the user's own
// home, so it applies in-process.
func needsSudoReexec(goos string) bool {
	return goos != "darwin"
}

// detectionMethodOptions renders the methods available on goos as huh options.
func detectionMethodOptions(methods []detector.MethodInfo) []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(methods))
	for _, m := range methods {
		opts = append(opts, huh.NewOption(m.Label, m.Name))
	}
	return opts
}

// detectionMethodHelp renders the per-method descriptions shown beside the
// picker.
func detectionMethodHelp(methods []detector.MethodInfo) string {
	s := ""
	for i, m := range methods {
		if i > 0 {
			s += "\n"
		}
		s += m.Label + ": " + m.Description
	}
	return s
}
