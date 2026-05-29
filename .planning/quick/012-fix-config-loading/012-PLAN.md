# Quick Task 012 Plan: Fix config loading from onboard

## Tasks

### Task 1: Add /etc/on-a-meet as config search path

**Files:**
- `cmd/root.go`

**Action:**
Add `/etc/on-a-meet` as an additional viper config search path in `initConfig()`. Order: user home config first (user override), then system config from onboard, then current dir.

Change the config path block from:
```go
viper.AddConfigPath(home + "/.config/on-a-meet")
viper.AddConfigPath(".")
```
to:
```go
viper.AddConfigPath(home + "/.config/on-a-meet")
viper.AddConfigPath("/etc/on-a-meet")
viper.AddConfigPath(".")
```

### Task 2: Add verbose/silent viper defaults

**Files:**
- `cmd/root.go`

**Action:**
Add `silent` and `verbose` to the viper defaults block so they're available via `viper.GetBool()`.

After the existing defaults:
```go
viper.SetDefault("silent", false)
viper.SetDefault("verbose", false)
```

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All tests pass

**Done:**
- Config at `/etc/on-a-meet/config.yaml` (from onboard) is found by `detect` command
- silent/verbose defaults available in viper
