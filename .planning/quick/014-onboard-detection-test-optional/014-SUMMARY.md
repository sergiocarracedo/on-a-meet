# Quick Task 014 Summary

**Task:** Camera test must be optional, and the question about keep the selected method or continue should only happen after a failing test (if user doesn't skip it)

**Completed:** 2026-05-29

## What was done

1. Added "Run detection test?" confirm with "Skip / Test" buttons after the form — user can skip the entire detection test
2. Moved the "Keep V4L2?" method change prompt to only appear when a detection test fails and the user chooses to skip
3. When skip-test is chosen, `testFailed = true` triggers the method change flow

## Files changed

- `cmd/onboard.go`: Restructured post-form flow — detection test is now optional, method change only on test failure

## Commit

`895949c`
