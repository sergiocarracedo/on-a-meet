package cmd

import (
	"strings"
	"testing"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

func TestRequirePrivileges(t *testing.T) {
	const nonRoot = 501

	tests := []struct {
		name    string
		goos    string
		euid    int
		wantErr bool
	}{
		{"linux as root is allowed", "linux", 0, false},
		{"linux without root is refused", "linux", nonRoot, true},
		// A LaunchAgent's plist path comes from the *current* user's home, so
		// under sudo it would be installed for root and never run.
		{"darwin as root is refused", "darwin", 0, true},
		{"darwin as the user is allowed", "darwin", nonRoot, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requirePrivileges(tt.goos, tt.euid, "service install")
			if tt.wantErr && err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// The macOS message has to tell the user what to do differently, not just
// that something went wrong.
func TestRequirePrivilegesDarwinRootMessageIsActionable(t *testing.T) {
	err := requirePrivileges("darwin", 0, "onboard --apply x.json")
	if err == nil {
		t.Fatal("expected an error")
	}
	msg := err.Error()
	for _, want := range []string{"sudo", "without sudo"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message should mention %q, got: %s", want, msg)
		}
	}
}

func TestNeedsSudoReexec(t *testing.T) {
	if !needsSudoReexec("linux") {
		t.Error("linux writes to /etc and installs a system unit; it needs the sudo re-exec")
	}
	if needsSudoReexec("darwin") {
		t.Error("darwin writes only inside the user's home; it must not re-exec under sudo")
	}
}

func TestDetectionMethodOptionsMatchMethods(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		methods := detector.MethodsFor(goos)
		opts := detectionMethodOptions(methods)
		if len(opts) != len(methods) {
			t.Fatalf("%s: %d options for %d methods", goos, len(opts), len(methods))
		}
		for i, m := range methods {
			if opts[i].Value != m.Name {
				t.Errorf("%s: option %d value = %q, want %q", goos, i, opts[i].Value, m.Name)
			}
		}
	}
}

func TestDetectionMethodHelpMentionsEveryMethod(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		help := detectionMethodHelp(detector.MethodsFor(goos))
		for _, m := range detector.MethodsFor(goos) {
			if !strings.Contains(help, m.Label) {
				t.Errorf("%s: help text omits %q", goos, m.Label)
			}
		}
	}
}
