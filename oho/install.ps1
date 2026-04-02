#Requires -Version 5.1
<#
.SYNOPSIS
    oho Windows 安装脚本

.DESCRIPTION
    从 GitHub Releases 下载预编译的 oho 二进制文件，或从源码编译
    支持 Windows x64 和 ARM64 架构

.PARAMETER Version
    指定安装版本 (默认: latest)

.PARAMETER InstallDir
    指定安装目录 (默认: $env:LOCALAPPDATA\Programs\oho)

.PARAMETER FromSource
    强制从源码编译

.PARAMETER SkipConfig
    跳过配置文件创建

.EXAMPLE
    .\install.ps1

.EXAMPLE
    .\install.ps1 -Version v1.1.0 -InstallDir "D:\Tools\oho"

.EXAMPLE
    .\install.ps1 -FromSource
#>

param(
    [string]$Version,
    [string]$InstallDir,
    [switch]$FromSource,
    [switch]$SkipConfig
)

# ============================================
# 配置变量
# ============================================
$REPO_OWNER = "tornado404"
$REPO_NAME = "opencode_cli"
$BINARY_NAME = "oho"

# 默认安装目录
if ([string]::IsNullOrEmpty($InstallDir)) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\oho"
}

# 配置文件目录
$ConfigDir = Join-Path $env:APPDATA "oho"
$ConfigFile = Join-Path $ConfigDir "config.json"

# 颜色定义
function Get-ColorCodes {
    return @{
        "Red" = "`e[31m"
        "Green" = "`e[32m"
        "Yellow" = "`e[33m"
        "Reset" = "`e[0m"
    }
}
$Colors = Get-ColorCodes

# ============================================
# 辅助函数
# ============================================

function Write-Success {
    param([string]$Message)
    Write-Host "${Colors.Green}[成功]${Colors.Reset} $Message" -NoNewline
    Write-Host ""
}

function Write-ErrorMsg {
    param([string]$Message)
    Write-Host "${Colors.Red}[错误]${Colors.Reset} $Message" -NoNewline
    Write-Host ""
}

function Write-Warning {
    param([string]$Message)
    Write-Host "${Colors.Yellow}[警告]${Colors.Reset} $Message" -NoNewline
    Write-Host ""
}

function Write-Info {
    param([string]$Message)
    Write-Host "[信息] $Message"
}

# ============================================
# 检测操作系统和架构
# ============================================
function Get-SystemArch {
    Write-Info "检测系统架构..."

    # 检测架构
    $Arch = $env:PROCESSOR_ARCHITECTURE
    if ([string]::IsNullOrEmpty($Arch)) {
        $Arch = (Get-CimInstance Win32_OperatingSystem).OSArchitecture
    }

    # 标准化架构名称
    switch -Regex ($Arch) {
        "AMD64|x64|64" {
            $Arch = "amd64"
        }
        "ARM64|aarch64|ARM" {
            $Arch = "arm64"
        }
        default {
            Write-ErrorMsg "不支持的架构: $Arch"
            exit 1
        }
    }

    Write-Info "架构: $Arch"
    return $Arch
}

# ============================================
# 获取最新版本
# ============================================
function Get-LatestVersion {
    Write-Info "获取最新版本..."

    try {
        $Response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" -TimeoutSec 10 -ErrorAction Stop
        $Version = $Response.tag_name

        if ([string]::IsNullOrEmpty($Version)) {
            throw "版本为空"
        }

        Write-Info "最新版本: $Version"
        return $Version
    }
    catch {
        Write-Warning "无法获取最新版本: $($_.Exception.Message)"
        Write-Info "将尝试从源码编译"
        return "dev"
    }
}

