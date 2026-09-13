//go:build !darwin

package detector

import "errors"

var errDarwinOnly = errors.New("macOS detection is only supported on Darwin")

type MacOSDetector struct{}

func NewMacOSDetector() *MacOSDetector {
	return &MacOSDetector{}
}

func (d *MacOSDetector) ListDevices() ([]DeviceInfo, error) {
	return nil, errDarwinOnly
}

func (d *MacOSDetector) Detect(devicePath string) (DeviceStatus, error) {
	return DeviceStatus{}, errDarwinOnly
}

func (d *MacOSDetector) Close() error { return nil }
