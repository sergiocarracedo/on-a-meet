# STACK.md — Technologies and Dependencies Overview

## Language and Go Version

- **Language:** Go
- **Go version required:** `go 1.25.0` (as declared in `go.mod` line 3)
- **Module path:** `github.com/sergiocarracedo/on-a-meet` (`go.mod` line 1)

The Go version in use tracks the latest Go toolchain. The GitHub Actions release workflow (`/.github/workflows/release.yml` line 21) pins `go-version: "1.22"` for build reproducibility, but the actual module directive has been bumped to 1.25.0.

---

## Direct Dependencies (from `go.mod`)

| Dependency | Version | Source File | Purpose |
|---|---|---|---|
| `github.com/spf13/cobra` | `v1.10.2` | `go.mod:8` | CLI framework — root command, subcommands, persistent flags, help generation. Used in `cmd/root.go`, `cmd/detect.go`, `cmd/list.go`, `cmd/onboard.go`, `cmd/service.go`, `cmd/install.go`, `cmd/uninstall.go`, `cmd/start.go`, `cmd/stop.go`, `cmd/restart.go`. |
| `github.com/spf13/viper` | `v1.21.0` | `go.mod:9` | Configuration management — merges YAML config files, CLI flags, environment variables. Config lookup paths: `~/.config/on-a-meet/`, `/etc/on-a-meet/`, `.` (current dir). Environment prefix: `ON_A_MEET`. See `cmd/root.go:49-81`. |
| `github.com/pterm/pterm` | `v0.12.83` | `go.mod:7` | Terminal output — Info, Success, Warning, Error, Debug loggers and table rendering. Wrapped in `internal/output/output.go`. |
| `github.com/kardianos/service` | `v1.2.4` | `go.mod:6` | Cross-platform service management — Install(), Uninstall(), Start(), Stop() for systemd (Linux) and launchd (macOS). Used in `cmd/install.go`, `cmd/uninstall.go`, `cmd/start.go`, `cmd/stop.go`, `cmd/restart.go`. |
| `golang.org/x/sys` | `v0.45.0` | `go.mod:10` | Unix syscall wrappers — `unix.Open`, `unix.Close` for V4L2 device detection. Used in `internal/detector/v4l2_linux.go`. |

### Charmbracelet / Huh (indirect but essential)

| Dependency | Version | Source File | Purpose |
|---|---|---|---|
| `github.com/charmbracelet/huh` | `v1.0.0` | `go.mod:23` | Interactive TUI forms — Select, MultiSelect, Input, Confirm for the `onboard` wizard (`cmd/onboard.go`). Indirectly pulled via `go.mod:23`. |
| `github.com/charmbracelet/bubbletea` | `v1.3.6` | `go.mod:21` | TUI framework runtime — powers huh forms. Indirect. |
| `github.com/charmbracelet/bubbles` | `v0.21.1-0.20250623103423-23b8fd6302d7` | `go.mod:20` | Bubbletea components. Indirect. |
| `github.com/charmbracelet/lipgloss` | `v1.1.0` | `go.mod:24` | Terminal styling. Indirect. |

### Other Key Indirect Dependencies

