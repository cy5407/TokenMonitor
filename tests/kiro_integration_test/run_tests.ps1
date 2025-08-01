#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Kiro IDE 整合測試執行腳本

.DESCRIPTION
    測試 Kiro IDE Token 監控系統是否正常運作
#>

param(
    [switch]$Verbose
)

$testDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent (Split-Path -Parent $testDir)
$logFile = Join-Path $projectRoot "data/kiro-usage.log"

Write-Host "🧪 開始 Kiro IDE 整合測試..." -ForegroundColor Cyan

# 檢查日誌檔案是否存在
if (-not (Test-Path $logFile)) {
    Write-Host "❌ 找不到日誌檔案: $logFile" -ForegroundColor Red
    exit 1
}

# 記錄測試開始前的日誌行數
$initialLogCount = (Get-Content $logFile | Measure-Object).Count
Write-Host "📊 測試前日誌記錄數: $initialLogCount" -ForegroundColor Gray

Write-Host "💬 請在 Kiro IDE 中與 AI 進行對話..." -ForegroundColor Yellow
Write-Host "   (測試將等待 30 秒檢查新記錄)" -ForegroundColor Gray

# 等待用戶進行對話
Start-Sleep -Seconds 30

# 檢查是否有新的日誌記錄
$finalLogCount = (Get-Content $logFile | Measure-Object).Count
$newRecords = $finalLogCount - $initialLogCount

Write-Host "📊 測試後日誌記錄數: $finalLogCount" -ForegroundColor Gray
Write-Host "📈 新增記錄數: $newRecords" -ForegroundColor Gray

if ($newRecords -gt 0) {
    Write-Host "✅ 檢測到新的監控記錄" -ForegroundColor Green
    
    # 顯示最新的記錄
    $recentLogs = Get-Content $logFile | Select-Object -Last $newRecords
    Write-Host "`n📝 最新記錄:" -ForegroundColor Cyan
    
    foreach ($log in $recentLogs) {
        try {
            $record = $log | ConvertFrom-Json
            $eventType = $record.event
            $timestamp = $record.timestamp
            
            if ($eventType -eq "chat_message") {
                Write-Host "  ✅ AI 對話記錄: $timestamp" -ForegroundColor Green
            } elseif ($eventType -eq "tool_execution") {
                Write-Host "  🔧 工具執行記錄: $timestamp" -ForegroundColor Blue
            } else {
                Write-Host "  📄 其他記錄: $eventType - $timestamp" -ForegroundColor Gray
            }
        }
        catch {
            Write-Host "  📄 記錄: $log" -ForegroundColor Gray
        }
    }
} else {
    Write-Host "❌ 沒有檢測到新的監控記錄" -ForegroundColor Red
    Write-Host "   可能的原因:" -ForegroundColor Yellow
    Write-Host "   - Kiro IDE Hook 未啟用" -ForegroundColor Gray
    Write-Host "   - 監控系統未正常運作" -ForegroundColor Gray
    Write-Host "   - 沒有進行 AI 對話" -ForegroundColor Gray
}

Write-Host "`n🏁 測試完成" -ForegroundColor Cyan