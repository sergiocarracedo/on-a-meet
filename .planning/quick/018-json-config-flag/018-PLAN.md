# Quick Task 018: JSON config flag for onboard

## Task

**DESCRIPTION:** Add `--config <file>` flag to onboard command. When provided, read all config values from a JSON file and skip the interactive wizard entirely.

### Must haves
- `onboard --config <file>` reads JSON and skips all interactive prompts (form + detection test)
- Missing fields default to v4l2, 1s, 2
- `--dry-run` works with `--config`
- Shows confirm dialog before sudo re-exec

### Files
- `cmd/onboard.go`

### Action
1. Add `onboardConfigFile` variable and `--config` flag
2. In RunE, between `--apply` block and interactive wizard: parse JSON, apply defaults, skip to confirm+sudo flow
3. Empty cameras array → camera field omitted (monitor all)

### Verify
```bash
go build ./... && go vet ./... && go test ./...
# Dry-run with a JSON config produces correct output
go run . onboard --config /tmp/test-config.json --dry-run
```

### Done
When `onboard --config myconfig.json` skips all prompts and goes directly to install confirmation.
