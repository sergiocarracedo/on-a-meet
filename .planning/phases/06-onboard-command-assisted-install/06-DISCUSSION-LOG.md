# Phase 6 Discussion Log

**Date:** 2026-05-29
**Phase:** 06 — Onboard Command — Assisted Install
**Mode:** Standard

## Gray Areas Discussed

### 1. Interactive UI Library

**Options considered:**
- `github.com/charmbracelet/huh` — modern form library, multi-select with keyboard+space (selected)
- `github.com/AlecAivazis/survey/v2` — well-known, simpler, less polished
- `github.com/charmbracelet/bubbletea` — full TUI framework, overkill
- pterm only — already dep, but no multi-select

**User choice:** huh (Recommended)

### 2. Command Structure

**Options considered:**
- Pure interactive — no flags (Recommended)
- Minimal flags — --camera, --method, --interval to skip steps
- `--dry-run` flag — generate config without installing (selected)

**User choice:** `--dry-run` flag

### 3. Config File & Sudo Flow

**Options considered:**
- Re-exec with sudo (Recommend) — collected input → sudo re-exec → write config + install (selected)
- Write instructions — print sudo instructions for manual copy-paste
- Split: sudo at start — require root from the beginning

**User choice:** Re-exec with sudo

### 4. Detection Test UX

**Options considered:**
- Simple test (Recommended) — "Enable camera, press Enter" → Detect() → show; "Disable, press Enter" → show (selected)
- Brief poll test — 3-5 cycles showing live status
- Skip test — just textual explanation

**User choice:** Simple test

### Agent's Discretion Areas
- Prompt text / wording
- Error handling for failed tests

### Deferred Ideas
None
