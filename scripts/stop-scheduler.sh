#!/bin/bash
set -euo pipefail

SUPPORT_DIR="$HOME/Library/Application Support/GoScheduler"
PID_FILE="$SUPPORT_DIR/scheduler.pid"
PORT="${TODO_PORT:-7540}"

if [[ -f "$PID_FILE" ]]; then
	pid="$(cat "$PID_FILE")"
	if kill -0 "$pid" 2>/dev/null; then
		kill "$pid"
		echo "Сервер остановлен (PID $pid)."
	else
		echo "Процесс $pid не найден."
	fi
	rm -f "$PID_FILE"
else
	pids="$(lsof -i ":${PORT}" -sTCP:LISTEN -t 2>/dev/null || true)"
	if [[ -n "$pids" ]]; then
		kill $pids
		echo "Сервер на порту ${PORT} остановлен."
	else
		echo "Сервер не запущен."
	fi
fi
