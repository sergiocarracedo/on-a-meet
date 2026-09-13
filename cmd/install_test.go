package cmd

import (
	"strings"
	"testing"
)

// nil options mean kardianos renders exactly the systemd unit this tool has
// always produced.
func TestServiceOptionsLinuxIsNil(t *testing.T) {
	if got := serviceOptions("linux", testHome); got != nil {
		t.Errorf("linux options must stay nil, got %v", got)
	}
}

func TestServiceOptionsDarwinIsUserAgent(t *testing.T) {
	opts := serviceOptions("darwin", testHome)
	if opts == nil {
		t.Fatal("darwin needs explicit options")
	}
	// Without UserService the plist would become a system LaunchDaemon with
	// no GUI session, where notification commands silently do nothing.
	if opts["UserService"] != true {
		t.Error("UserService must be true so this installs as a per-user LaunchAgent")
	}
	// kardianos defaults RunAtLoad to false on darwin, unlike systemd's
	// Install() which also enables the unit.
	if opts["RunAtLoad"] != true {
		t.Error("RunAtLoad must be true or the agent will not return after logout")
	}
	if opts["KeepAlive"] != true {
		t.Error("KeepAlive must be true")
	}
	if dir, _ := opts["LogDirectory"].(string); !strings.HasPrefix(dir, testHome) {
		t.Errorf("LogDirectory should be under the user's home, got %q", dir)
	}
}

func TestServiceConfigArguments(t *testing.T) {
	linux := serviceConfig("linux", testHome, "someone")
	want := []string{"detect", "--config", "/etc/on-a-meet/config.yaml"}
	if len(linux.Arguments) != len(want) {
		t.Fatalf("linux arguments = %v", linux.Arguments)
	}
	for i := range want {
		if linux.Arguments[i] != want[i] {
			t.Errorf("linux argument %d = %q, want %q", i, linux.Arguments[i], want[i])
		}
	}

	darwin := serviceConfig("darwin", testHome, "someone")
	if darwin.Arguments[2] != onboardConfigPath("darwin", testHome) {
		t.Errorf("darwin config argument = %q", darwin.Arguments[2])
	}
}

func TestServiceConfigLinuxKeepsUserName(t *testing.T) {
	if got := serviceConfig("linux", testHome, "someone"); got.UserName != "someone" {
		t.Errorf("UserName = %q, want %q", got.UserName, "someone")
	}
}

// A LaunchAgent already runs as its owning user; a UserName key there is at
// best ignored.
func TestServiceConfigDarwinDropsUserName(t *testing.T) {
	if got := serviceConfig("darwin", testHome, "someone"); got.UserName != "" {
		t.Errorf("UserName should be empty on darwin, got %q", got.UserName)
	}
}

func TestRewriteSystemdEnvironmentFile(t *testing.T) {
	unit := "[Service]\nEnvironmentFile=-/etc/sysconfig/on-a-meet\nExecStart=/usr/local/bin/on-a-meet\n"

	got, changed := rewriteSystemdEnvironmentFile(unit, "/etc/default/on-a-meet")
	if !changed {
		t.Fatal("expected the directive to be rewritten")
	}
	if !strings.Contains(got, "EnvironmentFile=-/etc/default/on-a-meet") {
		t.Errorf("rewritten unit missing the new path:\n%s", got)
	}
	if strings.Contains(got, "/etc/sysconfig/on-a-meet") {
		t.Errorf("old path should be gone:\n%s", got)
	}

	other := "[Service]\nExecStart=/usr/local/bin/on-a-meet\n"
	got, changed = rewriteSystemdEnvironmentFile(other, "/etc/default/on-a-meet")
	if changed {
		t.Error("a unit without the directive should report no change")
	}
	if got != other {
		t.Error("a unit without the directive must be returned unmodified")
	}
}

// The rewrite is a literal string replacement, so regex and expansion
// metacharacters in the path must survive untouched.
func TestRewriteSystemdEnvironmentFileHandlesMetacharacters(t *testing.T) {
	unit := "EnvironmentFile=-/etc/sysconfig/on-a-meet\n"
	weird := `/etc/default/on-a-meet$1.[a-z]*`
	got, changed := rewriteSystemdEnvironmentFile(unit, weird)
	if !changed {
		t.Fatal("expected a change")
	}
	if !strings.Contains(got, weird) {
		t.Errorf("path was mangled: %q", got)
	}
}

// launchd has no EnvironmentFile; on-a-meet implements the key itself at
// command-exec time, so there is nothing to patch. This must return before
// touching the filesystem, which is also what makes it safe to assert here.
func TestApplyEnvironmentFileDarwinIsNoop(t *testing.T) {
	if err := applyEnvironmentFile("darwin", "/definitely/not/a/real/path"); err != nil {
		t.Errorf("darwin should be a no-op, got %v", err)
	}
}

func TestApplyEnvironmentFileEmptyIsNoop(t *testing.T) {
	if err := applyEnvironmentFile("linux", ""); err != nil {
		t.Errorf("empty env file should be a no-op, got %v", err)
	}
}
