@echo off
chcp 65001 >nul
setlocal

set "ROOT=%~dp0.."
set "DIST=%ROOT%\dist\GoScheduler"
set "ARCH=%~1"
if "%ARCH%"=="" set "ARCH=amd64"

if /I not "%ARCH%"=="amd64" if /I not "%ARCH%"=="arm64" (
    echo Использование: %~nx0 [amd64^|arm64]
    exit /b 1
)

echo Создание папки Windows-приложения (%ARCH%)...
if exist "%DIST%" rmdir /s /q "%DIST%"
mkdir "%DIST%"

echo Сборка scheduler.exe...
pushd "%ROOT%"
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=%ARCH%
go build -o "%DIST%\scheduler.exe" .
if errorlevel 1 exit /b 1
popd

xcopy /E /I /Y "%ROOT%\web" "%DIST%\web" >nul
copy /Y "%ROOT%\scripts\windows\launcher.ps1" "%DIST%\" >nul
copy /Y "%ROOT%\scripts\windows\stop-scheduler.ps1" "%DIST%\" >nul
copy /Y "%ROOT%\scripts\windows\Планировщик задач.bat" "%DIST%\" >nul
copy /Y "%ROOT%\scripts\windows\stop-scheduler.bat" "%DIST%\" >nul

echo.
echo Готово: %DIST%
echo Запуск: "%DIST%\Планировщик задач.bat"
