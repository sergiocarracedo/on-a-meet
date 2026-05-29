# Phase 8: macOS Detection Backend — Discussion Log

## 2026-05-29 — Planning session

- Decision: Use `log stream` as primary detection method (follows existing exec-based pattern)
- Decision: Use `system_profiler SPCameraDataType -json` for device enumeration
- Decision: Pure Go, no cgo, exec-based (same as lsof backend)
- Decision: `darwin.go` + `darwin_stub.go` pattern (same as v4l2/lsof platform pairs)
- Decision: State tracking in memory — initial state is OFF, updated on transitions
- Decision: Wave 01-01 = macOSDetector implementation; Wave 01-02 = README docs polish
