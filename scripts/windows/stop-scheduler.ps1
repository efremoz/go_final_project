$SupportDir = Join-Path $env:APPDATA "GoScheduler"
$PidFile = Join-Path $SupportDir "scheduler.pid"
$Port = if ($env:TODO_PORT) { [int]$env:TODO_PORT } else { 7540 }

if (Test-Path $PidFile) {
    $pid = [int](Get-Content $PidFile)
    $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
    if ($proc) {
        Stop-Process -Id $pid -Force
        Write-Host "Сервер остановлен (PID $pid)."
    } else {
        Write-Host "Процесс $pid не найден."
    }
    Remove-Item $PidFile -Force
} else {
    $conns = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
    if ($conns) {
        $conns | ForEach-Object {
            Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
        }
        Write-Host "Сервер на порту $Port остановлен."
    } else {
        Write-Host "Сервер не запущен."
    }
}
