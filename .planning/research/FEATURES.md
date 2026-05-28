# Feature Research

**Domain:** Linux camera state detection CLI tool
**Researched:** 2026-05-28
**Confidence:** HIGH

## Table Stakes

Features users expect by default. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Detect camera on/off state | Primary purpose of the tool | MEDIUM | V4L2 ioctl + lsof + udev backends, user-selectable |
| Execute command on state change | Core use case (home automation triggers) | LOW | Template variable interpolation ({{.CameraID}}, {{.Device}}, {{.State}}) |
| CLI flags for all options | Standard CLI behavior | LOW | Cobra handles this |
| YAML config file | Config persistence without repeated flags | LOW | Viper reads YAML, CLI flags override |
| Polling interval configuration | User controls responsiveness vs CPU | LOW | Default 1s, configurable via --interval |
| Multi-camera support (OR logic) | Multiple cameras common on laptops (internal + external) | LOW | Any on = on, all off = off |
| Specific camera targeting | User wants only e.g. /dev/video2 | LOW | --camera flag filters to one device |
| Print device list | User needs to know what cameras exist | LOW | List devices with /dev/videoN paths |
| Service installation | Must run in background autonomously | MEDIUM | kardianos/service handles this. Auto-gen systemd/launchd. |
| Cross-platform (Linux + macOS) | User requirement | MEDIUM | Linux first (V4L2), macOS deferred (AVFoundation) |

## Differentiators

Features that set the product apart. Not required for launch, but valuable.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Template variables in commands | Pass camera ID, device path, state to spawned commands | LOW | Go text/template, very simple to implement |
| Detection method selection | User chooses V4L2 / lsof / udev | MEDIUM | Plugin-like backend interface |
| Config file + CLI flag merge | CLI flags win over config, intuitive precedence | LOW | Viper does this natively |

## Anti-Features

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| eBPF-based event detection | Zero-polling, instant response | Requires root, kernel BTF, C compilation, only Linux 5.10+. Massive complexity for marginal gain. | Polling with 500ms-1s interval is good enough |
| GUI or TUI dashboard | Visual camera status | Scope creep. This is a headless service tool. | pterm output for CLI interaction |
| Webhook instead of command execution | "More modern" | Adds HTTP server complexity, security surface. Users can wrap with curl in their command. | --on/--off flags execute arbitrary commands |
| Automatic camera detection for non-V4L2 cameras (e.g. libcamera on Raspberry Pi) | More coverage | libcamera uses different API entirely, not V4L2 | Document limitation, use lsof backend as fallback |

## Feature Dependencies

```
[Detect camera state]
    └──requires──> [List V4L2 devices]
                       └──requires──> [V4L2 device access]

[Service installation]
    └──requires──> [Detect camera state working correctly]

[Command spawning with templates]
    └──requires──> [Detect camera state (need state change events)]

[Config file support]
    └──enhances──> [CLI flags (config can set all flag values)]
```

### Dependency Notes

- **State detection** is the core primitive everything depends on. Must be reliable first.
- **Service installation** should be one of the last features — you need a working detector first.
- **Template variables** are easy to add after command execution works.

## MVP Definition

### Launch With (v1)

- [ ] V4L2 polling detection for camera on/off state
- [ ] --on and --off command execution with template variables
- [ ] Multi-camera OR logic (any camera on = state on)
- [ ] --camera flag for specific device
- [ ] --interval flag for polling rate
- [ ] YAML config file with CLI override
- [ ] Print device list command
- [ ] systemd service installation (--install)

### Add After Validation (v1.x)

- [ ] lsof detection backend (process-based fallback)
- [ ] macOS support (AVFoundation backend)
- [ ] udev detection backend

### Future Consideration (v2+)

- [ ] eBPF-based event detection
- [ ] Multiple detection methods running simultaneously
- [ ] Logging to file/journald

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| V4L2 polling detection | HIGH | MEDIUM | P1 |
| --on/--off command execution | HIGH | LOW | P1 |
| Multi-camera OR logic | HIGH | LOW | P1 |
| Template variables | MEDIUM | LOW | P1 |
| YAML config file | MEDIUM | LOW | P1 |
| Service installation | HIGH | MEDIUM | P1 |
| lsof detection backend | MEDIUM | LOW | P2 |
| macOS support | MEDIUM | HIGH | P2 |
| udev detection backend | LOW | LOW | P3 |

---
*Feature research for: on-a-meet camera detection CLI*
*Researched: 2026-05-28*
