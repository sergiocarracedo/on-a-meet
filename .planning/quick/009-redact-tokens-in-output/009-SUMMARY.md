# Quick Task 009 Summary

**Task:** Detect tokens (strings starting with `ey...`) in output and replace with asterisks

**Completed:** 2026-05-29

## What was done

Added `output.RedactSecrets()` function that detects JWT-like tokens (3-part base64url strings starting with `ey`) and replaces them with `ey***.***.***`. Applied it to:
- Config line: on/off commands are redacted before display
- Executor debug output: rendered command and output are redacted before display

## Files changed

- `internal/output/output.go`: Added `jwtRe` regex and `RedactSecrets()` function
- `cmd/detect.go`: Wrapped `cfg.OnCmd` and `cfg.OffCmd` with `output.RedactSecrets()`
- `internal/executor/executor.go`: Wrapped `rendered` and `outStr` with `output.RedactSecrets()`

## Commit

`cdf32b8`
