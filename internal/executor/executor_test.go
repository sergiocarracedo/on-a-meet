package executor

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestExecOnSuccess(t *testing.T) {
	exe := New(5 * time.Second)
	err := exe.ExecOn(context.Background(), "echo hello", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "on",
	})
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestExecOffSuccess(t *testing.T) {
	exe := New(5 * time.Second)
	err := exe.ExecOff(context.Background(), "echo goodbye", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "off",
	})
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestTemplateSubstitution(t *testing.T) {
	exe := New(5 * time.Second)
	err := exe.ExecOn(context.Background(), "echo {{.State}}-{{.CameraID}}", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "on",
	})
	if err != nil {
		t.Errorf("expected nil error with template substitution, got: %v", err)
	}
}

func TestExecTimeout(t *testing.T) {
	exe := New(50 * time.Millisecond)
	err := exe.ExecOn(context.Background(), "sleep 10", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "on",
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestSameStateSkip(t *testing.T) {
	exe := New(10 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		exe.ExecOn(context.Background(), "sleep 10", TemplateData{
			CameraID: "video0",
			Device:   "/dev/video0",
			State:    "on",
		})
	}()

	time.Sleep(50 * time.Millisecond)

	err := exe.ExecOn(context.Background(), "echo should-skip", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "on",
	})
	if err != nil {
		t.Errorf("expected nil (skip), got: %v", err)
	}
}

func TestCrossStateAllow(t *testing.T) {
	exe := New(10 * time.Second)

	go func() {
		exe.ExecOn(context.Background(), "sleep 10", TemplateData{
			CameraID: "video0",
			Device:   "/dev/video0",
			State:    "on",
		})
	}()

	time.Sleep(50 * time.Millisecond)

	err := exe.ExecOff(context.Background(), "echo should-run", TemplateData{
		CameraID: "video0",
		Device:   "/dev/video0",
		State:    "off",
	})
	if err != nil {
		t.Errorf("expected nil (cross-state runs), got: %v", err)
	}
}
