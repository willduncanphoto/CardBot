# CardBot — TODO

## Current Version: 0.1.5

Detection, analysis, EXIF, config, UI polish, and basic copy complete.

---

## 0.1.6 — Copy Robustness (Next Up)

- [ ] Handle card removed during copy
- [ ] Handle destination disk full
- [ ] Cancel copy in progress (with cleanup of partial files)
- [ ] File collision logic (skip if dest file exists and size matches)
- [ ] Handle "no DCIM" case — detect as volume, warn "not a camera card?"
- [ ] Handle read-only cards — warn before copy that dotfile can't be written
- [ ] Output mutex — add `outputMu sync.Mutex` to `app`; copy progress + scan
      goroutine will interleave without it
- [ ] Cancel in-flight scan on removal — `displayCard` goroutine currently finishes
      and prints results even if card was removed mid-scan; needs context/cancellation
- [ ] Better error messages for common failures (permissions, full disk, network paths)

---

## Code Cleanup

- [ ] Split `main.go` (643 lines) — extract display/prompt/UI logic into separate package
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
