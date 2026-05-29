# Roadmap — on-a-meet v1.1.0

## Phase 1: macOS Detection Backend

**Goal:** Add macOS camera detection via `log stream` and device enumeration via `system_profiler`.

**Requirements:** REQ-014, REQ-015, REQ-016

**Tasks:**
1. Create `internal/detector/darwin.go` — macOSDetector with:
   - `Detect()`: exec `log stream --predicate` to detect camera stream events
   - `ListDevices()`: exec `system_profiler SPCameraDataType` for camera enumeration
2. Create `internal/detector/darwin_stub.go` — non-darwin stub
3. Update `internal/detector/detector.go` — add "darwin" detection method to factory
4. Update `cmd/list.go` — darwin uses default method
5. Tests: darwin detector smoketest (skipped without macOS)

**Success criteria:** `on-a-meet detect --detect darwin` compiles on macOS and detects camera on/off.

---

## Phase 2: README Documentation & Polish

**Goal:** Fix README gaps — add onboard command docs, macOS install instructions.

**Requirements:** REQ-017

**Tasks:**
1. Add `onboard` command usage section to README
2. Add macOS-specific install and permissions notes
3. Update `--help` output consistency

**Success criteria:** README covers all subcommands (detect, list, onboard, service).

---

*Last updated: 2026-05-29*
