# Quick Task 009 Plan: Redact JWT tokens from CLI output

When commands contain Bearer tokens (JWT), the config line and debug output leak secrets. Add a function to detect and redact JWT-like tokens (strings starting with `ey`, 3 base64url segments separated by dots) before displaying.

**Files to modify:**
- `internal/output/output.go` — add `RedactSecrets(s string) string`
- `internal/executor/executor.go` — apply redaction to rendered command and output
- `cmd/detect.go` — apply redaction to on/off config values in config line

## Tasks

### Task 1: Add RedactSecrets to output package

Add a function that uses a regex to find JWT-like tokens (`ey[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`) and replaces them with `ey***.***.***`.

Import `regexp`, compile once at package level.

### Task 2: Redact in config line

Wrap `cfg.OnCmd` and `cfg.OffCmd` with `output.RedactSecrets()` in the config printfln at `cmd/detect.go:65`.

### Task 3: Redact in executor debug output

Wrap `rendered` and `outStr` with `output.RedactSecrets()` in the debug printfln calls at `internal/executor/executor.go:89-90`.

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass
