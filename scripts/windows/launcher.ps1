$ErrorActionPreference = "Stop"

$AppDir = $PSScriptRoot
$SupportDir = Join-Path $env:APPDATA "GoScheduler"
$Port = if ($env:TODO_PORT) { [int]$env:TODO_PORT } else { 7540 }
$Url = "http://localhost:$Port"
$PidFile = Join-Path $SupportDir "scheduler.pid"
$LogFile = Join-Path $SupportDir "scheduler.log"
$LockFile = Join-Path $SupportDir "scheduler.lock"
$ServerExe = Join-Path $SupportDir "scheduler.exe"
$WebDir = Join-Path $SupportDir "web"
$WebSource = Join-Path $AppDir "web"

New-Item -ItemType Directory -Force -Path (Join-Path $SupportDir "data") | Out-Null

if (-not (Test-Path $WebDir)) {
    cmd /c mklink /J "`"$WebDir`"" "`"$WebSource`"" 2>$null | Out-Null
    if (-not (Test-Path $WebDir)) {
        Copy-Item -Recurse -Force $WebSource $WebDir
    }
}

Copy-Item -Force (Join-Path $AppDir "scheduler.exe") $ServerExe

function Test-PortListening([int]$ListenPort) {
    $found = Get-NetTCPConnection -LocalPort $ListenPort -State Listen -ErrorAction SilentlyContinue
    return $null -ne $found
}

function Test-ServerRunning {
    if (Test-Path $PidFile) {
        $pid = [int](Get-Content $PidFile -ErrorAction SilentlyContinue)
        $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
        if ($proc -and (Test-PortListening $Port)) {
            return $true
        }
        Remove-Item $PidFile -Force -ErrorAction SilentlyContinue
    }
    return Test-PortListening $Port
}

function Open-Browser {
    Start-Process $Url
}

function Wait-ForServer {
    for ($i = 0; $i -lt 50; $i++) {
        try {
            $response = Invoke-WebRequest -Uri $Url -UseBasicParsing -TimeoutSec 1
            if ($response.StatusCode -eq 200) {
                return $true
            }
        } catch {}
        Start-Sleep -Milliseconds 200
    }
    return $false
}

if (Test-ServerRunning) {
    Open-Browser
    exit 0
}

$lockStream = $null
try {
    $lockStream = [System.IO.File]::Open(
        $LockFile,
        [System.IO.FileMode]::OpenOrCreate,
        [System.IO.FileAccess]::ReadWrite,
        [System.IO.FileShare]::None
    )
} catch {
    for ($i = 0; $i -lt 25; $i++) {
        if (Test-ServerRunning) {
            Open-Browser
            exit 0
        }
        Start-Sleep -Milliseconds 200
    }
    Write-Error "Не удалось запустить сервер: другой экземпляр занят."
    exit 1
}

try {
    if (Test-ServerRunning) {
        Open-Browser
        exit 0
    }

    $env:TODO_PORT = "$Port"
    $proc = Start-Process `
        -FilePath $ServerExe `
        -WorkingDirectory $SupportDir `
        -WindowStyle Hidden `
        -RedirectStandardOutput $LogFile `
        -RedirectStandardError $LogFile `
        -PassThru

    Set-Content -Path $PidFile -Value $proc.Id

    if (Wait-ForServer) {
        Open-Browser
        exit 0
    }

    Write-Error "Сервер не запустился за 10 секунд. См. $LogFile"
    exit 1
} finally {
    if ($null -ne $lockStream) {
        $lockStream.Close()
    }
}
