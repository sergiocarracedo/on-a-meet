# macOS Camera Detection — Tech Stack

## Recommended Approach: `log stream` + `lsof`

No existing pure-Go library reliably detects macOS camera on/off state. The ecosystem is fragmented across OS versions.

## Primary Options

| Approach | Type | cgo? | Reliable? | Fragile? | Notes |
|----------|------|------|-----------|----------|-------|
| `log stream --predicate ...` | exec `log` | No | Medium | High | Different predicates per macOS version. No cgo. |
| `lsof` on USB camera IO services | exec `lsof` | No | Medium | Medium | Lists processes with open USB camera interfaces |
| AVFoundation via cgo | cgo + ObjC | Yes | High | Low | Full Apple API access but requires cgo + Xcode |
| IOKit via darwinkit | darwinkit Go bindings | Yes | High | Medium | Requires cgo, 3rd-party binding layer |
| `appleh13camerad` process check | exec `ps` | No | Low | Low | Daemon always runs, not reliable |

## Recommendation for v1.1.0

Use `log stream` as primary detection on macOS (similar to Linux's `lsof` backend — lightweight, no cgo). The `log` command is built into macOS, so no dependencies. Monitor the unified log for camera stream events.
