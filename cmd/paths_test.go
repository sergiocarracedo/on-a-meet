package cmd

import (
	"strings"
	"testing"
)

const testHome = "/Users/someone"

// Locks the historical path: already-installed systemd units point at this
// exact file, so it must not move.
func TestServiceConfigPathLinuxIsUnchanged(t *testing.T) {
	if got := serviceConfigPath("linux", testHome); got != "/etc/on-a-meet/config.yaml" {
		t.Errorf("got %q, want %q", got, "/etc/on-a-meet/config.yaml")
	}
}

func TestSystemConfigDirLinuxIsUnchanged(t *testing.T) {
	if got := systemConfigDir("linux"); got != "/etc/on-a-meet" {
		t.Errorf("got %q, want %q", got, "/etc/on-a-meet")
	}
}

func TestServiceConfigPathDarwinUsesGivenHome(t *testing.T) {
	got := serviceConfigPath("darwin", testHome)
	if !strings.HasPrefix(got, testHome) {
		t.Errorf("got %q, want a path under %q", got, testHome)
	}
	if strings.HasPrefix(got, "/etc") {
		t.Errorf("macOS must not use /etc, got %q", got)
	}
}

func TestOnboardConfigPath(t *testing.T) {
	if got := onboardConfigPath("linux", testHome); got != "/etc/on-a-meet/config.yaml" {
		t.Errorf("linux: got %q", got)
	}
	want := testHome + "/.config/on-a-meet/config.yaml"
	if got := onboardConfigPath("darwin", testHome); got != want {
		t.Errorf("darwin: got %q, want %q", got, want)
	}
}

// The macOS LaunchAgent and a hand-run `on-a-meet detect` must read the same
// file, so onboarding has to write into the search path's first entry.
func TestDarwinOnboardPathIsFirstInSearchPath(t *testing.T) {
	paths := configSearchPaths("darwin", testHome)
	if len(paths) == 0 {
		t.Fatal("no search paths")
	}
	if paths[0] != onboardConfigDir("darwin", testHome) {
		t.Errorf("onboard writes to %q but the search path starts at %q",
			onboardConfigDir("darwin", testHome), paths[0])
	}
}

func TestConfigSearchPathsLinuxUnchanged(t *testing.T) {
	want := []string{testHome + "/.config/on-a-meet", "/etc/on-a-meet", "."}
	got := configSearchPaths("linux", testHome)
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("path %d = %q, want %q", i, got[i], want[i])
		}
	}
}

// A config hand-placed in /etc on a Mac before macOS had its own location
// must keep resolving.
func TestConfigSearchPathsDarwinStillIncludesEtc(t *testing.T) {
	found := false
	for _, p := range configSearchPaths("darwin", testHome) {
		if p == "/etc/on-a-meet" {
			found = true
		}
	}
	if !found {
		t.Error("darwin search path should still include /etc/on-a-meet")
	}
}

func TestNoDevicesHelpIsOSAppropriate(t *testing.T) {
	linux := strings.Join(noDevicesHelp("linux"), "\n")
	if !strings.Contains(linux, "video") {
		t.Error("linux help should mention the video group")
	}
	darwin := strings.Join(noDevicesHelp("darwin"), "\n")
	// macOS has no /dev/video* and no video group; that advice would be noise.
	if strings.Contains(darwin, "usermod") || strings.Contains(darwin, "newgrp") {
		t.Errorf("macOS help must not give Linux group advice: %s", darwin)
	}
}
