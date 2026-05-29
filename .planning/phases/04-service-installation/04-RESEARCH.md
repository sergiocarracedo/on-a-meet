# Phase 4: Service Installation — Research

**Gathered:** 2026-05-28

## Don't Hand-Roll

- **kardianos/service** (v1.2.4) handles all OS-specific service management (systemd on Linux, launchd on macOS, etc.). Do NOT write systemd unit templates manually — the library generates them from a Config struct.
- **Sudo detection:** `os.Geteuid() != 0` is standard. Do NOT implement auto-elevation (too complex, fragile across platforms).
- **Binary path:** `os.Executable()` from the standard library (Go 1.8+) is sufficient. No need for kardianos/osext.

## Common Pitfalls

- **Run() vs one-shot API:** The typical kardianos/service example calls `s.Run()` which blocks in the service loop. For install/uninstall commands, use the one-shot API: `s.Install()`, `s.Start()`, `s.Stop()`, `s.Uninstall()` directly — without calling `Run()`.
- **Arguments in Config:** The `service.Config.Arguments` field appends to the ExecStart command line. The binary path is automatically detected by kardianos/service from the running binary.
- **Interface requirement:** `service.New()` requires a `service.Interface` implementation. For install/uninstall commands, provide a minimal no-op program struct — it won't be used since we never call `Run()`.
- **Permission error handling:** Install/Uninstall operations typically fail with "access denied" when not running as root. Check `os.Geteuid()` first, print clear instructions, exit.
- **Start() after Install():** kardianos/service Install() only creates the unit file. Must call Start() separately to start the service immediately.

## Existing Patterns in This Codebase

- Cobra subcommand stubs already exist at `cmd/install.go` and `cmd/uninstall.go` with placeholder messages.
- Root command structure follows standard cobra pattern with `rootCmd.AddCommand()` in init().
- No existing kardianos/service dependency in go.mod — will be added.

## Recommended Approach

1. **Single plan** (or 2 plans in same wave): implement install + uninstall
2. **Shared service config generation** — function that returns `*service.Config` with Name, DisplayName, Description, Arguments
3. **Install:** Validate root → generate config → `service.New(noopProgram, cfg)` → `s.Install()` → `s.Start()`
4. **Uninstall:** Validate root → generate config → `service.New(noopProgram, cfg)` → `s.Stop()` → `s.Uninstall()`
5. **Arguments include** all detect flags: `detect`, `--config`, `/etc/on-a-meet/config.yaml`, `--interval`, `--debounce`, `--on-command`, `--off-command`, `--timeout`, `--camera`, `--silent`
6. **Working directory:** `/` (root) — set via Config.WorkingDir field
7. **No status** subcommand — users use `systemctl status on-a-meet`
