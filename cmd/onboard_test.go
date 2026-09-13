package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
)

func TestNormalizeOnboardConfigDefaults(t *testing.T) {
	for _, goos := range []string{"linux", "darwin"} {
		cfg := onboardConfig{}
		if err := normalizeOnboardConfig(&cfg, goos); err != nil {
			t.Fatalf("%s: unexpected error: %v", goos, err)
		}
		if cfg.Method != detector.DefaultMethod(goos) {
			t.Errorf("%s: method = %q, want %q", goos, cfg.Method, detector.DefaultMethod(goos))
		}
		if cfg.Interval != defaultOnboardInterval {
			t.Errorf("%s: interval = %q", goos, cfg.Interval)
		}
		if cfg.Debounce != defaultOnboardDebounce {
			t.Errorf("%s: debounce = %d", goos, cfg.Debounce)
		}
	}
}

// Defaulting to v4l2 on macOS was the old behaviour and produced a config
// that could never detect anything.
func TestNormalizeOnboardConfigRejectsWrongOSMethod(t *testing.T) {
	cfg := onboardConfig{Method: detector.MethodV4L2}
	err := normalizeOnboardConfig(&cfg, "darwin")
	if err == nil {
		t.Fatal("expected v4l2 to be rejected on darwin")
	}
	if !strings.Contains(err.Error(), detector.MethodDarwin) {
		t.Errorf("error should name the usable method, got: %v", err)
	}

	cfg = onboardConfig{Method: detector.MethodDarwin}
	if err := normalizeOnboardConfig(&cfg, "linux"); err == nil {
		t.Error("expected the darwin method to be rejected on linux")
	}
}

func TestNormalizeOnboardConfigKeepsValidValues(t *testing.T) {
	cfg := onboardConfig{Method: detector.MethodLsof, Interval: "500ms", Debounce: 7}
	if err := normalizeOnboardConfig(&cfg, "linux"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Method != detector.MethodLsof || cfg.Interval != "500ms" || cfg.Debounce != 7 {
		t.Errorf("explicit values were overwritten: %+v", cfg)
	}
}

func TestCameraListUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"bare string", `{"cameras": "/dev/video0"}`, []string{"/dev/video0"}},
		{"array", `{"cameras": ["/dev/video0", "/dev/video1"]}`, []string{"/dev/video0", "/dev/video1"}},
		{"macOS display name", `{"cameras": "Logitech StreamCam"}`, []string{"Logitech StreamCam"}},
		{"empty array", `{"cameras": []}`, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg onboardConfig
			if err := json.Unmarshal([]byte(tt.in), &cfg); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if len(cfg.Cameras) != len(tt.want) {
				t.Fatalf("got %v, want %v", cfg.Cameras, tt.want)
			}
			for i := range tt.want {
				if cfg.Cameras[i] != tt.want[i] {
					t.Errorf("camera %d = %q, want %q", i, cfg.Cameras[i], tt.want[i])
				}
			}
		})
	}
}

func TestToWriteConfigCameraSelection(t *testing.T) {
	tests := []struct {
		name    string
		cameras []string
		want    string
	}{
		{"single camera pins to it", []string{"/dev/video0"}, "/dev/video0"},
		{"the all sentinel means everything", []string{"all"}, ""},
		{"several cameras mean everything", []string{"/dev/video0", "/dev/video1"}, ""},
		{"none means everything", nil, ""},
		{"macOS display name", []string{"Logitech StreamCam"}, "Logitech StreamCam"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toWriteConfig(onboardConfig{Cameras: cameraList(tt.cameras)})
			if got.Camera != tt.want {
				t.Errorf("Camera = %q, want %q", got.Camera, tt.want)
			}
		})
	}
}

// Commands routinely contain ":" and "{{...}}", both of which change meaning
// in unquoted YAML.
func TestWriteConfigQuotesCommands(t *testing.T) {
	wc := toWriteConfig(onboardConfig{
		Method:   detector.MethodDarwin,
		Interval: "1s",
		OnCmd:    `curl -H "Bearer x" http://host:8123 {{.CameraID}}`,
		OffCmd:   "echo off: {{.State}}",
	})
	data, err := yaml.Marshal(wc)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var round writeConfig
	if err := yaml.Unmarshal(data, &round); err != nil {
		t.Fatalf("the generated YAML does not parse: %v\n%s", err, data)
	}
	if string(round.OnCmd) != string(wc.OnCmd) {
		t.Errorf("on-command changed through YAML:\n got %q\nwant %q", round.OnCmd, wc.OnCmd)
	}
	if string(round.OffCmd) != string(wc.OffCmd) {
		t.Errorf("off-command changed through YAML:\n got %q\nwant %q", round.OffCmd, wc.OffCmd)
	}
}

