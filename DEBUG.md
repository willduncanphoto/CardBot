# Debugging CardBot Daemon + Auto-Launch

This guide documents all daemon debug controls, CLI helpers, and expected debug output.

## Quick Start

```bash
# Check current daemon debug setting
cardbot daemon-debug status

# Enable verbose daemon/launcher debug logging
cardbot daemon-debug on

# Disable verbose daemon/launcher debug logging
cardbot daemon-debug off
```

After changing debug mode, restart a running daemon process so the new setting is applied.

---

## Debug-Related CLI Commands

### `cardbot daemon-debug [status|on|off]`

Controls `daemon.debug` in config.

- `status` (default): show current state
- `on`: set `daemon.debug = true`
- `off`: set `daemon.debug = false`

Examples:

```bash
cardbot daemon-debug
cardbot daemon-debug status
cardbot daemon-debug on
cardbot daemon-debug off
```

### `cardbot daemon-status`

Shows daemon configuration and runtime environment in human-readable format.

### `cardbot daemon-status --json`

JSON output includes daemon debug state under:

```json
"daemon": {
  "debug": false
}
```

---

## Daemon Config Debug Option

Config file path:

- macOS: `~/Library/Application Support/cardbot/config.json`
- Linux: `~/.config/cardbot/config.json`

Daemon section:

```json
"daemon": {
  "enabled": false,
  "start_at_login": false,
  "terminal_app": "Terminal",
  "launch_args": [],
  "debug": false
}
```

When `debug` is `true`, daemon mode prints verbose launch diagnostics.

---

## Expected Debug Output

### Interactive startup (`cardbot`)

When debug is enabled, startup prints:

- `Daemon debug: enabled`

### Daemon startup (`cardbot --daemon`)

When debug is enabled, startup prints:

- `Daemon debug logging: enabled`

### Daemon launch flow logs

Verbose debug lines include (examples):

- daemon startup config summary
- card insert callback mount path
- single-instance guard block reason
- launcher branch selection (Default/Terminal/Ghostty/custom)
- exact executed command arguments
- generated `.command` script path (system-default terminal branch)

---

## Typical Debug Workflow

1. Enable debug:
   ```bash
   cardbot daemon-debug on
   ```
2. Restart daemon process (or re-run `cardbot --daemon`).
3. Insert card and reproduce issue.
4. Inspect terminal output and log file (`advanced.log_file`) for `Debug:` lines.
5. Disable debug when done:
   ```bash
   cardbot daemon-debug off
   ```
