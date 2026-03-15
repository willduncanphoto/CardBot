# CardBot 0.3.0 — File Renaming Design Plan

## Overview

CardBot 0.3.0 introduces file renaming on copy. Transform cryptic camera filenames
(`DSC_0001.NEF`, `_DSC0002.JPG`) into organized, meaningful names using EXIF data,
camera identification, and custom sequences.

**Key principle:** The original filename on the card never changes. Renaming happens
only during copy to the destination. The `.cardbot` dotfile tracks the mapping for
recovery/debugging.

---

## User Stories

1. **Event photographer:** "I want files named `2026-03-14_Wedding_Smith_001.NEF`
   so I can find shots by date and client without digging through folders."

2. **Travel shooter:** "I want `Z9_20260314_0001.NEF` so I know which camera body
   and when each shot was taken."

3. **Studio workflow:** "I need sequential numbering per session that continues
   across multiple cards shot the same day."

4. **Safety first:** "I want to preview what files will be named before I commit
to the copy."

---

## Configuration Schema

```json
{
  "$schema": "cardbot-config-v2",
  "destination": {
    "path": "~/Pictures/CardBot"
  },
  "renaming": {
    "enabled": false,
    "template": "{date}_{camera}_{seq}.{ext}",
    "date_format": "YYYY-MM-DD",
    "sequence": {
      "scope": "card",
      "start": 1,
      "padding": 4
    },
    "collision": "skip",
    "templates": [
      {
        "name": "Date + Camera + Sequence",
        "template": "{date}_{camera}_{seq}.{ext}"
      },
      {
        "name": "Date Only",
        "template": "{date}_{seq}.{ext}"
      },
      {
        "name": "Camera + Sequence",
        "template": "{camera}_{seq}.{ext}"
      },
      {
        "name": "Original + Date",
        "template": "{original}_{date}.{ext}"
      }
    ]
  }
}
```

---

## Template Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `{date}` | EXIF DateTimeOriginal | `2026-03-14` |
| `{camera}` | Clean camera brand/model | `Z9`, `R5`, `A7IV` |
| `{seq}` | Sequence number | `0001`, `0042` |
| `{ext}` | Lowercase extension | `nef`, `jpg`, `mov` |
| `{original}` | Original filename (no ext) | `DSC_0001`, `_DSC0002` |

### Date Formats

| Format | Output |
|--------|--------|
| `YYYY-MM-DD` | `2026-03-14` |
| `YYYYMMDD` | `20260314` |
| `YYMMDD` | `260314` |
| `YYYY-MM` | `2026-03` |

### Sequence Scopes

| Scope | Behavior |
|-------|----------|
| `card` | Resets to `start` for each new card |
| `global` | Persists across cards in the same session |
| `date` | Resets when EXIF date changes |

---

## UI/UX Design

### New Prompt Keys

| Key | Action |
|-----|--------|
| `r` | **Rename settings** — open template chooser/interactive config |
| `n` | **Toggle renaming** — enable/disable for current session (default: off) |

### Card Display (when renaming enabled)

```
  Status:   New
  Path:     /Volumes/NIKON Z 9
  Storage:  96.4 GB / 476.9 GB (20%)
  Camera:   Nikon Z 9
  Starred:  1
  Content:  2026-02-27      12.9 GB    418   NEF
            2026-02-26      28.4 MB      1   NEF

  Total:    3048 photos, 0 videos, 96.0 GB
  Renaming: ON → {date}_{camera}_{seq}.{ext}
────────────────────────────────────────
[a] Copy All  [s] Copy Selects  [p] Copy Photos  [v] Copy Videos  [r] Rename Settings  [n] Toggle Renaming  [e] Eject  [x] Exit
```

### Rename Settings Menu (`[r]`)

```
Rename Settings
===============
Current: OFF

[1] Date + Camera + Sequence  →  2026-03-14_Z9_0001.nef
[2] Date Only                 →  2026-03-14_0001.nef
[3] Camera + Sequence         →  Z9_0001.nef
[4] Original + Date           →  DSC_0001_2026-03-14.nef
[c] Custom template...
[d] Date format: YYYY-MM-DD
[s] Sequence scope: card
[p] Padding: 4 digits
[o] Collision: skip

[x] Back
> 1
Renaming enabled with template: {date}_{camera}_{seq}.{ext}
```

### Custom Template Input

```
Enter custom template (variables: {date} {camera} {seq} {ext} {original}):
> {date}_{client}_{seq}.{ext}

Preview with sample file:
  DSC_0001.NEF  →  2026-03-14_Wedding_0001.nef

Save this template? [y/n] > y
```

### During Copy — Filename Ticker (`[t]`)

Press `[t]` during copy to toggle filename display:

```
[2026-03-14T10:23:15] Copying (All)...
  156/3048 files  4.2 GB  85 MB/s
  Current: 2026-03-14_Z9_0156.nef

[\] Cancel  [t] Hide filename  [p] Pause
```

---

