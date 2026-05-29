# Phase 7: Release Automation — Context

## Overview

Set up GoReleaser to build cross-platform binaries and upload them to GitHub Releases on tag push.

## Release Workflow

1. Push a tag matching `v*`:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
2. GitHub Actions runs GoReleaser which builds linux/darwin + amd64/arm64 binaries
3. Binaries are uploaded to the GitHub Release as tarballs with checksums

## Distribution

Users download the appropriate tarball from the GitHub Releases page, or install from source via `go install`.
