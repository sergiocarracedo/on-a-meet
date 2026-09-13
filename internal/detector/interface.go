package detector

import (
	"errors"
	"time"
)

// ErrDeviceNotObservable reports that a device was enumerated but the active
// detection method has no way to observe its on/off state. It is returned
// instead of a plain "off" so the failure is visible rather than silent.
var ErrDeviceNotObservable = errors.New("device state cannot be observed with this detection method")

type DeviceStatus struct {
	On        bool
	CheckedAt time.Time
}

type DeviceInfo struct {
	// Path identifies the device to Detect. On Linux it is the device node
	// (/dev/video0); on macOS there are no device nodes, so it is the display
	// name reported by system_profiler ("Logitech StreamCam").
	Path string
	// ID is a short, stable identifier for the device, exposed to user
	// commands as {{.CameraID}}. Linux uses the node basename ("video0"),
	// macOS the CoreMediaIO unique-id.
	ID     string
	Driver string
	Card   string
	Bus    string
}

type Detector interface {
	ListDevices() ([]DeviceInfo, error)
	Detect(devicePath string) (DeviceStatus, error)
}
