# 0.5.1 Sleep/Wake QA Run — 2026-03-19

## Scope

Manual validation for follow-up item **051-002**:
- Sleep/wake duplicate launch behavior on real hardware.

## Environment

- Terminal app: Ghostty
- Daemon mode: `cardbot --daemon`
- Test log: `/tmp/cardbot-qa-051-sleepwake-manual.log`

## Procedure

1. Started daemon capture.
2. Inserted card (baseline).
3. Put machine to sleep and woke.
4. Removed/reinserted card multiple times.
5. Stopped daemon and reviewed log.

## Observed log evidence

```text
[2026-03-18T18:31:48] Card detected: NIKON Z 9   (/Volumes/NIKON Z 9  )
[2026-03-18T18:31:48] Launched Ghostty for /Volumes/NIKON Z 9
[2026-03-18T18:33:28] Card removed: /Volumes/NIKON Z 9
[2026-03-18T18:33:38] Card detected: NIKON Z 9   (/Volumes/NIKON Z 9  )
[2026-03-18T18:33:38] Launched Ghostty for /Volumes/NIKON Z 9
[2026-03-18T18:36:46] Card removed: /Volumes/NIKON Z 9
[2026-03-18T18:36:57] Card detected: NIKON Z 9   (/Volumes/NIKON Z 9  )
[2026-03-18T18:36:57] Launched Ghostty for /Volumes/NIKON Z 9
```

## Result

Status: ✅ **PASS**

Interpretation:
- No launch storms observed.
- One launch per reinsert after removal.
- Sleep/wake cycle did not produce duplicate rapid launch behavior.
