# Quick Task 005 Summary

**Task:** The CLI should emit the initial cameras status on start

**Completed:** 2026-05-29

## What was done

Added a loop in `cmd/detect.go` that calls `det.Detect()` for each enumerated device after the startup banner, printing initial ON/OFF status per device. Uses the same `output.Info.Printfln` pattern and indentation.

## Files changed

- `cmd/detect.go`: Added initial status print block after device banner line

## Commit

`d33a9df`
