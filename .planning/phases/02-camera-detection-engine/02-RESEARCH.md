# Phase 2: Camera Detection Engine — Research

**Date:** 2026-05-28
**Sources:** Kernel V4L2 documentation, go4vl, Go V4L2 examples

## Don't Hand-Roll

- **V4L2 structs/ioctl codes:** Must define `v4l2_capability` struct and ioctl request codes (VIDIOC_QUERYCAP etc.) manually — `golang.org/x/sys/unix` does not include V4L2 constants. Alternatives:
  - `go4vl/v4l2` (v0.5.0, March 2026) provides full V4L2 bindings but adds unnecessary dependency for simple open/EBUSY + querycap
  - `thinkski/go-v4l2` (unmaintained, 2022) — avoid
  - Decision: hand-roll the minimal consts/structs needed (confirmed by CONTEXT.md)
- **go4vl:** Would be useful for streaming capture, but for on/off detection it's overkill. Simple syscalls suffice.

## Common Pitfalls

- **EACCES on O_RDWR**: Users without `video` group membership get permission denied. Fallback to O_RDONLY (allows capability query but not capture) as decided in CONTEXT.md. Provide clear fix instructions.
- **Not all /dev/videoN are cameras**: Many non-camera devices appear as /dev/videoN (TV tuners, codecs, ISP pipelines on Raspberry Pi). V4L2_CAP_VIDEO_CAPTURE filter is essential.
- **Device numbering changes**: /dev/video0 may become /dev/video1 after reboot. Polling engines must handle this gracefully (hotplug detect).
- **EBUSY semantics**: EBUSY only fires when a process holds the device open (e.g., Zoom, Chrome). Simply having the camera powered on but not claimed by any process returns O_RDWR success (camera is "off").
- **Race conditions**: Between open() checking EBUSY and querying VIDIOC_QUERYCAP, another process could grab the device. Handle transient errors gracefully — retry next cycle.
- **O_NONBLOCK not needed**: For simple open/close detection, blocking open is fine. O_NONBLOCK is for streaming read() loops.

## Existing Patterns in This Codebase

- **internal/detector/**: New package — no existing code. Must create interface.go (Detector interface, DeviceStatus, DeviceInfo types) and v4l2.go (V4L2Detector implementation).
- **internal/engine/**: New package — no existing code. Must create engine.go with polling loop, state tracking, hotplug, debounce.
- **internal/output/output.go**: Existing Banner() function shows device count. Will be extended with startup device listing. Info/Warning/Success wrappers available for state change lines.
- **internal/config/config.go**: Config struct has DetectMethod, Debounce, Interval, Camera fields — already Phase 2 ready.
- **cmd/detect.go**: Stub that prints "not yet implemented". Will be wired to create V4L2Detector + engine.Run() with OnChange callback.
- **cmd/list.go**: Already exists as a stub — Phase 2 may partially implement device listing (REQ-008 is Phase 5, but ListDevices() will be needed).

## Recommended Approach

1. **Define V4L2 constants in `internal/detector/v4l2_linux.go`**: Only what's needed — VIDIOC_QUERYCAP request code, v4l2_capability struct, V4L2_CAP_VIDEO_CAPTURE flag. Use golang.org/x/sys/unix for open(), close(), ioctl().

2. **Create `interface.go` with Detector interface**: DeviceStatus (On bool, CheckedAt time.Time), DeviceInfo (Path, Driver, Card, Bus), Detector interface (ListDevices, Detect). This is the contract.

3. **V4L2Detector implementation**:
   - `NewV4L2Detector() *V4L2Detector` — simple constructor
   - `ListDevices()` — glob /dev/video*, try open+VIDIOC_QUERYCAP, filter V4L2_CAP_VIDEO_CAPTURE, collect DeviceInfo
   - `Detect(path)` — open O_RDWR, if success → close, return OFF. If EBUSY → return ON. If EACCES → try O_RDONLY fallback.
   - Use build tag `//go:build linux` for the implementation file.

4. **Polling engine in `internal/engine/engine.go`**:
   - `Engine` struct with Detector, interval, debounce count, camera filter, OnChange callback
   - `Run(ctx context.Context)` — single goroutine loop
   - Per-device state tracking: map[string]*deviceState (current, previous, debounceCounter)
   - Each cycle: re-enumerate, compare for hotplug, poll each device, check debounce, fire OnChange
   - OR logic: track overall on/off, fire when any device transitions

5. **Wire in `cmd/detect.go`**:
   - Parse config, create V4L2Detector
   - Create engine with OnChange callback that prints to pterm
   - Call engine.Run(ctx) with signal-based context cancellation
   - Startup banner: show detected cameras
   - Permission check: if no devices accessible, print fix instructions and exit

6. **Testing**: Mock Detector interface for engine tests. V4L2-specific tests need /dev/video* or can be integration-only.
