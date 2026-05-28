//go:build !linux

package detector

import "errors"

type V4L2Detector struct{}

func NewV4L2Detector() *V4L2Detector {
	return &V4L2Detector{}
}

func (d *V4L2Detector) ListDevices() ([]DeviceInfo, error) {
	return nil, errors.New("V4L2 detection is only supported on Linux")
}

func (d *V4L2Detector) Detect(devicePath string) (DeviceStatus, error) {
	return DeviceStatus{}, errors.New("V4L2 detection is only supported on Linux")
}
