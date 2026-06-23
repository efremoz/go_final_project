#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIST="$ROOT/dist/GoScheduler"
ARCH="${1:-amd64}"

if [[ "$ARCH" != "amd64" && "$ARCH" != "arm64" ]]; then
	echo "Использование: $0 [amd64|arm64]" >&2
	exit 1
fi

echo "Создание папки Windows-приложения ($ARCH)..."
rm -rf "$DIST"
mkdir -p "$DIST"

echo "Сборка scheduler.exe..."
cd "$ROOT"
CGO_ENABLED=0 GOOS=windows GOARCH="$ARCH" go build -o "$DIST/scheduler.exe" .

cp -R "$ROOT/web" "$DIST/web"
cp "$ROOT/scripts/windows/launcher.ps1" "$DIST/"
cp "$ROOT/scripts/windows/stop-scheduler.ps1" "$DIST/"
cp "$ROOT/scripts/windows/Планировщик задач.bat" "$DIST/"
cp "$ROOT/scripts/windows/stop-scheduler.bat" "$DIST/"

echo ""
echo "Готово: $DIST"
echo "Скопируйте папку GoScheduler на Windows и запустите «Планировщик задач.bat»"
