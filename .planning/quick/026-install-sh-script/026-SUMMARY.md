# Quick Task 026 Summary

**Task:** Add install.sh script with OS detection + curl-pipe-install in README
**Completed:** 2026-05-29

## What was done
Created `install.sh` that detects OS/architecture, fetches the latest release from GitHub API, downloads the correct tar.gz archive, and installs the binary to `/usr/local/bin/`. Added a one-line curl-pipe-sudo-bash quick install command at the top of the README Installation section.

## Files changed
- `install.sh` (new): OS/arch detection, GitHub API version fetch, download+extract+install
- `README.md`: Added "Quick install" section with curl-pipe-command

## Commit
`32b284d`
