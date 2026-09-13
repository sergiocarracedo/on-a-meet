//go:build !darwin

package detector

import "testing"

func TestMacOSStubReturnsError(t *testing.T) {
	d := MacOSDetector{}

	devices, err := d.ListDevices()
	if err == nil {
		t.Error("expected error on non-Darwin, got nil")
	}
	if devices != nil {
		t.Error("expected nil devices on non-Darwin")
	}

	status, err := d.Detect("Logitech StreamCam")
	if err == nil {
		t.Error("expected error on non-Darwin, got nil")
	}
	if status.On {
		t.Error("expected On=false on non-Darwin")
	}

	if err := d.Close(); err != nil {
		t.Errorf("Close on non-Darwin should be a no-op, got %v", err)
	}
}
