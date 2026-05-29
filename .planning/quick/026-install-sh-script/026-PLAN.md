# Quick Task 026: install.sh script + curl-pipe-install in README

## Tasks

### Task 1: Create install.sh

- **Files:** `install.sh`
- **Action:** Create a bash script that:
  1. Detects OS (linux/darwin) via `uname -s`
  2. Detects architecture (amd64/arm64) via `uname -m`
  3. Fetches latest release tag from GitHub API
  4. Downloads the matching tar.gz archive from GitHub releases
  5. Extracts the `on-a-meet` binary
  6. Installs to `/usr/local/bin/` (uses `sudo` if not root)
- **Verify:** `bash -n install.sh` (syntax check), `chmod +x install.sh`
- **Done:** `install.sh` exists, is executable, has clean syntax

### Task 2: Add curl-pipe-install command to README

- **Files:** `README.md`
- **Action:** Add a one-line curl-pipe-sudo-bash install command at the top of the Installation section, before the platform-specific sections
- **Verify:** README renders correctly, command uses the new install.sh URL
- **Done:** README has a quick install section with the curl command
