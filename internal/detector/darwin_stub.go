//go:build !darwin

package detector

import "errors"

type MacOSDetector struct{}

func NewMacOSDetector() *MacOSDetector {
	return &MacOSDetector{}
}

func (d *MacOSDetector) ListDevices() ([]DeviceInfo, error) {
	return nil, errors.New("macOS detection is only supported on Darwin")
}

func (d *MacOSDetector) Detect(devicePath string) (DeviceStatus, error) {
	return DeviceStatus{}, errors.New("macOS detection is only supported on Darwin")
}
