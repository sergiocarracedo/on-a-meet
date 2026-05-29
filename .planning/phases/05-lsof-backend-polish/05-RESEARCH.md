# Phase 5: lsof Backend & Polish — Research

**Completed:** 2026-05-29

## Don't Hand-Roll

- **Command execution:** `os/exec` is sufficient for running `lsof`. No need for a process management library.
- **Release builds:** goreleaser v2 is the standard for Go releases. Avoid shell scripts for cross-compilation.
- **Config:** Viper already handles config precedence — no new config infrastructure needed.

## Common Pitfalls

- **lsof not installed:** lsof may not be present on minimal Docker/alpine systems. The error message should clearly indicate this.
- **lsof exit codes:** lsof returns exit code 1 when no process has the file open — this is the OFF state, not an error. Must distinguish exit code 1 from other errors (lsof not found, permission denied).
- **lsof output parsing:** The plan uses exit code only (not output parsing) — simpler and more reliable. Exit code 0 = ON, exit code 1 = OFF. No need to parse process name/PID for the detection logic.
- **Factory pattern:** Simple switch/case is sufficient. No reflection or plugin system needed.
- **goreleaser $GOPATH:** goreleaser v2 config uses `version: 2` format. `CGO_ENABLED=0` is essential for static binaries.

## Existing Patterns in This Codebase

- **Platform-specific files:** `_linux.go` + `_stub.go` pattern with build tags from Phase 2 (v4l2_linux.go/v4l2_stub.go) — lsof follows the same pattern.
- **Shared package functions:** `openAndQueryCap()` in v4l2_linux.go is used by LsofDetector.ListDevices() since lsof cannot enumerate devices.
- **Factory integration:** detect.go currently hardcodes `NewV4L2Detector()` at line 42. Replace with `detector.New(cfg.DetectMethod)`.
- **Output:** `output.Table()` in internal/output already accepts `pterm.TableData`.

## Recommended Approach

1. **lsof backend:** LsofDetector struct. Detect() runs `exec.Command("lsof", path)`, checks exit code. ListDevices() calls openAndQueryCap() from v4l2 package.
2. **Factory:** `detector.New(method string)` in new detector.go file. Switch on "v4l2" | "lsof".
3. **List command:** Use `detector.New()` like detect.go does. Enumeration + Detect() per device for Status column.
4. **goreleaser:** `.goreleaser.yaml` with linux/darwin + amd64/arm64 targets.
5. **README:** Standard install/config/usage docs.
