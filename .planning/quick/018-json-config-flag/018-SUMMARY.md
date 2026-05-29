# Quick Task 018 Summary

**Task:** JSON config flag for onboard
**Completed:** 2026-05-29

## What was done

Added `--config <file>` flag to onboard command. Reads all config values from a JSON file (cameras, method, interval, debounce, on-command, off-command), skips the interactive wizard entirely, and goes straight to the install confirmation. Missing fields default to v4l2/1s/2.

## Files changed

- `cmd/onboard.go`: Added `onboardConfigFile` variable, `--config` flag, JSON parsing with defaults, shared confirm+sudo re-exec flow

## Usage

```bash
onboard --config myconfig.json            # read config, skip wizard, confirm, install
onboard --config myconfig.json --dry-run  # preview only
```

## Commit

`<commit-sha>`
