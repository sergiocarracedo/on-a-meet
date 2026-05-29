# Quick Task 016: Add on/off command inputs to onboard wizard

## Done

Added `--on` / `--off` command inputs to the interactive onboard wizard:

- Added `OnCmd` / `OffCmd` fields to `onboardConfig` and `writeConfig` structs
- Added two new `huh.NewInput` fields in the form (optional, with template variable hints)
- Wired through the `--apply` path so on/off commands are written to `/etc/on-a-meet/config.yaml`
- Displayed in dry-run preview
- Committed as `<commit-sha>`

19 tests passing, build + vet clean.
