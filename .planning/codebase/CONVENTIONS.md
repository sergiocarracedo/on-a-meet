# CONVENTIONS.md — Code Style and Patterns

## 1. Project Structure

```
./
├── main.go                    # Entry point with ldflags version injection
├── cmd/                       # Cobra command definitions (one file per command)
├── internal/                  # Non-reusable application logic
│   ├── config/                # Config struct, defaults, mapstructure tags
│   ├── detector/              # Detector interface + implementations + factory
│   ├── engine/                # Polling engine with state machine
│   ├── executor/              # Command execution with templates
│   └── output/                # pterm wrapper functions
├── .goreleaser.yaml           # GoReleaser cross-platform build config
├── config.yaml.example        # Example YAML config
├── go.mod / go.sum
└── packages/                  # Reserved for dist packages (currently empty)
```

**Key rule:** `cmd/` owns CLI surface (Cobra commands, flag binding, Viper wiring). `internal/` owns domain logic. No `internal/` package imports from `cmd/`.

**Test files** always live alongside the source they test (e.g., `engine.go` and `engine_test.go` in the same directory).

---

## 2. Package Naming and Grouping

| Package | Import Path | Responsibility |
|---------|-------------|----------------|
| `cmd` | `github.com/sergiocarracedo/on-a-meet/cmd` | All Cobra commands |
| `detector` | `.../internal/detector` | Interface, factory, backends |
| `engine` | `.../internal/engine` | Polling state machine |
| `executor` | `.../internal/executor` | Command execution with templates |
| `config` | `.../internal/config` | Config struct + defaults |
| `output` | `.../internal/output` | pterm wrappers |

**Source files in `cmd/` follow a flat file-per-command pattern:**
- `cmd/detect.go` — `detectCmd`
- `cmd/list.go` — `listCmd`
- `cmd/onboard.go` — `onboardCmd`
- `cmd/service.go` — `serviceCmd` (parent, adds `svc` alias)
- `cmd/install.go`, `cmd/uninstall.go`, `cmd/start.go`, `cmd/stop.go`, `cmd/restart.go` — subcommands of `serviceCmd`

---

## 3. Import Aliasing

- **No import aliases** for internal packages. Use the package name directly.
- **Standard library imports** grouped first, then third-party, then internal (separated by blank lines). See `cmd/detect.go` lines 4–17 for canonical grouping.

---

## 4. Interface Conventions

### 4.1 Single Package Interface

The `Detector` interface is defined in the same package as its implementations (`internal/detector/interface.go`):

```go
// internal/detector/interface.go
type Detector interface {
    ListDevices() ([]DeviceInfo, error)
    Detect(devicePath string) (DeviceStatus, error)
}
```

**Pattern:** Interface, domain types (`DeviceStatus`, `DeviceInfo`), factory (`New`), and all implementations live in the same package. Consumers depend on the interface, not on concrete types.

### 4.2 Compile-Time Interface Check

Used in tests to verify implementations satisfy the interface:

```go
// internal/detector/detector_test.go line 39
var _ Detector = d
```

**DON'T** put these checks in production code; they belong in tests only.

---

## 5. Constructor / Factory Patterns

### 5.1 Simple Constructors

Return concrete pointer types for non-interface consumers:

```go
// internal/detector/v4l2_linux.go line 38
func NewV4L2Detector() *V4L2Detector {
    return &V4L2Detector{}
}

// internal/detector/lsof_linux.go line 14
func NewLsofDetector() *LsofDetector {
    return &LsofDetector{}
}

// internal/executor/executor.go line 34
func New(timeout time.Duration) *Executor {
    return &Executor{timeout: timeout}
}
```

**Pattern:** Constructors named `New` or `New<Type>`. Struct zero-initialization where possible; required parameters as arguments; optional parameters via functional options.

### 5.2 Factory Function

Returns the interface, not the concrete type:

