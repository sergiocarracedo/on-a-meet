# Phase 6 Verification

## Goal
Interactive `onboard` command that walks users through camera selection, detection method selection with live test, configuration, and automatic service installation.

## Must-Haves Check
- [x] Camera selection via MultiSelect (huh) — works with "All cameras" option
- [x] Detection method selection (V4L2/lsof) with explanation
- [x] Live detection test (enable + disable with detection verification)
- [x] Debounce + interval configuration inputs
- [x] `--dry-run` flag to preview config
- [x] Config saved to `/tmp/on-a-meet-onboard.json`
- [x] Auto sudo re-exec to apply config and install service
- [x] `--apply <file>` flag for sudo path
- [x] YAML config written to `/etc/on-a-meet/config.yaml`
- [x] Reuses `installService()` from install command

## Build & Tests
- `go build ./...` — ✅ passes
- `go vet ./...` — ✅ passes
- `go test ./...` — ✅ 19 tests passing

## Dry-run verification
```
./on-a-meet onboard --dry-run
```
Runs the wizard and prints config to stdout without sudo.

**Status:** ✅ PASSED (2026-05-29)
