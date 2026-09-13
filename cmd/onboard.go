package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

type yamlQuotedString string

func (s yamlQuotedString) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: string(s),
		Style: yaml.DoubleQuotedStyle,
	}, nil
}

type cameraList []string

func (c *cameraList) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = []string{s}
		return nil
	}
	return json.Unmarshal(data, (*[]string)(c))
}

type onboardConfig struct {
	Cameras  cameraList `json:"cameras"`
	Method   string     `json:"method"`
	Debounce int        `json:"debounce"`
	Interval string     `json:"interval"`
	OnCmd    string     `json:"on-cmd"`
	OffCmd   string     `json:"off-cmd"`
	EnvFile  string     `json:"environment-file,omitempty"`
	Verbose  bool       `json:"verbose,omitempty"`
}

type writeConfig struct {
	Camera       string           `yaml:"camera,omitempty"`
	Interval     string           `yaml:"interval"`
	DetectMethod string           `yaml:"detect-method"`
	Debounce     int              `yaml:"debounce"`
	OnCmd        yamlQuotedString `yaml:"on-command"`
	OffCmd       yamlQuotedString `yaml:"off-command"`
	EnvFile      string           `yaml:"environment-file,omitempty"`
	Verbose      bool             `yaml:"verbose"`
}

var (
	onboardDryRun     bool
	onboardApply      string
	onboardConfigFile string
	onboardAssumeYes  bool
)

const (
	defaultOnboardInterval = "1s"
	defaultOnboardDebounce = 2
)

// normalizeOnboardConfig fills in defaults and rejects a method that cannot
// work on this OS. Shared by --apply and --config so the two cannot drift.
func normalizeOnboardConfig(cfg *onboardConfig, goos string) error {
	if cfg.Method == "" {
		cfg.Method = detector.DefaultMethod(goos)
	}
	if err := detector.ValidateMethod(goos, cfg.Method); err != nil {
		return err
	}
	if cfg.Interval == "" {
		cfg.Interval = defaultOnboardInterval
	}
	if cfg.Debounce == 0 {
		cfg.Debounce = defaultOnboardDebounce
	}
	return nil
}

func toWriteConfig(cfg onboardConfig) writeConfig {
	// A single named camera pins monitoring to it; "all" or several cameras
	// leave the key empty, which means "monitor everything".
	camera := ""
	if len(cfg.Cameras) == 1 && cfg.Cameras[0] != "all" {
		camera = cfg.Cameras[0]
	}
	return writeConfig{
		Camera:       camera,
		Interval:     cfg.Interval,
		DetectMethod: cfg.Method,
		Debounce:     cfg.Debounce,
		OnCmd:        yamlQuotedString(cfg.OnCmd),
		OffCmd:       yamlQuotedString(cfg.OffCmd),
		EnvFile:      cfg.EnvFile,
		Verbose:      cfg.Verbose,
	}
}

func printConfigPreview(cfg onboardConfig) {
	output.Success.Println("Configuration preview:")
	fmt.Printf("detect-method: %s\n", cfg.Method)
	fmt.Printf("interval: %s\n", cfg.Interval)
	fmt.Printf("debounce: %d\n", cfg.Debounce)
	if cfg.OnCmd != "" {
		fmt.Printf("on-command: %s\n", cfg.OnCmd)
	}
	if cfg.OffCmd != "" {
		fmt.Printf("off-command: %s\n", cfg.OffCmd)
	}
	if cfg.EnvFile != "" {
		fmt.Printf("environment-file: %s\n", cfg.EnvFile)
	}
	fmt.Printf("verbose: %t\n", cfg.Verbose)
	switch {
	case len(cfg.Cameras) == 1 && cfg.Cameras[0] != "all":
		fmt.Printf("camera: %s\n", cfg.Cameras[0])
	case len(cfg.Cameras) == 0:
		fmt.Println("cameras: all")
	default:
		fmt.Println("cameras:")
		for _, cam := range cfg.Cameras {
			fmt.Printf("  - %s\n", cam)
		}
	}
	fmt.Println("\nRun without --dry-run to write config and install the service.")
}

