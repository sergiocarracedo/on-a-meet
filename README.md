# on-a-meet

CLI tool that detects camera on/off state and triggers user-defined commands.

## Installation

### Homebrew (recommended)

```bash
brew tap sergiocarracedo/homebrew-tap
brew install on-a-meet
```

Works on both macOS and Linux (via the Homebrew prefix). Updates are handled by Homebrew — just run `brew upgrade on-a-meet`.

### Quick install (Linux & macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/sergiocarracedo/on-a-meet/main/install.sh | sudo bash
```

Detects your OS and architecture, downloads the correct binary, and installs it to `/usr/local/bin/`.

## Manual installation

### Binary (Linux)

Download the latest binary for your platform from the
[releases page](https://github.com/sergiocarracedo/on-a-meet/releases),
then make it executable:

```bash
chmod +x on-a-meet
sudo mv on-a-meet /usr/local/bin/
```

### Binary (macOS)

Download the latest darwin binary:

```bash
curl -L -o on-a-meet https://github.com/sergiocarracedo/on-a-meet/releases/latest/download/on-a-meet_darwin_amd64
chmod +x on-a-meet
sudo mv on-a-meet /usr/local/bin/
```

> **Permissions:** macOS detection uses the built-in `log` and `system_profiler`
> commands — no additional permissions and no `sudo`. The tool reads the unified
> system log rather than opening the camera, so no `NSCameraUsageDescription`
> entitlement is needed and the green indicator light is never turned on.

> **Which cameras are detected on macOS:** detection follows the camera power
> events that `UVCAssistant` writes to the unified log, which covers **USB/UVC
> webcams**. The **built-in MacBook camera** goes through a different subsystem
> that does not publish equivalent events, so it is reported as
> `not observable` rather than silently as "off". If you need built-in camera
> support, the CoreMediaIO property `kCMIODevicePropertyDeviceIsRunningSomewhere`
> reports it correctly, but reaching it requires cgo and a macOS build runner.

### From source

```bash
go install github.com/sergiocarracedo/on-a-meet@latest
```

Requires Go 1.22+. The binary is placed in `$GOPATH/bin` (or `$HOME/go/bin`).

### Linux Permissions

Access to `/dev/video*` devices is restricted. You have a few options:

#### Option A: Run with sudo (recommended for occasional use)

```bash
sudo on-a-meet detect
```

Root can access any device. No setup needed. Best for ad-hoc monitoring.

#### Option B: Install as a system service

```bash
sudo on-a-meet service install
```

The service runs as root by default and has full camera access. Best for
persistent background monitoring.

#### Option C: Add your user to the `video` group (convenient, has trade-offs)

```bash
sudo usermod -a -G video $USER
# Log out and back in, or run:
newgrp video
```

> **Trade-off:** Adding your user to the `video` group grants permanent read/write
> access to **all** camera devices on the system to every process you run — not
> just on-a-meet. This is generally safe on a personal machine but may be
> undesirable in multi-user or security-conscious environments. For those cases,
> use Option A (sudo) or Option B (service) instead.

## Configuration

Create `~/.config/on-a-meet/config.yaml`:

```yaml
interval: '1s'
debounce: 3
# Optional: defaults to v4l2 on Linux and darwin on macOS.
detect-method: 'v4l2'
timeout: '30s'
on-command: 'notify-send "Camera" "Camera turned ON"'
off-command: 'notify-send "Camera" "Camera turned OFF"'
```

On macOS the config lives at `~/.config/on-a-meet/config.yaml` and the service
runs as a per-user LaunchAgent, so nothing here needs `sudo`.

All options can also be set via CLI flags, which override config values.

### Environment file

Optionally source variables from a file before running `--on`/`--off` commands.
Variables are loaded as `KEY=VALUE` pairs (shell format, `export` prefix is allowed):

```yaml
environment-file: '/etc/default/on-a-meet'
```

Example env file (`/etc/default/on-a-meet`):

```
export HASS_TOKEN="eyJ..."  # API tokens, URLs, etc.
HASS_SERVER="http://hass.local:8123"
```

The file is re-read on **every** command execution, so edits take effect without
restarting anything. If the file cannot be read, or a line has no `KEY=VALUE`,
on-a-meet warns rather than failing quietly — a value wrapped across two lines
is a common cause of a token that silently arrives truncated.

These variables are merged into the environment of executed commands and can be
referenced as shell variables:

```bash
on-a-meet detect \
  --on 'curl -H "Bearer $HASS_TOKEN" $HASS_SERVER/api/webhook/camera-on'
