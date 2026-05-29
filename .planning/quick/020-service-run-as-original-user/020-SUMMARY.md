# Quick Task 020 Summary

**Task:** Service runs commands as original user
**Completed:** 2026-05-29

## What was done

The installed systemd service now runs as the original user (who ran `sudo on-a-meet install` or `onboard`) instead of root. Uses `os.Getenv("SUDO_USER")` and kardianos/service's `UserName` config field.

## Files changed

- `cmd/install.go`: `serviceConfig()` takes a user param, sets `UserName`; `installService()` reads `SUDO_USER` env var
- `cmd/uninstall.go`: Updated `serviceConfig("")` call

## Commit

`789003b`
