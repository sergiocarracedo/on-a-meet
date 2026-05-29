---
objective: "Fix V4L2 detection: when open(O_RDWR) succeeds, check /proc for open handles to detect camera in use"
must_haves:
  truths:
    - "internal/detector/v4l2_linux.go Detect() returns ON when open succeeds but /proc shows open handles"
    - "internal/detector/v4l2_linux.go Detect() returns OFF when open succeeds and /proc shows no open handles"
    - "go build passes, all tests pass"
  artifacts:
    - internal/detector/v4l2_linux.go
  key_links:
    - "Detect() uses hasOpenHandle() to check /proc for processes holding the device open"
---

# Quick Task 003: Fix V4L2 Camera ON Detection

## Tasks

<task id="003-01">
<title>Fix V4L2 Detect() to check /proc for open handles</title>
<files>
- internal/detector/v4l2_linux.go
</files>
<action>
Modify `internal/detector/v4l2_linux.go`:

1. Add `"os"` to imports (for reading /proc entries).

2. Add a helper function `hasOpenHandle(path string) bool` that scans `/proc/*/fd/` for symlinks pointing to the device path:

```go
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
```

3. Modify `Detect()` — when `open(O_RDWR)` succeeds, check `/proc` for open handles before declaring OFF:

```go
func (d *V4L2Detector) Detect(devicePath string) (DeviceStatus, error) {
    fd, err := unix.Open(devicePath, unix.O_RDWR, 0)
    if err == nil {
        unix.Close(fd)
        if hasOpenHandle(devicePath) {
            return DeviceStatus{On: true, CheckedAt: time.Now()}, nil
        }
        return DeviceStatus{On: false, CheckedAt: time.Now()}, nil
    }
    // ... rest unchanged
}
```

4. Update the package's import block to include `"os"`, `"path/filepath"`.
</action>
<verify>
grep -q "hasOpenHandle" internal/detector/v4l2_linux.go && grep -q "/proc" internal/detector/v4l2_linux.go
</verify>
<done>
V4L2 Detect() uses /proc scan to detect camera in use when open succeeds
</done>
</task>
