# Roadmap

Work items grouped by milestone.

---

## Completed

### 0.1.0 — Detection
- [x] macOS native detection (CGO + DiskArbitration)
- [x] macOS fallback detection (polling, no Xcode)
- [x] Linux detection (polling-based)
- [x] DCIM folder detection
- [x] Brand guess from folder names
- [x] Basic display (path, storage, brand)
- [x] `[e] Eject` and `[c] Cancel` actions
- [x] Card queue for multiple cards
- [x] Graceful shutdown

### 0.1.1 — Card Analysis
- [x] Walk DCIM tree recursively
- [x] Skip hidden files: `.*`, `._*`, `.DS_Store`, `.Trashes`
- [x] Group files by date
- [x] Count files per date, sum sizes
- [x] List file extensions per date
- [x] Fixed-width column formatting
- [x] Handle empty cards gracefully
- [x] Summary line with totals

### 0.1.2 — EXIF, Gear & Hardware
- [x] Add `evanoberholster/imagemeta` dependency
- [x] Extract camera model from EXIF
- [x] Display camera model
- [x] Use `DateTimeOriginal` for date grouping
- [x] Extract star ratings from XMP
- [x] Display starred count
- [x] Photo/video split in totals
- [x] Hardware info query via IOKit (macOS) / sysfs (Linux)
- [x] Raw device size vs filesystem size
- [x] CID parsing on Linux with direct SD slot
- [x] `[i]` key for hardware info (hidden command)

### 0.1.3 — Config & Destination
- [x] Config file: `~/.config/cardbot/config.json`
- [x] First-run setup — prompt for destination
- [x] Store: default destination path
- [x] CLI flags: `--dest`, `--version`, `--dry-run`, `--setup`, `--reset`
- [x] Logging: `~/.cardbot/cardbot.log` (5MB rotation)

### 0.1.4 — UI Polish
- [x] Merge brand + camera into single `Camera:` line
- [x] Clean brand names ("NIKON Z 9" → "Nikon Z 9")
- [x] Brand colors (ANSI): Nikon yellow, Canon red, Sony white, etc.
- [x] Card status line (New / Copied via `.cardbot` dotfile)
- [x] Parallel EXIF workers (4 default, 3.7x faster)
- [x] Progress callback throttled to every 100 files
- [x] Config save skipped when destination unchanged
- [x] Native macOS folder picker via `osascript`
- [x] `~` shorthand for paths (display and storage)
- [x] `[t]` hidden speed test command

### 0.1.5 — Copy (Current)
- [x] `[a] Copy All` mode
- [x] Dated folders: `YYYY-MM-DD/<original-structure>`
- [x] Buffered copy with configurable buffer size
- [x] Progress updates every 2 seconds
- [x] Size verification after each file
- [x] `.cardbot` dotfile written on success
- [x] "Copied on" status on re-insert
- [x] Destination write probe on first creation
- [x] Dry-run aware

---

## Upcoming

### 0.1.6 — Copy Robustness
- [ ] Handle card removed during copy
- [ ] Handle destination disk full
- [ ] Cancel copy in progress (with cleanup)
- [ ] File collision logic (skip existing)
- [ ] Handle "no DCIM" case with warning
- [ ] Handle read-only cards (warn before copy)
- [ ] Output mutex for concurrent progress/scan output
- [ ] Cancel in-flight scan on card removal
- [ ] Better error messages

### 0.1.7 — Linux
- [ ] Test on Ubuntu, Fedora, Debian
- [ ] Document mount point behavior
- [ ] Linux build marked stable

### 0.1.8 — Polish
- [ ] Startup <100ms
- [ ] Single-key input (raw terminal mode)
- [ ] Estimated time remaining during copy
- [ ] Performance benchmarks

### 0.1.9 — Distribution
- [ ] README with install/usage
- [ ] GitHub releases (macOS Intel/ARM, Linux AMD64/ARM64)
- [ ] Homebrew formula
- [ ] `--help` with examples
- [ ] LICENSE file

---

## Later

- Windows support
- File renaming on copy (date-based, camera+date, sequence numbering)
- `[s] Selects` copy mode (starred only)
- `[v] Videos` / `[p] Photos` copy modes
- Incremental copy (only new/changed files)
- Resume interrupted copies
- Video metadata (duration, resolution)
- Network destinations
- TOML/YAML config
- JSON output mode
- Toggle flat vs preserve DCIM structure
- Auto-update: check GitHub Releases for new version at startup, `--update` flag to self-upgrade
  - Prereqs: GoReleaser for multi-platform builds, checksums, `-ldflags` version injection
  - Consider `go-selfupdate` library or DIY (~200 lines for single binary)
  - UX: non-blocking check once/day, notify user, don't force
