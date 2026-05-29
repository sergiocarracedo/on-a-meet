# Quick Task 013 Summary

**Task:** Improve onboard detection test UX

**Completed:** 2026-05-29

## What was done

1. Detection method selection is now optional — removed from main form, added a `huh.NewConfirm("Keep V4L2?")` with "Keep V4L2" / "Change method" buttons. Method selection only shown when user explicitly chooses to change.
2. Multi-camera ON test now requires **at least one** camera ON (anyOn) instead of all cameras ON (allOK), which is more practical for multi-camera setups.

## Files changed

- `cmd/onboard.go`: Added method confirmation, changed ON test logic to anyOn, renamed `changeMethod` to `keepV4L2` for clarity

## Commit

`c747722`
