# Quick Task 012 Summary

**Task:** Fix config loading — `--camera` not overriding onboard choices, config not printed

**Completed:** 2026-05-29

## What was done

1. Added `/etc/on-a-meet` as an additional viper config search path in `initConfig()` so that configs created by `onboard --apply` (which writes to `/etc/on-a-meet/config.yaml`) are found when running `detect` manually
2. Added viper defaults for `silent` and `verbose` keys so they're available via `viper.GetBool("silent")` / `viper.GetBool("verbose")` if ever needed

## Files changed

- `cmd/root.go`: Added `/etc/on-a-meet` to search paths; added silent/verbose viper defaults

## Commit

`60dbd9d`
