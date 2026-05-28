# Stack Research

**Domain:** Linux camera state detection CLI tool
**Researched:** 2026-05-28
**Confidence:** HIGH

## Recommended Stack

### Core Technologies

| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go | 1.22+ | Language runtime | Cross-compilation to single binary, excellent stdlib for syscalls, V4L2 access via syscall package without CGo |
| spf13/cobra | v1.10.2 | CLI framework | Battle-tested (Kubernetes, Docker, Hugo, GitHub CLI), POSIX flags, subcommands, shell completion |
| spf13/viper | v1.19+ | Configuration management | Reads YAML, env vars, CLI flags with clear precedence hierarchy — exactly what we need for flag-over-config pattern |

### Supporting Libraries

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| pterm/pterm | v0.12.83 | Terminal output | All user-facing output — colored status, spinners, tables for device list |
| kardianos/service | v1.2.4 | Service management | --install/--uninstall flags; supports systemd, launchd, OpenRC out of the box |
| vladimirvivien/go4vl | v0.x | V4L2 abstraction | Optional V4L2 backend — idiomatic Go channels, zero-copy MMAP. Alternative: implement raw ioctl via syscall. |
| coreos/go-systemd | v22 | systemd D-Bus | Optional: notify systemd of readiness, watchdog support if running as user service |

### Development Tools

| Tool | Purpose | Notes |
|------|---------|-------|
| golangci-lint | Linting | Standard Go linting |
| goreleaser | Cross-platform builds | Build for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64 |
| go-task or Makefile | Build automation | Simple build targets |

## Alternatives Considered

| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| spf13/cobra | urfave/cli v3 | urfave/cli is simpler if you don't need subcommands or nested help. Cobra wins for --install/--uninstall service subcommands. |
| kardianos/service | Custom systemd unit template | kardianos/service is worth the dependency — handles install paths, user detection, sudo prompts across platforms |
| go4vl / raw ioctl | blackjack/webcam | blackjack/webcam (v0.6.1) is older and less maintained. go4vl is more active and idiomatic. Raw ioctl is also valid for simple state checking. |
| eBPF (cilium/ebpf) | V4L2 polling | eBPF would give event-driven detection but adds significant complexity (kernel BTF dependency, root, C compilation). Overkill for polling approach. |

## What NOT to Use

| Avoid | Why | Use Instead |
|-------|-----|-------------|
| Python or shell script | Can't produce single binary, no cross-compile, slower startup | Go |
| Direct CGo for V4L2 | Adds build complexity, cross-compilation pain | Go syscall package or go4vl |
| systemd unit hardcoded paths | Breaks when binary moves | kardianos/service handles this |
| TOML or JSON for config | Less human-friendly for a system utility | YAML |

## Versions

### Version Compatibility

| Package A | Compatible With | Notes |
|-----------|-----------------|-------|
| spf13/cobra v1.10.2 | Go 1.18+ | Latest stable |
| spf13/viper v1.19+ | Go 1.21+ | YAML support via mapstructure |
| pterm/pterm v0.12.83 | Go 1.21+ | Latest as of Feb 2026 |
| kardianos/service v1.2.4 | Go 1.21+ | Supports systemd, launchd, OpenRC |
| vladimirvivien/go4vl | Linux only | Go 1.16+, V4L2 kernel headers |

### Installation

```bash
# Initialize module
go mod init github.com/yourusername/on-a-meet

# Core dependencies
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/pterm/pterm@latest

# Service management
go get github.com/kardianos/service@latest

# V4L2 (Linux only, optional - can use syscall directly)
go get github.com/vladimirvivien/go4vl@latest
```

---
*Stack research for: on-a-meet camera detection CLI*
*Researched: 2026-05-28*
