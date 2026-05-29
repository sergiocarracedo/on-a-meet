# Quick Task 022 Summary

**Task:** Add service subcommand group with start/stop/restart + verbose config
**Completed:** 2026-05-29

## What was done

1. Grouped install/uninstall/restart under a new `service` parent command (`on-a-meet service install`, `service uninstall`, `service restart`)
2. Added `start` and `stop` service subcommands
3. Added `Verbose bool` to the Config struct so `verbose: true` in config.yaml is respected by detect
4. Fixed `output.Init()` ordering in `initConfig()` — now reads viper (merged config+flags) instead of raw flag vars, so config file `verbose`/`silent` values take effect

## Files changed

- `cmd/service.go` — new: `service` parent command (alias: `svc`)
- `cmd/start.go` — new: `start` subcommand
- `cmd/stop.go` — new: `stop` subcommand
- `cmd/install.go` — register under `serviceCmd` instead of `rootCmd`
- `cmd/uninstall.go` — register under `serviceCmd`; updated sudo hint
- `cmd/restart.go` — register under `serviceCmd`; updated sudo hint
- `cmd/root.go` — moved `output.Init()` after `viper.ReadInConfig()` so config file values merge with flags
- `internal/config/config.go` — added `Verbose bool` field
- `README.md` — updated service management commands
- `AGENTS.md` — updated commands table

## Commit

`e81f640`
