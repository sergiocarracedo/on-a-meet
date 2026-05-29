# macOS Camera Detection — Pitfalls

## Don't Hand-Roll

1. **Don't use `appleh13camerad`** — it runs constantly on Apple Silicon. Not a reliable indicator of camera-in-use.
2. **Don't poll AVCaptureDevice** — opening the camera for detection forces the green light on, defeating the purpose.
3. **Don't parse `ioreg`** — IOKit registry access has changed significantly across macOS versions. Fragile.

## Common Pitfalls

### 1. Fragile log predicates
The `log stream` predicate `com.apple.UVCExtension` changed between Big Sur and Monterey. Apple changes these without notice. Build version-specific fallbacks.

### 2. Permission model
macOS requires camera permission (`NSCameraUsageDescription` in Info.plist) for ANY camera access. The `log` approach bypasses this (reads system logs, not hardware), but `lsof` and IOKit approaches may be blocked.

### 3. No /dev/video* on macOS
macOS cameras don't appear as device files. The existing Linux device path logic will not work. Must use different device discovery.

### 4. Apple Silicon vs Intel differences
M-series Macs have different camera hardware and daemon behavior. Test on both architectures.

### 5. USB camera naming
USB cameras appear as CoreMediaIO DAL plugins, not as /dev/video*. Device naming is inconsistent.

## Recommended Approach

- **Primary:** `log stream --predicate 'subsystem contains "com.apple.UVCExtension" and composedMessage contains "Post PowerLog"'`
- **Fallback:** `lsof` on camera-related IOKit services or `VDCAssistant` process
- **Device list:** `system_profiler SPCameraDataType` for camera enumeration
- **Keep it simple:** Pure Go, no cgo, exec-based

Approach is same pattern as Linux lsof backend — lightweight, external command parsing, no cgo.
