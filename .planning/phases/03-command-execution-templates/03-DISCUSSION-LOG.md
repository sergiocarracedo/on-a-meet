# Phase 3: Command Execution & Templates — Discussion Log

**Date:** 2026-05-28
**Mode:** standard
**Previous context files found:** None (fresh discussion)

## Gray Areas Discussed

### 1. Command Executor Package
**Options considered:**
- `internal/executor/` package with Executor struct ✓ (selected)
- Directly in `detect.go`
- Methods on Engine

**User chose:** `internal/executor/` package — clean separation, testable, reusable.

### 2. Shell vs No Shell
**Options considered:**
- `sh -c` via `os/exec` ✓ (selected)
- Direct `exec.Command` (no shell)

**User chose:** Shell — users expect pipes, redirects, chaining. Shell injection accepted as intentional.

### 3. Template Variables
**Options considered:**
- CameraID=device short, Device=full path, State=on/off ✓ (selected)
- Same + Driver + Card

**User chose:** The four-field set (CameraID, Device, State). Driver and Card available via DeviceInfo if needed later.

### 4. Overlap Prevention Policy
**Options considered:**
- Same-state skip, cross-state allow ✓ (selected)
- Sequential queue
- Kill previous same-state

**User chose:** Same-state skip + cross-state allow. If ON fires while ON is running, skip. ON→OFF always starts the off-command.

### 5. Timeout & Cancellation
**Options considered:**
- 30s default, configurable, 0 = no timeout ✓ (selected)
- No timeout
- 60s default

**User chose:** 30s default, configurable via config/--timeout, 0 = no timeout, kill+log on timeout.

### 6. Error Handling
**Options considered:**
- Print stderr + non-zero exit warning ✓ (selected)
- Silent (debug only)
- Forward stderr directly

**User chose:** Capture stdout/stderr, print warning via pterm with exit code. Never block polling loop.

## Areas Delegated to Agent's Discretion
- Executor struct internals (method signatures, field names)
- Test approach
- Whether `--on`/`--off` stays in detect.go config or moves to executor

## Deferred Ideas
None.

---
*Phase: 03-command-execution-templates*
*Discussion logged: 2026-05-28*
