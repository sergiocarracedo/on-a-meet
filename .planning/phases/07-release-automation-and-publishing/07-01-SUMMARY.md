# Plan 07-01: Release Automation — Summary

## Built

- `.goreleaser.yaml` — cross-platform builds for linux/darwin (amd64 + arm64), CGO_ENABLED=0, ldflags with version/commit/date, checksums
- `.github/workflows/release.yml` — GitHub Actions workflow triggered on tag push `v*`, runs GoReleaser to build binaries and upload to GitHub Releases

## Removed

- npm wrapper package (`packages/npm/`) — not needed for a Go CLI tool
- Homebrew `brews` section from .goreleaser.yaml — can be re-added later if needed
- `publish-npm` job from release workflow — simplifies to single release job

## Verification

- `go build ./...` and `go vet ./...` pass
- All tests pass
- Release workflow is manual (tag push only)