| Dependency | Version | Source File | Purpose |
|---|---|---|---|
| `github.com/spf13/pflag` | `v1.0.10` | `go.mod:52` | Flag parsing (underlying cobra). |
| `gopkg.in/yaml.v3` | `v3.0.1` | `go.mod:59` | YAML marshaling — used in `cmd/onboard.go` for config writing. |
| `go.yaml.in/yaml/v3` | `v3.0.4` | `go.mod:55` | YAML v3 (alternate fork). |
| `github.com/fsnotify/fsnotify` | `v1.9.0` | `go.mod:33` | File system notifications (viper uses this for live config reload). |
| `github.com/subosito/gotenv` | `v1.6.0` | `go.mod:53` | Env file parsing (viper). |
| `github.com/pelletier/go-toml/v2` | `v2.2.4` | `go.mod:46` | TOML support (viper). |
| `github.com/spf13/afero` | `v1.15.0` | `go.mod:50` | Virtual filesystem abstraction (viper). |
| `github.com/spf13/cast` | `v1.10.0` | `go.mod:51` | Type casting helpers (viper). |
| `github.com/go-viper/mapstructure/v2` | `v2.4.0` | `go.mod:34` | Map-to-struct decoding (viper). |
| `github.com/mattn/go-isatty` | `v0.0.20` | `go.mod:39` | Terminal detection. |
| `github.com/mattn/go-runewidth` | `v0.0.20` | `go.mod:40` | Unicode string width. |
| `github.com/rivo/uniseg` | `v0.4.7` | `go.mod:47` | Unicode segmentation. |
| `golang.org/x/text` | `v0.34.0` | `go.mod:58` | Text processing. |
| `golang.org/x/term` | `v0.40.0` | `go.mod:57` | Terminal handling. |
| `golang.org/x/sync` | `v0.19.0` | `go.mod:56` | Concurrency primitives. |
| `atomicgo.dev/cursor` | `v0.2.0` | `go.mod:14` | Cursor control (pterm). |
| `atomicgo.dev/keyboard` | `v0.2.9` | `go.mod:15` | Keyboard input (pterm). |
| `atomicgo.dev/schedule` | `v0.1.0` | `go.mod:16` | Scheduling (pterm). |
| `github.com/gookit/color` | `v1.6.0` | `go.mod:35` | Terminal colors (pterm). |
| `github.com/xo/terminfo` | `v0.0.0-20220910002029-abceb7e1c41e` | `go.mod:54` | Terminal info queries. |
| `github.com/lithammer/fuzzysearch` | `v1.1.8` | `go.mod:37` | Fuzzy text matching. |
| `github.com/atotto/clipboard` | `v0.1.4` | `go.mod:17` | Clipboard access. |
| `github.com/catppuccin/go` | `v0.3.0` | `go.mod:19` | Catppuccin theme colors. |
| `github.com/aymanbagabas/go-osc52/v2` | `v2.0.1` | `go.mod:18` | OSC52 escape sequences. |
| `github.com/charmbracelet/x/ansi` | `v0.9.3` | `go.mod:25` | ANSI sequence handling. |
| `github.com/charmbracelet/x/cellbuf` | `v0.0.13` | `go.mod:26` | Cell buffer for terminal rendering. |
| `github.com/muesli/termenv` | `v0.16.0` | `go.mod:45` | Terminal environment. |
| `github.com/muesli/ansi` | `v0.0.0-20230316100256-276c6243b2f6` | `go.mod:43` | ANSI parsing. |
| `github.com/muesli/cancelreader` | `v0.2.2` | `go.mod:44` | Cancellable reader. |
| `github.com/erikgeiser/coninput` | `v0.0.0-20211004153227-1c3628e74d0f` | `go.mod:32` | Console input. |
| `github.com/containerd/console` | `v1.0.5` | `go.mod:30` | Console handling (pterm). |
| `github.com/mattn/go-localereader` | `v0.0.1` | `go.mod:41` | Locale-aware reader. |
| `github.com/lucasb-eyer/go-colorful` | `v1.2.0` | `go.mod:38` | Color space conversions. |
| `github.com/clipperhouse/uax29/v2` | `v2.7.0` | `go.mod:29` | Unicode text segmentation. |
| `github.com/dustin/go-humanize` | `v1.0.1` | `go.mod:31` | Human-readable formatting. |
| `github.com/mitchellh/hashstructure/v2` | `v2.0.2` | `go.mod:42` | Hash struct values. |
| `github.com/sagikazarmark/locafero` | `v0.11.0` | `go.mod:49` | File location (viper). |
| `github.com/sourcegraph/conc` | `v0.3.1-0.20240121214520-5f936abd7ae8` | `go.mod:50` | Structured concurrency. |
| `github.com/inconshreveable/mousetrap` | `v1.1.0` | `go.mod:36` | CLI trap detection (cobra). |

---

## Build System

### Local Build

```bash
go build -o on-a-meet .
```

Build tags: `netgo` (statically linked net, set in `.goreleaser.yaml` line 19).

The binary is compiled with `CGO_ENABLED=0` to produce fully static binaries (`.goreleaser.yaml` line 9).

### ldflags (version injection)

Three variables in `main.go` (lines 6-8) are injected at build time:

```
  -s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
```

Defined in `.goreleaser.yaml` lines 16-17. These are goreleaser template variables; `main.version` is also shadowed in `cmd/root.go:14` with a `"dev"` fallback.

### GoReleaser Configuration

**File:** `/.goreleaser.yaml`

- **Version:** 2 (goreleaser config API v2)
- **Before hooks:** `go mod tidy`
- **Builds:**
  - `CGO_ENABLED=0`
  - `goos`: `linux`, `darwin`
  - `goarch`: `amd64`, `arm64`
  - `tags`: `netgo`
  - `ldflags`: version/commit/date injection
- **Archives:** `tar.gz` format with name template `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}`
- **Checksum:** `checksums.txt`
- **Snapshot:** `{{ incpatch .Version }}-next`
- **Changelog:** ascending sort, filtering out `docs:` and `test:` commits

Note: Windows is **not** targeted. Only Linux and Darwin (macOS) builds are produced.

### Build Tags / Cross-Platform Considerations

The codebase uses Go build constraints for platform-specific code:

1. **`//go:build linux`** — `internal/detector/v4l2_linux.go` (lines 1, 155 lines): Full V4L2 implementation using `golang.org/x/sys/unix` for syscalls (`SYS_IOCTL`, `VIDIOC_QUERYCAP`), `/dev/video*` glob, `/proc` walk for open file handle detection.

2. **`//go:build linux`** — `internal/detector/lsof_linux.go` (line 1, 52 lines): lsof-based detection backend using `os/exec.Command("lsof", devicePath)`.

3. **`//go:build !linux`** — `internal/detector/v4l2_stub.go` (line 1, 19 lines): stub that returns errors like `"V4L2 detection is only supported on Linux"`. This ensures compilation on macOS but provides no functionality.

