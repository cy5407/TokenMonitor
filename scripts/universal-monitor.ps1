#!/usr/bin/env pwsh

<#
.SYNOPSIS
    通用 Token 監控服務 - 支援任何 IDE

.DESCRIPTION
    這個服務可以監控任何 IDE 或編輯器的文件變化，自動計算 Token 使用量。
    無論使用 VS Code、Visual Studio、IntelliJ IDEA、Sublime Text 還是其他任何編輯器都能正常工作。

.PARAMETER Action
    要執行的動作: start, stop, status, install, uninstall

.PARAMETER Background
    是否在後台運行

.EXAMPLE
    .\universal-monitor.ps1 start
    開始監控服務

.EXAMPLE
    .\universal-monitor.ps1 status
    查看監控狀態
#>

param(
    [Parameter(Position = 0)]
    [ValidateSet("start", "stop", "status", "install", "uninstall", "test")]
    [string]$Action = "start",
    
    [switch]$Background
)

# 服務配置
$ServiceName = "UniversalTokenMonitor"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$MonitorScript = Join-Path $ScriptDir "universal-token-monitor.js"
$PidFile = Join-Path $ScriptDir "monitor.pid"
$LogFile = Join-Path $ScriptDir "monitor.log"

function Write-StatusMessage {
    param([string]$Message, [string]$Type = "Info")
    
    $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    $color = switch ($Type) {
        "Success" { "Green" }
        "Warning" { "Yellow" }
        "Error" { "Red" }
        default { "Cyan" }
    }
    
    Write-Host "[$timestamp] $Message" -ForegroundColor $color
}

function Test-NodeExists {
    try {
        $null = Get-Command node -ErrorAction Stop
        return $true
    } catch {
        return $false
    }
}

function Test-MonitorRunning {
    if (Test-Path $PidFile) {
        $pid = Get-Content $PidFile -ErrorAction SilentlyContinue
        if ($pid) {
            try {
                $process = Get-Process -Id $pid -ErrorAction Stop
                return $true
            } catch {
                Remove-Item $PidFile -Force -ErrorAction SilentlyContinue
                return $false
            }
        }
    }
    return $false
}

function Start-Monitor {
    if (Test-MonitorRunning) {
        Write-StatusMessage "監控服務已在運行中" "Warning"
        return
    }

    if (-not (Test-NodeExists)) {
        Write-StatusMessage "錯誤: 找不到 Node.js，請先安裝 Node.js" "Error"
        return
    }

    if (-not (Test-Path $MonitorScript)) {
        Write-StatusMessage "錯誤: 找不到監控腳本 $MonitorScript" "Error"
        return
    }

    Write-StatusMessage "🚀 啟動通用 Token 監控服務..." "Info"

    try {
        if ($Background) {
            # 後台運行
            $process = Start-Process -FilePath "node" -ArgumentList $MonitorScript -WindowStyle Hidden -PassThru
            $process.Id | Out-File -FilePath $PidFile -Encoding UTF8
            Write-StatusMessage "✅ 監控服務已在後台啟動 (PID: $($process.Id))" "Success"
        } else {
            # 前台運行
            Write-StatusMessage "監控服務運行中... (按 Ctrl+C 停止)" "Info"
            & node $MonitorScript
        }
    } catch {
        Write-StatusMessage "❌ 啟動失敗: $($_.Exception.Message)" "Error"
    }
}

function Stop-Monitor {
    if (-not (Test-MonitorRunning)) {
        Write-StatusMessage "監控服務未運行" "Warning"
        return
    }

    try {
        $pid = Get-Content $PidFile
        Stop-Process -Id $pid -Force
        Remove-Item $PidFile -Force
        Write-StatusMessage "🛑 監控服務已停止" "Success"
    } catch {
        Write-StatusMessage "❌ 停止失敗: $($_.Exception.Message)" "Error"
    }
}

function Show-Status {
    Write-Host "`n📊 ===== 通用 Token 監控狀態 =====" -ForegroundColor Cyan
    
    # 服務狀態
    if (Test-MonitorRunning) {
        $pid = Get-Content $PidFile
        Write-Host "🔄 服務狀態: " -NoNewline -ForegroundColor Gray
        Write-Host "運行中 (PID: $pid)" -ForegroundColor Green
    } else {
        Write-Host "🔄 服務狀態: " -NoNewline -ForegroundColor Gray
        Write-Host "已停止" -ForegroundColor Red
    }

    # Node.js 狀態
    Write-Host "🟢 Node.js: " -NoNewline -ForegroundColor Gray
    if (Test-NodeExists) {
        $nodeVersion = & node --version
        Write-Host "已安裝 ($nodeVersion)" -ForegroundColor Green
    } else {
        Write-Host "未安裝" -ForegroundColor Red
    }

    # 腳本文件狀態
    Write-Host "📄 監控腳本: " -NoNewline -ForegroundColor Gray
    if (Test-Path $MonitorScript) {
        Write-Host "存在" -ForegroundColor Green
    } else {
        Write-Host "缺失" -ForegroundColor Red
    }

    # 日誌文件統計
    $logPath = Join-Path $ScriptDir "data\kiro-usage.log"
    Write-Host "📝 日誌文件: " -NoNewline -ForegroundColor Gray
    if (Test-Path $logPath) {
        $content = Get-Content $logPath -ErrorAction SilentlyContinue
        $lineCount = ($content | Measure-Object).Count
        Write-Host "$logPath ($lineCount 筆記錄)" -ForegroundColor Green
    } else {
        Write-Host "尚未建立" -ForegroundColor Yellow
    }

    # 監控的文件類型
    Write-Host "🔍 監控檔案: " -NoNewline -ForegroundColor Gray
    Write-Host ".md, .txt, .js, .ts, .py, .java, .cpp, .html, .css, .json..." -ForegroundColor Cyan

    # 支援的 IDE
    Write-Host "🛠️  支援的 IDE: " -NoNewline -ForegroundColor Gray
    Write-Host "VS Code, Visual Studio, IntelliJ IDEA, Sublime Text, Notepad++, Vim, 等等" -ForegroundColor Cyan

    Write-Host "=====================================`n" -ForegroundColor Cyan
}

