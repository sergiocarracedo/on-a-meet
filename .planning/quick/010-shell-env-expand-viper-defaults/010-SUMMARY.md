# Quick Task 010 Summary

**Task:** ${HASS_TOKEN} is not being replaced, CLI not showing when call is executed or output

**Completed:** 2026-05-29

## What was done

1. Added `os.ExpandEnv()` in executor after template rendering — `${HASS_TOKEN}` and `$VAR` in command strings are now expanded from environment variables at runtime
2. Added missing viper defaults for `interval`, `debounce`, `timeout`, `camera`, `on-command`, `off-command` — explains why `debounce=0` showed instead of `debounce=3`

## Files changed

- `internal/executor/executor.go`: Added `os.ExpandEnv(rendered)` after template render + `"os"` import
- `cmd/root.go`: Added `viper.SetDefault` for all 6 missing config fields

## Commit

`f7fdbe9`

## Shell quoting note

The `\${HASS_TOKEN}` in the user's command prevents shell expansion. Either:
1. Use `"..."` with `$HASS_TOKEN` (no backslash), shell expands it before the CLI sees it
2. Use the new env var expansion: store the token in an env var (e.g. `export HASS_TOKEN=...`) and the CLI will expand `${HASS_TOKEN}` at runtime
