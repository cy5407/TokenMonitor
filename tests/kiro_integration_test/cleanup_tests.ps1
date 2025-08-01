#!/usr/bin/env pwsh

<#
.SYNOPSIS
    清理 Kiro IDE 整合測試產生的檔案

.DESCRIPTION
    清理測試過程中產生的臨時檔案和測試記錄
#>

$testDir = Split-Path -Parent $MyInvocation.MyCommand.Path

Write-Host "🧹 清理測試檔案..." -ForegroundColor Cyan

# 清理可能的測試檔案
$testFiles = @(
    "test_*.md",
    "temp_*.log",
    "*.tmp"
)

foreach ($pattern in $testFiles) {
    $files = Get-ChildItem -Path $testDir -Filter $pattern -ErrorAction SilentlyContinue
    foreach ($file in $files) {
        Remove-Item $file.FullName -Force
        Write-Host "  🗑️  已刪除: $($file.Name)" -ForegroundColor Gray
    }
}

Write-Host "✅ 清理完成" -ForegroundColor Green