# INTEGRATIONS.md — External Services and Integrations

## Overview

`on-a-meet` is a **standalone, self-contained CLI tool** with no dependency on external web APIs, databases, message queues, or third-party cloud services. All integrations are with the local operating system — Linux kernel interfaces, system services, and the filesystem.

---

## External APIs

**None.** The tool does not call any external HTTP APIs, REST endpoints, web services, or cloud platforms. It has no network dependencies at runtime. There are no webhook integrations (neither outbound nor inbound), no OAuth flows, and no telemetry/analytics endpoints.

---

## System Service Integrations

### 1. Service Manager (systemd / launchd)

**Library:** `github.com/kardianos/service` v1.2.4

**Files:**
- `cmd/install.go` (lines 1-104)
- `cmd/uninstall.go` (lines 1-46)
- `cmd/start.go` (lines 1-39)
- `cmd/stop.go` (lines 1-39)
- `cmd/restart.go` (lines 1-50)

The `kardianos/service` library provides a cross-platform abstraction over:
- **Linux:** systemd (generates `/etc/systemd/system/on-a-meet.service`)
- **macOS:** launchd (generates a launchd plist)

The service runs the `detect` subcommand with `--config /etc/on-a-meet/config.yaml`.

**Service configuration** (`cmd/install.go:21-33`):
- Service name: `on-a-meet`
- Display name: `on-a-meet`
- Description: `Camera state monitoring service`
- Arguments: `["detect", "--config", "/etc/on-a-meet/config.yaml"]`
- Working directory: `/`
- User: the original `SUDO_USER` (extracted from env, line 57)

**systemd unit patching** (`cmd/install.go:35-54`):
After kardianos/service installs the unit file, the tool optionally patches the `EnvironmentFile` directive. It reads `/etc/systemd/system/on-a-meet.service`, replaces the default line `EnvironmentFile=-/etc/sysconfig/on-a-meet` with the user-configured path (e.g., `/etc/default/on-a-meet`), and runs `systemctl daemon-reload`. This is handled in the `patchUnitEnvironmentFile()` function.

### 2. V4L2 (Video4Linux2) Kernel Interface

**File:** `internal/detector/v4l2_linux.go` (lines 1-155)

**Build constraint:** `//go:build linux`

The primary detection backend uses direct Linux kernel syscalls:

- **Device enumeration:** `filepath.Glob("/dev/video*")` — scans the V4L2 device filesystem.
- **I/O control (ioctl):** Uses `syscall.Syscall(syscall.SYS_IOCTL, ...)` with the `VIDIOC_QUERYCAP` command (`_IOR('V', 0, ...)`) to query device capabilities. Filters for devices with `V4L2_CAP_VIDEO_CAPTURE` capability bit (0x00000001).
- **Open handle detection:** Walks `/proc/[pid]/fd/` symlinks to determine if any process currently has a file descriptor open on a given `/dev/video*` device (`hasOpenHandle()` function, lines 100-130). This is the secondary check.
- **Device busy detection:** When `unix.Open()` returns `EBUSY`, the device is reported as ON (a process has it exclusively open). When the open succeeds but no process has an fd open, the device is reported as OFF.
- **Library:** Uses `golang.org/x/sys/unix` for `unix.Open`, `unix.Close`, and error constants.

**Permissions note:** Access to `/dev/video*` devices requires membership in the `video` group (documented in README and output hints).

### 3. lsof (List Open Files) Command

**File:** `internal/detector/lsof_linux.go` (lines 1-52)

**Build constraint:** `//go:build linux`

An alternative detection backend that:

- **Executes:** `exec.Command("lsof", devicePath)` — runs the system `lsof` command to check if any process has the device open.
- **Exit code interpretation:**
  - Exit code 0 → device is in use (ON)
  - Exit code 1 → device is not in use (OFF)
  - Any other exit code → wraps as a detection error
- **Device enumeration:** Same as V4L2 backend — uses `filepath.Glob("/dev/video*")` and `openAndQueryCap()` to gather device metadata.

**Prerequisite:** The `lsof` command must be installed on the system.

### 4. Detection Factory

**File:** `internal/detector/detector.go` (lines 1-14)

The `New(method string)` function is the factory that selects the detection backend:

