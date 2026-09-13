package cmd

import (
	"path/filepath"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

const appName = "on-a-meet"

// These helpers all take goos and home explicitly rather than reading
// runtime.GOOS or the environment, so every platform's behaviour is testable
// from every platform. Only call sites pass runtime.GOOS.

// userConfigDir is the per-user config location, the same on all platforms.
func userConfigDir(home string) string {
	return filepath.Join(home, ".config", appName)
}

// systemConfigDir is the machine-wide config location.
func systemConfigDir(goos string) string {
	if goos == "darwin" {
		return filepath.Join("/Library/Application Support", appName)
	}
	return filepath.Join("/etc", appName)
}

// onboardConfigDir is where the onboarding wizard writes.
//
// On Linux the service runs as root from /etc. On macOS it is a per-user
// LaunchAgent, so the config belongs in the user's home: it is already the
// first entry in the search path, which means the agent and a hand-run
// `on-a-meet detect` read the exact same file, and the user can edit it
// without sudo.
func onboardConfigDir(goos, home string) string {
	if goos == "darwin" {
		return userConfigDir(home)
	}
	return systemConfigDir(goos)
}

func onboardConfigPath(goos, home string) string {
	return filepath.Join(onboardConfigDir(goos, home), "config.yaml")
}

// serviceConfigPath is the --config value baked into the service definition.
// On Linux it must stay byte-identical to the historical hardcoded string so
// already-installed units keep resolving.
func serviceConfigPath(goos, home string) string {
	if goos == "darwin" {
		return onboardConfigPath(goos, home)
	}
	return "/etc/on-a-meet/config.yaml"
}

// configSearchPaths lists config locations in precedence order.
func configSearchPaths(goos, home string) []string {
	paths := []string{userConfigDir(home)}
	if goos == "darwin" {
		// Keep /etc in the list so any config hand-placed there before
		// macOS had its own location still resolves.
		paths = append(paths, systemConfigDir(goos), "/etc/"+appName)
	} else {
		paths = append(paths, systemConfigDir(goos))
	}
	return append(paths, ".")
}

// noDevicesHelp returns the troubleshooting text shown when no cameras are
// found. The Linux advice (the video group) is meaningless on macOS, which has
// no /dev/video* nodes and needs no elevated privileges at all.
func noDevicesHelp(goos string) []string {
	if goos == "darwin" {
		return []string{
			"Make sure your camera is connected and not disabled.",
			"Tip: run `system_profiler SPCameraDataType` to see what macOS reports.",
		}
	}
	return []string{
		"Make sure your camera is connected and you have the right permissions.",
		"Tip: add your user to the 'video' group: sudo usermod -a -G video $USER",
		"Then log out and back in, or run: newgrp video",
	}
}

func printNoDevicesHelp(goos string) {
	for _, line := range noDevicesHelp(goos) {
		output.Info.Println(line)
	}
}
