# Architecture Research

**Domain:** Linux camera state detection CLI tool
**Researched:** 2026-05-28
**Confidence:** HIGH

## Component Boundaries

### System Overview

```
┌──────────────────────────────────────────────────────────────┐
│                        CLI Layer (cobra)                       │
├──────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐  │
│  │ detect cmd   │  │ list cmd     │  │ install/uninstall  │  │
│  │ (monitor)    │  │ (devices)    │  │ (service mgmt)     │  │
│  └──────┬───────┘  └──────┬───────┘  └─────────┬──────────┘  │
├─────────┴─────────────────┴────────────────────┴────────────┤
│                    Config Layer (viper)                       │
│         YAML file ← CLI flags ← Env vars ← Defaults          │
├──────────────────────────────────────────────────────────────┤
│                    Detection Engine                            │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐  │
│  │ V4L2 Backend │  │ lsof Backend │  │ udev Backend (v2)  │  │
│  └──────┬───────┘  └──────┬───────┘  └────────────────────┘  │
│         └────────┬────────┘                                    │
│                  ▼                                             │
│  ┌─────────────────────────────────────────────────────┐      │
│  │              State Machine                           │      │
│  │  Polls → Compares → Detects change → Fires command   │      │
│  └──────────────────────┬──────────────────────────────┘      │
├─────────────────────────┴────────────────────────────────────┤
│                    Service Layer (kardianos/service)           │
│         systemd unit generation, start/stop/install            │
└──────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| CLI Layer | Parse commands, flags, display output | spf13/cobra commands + pterm output |
| Config Layer | Read YAML, env vars, merge with CLI flags | spf13/viper with cobra binding |
| Detection Engine | Orchestrate polling, backends, state machine | Go interface pattern |
| V4L2 Backend | Query /dev/videoN via ioctl for device status | go4vl or raw syscall; check VIDIOC_QUERYCAP + VIDIOC_STREAMON status |
| lsof Backend | Check which process holds /dev/videoN open | os/exec to run lsof, or parse /proc |
| udev Backend | Monitor udev events for camera plug/unplug (v2) | github.com/gentlemanautomaton/udev or raw netlink |
| State Machine | Track per-camera state, detect transitions, fire hooks | struct with map[device]state + callback |
| Service Layer | Install/uninstall systemd/launchd unit | kardianos/service |
| Command Executor | Spawn user command with template var substitution | os/exec with text/template |
| Output Layer | Colorized terminal output, spinners, tables | pterm |

## Data Flow

### Primary Data Flow

```
Poll timer tick (every N seconds)
    ↓
Detection Engine iterates configured backends
    ↓
Each backend checks camera device X
    ↓
Engine compares result with last known state
    ├── No change → wait for next tick
    └── State changed → run template → exec command
```

### Data Flow Description

| Flow | Source | Destination | Format | Notes |
|------|--------|-------------|--------|-------|
| Device poll | /dev/videoN | Detection Engine | bool (in use or not) | V4L2 ioctl or lsof output |
| State change | Detection Engine | Command Executor | {CameraID, Device, State} | Go struct |
| Command exec | Command Executor | os/exec | templated string | text/template with struct vars |
| Config read | config.yaml | Viper | YAML → Go struct | Read at startup, CLI flags override |
| Output | All components | stdout/stderr | pterm formatted | Colored, structured |

## Build Order

Suggested implementation sequence based on dependencies:

| Order | Component | Dependencies | Rationale |
|-------|-----------|--------------|-----------|
| 1 | Go module + CLI scaffold | None | cobra init, project structure |
| 2 | Config (Viper + YAML) | #1 | All features need config |
| 3 | Device listing (V4L2) | #1 | Foundation: enumerate cameras |
| 4 | V4L2 state detection backend | #1, #2, #3 | Core detection logic |
| 5 | State machine + polling loop | #1, #2, #4 | Wire detection to state transitions |
| 6 | Command execution with templates | #1, #5 | --on/--off flags, template vars |
| 7 | lsof detection backend | #1, #2 | Fallback detection method |
| 8 | Service installation | #1, #2, #5, #6 | --install flag wraps working tool |
| 9 | macOS detection (v2) | #1, #2, #5, #6 | AVFoundation backend |
| 10 | udev detection backend (v2) | #1, #2 | Event-based detection |

## Integration Points

### External Integrations

| Integration | Type | Protocol | Auth | Notes |
|------------|------|----------|------|-------|
| /dev/videoN | Device file | V4L2 ioctl syscalls | udev permissions | /dev/video* requires read access or root |
| lsof | CLI command | stdout parsing | None | Standard Linux tool |
| systemd | Service manager | unit file + dbus | systemctl | /etc/systemd/system/ requires root |
| launchd | Service manager | plist file | launchctl | ~/Library/LaunchAgents/ for user agents |

### Internal Boundaries

| Boundary | Left Side | Right Side | Contract |
|----------|-----------|------------|----------|
| Detection backend | Detection Engine | Backend implementation | Backend interface: Detect(devicePath) → (bool, error) |
| State change | Detection Engine | Command Executor | StateChange{CameraID, Device, State string} |
| Config | CLI | Config Layer | viper.Viper instance bound to cobra flags |

## Recommended Project Structure

```
on-a-meet/
├── cmd/
│   └── root.go              # Cobra root command + config init
│   └── detect.go            # 'detect' subcommand (main monitoring)
│   └── list.go              # 'list' subcommand (show cameras)
│   └── install.go           # 'install'/'uninstall' subcommand
├── internal/
│   ├── detector/
│   │   ├── interface.go     # Detector interface
│   │   ├── v4l2.go          # V4L2 backend
│   │   ├── lsof.go          # lsof backend
│   │   └── udev.go          # udev backend (v2)
│   ├── engine/
│   │   ├── engine.go        # Polling loop, state machine
│   │   └── engine_test.go
│   ├── exec/
│   │   ├── command.go       # Command execution with templates
│   │   └── command_test.go
│   ├── config/
│   │   └── config.go        # Viper setup, struct types
│   └── service/
│       └── service.go       # kardianos/service wrapper
├── main.go                  # Entry point
├── config.yaml.example      # Example config file
└── go.mod
```

---
*Architecture research for: on-a-meet camera detection CLI*
*Researched: 2026-05-28*
