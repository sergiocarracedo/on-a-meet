# Quick Task 001 Summary

**Task:** Fix --detect flag defaulting to empty string
**Completed:** 2026-05-29

## What was done
Changed `--detect` flag default from `""` to `"v4l2"` in both detect and list commands. When the flag was not explicitly passed, Viper returned the flag's default (`""`) instead of the config struct's `Defaults()` value, causing `"unknown detection method"` error.

## Files changed
- `cmd/detect.go`: Changed `--detect` flag default from `""` to `"v4l2"` (line 142)
- `cmd/list.go`: Changed `--detect` flag default from `""` to `"v4l2"` (line 59)

## Commit
`4b83107` fix(quick-001): set --detect flag default to v4l2
