# CardBot — Kimi Review Notes

Fresh code review with different observations from the previous review.

---

## Architecture Observations

### 1. The `app` struct is a god object
It holds 13 fields including detector, config, logger, multiple channels, and state flags. This makes the code hard to test because you need to construct a full app to test any method.

**Suggestion:** Extract pure functions that don't need the full app struct. For example, `printCardInfo` could take `(card, result, cfg)` instead of being a method on `*app`.

### 2. Channel buffering is inconsistent
- `inputChan` is buffered (size 10)
- `sigChan` is buffered (size 1)
- `doneCh` in copy is buffered (size 1)
- `detector.Events()` and `detector.Removals()` are unbuffered

This could lead to different backpressure behaviors. The event channels from detector should probably be buffered to avoid blocking the detection goroutine.

### 3. Magic numbers scattered throughout
- `150 * time.Millisecond` × 3 in startup
- `500 * time.Millisecond` before displayCard
- `2 * time.Second` removal delay
- `256` default buffer size (also in config)
- `100` files progress throttle
- `256 * 1024` XMP buffer

These should be named constants with explanations of why those values were chosen.

---

## Concurrency Issues

### 4. The `displayCard` goroutine leak risk
```go
go func() {
    time.Sleep(500 * time.Millisecond)
    a.displayCard(cardPath)
}()
```

If the card is removed during that 500ms sleep, the goroutine still runs and may print to stdout after the removal message. The `isCurrentCard` check helps but it's a race.

**Fix:** Pass a context that can be cancelled on card removal.

### 5. `printMu` is only used in copy
The output mutex protects copy progress but not scan progress. If a card is inserted during analysis of another card, their outputs can interleave.

**Fix:** Use `printMu` in `displayCard` progress callback too.

### 6. `lastUpdate` capture in copy progress
The progress callback in `copyAll` captures `lastUpdate` from the outer scope:
```go
lastUpdate := time.Now()
go func() {
    // ... mutates lastUpdate from here
}()
```

This works but is fragile. If the goroutine panics, `lastUpdate` may be in an undefined state.

---

## Error Handling

### 7. Silent failures in filepath.WalkDir
Both analyze and copy use:
```go
if err != nil {
    return nil  // silently skip
}
```

Permission denied, broken symlinks, or I/O errors are silently ignored. This could hide real problems (dying card, corrupted filesystem).

**Suggestion:** Log warnings for walk errors, or collect them in the result.

### 8. `friendlyErr` is underutilized
Many error paths still show raw errors:
- Eject errors use `friendlyErr` ✓
- Copy errors use `friendlyErr` ✓
- Config load errors show raw `%v`
- Dotfile write errors show raw `%v`

**Fix:** Route all user-facing errors through `friendlyErr`.

### 9. No validation of destination path
Config accepts any string. `/dev/null`, ``, or `/root/foo` (without permissions) will all fail confusingly at copy time.

**Suggestion:** Validate destination on config load — check it's a valid path format and writable.

---

## Code Duplication

### 10. Stub command pattern is repeated 3 times
```go
case "s":
    fmt.Println("\nCopy Selects is not yet available.")
    a.printPrompt()
case "p":
    fmt.Println("\nCopy Photos is not yet available.")
    a.printPrompt()
case "v":
    fmt.Println("\nCopy Videos is not yet available.")
    a.printPrompt()
```

**Suggestion:** Extract a `notYetAvailable(cmd string)` helper.

### 11. Card header printing duplicated
`printCardInfo` and `printInvalidCardInfo` both print:
- Status
- Path
- Storage
- Camera

But with slightly different formatting. This is a maintenance risk.

**Suggestion:** Extract `printCardHeader(card)` that's used by both.

---

## UX Polish

### 12. Progress format could be clearer
Current:
```
Copying... 1247/3051 files  48.2 GB/96.4 GB (50%)
```

Issues:
- Files and bytes are both shown — pick one primary metric
- No ETA (you have the data to calculate it)
- No indication of speed

**Suggestion:** Simplify to one line with the most relevant info:
```
Copying... 1247/3051 files, ~3m remaining (150 MB/s)
```

### 13. The "Already copied" message is abrupt
```
Already copied.
```

No guidance on what to do next. Consider:
```
Already copied on 2026-03-12. [e] Eject  [x] Exit  [?]
```

### 14. No visual distinction between copy modes
When `[s]`, `[p]`, `[v]` are implemented, the user won't know which mode they're in from the prompt. The prompt always shows `[a] Copy All`.

**Suggestion:** Dynamic prompt that shows the active mode.

---

## Testing Gaps

### 15. `main.go` has 0% coverage
27 functions, 941 lines, not testable without refactoring. The business logic is intertwined with side effects (fmt.Print, os.Exit, signal handling).

**Path forward:**
1. Extract `runApp()` that returns errors instead of exiting
2. Make `app` methods take interfaces instead of concrete types
3. Use dependency injection for detector, logger, etc.

### 16. Race conditions are hard to test
The copy cancellation, card removal mid-copy, and queue handling all involve concurrent code. Currently only tested manually.

**Suggestion:** Add race-detection tests with `t.Parallel()` and goroutine synchronization.

### 17. Platform-specific code is untested
- `detect_darwin.go` — uses CGO, DiskArbitration
- `detect_linux.go` — uses sysfs
- `speedtest_darwin.go` — uses IOKit

These need integration tests on actual hardware, or mock interfaces.

---

## Minor Issues

### 18. `version` constant is untyped
```go
const version = "0.1.6"
```

Should be:
```go
const version string = "0.1.6"
```

(Though Go infers it correctly, explicit types are clearer.)

### 19. Comment drift: "remove in 0.4.0"
```go
// UX delays — remove in 0.4.0 when real startup and analysis timings replace them.
const (
    removalDelay = 2 * time.Second
)
```

But the sleeps at lines 198-202 are the ones to remove. `removalDelay` is actually reasonable UX.

### 20. `cardInvalid` is a negative name
```go
cardInvalid bool // true when current card has no DCIM directory
```

Double negatives are confusing: `if !invalid`. Better: `cardValid` or `hasDCIM`.

---

## Summary — Top 5 to Address

| Priority | Issue | Why It Matters |
|----------|-------|----------------|
| 1 | `displayCard` goroutine race | Can print after card removal, confusing UX |
| 2 | Silent walk errors | Hides filesystem/corruption issues |
| 3 | No destination validation | Fails late with confusing errors |
| 4 | 450ms startup delay | Unnecessary, hurts perceived performance |
| 5 | `main.go` not testable | Blocks confidence in refactoring |
