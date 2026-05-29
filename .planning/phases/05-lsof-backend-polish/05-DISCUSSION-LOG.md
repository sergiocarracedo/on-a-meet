# Phase 5: lsof Backend & Polish — Discussion Log

**Date:** 2026-05-29
**Mode:** standard
**Participants:** User + Agent

## Areas Discussed

### 1. lsof Detection Implementation

**Options considered:**
- (A) Check lsof exit code + parse output — run `lsof /dev/videoN`, exit 0 = ON, exit 1 = OFF. Parse for process name/PID.
- (B) lsof exit code only — simpler, no process info in logs
- (C) Parse /proc/PID/fd/ instead — avoids external lsof dep but more complex

**Chosen:** A (Check lsof exit code + parse output)
**Rationale:** Process name/PID useful for logging. No external deps beyond lsof which is standard.

**Implications:**
- ListDevices() falls back to V4L2 enumeration (lsof can't enumerate)
- Files: `lsof_linux.go` + `lsof_stub.go` in `internal/detector/`

---

### 2. Backend Selection & Factory

**Options considered:**
- (A) `detector.New("v4l2")` factory — clean, extensible
- (B) Switch in detect command — simpler but less clean
- (C) Config-driven auto-detection — magic, less predictable

**Chosen:** A (Factory function in detector package)

**Rationale:** Clean separation, easy to add future backends. Config already has detect-method field.

**Implications:**
- New file `internal/detector/detector.go` with `New(method string) (Detector, error)`
- detect.go replaces `NewV4L2Detector()` with `detector.New(cfg.DetectMethod)`

---

### 3. List Command UX

**Options considered:**
- (A) Full table (Path, Driver, Card, Bus, Status) — most useful
- (B) Path + Driver only — too minimal
- (C) Full table without Status — misses key info

**Chosen:** A (Full table with Status column)
**Rationale:** Users want to know which cameras are currently in use. Running Detect() per device gives immediate feedback.

---

### 4. Release Scripting

**Options considered:**
- (A) Makefile — simpler, no external tooling
- (B) goreleaser — proper GitHub releases
- (C) Both — most complete but most maintenance

**Chosen:** B (goreleaser)
**Rationale:** Proper releases for binary distribution. Cross-compilation targets: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64.

---

### 5. README & --help Refinement

**Options considered:**
- (A) README + improved --help — balanced
- (B) README + man page + improved --help — too heavy for v1
- (C) README only (minimal) — insufficient

**Chosen:** A (README + improved --help)
**Rationale:** Covers all common needs. man page can be added post-v1 if needed.

---

### 6. lsof Package Layout (bonus clarification)

**Options considered:**
- (A) `lsof_linux.go` + `lsof_stub.go` in `internal/detector/` — same pattern as V4L2
- (B) Separate `internal/detector/lsof/` package — cleaner but more imports

**Chosen:** A (Files in detector package)
**Rationale:** Consistent with existing V4L2 pattern. Fewer imports, simpler integration.

## Deferred Ideas

None.
