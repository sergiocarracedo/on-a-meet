# TESTING.md — Test Structure and Practices

## 1. Test Framework

The project uses the standard Go `testing` package exclusively. No third-party test frameworks (testify, ginkgo, etc.) are used. The project also does not use `testify` for assertions -- all assertions are hand-written with `t.Error`, `t.Errorf`, `t.Fatal`, and `t.Fatalf`.

---

## 2. Test File Locations

Test files are placed **alongside the source files they test**, in the same package directory:

| Source | Test | Package |
|--------|------|---------|
| `internal/detector/detector.go` | `internal/detector/detector_test.go` | `detector` |
| `internal/detector/v4l2_stub.go` | `internal/detector/v4l2_stub_test.go` | `detector` |
| `internal/engine/engine.go` | `internal/engine/engine_test.go` | `engine` |
| `internal/executor/executor.go` | `internal/executor/executor_test.go` | `executor` |
| `internal/config/config.go` | `internal/config/config_test.go` | `config` |
| `internal/output/output.go` | `internal/output/output_test.go` | `output` |

All test files use the same package name as their source (not `_test` suffix for external tests), meaning they have access to unexported types and functions.

---

## 3. How to Run Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/engine/...

# Run tests with verbose output
go test -v ./...

# Run tests with race detector
go test -race ./...

# Run a specific test function
go test -run TestDebounce ./internal/engine/...

# Run tests excluding Linux-only files (on non-Linux)
go test -tags '!linux' ./...
```

There is no Makefile -- use `go test` directly.

---

## 4. Coverage State

Current coverage is **limited** and concentrated in specific areas:

| Package | Test File | Test Functions | Coverage Focus |
|---------|-----------|----------------|----------------|
| `detector` | `detector_test.go` | 5 | Factory: `New()` for valid/invalid methods, compile-time interface checks |
| `detector` | `v4l2_stub_test.go` | 1 | Non-Linux stub: `ListDevices`/`Detect` return errors |
| `engine` | `engine_test.go` | 5 | Startup, debounce, hotplug add, camera filter, graceful shutdown |
| `executor` | `executor_test.go` | 6 | Success (on/off), template substitution, timeout, same-state skip, cross-state allow |
| `config` | `config_test.go` | 1 | `Defaults()` returns expected values |
| `output` | `output_test.go` | 2 | `Init(silent)` and `Init(verbose)` don't panic |

**Notable gaps with no tests:**
- `internal/detector/v4l2_linux.go` -- V4L2 syscall detect logic (Linux-only, requires real hardware)
- `internal/detector/lsof_linux.go` -- Lsof detect logic (Linux-only, requires lsof)
- `cmd/` -- All Cobra commands have zero tests
- `internal/executor/executor.go` -- `parseEnvFile()` and `SetEnvFile()`
- `internal/output/output.go` -- `Table()`, `Banner()`, `RedactSecrets()`

---

## 5. Test Patterns

### 5.1 Mock Pattern (engine_test.go)

The most sophisticated mock is `mockDetector` for the `Detector` interface:

```go
// internal/engine/engine_test.go lines 12-45
type mockDetector struct {
    mu          sync.Mutex
    devices     []detector.DeviceInfo
    detectResp  map[string]mockDetectResponse
    callCount   int
    detectCalls []string
}

type mockDetectResponse struct {
    status detector.DeviceStatus
    err    error
}

func (m *mockDetector) ListDevices() ([]detector.DeviceInfo, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.devices, nil
}

func (m *mockDetector) Detect(devicePath string) (detector.DeviceStatus, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.callCount++
    m.detectCalls = append(m.detectCalls, devicePath)
    resp, ok := m.detectResp[devicePath]
    if ok {
        return resp.status, resp.err
    }
    return detector.DeviceStatus{On: false, CheckedAt: time.Now()}, nil
}
```

**Pattern details:**
- The mock embeds `sync.Mutex` for thread-safe mutation from goroutines (tests are concurrent)
- `detectResp` is a `map[string]mockDetectResponse` keyed by device path
- `detectCalls` slice records which paths were detected (used for filter verification)
- `callCount` tracks total detect invocations
- Default response (if path not in map) returns `{On: false}`

### 5.2 Fire Tracker Pattern (engine_test.go)

A helper struct for recording and inspecting callback invocations:

```go
// internal/engine/engine_test.go lines 140-156
type fireCall struct {
    path     string
    oldState bool
    newState bool
    info     detector.DeviceInfo
}

