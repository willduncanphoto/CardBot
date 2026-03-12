# CardBot — TODO

## Current Version: 0.1.6

Detection, analysis, EXIF, config, UI polish, copy with robustness, and UX improvements complete.

**Target: 0.2.0 — Daily Driver.** The version you hand to another photographer and say "try this."

---

## 0.1.7 — Polish (Next Up)

- [ ] Estimated time remaining during copy
- [ ] Show current filename during copy (deferred to renaming milestone)

---

## Wishlist

These are "nice to have" features that aren't on the immediate roadmap:

- Single-key input (raw terminal mode, no Enter required)
- Auto-update: check GitHub Releases for new version at startup
- Network destination support
- Windows support
- JSON output mode for scripting
- Star rating filters: `[2]` Copy 2★, `[3]` Copy 3★, `[4]` Copy 4★, `[5]` Copy 5★

---

## 0.1.8 — Selective Copy

- [ ] `[s]` Copy Selects — copy starred/picked files only (XMP rating > 0)
- [ ] `[p]` Copy Photos — copy photo files only (RAW + JPEG, no video)
- [ ] `[v]` Copy Videos — copy video files only (MOV, MP4, MXF, etc.)
- [ ] Dotfile tracks copy mode per operation (`"mode": "selects"`)
- [ ] Status line reflects partial copy (`Selects copied on ...`)
- [ ] Re-copy guard per mode — don't skip if previous copy was a different mode
- [ ] Disk space preflight scoped to selected file subset
- [ ] Help removes strikethrough from `[s]`, `[p]`, `[v]` once implemented

### Partial Copy State — Dotfile Design

The `.cardbot` dotfile currently tracks a single copy event. With selective copy modes,
a card may be partially copied in multiple independent passes (e.g. videos first in the
field, photos later in the studio). The dotfile needs to track each mode independently.

Questions to answer before implementation:
- [ ] Store a `copies` array in the dotfile — one entry per mode with timestamp, dest, file count, bytes
- [ ] Status line logic: what to show when multiple modes have been copied? e.g.
      `Photos + Videos copied` vs `All copied` vs `Selects copied on ...`
- [ ] Should `[a] Copy All` mark all selective modes as complete, or only the "all" mode?
- [ ] If photos were copied and user runs `[a]`, should photo files be skipped (size check) or re-evaluated?
- [ ] Consider a `completed_modes []string` field so the UI can show checkmarks per mode

---

## Code Cleanup

- [ ] Split `main.go` (~941 lines) — extract display/prompt, copy orchestration, app logic
- [ ] Drop `app.printf()` — use explicit `fmt.Printf` + `a.logf` pairs
- [ ] Add `context.Context` to `displayCard` and analyzer for clean cancellation
- [ ] Move `FormatBytes` to unguarded file (currently darwin/linux only via build tag)
- [ ] Extract `printCardHeader` helper — shared between `printCardInfo` and `printInvalidCardInfo`
- [ ] Remove startup `time.Sleep` calls (3 × 150ms) — conflicts with 0.1.7 startup goal
- [ ] Add `df.Sync()` before close in `copyFile` — correctness on Linux removable media
- [ ] Standardize error handling — `friendlyErr` everywhere user-facing, raw `%v` log only

---

## Future Features

- [ ] **File renaming** — configurable rename patterns on copy (e.g. date-based,
      camera+date, sequence numbering). Current behavior: keep original filenames.

---

## Speed Test — Future Improvements

Current implementation is a synthetic 256MB sequential read/write benchmark.

- [ ] Write test files sized like actual RAW photos (e.g. 50MB for Z9 NEF)
- [ ] Write test files sized like actual video clips (e.g. 500MB–2GB for N-RAW/ProRes)
- [ ] Multi-file burst test (simulate ingesting a full card worth of files)
- [ ] Report burst speed vs sustained speed separately
- [ ] Compare results against card's rated spec
- [ ] Warn if measured speed is significantly below rated speed
- [ ] Bypass OS page cache for accurate read speeds (`F_NOCACHE` / `fcntl`)

---

## Testing Notes

- [ ] **Destination path display** — verify `~` shorthand across all setup flows:
      first run, `--setup`, folder picker, manual text entry
- [ ] **Destination write probe** — test copy to: new directory, existing directory,
      read-only path, full disk, network volume
- [ ] **Copy verification** — test with real card (Z9 NEF files), verify EXIF dates
      drive folder grouping instead of mtime
- [ ] **Dotfile round-trip** — copy card, re-insert, verify "Copied on" status displays
- [ ] **Re-copy behavior** — copy same card twice, verify no file collisions or errors
- [ ] **Selective copy** — copy selects, re-insert, verify status and re-copy guard
