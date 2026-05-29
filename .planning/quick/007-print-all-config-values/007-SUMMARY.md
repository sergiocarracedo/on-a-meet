# Quick Task 007 Summary

**Task:** Check all flags loading — print all config values on startup to compare with passed flags

**Completed:** 2026-05-29

## What was done

Expanded the startup config print line from 4 fields (method, interval, debounce, timeout) to all 7 fields including camera, on-command, and off-command. This helps debug which values the CLI is actually using versus what was passed.

## Files changed

- `cmd/detect.go`: Expanded config line to include all 7 config fields

## Commit

`45c7fa6`
