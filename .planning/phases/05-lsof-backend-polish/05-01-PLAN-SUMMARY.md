# Plan 05-01 Summary

**Completed:** 2026-05-29

## What was built
LsofDetector backend and factory function. Detection runs `lsof /dev/videoN` (exit 0=ON, exit 1=OFF). Device enumeration falls back to V4L2 openAndQueryCap. Factory `detector.New(method)` supports "v4l2" and "lsof".

## Key files
- `internal/detector/lsof_linux.go`: LsofDetector with Detect() and ListDevices()
- `internal/detector/lsof_stub.go`: Non-Linux stub
- `internal/detector/detector.go`: New(method) factory
- `internal/detector/detector_test.go`: 5 tests

## Decisions made
- lsof exit code 1 = OFF (distinguished from other errors via ExitError check)
- ListDevices reuses existing openAndQueryCap from v4l2_linux.go (same package)
