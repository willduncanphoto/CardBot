# CardBot — Active Work

See [ROADMAP.md](ROADMAP.md) for full future planning.

## Current: 0.3.2

**Status:** Ready for real-world validation.

### Before Tagging 0.3.2
- [ ] Real-world test: Z9 card, timestamp mode
- [ ] Verify 3-digit sequence behavior (001-999)
- [ ] Verify dry-run preview output
- [ ] Verify re-copy behavior (expected: may re-copy due to no mapping log)

### Known Limitations (Acceptable for 0.3.2)
- 1000+ files/day: sequence loops (001→999→001). Rare, documented.
- Multi-camera same second: collision risk. See ROADMAP "Stuff to Think About".
- Re-copy: uses size check, may re-copy renamed files.

---

## Next Version: TBD

See ROADMAP.md for candidates:
- Video workflow separation
- Config schema v3 migration harness
- Linux platform support

Multi-camera collision prevention is **parked** — not 0.4.0 priority.

---

## UI/UX Updates (Future)

- [ ] **Technical EXIF Display Mode**
  - Show raw EXIF values in card info:
    ```
    Make    : NIKON CORPORATION
    Model   : NIKON Z 9
    ```
  - Instead of cleaned "Nikon Z 9"
  - Toggle or config option for technical vs friendly display
  
- [ ] **Card Info Layout Refresh**
  - More technical/professional appearance
  - Raw EXIF values where meaningful
  - Cleaner alignment

## Quick Fixes (Any Release)

- [ ] Add "OM System" brand color (when confirmed)
- [ ] Config path display command (`cardbot --config`)

---

## Done

See ROADMAP.md for completed milestones.
