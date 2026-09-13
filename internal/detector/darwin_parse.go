package detector

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// logPredicate selects the UVCAssistant power-state events used to track
// camera state on macOS. Apple has changed this subsystem between releases
// (see .planning/research/PITFALLS.md), so it is kept in one place.
const logPredicate = `subsystem == "com.apple.UVCExtension" AND composedMessage CONTAINS "Post PowerLog"`

// A matching event looks like this, as a single eventMessage:
//
//	UVCExtensionDevice:0x10001af33 [0xa650a0780] Post PowerLog {
//	    "VDCAssistant_Device_GUID" = "00000000-0111-4300-046D-000008930000";
//	    "VDCAssistant_Power_State" = On;
//	}
//
// The state is on a different line than the "Post PowerLog" marker, which is
// why this must be matched against the whole message rather than line by line.
var (
	guidRe  = regexp.MustCompile(`"VDCAssistant_Device_GUID"\s*=\s*"([0-9A-Fa-f-]+)"`)
	stateRe = regexp.MustCompile(`"VDCAssistant_Power_State"\s*=\s*(On|Off)\b`)
)

// logLine is one record of `log show`/`log stream --style ndjson`.
type logLine struct {
	EventMessage string `json:"eventMessage"`
}

// powerEvent is a decoded UVCAssistant power-state transition.
type powerEvent struct {
	GUID string
	On   bool
}

// parsePowerEvent extracts a device GUID and its new power state from a log
// event message. ok is false unless both were found, so partial or unrelated
// messages (notably "Teardown Streaming Pipeline") are ignored rather than
// guessed at.
func parsePowerEvent(msg string) (ev powerEvent, ok bool) {
	g := guidRe.FindStringSubmatch(msg)
	if g == nil {
		return powerEvent{}, false
	}
	s := stateRe.FindStringSubmatch(msg)
	if s == nil {
		return powerEvent{}, false
	}
	return powerEvent{GUID: normalizeGUID(g[1]), On: s[1] == "On"}, true
}

func normalizeGUID(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}

// guidFromUniqueID converts a system_profiler spcamera_unique-id into the
// device GUID that UVCAssistant logs.
//
// A UVC camera's unique-id is a plain hex integer ("0x1114300046d0893"), which
// zero-padded to 16 digits maps onto the GUID positionally:
//
//	01114300046D0893 -> 00000000-0111-4300-046D-000008930000
//
// Built-in Apple cameras instead report a UUID
// ("6C707041-05AC-0011-0006-000000000001"), which is not a single hex integer.
// So a parse failure here is exactly "this is not a UVC device", and callers
// use it as such rather than needing a separate heuristic.
func guidFromUniqueID(uid string) (string, error) {
	s := strings.TrimSpace(uid)
	s = strings.TrimPrefix(strings.TrimPrefix(s, "0x"), "0X")
	if s == "" {
		return "", fmt.Errorf("empty camera unique-id")
	}
	v, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return "", fmt.Errorf("unique-id %q is not a UVC hex identifier: %w", uid, err)
	}
	u := fmt.Sprintf("%016X", v)
	return fmt.Sprintf("00000000-%s-%s-%s-0000%s0000", u[0:4], u[4:8], u[8:12], u[12:16]), nil
}

// parseNDJSONLine decodes one `log ... --style ndjson` record. It tolerates the
// human-readable header line that `log stream` prints first, blank lines, and
// array punctuation, returning ok=false for anything that is not an event.
func parseNDJSONLine(line []byte) (powerEvent, bool) {
	s := strings.TrimSpace(string(line))
	s = strings.TrimSuffix(s, ",")
	if !strings.HasPrefix(s, "{") {
		return powerEvent{}, false
	}
	var rec logLine
	if err := json.Unmarshal([]byte(s), &rec); err != nil {
		return powerEvent{}, false
	}
	return parsePowerEvent(rec.EventMessage)
}

// lastStateByGUID replays a block of ndjson log records and returns the final
// power state of each device.
//
// This is what lets on-a-meet know a camera is already streaming when it
// starts: the unified log only reports transitions, so current state has to be
// reconstructed from the most recent event per device.
func lastStateByGUID(data []byte) map[string]bool {
	states := make(map[string]bool)
	for _, line := range strings.Split(string(data), "\n") {
		if ev, ok := parseNDJSONLine([]byte(line)); ok {
			states[ev.GUID] = ev.On
		}
	}
	return states
}
