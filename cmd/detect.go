package cmd

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
	"github.com/sergiocarracedo/on-a-meet/internal/engine"
	"github.com/sergiocarracedo/on-a-meet/internal/executor"
	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var (
	detectCamera   string
	detectInterval string
	detectOnCmd    string
	detectOffCmd   string
	detectTimeout  string
	detectMethod   string
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect camera on/off state and execute commands",
	Long: `Continuously monitors camera devices and fires
user-defined commands when camera state changes.

Uses V4L2 by default on Linux, darwin on macOS (lsof also available)
to check camera device status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromViper(cmd.Flags().Lookup("detect").Changed)

		interval, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			return err
		}

		det, err := detector.New(cfg.DetectMethod)
		if err != nil {
			return err
		}
		if c, ok := det.(io.Closer); ok {
			defer c.Close()
		}

		devices, err := det.ListDevices()
		if err != nil {
			output.Error.Println("Failed to enumerate camera devices:", err)
			return err
		}
		if len(devices) == 0 {
			output.Warning.Println("No camera devices detected.")
			printNoDevicesHelp(runtime.GOOS)
			return nil
		}

		output.Banner(len(devices))
		for _, d := range devices {
			output.Info.Printfln("  %s — %s (driver: %s)", d.Path, d.Card, d.Driver)
		}
		output.Info.Printfln("Config: method=%s interval=%s debounce=%d timeout=%s camera=%s on=%s off=%s silent=%t verbose=%t", cfg.DetectMethod, cfg.Interval, cfg.Debounce, cfg.Timeout, cfg.Camera, output.RedactSecrets(cfg.OnCmd), output.RedactSecrets(cfg.OffCmd), cfgSilent, cfgVerbose)

		for _, d := range devices {
			status, err := det.Detect(d.Path)
			if err != nil {
				output.Warning.Printfln("  %s ⟶ unavailable: %v", d.Path, err)
				continue
			}
			stateStr := "OFF"
			if status.On {
				stateStr = "ON"
			}
			output.Info.Printfln("  %s ⟶ %s", d.Path, stateStr)
		}

		timeout, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			return err
		}

		exec := executor.New(timeout)
		if cfg.EnvironmentFile != "" {
			exec.SetEnvFile(cfg.EnvironmentFile)
		}

		eng := engine.New(det,
			engine.WithInterval(interval),
			engine.WithDebounce(cfg.Debounce),
			engine.WithOnChange(func(path string, oldState, newState bool, info detector.DeviceInfo) {
				switch {
				case oldState == newState && !newState:
					output.Info.Printfln("[+] %s detected (%s)", path, info.Driver)
				case oldState == newState && newState:
					output.Warning.Printfln("[-] %s disconnected", path)
				case newState:
					output.Warning.Printfln("%s ⟶ ON  (driver: %s)", path, info.Driver)
					if cfg.OnCmd != "" {
						go func() {
							data := executor.TemplateData{
								CameraID: cameraID(path, info),
								Device:   path,
								State:    "on",
							}
							if err := exec.ExecOn(context.Background(), cfg.OnCmd, data); err != nil {
								output.Warning.Printfln("on-command failed: %v", err)
							}
						}()
					}
				default:
					output.Info.Printfln("%s ⟶ OFF  (driver: %s)", path, info.Driver)
					if cfg.OffCmd != "" {
						go func() {
							data := executor.TemplateData{
								CameraID: cameraID(path, info),
								Device:   path,
								State:    "off",
							}
							if err := exec.ExecOff(context.Background(), cfg.OffCmd, data); err != nil {
								output.Warning.Printfln("off-command failed: %v", err)
							}
						}()
					}
				}
			}),
		)

		if cfg.Camera != "" {
			output.Info.Printfln("Monitoring only: %s", cfg.Camera)
			engine.WithCameraFilter(cfg.Camera)(eng)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			output.Info.Println("Shutting down...")
			cancel()
		}()

		// A cancelled context is how Ctrl-C and SIGTERM arrive; that is a
		// clean shutdown, not a failure worth printing usage for.
		if err := eng.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

	detectCmd.Flags().StringVarP(&detectCamera, "camera", "", "", "target specific camera (e.g., /dev/video0)")
	detectCmd.Flags().StringVarP(&detectInterval, "interval", "i", "1s", "polling interval")
	detectCmd.Flags().StringVarP(&detectOnCmd, "on", "", "", "command to run when camera turns on")
	detectCmd.Flags().StringVarP(&detectOffCmd, "off", "", "", "command to run when camera turns off")
	detectCmd.Flags().StringVarP(&detectTimeout, "timeout", "t", "30s", "command execution timeout (0 for no timeout)")
	detectCmd.Flags().StringVarP(&detectMethod, "detect", "d", "v4l2", "detection method (v4l2, lsof, darwin)")

	viper.BindPFlag("camera", detectCmd.Flags().Lookup("camera"))
	viper.BindPFlag("interval", detectCmd.Flags().Lookup("interval"))
	viper.BindPFlag("on-command", detectCmd.Flags().Lookup("on"))
	viper.BindPFlag("off-command", detectCmd.Flags().Lookup("off"))
	viper.BindPFlag("timeout", detectCmd.Flags().Lookup("timeout"))
}

type detectConfig struct {
	Camera          string
	Interval        string
	Debounce        int
	OnCmd           string
	OffCmd          string
	Timeout         string
	DetectMethod    string
	EnvironmentFile string
}

func configFromViper(detectChanged bool) detectConfig {
	method := detectMethod
	if !detectChanged {
		method = viper.GetString("detect-method")
	}
	return detectConfig{
		Camera:          viper.GetString("camera"),
		Interval:        viper.GetString("interval"),
		Debounce:        viper.GetInt("debounce"),
		OnCmd:           viper.GetString("on-command"),
		OffCmd:          viper.GetString("off-command"),
		Timeout:         viper.GetString("timeout"),
		DetectMethod:    method,
		EnvironmentFile: viper.GetString("environment-file"),
	}
}

// cameraID renders {{.CameraID}}. It prefers the detector-supplied stable ID
// and falls back to trimming the Linux device-node prefix. It must never slice
// the path blindly: on macOS the path is a display name, so a fixed-width trim
// mangles it and panics outright for names shorter than the prefix.
func cameraID(path string, info detector.DeviceInfo) string {
	if info.ID != "" {
		return info.ID
	}
	return strings.TrimPrefix(path, "/dev/")
}
