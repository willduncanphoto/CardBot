# CardBot 0.3.2 — Simplified Scope

## What's In (Definite)

### 1. UX Polish Only
- Better setup prompt text (problem/solution framing)
- Clearer labels: "Camera original" vs "Timestamp + sequence"
- Simplified card info display
- Keep "timestamp" terminology

### 2. Fixed 3-Digit Sequence
- All cards use 3 digits (001-999)
- Simple, predictable, works for 99% of use cases
- No per-date detection complexity

## What's Out (Deferred)

### Per-Date Dynamic Digits
**Future idea:** If a single calendar day has >999 files, auto-expand to 4 digits.
**Why deferred:** Rare edge case, adds complexity, can be addressed later if users hit it.

### Multi-Camera Same Day
**Future idea:** When two cameras shoot the same event, timestamps can collide.
**Solutions to discuss later:**
- Camera ID prefix: `Z9_260314T103045_001.NEF`
- Subfolder per camera: `2026-03-14/Z9/`
- Per-camera sequence offset (camera A 001-500, camera B 501-999)

### Re-Copy Prevention via Mapping
**Future idea:** Store original→renamed mapping in dotfile to skip on re-copy.
**Why deferred:** With 3-digit fixed sequences, re-copy behavior is at least predictable (001-999 loop). Full mapping log is v3 feature.

---

## 0.3.2 Implementation Checklist

- [ ] Update `setup_naming.go` prompt text
- [ ] Update `namingStartupLine()` labels  
- [ ] Update `namingDisplayLine()` labels
- [ ] Fix sequence digits to always 3 (remove dynamic detection)
- [ ] Update tests for fixed 3-digit
- [ ] Update README example output

## Notes for Future (v4+)

### Sequence Behavior (Current: 3-digit, reset at 999)
```
Files 1-999:   _001.NEF through _999.NEF
File 1000:     _000.NEF (loops back to 000, not 001)
```
This is acceptable because:
- 999+ files from one camera in one day is rare
- Timestamp prefix still provides uniqueness per second
- Loop is predictable and documented

### Subsecond Timestamp Idea
If EXIF provides subsecond precision:
```
Current:   260314T143052_001.NEF
Subsecond: 260314T143052.947.NEF  ← no sequence needed?
```
**Blockers:**
- Not all cameras write subsecond EXIF (need to verify Z9 behavior)
- Bursts within same 1/1000s would still collide
- Changes filename format significantly

**Verdict:** Investigate for v4, not 0.3.x

### Multi-Camera Collision (Moved to 0.4.0 Roadmap)
See main TODO.md for full collision prevention strategy options.