// stdinIsTerminal reports whether prompts can actually be answered. Without
// this check --config and --apply are not really non-interactive: they block
// on a confirmation nobody can see.
//
// This needs a real tty test rather than a character-device check, because
// /dev/null is itself a character device and would otherwise look interactive.
func stdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func assumeYes() bool {
	return onboardAssumeYes || !stdinIsTerminal()
}

// confirmExplicit is for prompts where defaulting to "yes" would destroy
// something — overwriting an existing config. Being piped is not consent, so
// only the explicit --yes flag skips it; without a terminal it declines.
func confirmExplicit(title, description, affirmative, negative string) bool {
	if onboardAssumeYes {
		return true
	}
	if !stdinIsTerminal() {
		return false
	}
	return promptConfirm(title, description, affirmative, negative)
}

func confirm(title, description, affirmative, negative string) bool {
	if assumeYes() {
		return true
	}
	return promptConfirm(title, description, affirmative, negative)
}

func promptConfirm(title, description, affirmative, negative string) bool {
	var ok bool
	if err := huh.NewConfirm().
		Title(title).
		Description(description).
		Affirmative(affirmative).
		Negative(negative).
		Value(&ok).Run(); err != nil {
		return false
	}
	return ok
}

// writeOnboardConfig writes the YAML config, asking before clobbering an
// existing one.
func writeOnboardConfig(cfg onboardConfig, configPath string) error {
	yamlData, err := yaml.Marshal(toWriteConfig(cfg))
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(configPath); err == nil {
		if !confirmExplicit("Config already exists",
			"The configuration file already exists. Overwrite?",
			"Overwrite", "Keep existing") {
			output.Info.Printfln("Keeping existing config at %s (pass --yes to overwrite)", configPath)
			return nil
		}
	}

	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	output.Success.Printfln("Config written to %s", configPath)
	return nil
}

// finalizeOnboard writes the config and installs the service.
//
// Linux needs root to write /etc and install a systemd unit, so it re-execs
// itself under sudo. macOS writes only inside the user's own home and installs
// a per-user LaunchAgent, so it must stay unprivileged and applies in-process.
func finalizeOnboard(cfg onboardConfig, goos string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolving home directory: %w", err)
	}

	if !needsSudoReexec(goos) {
		configPath := onboardConfigPath(goos, home)
		if err := writeOnboardConfig(cfg, configPath); err != nil {
			return err
		}
		if err := installService(); err != nil {
			return err
		}
		output.Success.Printfln("Setup complete! Config: %s", configPath)
		return nil
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	tmpPath := "/tmp/on-a-meet-onboard.json"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temp config: %w", err)
	}

	binary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to find binary path: %w", err)
	}

	output.Info.Printfln("Running with elevated privileges to complete setup...")
	sudoCmd := exec.Command("sudo", binary, "onboard", "--apply", tmpPath)
	sudoCmd.Stdout = os.Stdout
	sudoCmd.Stderr = os.Stderr
	sudoCmd.Stdin = os.Stdin
	if err := sudoCmd.Run(); err != nil {
		return fmt.Errorf("sudo install failed: %w", err)
	}
	return nil
}

func readOnboardFile(path string) (onboardConfig, error) {
	var cfg onboardConfig
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}
	if err := normalizeOnboardConfig(&cfg, runtime.GOOS); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// runOnboardApply is the privileged primitive: write the config and install
// the service from an already-prepared JSON file.
func runOnboardApply(path string) error {
	if err := requirePrivileges(runtime.GOOS, os.Geteuid(), "onboard --apply "+path); err != nil {
		return err
	}

	cfg, err := readOnboardFile(path)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolving home directory: %w", err)
	}

	configPath := onboardConfigPath(runtime.GOOS, home)
	if err := writeOnboardConfig(cfg, configPath); err != nil {
		return err
	}
	if err := installService(); err != nil {
		return err
	}
	output.Success.Printfln("Setup complete! Config: %s", configPath)
	return nil
}

// runOnboardFromFile is the non-interactive wizard replacement.
func runOnboardFromFile(path string) error {
	cfg, err := readOnboardFile(path)
	if err != nil {
		return err
	}

	if onboardDryRun {
		printConfigPreview(cfg)
		return nil
	}

	if !confirm("Ready to install",
		fmt.Sprintf("Method: %s | Interval: %s | Debounce: %d | Cameras: %d\n%s",
			cfg.Method, cfg.Interval, cfg.Debounce, len(cfg.Cameras), installNote(runtime.GOOS)),
		"Install", "Abort") {
		output.Info.Println("Install cancelled.")
		return nil
	}

	return finalizeOnboard(cfg, runtime.GOOS)
}

