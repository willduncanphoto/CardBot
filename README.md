# CardBot

A CLI tool for camera memory cards.

![CardBot screenshot](screenshot.png)

## What CardBot does

- Detects camera memory cards on macOS (and eventually linux)
- Card content overview (file count, type, dates, equiptment data)
- Fast and safe copy operations. Time is money. Safety is life.
- Dated folder structures for local copy location
- Copy modes: all, selects (starred), photos only, videos only, etc.

## Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| macOS | Supported | Primary platform |
| Linux | It might work | Untested |
| Windows | Ugh | Someday, Maybe |

## Installation

One-liner installer:

```bash
curl -fsSL https://raw.githubusercontent.com/willduncanphoto/CardBot/main/install.sh | sh
```

## Usage

Start interactive mode:

```bash
cardbot
```

CardBot will automatically run the setup the first time you run it.

To run the setup again (dest, naming prefs, startup and auto-detect in prefs):

```bash
cardbot --setup
```

## CLI flags

| Flag | Description |
|------|-------------|
| `--dest <path>` | Override destination path for this run |
| `--dry-run` | Analyze only; do not copy |
| `--daemon` | Run headless background watcher |
| `--setup` | Re-run setup prompts |
| `--reset` | Clear saved config |
| `--version` | Print version |

## Interactive commands

| Key | Action |
|-----|--------|
| `a` | Copy all |
| `s` | Copy selects (starred/picked) |
| `p` | Copy photos |
| `v` | Copy videos |
| `e` | Eject card |
| `x` | Exit current card |
| `i` | Show card hardware info |
| `t` | Run speed test |
| `\` | Cancel active copy |
| `?` | Help |

## Roadmap

| Version | Focus | Status |
|---------|-------|--------|
| **0.5.2** | Launcher diagnostics + uninstall workflow | Current |
| **0.6.0** | General improvements | Next |
| **0.7.0** | Homebrew support | Planned |
| **0.8.0** | Copyright check | Planned |

## License

MIT — see [LICENSE](LICENSE).
