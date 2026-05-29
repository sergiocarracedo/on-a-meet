# Quick Task 013: Improve onboard detection test UX

## Tasks

### Task 1: Make detection method optional, fix multi-camera ON test

**Files:**
- `cmd/onboard.go`

**Action:**
1. Remove detection method from main form. Add `huh.NewConfirm("Keep V4L2?")` after form with default Yes. Show V4L2/lsof select only when user chooses "Change method".
2. Fix multi-camera ON test — require **at least one** ON (anyOn) instead of all ON (allOK).

**Result:**
- Conversation flow faster — method selection is one click away
- Multi-camera ON test passes when any camera is ON (not all)

**Commit:** `c747722`
