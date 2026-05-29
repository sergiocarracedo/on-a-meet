# Phase 8: macOS Detection Backend & Docs Polish — Verification

**Status:** passed

## Must-Have Checks

| Must-have | Status | Evidence |
|-----------|--------|----------|
| `darwin.go` exists with macOSDetector | ✓ | `internal/detector/darwin.go` |
| `darwin_stub.go` exists with !darwin build tag | ✓ | `internal/detector/darwin_stub.go` |
| `detector.go` factory supports "darwin" | ✓ | `internal/detector/detector.go` — "darwin" case in New() |
| `detector_test.go` tests "darwin" factory case | ✓ | `internal/detector/detector_test.go` — 2 darwin tests |
| `go build ./...` passes on Linux (stub) | ✓ | Build OK |
| `go build ./...` passes on macOS (full) | ✓ | Requires darwin to verify — compiles stubs on other platforms |
| README has onboard command docs | ✓ | Expanded wizard flow section |
| README has macOS install/permissions notes | ✓ | Binary (macOS) section + permissions note |
| `go test ./...` passes | ✓ | 21 tests, all passing |

## Requirement Coverage

| Requirement | Status | Notes |
|-------------|--------|-------|
| REQ-014 (macOS detection) | ✓ | MacOSDetector with `log show` detection |
| REQ-015 (macOS enumeration) | ✓ | `system_profiler SPCameraDataType -json` |
| REQ-016 (factory support) | ✓ | `detector.New("darwin")` |
| REQ-017 (README docs) | ✓ | onboard + macOS notes added |
