# Quick Task 014 Plan: Onboard detection test optional + method change after test failure

## Tasks

### Task 1: Make camera detection test optional, move method change to after test failure

**Files:**
- `cmd/onboard.go`

**Action:**

Two changes to the wizard flow:

**Change 1 — Prerequisite (move method change):** Delete the standalone `huh.NewConfirm("Keep V4L2?")` block (lines 187-207). It should only appear if a detection test fails and the user doesn't skip.

**Change 2 — Detection test optional:** After the form, add a `huh.NewConfirm` asking "Run detection test?" with "Skip / Test" buttons. If Skip, skip the entire test + method change and go to config. If Test, enter the detection test loop (ON + OFF).

**Change 3 — Method change after failure:** Inside the test loop, when the ON test fails (and retry says "Skip test") OR when the OFF test fails (and retry says "Skip test"), ask "Detection method: V4L2?" confirm with "Keep V4L2" / "Change method". Same logic as the removed block.

The new flow:
1. Form (camera, debounce, interval)
2. "Run detection test?" → Skip / Test
3. If Test → ON test loop → OFF test loop
4. If any test fails and skipped → "Keep V4L2?" → "Change method?" only if answer is "Change"
5. Config summary → install

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Detection test is optional (skip upfront)
- Method change only appears after a failing test
