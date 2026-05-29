package detector

import (
	"testing"
)

func TestNewV4L2(t *testing.T) {
	d, err := New("v4l2")
	if err != nil {
		t.Fatalf("New('v4l2') failed: %v", err)
	}
	if d == nil {
		t.Fatal("New('v4l2') returned nil")
	}
}

func TestNewLsof(t *testing.T) {
	d, err := New("lsof")
	if err != nil {
		t.Fatalf("New('lsof') failed: %v", err)
	}
	if d == nil {
		t.Fatal("New('lsof') returned nil")
	}
}

func TestNewUnknown(t *testing.T) {
	_, err := New("unknown")
	if err == nil {
		t.Fatal("New('unknown') should return error")
	}
}

func TestNewV4L2ImplementsDetector(t *testing.T) {
	d, err := New("v4l2")
	if err != nil {
		t.Fatal(err)
	}
	var _ Detector = d
}

func TestNewLsofImplementsDetector(t *testing.T) {
	d, err := New("lsof")
	if err != nil {
		t.Fatal(err)
	}
	var _ Detector = d
}

func TestNewDarwin(t *testing.T) {
	d, err := New("darwin")
	if err != nil {
		t.Fatalf("New('darwin') failed: %v", err)
	}
	if d == nil {
		t.Fatal("New('darwin') returned nil")
	}
}

func TestNewDarwinImplementsDetector(t *testing.T) {
	d, err := New("darwin")
	if err != nil {
		t.Fatal(err)
	}
	var _ Detector = d
}
