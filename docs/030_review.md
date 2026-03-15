# 0.3.0 Phase 1 — Review Notes (Next Pass)

## What We Did

### New files (7)
- `setup_naming.go` — naming prompt, labels, display helpers
- `setup_naming_test.go` — prompt parsing, IO, label tests
- `setup_flow.go` — testable `runSetup()` wiring dest + naming prompts
- `setup_flow_test.go` — integration: setup writes naming mode to config file
- `internal/copy/naming.go` — timestamp formatting, sequence digits, rename path
- `internal/copy/naming_test.go` — unit + copy integration tests for naming
- `docs/030_phase1.md` — Phase 1 design spec

### Modified files (9)
- `main.go` — version bump, setup flow, startup log line
- `internal/config/config.go` — `Naming` struct, constants, normalization
- `internal/config/config_test.go` — naming defaults, round-trip, invalid mode
- `internal/analyze/analyze.go` — `FileDateTimes` map for EXIF capture time
- `internal/analyze/analyze_test.go` — `FileDateTimes` assertions
- `internal/copy/copy.go` — `NamingMode` option, `captureTime` field, rename during copy
- `copy_cmd.go` — passes `NamingMode` to copy options
- `display.go` — shows naming mode + detected digit count in card info
- `docs/TODO.md` — roadmap updated (0.3.0 renaming, 0.4.0 video, 0.5.0 multi-cam)

### Test results
- 79 tests across 9 packages
- All pass with `-race`
- No regressions

---

## Issues Found During Review

### 1. Re-copy skip broken for timestamp mode
**Severity: Medium**
`copyFile()` skips if `dest exists && size matches`. But when naming mode is `timestamp`, the dest path changes from `DSC_0001.NEF` to `260314T143052_001.NEF`. Re-copying the same card will copy everything again because the renamed file doesn't match the original-named file.

**Fix options:**
- Accept it for now (re-copy is rare, files are small to re-verify)
- Track renames in dotfile so re-copy can check renamed paths
- Check by size across all files in the date folder (fragile)

**Decision needed:** Is this acceptable for Phase 1?

### 2. README still hardcodes "Keep original filenames"
**Severity: Low**
Line 122 of README.md shows the old startup output. Should reflect dynamic naming mode.

**Fix:** Update README output example to show `Naming:` line.

### 3. `--setup` flag description is stale
**Severity: Low**
Flag says "re-run destination setup" but now also prompts for naming.

**Fix:** Change to "re-run setup (destination and naming)".

### 4. Naming line not shown for invalid cards
**Severity: Low**
`printInvalidCardInfo()` doesn't show the naming line. Consistent display would include it.

**Fix:** Add `Naming:` line to invalid card display too — or skip it since you can't copy invalid cards anyway.

### 5. DryRun doesn't preview renamed filenames
**Severity: Low**
`--dry-run` returns early before Phase 2 (the copy loop), so renamed filenames are never computed or displayed.

**Fix (future):** Add a dry-run rename preview that shows the mapping without copying.

### 6. Sequence ordering depends on walk order
**Severity: Low**
`filepath.WalkDir` returns files in lexicographic order within each directory. Sequence numbers follow walk order, not EXIF timestamp order. Two files shot seconds apart in different DCIM subfolders may get non-chronological sequences.

**Fix (future):** Sort `files` slice by `captureTime` before the copy loop. Simple change but needs a test.

### 7. Extension case in renamed files
**Severity: Cosmetic**
`renamedRelativePath()` uppercases the extension (`.NEF`, `.MOV`). Original files on Nikon cards already use uppercase. Other cameras may use lowercase. Current behavior is consistent but worth noting.

### 8. No test for sequence looping
**Severity: Low**
We document that sequence loops at max (999/9999/99999) but there's no test that verifies the wrap-around.

**Fix:** Add a test with > max files to confirm loop behavior.

### 9. `--setup` doesn't re-prompt naming if config already has it
**Severity: Expected**
`--setup` always prompts for naming, even if already configured. The default is pre-filled from config so user can just press Enter. This is correct behavior.

---

## Quick Fixes for Next Pass

- [ ] Update `--setup` flag description
- [ ] Update README output example
- [ ] Add sequence loop test
- [ ] Sort files by captureTime before copy loop (optional)

## Deferred to Future

- Re-copy skip with renamed files (needs dotfile v3)
- DryRun rename preview
- Naming line on invalid card display
