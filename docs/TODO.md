# CardBot — TODO

## Current Version: 0.1.7

Detection, analysis, EXIF, config, UI polish, copy with robustness, UX improvements,
and bug fixes complete. 100 tests across 8 packages, all passing with `-race`.

**Target: 0.2.0 — Daily Driver.** The version you hand to another photographer and say "try this."

---

## Quick Fixes (land in any release)

- [x] Add "OM System" to `BrandColor` in `ui/color.go` — `cleanGear` already maps
      `"OM DIGITAL"` → `"OM System"` but the color map only has `"Olympus"` → cyan
- [x] `go build -ldflags="-s -w"` — strips debug info, saves ~1.8MB on binary (Makefile added)

---

## 0.1.8 — Selective Copy

Core feature: let users copy subsets of a card instead of everything.

- [ ] `[s]` Copy Selects — copy starred/picked files only (XMP rating > 0)
- [ ] `[p]` Copy Photos — copy photo files only (RAW + JPEG, no video)
- [ ] `[v]` Copy Videos — copy video files only (MOV, MP4, MXF, etc.)
- [ ] Dotfile tracks copy mode per operation (`"mode": "selects"`)
- [ ] Status line reflects partial copy (`Selects copied on ...`)
- [ ] Re-copy guard per mode — don't skip if previous copy was a different mode
- [ ] Disk space preflight scoped to selected file subset
- [ ] Help removes strikethrough from `[s]`, `[p]`, `[v]` once implemented

### Dotfile Design (decide before implementation)

The `.cardbot` dotfile currently tracks a single copy event. With selective copy modes,
a card may be partially copied in multiple independent passes (e.g. videos first in the
field, photos later in the studio). The dotfile needs to track each mode independently.

- [ ] Store a `copies` array — one entry per mode with timestamp, dest, file count, bytes
- [ ] Status line logic: what to show when multiple modes have been copied?
      e.g. `Photos + Videos copied` vs `All copied` vs `Selects copied on ...`
- [ ] Should `[a] Copy All` mark all selective modes as complete, or only the "all" mode?
- [ ] If photos were copied and user runs `[a]`, should photo files be skipped (size check)
      or re-evaluated?
- [ ] Consider a `completed_modes []string` field so the UI can show checkmarks per mode

---

## 0.1.9 — Code Health

Cleanup pass with only verified, real issues from four model reviews cross-checked
against the actual codebase.

### Split main.go (~995 lines)

`main.go` does too much: CLI flags, config, signal handling, event loop, card display,
copy orchestration, prompts, help, input, eject, speed test. Split into:

- **`main.go`** — flag parsing, config, logger setup, signal handling, `main()` (~100 lines)
- **`app.go`** — `app` struct, event loop, card/queue management, `handleInput`
- **`display.go`** — `printCardInfo`, `printInvalidCardInfo`, `printPrompt`, `showHelp`,
  `showHardwareInfo`, `friendlyErr`
- **`copy_cmd.go`** — `copyAll` method (the 120-line copy event loop)

### Extract `printCardHeader` helper

`printCardInfo` and `printInvalidCardInfo` both render Status, Path, Storage, Camera
with slightly different formatting. Extract shared `printCardHeader(card, result)`.

### Add `context.Context` to `displayCard` and analyzer

`displayCard` runs in a goroutine but can't be cancelled if the card is removed mid-scan.
The `isCurrentCard` check after `Analyze()` catches removal but analysis keeps running
until completion. Threading `context.Context` through the analyzer's `WalkDir` enables
clean cancellation.

### Log walk errors instead of swallowing them

Both `analyze.go` and `copy.go` `WalkDir` callbacks return `nil` on all errors —
permission denied, I/O errors, and broken symlinks are silently skipped. Permission
denied on a dying card would be completely invisible.

Fix: log warnings for permission/IO errors via the progress callback or collected
warnings slice. Broken symlinks and hidden files can stay silent.

### Standardize `friendlyErr` for all user-facing errors

Most error paths use `friendlyErr` but a few don't:
- Dotfile write warning (line ~714) shows raw `%v`
- Config load errors show raw `%v`

Route all user-facing errors through `friendlyErr`.

### Validate destination path

