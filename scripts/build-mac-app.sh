#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
APP_DIR="$ROOT/dist/yandex-task-to-do-go-final.app"
BINARY="$APP_DIR/Contents/Resources/scheduler"

echo "Создание .app..."
rm -rf "$APP_DIR"
mkdir -p "$APP_DIR/Contents/MacOS" "$APP_DIR/Contents/Resources"

echo "Сборка Go-бинарника..."
cd "$ROOT"
CGO_ENABLED=0 go build -o "$BINARY" .

cp -R "$ROOT/web" "$APP_DIR/Contents/Resources/web"
cp "$ROOT/scripts/launcher.sh" "$APP_DIR/Contents/MacOS/launcher"
chmod +x "$APP_DIR/Contents/MacOS/launcher" "$BINARY"

cat >"$APP_DIR/Contents/Info.plist" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>launcher</string>
	<key>CFBundleIdentifier</key>
	<string>com.goscheduler.app</string>
	<key>CFBundleName</key>
	<string>Планировщик задач</string>
	<key>CFBundleDisplayName</key>
	<string>Планировщик задач</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSMinimumSystemVersion</key>
	<string>11.0</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>LSMultipleInstancesProhibited</key>
	<false/>
</dict>
</plist>
EOF

echo ""
echo "Готово: $APP_DIR"
echo "Запуск: open \"$APP_DIR\""
