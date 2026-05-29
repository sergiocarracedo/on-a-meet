# on-a-meet — Directory Layout and Organization

## Full Directory Tree

```
/works/opensource/on-a-meet/
├── main.go                                  # Binary entry point — calls cmd.Execute()
│
├── cmd/
│   ├── root.go                              # Root Cobra command, Viper config init, persistent flags
│   ├── detect.go                            # `detect` subcommand — continuous polling + command execution
│   ├── list.go                              # `list` subcommand — one-shot device table
│   ├── onboard.go                           # `onboard` subcommand — interactive TUI setup wizard
│   ├── service.go                           # `service` parent command (alias: `svc`)
│   ├── install.go                           # `service install` — systemd/launchd unit creation
│   ├── uninstall.go                         # `service uninstall` — stop + remove unit
│   ├── start.go                             # `service start` — start existing service
│   ├── stop.go                              # `service stop` — stop running service
│   └── restart.go                           # `service restart` — stop + start to reload config
│
├── internal/
│   ├── config/
│   │   ├── config.go                        # Config struct with mapstructure tags, Defaults()
│   │   └── config_test.go                   # Tests for default config values
│   │
│   ├── detector/
│   │   ├── interface.go                     # Detector interface, DeviceInfo, DeviceStatus types
│   │   ├── detector.go                      # Factory: New(method) -> Detector
│   │   ├── v4l2_linux.go                    # V4L2Detector — ioctl-based (Linux only, build tag)
│   │   ├── v4l2_stub.go                     # V4L2 stub — returns errors on non-Linux
│   │   ├── lsof_linux.go                    # LsofDetector — exec.Command("lsof") (Linux only)
│   │   ├── lsof_stub.go                     # Lsof stub — returns errors on non-Linux
│   │   ├── detector_test.go                 # Factory tests (compile-time interface checks)
│   │   └── v4l2_stub_test.go                # Stub behavior tests (non-Linux build tag)
│   │
│   ├── engine/
│   │   ├── engine.go                        # Engine — polling loop, debounce, hotplug, filtering
│   │   └── engine_test.go                   # Engine tests: startup, debounce, hotplug, filter, shutdown
│   │
│   ├── executor/
│   │   ├── executor.go                      # Executor — template rendering, subprocess, overlap prevention
│   │   └── executor_test.go                 # Executor tests: success, timeout, overlap, cross-state
│   │
│   └── output/
│       ├── output.go                        # pterm wrappers: Info, Success, Warning, Error, Debug, RedactSecrets
│       └── output_test.go                   # Output tests: Init silent, Init verbose
│
├── config.yaml.example                      # Example user-facing YAML config with comments
├── go.mod                                   # Go module definition (github.com/sergiocarracedo/on-a-meet)
├── go.sum                                   # Go module checksums
├── README.md                                # User-facing documentation: install, usage, examples
├── AGENTS.md                                # Agent guide (phase tracking, project summary, commands)
├── .gitignore                               # Ignore: .opencode/, on-a-meet binary
├── .goreleaser.yaml                         # GoReleaser config: cross-compile, archive, checksum
│
├── .github/
│   └── workflows/
│       └── release.yml                      # GitHub Actions: release on tag push (v*)
│
├── .planning/                               # Planning artifacts (not part of the binary)
│   ├── codebase/                            # THIS FILE lives here
│   │   ├── ARCHITECTURE.md
│   │   └── STRUCTURE.md
│   ├── STATE.md                             # Session tracking
│   ├── PROJECT.md                           # Project scope
│   ├── ROADMAP.md                           # Archived roadmap
│   ├── milestones/                          # Archived requirement/roadmap docs
│   ├── research/                            # Pre-development research
│   ├── phases/                              # Per-phase plans, summaries, discussion logs
│   └── quick/                               # Quick-fix plans and summaries
│
└── packages/                                # Empty directory (placeholder for future packaging)
```

---

## File-by-File Purpose

### Root

| File | Purpose |
|------|---------|
| `main.go` | Minimal entry point. Calls `cmd.Execute()`. Declares `version`, `commit`, `date` vars injected by goreleaser ldflags. |
| `go.mod` | Module definition. Direct deps: cobra, viper, pterm, kardianos/service, charmbracelet/huh, golang.org/x/sys. |
| `go.sum` | Auto-generated dependency checksums. |
| `config.yaml.example` | Documented example config with all keys and explanations. Users copy this to `~/.config/on-a-meet/config.yaml`. |
| `README.md` | User documentation: install (binary/source), permissions (video group), config, usage, templates, service mgmt. |
| `AGENTS.md` | Agent/LLM guide: current phase, project summary, tech stack, architecture overview, key files, example commands. |
| `.goreleaser.yaml` | Build config: `CGO_ENABLED=0`, GOOS=linux+darwin, GOARCH=amd64+arm64, tar.gz archives, checksums. |
| `.gitignore` | Ignore local binary and `.opencode/` directory. |
| `on-a-meet` | Built binary (gitignored). |

