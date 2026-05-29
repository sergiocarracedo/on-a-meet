# Quick Task 015 Summary

**Task:** If the service is running stop it before configure to avoid this error, and if the config already exists, ask the user if we want to overwrite the config

**Completed:** 2026-05-29

## What was done

1. **installService()** (`cmd/install.go`): Added `svc.Stop()` before install (ignore error — service might not be running). If `svc.Install()` fails (e.g., "Init already exists"), call `svc.Uninstall()` and retry. This prevents re-install errors from onboard --apply.
2. **onboard --apply** (`cmd/onboard.go`): Added `os.Stat()` check before writing config. If `/etc/on-a-meet/config.yaml` exists, prompts with "Config already exists. Overwrite?" dialog. User can abort to keep existing config.

## Files changed

- `cmd/install.go`: Added stop-before-install + uninstall-retry on failure
- `cmd/onboard.go`: Added config overwrite confirmation prompt

## Commit

`ee1987b`
