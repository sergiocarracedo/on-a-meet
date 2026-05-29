# Phase 8: macOS Detection Backend & Docs Polish — Context

**Gathered:** 2026-05-29
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Implement macOS camera detection via `log stream` and `system_profiler` as a new Detector backend, add "darwin" detection method to the factory (default on darwin builds), add darwin platform support in goreleaser, and document the onboard command and macOS setup in README.

</domain>

<decisions>
## Implementation Decisions

### macOSDetector Implementation

- **Detection method:** Use `log stream --timeout 2 --predicate 'subsystem contains "com.apple.UVCExtension"' --style compact` — capture camera stream events. Track state transitions in memory: start event → ON, stop event → OFF. Initial state assumed OFF (polling captures subsequent transitions).
- **Device enumeration:** Use `system_profiler SPCameraDataType -json` — parse JSON output for device name, model, UID.
- **Fallback:** Use `lsof -c VDCAssistant` as supplementary check if log stream approach yields no events.
- **File organization:** `darwin.go` + `darwin_stub.go` in `internal/detector/` — following the platform-stub pattern used by V4L2 and lsof backends.
- **Detect("") behavior:** On macOS, empty device path checks ALL cameras and returns ON if any camera is active.

### Factory Integration

- Add `"darwin"` case to `detector.New()`.
- On darwin builds, default detection method = `"darwin"` (set in root.go config defaults).

</decisions>

<specifics>
## Specific Ideas
- log stream approach is inherently reactive (transitions, not absolute state) — document this limitation
- system_profiler provides built-in + USB camera enumeration without cgo
- darwin_stub.go returns errors on non-macOS (same pattern as v4l2_stub.go, lsof_stub.go)
</specifics>

<canonical_refs>
## Canonical References

- `.planning/REQUIREMENTS.md` — REQ-014 (macOS detection), REQ-015 (enumeration), REQ-016 (factory), REQ-017 (README)
- `.planning/research/ARCHITECTURE.md` — macOSDetector design, detector interface
- `.planning/research/PITFALLS.md` — macOS pitfalls (no /dev/video*, appleh13camerad, log stream fragility)
- `internal/detector/interface.go` — Detector interface definition
- `internal/detector/lsof_linux.go` — Existing backend pattern (exec-based, no cgo)
- `internal/detector/detector.go` — Factory function (add "darwin" case)
- `cmd/detect.go` — Detection command (backend selection)
- `cmd/list.go` — List command (device enumeration)
</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/detector/interface.go` — Detector interface, DeviceStatus, DeviceInfo types
- `internal/detector/v4l2_stub.go` — Non-Linux stub pattern (follow for darwin stub)
- `internal/detector/lsof_stub.go` — Another stub pattern reference
- `internal/detector/detector.go` — Factory function to extend

### Established Patterns
- `_linux.go` / `_stub.go` build-tag pairs → now `_darwin.go` / `_stub.go`
- pterm wrappers for terminal output
- Test files alongside source

### Integration Points
- `internal/detector/detector.go` — Add "darwin" case
- `cmd/detect.go` — No change needed (already uses detector.New())
- `cmd/list.go` — No change needed
- `.goreleaser.yaml` — Already builds darwin/amd64 + darwin/arm64 (no change)
- Config defaults — default detect method should be "darwin" on darwin builds
</code_context>

<deferred>
## Deferred Ideas

- AVFoundation cgo backend for more reliable detection (P2)
- IOKit-based device enumeration for more detail
- macOS green light hardware query (no public API)
</deferred>

---

*Phase: 08-macos-detection-backend*
*Context gathered: 2026-05-29*
