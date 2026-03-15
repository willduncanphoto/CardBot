# CardBot 0.3.0 — Phase 1: Basic Renaming (Proof of Concept)

## Goal
Add a simple binary choice for file naming that appears during first-time setup.
Two options only:
1. **Keep original filenames** (current behavior, default)
2. **Timestamp + Sequence** → `YYMMDDTHHMMSS_001.NEF` (3-5 digit auto)

## Note on File Type Examples

All examples in this spec use:
- **NEF** — Nikon RAW photo files
- **MOV** — Video files (generic/professional format)

NEV (Nikon N-RAW video) and JPEG support exist in CardBot but are not emphasized — this tool is built for RAW photo workflows.

## Engineering Note: Subsecond Support

**Concern:** Cameras shooting high-speed burst (Z9 does 20fps RAW) can have multiple frames within the same second. Without subsecond precision, collisions conceptually occur: `260314T143052_001.NEF`, `260314T143052_002.NEF` — both RAW frames from the same second.

**Options considered:**
- **Subsecond digits** (`260314T143052_42_001.NEF`) — EXIF doesn't reliably expose subsecond data across camera brands
- **Frame counter from EXIF** — not consistently available
- **Sequence per second reset** — complex, error-prone
- **Original filename hash suffix** — ugly, hard to read

**Decision for Phase 1:** Ignore subsecond. Use sequence number as the disambiguator. A 20fps burst of 50 frames in one second gets `_001` through `_050` (on a typical 3-digit card). Good enough for photography workflows, not science timing. Revisit if users report issues with burst sequences.

**Known limitation (Phase 1):** Multi-camera collision risk. If two cameras shoot at `260314T143052` and both generate `_001`, files collide in the same dated folder. Mitigation: use original filenames (unique per camera) or copy cameras to different destinations. Full solution (camera ID prefix) deferred to 0.5.0.

## User Flow

### First Run Experience

```
# Existing folder picker appears first
Select destination folder: ~/Pictures/CardBot

# NEW: Naming preference prompt
────────────────────────────────────────
File Naming
────────────────────────────────────────
How would you like files named when copying?

[1] Keep original filenames (DSC_0001.NEF)
[2] Timestamp + sequence (260314T143052_001.NEF)

The timestamp comes from when each photo was taken.
Sequence digits adjust automatically: 3 digits for typical cards,
4 or 5 digits for cards with 1000+ or 10000+ files.
You can change this later in settings.

> 2

Naming set to: Timestamp + sequence (auto-detected digits)
────────────────────────────────────────
[2026-03-14T10:23:15] Starting CardBot 0.3.0...
```

### Subsequent Runs

The choice is saved to config. On future runs, CardBot starts normally but displays the current naming mode:

```
  Status:   New
  Path:     /Volumes/NIKON Z 9
  Storage:  96.4 GB / 476.9 GB (20%)
  Camera:   Nikon Z 9
  Starred:  1
  Content:  2026-02-27      12.9 GB    418   NEF

  Total:    3048 photos, 0 videos, 96.0 GB
  Naming:   Timestamp sequence (260314T143052_0001.NEF) [4-digit]
────────────────────────────────────────
[a] Copy All  [s] Copy Selects  [p] Copy Photos  [v] Copy Videos  [e] Eject  [x] Exit
```

Or if keeping original:
```
  Naming:   Keep original (DSC_0001.NEF)
```

## Configuration

```json
{
  "$schema": "cardbot-config-v2",
  "destination": {
    "path": "~/Pictures/CardBot"
  },
  "naming": {
    "mode": "timestamp"
  }
}
```

**Modes:**
- `"original"` — keep source filenames (default)
- `"timestamp"` — YYMMDDTHHMMSS + _sequence (auto-detected digits)

**Sequence digits:** auto-detected from card scan (3, 4, or 5 digits)
- ≤ 999 files → 3 digits (001-999)
- ≤ 9,999 files → 4 digits (0001-9999)  
- ≤ 99,999 files → 5 digits (00001-99999)
- > 99,999 files → 5 digits (loops at 99999, extremely rare)

