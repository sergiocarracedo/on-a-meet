# Phase 8: macOS Detection Backend & Docs Polish — Verification

**Status:** pending

## Must-Have Checks

| Must-have | Status | Evidence |
|-----------|--------|----------|
| `darwin.go` exists with macOSDetector | ☐ | — |
| `darwin_stub.go` exists with !darwin build tag | ☐ | — |
| `detector.go` factory supports "darwin" | ☐ | — |
| `detector_test.go` tests "darwin" factory case | ☐ | — |
| `go build ./...` passes on Linux (stub) | ☐ | — |
| `go build ./...` passes on macOS (full) | ☐ | — |
| README has onboard command docs | ☐ | — |
| README has macOS install/permissions notes | ☐ | — |
| `go test ./...` passes | ☐ | — |

## Requirement Coverage

| Requirement | Status | Notes |
|-------------|--------|-------|
| REQ-014 (macOS detection) | ☐ | macOSDetector with log stream |
| REQ-015 (macOS enumeration) | ☐ | system_profiler SPCameraDataType |
| REQ-016 (factory support) | ☐ | detector.New("darwin") |
| REQ-017 (README docs) | ☐ | onboard + macOS notes |
