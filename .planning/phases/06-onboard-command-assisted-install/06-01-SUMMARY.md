# Plan 06-01 Summary

**Wave:** 1
**Objective:** Interactive onboard wizard with huh

## What was done

- Added `charmbracelet/huh` dependency
- Created `cmd/onboard.go` with:
  - Camera MultiSelect ("All cameras" + per-device options)
  - Detection method Select (V4L2/lsof with explanation)
  - Live detection test (enable → Detect → show, disable → Detect → show, retry support)
  - Debounce + Interval inputs with validation
  - `--dry-run` flag to preview config as YAML
  - Config saved to `/tmp/on-a-meet-onboard.json`, auto sudo re-exec to apply

## Files changed

- `go.mod`, `go.sum`: Added huh dependency
- `cmd/onboard.go`: Interactive wizard (new file)

## Commit

`4566166`
