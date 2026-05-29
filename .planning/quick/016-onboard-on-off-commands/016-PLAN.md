# Quick Task 016 Plan: Add on/off command inputs to onboard wizard

## Tasks

### Task 1: Add on/off commands to structs, form, apply path, and preview

**Files:**
- `cmd/onboard.go`

**Action:**

Six edits:

1. **Structs** — Add OnCmd/OffCmd to `onboardConfig` (JSON tags `on-cmd`, `off-cmd`) and `writeConfig` (YAML tags `on-command`, `off-command`, omitempty)

2. **Form** — Add a new `huh.NewGroup` after the debounce/interval group with two `huh.NewInput` fields for on/off commands (optional, placeholders with examples)

3. **Apply path** — Add `cfg.OnCmd` / `cfg.OffCmd` to the `writeConfig` initialization

4. **Config assembly** — Add `OnCmd`/`OffCmd` to the `onboardConfig` struct literal after the test flow

5. **Dry-run preview** — Print `on-command` and `off-command` values

6. **Confirm dialog** — Add on/off to the summary description

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Onboard wizard asks for on/off commands and includes them in the written config
