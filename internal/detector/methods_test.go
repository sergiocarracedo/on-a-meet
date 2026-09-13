package detector

import (
	"strings"
	"testing"
)

var allOSes = []string{"linux", "darwin", "windows"}

func TestMethodsForPerOS(t *testing.T) {
	linux := MethodNames("linux")
	if len(linux) == 0 {
		t.Fatal("linux must offer at least one method")
	}
	for _, n := range linux {
		if n == MethodDarwin {
			t.Errorf("linux must not offer %q", MethodDarwin)
		}
	}

	darwin := MethodNames("darwin")
	if len(darwin) == 0 {
		t.Fatal("darwin must offer at least one method")
	}
	for _, n := range darwin {
		// lsof_linux.go is linux-only, so offering it on macOS would hand the
		// user a method that can only fail.
		if n == MethodV4L2 || n == MethodLsof {
			t.Errorf("darwin must not offer %q", n)
		}
	}

	if got := MethodsFor("windows"); got != nil {
		t.Errorf("unsupported OS should offer nothing, got %v", got)
	}
}

func TestExactlyOneRecommendedPerOS(t *testing.T) {
	for _, goos := range allOSes {
		ms := MethodsFor(goos)
		if len(ms) == 0 {
			continue
		}
		n := 0
		for _, m := range ms {
			if m.Recommended {
				n++
			}
		}
		if n != 1 {
			t.Errorf("%s: want exactly 1 recommended method, got %d", goos, n)
		}
	}
}

// Guards against the two tables drifting apart.
func TestDefaultMethodIsAlwaysOffered(t *testing.T) {
	for _, goos := range allOSes {
		def := DefaultMethod(goos)
		if len(MethodsFor(goos)) == 0 {
			if def != "" {
				t.Errorf("%s: expected no default, got %q", goos, def)
			}
			continue
		}
		if !IsValidMethod(goos, def) {
			t.Errorf("%s: default %q is not in the method list", goos, def)
		}
	}
}

// Every advertised method must actually be constructible. New is
// OS-independent thanks to the stubs, so a typo in any OS's table is caught
// from every platform.
func TestEveryAdvertisedMethodIsConstructible(t *testing.T) {
	for _, goos := range allOSes {
		for _, name := range MethodNames(goos) {
			det, err := New(name)
			if err != nil {
				t.Errorf("%s: New(%q) failed: %v", goos, name, err)
			}
			if det == nil {
				t.Errorf("%s: New(%q) returned nil", goos, name)
			}
		}
	}
}

func TestIsValidMethodCrossOS(t *testing.T) {
	tests := []struct {
		goos, method string
		want         bool
	}{
		{"linux", MethodV4L2, true},
		{"linux", MethodLsof, true},
		{"linux", MethodDarwin, false},
		{"darwin", MethodDarwin, true},
		{"darwin", MethodV4L2, false},
		{"darwin", MethodLsof, false},
		{"windows", MethodV4L2, false},
		{"linux", "nonsense", false},
	}
	for _, tt := range tests {
		if got := IsValidMethod(tt.goos, tt.method); got != tt.want {
			t.Errorf("IsValidMethod(%q, %q) = %v, want %v", tt.goos, tt.method, got, tt.want)
		}
	}
}

func TestValidateMethodMessages(t *testing.T) {
	if err := ValidateMethod("darwin", MethodDarwin); err != nil {
		t.Errorf("valid pairing should not error: %v", err)
	}
	err := ValidateMethod("darwin", MethodV4L2)
	if err == nil {
		t.Fatal("expected an error for a wrong-OS method")
	}
	// The message must name what the user can actually pick.
	if !strings.Contains(err.Error(), MethodDarwin) {
		t.Errorf("error should list available methods, got: %v", err)
	}
	if err := ValidateMethod("windows", MethodV4L2); err == nil {
		t.Error("expected an error on an unsupported OS")
	}
}