## Naming Format: `YYMMDDTHHMMSS_NNN.EXT`

| Component | Source | Example |
|-----------|--------|---------|
| YY | EXIF year (2 digits) | 26 |
| MM | EXIF month | 03 |
| DD | EXIF day | 14 |
| T | ISO 8601 separator | T |
| HH | EXIF hour | 14 |
| MM | EXIF minute | 30 |
| SS | EXIF second | 52 |
| _NNN | Sequence per card | _001, _0001, _00001 (auto-detected) |
| EXT | Original extension | .NEF |

**Sequence behavior:**
- Starts at 1 for each new card (padded based on card size)
- Increments per file copied (regardless of type)
- Loops back to 1 after max (999, 9999, or 99999 based on detected digits)
- Separate sequence per copy session

**Why auto-detect:** Most cards need only 3 digits (under 1000 shots). Large event cards may need 4. Multi-day shoots may need 5. CardBot inspects the total file count and picks the minimum digits needed for uniqueness.

**Fallback:** If EXIF date is missing, use file modification time.

## Files to Modify

### 1. `internal/config/config.go`
- Add `NamingConfig` struct with `Mode` and `Padding`
- Default: `Mode: "original"`, `Padding: 3`
- Config schema version stays v2 (additive change)

### 2. `internal/pick/pick_darwin.go` OR new `internal/setup/setup.go`
- Create `RunFirstTimeSetup()` that:
  1. Calls existing folder picker
  2. NEW: Prompts for naming preference
  3. Saves complete config

### 3. `main.go`
- After folder picker returns, check if naming is configured
- If not set, prompt for naming preference before starting

### 4. `display.go`
- Add `printNamingInfo()` to show current naming mode in card display
- Shows either "Keep original" or "Timestamp sequence" with example

### 5. `internal/copy/copy.go`
- Add `RenameFunc` to `Options` struct
- If set, apply rename during copy
- Sequence counter managed per copy session

### 6. `copy_cmd.go`
- Create `buildRenameFunc(card, mode)` that returns rename function
- Function takes: original path, EXIF data → returns new filename
- Handles sequence increment atomically

### 7. `app.go`
- Pass naming config to copy command
- Wire up rename function when mode != "original"

## UX Display Examples

### With timestamp naming enabled (3-digit card):
```
[2026-03-14T10:23:15] Copying (All)...
  156/3048 files  4.2 GB  85 MB/s
  260314T143052_156.NEF
```

### File in destination (typical card, 3-digit):
```
~/Pictures/CardBot/
└── 2026-03-14/
    └── 100NIKON/
        ├── 260314T143052_001.NEF
        ├── 260314T143052_002.NEF
        └── 260314T143101_003.MOV
```

### Large card example (4-digit, 3000+ files):
```
~/Pictures/CardBot/
└── 2026-03-14/
    └── 100NIKON/
        ├── 260314T143052_0001.NEF
        ├── 260314T143052_1500.NEF
        └── 260314T163022_3001.NEF
```

## Implementation Checklist

- [ ] Add `NamingConfig` to config package
- [ ] Create setup flow for naming preference (post-folder-picker)
- [ ] Display current naming mode in card info (show detected digit count)
- [ ] Build rename function with sequence counter
- [ ] Auto-detect sequence digits from card file count (3/4/5)
- [ ] Integrate rename into copy engine
- [ ] Handle missing EXIF date (fallback to modtime)
- [ ] Sequence resets per card
- [ ] Test with real files

## Out of Scope for Phase 1

- Changing naming mode without `--reset` (use `cardbot --reset` to re-run setup)
- Custom templates
- Preview mode
- Per-file naming toggle
- Global sequence across cards
- Dotfile tracking of renames (keep it simple)

## Future Phases Will Add

- `[r]` key to change naming interactively
- `[n]` toggle during session
- Custom templates
- Preview before copy
- Dotfile v3 with rename mapping
