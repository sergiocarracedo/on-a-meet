# Requirements — on-a-meet v1.1.0

## v1.1.0 Requirements

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| REQ-014 | macOS detection backend — detect camera on/off state on macOS using `log stream` | P1 | Follow existing exec-based pattern (like lsof backend). Call `log stream --predicate` and parse output for camera stream events. |
| REQ-015 | macOS device enumeration — list built-in and USB cameras on macOS | P1 | Use `system_profiler SPCameraDataType` (no cgo) or enumerate via IOKit if needed. |
| REQ-016 | Backend factory support for darwin — `detector.New("darwin")` selects macOSDetector | P1 | Add "darwin" detection method to existing factory. Default on darwin builds. |
| REQ-017 | README documentation — add onboard command docs and macOS install/permissions notes | P1 | Fix the missing `onboard` docs gap. Add macOS-specific notes. |

## v2 Requirements (Next Milestone Candidates)

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| — | macOS AVFoundation cgo backend | P2 | More reliable but requires cgo + Xcode |
| — | udev detection backend | P2 | netlink-based for hotplug events on Linux |

## Out of Scope

| Item | Reasoning |
|------|-----------|
| Windows support | Not requested. Different API (DirectShow/MediaFoundation). |
| GUI or web dashboard | CLI/daemon tool only. |
| Real eBPF event detection | Too complex. Polling is sufficient for v1. |
| macOS green light hardware API | No public API exists. |
