# CardBot Roadmap

## 0.4.0 (Current)

**Focus:** Startup UX polish

### Completed
- ✅ Animated update check spinner on startup
- ✅ `Up to date` / `UPDATE AVAILABLE (x.x.x)` messaging
- ✅ `NO SIGNAL` with error code for failed update checks
- ✅ Classic `| / - \` spinner for update check and scanning
- ✅ Simplified naming display (removed verbose explanations)
- ✅ Silent startup when already up to date

### Inherited from 0.3.x
- ✅ Spinner animation for idle scanning
- ✅ Clean output formatting
- ✅ Cross-platform builds (darwin arm64/amd64, linux amd64/arm64)
- ✅ Architecture refactor: move root logic to `internal/app/`
- ✅ Self-update system with checksum verification

---

## 0.4.1 (Next)

**Focus:** Housekeeping

- [ ] Code style consistency pass
- [ ] Remove dead code or unused helpers
- [ ] Standardize error message formatting
- [ ] Improve test coverage for edge cases
- [ ] Documentation cleanup

---

## 0.5.0 (Future)

**Candidates:**
- [ ] Video workflow separation (photos → Pictures, videos → Movies)
- [ ] Batch operations for multiple cards
- [ ] Configuration presets/profiles
- [ ] Performance profiling of large card scans
- [ ] Terminal resize handling improvements

---

## Parking Lot

### Multi-Camera Collision Prevention
Two cameras shooting same event can produce identical filenames in timestamp mode.
Needs design work. Not blocking current releases.

### Linux/Windows Support  
Linux implemented but needs real-world testing. Windows is long-term.

### Single-Key Input
Raw terminal mode (no Enter required). Power user polish.

---

## Completed Milestones

| Version | Date | Highlights |
|---------|------|------------|
| 0.4.0 | 2026-03-16 | Startup UX polish, update check spinner, simplified display |
| 0.3.x | 2026-03-15 | Timestamp renaming, architecture refactor, self-update |
| 0.2.x | 2026-03-12 | Selective copy (starred, photos, videos) |
| 0.1.x | 2026-03-08 | Core copy engine, card detection, EXIF analysis |
