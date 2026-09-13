package detector

import "fmt"

// Method names accepted by New and by the detect-method config key.
const (
	MethodV4L2   = "v4l2"
	MethodLsof   = "lsof"
	MethodDarwin = "darwin"
)

// MethodInfo describes a detection backend for a given operating system.
type MethodInfo struct {
	Name        string
	Label       string
	Description string
	Recommended bool
}

// MethodsFor returns the detection methods that actually work on goos, in
// presentation order.
//
// This is the single source of truth for OS/method pairing: the onboarding
// wizard, the config defaults and validation all read it, so they cannot drift
// apart. Note lsof is deliberately absent on darwin — lsof_linux.go is
// linux-only, so New("lsof") there returns an error stub, and offering it
// would be a trap.
func MethodsFor(goos string) []MethodInfo {
	switch goos {
	case "linux":
		return []MethodInfo{
			{
				Name:        MethodV4L2,
				Label:       "V4L2 (recommended)",
				Description: "Direct kernel ioctl on /dev/video*. No extra dependencies.",
				Recommended: true,
			},
			{
				Name:        MethodLsof,
				Label:       "lsof",
				Description: "Checks which processes hold /dev/video* open. Requires the lsof binary.",
			},
		}
	case "darwin":
		return []MethodInfo{
			{
				Name:        MethodDarwin,
				Label:       "macOS unified log (recommended)",
				Description: "Follows UVCAssistant camera power events via the `log` command. No root required. USB/UVC cameras only.",
				Recommended: true,
			},
		}
	default:
		return nil
	}
}

// DefaultMethod returns the recommended method for goos, or "" if the platform
// has no detection backend.
func DefaultMethod(goos string) string {
	for _, m := range MethodsFor(goos) {
		if m.Recommended {
			return m.Name
		}
	}
	return ""
}

// IsValidMethod reports whether id is usable on goos.
func IsValidMethod(goos, id string) bool {
	for _, m := range MethodsFor(goos) {
		if m.Name == id {
			return true
		}
	}
	return false
}

// MethodNames lists the method names available on goos.
func MethodNames(goos string) []string {
	ms := MethodsFor(goos)
	names := make([]string, 0, len(ms))
	for _, m := range ms {
		names = append(names, m.Name)
	}
	return names
}

// ValidateMethod returns a descriptive error when id cannot be used on goos.
func ValidateMethod(goos, id string) error {
	if IsValidMethod(goos, id) {
		return nil
	}
	names := MethodNames(goos)
	if len(names) == 0 {
		return fmt.Errorf("no camera detection method is available on %s", goos)
	}
	return fmt.Errorf("detection method %q is not available on %s (available: %v)", id, goos, names)
}