```go
// internal/detector/detector.go line 5
func New(method string) (Detector, error) {
    switch method {
    case "v4l2":
        return NewV4L2Detector(), nil
    case "lsof":
        return NewLsofDetector(), nil
    default:
        return nil, fmt.Errorf("unknown detection method: %q (supported: v4l2, lsof)", method)
    }
}
```

**Pattern:** A switch-based factory that returns the `Detector` interface. Error for unknown methods includes supported values.

### 5.3 Functional Options Pattern

Used in `engine.New` for optional configuration:

```go
// internal/engine/engine.go
type Option func(*Engine)

func WithInterval(d time.Duration) Option {
    return func(e *Engine) { e.interval = d }
}

func WithDebounce(n int) Option {
    return func(e *Engine) { e.debounce = n }
}

func WithCameraFilter(path string) Option {
    return func(e *Engine) { e.cameraFilter = path }
}

func WithOnChange(cb OnChange) Option {
    return func(e *Engine) { e.onChange = cb }
}

func New(det detector.Detector, opts ...Option) *Engine {
    e := &Engine{
        detector: det,
        interval: 1 * time.Second,      // sensible defaults
        debounce: 3,
        states:   make(map[string]*deviceState),
        logger:   log.New(log.Writer(), "", 0),
    }
    for _, opt := range opts {
        opt(e)                           // apply overrides
    }
    return e
}
```

**Pattern:**
- `type Option func(*Engine)` — function type
- `WithXxx(value) Option` — constructor-style option functions
- Constructor sets **sensible defaults** first, then applies options

**DON'T** expose exported fields for mutation; use options.

---

## 6. Error Handling

### 6.1 Error Wrapping with `%w`

```go
// internal/detector/v4l2_linux.go line 45
return nil, fmt.Errorf("enumerate video devices: %w", err)

// cmd/onboard.go line 84
return fmt.Errorf("failed to read config file: %w", err)

// internal/detector/detector.go line 12
return nil, fmt.Errorf("unknown detection method: %q (supported: v4l2, lsof)", method)
```

**Pattern:** Always use `fmt.Errorf("context message: %w", err)` to wrap errors. Context message is lowercase (no leading capital, no trailing period).

### 6.2 Sentinel Errors

```go
// internal/detector/v4l2_stub.go (non-Linux stubs)
return nil, errors.New("V4L2 detection is only supported on Linux")

// internal/detector/lsof_stub.go (non-Linux stubs)
return nil, errors.New("lsof detection is only supported on Linux")
```

**Pattern:** Use `errors.New` for static error messages on stub implementations. The error message is descriptive and includes "only supported on Linux" for platform-specific features.

### 6.3 Error Handling in Commands (RunE)

Commands use `RunE` (not `Run`) to return errors:

```go
// cmd/restart.go line 17
RunE: func(cmd *cobra.Command, args []string) error {
```

Root `Execute()` prints returned errors to stderr and exits with code 1:

```go
// cmd/root.go line 32
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

**Pattern:** Use `RunE` for proper error propagation. Non-error output uses `output.*`. Genuine errors are returned for the caller in `Execute()` to handle.

### 6.4 Error Type Assertion

```go
// internal/detector/lsof_linux.go line 45-48
if exitErr, ok := err.(*exec.ExitError); ok {
    if exitErr.ExitCode() == 1 {
        return DeviceStatus{On: false, CheckedAt: time.Now()}, nil
    }
}
```

---

## 7. Output Patterns (pterm Wrappers)

All user-facing output goes through `internal/output/output.go`:

```go
// Pre-configured loggers
var (
    Info    = pterm.Info
    Success = pterm.Success
    Warning = pterm.Warning
    Error   = pterm.Error
    Debug   = pterm.Debug
)
```

**Usage in commands:**
```go
output.Info.Println("Shutting down...")
output.Warning.Printfln("Failed to patch environment file path: %v", err)
output.Error.Println("Failed to enumerate camera devices:", err)
output.Success.Println("Service installed")
output.Debug.Printfln("%s-command: %s", state, output.RedactSecrets(rendered))
output.Table(rows)
output.Banner(len(devices))
output.Init(cfgSilent, cfgVerbose)
output.RedactSecrets(s string) string   // JWT redaction
```

**Patterns:**
- Use `Printfln(format, args...)` for formatted output
- Use `Println(args...)` for simple messages
- Do NOT use `fmt.Print*` or `pterm.*` directly in commands
- `RedactSecrets()` wraps strings to redact JWT tokens before logging
- `Init(silent, verbose)` controls output routing at startup
- `Table(data)` wraps pterm table rendering with error handling

---

## 8. Viper Config Pattern

### 8.1 Config Struct

```go
// internal/config/config.go
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

