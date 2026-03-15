# CardBot 0.3.2 — Naming Choice UX Review

## Current Flow (0.3.1)

### First-Run Setup

```
Welcome to CardBot!

Where should CardBot copy your work?

[folder picker opens, user selects ~/Pictures/CardBot]
Destination: /Users/alex/Pictures/CardBot
────────────────────────────────────────
File Naming
────────────────────────────────────────
How would you like files named when copying?

[1] Keep original filenames (DSC_0001.NEF)
[2] Timestamp + sequence (260314T143052_001.NEF)

The timestamp comes from when each photo was taken.
Sequence digits adjust automatically based on card size (3/4/5 digits).
You can change this later with cardbot --setup.

Choice [1]:
```

### Startup Confirmation

```
[2026-03-14T10:30:46] Starting CardBot 0.3.1...
[2026-03-14T10:30:46] Copy path ~/Pictures/CardBot
[2026-03-14T10:30:46] Keep original filenames
```

Or if timestamp mode:

```
[2026-03-14T10:30:46] Timestamp + sequence filenames (YYMMDDTHHMMSS_NNN.EXT, auto 3/4/5 digits)
```

### Card Info Display

```
  Status:   New
  Path:     /Volumes/NIKON Z 9
  Storage:  128.4 GB / 476.9 GB (27%)
  Camera:   Nikon Z 9
  Starred:  3
  Content:  2026-03-14      45.2 GB    842   NEF, MOV
            2026-03-13      12.1 GB    234   NEF

  Total:    1076 photos, 18 videos, 128.4 GB
  Naming:   Keep original (DSC_0001.NEF)
  Dest:     ~/Pictures/CardBot
────────────────────────────────────────
[a] Copy All  [e] Eject  [x] Exit  [?] Help  >
```

Or timestamp mode:

```
  Naming:   Timestamp sequence (260314T143052_0001.NEF) [4-digit]
```

---

## Issues Identified

### 1. Setup Prompt — Vague "How would you like"
**Current:** "How would you like files named when copying?"

**Problem:**
- Doesn't explain *why* you'd choose one over the other
- Photographer with 50k files may not grasp the filing benefit immediately
- No mention of burst handling (the actual pain point timestamp mode solves)

### 2. Option [1] Label — "Keep original" sounds like "do nothing"
**Current:** "Keep original filenames (DSC_0001.NEF)"

**Problem:**
- Sounds like the "safe" default, but doesn't explain when it's problematic
- No mention of the 9999 rollover problem (DSC_9999.NEF → DSC_0001.NEF)
- No mention of multiple camera collision

### 3. Option [2] Label — Technical, not benefit-driven
**Current:** "Timestamp + sequence (260314T143052_001.NEF)"

**Problem:**
- "Timestamp" implies wall-clock, not capture-time
- "Sequence" is vague — sequence of what?
- The example format is cryptic (YYMMDDTHHMMSS_NNN.EXT)
- No mention of the actual benefit: chronological ordering across cards/folders

### 4. The Explanation — Buried, dense
**Current:** "The timestamp comes from when each photo was taken. Sequence digits adjust automatically based on card size (3/4/5 digits)."

**Problem:**
- Two sentences, both parenthetical
- "3/4/5 digits" is unexplained mechanics
- Doesn't answer: "What problem does this solve for me?"

### 5. Card Info — "Naming" line is cramped
**Current:** "Naming: Timestamp sequence (260314T143052_0001.NEF) [4-digit]"

**Problem:**
- Timestamp format is noise
- "[4-digit]" is unexplained
- Doesn't show current mode in intuitive terms

### 6. No Quick Toggle
**Current:** Must run `cardbot --setup` to change

**Problem:**
- User may want to toggle per-card (wedding vs. personal)
- No in-session discovery of the option

---

## Proposed Flow (0.3.2)

### First-Run Setup — Revised

```
Welcome to CardBot!

Where should CardBot copy your work?
Destination: /Users/alex/Pictures/CardBot

────────────────────────────────────────
How should files be named?
────────────────────────────────────────

Camera filenames reset every 10,000 shots (DSC_9999.NEF → DSC_0001.NEF).
This causes duplicates and loses chronological order across cards.

[1] Keep camera filenames
    DSC_0001.NEF, DSC_0002.NEF, ...
    Use this if you rely on camera numbering for your workflow.

[2] Chronological capture time
    260314T103045_001.NEF, 260314T103046_002.NEF, ...
    Use this for automatic chronological order across all cards.
    Perfect for events and multi-camera shoots.

Choice [2]: 2
Naming set to: Chronological capture time
```

### Startup Confirmation — Revised

Original mode:
```
[2026-03-14T10:30:46] Naming: Camera original (DSC_xxxx.NEF)
```

Timestamp mode:
```
[2026-03-14T10:30:46] Naming: Capture time + sequence (auto-ordered)
```

### Card Info — Revised

Original mode:
```
  Naming:   Camera original (DSC_xxxx.NEF)
```

Timestamp mode:
```
  Naming:   Capture time + sequence (4-digit, ~3,000 files/card)
```

Or with actual card stats:
```
  Naming:   Capture time + sequence (3-digit, fits 999 files)
```

---

## Additional 0.3.2 Ideas

### In-Session Toggle (Experimental)

Add `[n]` command to toggle naming mode for current card only:

```
[a] Copy All  [n] Switch to capture-time  [e] Eject  [x] Exit  [?] Help  >
```

Pressing `[n]` would:
1. Toggle mode for this card only
2. Show confirmation: "Switched to capture-time naming for this card"
3. Update the `Naming:` line in card info
4. Not persist to config (session-only)

### Dry-Run Preview Context

Before showing mappings, add context line:

```
[2026-03-14T10:30:46] Dry-run: would copy 1,076 files to ~/Pictures/CardBot
Naming mode: Capture time + sequence (4-digit)

DSC_8234.NEF → 260314T103045_0001.NEF
DSC_8235.NEF → 260314T103046_0002.NEF
...
```

---

## Comparison Matrix

| Aspect | Current (0.3.1) | Proposed (0.3.2) |
|--------|-----------------|------------------|
| Setup prompt | "How would you like..." | "How should files be named?" |
| Option 1 label | "Keep original filenames" | "Keep camera filenames" + note on 9999 rollover |
| Option 2 label | "Timestamp + sequence" | "Chronological capture time" + event/multi-cam benefit |
| Setup explanation | Technical (timestamp source, digit count) | Problem/solution framing |
| Startup line | "Timestamp + sequence filenames (YYMMDDTHHMMSS_NNN.EXT, auto 3/4/5 digits)" | "Capture time + sequence (auto-ordered)" |
| Card info | "Timestamp sequence (260314T143052_0001.NEF) [4-digit]" | "Capture time + sequence (4-digit, ~3,000 files/card)" |
| Session toggle | Not available | `[n]` toggle per-card (optional) |

---

## Open Questions

1. **Terminology:** Is "capture time" clearer than "timestamp"? Alternatives: "shot time", "photo time", "EXIF time"

2. **Option order:** Should chronological be [1] and camera original [2]? (Currently reversed)

3. **Session toggle:** Worth adding `[n]` command or keep config-only?

4. **Per-card persistence:** Should a session toggle affect just this card, or all cards until restart?

5. **Visual indicator:** Should the prompt show current mode visually? E.g., `[a] Copy All (chronological)`
