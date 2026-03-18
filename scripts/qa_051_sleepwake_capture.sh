#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

TS="$(date +%Y%m%d-%H%M%S)"
OUT_DIR="${1:-/tmp/cardbot-qa-051-$TS}"
mkdir -p "$OUT_DIR"
LOG_FILE="$OUT_DIR/daemon.log"

echo "[qa-051] output dir: $OUT_DIR"
echo "[qa-051] building cardbot"
go build -o cardbot .

echo "[qa-051] starting daemon"
./cardbot --daemon >"$LOG_FILE" 2>&1 &
DAEMON_PID=$!

TAIL_PID=""
cleanup() {
  set +e
  if [[ -n "$TAIL_PID" ]]; then
    kill "$TAIL_PID" >/dev/null 2>&1 || true
  fi
  if kill -0 "$DAEMON_PID" >/dev/null 2>&1; then
    kill -INT "$DAEMON_PID" >/dev/null 2>&1 || true
    wait "$DAEMON_PID" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT INT TERM

sleep 1

echo "[qa-051] daemon pid: $DAEMON_PID"
echo "[qa-051] log file: $LOG_FILE"
echo
echo "Perform manual sleep/wake validation now:"
echo "  1) Insert a card once and verify one app launch"
echo "  2) Put Mac to sleep, then wake"
echo "  3) Remove/reinsert card quickly"
echo "  4) Watch for duplicate suppression instead of launch storms"
echo
echo "Live daemon log (Ctrl+C to stop):"

# tail exits when killed by cleanup
(tail -f "$LOG_FILE") &
TAIL_PID=$!

# keep script alive until user interrupts
wait "$TAIL_PID" || true

echo
echo "[qa-051] summary"
SUPPRESSIONS=$(grep -c "Suppressing duplicate card event" "$LOG_FILE" || true)
SKIPS=$(grep -c "already running in another process" "$LOG_FILE" || true)
LAUNCH_FAILS=$(grep -c "Launch failed:" "$LOG_FILE" || true)

echo "  duplicate suppressions: $SUPPRESSIONS"
echo "  single-instance skips: $SKIPS"
echo "  launch failures: $LAUNCH_FAILS"
echo
echo "[qa-051] done"
