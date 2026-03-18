# 0.5.1 QA Runbook — Sleep/Wake Duplicate Launch Validation

Goal: validate `0.5.0` daemon reliability after sleep/wake and confirm whether any patch fix is needed for `0.5.1`.

---

## Scope

Covers QA follow-up item:
- **051-002** Sleep/wake duplicate launch behavior on real hardware

This runbook does **not** replace automated tests; it captures physical device behavior.

---

## Prerequisites

- macOS machine
- CardBot repo checked out
- A real SD/CFexpress card or equivalent test media
- Terminal app allowed to run CardBot normally

Optional:
- Disable any extra CardBot instances before starting (for cleaner logs)

---

## Quick Start (capture script)

From repo root:

```bash
./scripts/qa_051_sleepwake_capture.sh
```

The script:
- builds `cardbot`
- starts `cardbot --daemon`
- tails daemon log live
- captures log to `/tmp/cardbot-qa-051-<timestamp>/daemon.log`

Stop with `Ctrl+C` when done.

---

## Manual Steps

1. Start capture script.
2. Insert card once.
   - Expect one launch event (or one skip if foreground CardBot already running).
3. Put Mac to sleep with card inserted.
4. Wake Mac.
5. Remove and reinsert card quickly (1–3 times).
6. Observe daemon logs and launched windows.

---

## Pass/Fail Criteria

### PASS

- No launch storm after wake.
- Duplicate mount churn is suppressed (log contains suppression line when appropriate).
- Reinsert after real removal still allows a fresh launch.

### FAIL

- Multiple unexpected windows spawn from one logical insert cycle.
- Daemon repeatedly triggers launch for identical rapid events despite cooldown.
- Reinsert never launches after removal (over-suppression).

---

## Log Signals to Look For

Expected useful lines:

- `CardBot daemon started — watching for cards...`
- `Card detected: ...`
- `Suppressing duplicate card event for ... (cooldown)`
- `CardBot already running in another process — skipping auto-launch`
- `Launch failed: ...` (if permissions/automation issues)

---

## Record Result

After run, capture:

- macOS version
- terminal app used (`daemon.terminal_app`)
- card type/media reader
- PASS/FAIL
- relevant log snippets from generated log file

If FAIL, file a 0.5.1 bug item in `agent/051_qa_fixes.md` with reproduction steps.
