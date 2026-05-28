# Phase 3: Command Execution & Templates — Context

**Gathered:** 2026-05-28
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Execute user-defined commands on camera state transitions with template variable substitution. Includes command timeout, overlap prevention for same-state commands, and proper error reporting. Phase 4 handles service installation — this phase delivers the command execution layer that the engine's OnChange callback fires into.

</domain>

<decisions>
## Implementation Decisions

### Command Executor Package
- New `internal/executor/` package with an `Executor` struct.
- `detect.go` creates the executor, the engine's `OnChange` callback invokes it.
- The executor is completely separate from the engine — clean separation of concerns, testable in isolation.
- Executor methods: `ExecOn(path string, info detector.DeviceInfo)` and `ExecOff(path string, info detector.DeviceInfo)`.

### Shell vs No Shell
- Use `sh -c` for command execution via `os/exec`.
- Users expect pipes (`|`), redirects (`>`), env vars (`$VAR`), and chaining (`&&`).
- **Note:** Shell injection is accepted — the user provides these commands intentionally. Document that commands run through the shell.

### Template Variables
- Four fields exposed:
  - `{{.CameraID}}` — short device name (e.g., `video0`)
  - `{{.Device}}` — full device path (e.g., `/dev/video0`)
  - `{{.State}}` — `on` or `off`
- Struct name: `TemplateData` in the executor package or common types file.
- Template parsing: use `text/template` with `template.Must` and cached templates.

### Overlap Prevention Policy
- **Same-state skip:** If an ON command is already running when another ON fires, skip the new one. Same for OFF.
- **Cross-state allow:** ON → OFF transition always starts the off-command, even if the on-command is still running.
- **Fire-and-forget model:** No queueing. If we skip a same-state command, we log a debug message.
- **State tracking:** Per-command running state tracked via a `sync.Map` or similar.

### Timeout & Cancellation
- Default timeout: 30 seconds.
- Configurable via `--timeout` flag and `Config.Timeout` field (overrides default).
- If set to `0`, no timeout (run until completion).
- On timeout: kill the process group, log a warning with the command and duration.

### Error Handling
- Capture both stdout and stderr from the command.
- On non-zero exit or timeout: print a warning via pterm with the exit code and captured stderr.
- Continue polling — command errors never block the detection loop.
- Logging respects `--silent`/`--verbose` flags (pterm handles this via output wrappers).

### Agent's Discretion
- Internal details of the executor struct (method signatures, field names) — follow existing conventions.
- Test approach: use a test script or mock exec — choose based on what's cleanest.
- Whether to add `--on`/`--off` as part of the executor or keep them in detect.go config — the executor receives the command strings from config.

</decisions>

<specifics>
## Specific Ideas

- The executor should be usable standalone (not just via engine) for future flexibility.
- Config struct needs a `Timeout` field added (duration string like `"30s"`).

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — REQ-002, REQ-003
- `.planning/phases/02-camera-detection-engine/02-CONTEXT.md` — Engine OnChange callback signature, DeviceInfo type
- `internal/config/config.go` — Config struct (needs Timeout field added)
- `cmd/detect.go` — detect command RunE, existing OnChange wiring, flag definitions
- `internal/engine/engine.go` — OnChange type, Engine struct, pollCycle

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/output/` — pterm wrappers for warning/error output, respects --silent/--verbose
- `internal/engine/engine.go` — OnChange callback with path, oldState, newState, DeviceInfo
- `internal/detector/interface.go` — DeviceInfo type with Path, Driver, Card, Bus

### Established Patterns
- `internal/` subpackages with single responsibility (detector, engine, config, output, executor)
- Option function pattern for engine configuration
- Test files alongside source code

### Integration Points
- `internal/executor/` — new package
- `cmd/detect.go` — inject executor into OnChange callback, parse --timeout flag
- `internal/config/config.go` — add Timeout field
- `config.yaml.example` — add timeout: "30s"
- `internal/engine/engine.go` — OnChange callback already supports firing commands (currently just logs)

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 03-command-execution-templates*
*Context gathered: 2026-05-28*
