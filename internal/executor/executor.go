package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	envFile string
}

func (e *Executor) SetEnvFile(path string) {
	e.envFile = path
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

func (e *Executor) parseEnvFile() map[string]string {
	vars := make(map[string]string)
	if e.envFile == "" {
		return vars
	}
	data, err := os.ReadFile(e.envFile)
	if err != nil {
		output.Debug.Printfln("env file %s: %v", e.envFile, err)
		return vars
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return vars
}

func (e *Executor) exec(ctx context.Context, cmdStr string, data TemplateData, state string) error {
	if _, loaded := e.running.Load(state); loaded {
		return nil
	}

	output.Info.Printfln("%s-command executing", state)

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

	fileVars := e.parseEnvFile()
	rendered = os.Expand(rendered, func(key string) string {
		if v, ok := fileVars[key]; ok {
			return v
		}
		return os.Getenv(key)
	})

	cmd := exec.CommandContext(cmdCtx, "sh", "-c", rendered)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if len(fileVars) > 0 {
		cmd.Env = os.Environ()
		for k, v := range fileVars {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

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

	output.Debug.Printfln("%s-command: %s", state, output.RedactSecrets(rendered))
	output.Debug.Printfln("%s-command output: %s", state, output.RedactSecrets(outStr))

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
