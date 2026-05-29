//go:build !linux

package detector

import "errors"

type LsofDetector struct{}

func NewLsofDetector() *LsofDetector {
	return &LsofDetector{}
}

func (d *LsofDetector) ListDevices() ([]DeviceInfo, error) {
	return nil, errors.New("lsof detection is only supported on Linux")
}

func (d *LsofDetector) Detect(devicePath string) (DeviceStatus, error) {
	return DeviceStatus{}, errors.New("lsof detection is only supported on Linux")
}
