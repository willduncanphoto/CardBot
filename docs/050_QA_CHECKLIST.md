# 050 QA Checklist — Background Auto-Launch

Use this before shipping v0.5.0.

## 1) Automated Checks

Run:

```bash
make qa-050
```

Expected:
- `go test ./...` passes
- `cardbot daemon-status` text output includes daemon/guard/launch-agent sections
- `cardbot daemon-status --json` output includes expected keys

Status: ✅ available

---

## 2) Manual macOS Functional Checks

### A. Setup Flow
- [ ] Run `cardbot --setup`
- [ ] Choose daemon enabled = yes
- [ ] Choose start at login = yes
- [ ] Choose terminal app = Terminal (then rerun with Ghostty/custom)
- [ ] Verify setup summary reflects choices

### B. LaunchAgent Commands
- [ ] Run `cardbot install-daemon`
- [ ] Verify plist exists: `~/Library/LaunchAgents/com.illwill.cardbot.plist`
- [ ] Run `cardbot daemon-status` and verify:
  - LaunchAgent installed = enabled
  - LaunchAgent loaded = enabled
- [ ] Run `cardbot uninstall-daemon`
- [ ] Verify plist removed and status reflects disabled

### C. Daemon Runtime
- [ ] Start daemon: `cardbot --daemon`
- [ ] Insert supported card
- [ ] Verify terminal app opens CardBot for mounted path
- [ ] Insert/remove same card rapidly and confirm duplicate launch suppression
- [ ] Verify single-instance guard:
  - keep foreground CardBot open
  - insert card
  - daemon should log skip and not open extra window

### D. Direct Path Launch
- [ ] Run `cardbot /Volumes/<CARD>`
- [ ] Verify immediate analysis starts for that path (no waiting scan spinner)

---

## 3) Sleep/Wake Checks (macOS)

- [ ] Start daemon
- [ ] Put Mac to sleep with card inserted, then wake
- [ ] Confirm no launch storm / duplicate windows
- [ ] Remove + reinsert card after wake
- [ ] Confirm single launch behavior

---

## 4) Permission Checks

### Automation
- [ ] Revoke terminal automation permission (if granted)
- [ ] Insert card while daemon running
- [ ] Verify helpful launch hint mentions Automation in output

### Full Disk Access
- [ ] Restrict access scenario (or test on fresh machine)
- [ ] Trigger daemon launch failure path
- [ ] Verify hint mentions Full Disk Access

---

## 5) Status Command Checks

- [ ] `cardbot daemon-status`
- [ ] `cardbot daemon-status --json`
- [ ] Verify fields:
  - `version`
  - `pid`
  - `daemon`
  - `single_instance_guard`
  - `launch_agent`

---

## 6) Release Gate

Ship only when all manual boxes above are checked and automated checks are green.
