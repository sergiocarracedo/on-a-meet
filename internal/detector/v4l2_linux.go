//go:build linux

package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

var vidiocQueryCap = _IOR('V', 0, unsafe.Sizeof(v4l2Capability{}))

func _IOR(t byte, nr byte, size uintptr) uint32 {
	return (2 << 30) | (uint32(size) << 16) | (uint32(t) << 8) | uint32(nr)
}

type v4l2Capability struct {
	Driver       [16]byte
	Card         [32]byte
	BusInfo      [32]byte
	Version      uint32
	Capabilities uint32
	DeviceCaps   uint32
	Reserved     [3]uint32
}

const (
	v4l2CapVideoCapture = 0x00000001
)

type V4L2Detector struct{}

func NewV4L2Detector() *V4L2Detector {
	return &V4L2Detector{}
}

func (d *V4L2Detector) ListDevices() ([]DeviceInfo, error) {
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

func openAndQueryCap(path string) (DeviceInfo, bool, error) {
	fd, err := unix.Open(path, unix.O_RDWR, 0)
	if err != nil {
		if err == unix.EACCES {
			fd, err = unix.Open(path, unix.O_RDONLY, 0)
		}
		if err != nil {
			return DeviceInfo{}, false, err
		}
	}
	defer unix.Close(fd)

	var cap v4l2Capability
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(vidiocQueryCap), uintptr(unsafe.Pointer(&cap))); err != 0 {
		return DeviceInfo{}, false, err
	}

	if cap.Capabilities&v4l2CapVideoCapture == 0 {
		return DeviceInfo{}, false, nil
	}

	return DeviceInfo{
		Path:   path,
		Driver: nullTerminatedString(cap.Driver[:]),
		Card:   nullTerminatedString(cap.Card[:]),
		Bus:    nullTerminatedString(cap.BusInfo[:]),
	}, true, nil
}

func nullTerminatedString(b []byte) string {
	for i, v := range b {
		if v == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

func hasOpenHandle(devicePath string) bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pid := entry.Name()
		if pid[0] < '0' || pid[0] > '9' {
			continue
		}
		fdDir := filepath.Join("/proc", pid, "fd")
		fds, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range fds {
			linkPath := filepath.Join(fdDir, fd.Name())
			link, err := os.Readlink(linkPath)
			if err != nil {
				continue
			}
			if link == devicePath {
				return true
			}
		}
	}
	return false
}

func (d *V4L2Detector) Detect(devicePath string) (DeviceStatus, error) {
	fd, err := unix.Open(devicePath, unix.O_RDWR, 0)
	if err == nil {
		unix.Close(fd)
		if hasOpenHandle(devicePath) {
			return DeviceStatus{On: true, CheckedAt: time.Now()}, nil
		}
		return DeviceStatus{On: false, CheckedAt: time.Now()}, nil
	}

	switch err {
	case unix.EBUSY:
		return DeviceStatus{On: true, CheckedAt: time.Now()}, nil
	case unix.EACCES:
		fd, err = unix.Open(devicePath, unix.O_RDONLY, 0)
		if err == nil {
			unix.Close(fd)
			return DeviceStatus{On: false, CheckedAt: time.Now()}, nil
		}
		return DeviceStatus{}, fmt.Errorf("cannot access %s: %w (try: sudo usermod -a -G video $USER)", devicePath, err)
	default:
		return DeviceStatus{}, fmt.Errorf("cannot detect %s: %w", devicePath, err)
	}
}