### `cmd/` — Cobra Commands

| File | Purpose |
|------|---------|
| `root.go` | Creates `rootCmd`, sets persistent flags (`--config`, `--silent`, `--verbose`), defines `initConfig()` for Viper setup (config paths, env prefix, defaults), and `Execute()` entry point. |
| `detect.go` | The main subcommand. Creates detector, executor, engine; wires OnChange callback to executor; handles SIGINT/SIGTERM; defines `detectConfig` struct and `configFromViper()` helper. |
| `list.go` | One-shot command: creates detector, enumerates devices, detects status per device, renders pterm table with Path/Driver/Card/Bus/Status columns. |
| `onboard.go` | Interactive TUI wizard using charmbracelet/huh. Three modes: (1) full interactive form, (2) `--config <file>` to skip wizard, (3) `--apply <file>` (root) to write config + install service. Implements sudo re-exec pattern. Contains `onboardConfig`, `writeConfig`, `yamlQuotedString`, `cameraList` types. |
| `service.go` | Parent command for service subcommands. Adds `service` and `svc` alias to root. No RunE — just a grouping node. |
| `install.go` | Root-gated install. Creates `noopProgram` (satisfies `service.Interface`), calls `service.New`, `Install()`, `Start()`. Also patches systemd unit EnvironmentFile path if set. Exports `serviceConfig()` and `installService()` for reuse by onboard. |
| `uninstall.go` | Root-gated uninstall. Calls `svc.Stop()` then `svc.Uninstall()`. |
| `start.go` | Root-gated start. Calls `svc.Start()`. |
| `stop.go` | Root-gated stop. Calls `svc.Stop()`. |
| `restart.go` | Root-gated restart. Calls `svc.Stop()` then `svc.Start()`. Also patches EnvironmentFile before restart. |

### `internal/config/`

| File | Purpose |
|------|---------|
| `config.go` | `Config` struct with 9 fields (`Camera`, `Interval`, `OnCommand`, `OffCommand`, `DetectMethod`, `Debounce`, `Timeout`, `Verbose`, `EnvironmentFile`), all tagged with `mapstructure`. `Defaults()` function returns sensible defaults. |
| `config_test.go` | One test: `TestDefaults` verifies interval, detect-method, and debounce defaults match expected values. |

### `internal/detector/`

| File | Purpose |
|------|---------|
| `interface.go` | Defines `DeviceStatus` (On bool, CheckedAt time.Time), `DeviceInfo` (Path, Driver, Card, Bus string), and `Detector` interface (ListDevices, Detect). No build tags — always compiled. |
| `detector.go` | Factory function `New(method string) (Detector, error)` that dispatches `"v4l2"` or `"lsof"` to the respective constructor. Returns error for unknown methods. No build tags. |
| `v4l2_linux.go` | **Build tag: `linux`**. `V4L2Detector` struct. `ListDevices()` globs `/dev/video*`, opens each via `unix.Open`, issues `VIDIOC_QUERYCAP` ioctl, filters for `V4L2_CAP_VIDEO_CAPTURE`. `Detect()` tries to open device: `EBUSY` = ON, success = check `/proc/*/fd/` for open handles. Includes `_IOR` macro, `v4l2Capability` struct, `hasOpenHandle()` helper. |
| `v4l2_stub.go` | **Build tag: `!linux`**. Same struct and constructor signatures, but `ListDevices()` and `Detect()` return "only supported on Linux" errors. |
| `lsof_linux.go` | **Build tag: `linux`**. `LsofDetector` struct. `ListDevices()` reuses the same `openAndQueryCap()` from v4l2 for metadata (shared in same package). `Detect()` runs `exec.Command("lsof", devicePath)`; exit code 0 = ON, exit code 1 = OFF. |
| `lsof_stub.go` | **Build tag: `!linux`**. Same pattern as v4l2 stub — returns "only supported on Linux" errors. |
| `detector_test.go` | Tests `New()` factory for all three cases (v4l2, lsof, unknown). Includes compile-time interface compliance checks (`var _ Detector = d`). |
| `v4l2_stub_test.go` | **Build tag: `!linux`**. Verifies that V4L2 stub returns expected errors and zero-value statuses on non-Linux platforms. |

