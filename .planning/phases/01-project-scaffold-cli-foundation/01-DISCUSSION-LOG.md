# Phase 1: Project Scaffold & CLI Foundation - Discussion Log

**Date:** 2026-05-28
**Mode:** Standard
**Workflow:** discuss-phase 1

## Gray Areas Discussed

### 1. Project Layout

**Options considered:**
- `cmd/ + internal/` layout (Recommended) — **SELECTED**
- Flat `cmd/` package
- Domain-driven packages at root

**Internal package granularity:**
- One package per concern (Recommended) — **SELECTED**
- Monolithic internal package
- Flat packages at root

**Package naming:**
- Full descriptive names (Recommended) — **SELECTED**
- Short abbreviations

**Test placement:**
- Alongside source code (Recommended) — **SELECTED**
- Separate test/ directory

**User moved to next area after 4 questions.**

---

### 2. Command Tree Design

**Options considered:**
- Root with subcommands (Recommended) — **SELECTED**
- Root as detect command
- Flat flags on root

**When to define detect flags:**
- Define on detect subcommand now (Recommended) — **SELECTED**
- Define only Phase 1 flags now

**--config flag scope:**
- Persistent on root (Recommended) — **SELECTED**
- Local to detect only

**Version handling:**
- `--version` flag on root (Recommended) — **SELECTED**
- `version` subcommand
- None for now

**User moved to next area after 4 questions.**

---

### 3. Config Struct & YAML Schema

**Default config path:**
- `$HOME/.config/on-a-meet/config.yaml` (Recommended) — **SELECTED**
- `$HOME/.on-a-meet.yaml`
- `./on-a-meet.yaml`

**Config struct fields:**
- Full v1 surface (Recommended) — **SELECTED**
- Minimal Phase 1 only

**Viper binding strategy:**
- `viper.BindPFlags` automatic (Recommended) — **SELECTED**
- Manual per-flag binding

**Config file resolution:**
- Viper AddConfigPath + automatic search (Recommended) — **SELECTED**
- Manual file read

**User moved to next area after 4 questions.**

---

### 4. pterm Output Patterns

**Helper architecture:**
- Thin wrapper functions (Recommended) — **SELECTED**
- Direct pterm calls
- Logger interface

**Log levels:**
- `--silent` + `--verbose` flags (Recommended) — **SELECTED**
- No log levels
- slog/Logrus

**Flag location:**
- Persistent flags on root (Recommended) — **SELECTED**
- Local to detect

**Color scheme:**
- Modern terminal colors (Recommended) — **SELECTED**
- Monochrome

**Startup banner:**
- Show startup banner with detected cameras — **SELECTED** (user chose this over "No banner")
- Minimal output (Recommended)

**User confirmed wrap-up after 5 questions.**

---

## Agent's Discretion Items

- Exact function signatures of output helpers
- Config struct field ordering
- Tab completion setup

## Deferred Ideas

None
