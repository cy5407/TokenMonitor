#!/usr/bin/env pwsh

<#
.SYNOPSIS
    終極 Token 獵人 - 多維度獲取 Kiro IDE 對話 Token

.DESCRIPTION
    使用多種策略嘗試獲取 Kiro IDE 中的對話數據和 Token 使用量
#>

param(
    [switch]$DeepScan,
    [switch]$NetworkMonitor,
    [switch]$MemoryAnalysis,
    [switch]$FileSystemHunt
)

class UltimateTokenHunter {
    [string]$LogFile
    [array]$KiroProcesses
    [hashtable]$FoundData
    
    UltimateTokenHunter() {
        $this.LogFile = "data/token_hunt_results.log"
        $this.FoundData = @{}
        $this.EnsureLogDirectory()
    }
    
    [void] EnsureLogDirectory() {
        $logDir = Split-Path $this.LogFile -Parent
        if (-not (Test-Path $logDir)) {
            New-Item -ItemType Directory -Path $logDir -Force | Out-Null
        }
    }
    
    [void] LogFindings([string]$Category, [string]$Finding) {
        $timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        $logEntry = "[$timestamp] [$Category] $Finding"
        
        Write-Host $logEntry -ForegroundColor Green
        Add-Content -Path $this.LogFile -Value $logEntry
        
        if (-not $this.FoundData.ContainsKey($Category)) {
            $this.FoundData[$Category] = @()
        }
        $this.FoundData[$Category] += $Finding
    }
    
    [void] StartHunt() {
        Write-Host "🔍 啟動終極 Token 獵人..." -ForegroundColor Cyan
        Write-Host "🎯 目標：找到 Kiro IDE 對話的 Token 數據" -ForegroundColor Yellow
        
        # 1. 基本環境掃描
        $this.ScanEnvironment()
        
        # 2. Kiro 程序分析
        $this.AnalyzeKiroProcesses()
        
        # 3. 檔案系統深度掃描
        if ($FileSystemHunt -or $DeepScan) {
            $this.DeepFileSystemScan()
        }
        
        # 4. 網路監控
        if ($NetworkMonitor -or $DeepScan) {
            $this.StartNetworkMonitoring()
        }
        
        # 5. 記憶體分析
        if ($MemoryAnalysis -or $DeepScan) {
            $this.AnalyzeProcessMemory()
        }
        
        # 6. 生成報告
        $this.GenerateReport()
    }
    
    [void] ScanEnvironment() {
        Write-Host "`n🌍 環境掃描..." -ForegroundColor Cyan
        
        # 檢查 Kiro IDE 安裝位置
        $possiblePaths = @(
            "$env:PROGRAMFILES\Kiro",
            "$env:PROGRAMFILES(X86)\Kiro",
            "$env:LOCALAPPDATA\Programs\Kiro",
            "$env:APPDATA\Kiro"
        )
        
        foreach ($path in $possiblePaths) {
            if (Test-Path $path) {
                $this.LogFindings("Environment", "Found Kiro installation: $path")
                
                # 掃描配置檔案
                $configFiles = Get-ChildItem -Path $path -Recurse -Include "*.json", "*.yaml", "*.config" -ErrorAction SilentlyContinue
                foreach ($file in $configFiles) {
                    $this.LogFindings("Config", "Config file: $($file.FullName)")
                }
            }
        }
        
        # 檢查使用者配置目錄
        $userConfigPaths = @(
            "$env:APPDATA\Kiro",
            "$env:LOCALAPPDATA\Kiro",
            "$env:USERPROFILE\.kiro"
        )
        
        foreach ($path in $userConfigPaths) {
            if (Test-Path $path) {
                $this.LogFindings("UserConfig", "User config directory: $path")
                
                # 尋找日誌檔案
                $logFiles = Get-ChildItem -Path $path -Recurse -Include "*.log", "*.txt" -ErrorAction SilentlyContinue
                foreach ($file in $logFiles) {
                    if ($file.Length -gt 0) {
                        $this.LogFindings("Logs", "Log file: $($file.FullName) ($(($file.Length/1KB).ToString('F2')) KB)")
                    }
                }
            }
        }
    }
    
