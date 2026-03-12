# CardBot — MiniMax Review Notes

These are my review notes from examining the codebase and tests. Not all items need to be done — some are optional improvements, others are potential issues to investigate.

---

## Quick Wins

### Binary Size
- Run `go build -ldflags="-s -w"` to strip debug info — saves ~1MB

### Test Coverage Gaps
| Package | Coverage | Notes |
|---------|----------|-------|
| main | 0% | Not testable until main.go is split |
| detect | 11.5% | Platform code, but `detectBrand` and `FormatBytes` are testable pure functions |
| pick | 0% | macOS only (osascript) — hard to unit test |
| speedtest | 0% | Needs real filesystem |
| ui | 0% | Simple, low priority |

**Action:** Add unit tests for `detectBrand()` and `detect.FormatBytes()` — easy coverage wins.

---

## Potential Issues

### 1. No fsync after copy
`copyFile()` writes data but doesn't call `df.Sync()` before close. On Linux with removable media, data could still be in the page cache when "copy complete" is reported.

**Recommendation:** Add `df.Sync()` before the size check in `copyFile()`.

### 2. Progress callback captures mutable variable
In `copyAll`, the progress callback captures `lastUpdate` (a `time.Time` on the stack) and mutates it from inside the goroutine. Works in practice but is fragile — could cause subtle timing bugs.

**Recommendation:** Move `lastUpdate` to a package-level variable or struct field, or use `atomic.Int64` for timestamp.

### 3. FormatBytes duplication
`internal/copy/copy.go` has its own `fmtBytes()` function that duplicates `detect.FormatBytes()`. The copy package does this intentionally to avoid a dependency, but they could drift.

**Options:**
- Create a tiny `internal/format` package
- Accept the duplication (only ~10 lines)
- Have copy import detect (tighter coupling)

### 4. displayCard goroutine can't be cancelled
`displayCard()` runs analysis in a goroutine but uses polling-style guards (`isCurrentCard()`) to check for cancellation. Should take a `context.Context` for clean cancellation.

---

## Code Cleanup

### 5. main.go is too large (941 lines)
Doing too much: CLI flags, config, signal handling, event loop, card display, copy orchestration, prompts, input handling, eject, speed test.

**Recommendation:** Split into:
- `main.go` — flag parsing, config, logger setup, main() only (~100 lines)
- `app.go` — app struct, event loop
- `display.go` — printCardInfo, printPrompt, showHelp, friendlyErr
- `copy_cmd.go` — copyAll (the 120-line mini event loop)

### 6. Duplicate card info rendering
`printCardInfo` and `printInvalidCardInfo` share Status/Path/Storage/Camera display code.

**Recommendation:** Extract a `printCardHeader(card)` helper.

### 7. printPrompt inconsistency
There are calls to `printPrompt()` in various states but also hardcoded prompts in some places (like "Already copied." messages that print their own prompts).

**Recommendation:** Ensure all code paths use `printPrompt()` consistently.

### 8. startup time.Sleep calls
Three `150ms` sleeps in main() total 450ms of artificial delay. These will conflict with any future startup timing goals.

**Recommendation:** Remove or gate behind a debug flag.

### 9. friendlyErr not used everywhere
`friendlyErr` exists but some error paths still use raw `%v`. Standardize: all user-facing errors should go through `friendlyErr`.

---

## Platform / Cross-Platform

### 10. FormatBytes has build tag
`detect/shared.go` has `//go:build darwin || linux` but `FormatBytes()` is used by `main.go` for display. This will break Windows support if added later.

**Recommendation:** Move `FormatBytes()` to an unguarded file.

---

## Testing Notes

### 11. Walk errors silently swallowed
Both analyze and copy `WalkDir` callbacks return `nil` on errors. Files with permission issues or broken symlinks are silently skipped.

**Recommendation:** Consider collecting skipped-file warnings and surfacing them in results.

### 12. dotfile atomic write on FAT32
The temp-file + rename pattern may not be truly atomic on FAT32/exFAT (common card filesystems). Low risk but worth noting.

---

## Summary — Recommended Priority

| # | Item | Effort | Impact |
|---|------|--------|--------|
| 1 | Add fsync after copy | Low | Correctness |
| 2 | Split main.go | High | Maintainability |
| 3 | Test detectBrand/FormatBytes | Low | Coverage |
| 4 | Extract printCardHeader | Low | DRY |
| 5 | Standardize friendlyErr | Low | UX consistency |
| 6 | Add go build -ldflags="-s -w" | Low | Binary size |
| 7 | Move FormatBytes to unguarded file | Low | Future-proofing |
| 8 | Remove startup sleeps | Low | Cleanup |