Config accepts any string for `destination.path`. Empty string, `/dev/null`, or an
unwritable path fails with a confusing error at copy time. Validate on config load
or at minimum on copy start with a clear message.

### Move `FormatBytes` to platform-agnostic file

`detect/shared.go` has `//go:build darwin || linux` but `FormatBytes` is a pure
function with no platform dependencies. Move it to an unguarded file so it compiles
on all platforms.

### Remove 500ms `displayCard` delay

`handleCardEvent` sleeps 500ms before calling `displayCard` (line 351). No comment
explains why. The card is already detected — analysis can start immediately.

### Test coverage improvements

| Package | Coverage | Action |
|---------|----------|--------|
| main | 0% | Blocked on split — becomes testable after refactor |
| detect | 11.5% | Add unit tests for `detectBrand` and `FormatBytes` (pure functions) |
| pick | 0% | macOS-only osascript — skip |
| speedtest | 0% | Needs real filesystem — skip or integration test |
| ui | covered | Already has tests |

Target: 80%+ across testable packages after split.

---

## 0.2.0 — Daily Driver

Everything from 0.1.x is solid, tested, and feels intentional.

- [ ] All 0.1.7, 0.1.8, 0.1.9 items complete
- [ ] Single-key input (raw terminal mode, no Enter required)
- [ ] Selective copy fully implemented with correct status tracking
- [ ] Partial copy state in dotfile — multi-mode copy history
- [ ] No known crashes or data loss scenarios
- [ ] Tested on personal gear across multiple shooting days
- [ ] Feedback from at least one other photographer
- [ ] README reflects actual current behavior
- [ ] First public-facing release candidate

---

## 0.3.0 — Linux Support

- [ ] Linux detection (polling-based, /run/media, /media, /mnt)
- [ ] Linux hardware info (sysfs, CID parsing)
- [ ] Linux speed test
- [ ] Linux eject (udisksctl / umount)
- [ ] Real-world testing (Ubuntu, Fedora, Debian)
- [ ] Stable build and CI

---

## Wishlist

Not on the immediate roadmap. Nice-to-have for someday.

- Estimated time remaining during copy
- Show current filename during copy (deferred to renaming milestone)
- Per-file copy logging (forensic/recovery audit trail)
- Single-key input (raw terminal mode, no Enter required) → promoted to 0.2.0
- Auto-update: check GitHub Releases for new version at startup
- Network destination support
- Windows support
- JSON output mode for scripting
- Star rating filters: `[2]` Copy 2★+, `[3]` Copy 3★+, `[4]` Copy 4★+, `[5]` Copy 5★ only
- File renaming on copy (date-based, camera+date, sequence numbering)
- Resume interrupted copies
- Video metadata (duration, resolution)

---

## Won't Fix

Items raised in code reviews that were investigated and rejected.

| Item | Why |
|------|-----|
| `lastUpdate` race in copy progress | Not a race — only the copy goroutine reads/writes it |
| `cardInvalid` naming ("negative name") | Reads fine: `if a.cardInvalid`. `hasDCIM` would be worse |
| Queue can grow unbounded | Photographers don't have 10 card readers. Never happens |
| Input channel size 10 vs 1 | Works fine with `drainInput()`. Not worth changing |
| FAT32 dotfile atomicity | Rename is atomic for metadata. Non-issue |
| XMP buffer too small/large | 256KB is correct for camera RAW headers |
| God object / extract pure functions | Premature abstraction. `app` struct is manageable |
| Version constant should be typed | Idiomatic Go. `const version = "0.1.7"` is correct |
| Log file needs fsync | CLI tool log doesn't need fsync on every write |
| `printf` vs `fmt.Printf` inconsistent | Actually consistent: `a.printf` = print+log, `fmt.Printf` = transient output |
| `FormatBytes` duplication in copy | Already fixed in 0.1.7 — copy imports `detect.FormatBytes` |
| Detector channels unbuffered | Wrong — they're buffered at size 10 |

---

## Review History

This file consolidates findings from four independent code reviews (Claude, Kimi,
MiniMax, GLM) conducted on 2026-03-12, cross-checked against the actual codebase
on 2026-03-13. Stale items (already fixed in 0.1.7) were removed. Disagreements
were resolved by reading the code. The individual review files have been retired.
