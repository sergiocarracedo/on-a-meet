# Phase 4: Service Installation - Context

**Gathered:** 2026-05-28
**Mode:** standard
**Status:** Ready for planning

<domain>
## Phase Boundary

Implement systemd/launchd service installation and removal using kardianos/service. install runs as cobra subcommand (already stubbed), installs+enables+starts the service. uninstall stops+disables+removes. The service unit itself re-runs `on-a-meet detect` with all detect flags explicitly passed. Sudo check at start with user-friendly instructions.

</domain>

<decisions>
## Implementation Decisions

### kardianos/service Integration & Lifecycle
- **Pattern:** Separate install/uninstall cobra commands using kardianos/service one-shot API (Install(), Uninstall(), Start(), Stop()). No Run() needed in these commands — the service unit re-runs `on-a-meet detect` with its config.
- **Library:** github.com/kardianos/service — not yet in go.mod, will be added during planning.

### Service Arguments
- **All detect flags passed explicitly:** The generated service unit includes all detect flags (--config, --interval, --debounce, --on-command, --off-command, --timeout, --camera, --silent) on the ExecStart line.
- **System config path:** /etc/on-a-meet/config.yaml for the --config flag.

### Lifecycle Commands
- **install:** Creates unit file + enables + starts the service.
- **uninstall:** Stops + disables + removes the service unit.
- **No status subcommand:** Not in scope. Users can use `systemctl status on-a-meet`.

### Binary & Working Directory
- **Binary path:** os.Executable() — finds the running binary at install time for the service unit.
- **Config path:** System-level at /etc/on-a-meet/config.yaml.
- **Working directory:** / (root).
- **Binary during uninstall:** The binary is running, so uninstall stops the service first, then removes unit. The binary itself stays on disk (user manages removal).

### Sudo Handling
- **Check at start, warn, exit:** If not root (os.Geteuid() != 0), print clear message: "Please re-run with sudo: sudo on-a-meet install". Exit with non-zero code.
- **No auto-elevation:** Too complex and fragile. Users handle sudo themselves.

</decisions>

<specifics>
## Specific Ideas

- Service unit ExecStart should have explicit flags visible in `ps aux` output for debugging
- kardianos/service handles OS detection automatically (systemd on Linux, launchd on macOS)
- The library generates correct unit files; no manual systemd unit template needed

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — REQ-009 (service install)
- `.planning/research/ARCHITECTURE.md` — Service layer overview, component diagram
- `cmd/install.go` — Existing stub
- `cmd/uninstall.go` — Existing stub

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `cmd/install.go` — Cobra command stub with Use/Short/Long defined
- `cmd/uninstall.go` — Cobra command stub with Use/Short/Long defined
- `internal/config/config.go` — Config struct with all detect options
- `cmd/detect.go` — Existing detect command that service will re-invoke

### Established Patterns
- Package-per-concern under `internal/` — new package expected for service logic
- Cobra subcommand pattern with RunE returning error

### Integration Points
- `cmd/install.go` RunE → creates kardianos/service Service, calls Install()+Start()
- `cmd/uninstall.go` RunE → creates kardianos/service Service, calls Stop()+Uninstall()
- The service unit will point to the compiled binary running `on-a-meet detect`

</code_context>

<deferred>
## Deferred Ideas

- Status subcommand — not in scope for this phase. Users can use `systemctl status`.
- System installer packages (.deb, .rpm) — future consideration.

</deferred>

---
*Phase: 04-service-installation*
*Context gathered: 2026-05-28*
