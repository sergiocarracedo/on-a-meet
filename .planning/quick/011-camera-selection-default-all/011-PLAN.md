# Quick Task 011 Plan: Fix camera selection UX

## Background

The camera selection MultiSelect has two problems:
1. All cameras are not selected by default
2. "All cameras" vs individual cameras have no mutual exclusion

## Approach

Replace the `__all__` MultiSelect option pattern with a cleaner two-step flow:
- Step 1: Ask "Monitor all cameras" vs "Choose specific cameras" (default: "Monitor all cameras")
- Step 2: Only show the camera MultiSelect when "Choose specific cameras" is selected, using `huh.NewGroup.WithHideFunc`

## Tasks

### Task 1: Rewrite camera selection with select-all toggle

**Files:**
- `cmd/onboard.go`

**Action:**
1. Add a `cameraChoice` string variable
2. Add a new `huh.NewGroup` BEFORE the MultiSelect group with a `huh.NewSelect[string]`:
   - Title: "Camera selection"
   - Options: "Monitor all cameras" (value "all"), "Choose specific cameras" (value "select")
   - Default: "all"
   - Value: `&cameraChoice`

3. Modify the MultiSelect group to be hidden when `cameraChoice == "all"`:
   - Remove `__all__` from deviceOpts (only individual device paths remain)
   - Add `.WithHideFunc(func() bool { return cameraChoice == "all" })` to the group

4. Update the post-processing:
   - If `cameraChoice == "all"`: use all device paths (ignore camSelections)
   - If `cameraChoice == "select"`: use camSelections as-is

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- All cameras selected by default (user sees "Monitor all cameras" pre-selected)
- When user chooses "Choose specific cameras", they get individual camera toggles
- No `__all__` vs individual camera confusion
