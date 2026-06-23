#!/bin/bash
set -euo pipefail

APP_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RESOURCES="$APP_ROOT/Resources"
SUPPORT_DIR="$HOME/Library/Application Support/GoScheduler"
PORT="${TODO_PORT:-7540}"
URL="http://localhost:${PORT}"
PID_FILE="$SUPPORT_DIR/scheduler.pid"
LOG_FILE="$SUPPORT_DIR/scheduler.log"
LOCK_FILE="$SUPPORT_DIR/scheduler.lock"

mkdir -p "$SUPPORT_DIR/data"
ln -sfn "$RESOURCES/web" "$SUPPORT_DIR/web"
ln -sfn "$RESOURCES/scheduler" "$SUPPORT_DIR/scheduler"

is_server_running() {
	if [[ -f "$PID_FILE" ]]; then
		local pid
		pid="$(cat "$PID_FILE")"
		if kill -0 "$pid" 2>/dev/null && lsof -i ":${PORT}" -sTCP:LISTEN -t >/dev/null 2>&1; then
			return 0
		fi
		rm -f "$PID_FILE"
	fi
	lsof -i ":${PORT}" -sTCP:LISTEN -t >/dev/null 2>&1
}

open_browser() {
	open "$URL"
}

if is_server_running; then
	open_browser
	exit 0
fi

exec 9>"$LOCK_FILE"
if ! flock -n 9; then
	for _ in $(seq 1 25); do
		if is_server_running; then
			open_browser
			exit 0
		fi
		sleep 0.2
	done
	echo "Не удалось запустить сервер: другой экземпляр занят." >&2
	exit 1
fi

if is_server_running; then
	open_browser
	exit 0
fi

cd "$SUPPORT_DIR"
export TODO_PORT="$PORT"
nohup ./scheduler >>"$LOG_FILE" 2>&1 &
echo $! >"$PID_FILE"
disown

for _ in $(seq 1 50); do
	if curl -sf "$URL" >/dev/null 2>&1; then
		open_browser
		exit 0
	fi
	sleep 0.2
done

echo "Сервер не запустился за 10 секунд. См. $LOG_FILE" >&2
exit 1
