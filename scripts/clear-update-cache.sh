#!/usr/bin/env bash
set -euo pipefail

CFG="$HOME/.config/cardbot/config.json"

if [[ ! -f "$CFG" ]]; then
  echo "Config not found: $CFG"
  exit 1
fi

python3 - <<'PY'
import json, pathlib, sys
p = pathlib.Path.home()/".config/cardbot/config.json"
try:
    d = json.loads(p.read_text())
except Exception as e:
    print(f"Failed to parse {p}: {e}")
    sys.exit(1)

d.setdefault("update", {})["last_check"] = ""
p.write_text(json.dumps(d, indent=2) + "\n")
print(f"Cleared update cache: {p}")
print("update.last_check = ''")
PY
