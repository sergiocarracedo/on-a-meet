# Plan 07-01: Release Automation & Publishing — Summary

## Built

- `.goreleaser.yaml` — cross-platform builds for linux/darwin (amd64 + arm64), CGO_ENABLED=0, ldflags with version/commit/date, commented-out Homebrew `brews` section
- `.github/workflows/release.yml` — GitHub Actions workflow triggered on tag push `v*` with two jobs: `goreleaser` (builds + uploads) and `publish-npm` (publishes npm wrapper)
- `packages/npm/package.json` — npm package `on-a-meet` with binary entry point and postinstall download script
- `packages/npm/bin/on-a-meet.js` — Node.js shebang that spawns the downloaded platform binary
- `packages/npm/scripts/postinstall.js` — downloads the correct platform binary from GitHub Release, extracts to `vendor/`, sets executable
- `07-CONTEXT.md` — context and steps for each platform (GitHub Releases, Homebrew, npm)
- `07-MANUAL-STEPS.md` — detailed step-by-step instructions for all three distribution channels

## Key Decisions

- npm trusted publishing via `id-token: write` (not token-based)
- Homebrew `brews` section kept commented-out until `homebrew-tap` repo is created
- npm version kept at `0.0.0` — CI replaces it dynamically from the Git tag

## Verification

- `go build ./...` and `go vet ./...` pass
- All 19 tests pass
- Release workflow is manual (tag push only) — no automatic CI runs without intentional action
