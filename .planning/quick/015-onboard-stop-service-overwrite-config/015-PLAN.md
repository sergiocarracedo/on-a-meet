# Quick Task 015 Plan: Stop existing service before install + ask before config overwrite

## Tasks

### Task 1: Stop existing service before re-installing

**Files:**
- `cmd/install.go`

**Action:**
In `installService()`, before `svc.Install()`:
1. Call `svc.Stop()` — ignore error (service may not be running)
2. Attempt `svc.Install()`. If it fails (already exists), call `svc.Uninstall()` then retry `svc.Install()`

This prevents the "Init already exists" error from `onboard --apply` when the service was already installed.

### Task 2: Ask before overwriting existing config in onboard --apply

**Files:**
- `cmd/onboard.go`

**Action:**
In the `--apply` path (line 87), before writing the config file, check if `/etc/on-a-meet/config.yaml` already exists using `os.Stat`. If it exists, prompt with `huh.NewConfirm("Config already exists. Overwrite?")`. If user says no, keep existing and skip install (return nil).

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Service install no longer fails on re-run (stops + reinstalls)
- Existing config triggers overwrite confirmation dialog
