# Quick Task 020: Service runs commands as original user

## Task

**DESCRIPTION:** Ensure on-command/off-command execute as the user who ran `onboard`, not as root.

### Files
- `cmd/install.go`
- `cmd/uninstall.go`

### Action
1. Modify `serviceConfig()` to accept a `user` string parameter and set `cfg.UserName`
2. In `installService()`, read `SUDO_USER` env var and pass to `serviceConfig()`
3. Update `uninstall.go` to pass `""` to `serviceConfig("")`

### Verify
```bash
go build ./... && go vet ./... && go test ./...
```
When installed via `sudo on-a-meet install`, the systemd unit has `User=<original_user>`.
When installed via `onboard` (which re-execs with sudo), same result.

### Done
Service runs detect + commands as the user who triggered the install.
