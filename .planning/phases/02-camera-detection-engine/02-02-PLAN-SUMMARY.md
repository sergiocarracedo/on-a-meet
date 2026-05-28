# Plan 02-02 Summary

**Completed:** 2026-05-28

## What was built

Polling engine that drives the V4L2 detector with configurable interval, debounce (N consecutive same-state polls), hotplug detection (re-enumerate each cycle), `--camera` filtering, and graceful shutdown via context cancellation. The detect command now shows a startup banner with detected cameras and prints real-time state change lines.

## Key files

- `internal/engine/engine.go`: `Engine` struct with `Run(ctx)`, `pollCycle()`, functional options (`WithInterval`, `WithDebounce`, `WithCameraFilter`, `WithOnChange`), per-device state tracking with debounce counter
- `internal/engine/engine_test.go`: Mock detector + 5 tests (startup, debounce, hotplug add, camera filter, graceful shutdown)
- `cmd/detect.go`: Full implementation — V4L2Detector creation, startup banner, signal-aware context, hotplug-aware OnChange handler

## Decisions made

- Hotplug add uses `onChange(path, false, false, info)` — oldState==newState signals hotplug event
- Hotplug remove uses `onChange(path, true, true, info)` — oldState==newState=true for disambiguation from ON→OFF
- `pollCycle()` does all map mutations under lock but dispatches callbacks outside lock
- Removed devices are detected before added devices in the cycle (so removal can't be mistaken for add)

## Notes for downstream

- The `OnChange` callback distinguishes events: `oldState != newState` = state transition, `oldState == newState` = hotplug
- Phase 3 will slot in command execution between state change detection and terminal output
