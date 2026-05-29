# Quick Task 007 Plan: Print all config values on startup

## Tasks

### Task 1: Expand config print line to show all fields

**Files:**
- `cmd/detect.go`

**Action:**
Replace the existing config print line at line 65:
```go
output.Info.Printfln("Config: method=%s interval=%s debounce=%d timeout=%s", cfg.DetectMethod, cfg.Interval, cfg.Debounce, cfg.Timeout)
```
with a comprehensive one that includes every config field:
```go
output.Info.Printfln("Config: method=%s interval=%s debounce=%d timeout=%s camera=%s on=%s off=%s", cfg.DetectMethod, cfg.Interval, cfg.Debounce, cfg.Timeout, cfg.Camera, cfg.OnCmd, cfg.OffCmd)
```

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Config line prints all 7 config values for CLI vs loaded comparison
