package detector

import (
	"encoding/json"
	"strings"
	"testing"
)

// Captured verbatim from `log show --style ndjson` on macOS 26 (Darwin 25.6).
const (
	msgPowerOn = `UVCExtensionDevice:0x10001af33 [0xa650a0780] Post PowerLog {
    "VDCAssistant_Device_GUID" = "00000000-0111-4300-046D-000008930000";
    "VDCAssistant_Power_State" = On;
}`
	msgPowerOff = `UVCExtensionDevice:0x10001af33 [0xa650a0780] Post PowerLog {
    "VDCAssistant_Device_GUID" = "00000000-0111-4300-046D-000008930000";
    "VDCAssistant_Power_State" = Off;
}`
	// This message sits next to the real ones and contains "Streaming". The
	// previous implementation matched it and reported a camera that was
	// shutting down as ON.
	msgTeardown = `UVCExtensionStream:0x10001af3a [0xa64d5c000] Teardown Streaming Pipeline`
)

const wantGUID = "00000000-0111-4300-046D-000008930000"

func TestParsePowerEvent(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		wantOK   bool
		wantOn   bool
		wantGUID string
	}{
		{"power on", msgPowerOn, true, true, wantGUID},
		{"power off", msgPowerOff, true, false, wantGUID},
		{"teardown line is not a state change", msgTeardown, false, false, ""},
		{"empty", "", false, false, ""},
		{"guid without state", `"VDCAssistant_Device_GUID" = "00000000-0111-4300-046D-000008930000";`, false, false, ""},
		{"state without guid", `"VDCAssistant_Power_State" = On;`, false, false, ""},
		{"lowercase guid hex is normalized", strings.ReplaceAll(msgPowerOn, "046D", "046d"), true, true, wantGUID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev, ok := parsePowerEvent(tt.msg)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if ev.On != tt.wantOn {
				t.Errorf("On = %v, want %v", ev.On, tt.wantOn)
			}
			if ev.GUID != tt.wantGUID {
				t.Errorf("GUID = %q, want %q", ev.GUID, tt.wantGUID)
			}
		})
	}
}

// The state lives on a different line than the "Post PowerLog" marker. Parsing
// line by line is what broke the original implementation.
func TestParsePowerEventNeedsWholeMessage(t *testing.T) {
	for _, line := range strings.Split(msgPowerOn, "\n") {
		if _, ok := parsePowerEvent(line); ok {
			return // a single line carried both fields; fine
		}
	}
	if _, ok := parsePowerEvent(msgPowerOn); !ok {
		t.Fatal("whole message must parse even though no single line does")
	}
}

func TestGUIDFromUniqueID(t *testing.T) {
	tests := []struct {
		name    string
		uid     string
		want    string
		wantErr bool
	}{
		// Verified against a live log event from this camera.
		{"logitech streamcam", "0x1114300046d0893", wantGUID, false},
		{"uppercase prefix", "0X1114300046D0893", wantGUID, false},
		{"no 0x prefix", "1114300046d0893", wantGUID, false},
		{"surrounding space", "  0x1114300046d0893 ", wantGUID, false},
		// Built-in Apple cameras report UUIDs, which are not UVC devices.
		{"builtin camera uuid", "6C707041-05AC-0011-0006-000000000001", "", true},
		{"desk view uuid", "6C707041-05AC-0011-0006-000000000002", "", true},
		{"empty", "", "", true},
		{"not hex", "zzzz", "", true},
		{"overflows uint64", "0x11111111111111111", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := guidFromUniqueID(tt.uid)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// The GUID derived from system_profiler must match the one the log reports,
// or per-device state lookups silently find nothing.
func TestGUIDFromUniqueIDMatchesLoggedGUID(t *testing.T) {
	derived, err := guidFromUniqueID("0x1114300046d0893")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ev, ok := parsePowerEvent(msgPowerOn)
	if !ok {
		t.Fatal("fixture failed to parse")
	}
	if derived != ev.GUID {
		t.Errorf("derived %q != logged %q", derived, ev.GUID)
	}
}

// ndjsonRecord wraps a message the way `log --style ndjson` emits it.
func ndjsonRecord(msg string) string {
	b, err := json.Marshal(map[string]string{
		"eventMessage": msg,
		"subsystem":    "com.apple.UVCExtension",
		"messageType":  "Default",
	})
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestParseNDJSONLineIgnoresNonEvents(t *testing.T) {
	// `log stream` prints this header before any records.
	header := `Filtering the log data using "composedMessage CONTAINS "Post PowerLog""`
	for _, line := range []string{header, "", "   ", "not json", "[", "]", "{bad json}"} {
		if _, ok := parseNDJSONLine([]byte(line)); ok {
			t.Errorf("line %q should not parse as an event", line)
		}
	}
}

func TestParseNDJSONLineTrailingComma(t *testing.T) {
	if _, ok := parseNDJSONLine([]byte(ndjsonRecord(msgPowerOn) + ",")); !ok {
		t.Error("a record with array punctuation should still parse")
	}
}

// This is the regression test for the reported bug: a camera that was already
// streaming before on-a-meet started must be reported as on, which requires
// replaying log history rather than only watching for new transitions.
func TestLastStateByGUIDRecoversCameraAlreadyOn(t *testing.T) {
	history := strings.Join([]string{
		`Filtering the log data using "x"`,
		ndjsonRecord(msgPowerOn),
		ndjsonRecord(msgTeardown),
		ndjsonRecord(msgPowerOff),
		ndjsonRecord(msgPowerOn), // most recent event wins
	}, "\n")

	states := lastStateByGUID([]byte(history))
	if len(states) != 1 {
		t.Fatalf("expected 1 device, got %d: %v", len(states), states)
	}
	if !states[wantGUID] {
		t.Errorf("camera should be reported ON; got %v", states[wantGUID])
	}
}

func TestLastStateByGUIDUsesMostRecentEvent(t *testing.T) {
	history := strings.Join([]string{
		ndjsonRecord(msgPowerOn),
		ndjsonRecord(msgPowerOff),
	}, "\n")
	if lastStateByGUID([]byte(history))[wantGUID] {
		t.Error("camera should be reported OFF after a trailing Off event")
	}
}

func TestLastStateByGUIDTracksDevicesIndependently(t *testing.T) {
	other := strings.ReplaceAll(msgPowerOff, "046D-000008930000", "AAAA-000011110000")
	otherGUID := "00000000-0111-4300-AAAA-000011110000"

	history := strings.Join([]string{
		ndjsonRecord(msgPowerOn),
		ndjsonRecord(other),
	}, "\n")

	states := lastStateByGUID([]byte(history))
	if !states[wantGUID] {
		t.Error("first camera should be ON")
	}
	if states[otherGUID] {
		t.Errorf("second camera should be OFF; states=%v", states)
	}
}

func TestLastStateByGUIDEmpty(t *testing.T) {
	if got := lastStateByGUID(nil); len(got) != 0 {
		t.Errorf("expected no devices, got %v", got)
	}
}
