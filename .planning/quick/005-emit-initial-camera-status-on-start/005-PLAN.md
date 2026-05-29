# Quick Task 005 Plan: Emit Initial Camera Status

## Tasks

### Task 1: Print initial camera ON/OFF status after startup banner

**Files:**
- `cmd/detect.go`

**Action:**
After the device banner (line 65 in detect.go, the `output.Info.Printfln("Config: ...")` line) and before creating the executor/engine, call `det.Detect()` for each device and print its initial status.

Add this block right after the config print (line 65):

```go
for _, d := range devices {
    status, err := det.Detect(d.Path)
    if err != nil {
        continue
    }
    stateStr := "OFF"
    if status.On {
        stateStr = "ON"
    }
    output.Info.Printfln("  %s ⟶ %s", d.Path, stateStr)
}
```

Use the same `output.Info.Printfln` pattern, same indentation.

**Verify:**
- `go build ./...` succeeds
- `go vet ./...` passes
- All 19 tests pass (regression)

**Done:**
- `cmd/detect.go` modified with the initial status print block