```go
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

### 5. Detector Interface

**File:** `internal/detector/interface.go` (lines 1-20)

```go
type Detector interface {
    ListDevices() ([]DeviceInfo, error)
    Detect(devicePath string) (DeviceStatus, error)
}
```

This interface is implemented by both `V4L2Detector` and `LsofDetector`, and is consumed by the polling engine.

---

## Filesystem Interactions

### Config File Reading (via Viper)

**File:** `cmd/root.go` (lines 49-79)

Viper searches for `config.yaml` in the following locations (in order):

1. `~/.config/on-a-meet/config.yaml` (user config)
2. `/etc/on-a-meet/config.yaml` (system config)
3. `.` (current directory, for development/testing)

Config can also be explicitly set via `--config` flag (line 41).

Default values (lines 66-76):

| Key | Default |
|---|---|
| `detect-method` | `"v4l2"` |
| `interval` | `"1s"` |
| `debounce` | `3` |
| `timeout` | `"30s"` |
| `camera` | `""` (all) |
| `on-command` | `""` |
| `off-command` | `""` |
| `silent` | `false` |
| `verbose` | `false` |
| `environment-file` | `""` |

### Config File Writing (via onboard wizard)

**File:** `cmd/onboard.go` (lines 106-133)

The `onboard --apply` path:
1. Reads a temporary JSON file (`/tmp/on-a-meet-onboard.json`)
2. Marshals it to YAML via `gopkg.in/yaml.v3`
3. Writes to `/etc/on-a-meet/config.yaml`
4. Creates directory `/etc/on-a-meet/` if it does not exist (mode 0755)

### Environment File Reading

**File:** `internal/executor/executor.go` (lines 46-68)

The `parseEnvFile()` function reads a user-specified environment file (e.g., `/etc/default/on-a-meet` or `/etc/sysconfig/on-a-meet`) before executing commands. Format: shell-style `KEY=VALUE` lines, supporting:
- Comments (`#`)
- Empty lines
- `export` prefix stripping

Variables from this file are merged into the command's environment.

### Process Execution

**File:** `internal/executor/executor.go` (lines 70-152)

User-defined commands are executed via:

```go
cmd := exec.CommandContext(cmdCtx, "sh", "-c", rendered)
cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
```

Key behavior:
- Commands run through `sh -c` (shell execution)
- `Setpgid: true` creates a new process group for clean cleanup
- Timeout: configurable (default 30s), enforced via `context.WithTimeout`
- Kill behavior: `SIGKILL` sent to the process group on timeout/cancel
- Overlap prevention: `sync.Map` keyed by state (`"on"` / `"off"`) prevents concurrent execution of the same state type
- Template variables (`{{.CameraID}}`, `{{.Device}}`, `{{.State}}`) are substituted via Go's `text/template`
- Environment variables from the env file and `os.Environ()` are available to the command
- JWT secrets in command output are redacted via regex (`output.RedactSecrets()`)

### Temp File Usage

**File:** `cmd/onboard.go` (lines 201-204, 516-518)

The onboard wizard writes a JSON config snapshot to `/tmp/on-a-meet-onboard.json` before re-executing itself with `sudo`. This temp file is read by the `--apply` path and not cleaned up after use.

---

## Webhooks

**None.** There is no webhook sending or receiving capability. The tool does not listen on any network port, nor does it make any outbound HTTP requests. State changes trigger local shell commands only.

---

## Message Queues

**None.** No message brokers (RabbitMQ, Kafka, NATS, etc.), no pub/sub systems, no IPC message queues.

---

## Databases

**None.** No SQL or NoSQL databases. No persistent storage beyond:
- The YAML config file (read/write)
- The environment file (read-only)
- The systemd unit file (created/patched by kardianos/service)

---

## Signal Handling

**File:** `cmd/detect.go` (lines 138-144)

The tool handles OS signals for graceful shutdown:

```go
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
```

On `SIGINT` or `SIGTERM`, the context is cancelled, which propagates to:
- The polling engine's `Run()` loop (stops polling)
- Any currently executing commands (process group is killed)

---

## Permissions Model

The tool has a two-tier permission model:

1. **Standard user:** Can run `detect` and `list` commands (requires `video` group for `/dev/video*` access). Can run the `onboard` wizard interactively (wizard runs v4l2 detection as the current user, then re-executes itself via `sudo` for installation).

2. **Root (sudo):** Required for `service install`, `service uninstall`, `service start`, `service stop`, `service restart`, and the `onboard --apply` path. The install command detects the original user via `SUDO_USER` environment variable and runs the service under that user identity.
