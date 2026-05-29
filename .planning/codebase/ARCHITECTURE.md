# on-a-meet — Architecture

## 1. Overall Architecture Pattern

**Monolithic CLI tool, single statically-linked binary.**

The entire application is a single Go binary with no runtime dependencies, plugins, or external processes. Cross-compilation with `CGO_ENABLED=0` produces platform-specific binaries that embed all functionality (V4L2 syscalls, lsof invocation, service management, TUI forms).

- **Module path:** `github.com/sergiocarracedo/on-a-meet`
- **Build target:** Single `main` package at root, everything else is either `cmd/` or `internal/`.
- **Go version:** 1.25.0 (go.mod)
- **Zero CGO:** `CGO_ENABLED=0` in goreleaser builds.

---

## 2. Entry Points

### `main.go` (`/works/opensource/on-a-meet/main.go`)

Minimal entry point:

```go
package main

import "github.com/sergiocarracedo/on-a-meet/cmd"

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    cmd.Execute()
}
```

- Calls `cmd.Execute()` which delegates to Cobra.
- Version/commit/date injected at build time via `ldflags` in `.goreleaser.yaml`:
  ```
  -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
  ```

### `cmd/` Package (`/works/opensource/on-a-meet/cmd/`)

All commands live in `package cmd`. Cobra drives the command tree:

```
rootCmd (on-a-meet)
├── detect          — Continuous camera monitoring + command execution
├── list            — One-shot device enumeration with status table
├── onboard         — Interactive TUI setup wizard
└── service (svc)   — Service management parent
    ├── install     — Install systemd/launchd unit
    ├── uninstall   — Stop + remove unit
    ├── start       — Start service
    ├── stop        — Stop service
    └── restart     — Stop + start (reload config)
```

Each file in `cmd/` follows the Cobra pattern:

```go
var commandNameCmd = &cobra.Command{
    Use:   "name",
    Short: "...",
    Long:  "...",
    RunE:  func(cmd *cobra.Command, args []string) error { ... },
}

func init() {
    parentCmd.AddCommand(subCmd)
    subCmd.Flags().StringVarP(...)
    viper.BindPFlag("config-key", subCmd.Flags().Lookup("flag"))
}
```

---

## 3. Module Structure

```
main.go                         # Binary entry point
cmd/                            # Cobra command definitions (import internal/ packages)
internal/
  config/                       # Config struct + defaults
  detector/                     # Detector interface + implementations (V4L2, lsof)
  engine/                       # Polling loop state machine
  executor/                     # Command dispatch with template rendering
  output/                       # pterm wrapper + helpers
```

### `internal/config/` (`/works/opensource/on-a-meet/internal/config/config.go`)

- **`Config` struct** — typed Go struct with `mapstructure` tags for Viper deserialization.
- **`Defaults()` function** — returns a `Config` with sensible defaults.
- Used by `cmd/root.go` where Viper sets defaults and reads from YAML/env/flags.
- The `detect` command uses a local `detectConfig` struct in `cmd/detect.go` via `configFromViper()`.

### `internal/detector/` (`/works/opensource/on-a-meet/internal/detector/`)

Core abstraction. Contains:

- **`interface.go`** — The `Detector` interface and shared types.
- **`detector.go`** — Factory function `New(method)` dispatching by name string.
- **`v4l2_linux.go`** — Linux-only V4L2 implementation using raw `ioctl` syscalls.
- **`v4l2_stub.go`** — Non-Linux stub returning errors.
- **`lsof_linux.go`** — Linux-only lsof implementation using `os/exec`.
- **`lsof_stub.go`** — Non-Linux stub returning errors.
- **`detector_test.go`** — Factory tests (interface compliance checks).
- **`v4l2_stub_test.go`** — Stub behavior verification (build-tag gated).

### `internal/engine/` (`/works/opensource/on-a-meet/internal/engine/engine.go`)

- **`Engine` struct** — polling loop with state machine, debounce, hotplug detection, camera filtering.
- **`Option` functional options pattern** — `WithInterval`, `WithDebounce`, `WithCameraFilter`, `WithOnChange`.
- **`OnChange` callback type** — `func(path string, oldState, newState bool, info detector.DeviceInfo)`.
- **`deviceState`** — per-device internal state: current/previous boolean, debounce counter.

### `internal/executor/` (`/works/opensource/on-a-meet/internal/executor/executor.go`)