## Dotfile v3 Schema

Track original → renamed mapping for recovery and debugging:

```json
{
  "version": 3,
  "copies": [
    {
      "mode": "all",
      "timestamp": "2026-03-14T10:23:15Z",
      "renamed": true,
      "template": "{date}_{camera}_{seq}.{ext}",
      "mappings": [
        {
          "original": "100NIKON/DSC_0001.NEF",
          "renamed": "2026-03-14_Z9_0001.nef",
          "hash": "sha256:abc123..."
        }
      ]
    }
  ]
}
```

**Migration:** v2 → v3 adds `renamed` boolean and optional `mappings` array.
Mappings are truncated after 1000 entries to prevent dotfile bloat (full mapping
available in log file).

---

## Collision Handling

When a destination file already exists:

| Mode | Behavior |
|------|----------|
| `skip` | Don't copy, count as skipped (existing behavior) |
| `overwrite` | Replace existing file (with confirmation if size differs) |
| `suffix` | Add incrementing suffix: `_0001`, `_0002` |

---

## Implementation Plan

### Phase 1: Core Infrastructure
- [ ] Add `renaming` section to config schema
- [ ] Create `internal/rename/` package with template parser
- [ ] Template variable extraction from EXIF data
- [ ] Sequence number management (per-card, global, per-date)

### Phase 2: UI Integration
- [ ] Add `[r]` rename settings command
- [ ] Add `[n]` toggle command
- [ ] Display rename status in card info
- [ ] Interactive template chooser

### Phase 3: Copy Integration
- [ ] Hook rename into copy engine
- [ ] Filename ticker (`[t]` during copy)
- [ ] Collision handling
- [ ] Preview mode (dry-run rename)

### Phase 4: Persistence
- [ ] Dotfile v3 with mappings
- [ ] Migration from v2
- [ ] Log file with full mapping

---

## Open Questions

1. **Client/job name variable?** Should `{client}` be a config setting or prompt
   at copy time? (Leaning: config setting, can change per session)

2. **Folder renaming?** Should we support renaming the dated folders too, or just
   files within them? (Leaning: keep folder structure, rename files only)

3. **Video file renaming?** Same templates or separate config? (Leaning: same,
   but video workflow 0.4.0 may split this)

4. **RAW+JPEG pairing?** Should sequences stay in sync for RAW+JPEG pairs shot
   together? (Leaning: yes, track base filename for pairing)

---

## Testing Checklist

- [ ] Template parsing with all variables
- [ ] Date format variations
- [ ] Sequence scopes (card/global/date)
- [ ] Padding (2, 3, 4, 5 digits)
- [ ] Collision modes
- [ ] Preview mode accuracy
- [ ] Dotfile v3 migration
- [ ] Filename ticker display
- [ ] Special characters in camera names
- [ ] Missing EXIF date fallback (file mod time?)

---

## Configuration Schema Migration Strategy

**Problem:** When we upgrade config schema (v1 → v2 → v3), existing user configs
break or get reset to defaults. We need a defined migration path.

**Current State:**
- v1: Original config (destination only, no schema field)
- v2: Current config with `$schema` field, naming mode, advanced settings
- v3: Future config with full renaming templates, sequence config, etc.

**Migration Rules:**
1. **Forward compatibility:** Newer CardBot reads old schemas and migrates on load
2. **No data loss:** User settings are preserved, new fields get defaults
3. **Version stamping:** Config is saved with current schema after migration
4. **Warnings:** User sees console message about migration:  
   `[2026-03-14T10:00:00] Config migrated from v2 to v3`

**Implementation Pattern:**
```go
// In config.Load()
cfg, migrated := migrateConfig(rawJSON)
if migrated {
    warnings = append(warnings, fmt.Sprintf("config migrated to %s", schemaVersion))
    // Save back immediately so user has migrated config on disk
    _ = Save(cfg, path)
}
```

**v2 → v3 Migration Mapping:**
| v2 Field | v3 Field | Transform |
|----------|----------|-----------|
| `naming.mode` | `renaming.enabled` + `renaming.template` | `"original"` → `enabled:false`; `"timestamp"` → `enabled:true, template:"{YY}{MM}{DD}T{HH}{mm}{ss}_{seq}.{ext}"` |
| (none) | `renaming.sequence.scope` | Default `"card"` |
| (none) | `renaming.collision` | Default `"skip"` |

**Open Issue:** What if user edited config while running older CardBot version?
- Option A: Always migrate on load (aggressive, may lose v3-specific edits)
- Option B: Preserve unknown fields and merge (complex)
- Option C: Warn and ask user to run `--setup` (safe, requires action)

**Recommendation:** Start with Option A for v2→v3 since v3 doesn't exist yet.
For future v3→v4, implement field-preserving merge.

---

## Future Enhancements (Post-0.3.0)

- Regex-based transformations on original filename
- Conditional templates (if starred, use different pattern)
- Batch rename existing copied files
- Export rename mapping as CSV