```

On Linux, when running as a system service, the path is additionally written into
the systemd unit's `EnvironmentFile=` directive, making the variables available to
the service process itself. launchd has no equivalent directive, so on macOS only
the in-process expansion above applies — which is all the on/off commands need.

Any config key can also be overridden from the environment with an `ON_A_MEET_`
prefix and hyphens replaced by underscores, e.g.
`ON_A_MEET_ENVIRONMENT_FILE=~/.env on-a-meet detect`.

## Usage

### Detect camera state changes

```bash
# Basic monitoring
on-a-meet detect

# With commands on state change
on-a-meet detect --on 'echo "Camera ON at {{.Device}}"'

# With a specific camera
# Linux takes a device node, macOS the camera's display name:
on-a-meet detect --camera /dev/video0
on-a-meet detect --camera 'Logitech StreamCam'

# Custom polling interval
on-a-meet detect --interval 500ms

# Select detection backend
on-a-meet detect --detect lsof
```

### Template Variables

Available in `--on` and `--off` commands:

| Variable        | Description                  | Linux         | macOS                 |
| --------------- | ---------------------------- | ------------- | --------------------- |
| `{{.CameraID}}` | Short, stable device id      | `video0`      | `0x1114300046d0893`   |
| `{{.Device}}`   | Device path / display name   | `/dev/video0` | `Logitech StreamCam`  |
| `{{.State}}`    | Camera state                 | `on` / `off`  | `on` / `off`          |

### List devices

```bash
on-a-meet list
on-a-meet list --detect lsof
```

### Service management

On **Linux** the service is a system-wide systemd unit and needs `sudo`:

```bash
sudo on-a-meet service install
sudo on-a-meet service start
sudo on-a-meet service stop
sudo on-a-meet service restart
sudo on-a-meet service uninstall
```

It reads `/etc/on-a-meet/config.yaml`.

On **macOS** the service is a per-user LaunchAgent, so run these **without**
`sudo` — the commands refuse to run as root, because a LaunchAgent installed by
root lands in `/var/root/Library/LaunchAgents` and would never run for you:

```bash
on-a-meet service install
on-a-meet service start
on-a-meet service stop
on-a-meet service restart
on-a-meet service uninstall
```

It reads `~/.config/on-a-meet/config.yaml`, installs
`~/Library/LaunchAgents/on-a-meet.plist`, and logs to `~/Library/Logs/`.

### Interactive setup

The `onboard` command guides you through an interactive wizard:

```bash
on-a-meet onboard                # Full wizard
on-a-meet onboard --dry-run      # Preview config without applying
on-a-meet onboard --yes          # Skip confirmations (implied when piped)
```

The wizard steps through:

1. **Detection method** — only the methods that work on your OS are offered
   (`v4l2` or `lsof` on Linux; the unified-log method on macOS). When a
   platform has a single method it is reported rather than prompted for.
2. **Camera selection** — Multi-select from detected cameras
3. **Live test** — Verify the method can detect camera state
4. **Commands** — Configure `--on` and `--off` commands with template support
5. **Environment file** *(optional)* — Point at a `KEY=VALUE` file holding API
   tokens or URLs for your commands to reference as `$NAME`. The path is read
   back immediately and the wizard reports how many variables it found, so a
   wrong path is caught here rather than surfacing later as a `401`.
6. **Verbose output** — Off by default; enable it to log each command and its
   output while you are setting things up.
7. **Apply** — Write config and optionally install the service

On Linux the wizard re-execs itself under `sudo` to write
`/etc/on-a-meet/config.yaml` and install the systemd unit. On macOS it writes
`~/.config/on-a-meet/config.yaml` and installs a per-user LaunchAgent in
process — it must **not** be run with `sudo`.
