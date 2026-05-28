# Research Summary

**Domain:** Linux camera state detection CLI tool
**Researched:** 2026-05-28
**Confidence:** HIGH

## Recommended Stack

### Primary Technologies

| Technology | Version | Role |
|------------|---------|------|
| Go | 1.22+ | Runtime — single binary, cross-compile, direct syscall access |
| spf13/cobra | v1.10.2 | CLI framework — commands, flags, help, completion |
| spf13/viper | v1.19+ | Config — YAML files + CLI flag override + env vars |
| pterm/pterm | v0.12.83 | Terminal output — colored status, spinners, tables |
| kardianos/service | v1.2.4 | Service management — systemd/launchd install/uninstall |

### Key Stack Decisions

- **Cobra over urfave/cli:** Need nested subcommands (detect, list, install). Cobra's command tree is better suited.
- **kardianos/service over custom templates:** Handles cross-platform service paths, root detection, user-mode services. Avoids maintaining systemd/launchd template logic ourselves.
- **Raw ioctl over go4vl:** For state detection we only need to query device status, not stream video. go4vl adds dependency weight for minimal benefit. Use Go's syscall package directly.
- **PTerm over lipgloss/bubbletea:** Lighter weight. We only need formatted output (spinners, tables, colors), not a full TUI framework.

## Table Stakes Features

Features that must be in v1:

- [x] V4L2 polling detection for camera on/off state
- [x] --on and --off command execution with Go text/template variables
- [x] Multi-camera OR logic (any camera on = state on)
- [x] --camera flag to target specific device
- [x] --interval flag for polling rate (default 1s)
- [x] YAML config file with CLI flag override (Viper precedence)
- [x] Device list command
- [x] systemd service installation via --install flag
- [ ] lsof backend (P2 — v1.x)
- [ ] macOS support (P2 — v1.x)

## Key Architecture Decisions

### System Shape

Single-binary Go CLI with three subcommands (`detect`, `list`, `install`), a polling engine with pluggable detection backends, and service management via kardianos/service. Config via Viper (YAML + CLI flags). Output via pterm.

### Critical Boundaries

| Boundary | What It Separates | Why It Matters |
|----------|-------------------|----------------|
| Detector interface | Detection Engine ↔ Backend impl | Pluggable backends (V4L2, lsof, udev) behind a common `Detect(device) → (bool, error)` |
| State machine | Detector output ↔ Command exec | Debounce prevents false positives; dedup prevents duplicate command fires |
| Service wrapper | CLI ↔ OS service manager | Clean separation: CLI handles flags, service wrapper handles lifecycle |

### Recommended Build Order

1. **Go module + cobra scaffold** — project structure, empty commands, pterm hello
2. **Config layer** — Viper reads YAML, binds cobra flags, merge precedence
3. **Device listing** — enumerate /dev/video*, print with pterm table
4. **V4L2 detection backend** — open device, check state via VIDIOC_QUERYCAP
5. **State machine + poll loop** — periodic detection, transition tracking, debounce
6. **Command execution** — text/template substitution, goroutine spawn with timeout
7. **lsof backend** — fallback detection when V4L2 gives false positives
8. **Service installation** — kardianos/service wraps tool for systemd/launchd
9. **macOS + udev** (v1.x)

## Top Pitfalls

| # | Pitfall | Severity | Prevention |
|---|---------|----------|------------|
| 1 | V4L2 false positives — opening device ≠ streaming | CRITICAL | Debounce window; combine V4L2 + lsof; don't treat open-file as "camera on" |
| 2 | Permission errors on /dev/video* | HIGH | Check at startup; print clear fix: "add user to video group" |
| 3 | Camera hotplug crashes | HIGH | Handle ENOENT per-device; re-scan device list periodically |
| 4 | Blocking command execution delays polling | MEDIUM | Spawn commands in goroutines with configurable timeout |
| 5 | StreamON/StreamOFF ambiguity — app opens device but not streaming | MEDIUM | Document as known limitation; recommend lsof backend for most accurate results |

## Primary Recommendation

**Build the polling engine with debounce and pluggable backends from day one.** The detection backend interface (`Detect(devicePath string) (inUse bool, err error)`) is the single most important architectural decision. Implement V4L2 first (via Go syscall, not go4vl), add lsof as the reliability fallback. Debounce (3 consecutive same-state polls before firing) is essential to avoid false positives from transient device queries.

---
*Research summary for: on-a-meet camera detection CLI*
*Researched: 2026-05-28*
*Sources: STACK.md, FEATURES.md, ARCHITECTURE.md, PITFALLS.md*
