# Project State

## Phase 1: Project Scaffold & CLI Foundation

**Status:** ✅ Complete (2026-05-28)
**Last action:** execute-phase 1 (2026-05-28)
**Context file:** `.planning/phases/01-project-scaffold-cli-foundation/01-CONTEXT.md`
**Plans:** 2 plans, 2 waves — all executed
**Verification:** Passed — 10/10 must-haves, 3/3 tests passing

## Phase 2: Camera Detection Engine

**Status:** ✅ Complete (2026-05-28)
**Last action:** execute-phase 2 (2026-05-28)
**Context file:** `.planning/phases/02-camera-detection-engine/02-CONTEXT.md`
**Plans:** 2 plans, 2 waves — all executed
**Verification:** Passed — 8/8 tests passing

## Phase 3: Command Execution & Templates

**Status:** ✅ Complete (2026-05-28)
**Last action:** execute-phase 3 (2026-05-28)
**Context file:** `.planning/phases/03-command-execution-templates/03-CONTEXT.md`
**Plans:** 2 plans, 2 waves — all executed
**Research:** Completed — `03-RESEARCH.md` (exec.CommandContext, template, process group management)
**Verification:** Passed — 22/22 tests passing, all must-haves met

## Phase 4: Service Installation

**Status:** ✅ Complete (2026-05-29)
**Last action:** execute-phase 4 (2026-05-29)
**Context file:** `.planning/phases/04-service-installation/04-CONTEXT.md`
**Decisions:** 5 areas discussed — integration pattern, service arguments, lifecycle, binary path, sudo handling
**Research:** Completed — `04-RESEARCH.md` (kardianos/service one-shot API, sudo check, binary path)
**Plans:** 2 plans, 2 waves — all executed
**Verification:** Passed — build passes, all must-haves met
| | Date | Workflow | Result |
| |------|----------|--------|
| | 2026-05-28 | new-project | Milestone v1 initialized, 5-phase roadmap, 12 requirements |
| | 2026-05-28 | discuss-phase 1 | 4 gray areas discussed, decisions captured |
| | 2026-05-28 | plan-phase 1 | 2 plans created across 2 waves |
| | 2026-05-28 | execute-phase 1 | Phase 1 implemented and verified |
| | 2026-05-28 | discuss-phase 2 | 4 gray areas discussed, decisions captured |
| | 2026-05-28 | plan-phase 2 | 2 plans created across 2 waves, research + verification passed |
| | 2026-05-28 | execute-phase 2 | Phase 2 implemented and verified — 8 tests passing |
| | 2026-05-28 | discuss-phase 3 | 6 gray areas discussed, decisions captured |
| | 2026-05-28 | plan-phase 3 | 2 plans created across 2 waves, research + verification passed |
| | 2026-05-28 | execute-phase 3 | Phase 3 implemented and verified — 22 tests passing |
| | 2026-05-28 | discuss-phase 4 | 5 gray areas discussed, decisions captured |
| | 2026-05-28 | plan-phase 4 | 2 plans created across 2 waves, research + verification passed |
| | 2026-05-29 | execute-phase 4 | Phase 4 implemented and verified — build passes, all tests passing |
| | 2026-05-29 | discuss-phase 5 | 5 gray areas discussed, decisions captured |
| | 2026-05-29 | plan-phase 5 | 2 plans created across 2 waves, research + verification passed |
| | 2026-05-29 | execute-phase 5 | Phase 5 implemented and verified — 27 tests passing, goreleaser + README |
| | 2026-05-29 | add-phase | Phase 6 added: Onboard Command — Assisted Install |
| | 2026-05-29 | discuss-phase 6 | 4 gray areas discussed: huh lib, --dry-run flag, sudo re-exec, simple detect test |
| | 2026-05-29 | plan-phase 6 | 2 plans, 2 waves — huh wizard (wave 1), sudo apply+install (wave 2) |
| | 2026-05-29 | execute-phase 6 | Phase 6 implemented — onboard command with wizard + sudo apply path |

## Phase 5: lsof Backend & Polish

