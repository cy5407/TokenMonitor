#!/usr/bin/env pwsh

<#
.SYNOPSIS
    Token Usage CLI - 模仿 ccusage 的專業 Token 使用量分析工具

.DESCRIPTION
    這個腳本提供類似 ccusage 的命令行介面，用於分析 Kiro IDE 的 Token 使用情況。

.PARAMETER Command
    要執行的命令: daily, weekly, monthly, summary

.PARAMETER Since
    分析的起始日期 (格式: YYYY-MM-DD)

.PARAMETER Model
    過濾特定的模型

.EXAMPLE
    .\tokusage.ps1 daily
    顯示每日 Token 使用報告

.EXAMPLE
    .\tokusage.ps1 daily --since 2025-07-01
    顯示自指定日期以來的每日報告
#>

param(
    [Parameter(Position = 0)]
    [ValidateSet("daily", "weekly", "monthly", "summary", "")]
    [string]$Command = "daily",
    
    [Parameter()]
    [string]$Since,
    
    [Parameter()]
    [string]$Model,
    
    [switch]$Help
)

# 顯示幫助信息
if ($Help) {
    Write-Host @"
Token Usage CLI - Kiro IDE Token 使用量分析工具

用法:
    tokusage [command] [options]

命令:
    daily       顯示每日 Token 使用報告 (預設)
    weekly      顯示每週彙總報告
    monthly     顯示每月彙總報告
    summary     顯示總體摘要

選項:
    --since DATE    從指定日期開始分析 (YYYY-MM-DD)
    --model MODEL   過濾特定模型
    --help          顯示此幫助訊息

範例:
    tokusage daily                    # 每日報告
    tokusage daily --since 2025-07-01 # 從指定日期的每日報告
    tokusage summary                  # 總體摘要
    tokusage weekly                   # 每週報告

"@ -ForegroundColor Cyan
    exit 0
}

# 設定工作目錄到專案根目錄
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir
Push-Location $ProjectRoot

try {
    # 檢查必要檔案
    if (-not (Test-Path "src/js/Professional-Token-Cli.js")) {
        Write-Error "找不到 src/js/Professional-Token-Cli.js 檔案"
        exit 1
    }

    if (-not (Test-Path "data/kiro-usage.log")) {
        Write-Warning "找不到使用記錄檔案: data/kiro-usage.log"
        Write-Host "請確保 Token 監控系統已啟用並產生了使用記錄。" -ForegroundColor Yellow
        exit 1
    }

    # 根據命令執行不同的分析
    switch ($Command) {
        "daily" {
            Write-Host "🔍 執行每日 Token 使用分析..." -ForegroundColor Green
            node "src/js/Professional-Token-Cli.js"
        }
        
        "weekly" {
            Write-Host "📅 執行每週 Token 使用分析..." -ForegroundColor Green
            # 這裡可以擴展為週報邏輯
            node "src/js/Professional-Token-Cli.js"
            Write-Host "`n📊 每週報告功能即將推出..." -ForegroundColor Yellow
        }
        
        "monthly" {
            Write-Host "📆 執行每月 Token 使用分析..." -ForegroundColor Green
            # 這裡可以擴展為月報邏輯
            node "src/js/Professional-Token-Cli.js"
            Write-Host "`n📊 每月報告功能即將推出..." -ForegroundColor Yellow
        }
        
        "summary" {
            Write-Host "📈 執行總體使用摘要分析..." -ForegroundColor Green
            node ".kiro/hooks/manual-token-calc.js"
        }
        
        default {
            Write-Host "🔍 執行預設每日分析..." -ForegroundColor Green
            node "src/js/Professional-Token-Cli.js"
        }
    }

    # 顯示額外信息
    Write-Host "`n💡 提示:" -ForegroundColor Cyan
    Write-Host "   • 使用 'tokusage --help' 查看更多選項" -ForegroundColor Gray
    Write-Host "   • 使用 'tokusage summary' 查看詳細統計" -ForegroundColor Gray
    Write-Host "   • Token 監控系統會自動記錄所有 Kiro 活動" -ForegroundColor Gray

} catch {
    Write-Error "執行失敗: $_"
    exit 1
} finally {
    Pop-Location
}
