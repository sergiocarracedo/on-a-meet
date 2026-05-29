# Quick Task 024 Summary

**Task:** Fire on/off commands when service or CLI starts
**Completed:** 2026-05-29

## What was done

Modified the engine to fire the OnChange callback (which triggers --on / --off commands) for each camera device on startup, based on its current state. Previously the engine silently initialized states without firing commands, so a camera already ON would never trigger the on-command until it cycled OFF→ON.

## Files changed

- `internal/engine/engine.go` — after initial state detection, fire `onChange` with `oldState = !actual` so the callback dispatches the correct on/off command
- `internal/engine/engine_test.go` — updated `TestEngineStartup` to expect 2 initial fires (one per device), updated `TestDebounce` to check 2nd call (initial fire + debounced transition)

## Commit

`14a10f0`
