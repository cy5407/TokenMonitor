#!/usr/bin/env pwsh

<#
.SYNOPSIS
    簡化版 Token 獵人 - 尋找 Kiro IDE 對話數據

.DESCRIPTION
    多角度搜尋 Kiro IDE 的對話和 Token 數據
#>

Write-Host "🔍 啟動 Token 獵人..." -ForegroundColor Cyan

# 1. 檢查 Kiro 程序
Write-Host "`n🔬 檢查 Kiro 程序..." -ForegroundColor Yellow
$kiroProcesses = Get-Process -Name "*Kiro*" -ErrorAction SilentlyContinue

if ($kiroProcesses) {
    foreach ($process in $kiroProcesses) {
        Write-Host "  ✅ 找到 Kiro 程序: PID $($process.Id), 記憶體: $(($process.WorkingSet64/1MB).ToString('F2')) MB" -ForegroundColor Green
        
        # 檢查程序路徑
        try {
            $processPath = $process.Path
            Write-Host "    📁 程序路徑: $processPath" -ForegroundColor Gray
            
            # 檢查程序目錄中的檔案
            $processDir = Split-Path $processPath -Parent
            $configFiles = Get-ChildItem -Path $processDir -Include "*.json", "*.db", "*.log" -Recurse -ErrorAction SilentlyContinue
            foreach ($file in $configFiles) {
                Write-Host "    📄 相關檔案: $($file.Name) ($(($file.Length/1KB).ToString('F2')) KB)" -ForegroundColor Cyan
            }
        }
        catch {
            Write-Host "    ⚠️ 無法存取程序路徑" -ForegroundColor Yellow
        }
    }
} else {
    Write-Host "  ❌ 沒有找到 Kiro 程序" -ForegroundColor Red
}

# 2. 搜尋使用者配置目錄
Write-Host "`n📂 搜尋使用者配置..." -ForegroundColor Yellow

$configPaths = @(
    "$env:APPDATA\Kiro",
    "$env:LOCALAPPDATA\Kiro", 
    "$env:USERPROFILE\.kiro",
    "$env:APPDATA\kiro",
    "$env:LOCALAPPDATA\kiro"
)

