# on-a-meet

## What This Is

A CLI tool that monitors camera state (on/off) on Linux and macOS, and triggers user-defined commands when state changes. Users wire camera activity to home automation, notifications, or any custom workflow — with no infrastructure beyond a single binary.

## Core Value

Reliably detect camera state changes and fire the correct command every time. When the camera turns on, the on-command runs. When it turns off, the off-command runs. Fast enough to feel immediate.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Camera detection via multiple backends: V4L2, lsof (process monitoring), udev — user-selectable via flag or YAML config, with a sensible default
- [ ] Command spawning with template variables (`{{.CameraID}}`, `{{.Device}}`, `{{.State}}`)
- [ ] Multi-camera support: OR logic by default (any camera on = on), user can target a specific camera via `--camera` flag
- [ ] Config file support (YAML) — CLI flags override config values
- [ ] Service installation: auto-generate systemd (Linux) / launchd (macOS), prompt for sudo elevation when needed
- [ ] Polling-based detection (configurable interval, default reasonable)

### Out of Scope

- Windows support — not in scope for v1. Can be added later if needed.
- macOS detection implementation — design for it, but actual macOS detection is deferred.
- GUI, web UI, or dashboard — CLI only.

## Context

The project is named `on-a-meet` after the primary use case: knowing when someone is on a video call. The tool polls camera device state at a configurable interval, compares it to the last known state, and executes user-provided commands when a transition occurs.

Platform detection differs significantly:
- **Linux**: `/dev/video*` devices via V4L2 ioctls, process monitoring via `lsof`, or udev events
- **macOS**: AVFoundation or IOKit (deferred)

The user wants polling (simplest implementation) as the default, with detection method selectable via flag or config.

## Constraints

- **Language**: Go — binary distribution, single-file output, cross-compilation
- **Output**: pterm — consistent terminal output formatting
- **OS**: Linux first (primary), macOS second
- **Service**: systemd on Linux, launchd on macOS — auto-generate unit files

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Polling (not event-driven) | Simplest implementation, fewer edge cases | ✓ Good |
| Config format: YAML | Familiar, readable, common in Go ecosystem | — Pending |
| Template variables for commands | Flexible metadata exposure without shell parsing | — Pending |
| Multi-camera: OR logic | Most intuitive: "is anyone on camera?" | — Pending |
| Detection method: selectable | User chooses the backend, V4L2 as default on Linux | — Pending |

---
*Last updated: 2026-05-28 after new-project questioning*