- **`Executor` struct** — manages command execution with timeout, overlap prevention, env file loading.
- **`TemplateData`** — struct with `CameraID`, `Device`, `State` fields for `text/template` rendering.
- **`sync.Map`** keyed by state ("on"/"off") to prevent concurrent execution of the same state command.
- **`SetEnvFile()`** — loads a key=value environment file for subprocess.

### `internal/output/` (`/works/opensource/on-a-meet/internal/output/output.go`)

- Wraps `pterm` styled loggers: `Info`, `Success`, `Warning`, `Error`, `Debug`.
- **`Init(silent, verbose bool)`** — configures pterm output.
- **`Table()`** — renders pterm table with header.
- **`Banner()`** — startup banner with device count.
- **`RedactSecrets()`** — regex-based JWT token redaction in output.

---

## 4. Key Abstractions

### Detector Interface (`/works/opensource/on-a-meet/internal/detector/interface.go`)

```go
type Detector interface {
    ListDevices() ([]DeviceInfo, error)
    Detect(devicePath string) (DeviceStatus, error)
}
```

- `DeviceInfo`: `Path`, `Driver`, `Card`, `Bus` — metadata about a camera device.
- `DeviceStatus`: `On bool`, `CheckedAt time.Time` — point-in-time state.
- Two implementations: **V4L2Detector** (direct `ioctl` syscall) and **LsofDetector** (runs `lsof`).
- Factory `New(method string)` in `detector.go` returns the appropriate implementation.

### Engine (`/works/opensource/on-a-meet/internal/engine/engine.go`)

```go
type Engine struct {
    detector     detector.Detector
    interval     time.Duration
    debounce     int
    cameraFilter string
    onChange     OnChange
    states       map[string]*deviceState
    mu           sync.Mutex
    logger       *log.Logger
}
```

- `Run(ctx context.Context) error` — infinite poll loop; returns when context is cancelled.
- Init phase: enumerates devices, records initial state, fires initial `onChange` callbacks.
- Poll cycle: re-enumerates (hotplug detection), checks each device, applies debounce logic, fires callbacks on transitions.
- Debounce: N consecutive same-state polls required before firing a transition. Configurable via `WithDebounce(n)`.

### Executor (`/works/opensource/on-a-meet/internal/executor/executor.go`)

```go
type Executor struct {
    timeout time.Duration
    running sync.Map
    envFile string
}
```

- `ExecOn(ctx, cmdStr, data)` / `ExecOff(ctx, cmdStr, data)` — render template, run `sh -c` subprocess.
- Overlap prevention: `sync.Map` stores cancel functions keyed by state name. If an "on" command is already running, a second `ExecOn` is silently skipped. Cross-state (on while off) is allowed.
- Timeout: context-based timeout via `context.WithTimeout`. Process group kill via `Setpgid` + `syscall.Kill(-pid, SIGKILL)`.
- Environment file: `SetEnvFile()` loads a key=value file; variables are available in both template expansion and subprocess environment.

### Config (`/works/opensource/on-a-meet/internal/config/config.go`)

```go
type Config struct {
    Camera          string `mapstructure:"camera"`
    Interval        string `mapstructure:"interval"`
    OnCommand       string `mapstructure:"on-command"`
    OffCommand      string `mapstructure:"off-command"`
    DetectMethod    string `mapstructure:"detect-method"`
    Debounce        int    `mapstructure:"debounce"`
    Timeout         string `mapstructure:"timeout"`
    Verbose         bool   `mapstructure:"verbose"`
    EnvironmentFile string `mapstructure:"environment-file"`
}
```

Viper configuration priority (highest to lowest):
1. CLI flags (`viper.BindPFlag`)
2. Environment variables (prefix `ON_A_MEET_`)
3. YAML config file (`~/.config/on-a-meet/config.yaml`, `/etc/on-a-meet/config.yaml`, `./config.yaml`)
4. Hard-coded defaults (`viper.SetDefault` in `cmd/root.go`)

---

## 5. Data Flow

