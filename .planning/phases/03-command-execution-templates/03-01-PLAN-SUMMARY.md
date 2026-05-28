# Plan 03-01 Summary

**Completed:** 2026-05-28

## What was built

Executor package (`internal/executor/`) with template rendering via `text/template`, shell command spawning via `sh -c`, configurable timeout with process group kill, and same-state overlap prevention. Added `Timeout` field to Config struct (default 30s).

## Key files

- `internal/executor/executor.go`: Executor struct, TemplateData, ExecOn/ExecOff, template rendering, timeout via context.WithTimeout, process group management (Setpgid + syscall.Kill(-pid)), WaitDelay
- `internal/executor/executor_test.go`: 6 tests — success, template substitution, timeout (50ms), same-state skip, cross-state allow
- `internal/config/config.go`: Added Timeout field with `"30s"` default

## Decisions made

- `sync.Map` tracks running commands per state ("on"/"off") for same-state skip
- `context.Background()` passed from detect command so in-flight commands survive engine shutdown
- Process groups used so `syscall.Kill(-pid, SIGKILL)` kills child processes on timeout
- `cmd.Cancel` + `WaitDelay` (5s) for Go 1.20+ clean process termination
- Combined stdout/stderr captured for error reporting

## Notes for downstream

- Commands fire in goroutines — errors are reported via pterm Warning but never block polling
- Hotplug events (oldState==newState) are ignored in OnChange — no commands fire on add/remove
