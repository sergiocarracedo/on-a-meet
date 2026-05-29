# Phase 5: lsof Backend & Polish — Verification

**Status:** passed

## Must-Have Checks

| Must-have | Status | Evidence |
|-----------|--------|----------|
| lsof_linux.go exists with LsofDetector | ✓ | `internal/detector/lsof_linux.go` |
| lsof_stub.go exists with !linux build tag | ✓ | `internal/detector/lsof_stub.go` |
| detector.go exists with New(method) factory | ✓ | `internal/detector/detector.go` |
| Factory tests pass | ✓ | go test ./internal/detector/ — 5 tests passing |
| go build passes | ✓ | Build OK |
| detect.go uses detector.New(cfg.DetectMethod) | ✓ | `cmd/detect.go` line 58 |
| list.go enumerates and displays table | ✓ | `cmd/list.go` — full implementation |
| .goreleaser.yaml exists | ✓ | linux/darwin + amd64/arm64 |
| README.md exists | ✓ | Install/config/usage/service docs |
| All tests pass | ✓ | 27 tests, all passing |

## Requirement Coverage

| Requirement | Status | Notes |
|-------------|--------|-------|
| REQ-013 (lsof backend) | ✓ | LsofDetector with factory selection |
| REQ-008 (list command) | ✓ | pterm table with device info and status |
