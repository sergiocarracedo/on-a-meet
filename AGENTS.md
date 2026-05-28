# on-a-meet — Agent Guide

## Current Phase

**Milestone:** v1 — Initial release
**Phase:** 3 — Command Execution & Templates ✅ complete → Phase 4 — Service Installation
**Status:** verifying
**Last updated:** 2026-05-28

## Project Summary

CLI tool in Go that detects camera on/off state and triggers user-defined commands. Monitors /dev/video* devices via V4L2 (or lsof as fallback), polls at configurable intervals, and fires commands with template variables on state transitions.

## Tech Stack

- **Language:** Go 1.22+
- **CLI framework:** spf13/cobra v1.10.2
- **Config:** spf13/viper v1.19+ (YAML, CLI flags, env vars)
- **Output:** pterm/pterm v0.12.83
- **Service mgmt:** kardianos/service v1.2.4

## Architecture

- Pluggable detection backends behind `Detector` interface: `Detect(devicePath) (DeviceStatus, error)` and `ListDevices() ([]DeviceInfo, error)`
- Polling engine with state machine, debounce (N consecutive same-state polls)
- Command execution via goroutines with text/template substitution
- Config merged via Viper: CLI flags > env vars > YAML config > defaults

## Key Files

- `.planning/PROJECT.md` — Project scope and context
- `.planning/REQUIREMENTS.md` — REQ-001 through REQ-017
- `.planning/ROADMAP.md` — 5 phases
- `.planning/research/` — Stack, Features, Architecture, Pitfalls, Summary
- `.planning/STATE.md` — Project state and session tracking
- `.planning/phases/01-project-scaffold-cli-foundation/` — Phase 1 context, plans, summaries

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
│   ├── detect.go         # detect subcommand — V4L2 polling + command execution
│   ├── list.go           # list subcommand (stub)
│   ├── install.go        # install subcommand (stub)
│   └── uninstall.go      # uninstall subcommand (stub)
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

# Run
./on-a-meet detect --interval 500ms
./on-a-meet detect --on "echo {{.State}}" --off "echo {{.State}}"
./on-a-meet list
./on-a-meet install
./on-a-meet uninstall
```
