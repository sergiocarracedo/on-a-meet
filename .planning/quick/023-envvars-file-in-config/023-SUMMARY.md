# Quick Task 023 Summary

**Task:** Allow env vars file config for service environment
**Completed:** 2026-05-29

## What was done

Added `environment-file` config field. When set, `service install` and `service restart` patch the systemd unit's `EnvironmentFile=` line to point to the configured path (instead of the default `/etc/sysconfig/on-a-meet`). Variables in that file become available to `--on`/`--off` commands via `os.ExpandEnv()`.

## Files changed

- `internal/config/config.go` — added `EnvironmentFile string` field
- `cmd/root.go` — added `environment-file` viper default
- `cmd/install.go` — added `patchUnitEnvironmentFile()` helper + wired into `installService()`
- `cmd/restart.go` — wired patching before restart
- `config.yaml.example` — documented the new field

## Commit

Will commit after docs.
