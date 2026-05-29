# Requirements — on-a-meet

## v1 Requirements (Current Milestone)

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| REQ-001 | V4L2 polling detection — detect camera on/off state via /dev/videoN ioctl | P1 | Use Go syscall directly (not go4vl). Check VIDIOC_QUERYCAP + device open status. |
| REQ-002 | Command execution on state change — --on and --off flags | P1 | Spawn in goroutine with timeout. Must not block polling loop. |
| REQ-003 | Template variable substitution in commands — {{.CameraID}}, {{.Device}}, {{.State}} | P1 | Go text/template. Available in both --on and --off. |
| REQ-004 | Multi-camera support — OR logic: any camera on = on, all off = off | P1 | Default behavior. --camera flag overrides to single device. |
| REQ-005 | --camera flag — target a specific camera by device path | P1 | Filter detected devices to matching path. |
| REQ-006 | --interval flag — configurable polling interval | P1 | Default: 1 second. Minimum: 100ms. |
| REQ-007 | YAML config file — all flags settable via config file | P1 | Viper-based. Path: --config flag or default locations (~/.config/on-a-meet/config.yaml). CLI flags override config values. |
| REQ-008 | Device list command — list available cameras | P1 | Enumerate /dev/video*. Show path, driver info. pterm table output. |
| REQ-009 | systemd service installation — --install and --uninstall flags | P1 | kardianos/service. Auto-detect OS, generate correct unit. Prompt for sudo if needed. |
| REQ-010 | Debounce window — prevent false positives from transient device queries | P1 | Configurable (default: 3 consecutive same-state polls before firing). |
| REQ-011 | Permission check at startup — verify /dev/video* accessibility | P1 | Clear error message + fix instructions if not accessible. |
| REQ-012 | Graceful camera hotplug — handle disconnected cameras without crash | P1 | Handle ENOENT per-device. Log warning, continue polling remaining devices. |

## v2 Requirements (Next Milestone Candidates)

| ID | Requirement | Priority | Notes |
|----|-------------|----------|-------|
| REQ-013 | lsof detection backend — fallback detection method | P2 | Most reliable method. Parses lsof output for /dev/video* ownership. User-selectable via --detect flag. |
| REQ-014 | macOS AVFoundation detection backend | P2 | Use AVFoundation or IOKit. Requires CGo or pure Go alternative. |
| REQ-015 | udev detection backend — monitor hotplug events | P2 | netlink-based udev monitoring for device add/remove events. |
| REQ-016 | Logging to file/stdout — --log flag | P2 | Structured logging (slog) with levels. Log to file when running as service. |
| REQ-017 | Multiple detection methods simultaneously | P3 | Vote-based: if 2/3 backends agree, fire the event. |

## Out of Scope

| Item | Reasoning |
|------|-----------|
| Windows support | Not requested. Would require DirectShow/MediaFoundation. Can be added if needed. |
| GUI or web dashboard | CLI/daemon tool only. Users wire their own dashboards via --on/--off commands. |
| Real eBPF event detection | Too complex for v1. Kernel dependency, root, BTF requirements. Polling is sufficient. |
| Built-in home assistant integration | Out of scope. Users integrate via the --on/--off curl commands. |
| Video streaming or recording | Camera detection only. Not a video processing tool. |
| Automatic camera detection for libcamera (RPi) | libcamera uses different API. Use lsof backend fallback if needed. |

---
*Last updated: 2026-05-28*
