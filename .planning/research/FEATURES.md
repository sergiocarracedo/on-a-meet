# macOS Camera Detection — Feature Landscape

## Existing tools that detect macOS camera state

| Tool | Approach | Language | Notes |
|------|----------|----------|-------|
| Home Assistant macOS app | AVFoundation sensors | Swift/ObjC | Camera-in-use sensor, unreliable due to appleh13camerad daemon |
| [CameraCoordinator](https://github.com/shuhaowu/CameraCoordinator) | eBPF (Linux only) | Go | Not applicable to macOS |
| Various shell scripts | `log stream` | Bash | Fragile across macOS versions |

## Key Challenges

1. **No public API** for querying camera streaming status — must use indirect methods
2. **`appleh13camerad`** camera daemon is always active on Apple Silicon Macs
3. **Green light indicator** is hardware-level with no query API
4. **Different macOS versions** use different unified log message formats
5. **USB camera interfaces** not exposed as `/dev/video*` — use IOKit registry
