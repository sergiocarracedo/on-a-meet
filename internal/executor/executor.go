package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

type TemplateData struct {
	CameraID string
	Device   string
	State    string
}

type Executor struct {
	timeout time.Duration
	running sync.Map
}

func New(timeout time.Duration) *Executor {
	return &Executor{timeout: timeout}
}

func (e *Executor) ExecOn(ctx context.Context, cmdStr string, data TemplateData) error {
	return e.exec(ctx, cmdStr, data, "on")
}

func (e *Executor) ExecOff(ctx context.Context, cmdStr string, data TemplateData) error {
	return e.exec(ctx, cmdStr, data, "off")
}

func (e *Executor) exec(ctx context.Context, cmdStr string, data TemplateData, state string) error {
	if _, loaded := e.running.Load(state); loaded {
		return nil
	}

	cmdCtx := ctx
	var cancel context.CancelFunc
	if e.timeout > 0 {
		cmdCtx, cancel = context.WithTimeout(ctx, e.timeout)
	} else {
		cmdCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	e.running.Store(state, cancel)
	defer e.running.Delete(state)

	tmpl, err := template.New("").Parse(cmdStr)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}
	rendered := buf.String()

	cmd := exec.CommandContext(cmdCtx, "sh", "-c", rendered)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Cancel = func() error {
		if cmd.Process != nil && cmd.Process.Pid > 0 {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}
	cmd.WaitDelay = 5 * time.Second

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err = cmd.Run()
	outStr := outBuf.String()

	exitCode := -1
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	output.Debug.Printfln("%s-command: %s", state, rendered)
	output.Debug.Printfln("%s-command output: %s", state, outStr)

	if err != nil {
		if cmdCtx.Err() != nil {
			output.Info.Printfln("%s-command exited with code %d", state, exitCode)
			return fmt.Errorf("command timed out after %v: %s", e.timeout, outStr)
		}
		output.Info.Printfln("%s-command exited with code %d", state, exitCode)
		return fmt.Errorf("command failed (exit code %d): %s", exitCode, outStr)
	}

	output.Info.Printfln("%s-command exited with code %d", state, exitCode)

	return nil
}
