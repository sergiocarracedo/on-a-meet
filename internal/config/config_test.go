package config

import "testing"

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	if cfg.Interval != "1s" {
		t.Errorf("Default interval = %q, want %q", cfg.Interval, "1s")
	}
	if cfg.DetectMethod != "v4l2" {
		t.Errorf("Default detect-method = %q, want %q", cfg.DetectMethod, "v4l2")
	}
	if cfg.Debounce != 3 {
		t.Errorf("Default debounce = %d, want %d", cfg.Debounce, 3)
	}
}
