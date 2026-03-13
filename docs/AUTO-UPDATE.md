# Auto-Update Plan

Goal: make updates easy for non-technical users without adding risky silent behavior.

Status: planning only (wishlist).

---

## Principles

- Keep updates **safe and explicit**.
- Prefer a clear prompt + one command over silent background replacement.
- Verify downloaded binaries before replacing anything.
- Fail gracefully with helpful instructions.

---

## Phase 1 — Update Check (low complexity)

Add startup/update-check support:

- Query latest version from GitHub Releases (`/releases/latest`)
- Compare against current `version`
- If newer, print a friendly message:
  - `Update available: 0.2.1 (you have 0.2.0)`
  - `Run: cardbot self-update`
- Add flags:
  - `--check-update` (manual check)
  - `--no-update-check` (skip check)

### Acceptance criteria

- Works without blocking app startup if network/API fails
- Times out quickly (e.g. 2–3s)
- No noisy errors for offline users

---

## Phase 2 — `self-update` Command (medium complexity)

Add `cardbot self-update` command:

1. Detect platform (`darwin-arm64`, `darwin-amd64`, etc.)
2. Download matching release binary
3. Download/verify SHA256 checksum
4. Replace current binary atomically
5. Keep executable permissions
6. Print clear success/failure output

### Permission behavior

- If binary path is not writable, print exact `sudo` command the user can run.
- Never leave a partially replaced binary.

### Acceptance criteria

- Safe rollback on failed replace
- Checksum mismatch aborts update
- Works for common install locations (`/usr/local/bin`, local folder)

---

## Phase 3 — Optional Auto-Apply (deferred)

Not needed now.

Potential future work:

- Auto-apply updates on startup
- Signed release verification
- Rollback cache and telemetry

This phase is intentionally deferred until user base + stability justify it.

---

## Rough Effort

- Phase 1: ~0.5 day
- Phase 2: ~1–2 days
- Phase 3: several days + ongoing maintenance

---

## Suggested Implementation Layout

- `internal/update/check.go` — release query + semver compare
- `internal/update/download.go` — fetch binary/checksum
- `internal/update/install.go` — atomic replace + permissions
- `internal/update/update_test.go` — unit tests for compare, URL selection, checksum

---

## Out of Scope (for now)

- GUI updater
- Homebrew formula auto-bump
- Delta/binary patch updates
- Windows installer support
