//go:build !linux

package detector

import "testing"

func TestV4L2StubReturnsError(t *testing.T) {
	d := V4L2Detector{}

	devices, err := d.ListDevices()
	if err == nil {
		t.Error("expected error on non-Linux, got nil")
	}
	if devices != nil {
		t.Error("expected nil devices on non-Linux")
	}

	status, err := d.Detect("/dev/video0")
	if err == nil {
		t.Error("expected error on non-Linux, got nil")
	}
	if status.On {
		t.Error("expected On=false on non-Linux")
	}
}
