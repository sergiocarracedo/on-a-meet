# Phase 2: Camera Detection Engine - Context

**Gathered:** 2026-05-28
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Implement V4L2-based camera on/off detection backend, device enumeration with capture-capability filtering, polling engine with state tracking, debounce (N consecutive same-state polls), multi-camera OR logic, --camera filter, hotplug handling, and terminal output showing state changes. Phase 3 handles command execution — this phase only shows state changes in terminal output.

</domain>

<decisions>
## Implementation Decisions

### V4L2 Detection Approach
- **On detection:** Open `/dev/videoN` with `O_RDWR`. If open succeeds → camera is off (close immediately). If `EBUSY` → another process holds the device → camera is on.
- **Permission fallback:** If `O_RDWR` fails with `EACCES`, fall back to `O_RDONLY`. If `O_RDONLY` succeeds, device exists and is off but user has insufficient permissions (warn). If `O_RDONLY` also fails, it's a real permission error.
- **Device filtering:** Use `VIDIOC_QUERYCAP` to filter devices with `V4L2_CAP_VIDEO_CAPTURE` capability. Exclude non-camera video devices (TV tuners, encoders, etc.).
- **Implementation:** Direct Go syscall (`golang.org/x/sys/unix`) — no go4vl dependency.

### Detector Interface Shape
- **Interface:** Rich interface with both `Detect()` and `ListDevices()`:

```go
type DeviceStatus struct {
    On        bool
    CheckedAt time.Time
}

type DeviceInfo struct {
    Path   string
    Driver string
    Card   string
    Bus    string
}

type Detector interface {
    ListDevices() ([]DeviceInfo, error)
    Detect(devicePath string) (DeviceStatus, error)
}
```

- **Package location:** `internal/detector/` with `interface.go` (interface definitions), `v4l2.go` (V4L2 implementation), `device_linux.go` (platform-specific).
- **V4L2 implementation:** Concrete `V4L2Detector` struct with constructor `NewV4L2Detector()`.

### Polling Engine Architecture
- **Goroutine model:** Single goroutine, sequential poll — all devices iterated per cycle, sleeps for `interval` between cycles. N=1-3 cameras makes concurrency unnecessary.
- **State change notification:** Callback function:
  ```go
  type OnChange func(path string, oldState, newState bool)
  ```
- **Graceful shutdown:** `Run(ctx context.Context)` with context cancellation. No explicit `Stop()` method.
- **State tracking:** Per-device struct with current state, previous state, and debounce counter. Map of `devicePath -> deviceState`.
- **Debounce:** Configurable N consecutive same-state polls before firing the callback (default: 3). Any mismatch resets the counter.
- **Multi-camera OR logic:** In the callback dispatch — if any device transitions to ON, overall state is ON. If all devices are OFF, overall state is OFF.
- **Package location:** `internal/engine/` with `engine.go`.
- **Hotplug:** On each poll cycle, re-enumerate devices. Detect added/removed devices by comparing against known set. New devices tracked from initial state. Removed devices removed from state map. Log add/remove events via callback.
- **`--camera` filtering:** At engine start, filter device list to matching path. Only poll the specified device.

### Terminal Output UX
- **Format:** Hybrid — startup banner + scrolling log.
  - **Startup:** pterm banner listing all detected cameras (path, driver, initial state).
  - **Scrolling log:** Print a line on each debounced state change.
- **Log line format:** `/dev/video0 ⟶ ON  (driver: uvcvideo)`
- **Hotplug display:** Log `[+] /dev/video2 detected (uvcvideo)` on device add, `[-] /dev/video1 disconnected` on device remove. Uses pterm Info/Warning styling.
- **Source of log lines:** The `detect` commands the engine's `OnChange` callback to pterm output wrappers.
- **--silent flag:** Suppress all output when set.

### Permission Check
- **Startup check:** Before entering poll loop, verify each device is accessible (or permission diagnostics are clear). If no devices are accessible, print error with fix instructions (add user to `video` group, check udev rules, or run with sudo).

</decisions>

<specifics>
## Specific Ideas

- Show driver name in startup banner and state change logs — helps users identify which camera is which in multi-camera setups
- Re-enumerate devices each poll cycle for hotplug detection — simple, no udev dependency for Phase 2
- `V4L2_CAP_VIDEO_CAPTURE` filter prevents false positives from non-camera /dev/video* devices

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — REQ-001, REQ-004, REQ-005, REQ-006, REQ-010, REQ-011, REQ-012
- `.planning/research/ARCHITECTURE.md` — Detector interface, engine design, component diagram
- `.planning/phases/01-project-scaffold-cli-foundation/01-CONTEXT.md` — Project layout, config schema, output patterns
- `internal/config/config.go` — Config struct (detect-method, debounce, interval fields)
- `cmd/detect.go` — Detect command stub with existing flags (--camera, --interval, --on, --off)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/output/` — pterm wrapper functions (Info, Success, Warning, Error, Table, StartupBanner) for all terminal output
- `internal/config/config.go` — `Config` struct with `DetectMethod`, `Debounce`, `Interval`, `Camera` fields already defined
- `cmd/detect.go` — Detect command stub, flags already registered with `viper.BindPFlag`

### Established Patterns
- Package-per-concern under `internal/` — new `internal/detector/` and `internal/engine/` packages follow convention
- pterm wrappers used instead of raw pterm calls to respect --silent/--verbose flags
- Test files alongside source (`config_test.go` next to `config.go`)

### Integration Points
- `cmd/detect.go` `RunE` → creates engine, passes config, starts `Run(ctx)`, wires `OnChange` callback to output
- `internal/detector/` → new package, no existing integration
- `internal/engine/` → new package, no existing integration
- Output wrapper functions from `internal/output/` will be used for terminal display

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 02-camera-detection-engine*
*Context gathered: 2026-05-28*
