package output

import (
	"io"

	"github.com/pterm/pterm"
)

var (
	Info    = pterm.Info
	Success = pterm.Success
	Warning = pterm.Warning
	Error   = pterm.Error
	Debug   = pterm.Debug
)

func Init(silent, verbose bool) {
	if silent {
		pterm.SetDefaultOutput(io.Discard)
	}
	if verbose {
		pterm.EnableDebugMessages()
	}
}

func Table(data pterm.TableData) {
	if err := pterm.DefaultTable.WithHasHeader(true).WithData(data).Render(); err != nil {
		pterm.Error.Println("Failed to render table:", err)
	}
}

func Banner(deviceCount int) {
	pterm.DefaultSection.Println("on-a-meet — Camera Monitor")
	if deviceCount > 0 {
		pterm.Info.Printfln("Detected %d camera device(s)", deviceCount)
	} else {
		pterm.Warning.Println("No camera devices detected")
	}
}