**Implication:** On macOS, the `detect` and `list` commands will fail with `"V4L2 detection is only supported on Linux"` when using the default `v4l2` method. The `lsof` method is also linux-only (no `lsof_darwin.go` equivalent exists). The service commands (`install`, `uninstall`, `start`, `stop`, `restart`) use `kardianos/service` which supports both systemd (Linux) and launchd (macOS) — so service management works cross-platform, but detection does not.

---

## Environment Variables

Defined in `cmd/root.go:64`:

```go
viper.AutomaticEnv()
viper.SetEnvPrefix("ON_A_MEET")
```

All Viper config keys are automatically mapped to environment variables with the prefix `ON_A_MEET_`. For example:

| Config Key | Env Variable |
|---|---|
| `camera` | `ON_A_MEET_CAMERA` |
| `interval` | `ON_A_MEET_INTERVAL` |
| `on-command` | `ON_A_MEET_ON_COMMAND` |
| `off-command` | `ON_A_MEET_OFF_COMMAND` |
| `detect-method` | `ON_A_MEET_DETECT_METHOD` |
| `debounce` | `ON_A_MEET_DEBOUNCE` |
| `timeout` | `ON_A_MEET_TIMEOUT` |
| `silent` | `ON_A_MEET_SILENT` |
| `verbose` | `ON_A_MEET_VERBOSE` |
| `environment-file` | `ON_A_MEET_ENVIRONMENT_FILE` |

Additionally, when executing user-defined commands, the `executor` package (`internal/executor/executor.go`) can optionally load variables from an environment file (specified by the `environment-file` config key). These are loaded via `parseEnvFile()` and merged into the command's environment. Variables available to executed commands also include `os.Environ()`.

---

## Key Configuration Files

| File | Purpose |
|---|---|
| `config.yaml.example` | Documented example config file (37 lines) with all available options and comments. Template for `~/.config/on-a-meet/config.yaml`. |
| `.goreleaser.yaml` | GoReleaser configuration for cross-platform release builds (37 lines). |
| `.github/workflows/release.yml` | GitHub Actions workflow for automated releases on tag push (26 lines). |
| `go.mod` | Go module definition with all dependencies (60 lines). |
| `go.sum` | Dependency checksums (213 lines). |
| `.gitignore` | Ignores `.opencode/` and `on-a-meet` binary. |

---

## CI/CD: GitHub Actions Workflows

**File:** `/.github/workflows/release.yml`

- **Trigger:** Push of tags matching `v*` (e.g., `v1.0.0`) OR `workflow_dispatch` (manual trigger).
- **Permissions:** `contents: write` (for creating GitHub releases and uploading artifacts).
- **Job:**
  - Runs on `ubuntu-latest`
  - Steps:
    1. `actions/checkout@v4` with `fetch-depth: 0` (full history for changelog).
    2. `actions/setup-go@v5` with `go-version: "1.22"`.
    3. `goreleaser/goreleaser-action@v7` with `args: release --clean`.
    4. `GITHUB_TOKEN` set to `${{ secrets.GITHUB_TOKEN }}` (automatic GitHub token for release creation).
- **Output:** GoReleaser creates cross-platform binaries (`linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`), archives them in `tar.gz`, generates checksums, and creates a GitHub Release with changelog.

---

## Entry Point

**File:** `/main.go`

```go
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    cmd.Execute()
}
```

These three variables are injected via ldflags at release build time. The `cmd.Execute()` call dispatches to the cobra root command defined in `cmd/root.go`.

---

## Project Directory Structure

```
/works/opensource/on-a-meet/
├── main.go                         # Entry point with ldflag vars
├── go.mod / go.sum                 # Go module
├── config.yaml.example             # Example config
├── .goreleaser.yaml                # Release build config
├── .github/workflows/release.yml   # CI/CD
├── cmd/
│   ├── root.go                     # Root command, Viper init, defaults
│   ├── onboard.go                  # Interactive TUI wizard (huh)
│   ├── detect.go                   # Polling engine + command execution
│   ├── list.go                     # Camera listing with pterm table
│   ├── service.go                  # Service subcommand group (+ alias svc)
│   ├── install.go                  # Service install + systemd unit patching
│   ├── uninstall.go                # Service uninstall
│   ├── start.go                    # Service start
│   ├── stop.go                     # Service stop
│   └── restart.go                  # Service restart (patches env file too)
├── internal/
│   ├── config/
│   │   ├── config.go               # Config struct, Defaults()
│   │   └── config_test.go
│   ├── detector/
│   │   ├── interface.go            # Detector interface + types
│   │   ├── detector.go             # New() factory
│   │   ├── v4l2_linux.go           # V4L2 via syscall (linux only)
│   │   ├── v4l2_stub.go            # Non-linux stub
│   │   └── lsof_linux.go           # lsof exec-based detection (linux)
│   ├── engine/
│   │   ├── engine.go               # Polling engine, debounce, hotplug
│   │   └── engine_test.go
│   ├── executor/
│   │   ├── executor.go             # Command execution with templates, timeout, overlap prevention
│   │   └── executor_test.go
│   └── output/
│       ├── output.go               # pterm wrappers, JWT redaction
│       └── output_test.go
└── packages/                       # Empty directory (reserved)
```