```
main.go
  └─ cmd.Execute()
       └─ rootCmd.Execute()
            ├─ initConfig()              # Viper: read config, env, flags
            │    └─ output.Init()        # pterm silent/verbose mode
            │
            ├─ detectCmd.RunE()
            │    ├─ configFromViper()    # Build detectConfig from Viper
            │    ├─ detector.New()       # Create V4L2Detector or LsofDetector
            │    ├─ det.ListDevices()    # Enumerate cameras
            │    ├─ det.Detect()         # Initial status display
            │    ├─ executor.New()       # Create Executor with timeout
            │    ├─ engine.New()         # Create Engine with detector + options
            │    │    └─ WithOnChange()  # Callback that fires Executor
            │    ├─ WithCameraFilter()   # Optional device filter
            │    ├─ ctx, cancel          # Signal handling (SIGINT/SIGTERM)
            │    └─ eng.Run(ctx)         # BLOCKING: poll loop
            │         ├─ ListDevices()   # [per cycle] re-enumerate
            │         ├─ Detect(path)    # [per device] check state
            │         ├─ Debounce logic  # N consecutive -> fire
            │         └─ onChange()      # Callback -> executor.ExecOn/ExecOff
            │              └─ go func()  # Async command dispatch
            │                   └─ exec.Command("sh", "-c", rendered)
            │
            ├─ listCmd.RunE()
            │    ├─ detector.New()
            │    ├─ det.ListDevices()
            │    ├─ det.Detect() per device
            │    └─ output.Table()       # pterm table render
            │
            ├─ onboardCmd.RunE()
            │    ├─ huh.Form             # Interactive TUI
            │    ├─ det.ListDevices()    # Show cameras
            │    ├─ det.Detect()         # Live test (ON then OFF)
            │    ├─ JSON marshal         # Write to /tmp
            │    └─ sudo re-exec         # Elevate, then --apply
            │         └─ YAML write      # Write /etc/on-a-meet/config.yaml
            │         └─ installService()# kardianos/service Install+Start
            │
            └─ service subcommands
                 ├─ installService()     # service.New -> Install() + Start()
                 ├─ svc.Stop()           # Stop service
                 ├─ svc.Uninstall()      # Remove unit
                 ├─ svc.Start()          # Start service
                 └─ svc.Restart()        # Stop + Start
```

### Poll Loop Detail (`engine.Run`)

```
Run(ctx)
  │
  ├── [Init Phase]
  │   ├── ListDevices()           → get available /dev/video* devices
  │   ├── Apply cameraFilter      → filter to target device if set
  │   ├── Detect(path) per device → initial state snapshot
  │   ├── Store in e.states map   → deviceState{current, previous, debounceCount}
  │   └── Fire onChange()         → emit initial ON/OFF (pretend transition from imaginary opposite)
  │
  └── [Loop]
      select {
        case <-ctx.Done():  return ctx.Err()
        case <-time.After(interval):
          pollCycle()
      }
```

### Poll Cycle Detail (`engine.pollCycle`)

```
pollCycle()
  │
  ├── ListDevices()               → current device list
  ├── Compare with e.states       → detect added/removed devices
  │   ├── Removed: fire onChange(path, true, true, info)  → "disconnected"
  │   └── Added:   fire onChange(path, false, false, info) → "detected"
  │
  └── For each remaining device:
      ├── Detect(path)            → get current status
      ├── Compare with deviceState.current
      │   ├── Same: reset debounceCount to 0
      │   └── Different: increment debounceCount
      │       └── debounceCount >= debounceTarget:
      │           ├── Update deviceState.current
      │           ├── Reset debounceCount
      │           └── Fire onChange(path, oldState, newState, info)
      │               └── In detect.go callback:
      │                   if newState: executor.ExecOn(ctx, onCmd, data)
      │                   if !newState: executor.ExecOff(ctx, offCmd, data)
      └── [goroutine per command]
```

---

## 6. State Management

### Per-Device State in Engine

The engine maintains a `map[string]*deviceState` protected by `sync.Mutex`:

```go
type deviceState struct {
    info           detector.DeviceInfo   // immutable device metadata
    current        bool                  // current debounced state (true=ON)
    previous       bool                  // previous state before last transition
    debounceCount  int                   // consecutive same-state observations
    debounceTarget int                   // threshold to fire (from config)
}
```

- Thread-safe: all map access is behind `e.mu.Lock()`.
- Updated in both `Run()` (init) and `pollCycle()` (steady state).
- Hotplug events add/remove entries from the map during `pollCycle()`.

### Overlap Prevention in Executor

```go
type Executor struct {
    timeout time.Duration
    running sync.Map        // key: "on"|"off", value: context.CancelFunc
    envFile string
}
```

- `sync.Map` stores a cancel function for each state ("on" and "off").
- On `ExecOn`/`ExecOff`: check `running.Load(state)` — if already present, skip (return nil).
- Store the cancel function before execution, delete after.
- Consequence: only one "on" command and one "off" command can run simultaneously. A new "on" while an "on" is running is skipped. A new "off" while an "on" is running _is_ allowed (cross-state).

