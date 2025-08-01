#!/usr/bin/env pwsh

<#
.SYNOPSIS
    改進的 AI Token 監控系統管理腳本

.DESCRIPTION
    管理改進的 AI Token 監控系統，包括啟動、停止、狀態查詢和測試功能

.PARAMETER Action
    要執行的動作: start, stop, status, test, diagnose

.EXAMPLE
    .\monitor_manager.ps1 start
    啟動改進的監控系統

.EXAMPLE
    .\monitor_manager.ps1 test
    測試監控系統功能
#>

param(
    [Parameter(Position = 0)]
    [ValidateSet("start", "stop", "status", "test", "diagnose")]
    [string]$Action = "status"
)

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

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

function Start-ImprovedMonitor {
    Write-StatusMessage "🚀 啟動改進的 AI 監控系統..." "Info"
    
    # 檢查 Node.js 是否可用
    try {
        $null = Get-Command node -ErrorAction Stop
    }
    catch {
        Write-StatusMessage "❌ 找不到 Node.js，請先安裝 Node.js" "Error"
        return
    }
    
    # 檢查監控腳本是否存在
    $monitorScript = Join-Path $ScriptDir "improved_ai_monitor.js"
    if (-not (Test-Path $monitorScript)) {
        Write-StatusMessage "❌ 找不到監控腳本: $monitorScript" "Error"
        return
    }
    
    # 啟動監控系統
    try {
        Write-StatusMessage "📊 啟動改進的監控系統..." "Info"
        Write-StatusMessage "💡 按 Ctrl+C 停止監控" "Warning"
        
        Push-Location $ProjectRoot
        & node $monitorScript
        
    }
    catch {
        Write-StatusMessage "❌ 啟動失敗: $($_.Exception.Message)" "Error"
    }
    finally {
        Pop-Location
    }
}

function Stop-ImprovedMonitor {
    Write-StatusMessage "🛑 停止改進的 AI 監控系統..." "Info"
    
    try {
        # 查找並停止 Node.js 監控程序
        $nodeProcesses = Get-Process -Name "node" -ErrorAction SilentlyContinue | 
                        Where-Object { $_.CommandLine -like "*improved_ai_monitor*" }
        
        if ($nodeProcesses) {
            foreach ($process in $nodeProcesses) {
                Stop-Process -Id $process.Id -Force
                Write-StatusMessage "✅ 已停止監控程序 (PID: $($process.Id))" "Success"
            }
        }
        else {
            Write-StatusMessage "⚠️ 沒有找到運行中的監控程序" "Warning"
        }
        
    }
    catch {
        Write-StatusMessage "❌ 停止失敗: $($_.Exception.Message)" "Error"
    }
}

function Show-MonitorStatus {
    Write-StatusMessage "📊 檢查改進的監控系統狀態..." "Info"
    
    # 檢查 Node.js 程序
    $nodeProcesses = Get-Process -Name "node" -ErrorAction SilentlyContinue
    $monitorProcesses = $nodeProcesses | Where-Object { $_.CommandLine -like "*improved_ai_monitor*" }
    
    Write-Host "`n📊 ===== 改進的 AI 監控系統狀態 =====" -ForegroundColor Cyan
    
    if ($monitorProcesses) {
        Write-Host "🟢 監控狀態: " -NoNewline -ForegroundColor Gray
        Write-Host "運行中" -ForegroundColor Green
        
        foreach ($process in $monitorProcesses) {
            Write-Host "   PID: $($process.Id)" -ForegroundColor Gray
        }
    }
    else {
        Write-Host "🔴 監控狀態: " -NoNewline -ForegroundColor Gray
        Write-Host "已停止" -ForegroundColor Red
    }
    
    # 檢查日誌檔案
    $logFile = Join-Path $ProjectRoot "data/kiro-usage.log"
    Write-Host "📄 日誌檔案: " -NoNewline -ForegroundColor Gray
    if (Test-Path $logFile) {
        $logStats = Get-Item $logFile
        $content = Get-Content $logFile -Raw -ErrorAction SilentlyContinue
        $lineCount = ($content -split "`n").Count
        
        Write-Host "存在 ($lineCount 筆記錄, $(($logStats.Length / 1KB).ToString('F2')) KB)" -ForegroundColor Green
        Write-Host "   最後修改: $($logStats.LastWriteTime.ToString('yyyy-MM-dd HH:mm:ss'))" -ForegroundColor Gray
    }
    else {
        Write-Host "不存在" -ForegroundColor Yellow
    }
    
    # 檢查 Node.js 版本
    try {
        $nodeVersion = & node --version
        Write-Host "🟢 Node.js: " -NoNewline -ForegroundColor Gray
        Write-Host "已安裝 ($nodeVersion)" -ForegroundColor Green
    }
    catch {
        Write-Host "🔴 Node.js: " -NoNewline -ForegroundColor Gray
        Write-Host "未安裝" -ForegroundColor Red
    }
    
    Write-Host "=====================================`n" -ForegroundColor Cyan
}

function Test-ImprovedMonitor {
    Write-StatusMessage "🧪 測試改進的監控系統..." "Info"
    
    $testScript = Join-Path $ScriptDir ".." "Tests" "monitor_diagnosis" "test_improved_monitor.js"
    
    if (-not (Test-Path $testScript)) {
        Write-StatusMessage "❌ 找不到測試腳本: $testScript" "Error"
        return
    }
    
    try {
        Push-Location $ProjectRoot
        Write-StatusMessage "🔄 執行監控系統測試..." "Info"
        & node $testScript
        
    }
    catch {
        Write-StatusMessage "❌ 測試失敗: $($_.Exception.Message)" "Error"
    }
    finally {
        Pop-Location
    }
}

function Start-Diagnosis {
    Write-StatusMessage "🔍 診斷現有監控系統..." "Info"
    
    $diagnosisScript = Join-Path $ScriptDir ".." "Tests" "monitor_diagnosis" "diagnose_current_system.js"
    
    if (-not (Test-Path $diagnosisScript)) {
        Write-StatusMessage "❌ 找不到診斷腳本: $diagnosisScript" "Error"
        return
    }
    
    try {
        Push-Location $ProjectRoot
        & node $diagnosisScript
        
    }
    catch {
        Write-StatusMessage "❌ 診斷失敗: $($_.Exception.Message)" "Error"
    }
    finally {
        Pop-Location
    }
}

# 主要邏輯
switch ($Action) {
    "start" {
        Start-ImprovedMonitor
    }
    "stop" {
        Stop-ImprovedMonitor
    }
    "status" {
        Show-MonitorStatus
    }
    "test" {
        Test-ImprovedMonitor
    }
    "diagnose" {
        Start-Diagnosis
    }
    default {
        Write-Host @"
改進的 AI Token 監控系統管理工具

用法: .\monitor_manager.ps1 [action]

動作:
  start      啟動改進的監控系統
  stop       停止監控系統
  status     查看監控狀態 (預設)
  test       測試監控功能
  diagnose   診斷現有系統

範例:
  .\monitor_manager.ps1 start
  .\monitor_manager.ps1 status
  .\monitor_manager.ps1 test

"@ -ForegroundColor Cyan
    }
}