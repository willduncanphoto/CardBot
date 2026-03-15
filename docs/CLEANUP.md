# CardBot Codebase Cleanup List

*Last updated after cleanup pass on 2026-03-15*

## Completed ✅

### Config Path (Critical)
- `config.Path()` now uses `os.UserConfigDir()` for platform-appropriate location

### Duplicate Code Eliminated
- `sequenceDigits` consolidated into exported `copy.SequenceDigits()`
- Duplicate `normalizeNamingMode` in copy package replaced with `isTimestampMode()` using `config.NormalizeNamingMode`
- Local naming mode constants removed from copy package — single source of truth in config

### Redundant Normalization Removed
- Removed extra `NormalizeNamingMode` call in `main.go` — `config.Load()` already normalizes, and `config.Defaults()` returns `NamingOriginal`

### Flag Description Updated
- `--setup` now says "re-run first-time setup (destination and naming)"

### Tests Improved
- All table-driven tests now use `t.Run()` subtests for failure isolation:
  - `TestSequenceDigits` (7 subtests)
  - `TestParseNamingChoice` (9 subtests, added empty + whitespace cases)
  - `TestNormalizeNamingMode` (6 subtests)
  - `TestNamingStartupLine` (2 subtests)
  - `TestNamingDisplayLine` (4 subtests, added 3-digit and 5-digit cases)
  - `TestRenamedRelativePath` (3 subtests: subdir, flat, extension case)
- New tests added:
  - `TestFormatSequence` — 8 cases including edge clamps
  - `TestSequenceRollover` — verifies max values and width for all digit counts
  - `TestIsTimestampMode` — 5 cases
  - `TestCopy_OriginalNaming_Unchanged` — confirms original mode passes through filenames
  - `TestPromptNamingModeIO_DefaultTimestamp` — confirms [2] default when timestamp configured
  - `TestPromptNamingModeIO_EOF` — confirms default returned on EOF
- Test count: 79 → 103 (all pass with `-race`)

---

## Remaining (Low Priority)

### Code Organization
- 9 files in main package — consider `internal/cli` for prompts when it grows further
- `ts()` function name is terse (but used everywhere, not worth churning)

### Test Coverage
- `main` package: 14.1% (interactive UI code, hard to unit test)
- `internal/detect`: 11.5% (hardware-dependent)

### Minor Polish
- `promptDestinationReadline()` ignores `ReadString` error (benign — returns default on EOF)
- Consider typed `NamingMode` string in config (adds safety, but churns all callsites)
- `modeDisplayName` in `app_logic.go` and `formatModeLabel` in `dotfile.go` do the same thing (title-case a mode string) — could share, but different packages
