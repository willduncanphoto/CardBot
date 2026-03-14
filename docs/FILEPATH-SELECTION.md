# Filepath Selection Design Notes

**Status:** Current implementation (as of 0.1.9)  
**Target:** 0.3.0 improvements wishlist

---

## Current Behavior

### First-Run Flow

When CardBot launches without an existing config file (or with `--setup` flag):

1. **Welcome message** — "Welcome to CardBot! Where should CardBot copy your work?"
2. **Native macOS folder picker** — Opens Finder-style folder selector
3. **Fallback to readline** — On Linux or if AppleScript fails, falls back to stdin prompt
4. **Save to config** — Selected path saved to `~/.config/cardbot/config.json`

### Implementation Details

#### Native Picker (`internal/pick/`)

**macOS (`pick_darwin.go`):**
```go
// Uses AppleScript osascript for native Finder dialog
script := `POSIX path of (choose folder with prompt "Where should CardBot copy your work?" default location POSIX file "...")`
exec.Command("osascript", "-e", script)
```

- Prompt: "Where should CardBot copy your work?"
- Default location: `~/Pictures/CardBot` (the default path)
- Returns POSIX path (e.g., `/Users/name/Pictures/CardBot`)
- Path escaping handles backslashes and quotes

**Other platforms (`pick_other.go`):**
```go
return "", errors.New("native folder picker not available on this platform")
```

#### Fallback Path (`main.go`)

```go
func promptDestinationReadline(defaultPath string) string {
    fmt.Printf("Destination [%s]: ", defaultPath)
    // bufio.NewReader + strings.TrimSpace
    // Empty input accepts the default
}
```

#### Default Path

```go
// From config.Defaults()
Destination: Destination{
    Path: "~/Pictures/CardBot",
},
```

- `~/Pictures/CardBot` is the built-in default
- Path expansion happens via `config.ExpandPath()` (handles `~` to home dir)
- Path contraction happens via `config.ContractPath()` (stores `~` shorthand in config)

---

## CLI Overrides

| Flag | Behavior |
|------|----------|
| `--dest <path>` | Skip setup, use this path directly |
| `--setup` | Force re-run destination setup even if config exists |
| `--reset` | Wipe config, then run setup |

---

## Edge Cases Handled

1. **Config file doesn't exist** → triggers setup
2. **AppleScript fails** → falls back to readline
3. **User cancels dialog** → empty string, falls back to readline
4. **Path with spaces** → properly escaped in AppleScript
5. **Path with quotes** → properly escaped in AppleScript
6. **Empty input on fallback** → accepts default path
7. **Custom path input** → validates and uses that

---

## 0.3.0 Improvement Ideas

### Volume Detection & Selection

Current limitation: User must navigate to folder, can't see available volumes at a glance.

**Potential improvements:**

1. **Volume list in picker** — Show mounted volumes first (SD cards, external drives)
2. **Smart default based on volumes** — If only one external volume mounted, default to it
3. **SD card auto-detection** — Pre-select if exactly one SD card detected
4. **Recent destinations** — Remember last 3-5 destinations, allow quick re-selection
5. **Bookmark favorites** — Pin frequently used destinations

### Better Fallback Experience

Current Linux fallback is basic stdin. Potential improvements:

1. **zenity/kdialog support** — Add native pickers for Linux
2. **Tab completion** — Readline with path tab completion
3. **Path validation** — Check permissions/disk space before accepting
4. **Path preview** — Show resolved full path before confirming

### Configuration Enhancements

1. **Multiple destination profiles** — "Studio", "Field", "Backup" modes
2. **Per-project destinations** — Remember destination per card/camera
3. **Default subfolder naming** — Configurable patterns: `YYYY-MM-DD/`, `PROJECT-NAME/`, etc.
4. **Destination validation** — Warn if destination is nearly full, read-only, etc.

### UX Flow Ideas

**Idea: Volume-first selection**
```
Select destination:
  [1] Macintosh HD (/Users/name/Pictures/CardBot)
  [2] SANDISK-128G (/Volumes/SANDISK-128G) ← external
  [3] LACIE-4TB (/Volumes/LACIE-4TB) ← external
  [4] Browse other location...
  [5] Type path manually...
```

**Idea: Post-detection destination change**
- Allow changing destination mid-session without restart
- `[d] Destination` key to re-open picker
- Warn if changing destination with pending copies

---

## Open Questions for 0.3.0

1. Should we keep the native Finder dialog or build a custom TUI volume selector?
2. How should we handle multiple cards with different destinations? (e.g., main + backup)
3. Should destination be per-card (remembered) or always global?
4. How to handle network destinations? (SMB, NFS, etc.)
5. Should we validate destination writability at selection time or copy time?
6. What about cloud destinations? (Dropbox, Google Drive, etc.)

---

## Related Files

| File | Purpose |
|------|---------|
| `internal/pick/pick_darwin.go` | macOS native folder picker |
| `internal/pick/pick_other.go` | Fallback for non-macOS |
| `internal/config/config.go` | Path expansion/contraction, defaults |
| `main.go` | `promptDestination()`, `promptDestinationReadline()` |
| `docs/CONFIG.md` | Config file documentation |
