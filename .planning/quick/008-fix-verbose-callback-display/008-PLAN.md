# Quick Task 008 Plan: Fix verbose detection and callback execution display

## Tasks

### Task 1: Add verbose/silent to config line

**Files:**
- `cmd/detect.go`

**Action:**
Add `cfgVerbose` and `cfgSilent` to the startup config print line so users can verify their --verbose/-V and --silent/-s flags are being picked up.

Change line 65 from:
```
Config: method=%s interval=%s debounce=%d timeout=%s camera=%s on=%s off=%s
```
to:
```
Config: method=%s interval=%s debounce=%d timeout=%s camera=%s on=%s off=%s silent=%t verbose=%t
```
Append `cfgSilent, cfgVerbose` to the args. These are package-level vars in root.go, accessible from detect.go.

### Task 2: Add "executing" message before command runs

**Files:**
- `internal/executor/executor.go`

**Action:**
In the `exec()` method, right after the `running` check passes (line 42) and before the timeout setup, add:
```go
output.Info.Printfln("%s-command executing", state)
```
This gives immediate feedback that a callback was triggered, before waiting for the command's output.

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Config line shows `silent` and `verbose` values
- Executor prints `"on-command executing"` / `"off-command executing"` when a callback starts
