# Phase 5: lsof Backend & Polish — Context

**Gathered:** 2026-05-29
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Implement lsof-based camera detection as a fallback backend, add runtime backend selection via factory pattern, implement `on-a-meet list` command with full device table including ON/OFF status, set up goreleaser for release builds, write README with install/config/usage docs, and refine CLI help text.

</domain>

<decisions>
## Implementation Decisions

### lsof Detection Implementation
- **Detection method:** Run `lsof /dev/videoN`, check exit code: 0 = ON (process has the device open), 1 = OFF (no process holds it). Parse output for process name and PID to include in logging.
- **Device enumeration:** lsof cannot enumerate devices — ListDevices() falls back to V4L2 enumeration (glob /dev/video* + VIDIOC_QUERYCAP filter).
- **File organization:** `lsof_linux.go` + `lsof_stub.go` in `internal/detector/` — same pattern as V4L2 backend. `LsofDetector` struct with `NewLsofDetector()` constructor.

### Backend Selection & Factory
- **Factory function:** `detector.New(method string) (Detector, error)` in `internal/detector/detector.go`. Accepts "v4l2" or "lsof", returns the appropriate Detector.
- **Config integration:** `detect.go` calls `detector.New(cfg.DetectMethod)` instead of hardcoded `NewV4L2Detector()`. Config already has `detect-method` field (default "v4l2").
- **--detect flag:** Already wired via config's `detect-method` field — no separate flag needed unless CLI override is desired (can be added to detect and list commands).

### List Command UX
- **Table columns:** Path, Driver, Card, Bus, Status (ON/OFF). Status determined by calling `Detect()` on each device.
- **No cameras:** Print "No camera devices detected" message with video group instructions (reuse existing output pattern from detect.go).
- **Backend:** Uses the same detector backend as detect command (from config).

### Release Scripting
- **Tool:** goreleaser with `.goreleaser.yaml` configuration.
- **Cross-compilation targets:** linux/amd64, linux/arm64, darwin/amd64, darwin/arm64.
- **Checksums:** Included in release artifacts.

### README & --help Refinement
- **README:** Install instructions (video group, sudo), config reference, usage examples (`--on`/`--off` with templates), service install/uninstall docs.
- **--help:** Updated descriptions for all flags across all commands. Man page not needed for v1.

</decisions>

<specifics>
## Specific Ideas
- List command shows Status column by running Detect() — gives immediate feedback on which cameras are in use
- Factory function makes it trivial to add future backends (udev, macOS AVFoundation)
- lsof output parsing for process name helps users identify which app is using the camera
- goreleaser enables proper GitHub releases with cross-compiled binaries

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — REQ-013 (lsof backend), REQ-008 (list command)
- `.planning/research/ARCHITECTURE.md` — Detector interface, component diagram
- `internal/detector/interface.go` — Detector interface definition
- `internal/detector/v4l2_linux.go` — Existing V4L2 implementation (pattern to follow)
- `cmd/detect.go` — Detect command (backend selection integration point)
- `cmd/list.go` — List command stub (to be implemented)
- `internal/config/config.go` — Config struct with DetectMethod field
- `.planning/phases/02-camera-detection-engine/02-CONTEXT.md` — Detector interface decisions from Phase 2

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/detector/interface.go` — Detector interface, DeviceStatus, DeviceInfo types
- `internal/detector/v4l2_linux.go` — V4L2Detector with ListDevices() (reused by LsofDetector for enumeration) and Detect() (pattern for lsof implementation)
- `internal/detector/v4l2_stub.go` — Non-Linux stub pattern (follow for lsof stub)
- `cmd/list.go` — Existing stub with cobra boilerplate
- `internal/output/output.go` — pterm wrapper functions (Table, Info, Warning, Error)
- `internal/config/config.go` — Config struct with DetectMethod field

### Established Patterns
- `internal/detector/` package with `_linux.go` + `_stub.go` build tags for platform-specific code
- pterm wrappers for all terminal output (respects --silent/--verbose)
- Test files alongside source (`_test.go` next to `.go`)

### Integration Points
- `internal/detector/detector.go` — New factory file for `New(method string)` function
- `cmd/detect.go:42` — Replace `detector.NewV4L2Detector()` with `detector.New(cfg.DetectMethod)`
- `cmd/list.go` — Complete stub, wire detector + output table
- `.goreleaser.yaml` — New file at project root

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 05-lsof-backend-polish*
*Context gathered: 2026-05-29*
