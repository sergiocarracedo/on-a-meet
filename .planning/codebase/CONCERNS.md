# Technical Debt, Security, and Fragile Areas

> **Audited:** 2026-05-29
> **Commit:** 26464c5 (HEAD, milestone v1.1.0 in definition)
> **Project version:** v1.0.0 shipped, v1.1.0 in definition

---

## 1. TODO / FIXME / HACK Comments

No `TODO`, `FIXME`, `HACK`, `XXX`, `BUG`, `WORKAROUND`, or `BROKEN` comments exist
in any `.go` file. The codebase is clean in this regard.

The only edge-case hit was a test helper:

- **`internal/output/output_test.go:13`** — calls `os.Unsetenv("PTERM_DEBUG")`
  inside a test. This is not a code-level marker but is a fragile side-effect
  pattern (see Section 6).

**Finding:** Zero todo/fixme markers. Good discipline, but means known limitations
(see Section 8) are undocumented _in code_.

---

## 2. File Sizes & Refactoring Candidates

All files are small (no single file > 550 lines), but a few are approaching
the threshold where splitting would improve maintainability:

| File | Lines | Concern |
|------|-------|---------|
| `cmd/onboard.go` | 543 | Largest file. Single 300+ line `RunE` closure with three code paths (--apply, --config, interactive). The interactive wizard (~250 LoC) could be extracted into helper methods. |
| `cmd/detect.go` | 194 | 80-line `RunE` closure with inline goroutine callbacks. On-change callbacks (lines 92-127) are defined inline, mixing orchestration with presentation. |
| `internal/engine/engine.go` | 198 | Contains both state-machine logic (debounce, transitions) and hotplug detection. Could split into `engine.go` + `hotplug.go` for clarity. |
| `internal/detector/v4l2_linux.go` | 155 | Packs ioctl structs, /proc scanning, and device enumeration into one file. The `hasOpenHandle` function (line 100-130) is an expensive /proc walk called on every `Detect()` poll. |

**Recommendation:** Refactor `cmd/onboard.go` when adding new wizard features
(milestone v1.1.0). Extract the inline callback in `cmd/detect.go` into a named
function or method.

---

## 3. Security Analysis

### 3.1 Hardcoded Secrets, Tokens, or Keys

**No hardcoded secrets found.** The `go.sum` file contains base64-encoded
hashes (e.g., `h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=`) — these are
cryptographic checksums for Go module verification, NOT secrets.

### 3.2 Shell Injection Risk — Executor (CRITICAL)

**File:** `internal/executor/executor.go:107`

```go
cmd := exec.CommandContext(cmdCtx, "sh", "-c", rendered)
```

