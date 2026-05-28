# Phase 1: Project Scaffold & CLI Foundation — Research

**Date:** 2026-05-28
**Confidence:** HIGH
**Sources:** Cobra user guide, Viper README, pterm.sh docs, pkg.go.dev/cobra, multiple Go CLI tutorials (2025-2026)

## Don't Hand-Roll

### Command-Line Parsing
- **Don't write your own CLI parser.** Cobra handles: command tree, POSIX flags (short/long), help generation, shell completions (bash/zsh/fish/powershell), version flag, PreRun/PostRun hooks, Levenshtein-based "did you mean" suggestions, and flag groups (required-together, mutually-exclusive) — [VERIFIED: cobra.dev/user_guide]
- **Use `RunE` not `Run`** for commands that can fail — Cobra propagates errors up and handles exit codes — [VERIFIED: cobra.dev/user_guide]
- **Use persistent flags for cross-command flags** (`--config`, `--verbose`, `--silent`) on root, local flags for command-specific options (`--on`, `--off`, `--interval` on detect) — [VERIFIED: cobra.dev/user_guide]

### Configuration Management
- **Don't write your own config loader.** Viper handles: YAML/JSON/TOML parsing, precedence chain (explicit set > flags > env > config file > defaults), multiple search paths, env var binding, live config watching — [VERIFIED: github.com/spf13/viper]
- **Use `viper.BindPFlags(cmd.Flags())`** for automatic flag-config binding rather than per-flag `BindPFlag` calls — reduces boilerplate and prevents missing bindings — [VERIFIED: viper README]
- **Use `cobra.OnInitialize(initConfig)`** pattern — this registers the Viper init function to run before any command executes, on the root command — [VERIFIED: cobra.dev/user_guide]

### Terminal Output
- **Don't write ANSI escape codes by hand.** pterm provides: colored output (Info/Success/Warning/Error with defaults), tables, sections/headers, spinners, progress bars, bullet lists, and theming — [VERIFIED: pterm.sh]
- **Use pterm's built-in printers** (`pterm.Info.Println`, `pterm.Success.Println`, `pterm.Warning.Println`, `pterm.Error.Println`) rather than custom color functions — [VERIFIED: github.com/pterm/pterm]
- **pterm respects `Output` and `RawOutput` flags** — setting `pterm.RawOutput = true` disables all styling for pipe-friendly output — [VERIFIED: pterm source]

## Common Pitfalls

### Cobra + Viper Binding Order
- **Pitfall:** Calling `BindPFlags` before flags are defined — results in silent no-op binding. [CITED: cobra.dev/user_guide]
- **Fix:** Define all flags in `init()`, then call `viper.BindPFlags(cmd.Flags())` in the root command's `PersistentPreRunE` or `initConfig` function — [ASSUMED: common Go CLI pattern]

### Viper Gets from Flags vs Config
- **Pitfall:** `viper.GetString("interval")` may return the flag default even when the config file has a different value, if `BindPFlag` binds the flag but `BindPFlags` hasn't been called yet. [CITED: viper README]
- **Fix:** Always use `viper.BindPFlags(cmd.Flags())` before accessing config values in a command's `RunE`. The `cobra.OnInitialize(initConfig)` pattern ensures this runs before command execution.

### pterm and Piped Output
- **Pitfall:** ANSI color codes in piped output (`on-a-meet list | grep something`) creates garbage. [ASSUMED]
- **Fix:** `pterm` doesn't auto-detect pipes. Check `isatty.IsTerminal(os.Stdout.Fd())` or use the `--silent`/`--verbose` flags as specified in CONTEXT.md. If not a TTY, set `pterm.RawOutput = true`.

### Viper _does not_ deep merge complex values
- **Pitfall:** If a config file defines `camera: /dev/video0` and env var defines `ON_A_MEET_CAMERA`, Viper replaces, not merges. For top-level scalars this is fine. For nested maps/slices, it replaces entirely. [VERIFIED: viper README]
- **Fix:** Our config is flat top-level scalars — no nested structures, so this doesn't apply to v1.