function Install-Dependencies {
    Write-StatusMessage "📦 檢查並安裝相依套件..." "Info"
    
    # 檢查 package.json
    $packageJson = Join-Path $ScriptDir "package.json"
    if (-not (Test-Path $packageJson)) {
        Write-StatusMessage "建立 package.json..." "Info"
        $packageContent = @{
            name = "universal-token-monitor"
            version = "1.0.0"
            description = "Universal Token Monitor for any IDE"
            main = "universal-token-monitor.js"
            dependencies = @{
                chokidar = "^3.5.3"
            }
            scripts = @{
                start = "node universal-token-monitor.js"
                monitor = "powershell -File universal-monitor.ps1 start"
            }
        } | ConvertTo-Json -Depth 3

        $packageContent | Out-File -FilePath $packageJson -Encoding UTF8
    }

    # 安裝相依套件
    try {
        Write-StatusMessage "安裝 Node.js 相依套件..." "Info"
        & npm install chokidar
        Write-StatusMessage "✅ 相依套件安裝完成" "Success"
    } catch {
        Write-StatusMessage "❌ 安裝失敗: $($_.Exception.Message)" "Error"
    }
}

function Test-Monitor {
    Write-StatusMessage "🧪 測試通用監控系統..." "Info"
    
    # 建立測試檔案
    $testFile = Join-Path $ScriptDir "test-monitor.md"
    $testContent = @"
# 監控系統測試

這是一個測試檔案，用來驗證通用 Token 監控系統是否正常工作。

## 測試內容

- 文件創建監控
- Token 計算功能
- 日誌記錄功能

測試時間: $(Get-Date)
"@

    try {
        # 先啟動監控（如果還沒運行）
        if (-not (Test-MonitorRunning)) {
            Write-StatusMessage "啟動監控系統進行測試..." "Info"
            Start-Process -FilePath "node" -ArgumentList $MonitorScript -WindowStyle Hidden
            Start-Sleep -Seconds 3
        }

        # 創建測試檔案
        Write-StatusMessage "創建測試檔案..." "Info"
        $testContent | Out-File -FilePath $testFile -Encoding UTF8
        
        Start-Sleep -Seconds 3
        
        # 修改測試檔案
        Write-StatusMessage "修改測試檔案..." "Info"
        Add-Content -Path $testFile -Value "`n## 修改測試`n`n這是一個修改測試。"
        
        Start-Sleep -Seconds 3
        
        # 檢查日誌
        $logPath = Join-Path $ScriptDir "data\kiro-usage.log"
        if (Test-Path $logPath) {
            $recentLogs = Get-Content $logPath | Select-Object -Last 5
            Write-StatusMessage "最近的監控記錄:" "Info"
            $recentLogs | ForEach-Object {
                try {
                    $record = $_ | ConvertFrom-Json
                    Write-Host "  📄 $($record.timestamp): $($record.file_name) ($($record.tokens) tokens)" -ForegroundColor Gray
                } catch {
                    Write-Host "  📄 $_" -ForegroundColor Gray
                }
            }
        }
        
        # 清理測試檔案
        Remove-Item $testFile -Force -ErrorAction SilentlyContinue
        Write-StatusMessage "✅ 測試完成" "Success"
        
    } catch {
        Write-StatusMessage "❌ 測試失敗: $($_.Exception.Message)" "Error"
    }
}

# 主要邏輯
switch ($Action) {
    "start" {
        Start-Monitor
    }
    "stop" {
        Stop-Monitor
    }
    "status" {
        Show-Status
    }
    "install" {
        Install-Dependencies
    }
    "test" {
        Test-Monitor
    }
    default {
        Write-Host @"
通用 Token 監控服務

用法: .\universal-monitor.ps1 [action] [options]

動作:
  start     啟動監控服務
  stop      停止監控服務  
  status    查看服務狀態
  install   安裝相依套件
  test      測試監控功能

選項:
  -Background    在後台運行服務
  -Verbose       顯示詳細信息

範例:
  .\universal-monitor.ps1 start -Background
  .\universal-monitor.ps1 status
  .\universal-monitor.ps1 test

"@ -ForegroundColor Cyan
    }
}
