# CardBot — TODO

## Current Version: 0.1.6

Detection, analysis, EXIF, config, UI polish, copy with robustness, and UX improvements complete.

---

## 0.1.7 — Polish (Next Up)

- [ ] Single-key input (raw terminal mode, no Enter required)
- [ ] Startup under 100ms
- [ ] Estimated time remaining during copy
- [ ] Copy Selects mode (`[s]` — starred files only)
- [ ] Show current filename during copy (deferred to renaming milestone)

---

## Code Cleanup

- [ ] Split `main.go` (~900 lines) — extract display/prompt/UI logic into separate package
- [ ] Drop `app.printf()` method — use explicit `fmt.Printf` + `a.logf` instead
- [ ] Review `OUTPUT.md` aspirational features vs reality — trim or mark as future

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
