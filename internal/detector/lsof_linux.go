//go:build linux

package detector

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

type LsofDetector struct{}

func NewLsofDetector() *LsofDetector {
	return &LsofDetector{}
}

func (d *LsofDetector) ListDevices() ([]DeviceInfo, error) {
	paths, err := filepath.Glob("/dev/video*")
	if err != nil {
		return nil, fmt.Errorf("enumerate video devices: %w", err)
	}

	var devices []DeviceInfo
	for _, path := range paths {
		info, ok, err := openAndQueryCap(path)
		if err != nil {
			continue
		}
		if ok {
			devices = append(devices, info)
		}
	}
	return devices, nil
}

func (d *LsofDetector) Detect(devicePath string) (DeviceStatus, error) {
	cmd := exec.Command("lsof", devicePath)
	err := cmd.Run()

	if err == nil {
		return DeviceStatus{On: true, CheckedAt: time.Now()}, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() == 1 {
			return DeviceStatus{On: false, CheckedAt: time.Now()}, nil
		}
	}

	return DeviceStatus{}, fmt.Errorf("lsof detection failed for %s: %w", devicePath, err)
}
