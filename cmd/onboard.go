package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
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

type onboardConfig struct {
	Cameras  []string `json:"cameras"`
	Method   string   `json:"method"`
	Debounce int      `json:"debounce"`
	Interval string   `json:"interval"`
	OnCmd    string   `json:"on-cmd"`
	OffCmd   string   `json:"off-cmd"`
}

type writeConfig struct {
	Camera       string           `yaml:"camera,omitempty"`
	Interval     string           `yaml:"interval"`
	DetectMethod string           `yaml:"detect-method"`
	Debounce     int              `yaml:"debounce"`
	OnCmd        yamlQuotedString `yaml:"on-command"`
	OffCmd       yamlQuotedString `yaml:"off-command"`
}

var (
	onboardDryRun bool
	onboardApply  string
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Guided camera monitor setup",
	Long: `Interactive wizard that walks through camera selection,
detection method configuration with live testing, and
automatic service installation.

Run without flags for the full interactive setup.
Use --dry-run to preview the config before installing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if onboardApply != "" {
			if os.Geteuid() != 0 {
				return fmt.Errorf("root privileges required — re-run with sudo: sudo on-a-meet onboard --apply %s", onboardApply)
			}

			data, err := os.ReadFile(onboardApply)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			var cfg onboardConfig
			if err := json.Unmarshal(data, &cfg); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}

			camera := ""
			if len(cfg.Cameras) == 1 {
				camera = cfg.Cameras[0]
			}

			wc := writeConfig{
				Camera:       camera,
				Interval:     cfg.Interval,
				DetectMethod: cfg.Method,
				Debounce:     cfg.Debounce,
				OnCmd:        yamlQuotedString(cfg.OnCmd),
				OffCmd:       yamlQuotedString(cfg.OffCmd),
			}

			yamlData, err := yaml.Marshal(&wc)
			if err != nil {
				return fmt.Errorf("failed to marshal yaml: %w", err)
			}

			configDir := "/etc/on-a-meet"
			if err := os.MkdirAll(configDir, 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			configPath := configDir + "/config.yaml"
			if _, err := os.Stat(configPath); err == nil {
				var overwrite bool
				huh.NewConfirm().
					Title("Config already exists").
					Description("The configuration file already exists. Overwrite?").
					Affirmative("Overwrite").
					Negative("Keep existing").
					Value(&overwrite).Run()
				if !overwrite {
					output.Info.Printfln("Keeping existing config at %s", configPath)
					return nil
				}
			}
			if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
			output.Success.Printfln("Config written to %s", configPath)

			if err := installService(); err != nil {
				return err
			}

			output.Success.Printfln("Setup complete! Config: %s", configPath)
			return nil
		}

		det, err := detector.New("v4l2")
		if err != nil {
			return fmt.Errorf("failed to create detector: %w", err)
		}

		devices, err := det.ListDevices()
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}
		if len(devices) == 0 {
			output.Warning.Println("No camera devices detected.")
			output.Info.Println("Connect a camera and try again.")
			os.Exit(1)
		}

		output.Banner(len(devices))
		for _, d := range devices {
			output.Info.Printfln("  %s — %s (driver: %s)", d.Path, d.Card, d.Driver)
		}

		deviceOpts := make([]huh.Option[string], 0, len(devices))
		for _, d := range devices {
			deviceOpts = append(deviceOpts, huh.NewOption(d.Path, d.Path))
		}

		var camSelections []string
		var method string
		var debounceStr string
		var intervalStr string
		var cameraChoice string
		var onCmdStr string
		var offCmdStr string

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
					Placeholder("2").
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
					Placeholder("1s").
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
				output.Warning.Println("No cameras selected. Exiting.")
				os.Exit(1)
			}
			cameras = camSelections
		}

		debounce := 2
		if debounceStr != "" {
			v, _ := strconv.Atoi(debounceStr)
			if v > 0 {
				debounce = v
			}
		}

		interval := intervalStr
		if interval == "" {
			interval = "1s"
		}

		method = "v4l2"
		var runTest bool
		huh.NewConfirm().
			Title("Run detection test?").
			Description("Verify that your camera detection is working correctly.\nYou'll be asked to turn your camera ON, then OFF.").
			Affirmative("Test").
			Negative("Skip").
			Value(&runTest).Run()

		if runTest {
			testDet, err := detector.New(method)
			if err == nil {
				testFailed := false
				reader := bufio.NewReader(os.Stdin)

				for {
					output.Warning.Println("Enable your camera now (open a video app), then press Enter to test detection...")
					reader.ReadString('\n')

					anyOn := false
					for _, cam := range cameras {
						status, err := testDet.Detect(cam)
						if err != nil {
							output.Warning.Printfln("  %s: detection error: %v", cam, err)
							continue
						}
						stateStr := "OFF"
						if status.On {
							stateStr = "ON"
							anyOn = true
						}
						output.Info.Printfln("  %s ⟶ %s", cam, stateStr)
					}

					if !anyOn {
						output.Warning.Println("No cameras detected as ON. Try again?")
						var retry bool
						err := huh.NewConfirm().
							Title("Detection test").
							Description("At least one camera was expected to show ON but none detected. Retry?").
							Affirmative("Retry").
							Negative("Skip test").
							Value(&retry).Run()
						if err != nil || !retry {
							testFailed = true
							break
						}
						continue
					}
					break
				}

				for {
					output.Warning.Println("Disable your camera now (close the video app), then press Enter to confirm detection...")
					reader.ReadString('\n')

					allOK := true
					for _, cam := range cameras {
						status, err := testDet.Detect(cam)
						if err != nil {
							output.Warning.Printfln("  %s: detection error: %v", cam, err)
							allOK = false
							continue
						}
						stateStr := "OFF"
						if status.On {
							stateStr = "ON"
						}
						output.Info.Printfln("  %s ⟶ %s", cam, stateStr)
						if status.On {
							allOK = false
						}
					}

					if !allOK {
						output.Warning.Println("Some cameras were not detected as OFF. Try again?")
						var retryOff bool
						err := huh.NewConfirm().
							Title("Off detection test").
							Description("Camera was expected to show OFF but is still detected as ON. Retry?").
							Affirmative("Retry").
							Negative("Skip test").
							Value(&retryOff).Run()
						if err != nil || !retryOff {
							testFailed = true
							break
						}
						continue
					}
					break
				}

				if testFailed {
					var keepV4L2 bool
					huh.NewConfirm().
						Title("Detection method: V4L2").
						Description("V4L2: Direct kernel syscall (Linux-only, no extra deps)\nLSOF: Uses 'lsof' command (cross-platform)").
						Affirmative("Keep V4L2").
						Negative("Change method").
						Value(&keepV4L2).Run()
					if !keepV4L2 {
						var methodSelect string
						huh.NewSelect[string]().
							Title("Detection method").
							Options(
								huh.NewOption("V4L2 (recommended)", "v4l2"),
								huh.NewOption("lsof", "lsof"),
							).
							Value(&methodSelect).Run()
						if methodSelect != "" {
							method = methodSelect
						}
					}
				}
			}
		}

		cfg := onboardConfig{
			Cameras:  cameras,
			Method:   method,
			Debounce: debounce,
			Interval: interval,
			OnCmd:    onCmdStr,
			OffCmd:   offCmdStr,
		}

		if onboardDryRun {
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
			if len(cfg.Cameras) == 1 {
				fmt.Printf("camera: %s\n", cfg.Cameras[0])
			} else {
				fmt.Println("cameras:")
				for _, cam := range cfg.Cameras {
					fmt.Printf("  - %s\n", cam)
				}
			}
			fmt.Println("\nRun without --dry-run to write config and install the service.")
			return nil
		}

		var confirm bool
		err = huh.NewConfirm().
			Title("Ready to install").
			Description(fmt.Sprintf("Method: %s | Interval: %s | Debounce: %d | Cameras: %d\nThis will write config and install the service (requires sudo).", method, interval, debounce, len(cameras))).
			Affirmative("Install").
			Negative("Abort").
			Value(&confirm).Run()
		if err != nil || !confirm {
			output.Info.Println("Install cancelled.")
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
		if err := sudoCmd.Run(); err != nil {
			return fmt.Errorf("sudo install failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
	onboardCmd.Flags().BoolVar(&onboardDryRun, "dry-run", false, "preview config without installing")
	onboardCmd.Flags().StringVar(&onboardApply, "apply", "", "apply collected config file and install service")
}
