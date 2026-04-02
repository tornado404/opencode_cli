#Requires -Version 5.1
<#
.SYNOPSIS
    OpenCode Server Windows 启动脚本

.DESCRIPTION
    从配置文件读取密码并启动 OpenCode Server
#>

param(
    [string]$ConfigFile
)

$ErrorActionPreference = "Stop"

# 默认配置文件路径
if ([string]::IsNullOrEmpty($ConfigFile)) {
    $ConfigFile = Join-Path $env:APPDATA "oho\config.json"
}

# 读取配置文件中的密码
$Password = ""
if (Test-Path $ConfigFile) {
    try {
        $Config = Get-Content -Path $ConfigFile -Raw | ConvertFrom-Json
        if ($Config.password) {
            $Password = $Config.password
        }
    }
    catch {
        Write-Warning "读取配置文件失败: $($_.Exception.Message)"
    }
}

if ([string]::IsNullOrEmpty($Password)) {
    Write-Host "❌ 配置文件未找到或密码为空。请先运行 install.ps1 安装并配置。" -ForegroundColor Red
    Write-Host ""
    Write-Host "运行以下命令安装:"
    Write-Host "  irm https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.ps1 | iex" -ForegroundColor Cyan
    exit 1
}

# 检查 opencode 是否在 PATH 中
$OhoCmd = Get-Command opencode -ErrorAction SilentlyContinue
if (-not $OhoCmd) {
    # 尝试从默认安装目录查找
    $DefaultPath = Join-Path $env:LOCALAPPDATA "Programs\oho\oho.exe"
    if (Test-Path $DefaultPath) {
        $env:PATH = "$((Split-Path $DefaultPath));$env:PATH"
    }
    else {
        Write-Host "❌ opencode 命令未找到。请确保 oho 已安装并添加到 PATH。" -ForegroundColor Red
        Write-Host ""
        Write-Host "运行以下命令安装:"
        Write-Host "  irm https://raw.githubusercontent.com/tornado404/opencode_cli/master/oho/install.ps1 | iex" -ForegroundColor Cyan
        exit 1
    }
}

# 设置环境变量
$env:OPENCODE_SERVER_PASSWORD = $Password

# 检查服务是否已运行
$RunningProcess = Get-Process -Name "opencode" -ErrorAction SilentlyContinue | Where-Object { $_.MainWindowTitle -like "*opencode*" }
if ($RunningProcess) {
    Write-Host "⚠️  OpenCode Server 已在运行 (PID: $($RunningProcess.Id))" -ForegroundColor Yellow
    Write-Host "   访问地址: http://localhost:4096" -ForegroundColor Cyan
    exit 0
}

Write-Host "🟢 正在启动 OpenCode Server..." -ForegroundColor Green

# 启动服务
$LogFile = Join-Path $env:TEMP "opencode.log"
$Process = Start-Process -FilePath "opencode" -ArgumentList "web","--hostname","0.0.0.0","--port","4096","--mdns","--mdns-domain","opencode.local" -PassThru -RedirectStandardOutput $LogFile -RedirectStandardError $LogFile

Start-Sleep -Seconds 2

# 验证服务是否启动成功
if (-not $Process.HasExited) {
    Write-Host "✅ OpenCode Server 已启动 (PID: $($Process.Id))" -ForegroundColor Green
    Write-Host "   访问地址: http://localhost:4096" -ForegroundColor Cyan
    Write-Host "   日志文件: $LogFile" -ForegroundColor Gray
}
else {
    Write-Host "❌ 服务启动失败" -ForegroundColor Red
    if (Test-Path $LogFile) {
        Write-Host "日志输出:" -ForegroundColor Yellow
        Get-Content $LogFile | Select-Object -First 10
    }
    exit 1
}