func TestInstallNoteMatchesPrivilegeModel(t *testing.T) {
	if !strings.Contains(installNote("linux"), "sudo") {
		t.Error("linux note should mention sudo")
	}
	if !strings.Contains(installNote("darwin"), "no sudo") {
		t.Errorf("macOS note should say no sudo is needed, got: %s", installNote("darwin"))
	}
}

func TestDetectionTroubleshootingIsOSAppropriate(t *testing.T) {
	darwin := strings.Join(detectionTroubleshooting("darwin"), " ")
	if !strings.Contains(darwin, "built-in") {
		t.Errorf("macOS troubleshooting should call out the built-in camera limit, got: %s", darwin)
	}
	if strings.Contains(darwin, "video' group") {
		t.Error("macOS troubleshooting must not mention the Linux video group")
	}
}

// Being piped is not consent to overwrite an existing config; only --yes is.
func TestConfirmExplicitRequiresTheFlag(t *testing.T) {
	orig := onboardAssumeYes
	defer func() { onboardAssumeYes = orig }()

	// Tests do not run on a terminal, so this exercises the piped case.
	onboardAssumeYes = false
	if confirmExplicit("t", "d", "yes", "no") {
		t.Error("without --yes and without a terminal, a destructive prompt must decline")
	}

	onboardAssumeYes = true
	if !confirmExplicit("t", "d", "yes", "no") {
		t.Error("--yes must approve the prompt")
	}
}

// The ordinary install confirmation should not block automation.
func TestConfirmAutoApprovesWhenNotATerminal(t *testing.T) {
	orig := onboardAssumeYes
	defer func() { onboardAssumeYes = orig }()
	onboardAssumeYes = false

	if !confirm("t", "d", "yes", "no") {
		t.Error("a non-destructive prompt should auto-approve when stdin is not a terminal")
	}
}

func TestToWriteConfigCarriesEnvFileAndVerbose(t *testing.T) {
	got := toWriteConfig(onboardConfig{
		Cameras: cameraList{"all"},
		EnvFile: "/etc/default/on-a-meet",
		Verbose: true,
	})
	if got.EnvFile != "/etc/default/on-a-meet" {
		t.Errorf("EnvFile = %q", got.EnvFile)
	}
	if !got.Verbose {
		t.Error("Verbose should be carried through")
	}
}

// Verbosity is a setup aid, not an everyday default.
func TestVerboseDefaultsOff(t *testing.T) {
	var cfg onboardConfig
	if err := normalizeOnboardConfig(&cfg, "darwin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Verbose {
		t.Error("verbose must default to off")
	}
	if cfg.EnvFile != "" {
		t.Error("environment file must default to unset")
	}
}

// The written YAML must use the key names the loader actually reads.
func TestWriteConfigKeyNames(t *testing.T) {
	data, err := yaml.Marshal(toWriteConfig(onboardConfig{
		Interval: "1s",
		EnvFile:  "/etc/default/on-a-meet",
		Verbose:  true,
	}))
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	for _, key := range []string{"environment-file:", "verbose:"} {
		if !strings.Contains(string(data), key) {
			t.Errorf("generated YAML is missing %q:\n%s", key, data)
		}
	}

	var round map[string]any
	if err := yaml.Unmarshal(data, &round); err != nil {
		t.Fatalf("generated YAML does not parse: %v", err)
	}
	if round["environment-file"] != "/etc/default/on-a-meet" {
		t.Errorf("environment-file round-tripped as %v", round["environment-file"])
	}
	if round["verbose"] != true {
		t.Errorf("verbose round-tripped as %v", round["verbose"])
	}
}

// An unset environment file must not emit an empty key that would override
// a value set elsewhere.
func TestWriteConfigOmitsEmptyEnvFile(t *testing.T) {
	data, err := yaml.Marshal(toWriteConfig(onboardConfig{Interval: "1s"}))
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if strings.Contains(string(data), "environment-file") {
		t.Errorf("empty environment file should be omitted:\n%s", data)
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home directory")
	}
	tests := []struct{ in, want string }{
		{"~/.hasscli.env", filepath.Join(home, ".hasscli.env")},
		{"~", home},
		{"/etc/default/on-a-meet", "/etc/default/on-a-meet"},
		{"relative/path", "relative/path"},
		// Only a leading "~/" is a home reference.
		{"/tmp/~/x", "/tmp/~/x"},
		{"~user/x", "~user/x"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := expandHome(tt.in); got != tt.want {
			t.Errorf("expandHome(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDefaultEnvFilePlaceholderIsOSAppropriate(t *testing.T) {
	if strings.HasPrefix(defaultEnvFilePlaceholder("darwin"), "/etc/") {
		t.Error("macOS placeholder should not point into /etc")
	}
	if !strings.HasPrefix(defaultEnvFilePlaceholder("linux"), "/etc/") {
		t.Error("linux placeholder should point into /etc")
	}
}