type fireTracker struct {
    mu    sync.Mutex
    calls []fireCall
}

func (f *fireTracker) record(path string, oldState, newState bool, info detector.DeviceInfo) {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.calls = append(f.calls, fireCall{path, oldState, newState, info})
}
```

This is used in tests by passing `fires.record` as the `OnChange` callback:
```go
eng := New(det,
    WithInterval(10*time.Millisecond),
    WithOnChange(fires.record),
)
```

### 5.3 Time-Based Concurrency Pattern (engine_test.go)

Tests manipulate the mock detector state from a goroutine running concurrently with the engine:

```go
// internal/engine/engine_test.go TestDebounce lines 105-119
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(15 * time.Millisecond)       // let engine start and settle
    det.mu.Lock()
    det.detectResp["/dev/video0"] = mockDetectResponse{
        status: detector.DeviceStatus{On: true, CheckedAt: time.Now()},
    }
    det.mu.Unlock()

    time.Sleep(50 * time.Millisecond)       // long enough for debounce to fire
    cancel()                                 // stop the engine
}()

eng.Run(ctx)
```

**Pattern:**
1. Create `context.WithCancel` (not `WithTimeout` for dynamic control)
2. Start a goroutine that: sleeps -> mutates mock -> sleeps more -> cancels context
3. Call `eng.Run(ctx)` (blocking)  
4. After `Run` returns, inspect captured state

For simpler tests, `context.WithTimeout` is used instead:
```go
// internal/engine/engine_test.go line 69
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()
eng.Run(ctx)
```

### 5.4 Test Structure (No Table-Driven Tests)

All tests are standalone functions. The project does **not** use table-driven tests (`[]struct{...}`). Each test function tests one specific scenario:

```go
func TestEngineStartup(t *testing.T) { ... }
func TestDebounce(t *testing.T) { ... }
func TestHotplugAdd(t *testing.T) { ... }
func TestCameraFilter(t *testing.T) { ... }
func TestGracefulShutdown(t *testing.T) { ... }
```

### 5.5 Assertion Style

No assertion library. Plain Go conditionals:

```go
// Fatalf for unrecoverable failures
if err != nil {
    t.Fatalf("New('v4l2') failed: %v", err)
}

// Errorf/Error for non-fatal failures
if fires != 2 {
    t.Errorf("expected 2 onChange fires ..., got %d", fires)
}

// Simple boolean checks
if d == nil {
    t.Fatal("New('v4l2') returned nil")
}
if !expectedPaths[p] {
    t.Errorf("unexpected path in onChange: %s", p)
}
if err == nil {
    t.Fatal("New('unknown') should return error")
}