### `internal/engine/`

| File | Purpose |
|------|---------|
| `engine.go` | `Engine` struct with detector, interval, debounce, cameraFilter, onChange callback, states map, mutex, logger. `Option` type and functional options: `WithInterval`, `WithDebounce`, `WithCameraFilter`, `WithOnChange`. `New()` constructor with defaults. `Run(ctx)` — init phase (list devices, initial detect, fire initial callbacks) + poll loop (`select` on ctx.Done or ticker). `pollCycle()` — re-list, hotplug detection (add/remove), per-device detect with debounce logic. |
| `engine_test.go` | Five tests: `TestEngineStartup` (verifies initial onChange fires), `TestDebounce` (verifies debounce count delay), `TestHotplugAdd` (verifies new device detection), `TestCameraFilter` (verifies only filtered device is polled), `TestGracefulShutdown` (verifies ctx cancellation). Uses `mockDetector` and `fireTracker` test helpers. |

### `internal/executor/`

| File | Purpose |
|------|---------|
| `executor.go` | `TemplateData` struct (CameraID, Device, State string). `Executor` struct with timeout, running sync.Map, envFile. `New(timeout)`, `SetEnvFile(path)`. `ExecOn()`/`ExecOff()` call private `exec()`. `exec()`: overlap check via `sync.Map.Load`, template parse+execute via `text/template`, env file parse (`parseEnvFile` with `export ` support), `os.Expand` for shell variable substitution, `exec.CommandContext("sh", "-c", rendered)` with `Setpgid` for process group kill, stdout/stderr capture, exit code reporting. |
| `executor_test.go` | Five tests: `TestExecOnSuccess`, `TestExecOffSuccess`, `TestTemplateSubstitution`, `TestExecTimeout` (verifies timeout kills process), `TestSameStateSkip` (verifies overlap prevention), `TestCrossStateAllow` (verifies on+off can run concurrently). |

### `internal/output/`

| File | Purpose |
|------|---------|
| `output.go` | Package-level vars aliasing pterm loggers: `Info`, `Success`, `Warning`, `Error`, `Debug`. `Init(silent, verbose)` configures pterm output (discard or debug). `Table()` renders pterm table. `Banner()` prints startup banner. `RedactSecrets()` uses regex to redact JWT tokens (pattern: `ey...`). |
| `output_test.go` | Two minimal tests: `TestInitSilent`, `TestInitVerbose`. Smoke tests that init doesn't panic. |

### CI/CD

| File | Purpose |
|------|---------|
| `.github/workflows/release.yml` | GitHub Actions workflow: triggers on `v*` tag push or manual dispatch. Checks out code, sets up Go 1.22, runs `goreleaser release --clean`. Requires `contents: write` permission. |

---

## File Naming Conventions

| Convention | Meaning | Examples |
|------------|---------|---------|
| `*_linux.go` | Linux-only implementation behind `//go:build linux` | `v4l2_linux.go`, `lsof_linux.go` |
| `*_stub.go` | Platform fallback behind `//go:build !linux` | `v4l2_stub.go`, `lsof_stub.go` |
| `*_test.go` | Test files alongside source | `engine_test.go`, `executor_test.go`, `detector_test.go`, `config_test.go`, `output_test.go`, `v4l2_stub_test.go` |
| `*.example` | Example/documentation files | `config.yaml.example` |
| `*.yaml` | Configuration (GoReleaser, CI, example) | `.goreleaser.yaml`, `config.yaml.example` |
| `*.yml` | GitHub Actions workflow | `.github/workflows/release.yml` |
| `*.md` | Documentation | `README.md`, `AGENTS.md`, planning files |

---

## Where to Find Specific Things