    [void] AnalyzeKiroProcesses() {
        Write-Host "`n🔬 Kiro 程序分析..." -ForegroundColor Cyan
        
        $this.KiroProcesses = Get-Process -Name "*Kiro*" -ErrorAction SilentlyContinue
        
        if ($this.KiroProcesses.Count -eq 0) {
            $this.LogFindings("Process", "No Kiro processes found")
            return
        }
        
        foreach ($process in $this.KiroProcesses) {
            $this.LogFindings("Process", "Kiro process: PID $($process.Id), Memory: $(($process.WorkingSet64/1MB).ToString('F2')) MB")
            
            # 分析程序模組
            try {
                $modules = $process.Modules
                foreach ($module in $modules) {
                    if ($module.ModuleName -like "*token*" -or $module.ModuleName -like "*api*") {
                        $this.LogFindings("Module", "Interesting module: $($module.ModuleName)")
                    }
                }
            }
            catch {
                $this.LogFindings("Process", "Cannot access modules for PID $($process.Id)")
            }
            
            # 檢查程序的工作目錄
            try {
                $processPath = $process.Path
                if ($processPath) {
                    $workingDir = Split-Path $processPath -Parent
                    $this.LogFindings("Process", "Working directory: $workingDir")
                    
                    # 掃描工作目錄中的相關檔案
                    $relevantFiles = Get-ChildItem -Path $workingDir -Include "*.db", "*.sqlite", "*.json" -ErrorAction SilentlyContinue
                    foreach ($file in $relevantFiles) {
                        $this.LogFindings("ProcessFiles", "Process file: $($file.FullName)")
                    }
                }
            }
            catch {
                $this.LogFindings("Process", "Cannot access path for PID $($process.Id)")
            }
        }
    }
    
    [void] DeepFileSystemScan() {
        Write-Host "`n🗂️ 檔案系統深度掃描..." -ForegroundColor Cyan
        
        # 掃描可能包含對話數據的檔案
        $searchPatterns = @(
            "*.db", "*.sqlite", "*.sqlite3",
            "*chat*.json", "*conversation*.json", "*session*.json",
            "*token*.log", "*usage*.log", "*api*.log",
            "*.cache", "*temp*"
        )
        
        $searchPaths = @(
            $env:APPDATA,
            $env:LOCALAPPDATA,
            $env:TEMP,
            $env:USERPROFILE
        )
        
        foreach ($basePath in $searchPaths) {
            foreach ($pattern in $searchPatterns) {
                try {
                    $files = Get-ChildItem -Path $basePath -Filter $pattern -Recurse -ErrorAction SilentlyContinue | 
                             Where-Object { $_.Name -like "*kiro*" -or $_.Directory.Name -like "*kiro*" }
                    
                    foreach ($file in $files) {
                        if ($file.Length -gt 0) {
                            $this.LogFindings("DeepScan", "Potential data file: $($file.FullName) ($(($file.Length/1KB).ToString('F2')) KB)")
                            
                            # 嘗試分析檔案內容
                            $this.AnalyzeFile($file.FullName)
                        }
                    }
                }
                catch {
                    # 忽略權限錯誤
                }
            }
        }
    }
    
