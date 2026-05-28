# Phase 1: Project Scaffold & CLI Foundation - Context

**Gathered:** 2026-05-28
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Initialize the Go module, set up the Cobra command tree with all subcommands (detect, list, install, uninstall), implement the Viper config layer (YAML reading + CLI flag binding + precedence), create config.yaml.example, and set up pterm output helpers. First code in the repo — green-field start.

</domain>

<decisions>
## Implementation Decisions

### Project Layout
- **Structure:** `cmd/ + internal/` layout — `cmd/on-a-meet/main.go` as entry point, `internal/` for all library code
- **Internal granularity:** One package per concern — `internal/config/`, `internal/output/`, `internal/detector/`, `internal/command/` (future phases add packages)
- **Package naming:** Full descriptive names (e.g., `cameradetector`, `commandexecutor`) rather than abbreviations
- **Test placement:** Alongside source code per Go convention — `config_test.go` next to `config.go`

### Command Tree Design
- **Root behavior:** Root command shows help, no action — all functionality via subcommands
- **Subcommands:** `detect` (main polling), `list` (enumerate cameras), `install` (service install), `uninstall` (service remove)
- **Flag definition timeline:** All detect subcommand flags defined in Phase 1 (`--on`, `--off`, `--interval`, `--camera`, `--config`, `--silent`, `--verbose`) — wiring happens in later phases
- **--config flag:** Persistent flag on root command, Viper reads before subcommand execution
- **--version:** Flag on root command, version set via ldflags at build time

### Config Struct & YAML Schema
- **Default path:** XDG-compliant `$HOME/.config/on-a-meet/config.yaml`
- **Schema surface:** Full v1 YAML surface defined in Phase 1 — all fields present even if wiring deferred:
  - `camera` (string, optional — target specific device)
  - `interval` (duration, default 1s)
  - `on-command` (string — command to run when camera turns on)
  - `off-command` (string — command to run when camera turns off)
  - `detect-method` (string, default "v4l2" — future use)
  - `debounce` (int, default 3 — future use)
- **Binding strategy:** `viper.BindPFlags(cmd.Flags())` in PersistentPreRunE — automatic, flag vs config file precedence via Viper
- **Config discovery:** Viper's `AddConfigPath` + `SetConfigName` + `SetConfigType` — automatic search of default paths

### pterm Output Patterns
- **Helper architecture:** Thin wrapper functions in `internal/output/` — `Info()`, `Success()`, `Warning()`, `Error()`, `Table()`, `StartupBanner()`
- **Quiet/verbose flags:** `--silent` (suppress output) and `--verbose` (debug info) as persistent flags on root
- **Color scheme:** Default pterm colors — Info (blue), Success (green), Warning (yellow), Error (red)
- **Startup behavior:** Show startup banner listing detected cameras on `detect` start

### Agent's Discretion
- Exact function signatures of output helpers — stick to pterm conventions
- Config struct field ordering — group by category
- Tab completion setup — add if straightforward

</decisions>

<specifics>
## Specific Ideas

- Full CLI surface area from Phase 1 even though some flags won't be wired until later phases — gives users immediate familiarity and prevents flag naming conflicts

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — v1 requirement definitions
- `.planning/research/ARCHITECTURE.md` — 5-component architecture with Detector interface design
- `.planning/research/STACK.md` — Go version, library versions, rationale

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- None — green-field project, no existing code

### Established Patterns
- None — first code to be written establishes all patterns

### Integration Points
- This phase creates all integration surfaces: command tree (where future phases register), config struct (where future fields are added), output helpers (used by all phases)

</code_context>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---
*Phase: 01-project-scaffold-cli-foundation*
*Context gathered: 2026-05-28*
