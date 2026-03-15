# CardBot 0.3.1 — Renaming Flow Polish

Focus: tighten 0.3.0 behavior after real-world use, no video workflow changes.

## Scope

- Stabilize timestamp renaming UX
- Fix dry-run behavior and side effects
- Improve confidence with app-level tests
- Keep config schema migration explicitly planned for future (deferred for now)

## P0 (Must Fix Before Release)

- [ ] **Dry-run must be side-effect free in app flow**
  - Current risk: dry-run path still goes through success handling in `copy_cmd.go`
  - Ensure dry-run does **not**:
    - write `.cardbot`
    - mark `copiedModes[mode] = true`
    - print "Copy complete ✓"
  - Should print `Dry-run complete` summary instead

- [ ] **Add app-level test coverage for dry-run side effects**
  - Test: no dotfile write attempted in dry-run
  - Test: mode is not marked copied after dry-run
  - Test: prompt returns to normal copy options

## P1 (High Value)

- [ ] **Large-card dry-run preview ergonomics**
  - Add preview line cap (e.g., first 200 mappings)
  - Print summary: `... +2848 more files`
  - Optional `--dry-run-full` later if needed

- [ ] **Timestamp re-copy behavior with renamed files**
  - Evaluate duplicate-copy behavior on re-inserted card in timestamp mode
  - Decide for 0.3.1:
    - document as known limitation, or
    - lightweight skip strategy

- [ ] **Refine copy completion language**
  - Distinguish clearly between:
    - real copy complete
    - dry-run preview complete

## Deferred (Not in 0.3.1, keep visible)

- [ ] **Config schema migration harness (prep for v3)**
  - Add explicit migration switch in `config.Load()`
  - Add migration tests:
    - unknown schema warning behavior
    - v2 -> v3 mapping (when introduced)
    - preserve user values where possible
  - Trigger only when schema changes are actually introduced

## P2 (Future-Proofing)

- [ ] **Docs consistency updates**
  - Ensure config path docs match `os.UserConfigDir()` behavior
  - Keep examples aligned to NEF/MOV and 0.3.x output

## Real-World Test Matrix (Before Tag)

- [ ] Z9 card, naming mode `timestamp`, copy all
- [ ] Z9 card, naming mode `timestamp`, copy photos only
- [ ] Z9 card, dry-run preview readability on a large card
- [ ] Mixed NEF + MOV card, verify sequence and extensions
- [ ] Reinsert same card and run again (check duplicate behavior)

## Acceptance Criteria

- Dry-run causes zero card/destination state mutation
- All existing tests pass + new dry-run tests
- At least one successful real-world validation run on actual card media