| What | Where |
|------|-------|
| **Config struct** | `/works/opensource/on-a-meet/internal/config/config.go` |
| **Config defaults** | `/works/opensource/on-a-meet/cmd/root.go` (viper.SetDefault) and `/works/opensource/on-a-meet/internal/config/config.go` (Defaults()) |
| **Config loading** | `/works/opensource/on-a-meet/cmd/root.go` — `initConfig()` function |
| **Detector interface** | `/works/opensource/on-a-meet/internal/detector/interface.go` |
| **Detector factory** | `/works/opensource/on-a-meet/internal/detector/detector.go` |
| **V4L2 implementation** | `/works/opensource/on-a-meet/internal/detector/v4l2_linux.go` |
| **lsof implementation** | `/works/opensource/on-a-meet/internal/detector/lsof_linux.go` |
| **Non-Linux stubs** | `/works/opensource/on-a-meet/internal/detector/v4l2_stub.go` and `lsof_stub.go` |
| **Engine (poll loop)** | `/works/opensource/on-a-meet/internal/engine/engine.go` |
| **Executor (command dispatch)** | `/works/opensource/on-a-meet/internal/executor/executor.go` |
| **Template variables** | `/works/opensource/on-a-meet/internal/executor/executor.go` — `TemplateData` struct |
| **Output/UI wrappers** | `/works/opensource/on-a-meet/internal/output/output.go` |
| **detect command** | `/works/opensource/on-a-meet/cmd/detect.go` |
| **list command** | `/works/opensource/on-a-meet/cmd/list.go` |
| **onboard wizard** | `/works/opensource/on-a-meet/cmd/onboard.go` |
| **Service install** | `/works/opensource/on-a-meet/cmd/install.go` |
| **Service uninstall** | `/works/opensource/on-a-meet/cmd/uninstall.go` |
| **Service start/stop/restart** | `/works/opensource/on-a-meet/cmd/start.go`, `stop.go`, `restart.go` |
| **Build configuration** | `/works/opensource/on-a-meet/.goreleaser.yaml` |
| **CI/CD release workflow** | `/works/opensource/on-a-meet/.github/workflows/release.yml` |
| **User docs** | `/works/opensource/on-a-meet/README.md` |
| **Example config** | `/works/opensource/on-a-meet/config.yaml.example` |
| **Dependencies** | `/works/opensource/on-a-meet/go.mod` |

---

## Test File Placement

All tests follow standard Go conventions: `*_test.go` files live alongside the source files they test, in the same package (white-box testing).

| Test File | What It Tests |
|-----------|---------------|
| `internal/config/config_test.go` | `Config.Defaults()` values |
| `internal/detector/detector_test.go` | Factory `New()` for v4l2, lsof, unknown; interface compliance |
| `internal/detector/v4l2_stub_test.go` | Stub behavior on non-Linux (build-tag gated) |
| `internal/engine/engine_test.go` | Startup initialization, debounce timing, hotplug add, camera filter, graceful shutdown |
| `internal/executor/executor_test.go` | Success execution, template substitution, timeout, same-state overlap prevention, cross-state concurrency |
| `internal/output/output_test.go` | Silent mode init, verbose mode init |

All test files use the `testing` standard library. No external test frameworks or assertion libraries. Mock implementations (e.g., `mockDetector` in `engine_test.go`, `fireTracker` in `engine_test.go`) are defined within the test files themselves.

---

## Dependency Graph (Import Hierarchy)

```
main.go
  └── cmd
        ├── cmd/root.go → internal/output
        ├── cmd/detect.go → internal/detector, internal/engine, internal/executor, internal/output
        ├── cmd/list.go → internal/detector, internal/output
        ├── cmd/onboard.go → internal/detector, internal/output
        ├── cmd/install.go → internal/output  (+ kardianos/service)
        ├── cmd/uninstall.go → internal/output (+ kardianos/service)
        ├── cmd/start.go → internal/output    (+ kardianos/service)
        ├── cmd/stop.go → internal/output     (+ kardianos/service)
        └── cmd/restart.go → internal/output  (+ kardianos/service)

internal/
  ├── internal/engine → internal/detector
  ├── internal/executor → internal/output
  └── internal/config   (standalone, no internal deps)
```

No circular dependencies. The `internal/` packages only depend on each other where explicitly noted above. External dependencies (cobra, viper, pterm, service, huh) are limited to `cmd/` and `internal/output`.

---

## Build Tags Summary

```
//go:build linux
  ├── internal/detector/v4l2_linux.go     — V4L2 ioctl-based detection
  └── internal/detector/lsof_linux.go     — lsof-based detection

//go:build !linux
  ├── internal/detector/v4l2_stub.go      — V4L2 stub (returns error)
  ├── internal/detector/lsof_stub.go      — lsof stub (returns error)
  └── internal/detector/v4l2_stub_test.go — Stub test
```

The `V4L2Detector` and `LsofDetector` types are declared once per platform pair with identical APIs but different implementations. The rest of the codebase (cmd/, engine/, executor/, output/, config/) is fully platform-agnostic.
