package output

import (
	"io"
	"regexp"

	"github.com/pterm/pterm"
)

var jwtRe = regexp.MustCompile(`ey[A-Za-z0-9_-]{20,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`)

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

func RedactSecrets(s string) string {
	return jwtRe.ReplaceAllString(s, "ey***.***.***")
}
