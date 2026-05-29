# macOS Camera Detection — Architecture

## macOSDetector Design

```
internal/detector/
├── interface.go          # Existing Detector interface — no change needed
├── darwin.go             # NEW: macOSDetector — log stream + lsof fallback
├── darwin_test.go        # NEW: tests (require macOS)
├── v4l2_linux.go         # Existing Linux V4L2
├── lsof_linux.go         # Existing Linux lsof
├── v4l2_stub.go          # Existing stub
├── lsof_stub.go          # Existing stub
├── detector.go           # Existing factory — add "darwin" case
└── detector_test.go      # Existing factory tests — add "darwin" case
```

## Detection Flow

1. `Detect(devicePath)` on macOS:
   - Run `log show --predicate ...` or `log stream --timeout N`
   - Parse output for camera stream ON/OFF events
   - Return ON if any process is actively streaming
   - Return OFF if no stream detected
   - Fall back to `lsof` on USB camera IO services if log approach unreliable

2. `ListDevices()` on macOS:
   - Run `system_profiler SPCameraDataType` for built-in USB cameras
   - Or use `lsof -c VDCAssistant` as indicator
   - Or enumerate via IOKit registry (requires cgo)

## Integration

- `detector.New("darwin")` selects macOS backend
- Default detection method on darwin builds: `"darwin"`
- The existing `list` command works without changes — just plug into Detector interface
