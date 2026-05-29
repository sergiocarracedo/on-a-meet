# Phase 6 Research

## Don't Hand-Roll
- **Interactive prompts:** Don't build your own terminal prompt system. `huh` (v1.0.0, import: `github.com/charmbracelet/huh`) handles all keyboard input, multi-select, validation, and rendering. Use `huh.NewMultiSelect` for cameras, `huh.NewSelect` for method, `huh.NewInput` for debounce/interval.
- **YAML marshaling:** Don't hand-roll YAML output. Use `gopkg.in/yaml.v3` with the existing `Config` struct's `mapstructure` tags.
- **Sudo re-exec:** Don't try to elevate privileges within the process. Use `os.Executable()` + `exec.Command("sudo", ...)` pattern — clean process boundary.

## Common Pitfalls
- **Sudo environment:** `sudo` may strip the `$PATH`, so the binary path must be absolute (use `os.Executable()`). Ensure `sudo` doesn't reset the user's terminal — `--preserve-env` flag may be needed.
- **huh import path:** v1 uses `github.com/charmbracelet/huh` (stable). v2 uses `charm.land/huh/v2`. Use v1 for simplicity.
- **Temp file security:** If writing collected answers to `/tmp/`, use a predictable path but don't store secrets in it. The onboard config doesn't contain secrets, so this is acceptable.
- **Terminal interaction with sudo:** `sudo` running a binary with huh prompts may inherit terminal issues if stdin is not properly forwarded. Test with and without `sudo -k` to ensure stdin passthrough.

## Existing Patterns in This Codebase
- Cobra subcommands in `cmd/` package with `init()` registration (see cmd/install.go)
- `os.Executable()` pattern: available from stdlib but not yet used
- `serviceConfig()` + `noopProgram{}` in cmd/install.go — can be extracted to shared helper
- `internal/config/config.go`: `Config` struct can be marshaled directly to YAML

## Recommended Approach
1. Add `github.com/charmbracelet/huh` dependency (v1.0.0)
2. Create `cmd/onboard.go` with the interactive wizard flow
3. Use `--dry-run` to preview config without sudo
4. Use `--apply <file>` on the sudo re-exec path to consume collected answers
5. Reuse `serviceConfig()` from install.go for the service installation step
