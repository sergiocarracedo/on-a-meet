package detector

import "time"

type DeviceStatus struct {
	On        bool
	CheckedAt time.Time
}

type DeviceInfo struct {
	Path   string
	Driver string
	Card   string
	Bus    string
}

type Detector interface {
	ListDevices() ([]DeviceInfo, error)
	Detect(devicePath string) (DeviceStatus, error)
}