**Status:** ✅ Complete (2026-05-29)
**Last action:** execute-phase 5 (2026-05-29)
**Context file:** `.planning/phases/05-lsof-backend-polish/05-CONTEXT.md`
**Decisions:** 5 areas discussed — lsof detection, backend factory, list command, goreleaser, README
**Research:** Completed — `05-RESEARCH.md` (lsof exit codes, factory pattern, goreleaser)
**Plans:** 2 plans, 2 waves — all executed
**Verification:** Passed — build passes, 27 tests passing

## Phase 6: Onboard Command — Assisted Install

**Status:** ✅ Complete (2026-05-29)
**Last action:** execute-phase 6 (2026-05-29)
**Context file:** `.planning/phases/06-onboard-command-assisted-install/06-CONTEXT.md`
**Decisions:** 4 areas discussed — huh library, --dry-run flag, sudo re-exec, simple detect test
**Research:** Completed — `06-RESEARCH.md` (huh API: NewForm, NewMultiSelect, NewSelect, NewInput)
**Plans:** 2 plans, 2 waves — all executed
**Verification:** Passed — build passes, 19 tests passing

### Quick Tasks Completed

| # | Description | Date | Commit | Directory |
|---|-------------|------|--------|-----------|
| 001 | Fix --detect flag defaulting to empty string | 2026-05-29 | `4b83107` | `quick/001-detect-flag-default/` |
| 002 | Print config on detect startup | 2026-05-29 | `196ad05` | `quick/002-print-config-on-startup/` |
| 003 | Fix V4L2 camera ON detection via /proc scan | 2026-05-29 | `1c2b662` | `quick/003-v4l2-detect-on-camera/` |
| 004 | Fix --detect flag Viper bind clash detect/list | 2026-05-29 | `5e9e50d` | `quick/004-detect-flag-bind-fix/` |
| 005 | Emit initial camera status on detect start | 2026-05-29 | `d33a9df` | `quick/005-emit-initial-camera-status-on-start/` |
| 006 | Log on/off callback exit code and verbose output | 2026-05-29 | `8ca2245` | `quick/006-log-on-off-callbacks-exit-code/` |
| 007 | Print all config values on startup | 2026-05-29 | `45c7fa6` | `quick/007-print-all-config-values/` |
| 008 | Fix verbose detection and callback execution display | 2026-05-29 | `a30790c` | `quick/008-fix-verbose-callback-display/` |
| 009 | Redact JWT tokens from CLI output | 2026-05-29 | `cdf32b8` | `quick/009-redact-tokens-in-output/` |
| 010 | Expand env vars in commands and fix viper defaults | 2026-05-29 | `f7fdbe9` | `quick/010-shell-env-expand-viper-defaults/` |
| 011 | Fix camera selection UX (all default, select-all toggle) | 2026-05-29 | `57ed001` | `quick/011-camera-selection-default-all/` |
| 012 | Fix config loading from onboard | 2026-05-29 | `60dbd9d` | `quick/012-fix-config-loading/` |
| 013 | Improve onboard detection test UX | 2026-05-29 | `c747722` | `quick/013-improve-onboard-detection-test-ux/` |
| 014 | Make detection test optional, show method change only after failure | 2026-05-29 | `895949c` | `quick/014-onboard-detection-test-optional/` |
| 015 | Stop service before reinstall, prompt before config overwrite | 2026-05-29 | `ee1987b` | `quick/015-onboard-stop-service-overwrite-config/` |
| 016 | Add on/off command inputs to onboard wizard | 2026-05-29 | `895163d` | `quick/016-onboard-on-off-commands/` |
| 017 | Fix YAML on/off command quoting in onboard --apply | 2026-05-29 | `19eb82d` | `quick/017-yaml-on-off-command-quoting/` |
| 018 | Add --config flag to onboard for JSON config input | 2026-05-29 | `762a831` | `quick/018-json-config-flag/` |
| 019 | Accept "all" string for cameras in --config JSON | 2026-05-29 | `0808bce` | `quick/019-cameras-all-string/` |
| 020 | Service runs commands as original user | 2026-05-29 | `789003b` | `quick/020-service-run-as-original-user/` |

## Next
- Milestone v1 complete — all 6 phases done
- Last activity: 2026-05-29 - Completed quick task 020: Service runs commands as original user

### Roadmap Evolution
- 2026-05-29: Phase 6 added + completed: Onboard Command — Assisted Install
- 2026-05-29: Phase 7 added: Release Automation & Publishing
