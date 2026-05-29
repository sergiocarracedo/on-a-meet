# on-a-meet

CLI tool that detects camera on/off state and triggers user-defined commands.

## Installation

```bash
# Download the latest release from GitHub
# or build from source:
go install github.com/sergiocarracedo/on-a-meet@latest
```

### From source

```bash
git clone https://github.com/sergiocarracedo/on-a-meet
cd on-a-meet
go build -o on-a-meet .
```

### Permissions

Access to `/dev/video*` devices requires the `video` group:

```bash
sudo usermod -a -G video $USER
# Log out and back in, or run:
newgrp video
```

## Configuration

Create `~/.config/on-a-meet/config.yaml`:

```yaml
interval: "1s"
debounce: 3
detect-method: "v4l2"
timeout: "30s"
on-command: 'notify-send "Camera" "Camera turned ON"'
off-command: 'notify-send "Camera" "Camera turned OFF"'
```

All options can also be set via CLI flags, which override config values.

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

| Variable       | Description          | Example              |
|----------------|----------------------|----------------------|
| `{{.CameraID}}` | Short device name   | `video0`             |
| `{{.Device}}`   | Full device path    | `/dev/video0`        |
| `{{.State}}`    | Camera state        | `on` or `off`        |

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
