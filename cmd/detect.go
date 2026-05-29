package cmd

import (
	"context"
	"os"
	"os/signal"
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

Uses V4L2 by default (lsof also available) to check /dev/video* device status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := configFromViper()

		interval, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			return err
		}

		det, err := detector.New(cfg.DetectMethod)
		if err != nil {
			return err
		}

		devices, err := det.ListDevices()
		if err != nil {
			output.Error.Println("Failed to enumerate camera devices:", err)
			return err
		}
		if len(devices) == 0 {
			output.Warning.Println("No camera devices detected.")
			output.Info.Println("Make sure your camera is connected and you have the right permissions.")
			output.Info.Println("Tip: add your user to the 'video' group: sudo usermod -a -G video $USER")
			output.Info.Println("Then log out and back in, or run: newgrp video")
			return nil
		}

		output.Banner(len(devices))
		for _, d := range devices {
			output.Info.Printfln("  %s — %s (driver: %s)", d.Path, d.Card, d.Driver)
		}

		timeout, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			return err
		}

		exec := executor.New(timeout)

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
								CameraID: path[5:],
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
								CameraID: path[5:],
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

		return eng.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

	detectCmd.Flags().StringVarP(&detectCamera, "camera", "", "", "target specific camera (e.g., /dev/video0)")
	detectCmd.Flags().StringVarP(&detectInterval, "interval", "i", "1s", "polling interval")
	detectCmd.Flags().StringVarP(&detectOnCmd, "on", "", "", "command to run when camera turns on")
	detectCmd.Flags().StringVarP(&detectOffCmd, "off", "", "", "command to run when camera turns off")
	detectCmd.Flags().StringVarP(&detectTimeout, "timeout", "t", "30s", "command execution timeout (0 for no timeout)")
	detectCmd.Flags().StringVarP(&detectMethod, "detect", "d", "v4l2", "detection method (v4l2, lsof)")

	viper.BindPFlag("camera", detectCmd.Flags().Lookup("camera"))
	viper.BindPFlag("interval", detectCmd.Flags().Lookup("interval"))
	viper.BindPFlag("on-command", detectCmd.Flags().Lookup("on"))
	viper.BindPFlag("off-command", detectCmd.Flags().Lookup("off"))
	viper.BindPFlag("timeout", detectCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("detect-method", detectCmd.Flags().Lookup("detect"))
}

type detectConfig struct {
	Camera       string
	Interval     string
	Debounce     int
	OnCmd        string
	OffCmd       string
	Timeout      string
	DetectMethod string
}

func configFromViper() detectConfig {
	return detectConfig{
		Camera:       viper.GetString("camera"),
		Interval:     viper.GetString("interval"),
		Debounce:     viper.GetInt("debounce"),
		OnCmd:        viper.GetString("on-command"),
		OffCmd:       viper.GetString("off-command"),
		Timeout:      viper.GetString("timeout"),
		DetectMethod: viper.GetString("detect-method"),
	}
}