foreach ($path in $configPaths) {
    if (Test-Path $path) {
        Write-Host "  ✅ 找到配置目錄: $path" -ForegroundColor Green
        
        # 搜尋相關檔案
        $files = Get-ChildItem -Path $path -Recurse -Include "*.log", "*.db", "*.json", "*.sqlite*" -ErrorAction SilentlyContinue
        foreach ($file in $files) {
            Write-Host "    📄 檔案: $($file.FullName) ($(($file.Length/1KB).ToString('F2')) KB)" -ForegroundColor Cyan
            
            # 檢查檔案內容
            if ($file.Extension -eq ".log" -and $file.Length -lt 1MB) {
                try {
                    $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue
                    if ($content -and ($content.Contains("token") -or $content.Contains("chat") -or $content.Contains("claude"))) {
                        Write-Host "      🎯 包含相關內容！" -ForegroundColor Magenta
                        
                        # 嘗試提取 token 數量
                        $tokenMatches = [regex]::Matches($content, "token[s]?[:\s]*(\d+)", [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
                        if ($tokenMatches.Count -gt 0) {
                            $totalTokens = 0
                            foreach ($match in $tokenMatches) {
                                $totalTokens += [int]$match.Groups[1].Value
                            }
                            Write-Host "      💎 發現 Token 數據: $totalTokens tokens" -ForegroundColor Green
                        }
                    }
                }
                catch {
                    Write-Host "      ⚠️ 無法讀取檔案內容" -ForegroundColor Yellow
                }
            }
        }
    }
}

# 3. 搜尋臨時檔案
Write-Host "`n🗂️ 搜尋臨時檔案..." -ForegroundColor Yellow

$tempPaths = @($env:TEMP, "$env:LOCALAPPDATA\Temp")
foreach ($tempPath in $tempPaths) {
    $kiroTempFiles = Get-ChildItem -Path $tempPath -Filter "*kiro*" -Recurse -ErrorAction SilentlyContinue
    foreach ($file in $kiroTempFiles) {
        Write-Host "  📄 臨時檔案: $($file.FullName)" -ForegroundColor Cyan
    }
}

# 4. 檢查網路連線
Write-Host "`n🌐 檢查網路連線..." -ForegroundColor Yellow

if ($kiroProcesses) {
    foreach ($process in $kiroProcesses) {
        try {
            $connections = Get-NetTCPConnection | Where-Object { $_.OwningProcess -eq $process.Id }
            foreach ($conn in $connections) {
                Write-Host "  🔗 連線: $($conn.LocalAddress):$($conn.LocalPort) -> $($conn.RemoteAddress):$($conn.RemotePort)" -ForegroundColor Cyan
                
                # 檢查是否為 AI 服務連線
                if ($conn.RemotePort -eq 443 -or $conn.RemoteAddress -like "*api*") {
                    Write-Host "    🤖 可能的 AI 服務連線！" -ForegroundColor Magenta
                }
            }
        }
        catch {
            Write-Host "  ⚠️ 無法取得程序 $($process.Id) 的網路連線" -ForegroundColor Yellow
        }
    }
}

# 5. 搜尋瀏覽器快取 (如果 Kiro 使用 Electron)
Write-Host "`n🌍 搜尋瀏覽器快取..." -ForegroundColor Yellow

$electronCachePaths = @(
    "$env:APPDATA\Kiro\User Data\Default\Cache",
    "$env:LOCALAPPDATA\Kiro\User Data\Default\Cache",
    "$env:APPDATA\kiro\User Data\Default\Cache"
)

foreach ($cachePath in $electronCachePaths) {
    if (Test-Path $cachePath) {
        Write-Host "  ✅ 找到快取目錄: $cachePath" -ForegroundColor Green
        
        $cacheFiles = Get-ChildItem -Path $cachePath -ErrorAction SilentlyContinue | Select-Object -First 10
        foreach ($file in $cacheFiles) {
            Write-Host "    📄 快取檔案: $($file.Name)" -ForegroundColor Cyan
        }
    }
}

# 6. 檢查 Windows 事件日誌
Write-Host "`n📋 檢查系統日誌..." -ForegroundColor Yellow

try {
    $recentEvents = Get-WinEvent -FilterHashtable @{LogName='Application'; StartTime=(Get-Date).AddHours(-1)} -MaxEvents 50 -ErrorAction SilentlyContinue |
                   Where-Object { $_.ProcessId -in $kiroProcesses.Id -or $_.Message -like "*kiro*" }
    
    if ($recentEvents) {
        Write-Host "  ✅ 找到 $($recentEvents.Count) 個相關事件" -ForegroundColor Green
        foreach ($event in $recentEvents | Select-Object -First 5) {
            Write-Host "    📝 事件: $($event.TimeCreated) - $($event.LevelDisplayName)" -ForegroundColor Cyan
        }
    } else {
        Write-Host "  ℹ️ 沒有找到相關的系統事件" -ForegroundColor Gray
    }
}
catch {
    Write-Host "  ⚠️ 無法存取系統事件日誌" -ForegroundColor Yellow
}

# 7. 總結建議
Write-Host "`n💡 獵取總結與建議:" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan

Write-Host "🎯 可能的 Token 數據來源:" -ForegroundColor Yellow
Write-Host "  1. Kiro IDE 的本地資料庫檔案 (.db, .sqlite)" -ForegroundColor White
Write-Host "  2. 應用程式日誌檔案 (.log)" -ForegroundColor White  
Write-Host "  3. 網路請求快取" -ForegroundColor White
Write-Host "  4. 程序記憶體中的即時數據" -ForegroundColor White

Write-Host "`n🔧 進階獵取策略:" -ForegroundColor Yellow
Write-Host "  1. 使用 Process Monitor 監控檔案存取" -ForegroundColor White
Write-Host "  2. 使用 Wireshark 攔截網路流量" -ForegroundColor White
Write-Host "  3. 使用 API Hook 攔截函數調用" -ForegroundColor White
Write-Host "  4. 分析 Electron 應用的 DevTools" -ForegroundColor White

Write-Host "`n⚡ 即時監控建議:" -ForegroundColor Yellow
Write-Host "  1. 在對話時同時運行檔案監控" -ForegroundColor White
Write-Host "  2. 監控網路流量變化" -ForegroundColor White
Write-Host "  3. 觀察程序記憶體使用變化" -ForegroundColor White

Write-Host "`n🎉 獵取完成！" -ForegroundColor Green