# Uninstalling CardBot

## Quick Teardown (instant)

```bash
# Kill background daemon
pkill -f "cardbot --daemon"

# Remove LaunchAgent
launchctl bootout gui/$(id -u)/com.illwill.cardbot
rm -f ~/Library/LaunchAgents/com.illwill.cardbot.plist

# Remove binary
rm ~/bin/cardbot

# Optional: purge config + logs
rm -rf ~/Library/Application\ Support/cardbot ~/.cardbot
```

---

## With Uninstall Script

```bash
# Full uninstall (daemon + binary)
sh uninstall.sh --install-dir ~/bin

# Full uninstall + purge config + logs
sh uninstall.sh --install-dir ~/bin --purge
```

### Uninstall Script Options

| Flag | Description |
|------|-------------|
| `--install-dir <path>` | Additional directory to remove `<path>/cardbot` |
| `--no-sudo` | Skip sudo attempts for protected files |
| `--purge` | Also remove config and log files |
| `--dry-run` | Print actions without deleting anything |
| `-h, --help` | Show help |

---

## What Gets Removed

| Item | Path |
|------|------|
| LaunchAgent plist | `~/Library/LaunchAgents/com.illwill.cardbot.plist` |
| Binary | `~/bin/cardbot` (or other install path) |
| Config | `~/Library/Application Support/cardbot/` |
| Logs | `~/.cardbot/` |

---

## After Uninstall

- CardBot will no longer start at login.
- No background daemon will be running.
- Config and logs are preserved unless `--purge` was used.
- To reinstall, see [INSTALL.md](INSTALL.md).