func installNote(goos string) string {
	if needsSudoReexec(goos) {
		return "This will write config and install the service (requires sudo)."
	}
	return "This will write config to your home directory and install a per-user LaunchAgent (no sudo)."
}

// chooseDetectionMethod presents the methods that work on this OS. With only
// one available it reports the choice rather than rendering a single-option
// prompt.
func chooseDetectionMethod(goos string) (string, error) {
	methods := detector.MethodsFor(goos)
	if len(methods) == 0 {
		return "", fmt.Errorf("no camera detection method is available on %s", goos)
	}

	method := detector.DefaultMethod(goos)
	if len(methods) == 1 {
		output.Info.Printfln("Detection method: %s", methods[0].Label)
		output.Info.Printfln("  %s", methods[0].Description)
		return method, nil
	}

	if err := huh.NewSelect[string]().
		Title("Detection method").
		Description(detectionMethodHelp(methods)).
		Options(detectionMethodOptions(methods)...).
		Value(&method).Run(); err != nil {
		return "", fmt.Errorf("method selection cancelled: %w", err)
	}
	return method, nil
}

// reconsiderDetectionMethod runs after a failed live test. On a platform with
// a single backend there is nothing to switch to, so it offers troubleshooting
// instead of a pointless picker.
func reconsiderDetectionMethod(goos, current string) string {
	methods := detector.MethodsFor(goos)
	if len(methods) <= 1 {
		output.Warning.Println("Detection did not work, and this platform has only one detection method.")
		for _, line := range detectionTroubleshooting(goos) {
			output.Info.Printfln("  %s", line)
		}
		return current
	}

	choice := current
	if err := huh.NewSelect[string]().
		Title("Detection failed — try another method?").
		Description(detectionMethodHelp(methods)).
		Options(append(detectionMethodOptions(methods),
			huh.NewOption("Keep "+current, current))...).
		Value(&choice).Run(); err != nil {
		return current
	}
	return choice
}

func detectionTroubleshooting(goos string) []string {
	if goos == "darwin" {
		return []string{
			"Only USB/UVC cameras report power events to the unified log.",
			"The built-in MacBook camera is not detectable this way and is reported as 'not observable'.",
			"Check the camera really started streaming (the green indicator light is on).",
			"If detection lags, lower `debounce` in the config.",
		}
	}
	return []string{
		"Check that your user can read the camera device (the 'video' group).",
		"Try running the wizard with sudo.",
	}
}

// runDetectionTest walks the user through turning the camera on and off,
// reporting whether the selected method actually sees the change.
func runDetectionTest(method string, cameras []string) bool {
	det, err := detector.New(method)
	if err != nil {
		output.Warning.Printfln("Could not create detector: %v", err)
		return false
	}
	if c, ok := det.(io.Closer); ok {
		defer c.Close()
	}

	reader := bufio.NewReader(os.Stdin)

	probe := func(wantOn bool) bool {
		for {
			if wantOn {
				output.Warning.Println("Enable your camera now (open a video app), then press Enter to test detection...")
			} else {
				output.Warning.Println("Disable your camera now (close the video app), then press Enter to confirm detection...")
			}
			reader.ReadString('\n')

			matched := !wantOn
			for _, cam := range cameras {
				status, err := det.Detect(cam)
				if errors.Is(err, detector.ErrDeviceNotObservable) {
					// Silence here would look exactly like "camera is off".
					output.Warning.Printfln("  %s ⟶ cannot be monitored with the %s method", cam, method)
					matched = false
					continue
				}
				if err != nil {
					output.Warning.Printfln("  %s: detection error: %v", cam, err)
					matched = false
					continue
				}
				stateStr := "OFF"
				if status.On {
					stateStr = "ON"
				}
				output.Info.Printfln("  %s ⟶ %s", cam, stateStr)
				if wantOn && status.On {
					matched = true
				}
				if !wantOn && status.On {
					matched = false
				}
			}

			if matched {
				return true
			}

			want := "ON"
			if !wantOn {
				want = "OFF"
			}
			if !confirm("Detection test",
				fmt.Sprintf("Expected at least one camera to read %s, but none did. Retry?", want),
				"Retry", "Skip test") {
				return false
			}
		}
	}

	if !probe(true) {
		return false
	}
	return probe(false)
}

