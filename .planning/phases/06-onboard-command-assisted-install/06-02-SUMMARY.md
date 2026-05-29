# Plan 06-02 Summary

**Wave:** 2
**Objective:** Sudo apply path — read collected answers, write config to /etc, install service

## What was done

- Extracted `installService()` from `cmd/install.go` for reuse by both `install` and `onboard` commands
- Added `--apply <file>` flag to `onboard` command for the sudo execution path
- `--apply` flow: reads JSON → marshals to YAML → writes `/etc/on-a-meet/config.yaml` → calls `installService()`
- Added `gopkg.in/yaml.v3` dependency for YAML marshaling
- Non-sudo flow now auto re-execs with `sudo on-a-meet onboard --apply /tmp/on-a-meet-onboard.json`
- Config handles single-camera (sets `camera:` field) and multi-camera (omits field = monitor all)

## Files changed

- `go.mod`, `go.sum`: Added yaml.v3 dependency
- `cmd/install.go`: Extracted `installService()` helper
- `cmd/onboard.go`: Added `--apply` flag, auto sudo re-exec, YAML write path

## Commit

`a106d4e`
