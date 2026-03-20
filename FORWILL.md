# CardBot Reboot Test Checklist

## Current State

- Daemon setup simplified to use **Terminal.app** via AppleScript (no app-choice prompt).
- Terminal app selection prompt removed from `--setup`.
- Setup always reinstalls LaunchAgent to keep plist aligned with current binary path.
- Target-path launches now populate filesystem usage (storage %).
- Stale `cardbot-daemon.sh` wrapper removed.

---

## Before Reboot

1. Build/install binary:
   ```bash
   go build -o ~/bin/cardbot .
   ```

2. Re-run setup (refreshes LaunchAgent plist):
   ```bash
   cardbot --setup
   ```

3. Confirm config:
   ```bash
   cardbot daemon-status
   ```

---

## Reboot + Test

1. Reboot macOS.
2. Log in normally.
3. Insert Nikon card.
4. Observe Terminal.app window opens with CardBot.

---

## Expected Results

- Terminal.app opens via AppleScript (no ugly `.command ; exit;` header).
- Storage usage shows correct values (not 0 B / 0 B).
- Single card launch, no duplicate windows.
- Card scan runs normally.

---

## If It Still Fails

1. `cardbot daemon-status --recent-launches 20`
2. `pgrep -fl cardbot`
3. Relevant lines from `~/.cardbot/cardbot.log`

---

## Backlog

- Ghostty terminal support (deferred — needs investigation).
- Branch-first workflow for future features.
