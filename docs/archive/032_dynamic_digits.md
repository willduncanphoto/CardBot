# CardBot 0.3.2 — Dynamic Per-Date Sequence Digits

## Current Behavior (0.3.1)

**Problem:** Uses single digit count for entire card based on total file count.

```
Card with 3,048 files:
- 2026-03-14: 842 files
- 2026-03-13: 234 files  
- 2026-03-12: 1,972 files (big event day)

Current: Uses 4-digit for ALL files because total > 999
         260314T103045_0001.NEF (even for dates with <1000 files)
         
Wasted:  260312T083022_0001.NEF for a day with only 234 files
```

## Proposed Behavior (0.3.2)

**Per-date digit detection:** Each date folder gets minimum digits needed.

```
Same card:
- 2026-03-14: 842 files → 3-digit (001-842)
- 2026-03-13: 234 files → 3-digit (001-234)
- 2026-03-12: 1,972 files → 4-digit (0001-1972)

Result: 260314T103045_001.NEF (clean 3-digit for small days)
        260312T083022_0001.NEF (4-digit only for big day)
```

## Implementation Requirements

### 1. Per-Date Analysis During Scan

Analyzer needs to compute per-date file counts:

```go
type DateGroup struct {
    Date       string
    Size       int64
    FileCount  int
    Extensions []string
    SeqDigits  int  // NEW: computed during analysis
}
```

During analysis:
```go
for each date group:
    group.SeqDigits = sequenceDigits(group.FileCount)
    // 3, 4, or 5 based on that day's count
```

### 2. Copy Engine Uses Per-Date Digits

Instead of single `seqDigits` for whole card:

```go
// Current (0.3.1)
seqDigits := SequenceDigits(len(files)) // one value for all

// Proposed (0.3.2)
seqDigitsByDate := make(map[string]int)
for _, g := range analyzeResult.Groups {
    seqDigitsByDate[g.Date] = SequenceDigits(g.FileCount)
}

// During copy:
destRelPath := renamedRelativePath(f.relPath, f.captureTime, seq, seqDigitsByDate[f.date])
```

### 3. Re-Copy Prevention via Mapping Log

**Current Problem:** Re-copying same card gets new sequence numbers because rename mapping isn't checked.

**Solution:** Use logged mapping to skip already-copied files.

Dotfile v2 structure (already implemented):
```json
{
  "$schema": "cardbot-dotfile-v2",
  "copies": [{
    "mode": "all",
    "timestamp": "2026-03-14T10:30:46Z",
    "files_copied": 3048,
    "bytes_copied": 128400000000
  }]
}
```

**Need:** Extend dotfile with rename mapping (or use separate log):

```json
{
  "$schema": "cardbot-dotfile-v3",
  "copies": [{
    "mode": "all",
    "timestamp": "2026-03-14T10:30:46Z",
    "files_copied": 3048,
    "renaming": {
      "enabled": true,
      "mappings": [
        {"src": "100NIKON/DSC_8234.NEF", "dst": "260314T103045_001.NEF"},
        {"src": "100NIKON/DSC_8235.NEF", "dst": "260314T103046_002.NEF"}
      ]
    }
  }]
}
```

Re-copy logic:
```go
// Check if source file already has a mapping
if existing, ok := dotfileMappings[srcRelPath]; ok {
    // Verify destination exists with correct size
    if destExists(existing.dst) && destSize(existing.dst) == srcSize {
        skip // already copied with this name
    }
}
```

### 4. Display Update

Card info shows per-date sequence info:

```
  Content:  2026-03-14      45.2 GB    842   NEF, MOV  [3-digit]
            2026-03-13      12.1 GB    234   NEF       [3-digit]
            2026-03-12      71.0 GB   1972   NEF, MOV  [4-digit]
```

## Edge Cases

### Mixed Digit Counts Same Card
User copies card once, gets:
- Day 1: 3-digit
- Day 2: 4-digit

Re-copies card next week:
- Same mappings should be detected via dotfile
- No duplicate files with different sequence numbers

### Card Re-insert Without Dotfile
If `.cardbot` was deleted or card was reformatted:
- Re-analysis computes digits fresh
- May get different sequence numbers for same files
- **Acceptable:** rare case, document limitation

### Sequence Overflow Mid-Copy
What if copy is interrupted and resumed:
- Interrupted at file 950 of 1000 (4-digit day)
- Resume: still 4-digit, sequence continues
- Dotfile mapping ensures no gaps/duplicates

## UX Clarification: Session Toggle

User asked: "3 session toggle for what? I dont get it?"

I was suggesting a command `[n]` during card interaction to toggle naming mode without running `--setup`. Example:

```
[a] Copy All  [n] Use timestamp naming  [e] Eject  [x] Exit  [?] Help  >
```

**Decision:** Skip this. Config-only via `--setup` is simpler. Per-card toggle adds complexity without clear benefit.

## Migration Path

**Dotfile v2 → v3:**
- Add optional `renaming` object to copy entries
- Backward compatible: v2 readers ignore new field
- Migration: copy engine checks for both formats

## Open Questions

1. **Mapping storage:** Store full mappings in dotfile, or separate `.cardbot-map` file?
   - Dotfile: simpler, one file per card
   - Separate: keeps dotfile small for quick status checks

2. **Mapping size limit:** Cap at N entries to prevent huge files?
   - 10,000 entries ≈ 500KB JSON
   - Acceptable for most cards

3. **Hash-based instead of name-based:** Should we hash file content to detect "same file, different name"?
   - More robust but slower (must read file)
   - Probably overkill for 0.3.2
