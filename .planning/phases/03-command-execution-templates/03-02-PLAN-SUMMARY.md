# Plan 03-02 Summary

**Completed:** 2026-05-28

## What was built

Wired executor into the detect command. Added `--timeout` flag with viper binding. OnChange callback launches goroutines calling ExecOn/ExecOff with template substitution. Updated config.yaml.example with timeout field and detailed template variable docs.

## Key files

- `cmd/detect.go`: imports executor, creates executor in RunE, wires into OnChange callback via goroutines
- `config.yaml.example`: added `timeout` field, expanded template variable docs

## Decisions made

- Executor is created once per detect run, shared across all camera state changes
- Commands fire in goroutines — never block the poll loop
- on-command/off-command checked for empty string before launching goroutine
- `path[5:]` strips "/dev/" to produce CameraID (e.g., "video0")
- `context.Background()` passed to executor — in-flight commands survive engine shutdown

## Notes for downstream

- Demo: `./on-a-meet detect --on 'echo cam-{{.CameraID}}-{{.State}}' --off 'echo cam-{{.CameraID}}-{{.State}}' --interval 500ms`
- Template variables work in both on and off commands
- Phase 4 (Service) will need to configure --on/--off through the config file
