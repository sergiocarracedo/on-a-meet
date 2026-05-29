# on-a-meet

## What This Is

A CLI tool that monitors camera state (on/off) on Linux and macOS, and triggers user-defined commands when state changes. Users wire camera activity to home automation, notifications, or any custom workflow — with no infrastructure beyond a single binary.

## Core Value

Reliably detect camera state changes and fire the correct command every time. When the camera turns on, the on-command runs. When it turns off, the off-command runs. Fast enough to feel immediate.

## Requirements

### Validated ✅ (v1.0.0)

- Camera detection via multiple backends: V4L2 and lsof — user-selectable via `--detect` flag or YAML config, V4L2 default on Linux. lsof with ENOENT fallback. REQ-001, REQ-013.
- Command execution on state change with template variables (`{{.CameraID}}`, `{{.Device}}`, `{{.State}}`). Timeout, overlap prevention, env file support, JWT redaction. REQ-002, REQ-003.
- Multi-camera support: OR logic (any camera on = on), `--camera` flag to target specific device. REQ-004, REQ-005.
- Configurable polling via `--interval`, debounce window (N consecutive same-state). REQ-006, REQ-010.
- YAML config file via Viper — CLI flags, env vars, config file with correct precedence. Service config at `/etc/on-a-meet/config.yaml`. REQ-007.
- Device list command with pterm table (Path, Driver, Card, Bus, Status). REQ-008.
- Service installation: systemd (Linux) / launchd (macOS) via kardianos/service. sudo check, one-shot API, runs as original user. Environment file support. REQ-009.
- Permission check at startup with clear error and fix instructions. REQ-011.
- Graceful camera hotplug — ENOENT handled per-device, periodic re-scan, log add/remove. REQ-012.
- Interactive onboard wizard (huh) — camera selection, method select with live test, config input, dry-run, auto sudo re-exec. All quick tasks 001-025.
- Release automation: GoReleaser cross-platform builds (linux/darwin + amd64/arm64), GitHub Actions workflow, GitHub Releases with checksums.

## Current Milestone: v1.1.0 — macOS Support & Docs Polish

**Goal:** Add macOS camera detection backend and fix README gaps (onboard docs).

**Target features:**
- macOS camera detection backend (AVFoundation or IOKit)
- README documentation for the `onboard` command
- Cross-platform test coverage

### Active

- macOS camera detection
- README onboard docs

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

- **Language**: Go 1.22+ — binary distribution, single-file output, cross-compilation
- **Output**: pterm — consistent terminal output formatting
- **OS**: Linux first (primary), macOS second
- **Service**: systemd on Linux, launchd on macOS — auto-generate unit files via kardianos/service
- **Version**: v1.1.0 — macOS Support & Docs Polish (in progress)

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Polling (not event-driven) | Simplest implementation, fewer edge cases | ✓ Implemented — configurable interval, debounce window |
| Config format: YAML | Familiar, readable, common in Go ecosystem | ✓ Viper-backed, CLI/env override precedence |
| Template variables for commands | Flexible metadata exposure without shell parsing | ✓ text/template — {{.CameraID}}, {{.Device}}, {{.State}} |
| Multi-camera: OR logic | Most intuitive: "is anyone on camera?" | ✓ Default behavior, --camera filter available |
| Detection method: selectable | User chooses the backend, V4L2 as default on Linux | ✓ New(method) factory — "v4l2" or "lsof" |
| Service mgmt: kardianos/service | Cross-platform, handles systemd + launchd | ✓ Install/Uninstall/Start/Stop/Restart, runs as original user |
| Command timeout and overlap prevention | Prevent hung commands blocking transitions | ✓ context.WithTimeout, sync.Map same-state skip, Setpgid process kill |

---
*Last updated: 2026-05-29 — v1.1.0 started*
