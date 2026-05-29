# Quick Task 011 Summary

**Task:** All cameras selected by default; mutual exclusion between "All cameras" and individual cameras

**Completed:** 2026-05-29

## What was done

Replaced the broken `__all__` MultiSelect option with a proper two-step flow:
1. `huh.NewSelect` with "Monitor all cameras" (default) / "Choose specific cameras"
2. Conditional MultiSelect (hidden via `WithHideFunc` when "Monitor all cameras" is chosen)

Removed `__all__` string matching entirely — camera selection is now clean and predictable.

## Files changed

- `cmd/onboard.go`: Restructured camera selection groups, added `cameraChoice` variable, updated post-processing logic

## Commit

`57ed001`
