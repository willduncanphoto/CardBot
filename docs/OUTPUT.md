# Output Format

## Startup

```
[2026-03-11 21:15:32] Starting CardBot 0.1.5...
[2026-03-11 21:15:32] Copy location is set to ~/Pictures/CardBot
[2026-03-11 21:15:32] Scanning for memory cards...
```

## First Run (No Config)

```
Welcome to CardBot!

Where should CardBot copy your work?

[macOS: native folder picker opens]

Destination: /Users/user/Pictures/CardBot

[2026-03-11 21:15:32] Starting CardBot 0.1.5...
[2026-03-11 21:15:32] Copy location is set to ~/Pictures/CardBot
[2026-03-11 21:15:32] Scanning for memory cards...
```

## Card Detected (New)

```
[2026-03-11 21:15:32] Scanning for memory cards...card found.
[2026-03-11 21:15:33] Scanning /Volumes/NIKON Z 9  ... 3051 files ✓
[2026-03-11 21:15:33] Scan completed in 0 seconds

  Status:   New
  Path:     /Volumes/NIKON Z 9  
  Storage:  96.4 GB / 476.9 GB (20%)
  Camera:   Nikon Z 9
  Starred:  1
  Content:  2026-02-27      12.9 GB    418   NEF
            2026-02-26      28.4 MB      1   NEF
            ...

  Total:    3048 photos, 0 videos, 96.0 GB
────────────────────────────────────────
[a] Copy All  [e] Eject  [c] Cancel  >
```

## Card Detected (Previously Copied)

```
  Status:   Copied on 2026-03-11 21:31
  Path:     /Volumes/NIKON Z 9  
  Storage:  96.4 GB / 476.9 GB (20%)
  Camera:   Nikon Z 9
  ...
────────────────────────────────────────
[a] Copy All  [e] Eject  [c] Cancel  >
```

## Copy Progress

```
[a] Copy All  [e] Eject  [c] Cancel  > a

[2026-03-11 21:15:35] Copying all files to ~/Pictures/CardBot
[2026-03-11 21:15:40] Copying... 1247/3051 files  48.2 GB/96.4 GB (50%)
...
[2026-03-11 21:22:18] Copy complete ✓
[2026-03-11 21:22:18] 3051 files, 96.0 GB copied in 8m32s (188.4 MB/s)

[a] Copy All  [e] Eject  [c] Cancel  >
```

## Destination Structure

Files are grouped by date, preserving original folder structure:

```
~/Pictures/CardBot/
├── 2026-02-26/
│   └── 100NIKON/
│       └── DSC_0001.NEF
├── 2026-02-27/
│   └── 100NIKON/
│       ├── DSC_0002.NEF
│       ├── DSC_0003.NEF
│       └── DSC_0004.JPG
└── 2026-03-08/
    ├── 100NIKON/
    │   └── DSC_0100.NEF
    └── 101NIKON/
        └── DSC_0200.MOV
```

## Eject

```
[a] Copy All  [e] Eject  [c] Cancel  > e
Ejecting NIKON Z 9  ...

[2026-03-11 21:20:15] Card ejected: /Volumes/NIKON Z 9  

[2026-03-11 21:20:18] Scanning for memory cards...
```

## Card Removal (Unexpected)

```
[2026-03-11 21:20:15] Card removed: /Volumes/NIKON Z 9  

[2026-03-11 21:20:18] Scanning for memory cards...
```

## Queue

When multiple cards are inserted:

```
[2026-03-11 21:15:33] Nikon detected (1 card in queue)
```

Queue is processed in insertion order. The queue count appears when additional cards are waiting.

## Commands

| Key | Action |
|-----|--------|
| `a` + Enter | Copy all files to destination |
| `e` + Enter | Eject the card |
| `c` + Enter | Cancel / dismiss card |

### Hidden Commands

| Key | Action |
|-----|--------|
| `i` + Enter | Show card hardware info (device, model, serial, firmware) |
| `t` + Enter | Run 256MB speed test (sequential write + read) |

## Content Layout

Fixed-width columns for consistent visual scanning:

```
  Content:  2026-03-08      12.9 GB    418   NEF
            2026-03-07      28.4 MB      1   NEF, JPG
```

| Column | Width | Alignment | Description |
|--------|-------|-----------|-------------|
| Date | 10 chars | Left | `YYYY-MM-DD` |
| Size | 10 chars | Right | `NNN.N GB` or `NNN.N MB` |
| Count | variable | Right | File count, right-aligned to widest |
| Extensions | variable | Left | Sorted alphabetically, uppercase |

## Dry-Run Mode

```bash
./cardbot --dry-run
```

Shows `(dry-run)` next to destination. Copy commands report what would happen without writing files.
