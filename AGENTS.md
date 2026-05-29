# on-a-meet — Agent Guide

## Current Phase

**Milestone:** v1 — Initial release
**Phase:** 6 — Onboard Command — Assisted Install
**Status:** ✅ complete
**Last updated:** 2026-05-29

## Project Summary

CLI tool in Go that detects camera on/off state and triggers user-defined commands. Monitors /dev/video* devices via V4L2 (or lsof as fallback), polls at configurable intervals, and fires commands with template variables on state transitions.

## Tech Stack

- **Language:** Go 1.22+
- **CLI framework:** spf13/cobra v1.10.2
- **Config:** spf13/viper v1.19+ (YAML, CLI flags, env vars)
- **Output:** pterm/pterm v0.12.83
- **Service mgmt:** kardianos/service v1.2.4
- **Interactive UI:** charmbracelet/huh v1.0.0

## Architecture

- Pluggable detection backends behind `Detector` interface: `Detect(devicePath) (DeviceStatus, error)` and `ListDevices() ([]DeviceInfo, error)`
- Polling engine with state machine, debounce (N consecutive same-state polls)
- Command execution via goroutines with text/template substitution
- Config merged via Viper: CLI flags > env vars > YAML config > defaults

## Key Files

- `.planning/PROJECT.md` — Project scope and context
- `.planning/REQUIREMENTS.md` — REQ-001 through REQ-017
- `.planning/ROADMAP.md` — 6 phases
- `.planning/research/` — Stack, Features, Architecture, Pitfalls, Summary
- `.planning/STATE.md` — Project state and session tracking
- `.planning/phases/06-onboard-command-assisted-install/` — Phase 6 context, plans, summaries

## Phase 6 ✅ Complete

| Plan | Wave | Depends | Objective | Key Files |
|------|------|---------|-----------|-----------|
| 06-01 | 1 | — | Interactive huh-based wizard — camera MultiSelect, method select+live test, config, dry-run | `cmd/onboard.go` |
| 06-02 | 2 | 06-01 | Sudo apply path — JSON→YAML, write /etc config, install service, auto sudo re-exec | `cmd/onboard.go`, `cmd/install.go` |

## Phase 5 ✅ Complete

| Plan | Wave | Depends | Objective | Key Files |
|------|------|---------|-----------|-----------|
| 05-01 | 1 | — | lsof backend + detector.New factory + tests | `internal/detector/lsof_linux.go`, `internal/detector/detector.go` |
| 05-02 | 2 | 05-01 | Wire factory into detect + list command + goreleaser + README | `cmd/detect.go`, `cmd/list.go`, `.goreleaser.yaml`, `README.md` |

## Phase 4 ✅ Complete

| Plan | Wave | Depends | Objective | Key Files |
|------|------|---------|-----------|-----------|
| 04-01 | 1 | — | Install command: sudo check, Install()+Start() via kardianos/service | `cmd/install.go` |
| 04-02 | 2 | 04-01 | Uninstall command: sudo check, Stop()+Uninstall() | `cmd/uninstall.go` |

## Phase 3 ✅ Complete

| Plan | Wave | Objective | Key Files |
|------|------|-----------|-----------|
| 03-01 | 1 | Executor package + config + tests | `internal/executor/executor.go`, `internal/config/config.go` |
| 03-02 | 2 | Wire executor into detect command | `cmd/detect.go`, `config.yaml.example` |

## Project Structure

```
├── main.go               # Entry point
├── cmd/
│   ├── root.go           # Root command, Viper config, flags
│   ├── onboard.go        # Interactive setup wizard (huh)
│   ├── detect.go         # detect subcommand — V4L2 polling + command execution
│   ├── list.go           # list subcommand (stub)
│   ├── install.go        # install subcommand — kardianos/service Install()+Start()
│   └── uninstall.go      # uninstall subcommand — kardianos/service Stop()+Uninstall()
├── internal/
│   ├── config/
│   │   ├── config.go     # Config struct & defaults
│   │   └── config_test.go
│   ├── detector/
│   │   ├── interface.go  # Detector interface, DeviceStatus, DeviceInfo
│   │   ├── v4l2_linux.go # V4L2Detector — syscall-based camera detection
│   │   └── v4l2_stub.go  # Non-Linux stub
│   ├── engine/
│   │   ├── engine.go     # Polling engine with debounce, hotplug, filtering
│   │   └── engine_test.go
│   ├── executor/
│   │   ├── executor.go   # Command execution with templates, timeout, overlap prevention
│   │   └── executor_test.go
│   └── output/
│       ├── output.go     # pterm wrapper functions
│       └── output_test.go
├── config.yaml.example
├── go.mod / go.sum
```

## Commands

```bash
# Build
go build -o on-a-meet .

# Interactive setup
./on-a-meet onboard                  # Full wizard
./on-a-meet onboard --dry-run        # Preview config only

# Run
./on-a-meet detect --interval 500ms
./on-a-meet detect --on "echo {{.State}}" --off "echo {{.State}}"
./on-a-meet list
./on-a-meet service install     # Requires sudo
./on-a-meet service uninstall   # Requires sudo
./on-a-meet service start       # Requires sudo
./on-a-meet service stop        # Requires sudo
./on-a-meet service restart     # Requires sudo
```
