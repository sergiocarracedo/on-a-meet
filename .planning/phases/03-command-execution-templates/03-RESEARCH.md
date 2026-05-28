# Phase 3 Research: Command Execution & Templates

**Gathered:** 2026-05-28

## Don't Hand-Roll

- **`os/exec` is sufficient** [VERIFIED: pkg.go.dev/os/exec]. No need for external libs like `go-cmd/cmd`. The standard library provides `CommandContext`, context cancellation, and pipe management.
- **`text/template` is sufficient** [VERIFIED: pkg.go.dev/text/template]. Template rendering with `Execute` and `bytes.Buffer` is straightforward. No external template engine needed.
- **Go 1.20+ features:** `cmd.Cancel` (defaults to `os.Process.Kill`) and `cmd.WaitDelay` (grace period after cancel before force-kill) handle timeout cleanup reliably [CITED: go.dev/go1.20].

## Common Pitfalls

1. **`exec.CommandContext` + `CombinedOutput()` can hang on long-running commands** [CITED: jarv.org]. When the context is cancelled, `SIGKILL` is sent to the shell process (`sh`), but `CombinedOutput` waits for the pipe to close — which doesn't happen until the child process (e.g., `sleep 10`) finishes. **Fix:** Use `exec.CommandContext` with `Run()` instead of `CombinedOutput()`, or use the `cmd.Cancel`/`cmd.WaitDelay` fields (Go 1.20+) and capture output manually.

2. **Setting `cmd.Stdout`/`cmd.Stderr` to buffers can break context cancellation** [CITED: StackOverflow]. When buffers are set, the context kill signal goes to `sh` but doesn't propagate to the spawned child. **Fix:** Use `cmd.Cancel` to kill the process group (`syscall.Kill(-cmd.Process.Pid, ...)`) rather than just the shell.

3. **Shell commands orphan child processes** [CITED: jarv.org]. When `sh -c` is killed, the child process becomes an orphan adopted by PID 1 — it continues running. **Fix:** Set `SysProcAttr{Setpgid: true}` and kill the process group (`syscall.Kill(-pid, syscall.SIGKILL)`).

4. **Stderr output may be lost** if only stdout is captured. **Fix:** Use `CombinedOutput` or capture both separately with `MultiWriter`.

## Existing Patterns in This Codebase

- **Functional options pattern** — used in engine (`WithInterval`, `WithDebounce`, etc.) and should be used in executor too
- **`internal/` packages** — executor follows the same convention
- **pterm output wrappers** — `output.Info`, `output.Warning`, `output.Error` for terminal output, all respect `--silent`/`--verbose`
- **Engine OnChange callback** — already provides `path string, oldState, newState bool, info detector.DeviceInfo`
- **Config struct** at `internal/config/config.go` — already has `OnCommand` and `OffCommand` fields; needs `Timeout` added

## Recommended Approach

### Executor Design

```go
type Executor struct {
    timeout     time.Duration
    running     sync.Map  // key: "on"/"off", value: context.CancelFunc
}

func New(timeout time.Duration) *Executor

// Exec executes a command with template substitution.
// Returns when the command completes or the context is cancelled.
func (e *Executor) Exec(ctx context.Context, command string, data TemplateData) error
```

### Template Data

```go
type TemplateData struct {
    CameraID string  // "video0"
    Device   string  // "/dev/video0"
    State    string  // "on" | "off"
}
```

### Command Execution (shell mode)

```go
cmd := exec.CommandContext(ctx, "sh", "-c", cmdStr)
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

var buf bytes.Buffer
cmd.Stdout = &buf
cmd.Stderr = &buf

// Go 1.20+: context cancellation kills process group
cmd.Cancel = func() error {
    return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
cmd.WaitDelay = 5 * time.Second

err := cmd.Run()
output := buf.String()
```

### Overlap Prevention

Track running commands per state using `sync.Map`:
- Before starting an ON command, check if an ON is already running → skip if so
- Same for OFF
- Store `context.CancelFunc` in the map, remove it when done
- Cross-state (ON while OFF running) always proceeds

### Integration

- `cmd/detect.go` creates `executor.Executor` and passes it to the engine's `OnChange` callback
- The OnChange callback checks `cfg.OnCommand`/`cfg.OffCommand`, renders the template, and calls `executor.Exec()`
- New `--timeout` flag added to detect command
- Config struct gets `Timeout` field added
- `config.yaml.example` gets `timeout: "30s"`

