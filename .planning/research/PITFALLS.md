# Pitfalls Research

**Domain:** Linux camera state detection CLI tool
**Researched:** 2026-05-28
**Confidence:** HIGH

## Common Mistakes

| # | Mistake | Severity | Frequency | Impact |
|---|---------|----------|-----------|--------|
| 1 | Device busy on every poll causes false positives | CRITICAL | Common | V4L2 device returns "busy" both when camera is in use AND when another process briefly queried it. Polling can cause false state changes. |
| 2 | Not handling device hotplug/unplug | HIGH | Common | USB cameras can be disconnected. Tool crashes or hangs on missing /dev/videoN. |
| 3 | Race between polling and kernel driver | MEDIUM | Occasional | StreamON/StreamOFF ioctl may not reflect actual user-space usage. A process can hold the device open without streaming. |
| 4 | Permission errors on /dev/video* | HIGH | Very Common | /dev/video* typically owned by video group. Tool fails silently or with cryptic error. |
| 5 | Command execution blocks the polling loop | MEDIUM | Common | If spawned command takes longer than poll interval, overlapping executions cause chaos. |

### Mistake Details

**1. Device busy false positives**
- **What happens:** V4L2 returns EBUSY when device is already opened by another process. But brief queries (v4l2-ctl, Cheese thumbnail) also open the device momentarily, causing false "camera on" events.
- **Why it happens:** Most tools check "can I open /dev/videoN?" rather than "is someone actively streaming?". Opening the device == false positive.
- **Example:** User runs `v4l2-ctl --all` to debug → triggers on-a-meet → fires --on command → HA toggles light on. User confused.
- **Fix cost:** MEDIUM — need smarter detection (check STREAMON status, not just open-file)

**2. Device hotplug/unplug**
- **What happens:** Camera disconnected mid-polling → /dev/videoN disappears → tool crashes or enters error loop.
- **Why it happens:** Not handling ENOENT gracefully in the poll loop.
- **Fix cost:** LOW — graceful degradation (log warning, continue polling other cameras)

**3. StreamON/StreamOFF vs open-file ambiguity**
- **What happens:** Some apps (Chrome, Zoom) open the V4L2 device at startup but only start streaming when the user joins a call. Result: "camera on" triggers when browser opens, not when call starts.
- **Why it happens:** V4L2's VIDIOC_STREAMON is not always called. Some drivers start streaming on first QBUF. Some apps map buffers without streaming.
- **Fix cost:** HIGH — truly reliable detection requires eBPF tracing of kernel vb2_ioctl_streamon/off

**4. Permission errors**
- **What happens:** Tool runs as non-root user without video group membership → cannot open /dev/videoN → reports "no cameras" or crashes.
- **Why it happens:** /dev/video* permissions default to root:video 0660. User must be in video group.
- **Fix cost:** LOW — check permissions at startup, print clear error with fix instructions

**5. Blocking command execution**
- **What happens:** User sets --on "sleep 60". Tool spawns it, blocks for 60s, misses state changes during that time.
- **Why it happens:** Sequential polling+execution without async spawn.
- **Fix cost:** LOW — spawn commands in goroutines with timeout/context

## Warning Signs

| Warning Sign | Indicates | Action |
|-------------|-----------|--------|
| Tool reports camera state flipping on/off rapidly | V4L2 false positives (Mistake #1) | Switch to lsof backend or add debounce |
| Tool crashes with "no such file or directory" mid-run | Camera disconnected (Mistake #2) | Add hotplug handling, graceful degrade |
| Tool says "camera on" but camera LED is off | StreamON ambiguity (Mistake #3) | Document limitation, suggest lsof backend |
| Tool says "no cameras found" | Permission issue (Mistake #4) | Check video group membership, print fix |
| Commands fire multiple times for one state change | Race condition or no debounce (#1, #5) | Add state debounce window |

## Prevention Strategies

| Strategy | Prevents | When to Apply | How |
|----------|----------|---------------|-----|
| Debounce window | #1 False positives | Detection engine design | Require N consecutive same-state polls before firing. Configurable debounce duration. |
| Graceful device removal | #2 Hotplug | Poll loop | Handle ENOENT per-device, not globally. Re-scan device list periodically. |
| Multiple detection backends | #1, #3 Ambiguity | Architecture | If V4L2 is unreliable, fall back to lsof (most reliable). Make backend user-selectable. |
| Goroutine command execution | #5 Blocking | Command execution | Always exec in goroutine. Add configurable timeout. Track running commands. |
| Permission check at startup | #4 Permission | CLI init | Check /dev/video* accessibility on start. Print actionable error. |
| Document video group requirement | #4 Permission | README/help | Clear docs: "Must be in video group. Run: sudo usermod -a -G video $USER" |
| Configurable debounce + interval | #1, #5 False positives | Config design | Users adjust for their hardware. Default: 1s interval, 2s debounce. |

## Domain-Specific Patterns

### Patterns That Look Right But Aren't

| Pattern | Why It Seems Good | Actual Problem | Better Approach |
|---------|-------------------|----------------|-----------------|
| Check if /dev/videoN is open (fuser/lsof) | Direct, simple | A process can open the device without streaming. Chrome opens webcam on page load, not on call join. | Combine lsof with V4L2 VIDIOC_QUERYCAP for more reliable check |
| Single polling interval for all cameras | Simple implementation | Different cameras may need different rates | Per-camera interval (v2) or single sensible default |
| Re-execute the tool for each poll | "Simple shell script" | Slow, no state persistence between runs | Daemon mode with internal state machine |

### Patterns That Look Wrong But Work

| Pattern | Why It Seems Bad | Why It Actually Works | When to Use |
|---------|------------------|----------------------|-------------|
| Polling (busy-wait) | Inefficient, not modern | Camera state changes happen at human scale (seconds). 1s polling = negligible CPU. No kernel deps. | Default detection mode |
| Using lsof in a Go tool | Parsing external command output is fragile | lsof has been stable for decades. Reliable, works across all Linux distros. | Fallback detection when V4L2 gives false positives |
| Config file + overrides | Two sources of truth is confusing | Viper's precedence is well-defined: CLI > env > config > default. Predictable. | Standard for production CLI tools |

---
*Pitfalls research for: on-a-meet camera detection CLI*
*Researched: 2026-05-28*