### Go Module Versioning with Indirect Dependencies
- **Pitfall:** `go get github.com/spf13/cobra@latest` and `go get github.com/pterm/pterm@latest` may pull incompatible indirect dependencies (e.g., different pflag versions). [ASSUMED]
- **Fix:** Use `go get <pkg>@<specific-version>` with the versions from STACK.md. Run `go mod tidy` after adding all deps.

### Only the first PersistentPreRun/PersistentPostRun in the command chain executes by default
- **Pitfall:** If `detect` defines its own `PersistentPreRun`, the root's `PersistentPreRun` (which sets up Viper) may not run. [VERIFIED: cobra.dev/user_guide]
- **Fix:** The root `initConfig` via `cobra.OnInitialize` runs regardless of which command is targeted. For Viper init, use `cobra.OnInitialize(initConfig)` in root's `init()`, not `PersistentPreRun`. The `initConfig` function will execute before any command.

## Existing Patterns in This Codebase

**None.** This is a green-field project. No existing Go modules, packages, or patterns to follow or reuse. All conventions established here will be the precedent for Phases 2-5.

The `.planning/` directory has project metadata only. The `AGENTS.md` at root has no codebase conventions beyond those set in the discuss-phase.

## Recommended Approach

### Project Structure

```
on-a-meet/
├── cmd/
│   ├── root.go          # Root command, initConfig, Viper setup
│   ├── detect.go        # detect subcommand (flags defined, RunE placeholder)
│   ├── list.go          # list subcommand (placeholder)
│   ├── install.go       # install subcommand (placeholder)
│   └── uninstall.go     # uninstall subcommand (placeholder)
├── internal/
│   ├── config/
│   │   └── config.go    # Config struct, defaults, YAML schema
│   └── output/
│       └── output.go    # pterm wrapper functions
├── main.go              # Entry point: cmd.Execute()
├── config.yaml.example  # Example config file
├── go.mod
└── go.sum
```

### Cobra + Viper Init Flow

1. `main.go` calls `cmd.Execute()`
2. **Before** any command runs, `cobra.OnInitialize(initConfig)` fires:
   - Set config name/type/paths (`on-a-meet`, `yaml`, XDG paths)
   - `viper.AutomaticEnv()` for env var binding
   - `viper.ReadInConfig()` (reads config file if exists, no error if missing)
3. If `--config` flag is set (in root's `PersistentPreRunE`):
   - Override Viper config path with the explicit file
4. Command-specific `RunE` accesses `viper.Get<Type>()` for all flag values

### Config Precedence (Viper default)
1. CLI flags (highest)
2. Environment variables (`ON_A_MEET_*`)
3. Config file (`~/.config/on-a-meet/config.yaml`)
4. Defaults (lowest)

### pterm Wrapper Pattern

```go
package output

import "github.com/pterm/pterm"

func Init(silent, verbose bool) {
    if silent {
        pterm.SetDefaultOutput(io.Discard)
    }
}

func Info(msg string)    { pterm.Info.Println(msg) }
func Success(msg string) { pterm.Success.Println(msg) }
func Warning(msg string) { pterm.Warning.Println(msg) }
func Error(msg string)   { pterm.Error.Println(msg) }
func Table(data pterm.TableData) { pterm.DefaultTable.WithData(data).Render() }
```

### Key Libraries and Versions

| Library | Version | Purpose |
|---------|---------|---------|
| spf13/cobra | v1.8+ | CLI framework |
| spf13/viper | v1.19+ | Config management |
| pterm/pterm | v0.12+ | Terminal output |
| kardianos/service | v1.2+ | Service management (scaffold only) |

### Testing Approach

- **Alongside source code** per Go convention (per CONTEXT.md)
- Test flag binding: use `cmd.SetArgs()` + capture output to verify flag/config interactions
- Test output helpers: redirect stdout, capture, assert content
- `go test ./internal/...` for package-level tests
- Integration test: build binary, run with test args, verify exit codes