The `rendered` string is the user's configured command after Go `text/template`
substitution and `os.Expand` environment variable expansion. Because the command
is executed via `sh -c`, any user-controlled values in the template data
(`{{.CameraID}}`, `{{.Device}}`, `{{.State}}`) that contain shell metacharacters
(`;`, `` ` ``, `$()`, `|`, etc.) could lead to arbitrary command injection.

**Risk factors:**
- `{{.Device}}` is user-visible device paths from `/dev/video*`. These are
  unlikely to contain shell metacharacters in practice (Linux device paths are
  alphanumeric), but an attacker who can control device naming (e.g., udev
  rules) could inject.
- `{{.CameraID}}` is derived via `path[5:]` — e.g., `/dev/video0` → `video0`.
  This is a safe subslice of `/dev/video*` paths, but still technically
  unsanitized.
- Users who write `{{.Device}}` or `{{.CameraID}}` directly into on/off commands
  are exposed.

**Mitigations present:**
- `text/template` (not `text/template` with `funcMap` that could eval) —
  only simple field access, no function calls.
- `os.Expand` only expands `$VAR` or `${VAR}` — no shell eval.
- The `--on`/`--off` commands come from config/CLI, not from user input at
  runtime.

**Recommendation:** Add shell escaping for template substitution values.
Use `strings.NewReplacer` to escape common shell metacharacters (`'`, `"`,
`$`, `` ` ``, `\`, `;`, `|`, `&`, `(`, `)`, `{`, `}`, `<`, `>`, newline)
in the rendered command, or use `syscall.Exec` with explicit argv instead of
`sh -c`.

### 3.3 JWT Token Handling

**File:** `internal/output/output.go:10-45`

```go
var jwtRe = regexp.MustCompile(`ey[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`)

func RedactSecrets(s string) string {
    return jwtRe.ReplaceAllString(s, "ey***.***.***")
}
```

**Finding:** JWT redaction is applied in debug output (executor.go:137-138)
and config display (detect.go:65). The regex matches the standard
`header.payload.signature` JWT format. This is a best-effort mitigation.

**Concerns:**
- The regex is compile-on-init — a panic if invalid, but it is valid.
- The regex may produce false positives (matching non-JWT base64url strings)
  and false negatives (non-standard JWT formats, truncated tokens).
- Only JWT tokens are redacted. Other secrets (API keys, passwords in env
  files, etc.) are NOT redacted. The `detect.go:65` line prints full config
  including on/off commands which may contain embedded secrets:
  ```go
  output.Info.Printfln("Config: ... on=%s off=%s ...",
      output.RedactSecrets(cfg.OnCmd), output.RedactSecrets(cfg.OffCmd))
  ```
- The env file variables are injected into `cmd.Env` (executor.go:110-115)
  and are visible in `/proc/PID/environ` to any process with appropriate
  privileges.

**Recommendation:** Add a more general secret redactor (e.g., match common
patterns like `key=`, `secret=`, `token=`, `password=`). Consider doing
redaction earlier in the pipeline (before the value reaches output).

### 3.4 Temporary File Race (Medium)

**File:** `cmd/onboard.go:201-212` and `cmd/onboard.go:516-527`

```go
tmpPath := "/tmp/on-a-meet-onboard.json"
if err := os.WriteFile(tmpPath, data, 0644); err != nil { ... }
sudoCmd := exec.Command("sudo", binary, "onboard", "--apply", tmpPath)
```

**Finding:** A fixed, predictable temp file path `/tmp/on-a-meet-onboard.json`
is used to pass configuration from the non-sudo wizard process to the
sudo-re-executed apply process. This creates a TOCTOU (time-of-check /
time-of-use) race condition:

1. A local attacker could pre-create a symlink at `/tmp/on-a-meet-onboard.json`
   pointing to `/etc/on-a-meet/config.yaml` (or any other file the user has
   write access to).
2. When the sudo process reads the file, it follows the symlink and
   overwrites the target with attacker-controlled content.
3. The file is written with `0644` permissions — world-readable.

**Mitigations possible:**
- Use `os.MkdirTemp` to create a directory with `0700` permissions and write
  the config inside it.
- Pass the config content via stdin to the sudo process instead of via a file:
  ```go
  sudoCmd.Stdin = strings.NewReader(string(data))
  ```
- Use a random suffix or PID in the filename.

**Recommendation:** High priority to fix due to the sudo escalation context.
Use `os.MkdirTemp("", "on-a-meet-*")` or pipe via stdin.

---

## 4. Permission Model

### 4.1 Root vs. User Permissions

The tool operates in two privilege contexts:

| Mode | User | Accesses | Concerns |
|------|------|----------|----------|
| `detect` / `list` (user) | Normal user | /dev/video* via V4L2 or lsof | Requires `video` group membership |
| `service install/start/stop/restart/uninstall` | root | systemd unit files | Sudo check at entry (e.g., `os.Geteuid() != 0`) |
| `onboard --apply` | root | /etc/on-a-meet/ + systemd | Sudo re-exec |
| `onboard` (interactive) | Normal user | /dev/video* only | Non-sudo |

**Video group check:** The tool prints a helpful tip when no cameras are
detected (`detect.go:56-57`, `list.go:38-39`), but does NOT proactively check
the user's group membership at startup. The error surfaces only when actual
detection fails (v4l2_linux.go:151).

**SUDO_USER assumption:** `cmd/install.go:57` uses `os.Getenv("SUDO_USER")`
to determine the service run-as user. This variable is set by `sudo` but may
not be set when running via `doas` or `pkexec`. The `serviceConfig("")` call
in start/stop/restart passes an empty string — the `kardianos/service` library
defaults to root in that case.

### 4.2 Service Runs as Original User

Quick task 020 set `UserName: user` from `SUDO_USER`. This means the service
runs commands as the original (non-root) user, which is correct for
permissions, but means the service inherits that user's environment and
permissions — including the `video` group membership requirement.

**Recommendation:** Add a proactive video group check at the start of
`detect` and `list` commands (before attempting any V4L2 calls). This
would give a clearer error message than the current fallback.

---

## 5. Cross-Platform Issues

### 5.1 Stub Files Are Complete No-Ops

| File | Build Tag | Behavior |
|------|-----------|----------|
| `internal/detector/v4l2_stub.go` | `!linux` | Always returns `"V4L2 detection is only supported on Linux"` |
| `internal/detector/lsof_stub.go` | `!linux` | Always returns `"lsof detection is only supported on Linux"` |

**Finding:** On macOS (Darwin), BOTH V4L2 and lsof backends return errors.
This means the tool is effectively **non-functional on macOS** despite
GoReleaser building `darwin/amd64` and `darwin/arm64` binaries. The `"lsof"
detection is only supported on Linux` message from the stub is misleading —
`lsof` itself exists and works on macOS.

**Impact:** Users downloading the macOS binary from GitHub Releases will get
a tool that compiles but cannot detect cameras. No macOS-specific backend
(AVFoundation, IOKit, or native `lsof` calling) has been implemented.

**Relevant files:**
- `.goreleaser.yaml:11-12` — builds `darwin/amd64` and `darwin/arm64`
- `.planning/PROJECT.md:29-30` — "macOS camera detection backend" is planned
  for v1.1.0

### 5.2 macOS lsof

The `LsofDetector` on Linux passes the device path as an argument to `lsof`.
On macOS, `lsof` works but `/dev/video*` does not exist — macOS uses a
different device model (e.g., `AVCaptureDevice`). A macOS lsof backend
would need to query different device paths or use `IOKit` directly.

**Recommendation:** Implement a `lsof_darwin.go` variant that uses macOS
`lsof` to detect camera processes (e.g., `lsof -c VDCAssistant`). Or
implement a proper AVFoundation-based detector. This is the stated goal
of milestone v1.1.0.

### 5.3 GoReleaser Go Version Mismatch

- `go.mod` line 3: `go 1.25.0`
- `.github/workflows/release.yml` line 21: `go-version: "1.22"`

The GitHub Actions workflow pins Go 1.22 for building, but the module
declares 1.25.0. If the project uses Go 1.24+ features (e.g., the new
`omitzero` struct tag, `crypto/tls` defaults, etc.), the build might fail
or behave differently. At minimum, `go mod tidy` may produce a different
`go.sum` on 1.22 vs 1.25.

**Recommendation:** Align the workflow Go version with `go.mod`. Either
update `go-version` to `"1.25"` or (if no 1.25 features are used) revert
`go.mod` to `go 1.22`.

---

## 6. Fragile Areas

### 6.1 Detector Factory

**File:** `internal/detector/detector.go:5-13`

```go
func New(method string) (Detector, error) {
    switch method {
    case "v4l2":
        return NewV4L2Detector(), nil
    case "lsof":
        return NewLsofDetector(), nil
    default:
        return nil, fmt.Errorf(...)
    }
}
```

**Fragility:**
- `NewV4L2Detector()` exists for all platforms (via stubs), but the stubs
  always return errors. The factory has no way to say "this method is
  registered but not supported on this OS" other than returning a runtime
  error.
- Adding a new backend (e.g., macOS AVFoundation) requires modifying this
  switch statement — an open-closed principle violation.
- No backend registration mechanism (no plugin, no init-time registration).

**Recommendation:** Consider a registry pattern (`Register(name, factory)`)
or a platform-aware decision table.

### 6.2 V4L2 Raw Syscall

**File:** `internal/detector/v4l2_linux.go:75`

```go
if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd),
    uintptr(vidiocQueryCap), uintptr(unsafe.Pointer(&cap))); err != 0 {
```

**Fragility:**
- Uses `unsafe` package — memory-unsafe by definition. A malformed struct
  could cause memory corruption.
- The `_IOR` macro (line 18-20) manually constructs the ioctl request code
  using bit shifts. This is architecture-dependent (size of `uintptr`,
  endianness on different platforms). The `v4l2Capability` struct (line 22-30)
  uses fixed-size `[16]byte` and `[32]byte` arrays — correct for x86_64 and
  arm64, but the struct alignment and padding must match the kernel's
  expectation exactly.
- Uses `golang.org/x/sys/unix` for `Open`/`Close` (clean) but falls back to
  raw `syscall.Syscall` for `SYS_IOCTL` (unclean). Inconsistent.

**Recommendation:** Use `unix.IoctlGetTermios`-style wrapping or
`unix.Ioctl` from `golang.org/x/sys/unix` instead of raw syscall. Consider
using the `golang.org/x/sys/unix` package's V4L2 constants if available, or
add a portable ioctl helper.

### 6.3 /proc Filesystem Walk

**File:** `internal/detector/v4l2_linux.go:100-130`

```go
func hasOpenHandle(devicePath string) bool {
    entries, err := os.ReadDir("/proc")
    ...
    for _, entry := range entries {
        ...
        fdDir := filepath.Join("/proc", pid, "fd")
        fds, err := os.ReadDir(fdDir)
        ...
        for _, fd := range fds {
            link, err := os.Readlink(linkPath)
            ...
        }
    }
}
```

**Fragility:**
- Called on **every** `Detect()` poll cycle — an O(N*M) walk of every PID
  and every file descriptor, on every poll. For a 1-second interval on a
  system with hundreds of processes, this is significant overhead.
- Permission-dependent: cannot read `/proc/PID/fd` for processes owned by
  other users. Returns `false` (no open handle) even when another user's
  process IS using the camera — leading to false negatives.
- /proc scanning is a Linux-specific pattern that does not exist on macOS.

**Recommendation:** The /proc walk is used to distinguish "device open by
another process" (ON) from "device just opened by our own poll" (OFF). A
less expensive approach: open the device, immediately close it, then check
if `EBUSY` on a second open attempt. If the first open succeeds and the
second fails with `EBUSY`, someone else is using it.

### 6.4 GoReleaser LDFLAGS Version Injection

**File:** `.goreleaser.yaml:17`

```yaml
ldflags:
  - -s -w -X main.version={{.Version}} -X main.commit={{.Commit}} -X main.date={{.Date}}
```

**Fragility:**
- The `{{.Version}}`, `{{.Commit}}`, `{{.Date}}` are GoReleaser template
  variables. If any of these are empty (e.g., dirty git state, missing tag),
  the binary will compile with empty/invalid version strings.
- The `main.version` variable (in `main.go:6`) has a fallback default `"dev"`,
  so this degrades gracefully for local builds, but `main.commit` and
  `main.date` will be empty strings.
- The `-s -w` flags strip symbol table and DWARF debug info — this is fine
  for release builds but makes debugging crashes harder.

### 6.5 CameraID Path Slicing

**File:** `cmd/detect.go:103,117`

```go
CameraID: path[5:],
```

`path` is a device path like `/dev/video0`. Slicing from index 5 gives
`video0`. This is correct for `/dev/` prefixed paths (5 characters) but
assumes the path always starts with `/dev/`. If the path is not a standard
`/dev/video*` device (e.g., a udev symlink in another directory), the
slicing produces unexpected results. If the path is shorter than 5
characters, this would panic (index out of range).

**Recommendation:** Use `strings.TrimPrefix(path, "/dev/")` with a fallback
to `filepath.Base(path)` for safer extraction.

### 6.6 Env File Parsing Fragility

**File:** `internal/executor/executor.go:56-66`

```go
for _, line := range strings.Split(string(data), "\n") {
    line = strings.TrimSpace(line)
    if line == "" || strings.HasPrefix(line, "#") { continue }
    line = strings.TrimPrefix(line, "export ")
    parts := strings.SplitN(line, "=", 2)
    if len(parts) == 2 {
        vars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }
}
```

**Fragility:**
- `strings.TrimSpace` on values strips intentional leading/trailing whitespace.
  A value like `PASSWORD=" secret "` becomes `" secret "` (with quotes) or
  `secret` (without quotes, but whitespace stripped). This is inconsistent.
- Quoting is NOT handled. `KEY="value with spaces"` will parse as:
  - key: `KEY`
  - value: `"value with spaces"` (with the literal double quotes).
- Lines with embedded `=` signs in values work correctly (SplitN with N=2),
  but unquoted values with trailing spaces get trimmed.
- The file is re-read and re-parsed on EVERY command execution (line 99),
  not cached. This is wasteful for a file that rarely changes.

**Recommendation:** Either document the quoting limitation, or use a proper
env file parser. Consider caching the parsed vars with a TTL or file
modification time check.

### 6.7 Goroutine Fire-and-Forget

**Files:**
- `cmd/detect.go:101-110,115-124` — on-change goroutines

```go
go func() {
    if err := exec.ExecOn(context.Background(), cfg.OnCmd, data); err != nil {
        output.Warning.Printfln("on-command failed: %v", err)
    }
}()
```

**Finding:** The goroutines spawned for command execution are fire-and-forget.
There is no:
- Error propagation back to the caller.
- WaitGroup tracking for graceful shutdown.
- Limit on the number of concurrent goroutines (though `Executor.exec` uses
  `running sync.Map` to prevent same-state overlap).

If the camera state flips rapidly (e.g., USB glitch), many goroutines could
stack up waiting on `ExecOn`/`ExecOff` (which block on `sh -c` command
execution). The `sync.Map` only prevents duplicates for the SAME state
(detect.go uses distinct `"on"` and `"off"` keys), so rapid on-off-on could
stack two goroutines.

**Recommendation:** Add a `sync.WaitGroup` in `detect.go` for graceful
shutdown. Use a bounded goroutine pool or semaphore channel to limit
concurrent command executions.

---

## 7. Input Validation

### 7.1 CLI Arguments

| Flag | Validation | Location |
|------|-----------|----------|
| `--interval` | `time.ParseDuration` | `cmd/detect.go:38-41` |
| `--timeout` | `time.ParseDuration` | `cmd/detect.go:79-81` |
| `--camera` | None (free string) | — |
| `--on` / `--off` | None (free string) | — |
| `--detect` / `--detect-method` | Validated by `detector.New()` factory | `cmd/detect.go:43-46` |
| `--debounce` (CLI) | None (passes through to Viper `GetInt`) | — |
| onboard debounce | `strconv.Atoi` + `>= 1` check | `cmd/onboard.go:278-289` |
| onboard interval | `time.ParseDuration` | `cmd/onboard.go:297-301` |

**Missing validation:**
- `--camera` accepts any string — no check that the path exists or is a
  valid device.
- `--on` / `--off` commands are not validated before use. Invalid templates
  are caught at execution time (`text/template.Parse` returns error), but
  references to non-existent env vars silently expand to empty strings.
- `--debounce` CLI flag value is not explicitly validated — Viper's
  `GetInt` returns 0 for invalid input, which disables debounce entirely
  (allowing every poll to trigger a state change).
- No path traversal protection on `--config` file paths (read elsewhere).

### 7.2 Template Input Validation

Template parsing is done at execution time, not at config load time:

```go
tmpl, err := template.New("").Parse(cmdStr)  // executor.go:89
```

An invalid template in the config file will only fail at runtime when the
camera transitions. This could be confusing: the service starts, appears
healthy, then fails on the first ON event.

**Recommendation:** Add a `--dry-run` or `--validate` mode that pre-parses
the configured command templates and checks for basic syntax errors. The
onboard `--dry-run` flag does this partially, but only for the interactive
setup path, not for the `detect` command itself.

---

## 8. Known Bugs and Edge Cases from Quick Tasks

The following issues were discovered and fixed across 25 quick tasks. Some
represent systemic fragility that may recur:

| # | Issue | Root Cause | Fix File(s) | Risk of Regression |
|---|-------|-----------|-------------|---------------------|
| 001 | `--detect` flag defaulted to empty string | Config priority conflict | `cmd/detect.go` | Medium — config merging is complex |
| 003 | V4L2 ON detection always returned OFF | Device open+close was too fast, never saw EBUSY | `internal/detector/v4l2_linux.go` | Low — /proc walk fixed it |
| 004 | `--detect` flag bound to both `detect` and `list` | Viper `BindPFlag` clash | `cmd/detect.go`, `cmd/list.go` | Medium — same flag name, different commands |
| 009 | JWT tokens in CLI output | `output.Info.Printfln` of config leaked secrets | `internal/output/output.go` | Low — redaction is active |
| 010 | Env vars not expanded in commands | `os.Expand` was not called | `internal/executor/executor.go` | Low — expansion added |
| 012 | Config not loading correctly | Viper `SetConfigFile` vs `AddConfigPath` conflict | `cmd/root.go` | Medium — path precedence |
| 017 | YAML on/off command quoting broken | `yaml.Marshal` unquotes strings | `cmd/onboard.go` | Low — fixed with `yamlQuotedString` |
| 019 | `"all"` cameras string not accepted | `[]string` unmarshal for single string | `cmd/onboard.go` | Low — custom `UnmarshalJSON` |
| 024 | Commands not firing on engine startup | No initial `onChange` call | `internal/engine/engine.go` | Medium — state machine edge case |

**Unresolved edge cases (noted from code analysis):**

- **Signal handling race:** `cmd/detect.go:139-144` sets up a signal handler
  that calls `cancel()`. If `eng.Run(ctx)` has already returned due to an
  internal error, `cancel()` is a no-op — fine. But if the signal arrives
  between `eng.Run(ctx)` returning and `defer cancel()` on line 136 running,
  behavior is correct.
- **Empty debounce from CLI:** If `--debounce` is set to an invalid value
  (e.g., `--debounce=abc`), Viper returns 0, and the engine's `debounceTarget`
  becomes 0. The debounce check `s.debounceCount >= s.debounceTarget` means
  every poll where `status.On != s.current` instantly fires — effectively
  disabling debounce. No error is reported.
- **Context leak in executor:** `cmd/onboard.go` passes `context.Background()`
  to `exec.ExecOn`/`exec.ExecOff` (lines 107, 121). These contexts are never
  cancelled — goroutines spawned for command execution will run until the
  command timeout (default 30s). During shutdown (`cancel()` on line 143),
  the `ctx` is cancelled, but the command goroutines use a new
  `context.Background()`, so they are NOT cancelled on shutdown.

---

## 9. Dependency Health

### Direct Dependencies

| Dependency | Version | Age | Notes |
|-----------|---------|-----|-------|
| `github.com/kardianos/service` | v1.2.4 | Released 2025-05-13 | Stable, widely used. Low maintenance burden. |
| `github.com/pterm/pterm` | v0.12.83 | Recent | Active development. Breaking changes unlikely in minor bumps. |
| `github.com/spf13/cobra` | v1.10.2 | Recent | Standard. Well-maintained. |
| `github.com/spf13/viper` | v1.21.0 | Recent | Stable release. v1.x. |
| `golang.org/x/sys` | v0.45.0 | Recent | Go standard extended library. |

### Transitive Dependencies of Note

| Dependency | Version | Notes |
|-----------|---------|-------|
| `github.com/charmbracelet/huh` | v1.0.0 | Major v1.0.0 — stable. Used for interactive wizard. |
| `gopkg.in/yaml.v3` | v3.0.1 | Stable. No changes expected. |
| `github.com/fsnotify/fsnotify` | v1.9.0 | Pulled by Viper for config file watching. Viper registers a watcher even though config file watching is not used by this tool. |
| `github.com/charmbracelet/bubbletea` | v1.3.6 | Underlying TUI framework for `huh`. Actively developed. |

**Concerns:**
- `fsnotify` is pulled by Viper but the tool does not use config file
  watching. This adds a small amount of binary size and a goroutine for
  each config file location. Not a security concern, but worth noting.
- `charmbracelet/huh v1.0.0` was released recently (likely Q2 2026). The
  import path indirection through `bubbles` → `bubbletea` → `x/term`
  creates a deep dependency tree.
- No known vulnerabilities in any dependency. The `go.sum` signatures
  are verified at build time.

### Go Version

As noted in Section 5.3: `go.mod` declares `go 1.25.0` but GitHub Actions
uses `go-version: "1.22"`. Run `go mod tidy` on 1.25 to update `go.mod`
and `go.sum` to match the actual build environment.

---

## 10. Testing Gaps

| Package | Tests | Coverage Notes |
|---------|-------|----------------|
| `internal/config` | 1 test, 3 assertions | Covers defaults only. No tests for config file loading or Viper binding. |
| `internal/detector` | 6 tests | Tests factory (`New`) and interface conformance. **No tests for actual V4L2 or lsof detection** (requires real hardware). Also no tests exercising the `detect.go:list.go` dual-bind scenario. |
| `internal/engine` | 6 tests | Good coverage of debounce, hotplug, camera filter, graceful shutdown. Uses mocks. |
| `internal/executor` | 6 tests | Covers execution, timeout, template substitution, overlap prevention. Missing: env file parsing tests, shell injection escaping tests. |
| `internal/output` | 2 tests | Minimal. Tests Init() calls only. **No tests for `RedactSecrets`**, `Table`, or `Banner`. |
| `cmd/*` | **0 tests** | No tests for any command: `detect`, `list`, `onboard`, `install`, `uninstall`, `start`, `stop`, `restart`, `service`. These contain the majority of orchestration logic. |

**Gaps:**
- **No `RedactSecrets` tests** — the primary security feature has no test
  coverage.
- **No executor env file tests** — the `parseEnvFile` function is untested.
- **No V4L2 backend tests** — the core detection logic cannot be unit-tested
  without `/dev/video*` devices. Consider adding a `V4L2TestDetector` or
  using build tags with a mock ioctl server.
- **No cmd package tests** — the Cobra command `RunE` functions are largely
  untested. These are where config wiring, error handling, and edge cases
  live.

---

## Summary of Recommended Actions

| Priority | Area | Action | Effort |
|----------|------|--------|--------|
| **CRITICAL** | Security — Shell injection | Escape template values or avoid `sh -c` | Medium |
| **HIGH** | Security — Temp file race | Use `os.MkdirTemp` or stdin pipe for sudo handoff | Small |
| **HIGH** | Cross-platform — macOS | Implement lsof_darwin.go or AVFoundation detector | Medium-Large |
| **HIGH** | CI — Go version mismatch | Align `.github/workflows/release.yml` go-version with `go.mod` | Small |
| **HIGH** | Fragile — /proc walk | Replace per-poll /proc scan with EBUSY retry logic | Medium |
| **HIGH** | Testing — No cmd tests | Write integration tests for detect/list/onboard commands | Large |
| **HIGH** | Testing — Security features | Add `RedactSecrets` test, env file parser test | Small |
| **MEDIUM** | Validation — debounce flag | Add explicit debounce range validation | Small |
| **MEDIUM** | Validation — template pre-check | Add `--validate` mode or pre-parse templates | Small |
| **MEDIUM** | Validation — CameraID path slice | Use `strings.TrimPrefix` instead of hardcoded `[5:]` | Trivial |
| **MEDIUM** | Fragile — env file quoting | Document or fix shell quoting in env file parser | Small |
| **MEDIUM** | Fragile — goroutine management | Add WaitGroup + bounded goroutine pool | Small |
| **LOW** | Tech debt — onboard.go | Extract wizard RunE into helper functions | Medium |
| **LOW** | Fragile — V4L2 raw syscall | Migrate to `unix.Ioctl` wrapper | Small |
| **LOW** | Security — general secret redaction | Expand regex beyond JWT tokens | Small |
| **LOW** | Deps — fsnotify overhead | Viper config file watch is unused but active | Trivial |
