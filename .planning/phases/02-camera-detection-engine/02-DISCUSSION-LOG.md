---
phase: 2
slug: camera-detection-engine
areas_discussed:
  - V4L2 Detection Approach
  - Detector Interface Shape
  - Polling Engine Architecture
  - Terminal Output UX
created: 2026-05-28
---

# Phase 2: Camera Detection Engine - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-05-28
**Phase:** 02-camera-detection-engine
**Areas discussed:** V4L2 Detection Approach, Detector Interface Shape, Polling Engine Architecture, Terminal Output UX

---

## V4L2 Detection Approach

### Decision 1: On-detection method

| Option | Description | Selected |
|--------|-------------|----------|
| Open O_RDWR, EBUSY = on | Try opening /dev/videoN. EBUSY means in use. Simple, proven. | ✓ (Recommended) |
| Open + streaming state check | Query VIDIOC_STREAMON for active streaming. More precise but complex. | |
| Combination | Check EBUSY first, then streaming as verification. Covers edge cases. | |

**User's choice:** Open O_RDWR, EBUSY = on
**Notes:** Accepted recommended option. Simplest and most proven approach.

### Decision 2: Permission fallback

| Option | Description | Selected |
|--------|-------------|----------|
| Fall back to O_RDONLY | O_RDWR EACCES → try O_RDONLY to distinguish "no permission" from "device available" | ✓ (Recommended) |
| Treat EACCES as error | If O_RDWR fails, report permission error immediately | |

**User's choice:** Fall back to O_RDONLY
**Notes:** Accepted recommended option.

### Decision 3: Device filtering

| Option | Description | Selected |
|--------|-------------|----------|
| Filter by V4L2_CAP_VIDEO_CAPTURE | Only monitor actual cameras via VIDIOC_QUERYCAP filtering | ✓ (Recommended) |
| Monitor all /dev/video* | Poll every entry, no capability filtering | |

**User's choice:** Filter by capture capability
**Notes:** Accepted recommended option. Prevents false positives from non-camera video devices.

---

## Detector Interface Shape

### Decision 1: Interface richness

| Option | Description | Selected |
|--------|-------------|----------|
| Simple: Detect only | `interface { Detect(path) (bool, error) }` + separate `ListDevices()` function | |
| Rich: Detect + ListDevices + DeviceInfo | Interface includes both + `DeviceInfo` struct with Path, Driver, Card, Bus | ✓ (Recommended) |

**User's choice:** Rich interface (Detect + ListDevices + DeviceInfo)
**Notes:** User preferred the richer interface, keeping backend logic self-contained.

### Decision 2: Detect return type

| Option | Description | Selected |
|--------|-------------|----------|
| Simple bool | `Detect(path) (bool, error)` | |
| Result struct | `Detect(path) (DeviceStatus, error)` where DeviceStatus has State + Timestamp | ✓ (Recommended) |

**User's choice:** Result struct (DeviceStatus with On + CheckedAt)
**Notes:** User chose the richer struct for future flexibility.

---

## Polling Engine Architecture

### Decision 1: Goroutine model

| Option | Description | Selected |
|--------|-------------|----------|
| Single goroutine, sequential | One loop polls all devices sequentially, sleeps between cycles | ✓ (Recommended) |
| Per-device goroutines | Each device gets own goroutine with independent ticker | |

**User's choice:** Single goroutine, sequential
**Notes:** Accepted recommended option. N=1-3 cameras doesn't need concurrency.

### Decision 2: Change notification mechanism

| Option | Description | Selected |
|--------|-------------|----------|
| Callback | `type OnChange func(path string, old, new bool)` passed to engine | ✓ (Recommended) |
| Channel-based event bus | Events sent on channel, consumer goroutine processes | |
| Direct coupling | Engine handles output and command dispatch internally | |

**User's choice:** Callback function
**Notes:** Accepted recommended option. Simple, testable, works for both Phase 2 (output) and Phase 3 (command dispatch).

### Decision 3: Graceful shutdown

| Option | Description | Selected |
|--------|-------------|----------|
| Context cancellation | `Run(ctx context.Context)` | ✓ (Recommended) |
| Stop() + done channel | Explicit Stop() method + Done() channel | |

**User's choice:** Context cancellation
**Notes:** Accepted recommended option. Standard Go pattern.

---

## Terminal Output UX

### Decision 1: Display format

| Option | Description | Selected |
|--------|-------------|----------|
| Hybrid: banner + scrolling log | Startup pterm banner + scrolling lines on state change | ✓ (Recommended) |
| Live-updating table | pterm table refreshed in-place | |
| Minimal: state changes only | No banner, just state change lines | |

**User's choice:** Hybrid (banner + scrolling log)
**Notes:** Accepted recommended option. Pipe-friendly and informative.

### Decision 2: Log line info

| Option | Description | Selected |
|--------|-------------|----------|
| Device path + state + driver | e.g. `/dev/video0 → ON (driver: uvcvideo)` | ✓ (Recommended) |
| Device path + state only | e.g. `/dev/video0 → ON` | |
| Index + state + device | e.g. `[1] /dev/video0 → ON (uvcvideo)` | |

**User's choice:** Device path + state + driver
**Notes:** Accepted recommended option. Driver name helps identify cameras in multi-camera setups.

### Decision 3: Hotplug display

| Option | Description | Selected |
|--------|-------------|----------|
| Log add/remove events | `[+] /dev/video2 detected`, `[-] /dev/video1 disconnected` | ✓ (Recommended) |
| Silent handling | Don't print device changes, only state transitions | |
| Re-scan banner on change | Reprint entire banner when device list changes | |

**User's choice:** Log add/remove events
**Notes:** Accepted recommended option.

---

## Agent's Discretion

- Exact pterm styling and color choices for hotplug logs (Info green for add, Warning yellow for remove)
- Specific polling engine constructor signature and configuration
- Device re-scanning implementation details (full re-enumerate each cycle vs. periodic)
- --camera flag validation (path format check)

## Deferred Ideas

None — discussion stayed within phase scope.

---

*Phase: 02-camera-detection-engine*
*Discussion log generated: 2026-05-28*
