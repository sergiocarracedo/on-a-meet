package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/sergiocarracedo/on-a-meet/internal/detector"
	"github.com/sergiocarracedo/on-a-meet/internal/output"
)

var listMethod string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available camera devices",
	Long: `Enumerates /dev/video* devices and shows driver information
and current on/off status for each camera.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		method := listMethod
		if !cmd.Flags().Lookup("detect").Changed {
			method = viper.GetString("detect-method")
		}

		det, err := detector.New(method)
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

		rows := pterm.TableData{
			{"Path", "Driver", "Card", "Bus", "Status"},
		}
		for _, d := range devices {
			status := "OFF"
			devStatus, err := det.Detect(d.Path)
			if err == nil && devStatus.On {
				status = "ON"
			}
			rows = append(rows, []string{d.Path, d.Driver, d.Card, d.Bus, status})
		}

		output.Table(rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listMethod, "detect", "d", "v4l2", "detection method (v4l2, lsof, darwin)")
}
