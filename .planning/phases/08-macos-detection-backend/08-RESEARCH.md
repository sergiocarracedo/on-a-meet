# Phase 8: macOS Detection Backend — Research

## Key Findings

### `log stream` Approach (Primary)
- Command: `log stream --timeout 2 --predicate 'subsystem contains "com.apple.UVCExtension"' --style compact`
- macOS unified log captures camera stream start/stop events
- Fragile across macOS versions — predicates may change
- No cgo, no external dependencies
- `log show --last Ns` for historical window check

### `system_profiler` for Enumeration
- Command: `system_profiler SPCameraDataType -json`
- Returns JSON with camera model, name, UID
- No cgo, built into macOS
- Covers both built-in and USB cameras

### `lsof` Supplementary Check
- `lsof -c VDCAssistant` may show active camera sessions
- Not all camera usage goes through VDCAssistant on modern macOS
- Use as fallback only

### Limitations
- log stream only shows *transitions* not absolute state
- First poll after startup assumes OFF (document this)
- Different macOS versions may have different log event formats

See `.planning/research/` for full research materials.