func Defaults() Config {
    return Config{
        Interval:     "1s",
        DetectMethod: "v4l2",
        Debounce:     3,
        Timeout:      "30s",
    }
}
```

**Pattern:** Config struct with `mapstructure` tags matching Viper/YAML keys. Separate `Defaults()` function.

### 8.2 Viper Initialization (cmd/root.go)

```go
func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "...")
    rootCmd.PersistentFlags().BoolVarP(&cfgSilent, "silent", "s", false, "...")
    rootCmd.PersistentFlags().BoolVarP(&cfgVerbose, "verbose", "V", false, "...")
    viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))
    viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
    // 1. If --config flag, use that file
    // 2. Else search: ~/.config/on-a-meet/, /etc/on-a-meet/, .
    // 3. viper.AutomaticEnv() with prefix "ON_A_MEET"
    // 4. viper.SetDefault(...) for all keys
    // 5. viper.ReadInConfig()
    // 6. output.Init(...)
}
```

**Merge order (lowest to highest priority):** YAML defaults < YAML config < env vars (ON_A_MEET_*) < CLI flags.

### 8.3 Per-Command Flag Binding

Each command's `init()` binds flags to Viper:

```go
// cmd/detect.go init()
viper.BindPFlag("camera", detectCmd.Flags().Lookup("camera"))
viper.BindPFlag("interval", detectCmd.Flags().Lookup("interval"))
viper.BindPFlag("detect-method", detectCmd.Flags().Lookup("detect"))
```

### 8.4 Local Config Read Pattern

The `configFromViper()` helper (cmd/detect.go) reads config for a specific command, with the detect-method flag-aware fallback:

```go
func configFromViper(detectChanged bool) detectConfig {
    method := detectMethod
    if !detectChanged {
        method = viper.GetString("detect-method")
    }
    return detectConfig{
        Camera:   viper.GetString("camera"),
        Interval: viper.GetString("interval"),
        // ...
    }
}
```

---

## 9. Engine Polling / State Machine Patterns

### 9.1 Callback Type

```go
// internal/engine/engine.go line 12
type OnChange func(path string, oldState, newState bool, info detector.DeviceInfo)
```

### 9.2 State Machine

- `deviceState` tracks `current`, `previous`, `debounceCount`, `debounceTarget`
- On each poll: if `status.On == s.current` -> reset debounce counter; else increment
- When `debounceCount >= debounceTarget` -> fire `onChange` callback

### 9.3 Hotplug

- `pollCycle` re-lists devices each cycle
- Devices in `states` but not in current list -> removed (fires onChange with oldState=newState=true)
- Devices in current list but not in `states` -> added (fires onChange with oldState=newState=false)

### 9.4 Thread Safety

```go
// internal/engine/engine.go
type Engine struct {
    // ...
    mu sync.Mutex
}
```

All access to `e.states` is protected by `e.mu.Lock()` / `e.mu.Unlock()`.

### 9.5 Context Cancellation

```go
// internal/engine/engine.go line 108-116
for {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(e.interval):
        e.pollCycle()
    }
}
```

---

## 10. Executor / Template Rendering

### 10.1 Template Data

```go
// internal/executor/executor.go line 18
type TemplateData struct {
    CameraID string  // e.g. "video0"
    Device   string  // e.g. "/dev/video0"
    State    string  // "on" or "off"
}
```

### 10.2 Template Rendering Pipeline

```go
// 1. Parse Go template
tmpl, err := template.New("").Parse(cmdStr)

