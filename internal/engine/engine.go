package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

type OnChange func(path string, oldState, newState bool, info detector.DeviceInfo)

type deviceState struct {
	info           detector.DeviceInfo
	current        bool
	previous       bool
	debounceCount  int
	debounceTarget int
}

type Engine struct {
	detector     detector.Detector
	interval     time.Duration
	debounce     int
	cameraFilter string
	onChange     OnChange

	states map[string]*deviceState
	mu     sync.Mutex
	logger *log.Logger
}

type Option func(*Engine)

func WithInterval(d time.Duration) Option {
	return func(e *Engine) { e.interval = d }
}

func WithDebounce(n int) Option {
	return func(e *Engine) { e.debounce = n }
}

func WithCameraFilter(path string) Option {
	return func(e *Engine) { e.cameraFilter = path }
}

func WithOnChange(cb OnChange) Option {
	return func(e *Engine) { e.onChange = cb }
}

func New(det detector.Detector, opts ...Option) *Engine {
	e := &Engine{
		detector: det,
		interval: 1 * time.Second,
		debounce: 3,
		states:   make(map[string]*deviceState),
		logger:   log.New(log.Writer(), "", 0),
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Engine) Run(ctx context.Context) error {
	devices, err := e.detector.ListDevices()
	if err != nil {
		return err
	}

	var filtered []detector.DeviceInfo
	for _, d := range devices {
		if e.cameraFilter != "" && d.Path != e.cameraFilter {
			continue
		}
		filtered = append(filtered, d)
	}

	for _, d := range filtered {
		status, err := e.detector.Detect(d.Path)
		if err != nil {
			continue
		}
		e.mu.Lock()
		e.states[d.Path] = &deviceState{
			info:           d,
			current:        status.On,
			previous:       status.On,
			debounceTarget: e.debounce,
		}
		e.mu.Unlock()
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(e.interval):
			e.pollCycle()
		}
	}
}

func (e *Engine) pollCycle() {
	devices, err := e.detector.ListDevices()
	if err != nil {
		return
	}

	filtered := make(map[string]detector.DeviceInfo)
	for _, d := range devices {
		if e.cameraFilter != "" && d.Path != e.cameraFilter {
			continue
		}
		filtered[d.Path] = d
	}

	type hotplugEvent struct {
		path string
		info detector.DeviceInfo
	}

	var added []hotplugEvent
	var removed []hotplugEvent

	e.mu.Lock()
	for path := range e.states {
		if _, exists := filtered[path]; !exists {
			removed = append(removed, hotplugEvent{path, e.states[path].info})
			delete(e.states, path)
		}
	}
	for path, info := range filtered {
		if _, exists := e.states[path]; !exists {
			e.states[path] = &deviceState{
				info:           info,
				debounceTarget: e.debounce,
			}
			added = append(added, hotplugEvent{path, info})
		}
	}
	e.mu.Unlock()

	for _, a := range added {
		if e.onChange != nil {
			e.onChange(a.path, false, false, a.info)
		}
	}
	for _, r := range removed {
		if e.onChange != nil {
			e.onChange(r.path, true, true, r.info)
		}
	}

	for path := range filtered {
		status, err := e.detector.Detect(path)
		if err != nil {
			continue
		}

		e.mu.Lock()
		s := e.states[path]
		e.mu.Unlock()

		if s == nil {
			continue
		}

		if status.On == s.current {
			s.debounceCount = 0
		} else {
			s.debounceCount++
			if s.debounceCount >= s.debounceTarget {
				oldState := s.current
				s.current = status.On
				s.previous = oldState
				s.debounceCount = 0
				if e.onChange != nil {
					e.onChange(path, oldState, status.On, s.info)
				}
			}
		}
	}
}
