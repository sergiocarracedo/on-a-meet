# Quick Task 008 Summary

**Task:** Verbose is not being detected, and also we are not displying when a callback is executed

**Completed:** 2026-05-29

## What was done

1. Added `silent` and `verbose` values to the startup config print line so users can verify their `--verbose/-V` and `--silent/-s` flags are being picked up
2. Added `output.Info.Printfln("%s-command executing", state)` at the start of the executor's `exec()` method so there's immediate feedback when a callback command is triggered (before waiting for output)

## Files changed

- `cmd/detect.go`: Added `silent=%t verbose=%t` to config line using `cfgSilent, cfgVerbose`
- `internal/executor/executor.go`: Added `{state}-command executing` log line before command runs

## Commit

`a30790c`
