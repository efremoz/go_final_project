#!/bin/bash
set -euo pipefail

APP_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RESOURCES="$APP_ROOT/Resources"
SUPPORT_DIR="$HOME/Library/Application Support/GoScheduler"
PORT="${TODO_PORT:-7540}"
URL="http://localhost:${PORT}"
PID_FILE="$SUPPORT_DIR/scheduler.pid"
LOG_FILE="$SUPPORT_DIR/scheduler.log"

mkdir -p "$SUPPORT_DIR/data"
ln -sfn "$RESOURCES/web" "$SUPPORT_DIR/web"
ln -sfn "$RESOURCES/scheduler" "$SUPPORT_DIR/scheduler"

cleanup() {
	if [[ -f "$PID_FILE" ]]; then
		kill "$(cat "$PID_FILE")" 2>/dev/null || true
		rm -f "$PID_FILE"
	fi
}
trap cleanup EXIT INT TERM

if lsof -i ":${PORT}" -sTCP:LISTEN -t >/dev/null 2>&1; then
	open "$URL"
	exit 0
fi

cd "$SUPPORT_DIR"
export TODO_PORT="$PORT"
./scheduler >>"$LOG_FILE" 2>&1 &
echo $! >"$PID_FILE"

for _ in $(seq 1 50); do
	if curl -sf "$URL" >/dev/null 2>&1; then
		open "$URL"
		wait "$(cat "$PID_FILE")"
		exit 0
	fi
	sleep 0.2
done

echo "Сервер не запустился за 10 секунд. См. $LOG_FILE" >&2
exit 1
