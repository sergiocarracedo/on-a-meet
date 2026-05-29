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

> **Permissions:** macOS detection uses the built-in `log` and `system_profiler` commands — no additional permissions required. The tool reads system logs rather than accessing camera hardware directly, so no `NSCameraUsageDescription` entitlement is needed.

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
detect-method: 'v4l2'
timeout: '30s'
on-command: 'notify-send "Camera" "Camera turned ON"'
off-command: 'notify-send "Camera" "Camera turned OFF"'
```

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

These variables are merged into the environment of executed commands and can be
referenced as shell variables:

```bash
on-a-meet detect \
  --on 'curl -H "Bearer $HASS_TOKEN" $HASS_SERVER/api/webhook/camera-on'
```

When running as a system service, the path is also written into the systemd unit's
`EnvironmentFile=` directive, making variables available to the service process itself.

## Usage

### Detect camera state changes

```bash
# Basic monitoring
on-a-meet detect

# With commands on state change
on-a-meet detect --on 'echo "Camera ON at {{.Device}}"'

# With a specific camera
on-a-meet detect --camera /dev/video0

# Custom polling interval
on-a-meet detect --interval 500ms

# Select detection backend
on-a-meet detect --detect lsof
```

### Template Variables

Available in `--on` and `--off` commands:

| Variable        | Description       | Example       |
| --------------- | ----------------- | ------------- |
| `{{.CameraID}}` | Short device name | `video0`      |
| `{{.Device}}`   | Full device path  | `/dev/video0` |
| `{{.State}}`    | Camera state      | `on` or `off` |

### List devices

```bash
on-a-meet list
on-a-meet list --detect lsof
```

### Service management

```bash
# Install as system service (requires sudo)
sudo on-a-meet service install

# Start / Stop / Restart (requires sudo)
sudo on-a-meet service start
sudo on-a-meet service stop
sudo on-a-meet service restart

# Uninstall (requires sudo)
sudo on-a-meet service uninstall
```

The service uses `/etc/on-a-meet/config.yaml` for configuration.

### Interactive setup

The `onboard` command guides you through an interactive wizard:

```bash
on-a-meet onboard                # Full wizard
on-a-meet onboard --dry-run      # Preview config without applying
```

The wizard steps through:

1. **Camera selection** — Multi-select from detected cameras
2. **Detection method** — Choose `v4l2`, `lsof`, or `darwin` (macOS)
3. **Live test** — Verify the method can detect camera state
4. **Commands** — Configure `--on` and `--off` commands with template support
5. **Apply** — Write config and optionally install as a system service

On Linux, the sudo apply path writes config to `/etc/on-a-meet/config.yaml`
and installs the system service automatically.
