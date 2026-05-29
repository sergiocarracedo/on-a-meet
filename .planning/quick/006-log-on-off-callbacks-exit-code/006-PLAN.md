# Quick Task 006 Plan: Log on/off callbacks with exit code and verbose output

## Tasks

### Task 1: Add Debug to output package

**Files:**
- `internal/output/output.go`

**Action:**
Add `Debug = pterm.Debug` alongside the existing log level aliases. `pterm.Debug` messages are only rendered when verbose mode is enabled (via `pterm.EnableDebugMessages()` in `Init()`).

### Task 2: Log command completion in executor

**Files:**
- `internal/executor/executor.go`

**Action:**
Import `github.com/sergiocarracedo/on-a-meet/internal/output`. After `cmd.Run()` completes (whether success or error), log:
- Always (using `output.Info`): `"on-command"` / `"off-command"` + `"exited with code X"`: right before returning
- Verbose (using `output.Debug`): the rendered command string + the captured output

The existing error return path should remain unchanged. The new log lines should use the same theme as the rest of the CLI output.

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All existing tests pass

**Done:**
- `output.go` has `Debug` variable
- Executor logs exit code (always) and command+output (verbose)
