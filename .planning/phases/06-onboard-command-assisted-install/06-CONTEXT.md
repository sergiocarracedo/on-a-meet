# Phase 6: Onboard Command — Assisted Install - Context

**Gathered:** 2026-05-29
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Interactive `onboard` command that walks users through camera selection, detection method selection with live test, debounce/interval configuration, config file generation, and automatic service installation.

</domain>

<decisions>
## Implementation Decisions

### Interactive UI Library
- Use `github.com/charmbracelet/huh` for all interactive prompts
- Keyboard+Space multi-select for camera selection
- huh provides form-style flows (groups, inputs, selects) that fit a guided install
- Replaces manual prompt handling; huh runs before sudo re-exec

### Command Structure
- `onboard` subcommand under root, no required arguments
- `--dry-run` flag: generate/show config YAML to stdout without writing to /etc or installing service
- Pure interactive mode (no flag overrides for individual steps)
- Sudo re-exec happens after all prompts: run `sudo on-a-meet onboard` with collected answers passed via flags or temp file

### Config File & Sudo Flow
- Re-exec with sudo using collected answers (pass as flags or write temp JSON to /tmp/on-a-meet-onboard.json)
- Config written to `/etc/on-a-meet/config.yaml` (sudo path)
- Re-exec calls: write config → install service → print config path
- `--dry-run`: skip sudo, just print config to stdout and exit

### Detection Test UX
- Simple two-step test: "Enable your camera, press Enter" → Detect() → show result. Then "Disable your camera, press Enter" → Detect() → show result
- Both V4L2 and lsof methods testable this way
- Test runs with the selected detection method
- If test fails, offer to go back and pick the other method

### Debounce Default
- Onboard prompts with debounce default of 2 (instead of engine default 3)

### Agent's Discretion
- Specific prompt text / wording for the interactive wizard
- Error handling for failed detect tests

</decisions>

<specifics>
## Specific Ideas

- "use choose but multiselection" — huh MultiSelect for cameras
- "explain the difference between the methods" — brief text description before method select
- "after that the service must be installed and we should print where is the config file"

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `cmd/install.go` — existing service install pattern (kardianos/service, sudo check, serviceConfig)
- `internal/config/config.go` — Config struct that must be serialized to YAML
- `internal/detector/detector.go` — detector.New(method) factory
- `cmd/detect.go` — existing OnChange wiring for detection test

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `cmd/install.go`: `serviceConfig()`, `noopProgram{}`, `svc.Install()` + `svc.Start()` — can be reused after config write
- `internal/config/config.go`: `Config` struct with mapstructure tags — can be marshaled to YAML for file write
- `internal/detector/detector.go`: `detector.New(method)` — used for the detection test

### Established Patterns
- Cobra commands in `cmd/` package, `init()` to register
- `output.*` wrappers for all user-facing output
- Sudo check with `os.Geteuid()`
- Viper for config, but onboard will write raw YAML to file

### Integration Points
- New file: `cmd/onboard.go` — the onboard command
- Reuses service install from `cmd/install.go`'s helpers (extract to shared function)
- Detection test calls `detector.New(method)` then `det.Detect(path)` in a loop

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 06-onboard-command-assisted-install*
*Context gathered: 2026-05-29*
