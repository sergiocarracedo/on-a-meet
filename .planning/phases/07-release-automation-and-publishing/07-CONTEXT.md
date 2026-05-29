# Phase 7: Release Automation & Publishing — Context

## Overview

Set up automated release pipelines for `on-a-meet` following the pattern from `skill-organizer`:
- GoReleaser for cross-platform builds + GitHub Releases
- Homebrew formula publishing (personal tap)
- npm wrapper package (downloads binary on install)

## Manual Steps

### 1. GitHub Releases (already partially set up)

`.goreleaser.yaml` exists from Phase 5. To trigger a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

GoReleaser Action builds binaries and uploads to GitHub Releases.

### 2. Homebrew

**Step 1:** Create personal tap repository:
- Go to https://github.com/new
- Name: `homebrew-tap`
- Description: "Homebrew tap for on-a-meet"
- Create empty repo

**Step 2:** Generate a GitHub PAT with `repo` scope for the tap repo.

**Step 3:** Add secret to on-a-meet repo:
- Settings → Secrets and variables → Actions
- Add `HOMEBREW_TAP_GITHUB_TOKEN` with the PAT

**Step 4:** Update `.goreleaser.yaml` to add `brews` section (see goreleaser config below).

### 3. npm Wrapper Package

**Step 1:** Create `packages/npm/` directory structure.

**Step 2:** Create `package.json`:
```json
{
  "name": "on-a-meet",
  "version": "0.0.0",
  "description": "CLI tool to detect camera on/off state",
  "bin": {
    "on-a-meet": "bin/on-a-meet.js"
  },
  "scripts": {
    "postinstall": "node scripts/postinstall.js"
  }
}
```

**Step 3:** Create `bin/on-a-meet.js` — shebang script that spawns the downloaded binary.

**Step 4:** Create `scripts/postinstall.js` — downloads the correct platform binary from GitHub Release, verifies SHA256, extracts to `vendor/`.

**Step 5:** Enable npm trusted publishing:
- Go to npmjs.com → Access → Packages → on-a-meet → Publishing access → Set up
- Or use `npm token create --publish` for initial setup

**Step 6:** Publish: `npm publish --access public`

### 4. GitHub Actions Workflow

Create `.github/workflows/release.yml` that:
- Triggers on tag push `v*`
- Runs GoReleaser
- Publishes npm package
