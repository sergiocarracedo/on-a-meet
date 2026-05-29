package engine

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

type mockDetector struct {
	mu          sync.Mutex
	devices     []detector.DeviceInfo
	detectResp  map[string]mockDetectResponse
	callCount   int
	detectCalls []string
}

func (m *mockDetector) trackDetect(path string) {
	m.detectCalls = append(m.detectCalls, path)
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

func TestEngineStartup(t *testing.T) {
	det := &mockDetector{
		devices: []detector.DeviceInfo{
			{Path: "/dev/video0", Driver: "uvcvideo", Card: "Camera 1"},
			{Path: "/dev/video1", Driver: "uvcvideo", Card: "Camera 2"},
		},
		detectResp: make(map[string]mockDetectResponse),
	}
	for _, d := range det.devices {
		det.detectResp[d.Path] = mockDetectResponse{
			status: detector.DeviceStatus{On: false, CheckedAt: time.Now()},
		}
	}

	var onChangeCalls []string
	eng := New(det,
		WithInterval(10*time.Millisecond),
		WithOnChange(func(path string, oldState, newState bool, info detector.DeviceInfo) {
			onChangeCalls = append(onChangeCalls, path)
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	eng.Run(ctx)

	if len(onChangeCalls) != 2 {
		t.Fatalf("expected 2 initial onChange calls on startup, got %d: %v", len(onChangeCalls), onChangeCalls)
	}
	expectedPaths := map[string]bool{"/dev/video0": true, "/dev/video1": true}
	for _, p := range onChangeCalls {
		if !expectedPaths[p] {
			t.Errorf("unexpected path in onChange: %s", p)
		}
	}
}

func TestDebounce(t *testing.T) {
	det := &mockDetector{
		devices: []detector.DeviceInfo{
			{Path: "/dev/video0", Driver: "uvcvideo", Card: "Camera"},
		},
		detectResp: make(map[string]mockDetectResponse),
	}
	det.detectResp["/dev/video0"] = mockDetectResponse{
		status: detector.DeviceStatus{On: false, CheckedAt: time.Now()},
	}

	var onChange fireTracker
	eng := New(det,
		WithInterval(10*time.Millisecond),
		WithDebounce(3),
		WithOnChange(func(path string, oldState, newState bool, info detector.DeviceInfo) {
			onChange.record(path, oldState, newState, info)
		}),
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(15 * time.Millisecond)
		det.mu.Lock()
		det.detectResp["/dev/video0"] = mockDetectResponse{
			status: detector.DeviceStatus{On: true, CheckedAt: time.Now()},
		}
		det.mu.Unlock()

		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	eng.Run(ctx)

	onChange.mu.Lock()
	fires := len(onChange.calls)
	onChange.mu.Unlock()

	if fires != 2 {
		t.Errorf("expected 2 onChange fires (1 initial + 1 after debounce), got %d", fires)
	}

	onChange.mu.Lock()
	call := onChange.calls[1]
	onChange.mu.Unlock()
	if call.newState != true {
		t.Errorf("expected newState=true, got %v", call.newState)
	}
	if call.oldState != false {
		t.Errorf("expected oldState=false, got %v", call.oldState)
	}
}

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

func TestHotplugAdd(t *testing.T) {
	det := &mockDetector{
		devices: []detector.DeviceInfo{
			{Path: "/dev/video0", Driver: "uvcvideo", Card: "Camera 1"},
		},
		detectResp: make(map[string]mockDetectResponse),
	}
	det.detectResp["/dev/video0"] = mockDetectResponse{
		status: detector.DeviceStatus{On: false, CheckedAt: time.Now()},
	}

	var fires fireTracker
	eng := New(det,
		WithInterval(10*time.Millisecond),
		WithOnChange(fires.record),
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(15 * time.Millisecond)
		det.mu.Lock()
		det.devices = append(det.devices, detector.DeviceInfo{
			Path: "/dev/video2", Driver: "uvcvideo", Card: "Camera 2",
		})
		det.detectResp["/dev/video2"] = mockDetectResponse{
			status: detector.DeviceStatus{On: false, CheckedAt: time.Now()},
		}
		det.mu.Unlock()

		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	eng.Run(ctx)

	fires.mu.Lock()
	hadAdd := false
	for _, c := range fires.calls {
		if c.path == "/dev/video2" && c.oldState == false && c.newState == false {
			hadAdd = true
		}
	}
	fires.mu.Unlock()

	if !hadAdd {
		t.Errorf("expected hotplug add event for /dev/video2")
	}
}

func TestCameraFilter(t *testing.T) {
	det := &mockDetector{
		devices: []detector.DeviceInfo{
			{Path: "/dev/video0", Driver: "uvcvideo", Card: "Camera 1"},
			{Path: "/dev/video1", Driver: "uvcvideo", Card: "Camera 2"},
			{Path: "/dev/video2", Driver: "uvcvideo", Card: "Camera 3"},
		},
		detectResp: make(map[string]mockDetectResponse),
	}
	for _, d := range det.devices {
		det.detectResp[d.Path] = mockDetectResponse{
			status: detector.DeviceStatus{On: false, CheckedAt: time.Now()},
		}
	}

	eng := New(det,
		WithInterval(10*time.Millisecond),
		WithCameraFilter("/dev/video1"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	eng.Run(ctx)

	det.mu.Lock()
	paths := make([]string, len(det.detectCalls))
	copy(paths, det.detectCalls)
	det.mu.Unlock()

	for _, p := range paths {
		if p != "/dev/video1" {
			t.Errorf("Detect called for unfiltered device: %s", p)
		}
	}
}

func TestGracefulShutdown(t *testing.T) {
	det := &mockDetector{
		devices:    []detector.DeviceInfo{},
		detectResp: make(map[string]mockDetectResponse),
	}

	eng := New(det, WithInterval(1*time.Second))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := eng.Run(ctx)
	if err == nil {
		t.Error("expected error on cancelled context, got nil")
	}
}
