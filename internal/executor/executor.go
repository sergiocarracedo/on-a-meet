package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	warned  sync.Map
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

// parseEnvFile reads KEY=VALUE pairs from the configured environment file.
//
// Problems here are reported as warnings, not debug output: a configured file
// that cannot be read means every ${VAR} in the user's command silently
// expands to empty, which typically shows up far away as an authentication
// failure rather than as a configuration error.
func (e *Executor) parseEnvFile() map[string]string {
	vars := make(map[string]string)
	if e.envFile == "" {
		return vars
	}
	data, err := os.ReadFile(e.envFile)
	if err != nil {
		e.warnOnce("envfile:"+e.envFile, func() {
			output.Warning.Printfln("environment-file %s could not be read: %v", e.envFile, err)
			output.Warning.Println("  variables from it will expand to empty in your on/off commands")
		})
		return vars
	}

	var malformed []int
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Most often a long value (an API token) wrapped across lines;
			// the continuation is dropped, silently truncating the value.
			malformed = append(malformed, i+1)
			continue
		}
		vars[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}

	if len(malformed) > 0 {
		e.warnOnce("envfile-malformed:"+e.envFile, func() {
			output.Warning.Printfln("environment-file %s: ignored %d line(s) with no KEY=VALUE (line %s)",
				e.envFile, len(malformed), joinInts(malformed))
			output.Warning.Println("  if a value wraps across lines, join it onto one line — the rest is discarded")
		})
	}
	return vars
}

func joinInts(ns []int) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ", ")
}

// warnOnce keeps a repeating misconfiguration from flooding the log, since
// commands run on every camera transition.
func (e *Executor) warnOnce(key string, emit func()) {
	if _, seen := e.warned.LoadOrStore(key, true); seen {
		return
	}
	emit()
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
	var undefined []string
	rendered = os.Expand(rendered, func(key string) string {
		if v, ok := fileVars[key]; ok {
			return v
		}
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		undefined = append(undefined, key)
		return ""
	})
	if len(undefined) > 0 {
		// Silently substituting empty here is how a missing token turns into
		// a puzzling 401 from the remote service.
		e.warnOnce("undefined:"+state+":"+strings.Join(undefined, ","), func() {
			output.Warning.Printfln("%s-command: %s is not defined, expanded to empty",
				state, strings.Join(dedupe(undefined), ", "))
		})
	}

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

func dedupe(in []string) []string {
	seen := make(map[string]bool, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
