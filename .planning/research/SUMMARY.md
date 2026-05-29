# macOS Camera Detection — Research Summary

## Conclusion

macOS camera detection is feasible via the `log stream` command, following the same pattern as the existing Linux `lsof` backend (exec-based, no cgo). The `log` command is built into macOS and can detect camera stream ON/OFF events from the unified log.

## Trade-offs

| Approach | Pro | Con |
|----------|-----|-----|
| `log stream` | No cgo, no deps, exec-based like lsof | Fragile across macOS versions |
| AVFoundation cgo | Most reliable | Requires cgo + Xcode, green light always on |
| darwinkit-bindings | Wraps AVFoundation in Go | cgo, 3rd-party dep, large import surface |

## Recommendation

Go with `log stream` for v1.1.0. It matches the project's existing pattern of exec-based backends, keeps the no-cgo constraint, and is serviceable for a first macOS implementation. Document the version-specific fragility and add fallback paths.
