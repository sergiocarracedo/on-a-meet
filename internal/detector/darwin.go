//go:build darwin

package detector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type MacOSDetector struct {
	isOn bool
}

type profilerCamera struct {
	Name    string `json:"_name"`
	ModelID string `json:"spcamera_model-id"`
}

func NewMacOSDetector() *MacOSDetector {
	return &MacOSDetector{}
}

func (d *MacOSDetector) ListDevices() ([]DeviceInfo, error) {
	cmd := exec.Command("system_profiler", "SPCameraDataType", "-json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("system_profiler: %w", err)
	}

	var result struct {
		Cameras []profilerCamera `json:"SPCameraDataType"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse system_profiler output: %w", err)
	}

	var devices []DeviceInfo
	for _, c := range result.Cameras {
		devices = append(devices, DeviceInfo{
			Path:   c.Name,
			Card:   c.Name,
			Driver: c.ModelID,
			Bus:    "macOS",
		})
	}
	return devices, nil
}

func (d *MacOSDetector) Detect(devicePath string) (DeviceStatus, error) {
	cmd := exec.Command(
		"log", "show",
		"--last", "3s",
		"--predicate", `subsystem contains "com.apple.UVCExtension" and composedMessage contains "Post PowerLog"`,
		"--style", "compact",
	)
	out, err := cmd.Output()
	if err != nil {
		return DeviceStatus{On: d.isOn, CheckedAt: time.Now()}, nil
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		lower := strings.ToLower(line)
		if strings.Contains(lower, "stream on") || strings.Contains(lower, "streaming") {
			d.isOn = true
		} else if strings.Contains(lower, "stream off") || strings.Contains(lower, "stop") {
			d.isOn = false
		}
	}

	return DeviceStatus{On: d.isOn, CheckedAt: time.Now()}, nil
}