---

## 7. Cross-Platform Strategy

Go build tags separate platform-specific implementations:

| File | Build Tag | Platform | Purpose |
|------|-----------|----------|---------|
| `v4l2_linux.go` | `//go:build linux` | Linux | V4L2 `ioctl` syscall via `golang.org/x/sys/unix` |
| `v4l2_stub.go` | `//go:build !linux` | Non-Linux | Returns error: "V4L2 detection is only supported on Linux" |
| `lsof_linux.go` | `//go:build linux` | Linux | `exec.Command("lsof", path)` — exit code 0 = ON |
| `lsof_stub.go` | `//go:build !linux` | Non-Linux | Returns error: "lsof detection is only supported on Linux" |
| `v4l2_stub_test.go` | `//go:build !linux` | Non-Linux | Tests that stub returns expected errors |

**Key details:**
- `goreleaser` cross-compiles for `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.
- On macOS, only lsof could potentially work (V4L2 is Linux-specific), but currently both backends return errors on non-Linux.
- The `Detector` interface itself is platform-agnostic (`interface.go` has no build tag).
- Stub types exist with identical signatures to satisfy compilation on all platforms.

---

## 8. Service Management Architecture

The `kardianos/service` library handles cross-platform service abstraction:

```
cmd/install.go
  ├── service.New(noopProgram{}, config{Name:"on-a-meet", Arguments:["detect","--config","/etc/on-a-meet/config.yaml"]})
  ├── svc.Install()   → Creates systemd unit on Linux, launchd plist on macOS
  ├── patchUnitEnvironmentFile()  → systemd-specific: patches EnvironmentFile path in unit
  ├── svc.Start()     → Enables and starts the unit
  └── Returns

cmd/uninstall.go
  ├── svc.Stop()      → Stops the service
  └── svc.Uninstall() → Removes the unit

cmd/restart.go
  ├── patchUnitEnvironmentFile()  → Re-apply env file path
  ├── svc.Stop()
  └── svc.Start()
```

The `noopProgram` struct satisfies `service.Interface` with empty Start/Stop methods since the detect command is self-contained (cobra RunE does all the work when the binary is launched by systemd).

---

## 9. Configuration Loading Flow

```
cmd/root.go initConfig()
  │
  ├── Determine config path:
  │   ├── --config flag → use explicitly
  │   └── default search:
  │       ├── ~/.config/on-a-meet/config.yaml
  │       ├── /etc/on-a-meet/config.yaml
  │       └── ./config.yaml (current directory)
  │
  ├── viper.SetConfigType("yaml")
  ├── viper.SetConfigName("config")
  │
  ├── viper.AutomaticEnv()
  ├── viper.SetEnvPrefix("ON_A_MEET")  → ON_A_MEET_INTERVAL, ON_A_MEET_ON_COMMAND, etc.
  │
  ├── viper.SetDefault(key, val) for all keys
  │
  ├── viper.ReadInConfig()
  │
  └── output.Init(silent, verbose)  → Must read config before output
```

Flag bindings in each `cmd/*.go` `init()` function wire CLI flags to Viper keys via `viper.BindPFlag()`.

---

## 10. Interactive Onboard Flow

```
onboardCmd.RunE()
  │
  ├── --config <file> mode:
  │   ├── Read JSON config file
  │   ├── Validate
  │   ├── --dry-run: print and exit
  │   └── Confirm → sudo re-exec with --apply
  │
  ├── --apply <file> mode (requires root):
  │   ├── Read JSON config
  │   ├── Marshal to YAML
  │   ├── Write /etc/on-a-meet/config.yaml
  │   ├── installService()
  │   └── Done
  │
  └── Interactive mode (no flags):
      ├── huh.NewForm with groups:
      │   ├── Camera selection (all vs specific)
      │   ├── Camera multi-select (if specific)
      │   ├── Debounce + Interval inputs
      │   └── On/Off command inputs
      ├── Live detection test:
      │   ├── "Turn camera ON → press Enter"
      │   ├── "Turn camera OFF → press Enter"
      │   ├── Retry or change method on failure
      ├── JSON marshal to /tmp
      ├── Confirm
      └── sudo re-exec --apply
```

The sudo re-exec pattern: the interactive part runs as the normal user (needs /dev/video access via `video` group), then escalates to root only for the write+install step.
