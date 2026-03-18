# v0.5.0 Draft Release Notes

## CardBot v0.5.0 — Background Auto-Launch

CardBot now supports a full background daemon workflow on macOS.

### Highlights

- **Daemon mode:** `cardbot --daemon`
  - Watches for card insertion events in the background
  - Launches configured terminal app when a card is detected
- **Direct path targeting:** `cardbot /Volumes/<CARD>`
  - Immediately analyzes a specific mounted card/path
- **LaunchAgent integration:**
  - `cardbot install-daemon`
  - `cardbot uninstall-daemon`
- **Setup improvements (`cardbot --setup`):**
  - Configure daemon enable/disable
  - Configure start-at-login
  - Choose terminal app for daemon launches
- **Status command:**
  - `cardbot daemon-status`
  - `cardbot daemon-status --json`
- **Reliability guards:**
  - Single-instance check to avoid opening duplicate CardBot windows
  - Duplicate-event cooldown to suppress rapid reinsert/sleep-wake event storms
- **Better diagnostics:**
  - Launch failure hints for Automation / Full Disk Access permission issues

### New Config Fields

```json
"daemon": {
  "enabled": false,
  "start_at_login": false,
  "terminal_app": "Terminal",
  "launch_args": []
}
```

### Notes

- LaunchAgent workflow is macOS-only.
- `daemon.launch_args` supports templates:
  - `{{cardbot_binary}}`
  - `{{mount_path}}`

### Upgrade Notes

- Re-run setup after upgrade:

```bash
cardbot --setup
```

- Check current daemon state:

```bash
cardbot daemon-status --json
```

### Known Follow-up

- Final manual sleep/wake + permission QA on physical hardware before tagging stable release.
