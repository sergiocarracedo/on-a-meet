# Roadmap — on-a-meet

## Phase 1: Project Scaffold & CLI Foundation ✅

**Goal:** Runnable binary with cobra commands, config layer, and pterm output.

**Requirements:** REQ-001 (partial: scaffold only), REQ-007 (config file)

**Tasks:**
1. ✅ Initialize Go module with dependencies (cobra, viper, pterm, kardianos/service)
2. ✅ Create cobra root command with subcommand structure (detect, list, install)
3. ✅ Implement Viper config layer — YAML reading, CLI flag binding, precedence
4. ✅ Create config.yaml.example
5. ✅ pterm output helpers for terminal formatting

**Success criteria:** `on-a-meet --help` shows all subcommands. `--config ./config.yaml` is accepted.
**Completed:** 2026-05-28

---

## Phase 2: Camera Detection Engine ✅

**Goal:** Detect camera on/off state and fire commands on transitions.

**Requirements:** REQ-001, REQ-004, REQ-005, REQ-006, REQ-010, REQ-011, REQ-012

**Tasks:**
1. V4L2 detection backend — syscall-based device status check on /dev/video*
2. Device enumeration — list /dev/video* with driver info
3. Polling engine — configurable interval, per-device state tracking
4. Multi-camera OR logic — any on = on, all off = off
5. --camera flag to filter to specific device
6. Debounce window — N consecutive same-state polls before firing
7. Camera hotplug handling — graceful ENOENT, periodic re-scan
8. Permission check at startup with user-friendly error

**Success criteria:** Running `on-a-meet detect --interval 500ms` shows camera state changes in terminal output.
**Completed:** 2026-05-28

---

## Phase 3: Command Execution & Templates ✅

**Goal:** Execute user commands on state transitions with template variables.

**Requirements:** REQ-002, REQ-003

**Tasks:**
1. ✅ --on and --off flag parsing (command strings)
2. ✅ State change event → command dispatch in goroutine
3. ✅ text/template substitution ({{.CameraID}}, {{.Device}}, {{.State}})
4. ✅ Command timeout and cancellation (via context)
5. ✅ Prevent overlapping command execution for same state

**Success criteria:** `on-a-meet detect --on "echo on-{{.CameraID}}" --off "echo off"` prints correct template output on transition.
**Completed:** 2026-05-28

---

## Phase 4: Service Installation ✅

**Goal:** Install/uninstall as systemd service with proper error handling.

**Requirements:** REQ-009

**Tasks:**
1. Integrate kardianos/service with cobra command (install, uninstall, status)
2. Auto-detect OS (systemd on Linux, launchd on macOS)
3. Sudo detection and prompt for system-level install
4. Service unit generation with correct binary path and config flags

**Success criteria:** `sudo on-a-meet install` creates working systemd service. `on-a-meet uninstall` removes it.
**Completed:** 2026-05-29

---

## Phase 5: lsof Backend & Polish

**Goal:** Fallback detection method, documentation, and release readiness.

**Requirements:** REQ-013, REQ-012, REQ-008

**Tasks:**
1. lsof detection backend — parse `lsof /dev/video*` output
2. --detect flag to select backend (v4l2, lsof)
3. `on-a-meet list` command — pterm table of detected cameras
4. Release scripting (goreleaser or Makefile)
5. README with usage examples, install instructions, video group docs
6. man page or --help refinement

**Success criteria:** `on-a-meet list` shows all cameras. `--detect lsof` works as alternative detection method.

---

## Phase Dependencies

```
Phase 1 (Scaffold) → Phase 2 (Detection) → Phase 3 (Commands) → Phase 4 (Service)
                                                                      ↕
Phase 5 (lsof + Polish) ←──────────────────────────────────────────────┘
```

Phase 5 can run in parallel with Phase 4 or independently after Phase 3.

---
*Last updated: 2026-05-28*
