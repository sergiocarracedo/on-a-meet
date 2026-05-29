# Plan 08-01: macOS Detection Backend — Summary

**Wave:** 1
**Depends on:** None
**Objective:** MacOSDetector with log stream detection, system_profiler enumeration, factory support

## Tasks

| ID | Title | Files | Key action |
|----|-------|-------|------------|
| 08-01-01 | Create darwin.go | `internal/detector/darwin.go` | MacOSDetector — `log show` detection + `system_profiler` enumeration |
| 08-01-02 | Create darwin_stub.go | `internal/detector/darwin_stub.go` | !darwin stub returning errors |
| 08-01-03 | Update factory | `internal/detector/detector.go` | Add "darwin" case to New() |
| 08-01-04 | Add test | `internal/detector/detector_test.go` | TestNewDarwin |
| 08-01-05 | Verify build | `go.mod` | go build + vet pass |
