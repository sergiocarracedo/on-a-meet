# Quick Task 013 Plan: Improve onboard detection test UX

## Tasks

### Task 1: Make detection method step optional

**Files:**
- `cmd/onboard.go`

**Action:**
Before the detection method form group, add a `huh.NewConfirm()` asking "Change detection method?" with default **No** (keyboard-accessible via arrow keys, just hit Enter to skip). If Yes, show the existing V4L2/lsof Select. If No, keep `method = "v4l2"` and skip the method Select.

Restructure: remove the static method group from the main form. Add the confirm + conditional select as inline code after the camera form but before debounce/interval.

### Task 2: ON test requires at least one camera ON (not all)

**Files:**
- `cmd/onboard.go`

**Action:**
In the ON detection test loop (the first `for {` block starting around line 232), change the logic from "all cameras must be ON" to "at least one camera must be ON". 

Change:
```go
allOK := true
for _, cam := range cameras {
    ...
    if !status.On {
        allOK = false
    }
}
```
To:
```go
anyOn := false
for _, cam := range cameras {
    ...
    if status.On {
        anyOn = true
    }
}
```
And update the retry prompt message accordingly. The OFF test should remain unchanged (all must be OFF).

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Detection method can be skipped (default V4L2)
- Multi-camera ON test passes with at least one camera ON
