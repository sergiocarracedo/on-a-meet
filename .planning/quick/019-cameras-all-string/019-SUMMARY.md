# Quick Task 019 Summary

**Task:** Accept "all" string for cameras in --config JSON
**Completed:** 2026-05-29

## What was done

`onboard --config` now accepts `"cameras": "all"` (string) to mean monitor all cameras, in addition to the existing array format `["/dev/video0"]`.

## Files changed

- `cmd/onboard.go`: Added `cameraList` type with `UnmarshalJSON` accepting both string and array; `"all"` → skip camera field in YAML output

## Commit

`0808bce`