// 2. Execute with TemplateData
tmpl.Execute(&buf, data)

// 3. Shell environment variable expansion
rendered = os.Expand(rendered, func(key string) string { ... })
```

### 10.3 Command Execution

```go
cmd := exec.CommandContext(cmdCtx, "sh", "-c", rendered)
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}   // process group isolation
```

### 10.4 Overlap Prevention

```go
// Track running commands by state ("on"/"off")
e.running sync.Map

// If a command for the same state is already running, skip
if _, loaded := e.running.Load(state); loaded {
    return nil    // silent skip
}
```

**Key:** on-commands and off-commands can run concurrently; two on-commands cannot overlap.

---

## 11. Platform Build Tags

Linux implementations use `//go:build linux` at the top of the file:

```go
//go:build linux
package detector
```

Stubs for non-Linux use `//go:build !linux`:

```go
//go:build !linux
package detector
```

**Pattern:** Non-Linux stubs return `errors.New("... only supported on Linux")` for every method.

---

## 12. Naming Conventions

| Convention | Example | Location |
|------------|---------|----------|
| Interface | `Detector` | `interface.go` |
| Exported types | `Engine`, `Executor`, `DeviceStatus`, `DeviceInfo`, `TemplateData` | Various |
| Unexported types | `deviceState`, `initDevice`, `onboardConfig`, `writeConfig`, `detectConfig`, `noopProgram`, `fireCall`, `fireTracker`, `mockDetector`, `mockDetectResponse` | Various |
| Unexported helpers | `configFromViper`, `patchUnitEnvironmentFile`, `openAndQueryCap`, `nullTerminatedString`, `hasOpenHandle`, `parseEnvFile`, `installService`, `serviceConfig` | Various |
| Unexported globals | `detectCamera`, `detectInterval`, `onboardDryRun`, `onboardApply`, `listMethod` | `cmd/` files |
| Exported but internal | `Info`, `Success`, `Warning`, `Error`, `Debug` | `output.go` |

**Prefix conventions in `cmd/`:**
- Command-local flag vars prefixed with command name: `detectCamera`, `detectInterval`, `listMethod`, `onboardDryRun`
- Command-local config structs: `detectConfig`, `onboardConfig`

---

## 13. What NOT to Do

1. **Do NOT import `cmd/` from `internal/`** -- this creates a circular dependency. `cmd/` imports `internal/`, never the reverse.

2. **Do NOT use `Run` instead of `RunE`** in Cobra commands -- all commands use `RunE` to return errors properly.

3. **Do NOT call `fmt.Print*` in internal packages** -- use `output.*` wrappers.

4. **Do NOT access Viper directly from `internal/`** -- config is read in `cmd/` and passed explicitly (e.g., `configFromViper` returns a struct, `executor.New(timeout)` takes a `time.Duration`).

5. **Do NOT use exported struct fields for optional configuration** -- use the functional options pattern (`WithInterval`, `WithDebounce`).

6. **Do NOT use `os.Exit(1)` in `internal/`** -- return errors to the Cobra command handler. Exit code 1 is the responsibility of `cmd.Execute()`.

7. **Do NOT return a concrete type from the factory** -- `detector.New()` returns `(Detector, error)`, not `(*V4L2Detector, error)`.

8. **Do NOT forget build tags for platform-specific files** -- `//go:build linux` / `//go:build !linux`.

9. **Do NOT export the `Option` type from a package that does not own options** -- `type Option func(*Engine)` lives in `package engine`, not in a shared package.

10. **Do NOT ignore `Sync.Map` return values** -- always check `loaded` from `Load`/`Store`.

11. **Do NOT use naked `go func()` without capturing closure variables** -- always pass as arguments (see `cmd/detect.go` lines 100-126 where `path`, `cfg`, `data` are captured correctly).

12. **Do NOT use `os.Exit(1)` in `onboard.go` for user-cancelled operations** -- use `return nil` to exit gracefully.