# ============================================
# 下载二进制文件
# ============================================
function Download-Binary {
    param(
        [string]$Version,
        [string]$Arch
    )

    Write-Info "准备下载 oho..."

    # 构建下载 URL
    $DownloadUrl = "https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$Version/$BINARY_NAME-windows-$Arch.zip"

    Write-Info "下载 URL: $DownloadUrl"

    # 创建临时目录
    $TempDir = Join-Path $env:TEMP "oho_install_$(Get-Random)"
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

    $ZipFile = Join-Path $TempDir "oho.zip"

    try {
        # 下载文件
        Write-Info "下载中..."

        # 使用 WebClient 支持进度
        $WebClient = New-Object System.Net.WebClient
        $WebClient.DownloadFile($DownloadUrl, $ZipFile)

        if (-not (Test-Path $ZipFile)) {
            throw "下载失败"
        }

        Write-Success "下载完成"

        # 解压
        Write-Info "解压中..."
        Expand-Archive -Path $ZipFile -DestinationPath $TempDir -Force

        # 查找解压后的 oho.exe
        $ExtractedExe = Get-ChildItem -Path $TempDir -Filter "*.exe" -Recurse | Select-Object -First 1

        if ($null -eq $ExtractedExe) {
            throw "解压后未找到 oho.exe"
        }

        # 创建安装目录
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        # 移动到安装目录
        $TargetPath = Join-Path $InstallDir "oho.exe"
        Copy-Item -Path $ExtractedExe.FullName -Destination $TargetPath -Force

        Write-Success "安装完成: $TargetPath"
        return $TargetPath
    }
    catch {
        Write-ErrorMsg "下载失败: $($_.Exception.Message)"
        return $null
    }
    finally {
        # 清理临时目录
        if (Test-Path $TempDir) {
            Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# ============================================
# 从源码编译
# ============================================
function Build-FromSource {
    Write-Info "从源码编译..."

    # 检查 Go 是否安装
    $GoCmd = Get-Command go -ErrorAction SilentlyContinue

    if ($null -eq $GoCmd) {
        Write-ErrorMsg "Go 未安装，无法从源码编译"
        Write-Info "请访问 https://golang.org/dl/ 安装 Go 1.21+"
        return $null
    }

    # 检查 Go 版本
    $GoVersion = (go version) -replace 'go', ''
    $GoVersionMatch = $GoVersion -match '(\d+)\.(\d+)'

    if ($GoVersionMatch) {
        $Major = [int]$Matches[1]
        $Minor = [int]$Matches[2]

        if ($Major -lt 1 -or ($Major -eq 1 -and $Minor -lt 21)) {
            Write-ErrorMsg "需要 Go 1.21+, 当前版本: $GoVersion"
            return $null
        }
    }

    Write-Info "Go 版本: $GoVersion"

    # 创建临时目录用于克隆
    $TempDir = Join-Path $env:TEMP "oho_build_$(Get-Random)"
    New-Item -ItemType Directory -Path $TempDir -Force | Out-Null

    try {
        Write-Info "克隆仓库..."
        $RepoUrl = "https://github.com/$REPO_OWNER/$REPO_NAME.git"

        # 使用 git clone
        $GitCmd = Get-Command git -ErrorAction SilentlyContinue
        if ($null -eq $GitCmd) {
            Write-ErrorMsg "git 未安装，无法克隆仓库"
            return $null
        }

        Push-Location $TempDir
        git clone --depth 1 $RepoUrl

        if ($LASTEXITCODE -ne 0) {
            throw "git clone 失败"
        }

        $OhoPath = Join-Path $TempDir "$REPO_NAME\oho"
        if (-not (Test-Path $OhoPath)) {
            $OhoPath = Join-Path $TempDir $REPO_NAME
        }

        Push-Location $OhoPath
        Write-Info "编译中..."

        # 创建安装目录
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        $TargetPath = Join-Path $InstallDir "oho.exe"

        # 编译 (带版本信息)
        $Version = Get-LatestVersion
        $Commit = "unknown"
        $Date = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")

        go build -ldflags "-s -w -X main.Version=$Version -X main.Commit=$Commit -X main.Date=$Date" -o $TargetPath ./cmd

        if ($LASTEXITCODE -ne 0) {
            throw "go build 失败"
        }

        Pop-Location
        Pop-Location

        Write-Success "编译完成: $TargetPath"
        return $TargetPath
    }
    catch {
        Write-ErrorMsg "编译失败: $($_.Exception.Message)"
        return $null
    }
    finally {
        # 清理临时目录
        if (Test-Path $TempDir) {
            Remove-Item -Path $TempDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# ============================================
# 创建配置文件 (带随机密码)
# ============================================
function New-ConfigFile {
    if ($SkipConfig) {
        Write-Info "跳过配置创建"
        return
    }

    Write-Info "创建配置文件..."

    # 创建配置目录
    if (-not (Test-Path $ConfigDir)) {
        New-Item -ItemType Directory -Path $ConfigDir -Force | Out-Null
    }

    # 检查配置文件是否已存在
    if (Test-Path $ConfigFile) {
        Write-Info "配置文件已存在: $ConfigFile"
        
        # 检查是否已有密码
        try {
            $ExistingConfig = Get-Content -Path $ConfigFile -Raw | ConvertFrom-Json
            if ($ExistingConfig.password) {
                Write-Info "配置文件已包含密码"
                return
            }
        }
        catch {
            Write-Warning "无法读取现有配置文件，将重新创建"
        }
    }

    # 生成 8 位随机密码
    $Password = -join ((65..90) + (97..122) + (48..57) | Get-Random -Count 8 | ForEach-Object { [char]$_ })
    Write-Info "生成随机密码: ********"

    # 创建默认配置 (包含随机密码)
    $ConfigContent = @{
        host = "127.0.0.1"
        port = 4096
        username = "opencode"
        password = $Password
        json = $false
    } | ConvertTo-Json -Depth 2

    $ConfigContent | Out-File -FilePath $ConfigFile -Encoding UTF8

    Write-Success "配置文件已创建: $ConfigFile"
    Write-Host "   密码: $Password (请妥善保管)" -ForegroundColor Yellow
}

# ============================================
# 检测并启动 OpenCode Server
# ============================================
function Start-OpenCodeServer {
    Write-Info "检测 OpenCode Server 服务状态..."

    # 检查 opencode 命令是否可用
    $OhoCmd = Get-Command opencode -ErrorAction SilentlyContinue
    if (-not $OhoCmd) {
        # 尝试从默认安装目录查找
        $DefaultPath = Join-Path $env:LOCALAPPDATA "Programs\oho\oho.exe"
        if (Test-Path $DefaultPath) {
            $env:PATH = "$((Split-Path $DefaultPath));$env:PATH"
            $OhoCmd = Get-Command opencode -ErrorAction SilentlyContinue
        }
        
        if (-not $OhoCmd) {
            Write-Warning "opencode 命令不可用，跳过服务启动"
            return
        }
    }

    # 检查服务是否已运行
    try {
        $Response = Invoke-RestMethod -Uri "http://127.0.0.1:4096/health" -TimeoutSec 3 -ErrorAction SilentlyContinue
        if ($Response) {
            Write-Success "OpenCode Server 已在运行"
            return
        }
    }
    catch {
        # 服务未运行，继续启动
    }

    Write-Info "OpenCode Server 未运行，准备启动..."

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
            Write-Warning "无法读取配置文件密码"
        }
    }

    if ([string]::IsNullOrEmpty($Password)) {
        Write-Warning "配置文件无密码，无法启动服务"
        Write-Info "请手动设置密码后运行: .\run.ps1"
        return
    }

    # 启动服务
    $env:OPENCODE_SERVER_PASSWORD = $Password
    
    Write-Info "启动 OpenCode Server..."
    $LogFile = Join-Path $env:TEMP "opencode.log"
    
    # 异步启动
    Start-Process -FilePath "opencode" -ArgumentList "web","--hostname","0.0.0.0","--port","4096","--mdns","--mdns-domain","opencode.local" -PassThru -WindowStyle Hidden
    
    Start-Sleep -Seconds 3

    # 验证启动
    try {
        $Response = Invoke-RestMethod -Uri "http://127.0.0.1:4096/health" -TimeoutSec 5 -ErrorAction SilentlyContinue
        if ($Response) {
            Write-Success "OpenCode Server 已启动"
            Write-Host "   访问地址: http://localhost:4096" -ForegroundColor Cyan
            Write-Host "   配置文件: $ConfigFile" -ForegroundColor Gray
        }
    }
    catch {
        Write-Warning "服务可能需要更长时间启动"
        Write-Info "请稍后手动检查: oho global health"
    }
}

# ============================================
# 添加到 PATH
# ============================================
function Add-ToPath {
    Write-Info "添加安装目录到 PATH..."

    # 获取当前用户 PATH
    $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")

    # 检查是否已存在
    if ($CurrentPath -like "*$InstallDir*") {
        Write-Info "安装目录已在 PATH 中"
        return
    }

    # 添加到 PATH
    $NewPath = "$InstallDir;$CurrentPath"
    [Environment]::SetEnvironmentVariable("PATH", $NewPath, "User")

    Write-Success "PATH 已更新"
    Write-Info "注意: 此更改将在新终端中生效"
    Write-Info "当前终端可以使用以下命令立即生效:"
    Write-Info "  `$env:PATH = `\"$NewPath`\""
}

# ============================================
# 验证安装
# ============================================
function Test-Install {
    param([string]$ExePath)

    Write-Info "验证安装..."

    if (-not (Test-Path $ExePath)) {
        Write-ErrorMsg "可执行文件不存在: $ExePath"
        return $false
    }

    try {
        # 尝试获取版本
        $Output = & $ExePath --version 2>&1

        if ($LASTEXITCODE -eq 0) {
            Write-Success "验证成功"
            Write-Host "版本信息: $Output"
            return $true
        }
        else {
            Write-Warning "无法获取版本信息，但文件已安装"
            return $true
        }
    }
    catch {
        Write-Warning "验证时发生错误: $($_.Exception.Message)"
        return $true
    }
}

# ============================================
# 主函数
# ============================================
function Main {
    Write-Host ""
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "         oho Windows 安装脚本" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""

    # 1. 检测架构
    $Arch = Get-SystemArch
    Write-Host ""

    # 2. 获取版本
    if ([string]::IsNullOrEmpty($Version)) {
        $Version = Get-LatestVersion
    }
    Write-Info "安装版本: $Version"
    Write-Info "安装目录: $InstallDir"
    Write-Host ""

    # 3. 安装二进制文件
    $ExePath = $null

    if ($FromSource) {
        # 从源码编译
        Write-Info "模式: 从源码编译"
        $ExePath = Build-FromSource
    }
    else {
        # 尝试下载
        Write-Info "模式: 下载预编译二进制"
        $ExePath = Download-Binary -Version $Version -Arch $Arch

        if ($null -eq $ExePath) {
            Write-Warning "下载失败，尝试从源码编译..."

            $GoCmd = Get-Command go -ErrorAction SilentlyContinue
            if ($null -ne $GoCmd) {
                $ExePath = Build-FromSource
            }
            else {
                Write-ErrorMsg "无法下载且 Go 未安装"
                exit 1
            }
        }
    }

    if ($null -eq $ExePath) {
        Write-ErrorMsg "安装失败"
        exit 1
    }

    Write-Host ""

    # 4. 创建配置文件
    New-ConfigFile
    Write-Host ""

    # 5. 添加到 PATH
    Add-ToPath
    Write-Host ""

    # 6. 验证安装
    $Verified = Test-Install -ExePath $ExePath
    Write-Host ""

    # 7. 检测并启动 OpenCode Server
    Start-OpenCodeServer
    Write-Host ""

    # 8. 显示使用说明
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "安装完成！" -ForegroundColor Green
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "使用说明:"
    Write-Host "  1. 重新打开一个终端窗口，或运行以下命令刷新 PATH:"
    Write-Host "     `$env:PATH = [Environment]::GetEnvironmentVariable(\"PATH\",\"User\")"
    Write-Host ""
    Write-Host "  2. 运行以下命令测试安装:"
    Write-Host "     oho --version"
    Write-Host ""
    Write-Host "  3. 配置 OpenCode Server 连接信息 (环境变量或配置文件):"
    Write-Host "     `$env:OPENCODE_SERVER_HOST = \"127.0.0.1\""
    Write-Host "     `$env:OPENCODE_SERVER_PORT = \"4096\""
    Write-Host "     `$env:OPENCODE_SERVER_PASSWORD = \"your-password\""
    Write-Host ""
    Write-Host "  4. 查看帮助:"
    Write-Host "     oho --help"
    Write-Host ""
}

# 执行主函数
Main