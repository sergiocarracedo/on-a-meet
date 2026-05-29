# Quick Task 019: Accept "all" string for cameras in --config JSON

## Task

Allow `"cameras": "all"` in the JSON config file to mean "monitor all cameras".

### Files
- `cmd/onboard.go`

### Action
1. Add `cameraList` type with `UnmarshalJSON` that accepts both `"all"` (string) and `["/dev/video0"]` (array)
2. Change `onboardConfig.Cameras` type from `[]string` to `cameraList`
3. In `--apply` path: `"all"` entry → camera = "" (monitor all)
4. Cast `cameraList(cameras)` in interactive path

### Done
When `{"cameras": "all"}` shows `cameras: all` in dry-run and produces a YAML with no camera field.