    [void] AnalyzeFile([string]$FilePath) {
        try {
            $extension = [System.IO.Path]::GetExtension($FilePath).ToLower()
            
            switch ($extension) {
                ".json" {
                    $content = Get-Content $FilePath -Raw -ErrorAction SilentlyContinue
                    if ($content -and ($content.Contains("token") -or $content.Contains("chat") -or $content.Contains("message"))) {
                        $this.LogFindings("FileAnalysis", "JSON file contains relevant data: $FilePath")
                        
                        # 嘗試解析 JSON
                        try {
                            $jsonData = $content | ConvertFrom-Json
                            if ($jsonData.tokens -or $jsonData.usage -or $jsonData.messages) {
                                $this.LogFindings("FileAnalysis", "JSON contains token/usage data: $FilePath")
                            }
                        }
                        catch {
                            $this.LogFindings("FileAnalysis", "JSON parse failed: $FilePath")
                        }
                    }
                }
                
                ".log" {
                    $content = Get-Content $FilePath -Raw -ErrorAction SilentlyContinue
                    if ($content -and ($content.Contains("token") -or $content.Contains("claude") -or $content.Contains("chat"))) {
                        $this.LogFindings("FileAnalysis", "Log file contains relevant data: $FilePath")
                        
                        # 計算可能的 token 數量
                        $tokenMatches = [regex]::Matches($content, "token[s]?[:\s]*(\d+)", [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
                        if ($tokenMatches.Count -gt 0) {
                            $totalTokens = 0
                            foreach ($match in $tokenMatches) {
                                $totalTokens += [int]$match.Groups[1].Value
                            }
                            $this.LogFindings("TokenData", "Found $($tokenMatches.Count) token references, total: $totalTokens tokens in $FilePath")
                        }
                    }
                }
                
                { $_ -in @(".db", ".sqlite", ".sqlite3") } {
                    $this.LogFindings("Database", "Database file found: $FilePath")
                    # 這裡可以加入 SQLite 分析邏輯
                }
            }
        }
        catch {
            $this.LogFindings("FileAnalysis", "Failed to analyze: $FilePath - $($_.Exception.Message)")
        }
    }
    
    [void] StartNetworkMonitoring() {
        Write-Host "`n🌐 網路監控..." -ForegroundColor Cyan
        
        # 檢查 Kiro 程序的網路連線
        if ($this.KiroProcesses.Count -gt 0) {
            foreach ($process in $this.KiroProcesses) {
                try {
                    $connections = Get-NetTCPConnection | Where-Object { $_.OwningProcess -eq $process.Id }
                    foreach ($conn in $connections) {
                        $this.LogFindings("Network", "Kiro connection: $($conn.LocalAddress):$($conn.LocalPort) -> $($conn.RemoteAddress):$($conn.RemotePort)")
                        
                        # 檢查是否連接到 AI 服務
                        if ($conn.RemoteAddress -like "*anthropic*" -or $conn.RemoteAddress -like "*openai*" -or $conn.RemotePort -eq 443) {
                            $this.LogFindings("AIConnection", "Potential AI service connection: $($conn.RemoteAddress):$($conn.RemotePort)")
                        }
                    }
                }
                catch {
                    $this.LogFindings("Network", "Cannot get connections for PID $($process.Id)")
                }
            }
        }
        
        # 檢查最近的網路活動
        $this.LogFindings("Network", "Network monitoring started - check Windows Event Log for detailed network activity")
    }
    
    [void] AnalyzeProcessMemory() {
        Write-Host "`n🧠 記憶體分析..." -ForegroundColor Cyan
        
        if ($this.KiroProcesses.Count -eq 0) {
            $this.LogFindings("Memory", "No Kiro processes to analyze")
            return
        }
        
        foreach ($process in $this.KiroProcesses) {
            $this.LogFindings("Memory", "Analyzing memory for PID $($process.Id)")
            
            # 基本記憶體資訊
            $memoryMB = ($process.WorkingSet64 / 1MB)
            $this.LogFindings("Memory", "Process memory usage: $($memoryMB.ToString('F2')) MB")
            
            # 這裡可以加入更進階的記憶體分析
            # 例如使用 Windows API 或第三方工具
            $this.LogFindings("Memory", "Advanced memory analysis requires additional tools")
        }
    }
    
    [void] GenerateReport() {
        Write-Host "`n📊 生成獵取報告..." -ForegroundColor Cyan
        
        $reportPath = "token_hunt_report_$(Get-Date -Format 'yyyyMMdd_HHmmss').md"
        
        $report = @"
# 終極 Token 獵人報告

**執行時間:** $(Get-Date)
**目標:** 獲取 Kiro IDE 對話 Token 數據

## 發現摘要

"@
        
        foreach ($category in $this.FoundData.Keys) {
            $report += "`n### $category`n`n"
            foreach ($finding in $this.FoundData[$category]) {
                $report += "- $finding`n"
            }
        }
        
        $report += @"

## 建議的下一步

1. **檢查發現的檔案** - 分析找到的配置和日誌檔案
2. **監控網路流量** - 使用 Wireshark 或 Fiddler 攔截 API 請求
3. **程序注入** - 開發 DLL 注入工具監控 API 調用
4. **逆向工程** - 分析 Kiro IDE 的執行檔案結構

## 技術限制

- 某些檔案可能需要管理員權限
- 記憶體分析需要專業工具
- 網路加密可能阻止內容分析
- 程序保護機制可能阻止注入

"@
        
        $report | Out-File -FilePath $reportPath -Encoding UTF8
        
        Write-Host "`n✅ 獵取完成！" -ForegroundColor Green
        Write-Host "📄 詳細報告: $reportPath" -ForegroundColor Yellow
        Write-Host "📋 日誌檔案: $($this.LogFile)" -ForegroundColor Yellow
        
        # 顯示摘要
        Write-Host "`n📊 發現摘要:" -ForegroundColor Cyan
        foreach ($category in $this.FoundData.Keys) {
            Write-Host "  $category`: $($this.FoundData[$category].Count) 項發現" -ForegroundColor Gray
        }
    }
}

# 主要執行邏輯
$hunter = [UltimateTokenHunter]::new()

if ($args.Count -eq 0) {
    Write-Host @"
終極 Token 獵人 - 多維度獲取 Kiro IDE 對話 Token

用法:
  .\ultimate_token_hunter.ps1 [選項]

選項:
  -DeepScan        執行所有掃描模式
  -NetworkMonitor  啟用網路監控
  -MemoryAnalysis  啟用記憶體分析
  -FileSystemHunt  啟用檔案系統深度掃描

範例:
  .\ultimate_token_hunter.ps1 -DeepScan
  .\ultimate_token_hunter.ps1 -FileSystemHunt -NetworkMonitor

"@ -ForegroundColor Cyan
} else {
    $hunter.StartHunt()
}