func runOnboardInteractive() error {
	goos := runtime.GOOS

	// The method decides what ListDevices returns, so it has to be chosen
	// before the camera list is built — not buried behind a failed test.
	method, err := chooseDetectionMethod(goos)
	if err != nil {
		return err
	}

	det, err := detector.New(method)
	if err != nil {
		return fmt.Errorf("failed to create detector: %w", err)
	}
	if c, ok := det.(io.Closer); ok {
		defer c.Close()
	}

	devices, err := det.ListDevices()
	if err != nil {
		return fmt.Errorf("failed to list devices: %w", err)
	}
	if len(devices) == 0 {
		output.Warning.Println("No camera devices detected.")
		printNoDevicesHelp(goos)
		return fmt.Errorf("no camera devices detected")
	}

	output.Banner(len(devices))
	for _, d := range devices {
		output.Info.Printfln("  %s — %s (driver: %s)", d.Path, d.Card, d.Driver)
	}

	deviceOpts := make([]huh.Option[string], 0, len(devices))
	for _, d := range devices {
		deviceOpts = append(deviceOpts, huh.NewOption(d.Path, d.Path))
	}

	var (
		camSelections []string
		debounceStr   string
		intervalStr   string
		cameraChoice  string
		onCmdStr      string
		offCmdStr     string
		useEnvFile    bool
		envFileStr    string
		verbose       bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Camera selection").
				Description("Choose whether to monitor all cameras or select specific ones").
				Options(
					huh.NewOption("Monitor all cameras", "all"),
					huh.NewOption("Choose specific cameras", "select"),
				).
				Value(&cameraChoice),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Cameras to monitor").
				Description("Select which cameras to monitor (Space to toggle)").
				Options(deviceOpts...).
				Value(&camSelections),
		).WithHideFunc(func() bool { return cameraChoice == "all" }),
		huh.NewGroup(
			huh.NewInput().
				Title("Debounce count").
				Description("Required consecutive same-state polls before firing (higher = less false positives)").
				Placeholder(strconv.Itoa(defaultOnboardDebounce)).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					v, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("must be a number")
					}
					if v < 1 {
						return fmt.Errorf("must be at least 1")
					}
					return nil
				}).
				Value(&debounceStr),
			huh.NewInput().
				Title("Poll interval").
				Description("How often to check camera state (e.g., 500ms, 1s, 2s)").
				Placeholder(defaultOnboardInterval).
				Validate(func(s string) error {
					if s == "" {
						return nil
					}
					_, err := time.ParseDuration(s)
					return err
				}).
				Value(&intervalStr),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("ON command").
				Description("Command to run when camera turns ON (optional; supports {{.State}}, {{.Device}}, {{.CameraID}})").
				Placeholder("e.g., echo 'Camera ON'").
				Value(&onCmdStr),
			huh.NewInput().
				Title("OFF command").
				Description("Command to run when camera turns OFF (optional)").
				Placeholder("e.g., echo 'Camera OFF'").
				Value(&offCmdStr),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Use an environment file?").
				Description("Loads KEY=VALUE pairs (API tokens, server URLs) that your\ncommands can reference as $NAME or ${NAME}.").
				Affirmative("Yes").
				Negative("No").
				Value(&useEnvFile),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Environment file path").
				Description("Absolute path, or ~/... for your home directory.").
				Placeholder(defaultEnvFilePlaceholder(runtime.GOOS)).
				Value(&envFileStr),
		).WithHideFunc(func() bool { return !useEnvFile }),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Verbose output?").
				Description("Logs each command and its output. Handy while setting things up,\nbut noisy for everyday use.").
				Affirmative("Enable").
				Negative("Keep off").
				Value(&verbose),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form cancelled: %w", err)
	}

	var cameras []string
	if cameraChoice == "all" {
		cameras = make([]string, len(devices))
		for i, d := range devices {
			cameras[i] = d.Path
		}
	} else {
		if len(camSelections) == 0 {
			return fmt.Errorf("no cameras selected")
		}
		cameras = camSelections
	}

	debounce := defaultOnboardDebounce
	if debounceStr != "" {
		if v, err := strconv.Atoi(debounceStr); err == nil && v > 0 {
			debounce = v
		}
	}

	interval := intervalStr
	if interval == "" {
		interval = defaultOnboardInterval
	}

	if confirm("Run detection test?",
		"Verify that your camera detection is working correctly.\nYou'll be asked to turn your camera ON, then OFF.",
		"Test", "Skip") {
		if !runDetectionTest(method, cameras) {
			method = reconsiderDetectionMethod(goos, method)
		}
	}

	envFile := ""
	if useEnvFile {
		envFile = expandHome(strings.TrimSpace(envFileStr))
		warnIfEnvFileUnusable(envFile)
	}

	cfg := onboardConfig{
		Cameras:  cameraList(cameras),
		Method:   method,
		Debounce: debounce,
		Interval: interval,
		OnCmd:    onCmdStr,
		OffCmd:   offCmdStr,
		EnvFile:  envFile,
		Verbose:  verbose,
	}
	if err := normalizeOnboardConfig(&cfg, goos); err != nil {
		return err
	}

	if onboardDryRun {
		printConfigPreview(cfg)
		return nil
	}

	if !confirm("Ready to install",
		fmt.Sprintf("Method: %s | Interval: %s | Debounce: %d | Cameras: %d\n%s",
			cfg.Method, cfg.Interval, cfg.Debounce, len(cfg.Cameras), installNote(goos)),
		"Install", "Abort") {
		output.Info.Println("Install cancelled.")
		return nil
	}

	return finalizeOnboard(cfg, goos)
}

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Guided camera monitor setup",
	Long: `Interactive wizard that walks through detection method selection,
camera selection with live testing, and automatic service installation.

Run without flags for the full interactive setup.
Use --config <file> to skip the wizard and apply a pre-made JSON config.
Use --dry-run to preview the config before installing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case onboardApply != "":
			return runOnboardApply(onboardApply)
		case onboardConfigFile != "":
			return runOnboardFromFile(onboardConfigFile)
		default:
			return runOnboardInteractive()
		}
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
	onboardCmd.Flags().BoolVar(&onboardDryRun, "dry-run", false, "preview config without installing")
	onboardCmd.Flags().StringVar(&onboardApply, "apply", "", "apply collected config file and install service")
	onboardCmd.Flags().StringVar(&onboardConfigFile, "config", "", "read config from JSON file (skip interactive wizard)")
	onboardCmd.Flags().BoolVarP(&onboardAssumeYes, "yes", "y", false, "assume yes for confirmation prompts (implied when stdin is not a terminal)")
}

func defaultEnvFilePlaceholder(goos string) string {
	if goos == "darwin" {
		return "~/.config/on-a-meet/env"
	}
	return "/etc/default/on-a-meet"
}

// expandHome resolves a leading ~ so a typed "~/.hasscli.env" does not become
// a literal path that silently fails to open later.
func expandHome(path string) string {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	return filepath.Join(home, path[2:])
}

// warnIfEnvFileUnusable checks the path up front, because the consequence of
// getting it wrong shows up much later and much less clearly: every ${VAR}
// expands to empty and the remote service answers 401.
func warnIfEnvFileUnusable(path string) {
	if path == "" {
		output.Warning.Println("No environment file path given — skipping.")
		return
	}
	if !filepath.IsAbs(path) {
		output.Warning.Printfln("%s is a relative path; it will be resolved against the service's working directory.", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		output.Warning.Printfln("Cannot read %s: %v", path, err)
		output.Warning.Println("  Saving it anyway, but variables from it will expand to empty until the file exists.")
		return
	}

	var keys, malformed int
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(strings.TrimPrefix(line, "export "), "=") {
			keys++
		} else {
			malformed++
		}
	}
	if malformed > 0 {
		output.Warning.Printfln("%s has %d line(s) with no KEY=VALUE — they will be ignored.", path, malformed)
		output.Warning.Println("  A long value wrapped across lines is the usual cause; join it onto one line.")
	}
	output.Success.Printfln("Read %s (%d variable(s)).", path, keys)
}
