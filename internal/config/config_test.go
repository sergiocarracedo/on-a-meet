package config

import (
	"runtime"
	"testing"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	if cfg.Interval != "1s" {
		t.Errorf("Default interval = %q, want %q", cfg.Interval, "1s")
	}
	// The default must match the host OS: v4l2 is meaningless on macOS,
	// which has no /dev/video* nodes.
	want := detector.DefaultMethod(runtime.GOOS)
	if cfg.DetectMethod != want {
		t.Errorf("Default detect-method = %q, want %q", cfg.DetectMethod, want)
	}
	if !detector.IsValidMethod(runtime.GOOS, cfg.DetectMethod) {
		t.Errorf("default method %q is not usable on %s", cfg.DetectMethod, runtime.GOOS)
	}
	if cfg.Debounce != 3 {
		t.Errorf("Default debounce = %d, want %d", cfg.Debounce, 3)
	}
}
