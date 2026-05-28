package output

import (
	"os"
	"testing"
)

func TestInitSilent(t *testing.T) {
	Init(true, false)
}

func TestInitVerbose(t *testing.T) {
	if err := os.Unsetenv("PTERM_DEBUG"); err != nil {
		t.Fatal(err)
	}
	Init(false, true)
}
