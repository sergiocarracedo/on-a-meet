# Quick Task 010 Plan: Expand env vars in commands + fix viper defaults

Two issues:
1. `${HASS_TOKEN}` is not expanded — the command template should support `$VAR` / `${VAR}` via `os.ExpandEnv` so secrets can be passed through environment variables without exposing them in the flag string
2. `debounce=0` shows in config — viper defaults for config fields (debounce, etc.) are missing, so users see zero values instead of actual defaults

## Tasks

### Task 1: Expand env vars in rendered command

**Files:**
- `internal/executor/executor.go`

**Action:**
After `rendered := buf.String()` (template rendering), call `os.ExpandEnv(rendered)` so environment variables like `$HASS_TOKEN` or `${HASS_TOKEN}` are expanded before the command runs. This allows users to pass secrets via env vars rather than embedding tokens in CLI flags.

Add `"os"` to imports.

### Task 2: Set viper defaults for all config fields

**Files:**
- `cmd/root.go`

**Action:**
Add `viper.SetDefault` calls for `interval`, `debounce`, `timeout`, `camera`, `on-command`, `off-command` alongside the existing `viper.SetDefault("detect-method", "v4l2")`. Use defaults from `config.Defaults()`.

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Executor expands `$VAR` and `${VAR}` from environment variables in command templates
- Viper defaults are set so config line shows correct values (e.g., `debounce=3`)
