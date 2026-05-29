# Quick Task 006 Summary

**Task:** Log when on/off callbacks happen with exit code (always) and command string + output (verbose mode)

**Completed:** 2026-05-29

## What was done

Added `Debug` log level to `output` package (maps to `pterm.Debug`, only shown in verbose mode). Modified `executor.exec()` to log the rendered command string and output at debug level, and the exit code at info level for both success and failure paths.

## Files changed

- `internal/output/output.go`: Added `Debug = pterm.Debug`
- `internal/executor/executor.go`: Added logging for command string, output (debug/verbose) and exit code (info/always)

## Commit

`8ca2245`
