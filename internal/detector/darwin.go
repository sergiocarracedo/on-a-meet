//go:build darwin

package detector

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// seedLookback is how far back the unified log is searched at startup to
// recover the current state of each camera. Without this, a camera that was
// already streaming before on-a-meet started would read as off until its next
// transition. A tight predicate keeps this well under two seconds.
const seedLookback = "12h"

// streamScanBufferBytes raises bufio.Scanner's 64KB default; ndjson log
// records carry a lot of metadata and a long one would otherwise abort the
// stream mid-run.
const streamScanBufferBytes = 1 << 20

// MacOSDetector tracks camera state on macOS by following UVCAssistant power
// events in the unified log.
//
// It keeps a single long-lived `log stream` process rather than polling `log
// show` on every tick: a poll only sees transitions that land inside its
// window, and each invocation costs roughly as long as the default poll
// interval.
type MacOSDetector struct {
	mu        sync.RWMutex
	state     map[string]bool   // device GUID -> streaming
	seen      map[string]bool   // device GUID -> reported by the live stream
	guidFor   map[string]string // DeviceInfo.Path -> device GUID
	streamErr error

	startOnce sync.Once
	cancel    context.CancelFunc
	done      chan struct{}
}

type profilerCamera struct {
	Name     string `json:"_name"`
	ModelID  string `json:"spcamera_model-id"`
	UniqueID string `json:"spcamera_unique-id"`
}

func NewMacOSDetector() *MacOSDetector {
	return &MacOSDetector{
		state:   make(map[string]bool),
		seen:    make(map[string]bool),
		guidFor: make(map[string]string),
		done:    make(chan struct{}),
	}
}

func (d *MacOSDetector) ListDevices() ([]DeviceInfo, error) {
	out, err := exec.Command("system_profiler", "SPCameraDataType", "-json").Output()
	if err != nil {
		return nil, fmt.Errorf("system_profiler: %w", err)
	}

	var result struct {
		Cameras []profilerCamera `json:"SPCameraDataType"`
	}
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("parse system_profiler output: %w", err)
	}

	devices := make([]DeviceInfo, 0, len(result.Cameras))
	guids := make(map[string]string, len(result.Cameras))
	for _, c := range result.Cameras {
		devices = append(devices, DeviceInfo{
			Path:   c.Name,
			ID:     c.UniqueID,
			Card:   c.Name,
			Driver: c.ModelID,
			Bus:    "macOS",
		})
		// Non-UVC cameras (the built-in ones) have no GUID and are simply
		// absent from the map; Detect reports them as not observable.
		if guid, err := guidFromUniqueID(c.UniqueID); err == nil {
			guids[c.Name] = guid
		}
	}

	d.mu.Lock()
	d.guidFor = guids
	d.mu.Unlock()

	return devices, nil
}

func (d *MacOSDetector) Detect(devicePath string) (DeviceStatus, error) {
	d.startOnce.Do(d.start)

	d.mu.RLock()
	err := d.streamErr
	guid, known := d.guidFor[devicePath]
	d.mu.RUnlock()

	if err != nil {
		return DeviceStatus{}, err
	}

	if !known {
		// Either the device list has not been read yet, or this camera was
		// hot-plugged since the last refresh.
		if _, lerr := d.ListDevices(); lerr != nil {
			return DeviceStatus{}, lerr
		}
		d.mu.RLock()
		guid, known = d.guidFor[devicePath]
		d.mu.RUnlock()
	}

	if !known {
		return DeviceStatus{}, fmt.Errorf("%q: %w", devicePath, ErrDeviceNotObservable)
	}

	d.mu.RLock()
	on := d.state[guid]
	d.mu.RUnlock()

	return DeviceStatus{On: on, CheckedAt: time.Now()}, nil
}

// Close stops the background log stream. It is safe to call more than once.
func (d *MacOSDetector) Close() error {
	d.mu.Lock()
	cancel := d.cancel
	d.cancel = nil
	d.mu.Unlock()

	if cancel == nil {
		return nil
	}
	cancel()
	<-d.done
	return nil
}

func (d *MacOSDetector) start() {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "log", "stream", "--style", "ndjson", "--predicate", logPredicate)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		close(d.done)
		d.setStreamErr(fmt.Errorf("log stream: %w", err))
		return
	}
	if err := cmd.Start(); err != nil {
		cancel()
		close(d.done)
		d.setStreamErr(fmt.Errorf("log stream: %w", err))
		return
	}

	d.mu.Lock()
	d.cancel = cancel
	d.mu.Unlock()

	go func() {
		defer close(d.done)
		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), streamScanBufferBytes)
		for scanner.Scan() {
			d.applyStreamLine(scanner.Bytes())
		}
		werr := cmd.Wait()
		// A cancelled context is an ordinary shutdown, not a failure.
		if ctx.Err() == nil {
			if werr == nil {
				werr = fmt.Errorf("log stream exited unexpectedly")
			}
			d.setStreamErr(fmt.Errorf("log stream: %w", werr))
		}
	}()

	d.seed()
}

// seed recovers current state from log history for any camera the live stream
// has not already reported, so a camera that was already on at startup is
// detected immediately.
func (d *MacOSDetector) seed() {
	out, err := exec.Command("log", "show", "--last", seedLookback, "--style", "ndjson", "--predicate", logPredicate).Output()
	if err != nil {
		// Seeding is best-effort: the live stream still works without it.
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	for guid, on := range lastStateByGUID(out) {
		// The live stream is authoritative; history only fills the gaps.
		if d.seen[guid] {
			continue
		}
		d.state[guid] = on
	}
}

// applyStreamLine records a state change reported by the live log stream.
func (d *MacOSDetector) applyStreamLine(line []byte) {
	ev, ok := parseNDJSONLine(line)
	if !ok {
		return
	}
	d.mu.Lock()
	d.state[ev.GUID] = ev.On
	d.seen[ev.GUID] = true
	d.mu.Unlock()
}

func (d *MacOSDetector) setStreamErr(err error) {
	d.mu.Lock()
	d.streamErr = err
	d.mu.Unlock()
}
