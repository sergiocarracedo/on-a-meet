# Plan 02-01 Summary

**Completed:** 2026-05-28

## What was built

V4L2 camera detection backend with device enumeration and on/off detection via syscalls. The `V4L2Detector` implements the `Detector` interface, using direct Linux syscalls (`unix.Open`/`unix.Close`, raw `syscall.Syscall` for ioctl) without external V4L2 libraries.

## Key files

- `internal/detector/interface.go`: `Detector` interface, `DeviceStatus`, `DeviceInfo` types
- `internal/detector/v4l2_linux.go`: `V4L2Detector` with `ListDevices()` and `Detect()` — V4L2 ioctl constants, `VIDIOC_QUERYCAP`-based device filtering, `EBUSY`→ON detection, `EACCES`→`O_RDONLY` fallback
- `internal/detector/v4l2_stub.go`: Non-Linux stub returning errors
- `internal/detector/v4l2_stub_test.go`: Tests for non-Linux stub

## Decisions made

- Used raw `syscall.Syscall(syscall.SYS_IOCTL, ...)` instead of `unix.Ioctl` (not available as generic function in this x/sys version)
- `vidiocQueryCap` computed via `_IOR` macro with `unsafe.Sizeof` (must be `var` not `const`)
- `openAndQueryCap` helper shares open+ioctl+close logic between `ListDevices` and filtering

## Notes for downstream

- The `Detect()` method uses `O_RDWR` with `EBUSY` detection — this is the standard V4L2 approach
- On `EACCES`, falls back to `O_RDONLY` to distinguish "no permission" from "device unavailable"
- Devices are filtered by `V4L2_CAP_VIDEO_CAPTURE` capability bit
