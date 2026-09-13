package cmd

import (
	"errors"
	"io"
	"runtime"

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

		rows := pterm.TableData{
			{"Path", "Driver", "Card", "Bus", "Status"},
		}
		for _, d := range devices {
			// A detector error must not render as "OFF" — that is
			// indistinguishable from a working camera nobody is using.
			var status string
			devStatus, err := det.Detect(d.Path)
			switch {
			case errors.Is(err, detector.ErrDeviceNotObservable):
				status = "not observable"
			case err != nil:
				status = "error"
			case devStatus.On:
				status = "ON"
			default:
				status = "OFF"
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
