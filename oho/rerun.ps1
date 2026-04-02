#Requires -Version 5.1
<#
.SYNOPSIS
    OpenCode Server Windows 重启脚本

.DESCRIPTION
    停止旧的 OpenCode Server 进程并重新启动
#>

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   OpenCode Server 重启脚本" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "🔴 正在停止 OpenCode Server 服务..." -ForegroundColor Yellow

# 查找并停止 opencode web 进程
$Processes = Get-Process -Name "opencode" -ErrorAction SilentlyContinue
$Stopped = $false

if ($Processes) {
    foreach ($Process in $Processes) {
        try {
            Stop-Process -Id $Process.Id -Force -ErrorAction Stop
            Write-Host "   已终止进程 PID: $($Process.Id)" -ForegroundColor Gray
            $Stopped = $true
        }
        catch {
            Write-Warning "无法终止进程 $($Process.Id): $($_.Exception.Message)"
        }
    }
    if ($Stopped) {
        Start-Sleep -Seconds 2
    }
}

Write-Host "✅ OpenCode Server 服务已停止" -ForegroundColor Green

# 获取脚本所在目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

if ($ScriptDir) {
    $RunScript = Join-Path $ScriptDir "run.ps1"
}
else {
    # 如果无法获取目录，尝试从当前目录查找
    $RunScript = Join-Path $PWD "run.ps1"
}

Write-Host ""
Write-Host "🟢 正在启动 OpenCode Server 服务..." -ForegroundColor Yellow

if (Test-Path $RunScript) {
    & $RunScript
}
else {
    Write-Host "❌ 未找到 run.ps1 脚本" -ForegroundColor Red
    Write-Host "   请确保 run.ps1 与 rerun.ps1 在同一目录" -ForegroundColor Gray
    exit 1
}