// Literal comparison
if cfg.Interval != "1s" {
    t.Errorf("Default interval = %q, want %q", cfg.Interval, "1s")
}
```

---

## 6. Test Details by Package

### 6.1 `internal/detector/detector_test.go`

Tests the `New()` factory function:

| Test | What it verifies |
|------|------------------|
| `TestNewV4L2` | `New("v4l2")` returns non-nil, no error |
| `TestNewLsof` | `New("lsof")` returns non-nil, no error |
| `TestNewUnknown` | `New("unknown")` returns an error |
| `TestNewV4L2ImplementsDetector` | Compile-time check: `var _ Detector = d` |
| `TestNewLsofImplementsDetector` | Compile-time check: `var _ Detector = d` |

### 6.2 `internal/detector/v4l2_stub_test.go`

Tests the non-Linux stub (build tag `!linux`):

| Test | What it verifies |
|------|------------------|
| `TestV4L2StubReturnsError` | `ListDevices()` returns error + nil devices; `Detect()` returns error + `On=false` |

### 6.3 `internal/engine/engine_test.go`

The most comprehensive test file:

| Test | What it verifies | Pattern |
|------|------------------|---------|
| `TestEngineStartup` | Two devices emit two initial `onChange` calls with correct paths | `WithTimeout` + deferred cancel |
| `TestDebounce` | After 3 consecutive ON polls, onChange fires with oldState=false, newState=true | Goroutine mutates mock + `WithCancel` |
| `TestHotplugAdd` | New device appearing mid-cycle emits onChange with oldState=false, newState=false | Goroutine appends to devices slice |
| `TestCameraFilter` | Only `/dev/video1` is detected; other paths never appear in `detectCalls` | `WithTimeout` + inspect mock state |
| `TestGracefulShutdown` | Pre-cancelled context causes `Run()` to return error | `context.WithCancel` + immediate `cancel()` |

### 6.4 `internal/executor/executor_test.go`

| Test | What it verifies |
|------|------------------|
| `TestExecOnSuccess` | `ExecOn("echo hello")` succeeds, no error |
| `TestExecOffSuccess` | `ExecOff("echo goodbye")` succeeds, no error |
| `TestTemplateSubstitution` | `ExecOn("echo {{.State}}-{{.CameraID}}")` with template data succeeds |
| `TestExecTimeout` | `ExecOn("sleep 10")` with 50ms timeout returns error |
| `TestSameStateSkip` | Two concurrent `ExecOn("sleep 10")` calls: the second returns nil (skipped) |
| `TestCrossStateAllow` | Overlapping `ExecOn` + `ExecOff` both execute (different states allowed) |

**Notable test details:**
- `TestSameStateSkip` uses `sync.WaitGroup` to ensure the goroutine starts first, then `time.Sleep(50ms)` before issuing the second call
- `TestCrossStateAllow` demonstrates that on and off command states are tracked independently via `sync.Map`

### 6.5 `internal/config/config_test.go`

| Test | What it verifies |
|------|------------------|
| `TestDefaults` | `Defaults()` returns `Interval="1s"`, `DetectMethod="v4l2"`, `Debounce=3` |

### 6.6 `internal/output/output_test.go`

| Test | What it verifies |
|------|------------------|
| `TestInitSilent` | `Init(true, false)` does not panic |
| `TestInitVerbose` | `Init(false, true)` after `os.Unsetenv("PTERM_DEBUG")` does not panic |

---

## 7. Executor Test Details (Overlap Prevention)

The `sync.Map`-based overlap prevention is tested explicitly:

```
Time:         0ms         50ms          100ms
TestSameStateSkip:
  goroutine:  ExecOn("sleep 10")  ──► (still running)
  main:                          ExecOn("echo skip")  ──► returns nil (skipped)

TestCrossStateAllow:
  goroutine:  ExecOn("sleep 10")  ──► (still running)
  main:                          ExecOff("echo run")  ──► returns nil (executed)
```

The key insight is that `ExecOn` and `ExecOff` track their state independently:
```go
// internal/executor/executor.go
func (e *Executor) ExecOn(...) error { return e.exec(ctx, cmdStr, data, "on") }
func (e *Executor) ExecOff(...) error { return e.exec(ctx, cmdStr, data, "off") }
```

Both call `e.exec()` with different state strings ("on" vs "off"), so they use different keys in `e.running`.

---

## 8. Testing Conventions Summary

| Convention | Practice |
|------------|----------|
| Test framework | Standard `testing` package only |
| Assertion library | None -- hand-written conditionals |
| Table-driven tests | Not used |
| External test package (`_test`) | Not used -- tests are in-package |
| Test helper functions | Minimal; `fireTracker.record` is the only notable helper |
| Mock framework | Hand-written mocks implementing the interface |
| Goroutine safety | Mocks embed `sync.Mutex`; test helpers (`fireTracker`) use `sync.Mutex` |
| Context usage | `context.WithTimeout` for fixed-duration tests; `context.WithCancel` for dynamic goroutine-driven tests |
| `t.Parallel()` | Not used anywhere |
| Test data/files | No test fixture files; all data is inlined |
| CI | No CI test runner configured in `.github/workflows/` |

---

## 9. What NOT to Do in Tests

1. **Do NOT import `testify` or other assertion packages** -- the project uses plain Go conditionals.

2. **Do NOT use table-driven tests unless adding many similar cases** -- the project convention is one function per scenario.

3. **Do NOT use `_test` package suffix** -- all tests are in-package (`package engine`, not `package engine_test`).

4. **Do NOT rely on real hardware** -- use `mockDetector` for engine tests.

5. **Do NOT add tests for Linux-only files that require real devices** -- `v4l2_linux.go` and `lsof_linux.go` have no tests because they require actual camera hardware.

6. **Do NOT add tests to `cmd/` without mocking Viper and Cobra** -- currently no `cmd/` tests exist.

7. **Do NOT use `t.Parallel()`** without ensuring mock thread safety -- but currently no tests use it.

8. **Do NOT forget build tags on test files for platform-specific code** -- `v4l2_stub_test.go` has `//go:build !linux`.
