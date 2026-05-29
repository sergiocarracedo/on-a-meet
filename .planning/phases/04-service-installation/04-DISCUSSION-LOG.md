# Phase 4: Service Installation — Discussion Log

**Gathered:** 2026-05-28
**Mode:** standard

## Areas Discussed

### 1. kardianos/service Integration & Lifecycle

**Options considered:**
- Separate install/uninstall as cobra commands, service runs via detect (Recommended) ✓
- Embed Run() in the install command itself
- Dedicated `service` subcommand with run/install/uninstall

**Chosen:** Separate cobra commands using one-shot API (Install/Uninstall/Start/Stop). No Run() needed.

### 2. Service Arguments & Config

**Options considered:**
- Minimal: --config pointing to config file only
- All detect flags passed explicitly (Recommended) ✓
- Config file + selective overrides

**Chosen:** All detect flags explicit on ExecStart line. System config at /etc/on-a-meet/config.yaml.

### 3. Status/Lifecycle Commands

**Options considered:**
- Install: install + start + enable, Uninstall: stop + disable + remove (Recommended) ✓
- Install only (no auto-start)
- Add status subcommand too

**Chosen:** install does full setup, uninstall does full teardown. No status command.

### 4. Binary & Working Directory

**Options considered:**
- os.Executable() + /etc/on-a-meet config path (Recommended) ✓
- argv[0] + user config path
- Prompt user for binary path on install

**Chosen:** os.Executable() for binary path, /etc/on-a-meet/config.yaml for system config, / as working directory.

### 5. Sudo & Permission Handling

**Options considered:**
- Check at start, warn with instructions, let user re-run (Recommended) ✓
- Auto-reinvoke with sudo
- Never check — assume user handles it

**Chosen:** Check os.Geteuid() at start, print clear instructions, exit. No auto-elevation.

## Deferred Ideas
- Status subcommand — not in scope. Users use `systemctl status`.
- System installer packages (.deb, .rpm) — future consideration.
