# Plan 05-02 Summary

**Completed:** 2026-05-29

## What was built
Wired backend factory into detect and list commands, implemented full list command with status table, created goreleaser config, wrote README with install/config/usage docs, added ldflags version vars to main.go.

## Key files
- `cmd/detect.go`: Uses detector.New(cfg.DetectMethod), added --detect flag
- `cmd/list.go`: Full implementation with pterm table (Path/Driver/Card/Bus/Status)
- `.goreleaser.yaml`: linux/darwin + amd64/arm64, CGO_ENABLED=0
- `README.md`: Install, config, usage, template vars, service management docs
- `main.go`: version/commit/date vars for ldflags injection
