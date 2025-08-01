#!/usr/bin/env pwsh

<#
.SYNOPSIS
    TokenMonitor 一鍵安裝腳本 - 從 GitHub 下載並部署

.DESCRIPTION
    這個腳本會從 GitHub 下載 TokenMonitor 並自動部署到指定位置

.PARAMETER TargetPath
    目標安裝路徑

.PARAMETER Mode
    部署模式: full (完整), lite (輕量), npm (套件)

.PARAMETER Version
    要安裝的版本 (預設: main)

.PARAMETER GitHubRepo
    GitHub 倉庫 (預設: cy5407/TokenMonitor)

.EXAMPLE
    .\install-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
    完整安裝到指定路徑

.EXAMPLE
    .\install-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode lite -Version "v1.0.0"
    安裝特定版本的輕量版
#>

param(
    [Parameter(Mandatory=$true, HelpMessage="目標安裝路徑")]
    [string]$TargetPath,
    
    [Parameter(HelpMessage="部署模式")]
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full",
    
    [Parameter(HelpMessage="版本")]
    [string]$Version = "main",
    
    [Parameter(HelpMessage="GitHub 倉庫")]
    [string]$GitHubRepo = "cy5407/TokenMonitor",
    
    [Parameter(HelpMessage="顯示幫助")]
    [switch]$Help
)

# 顯示幫助
if ($Help) {
    Write-Host @"
🚀 TokenMonitor 一鍵安裝工具

用法:
    install-tokenmonitor.ps1 -TargetPath <路徑> [選項]

參數:
    -TargetPath     目標安裝路徑 (必要)
    -Mode          部署模式 (full/lite/npm)
    -Version       版本標籤 (預設: main)
    -GitHubRepo    GitHub 倉庫 (預設: cy5407/TokenMonitor)
    -Help          顯示此幫助

部署模式:
    full    完整安裝 - 包含所有功能和工具
    lite    輕量安裝 - 只包含核心監控功能
    npm     NPM 套件 - 生成可發布的 NPM 套件

範例:
    .\install-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
    .\install-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode lite -Version "v1.0.0"

"@ -ForegroundColor Cyan
    exit 0
}

Write-Host "🚀 TokenMonitor 一鍵安裝工具" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Green
Write-Host "📁 目標路徑: $TargetPath" -ForegroundColor Gray
Write-Host "⚙️  部署模式: $Mode" -ForegroundColor Gray
Write-Host "🏷️  版本: $Version" -ForegroundColor Gray
Write-Host "📦 倉庫: $GitHubRepo" -ForegroundColor Gray
Write-Host ""

try {
    # 檢查目標路徑
    if (-not (Test-Path $TargetPath)) {
        Write-Host "📁 創建目標目錄: $TargetPath" -ForegroundColor Yellow
        New-Item -ItemType Directory -Path $TargetPath -Force | Out-Null
    }

    # 創建臨時目錄
    $tempDir = Join-Path $env:TEMP "TokenMonitor-Install-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    Write-Host "📂 創建臨時目錄: $tempDir" -ForegroundColor Gray

    # 構建下載 URL
    $zipUrl = "https://github.com/$GitHubRepo/archive/$Version.zip"
    $zipPath = Join-Path $tempDir "TokenMonitor.zip"
    
    Write-Host "📥 從 GitHub 下載中..." -ForegroundColor Yellow
    Write-Host "🔗 URL: $zipUrl" -ForegroundColor Gray
    
    # 下載檔案
    try {
        Invoke-WebRequest -Uri $zipUrl -OutFile $zipPath -UseBasicParsing
        Write-Host "✅ 下載完成" -ForegroundColor Green
    } catch {
        throw "下載失敗: $($_.Exception.Message)"
    }

    # 檢查下載的檔案
    if (-not (Test-Path $zipPath) -or (Get-Item $zipPath).Length -eq 0) {
        throw "下載的檔案無效或為空"
    }

    Write-Host "📦 解壓縮中..." -ForegroundColor Yellow
    
    # 解壓縮
    try {
        Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
        Write-Host "✅ 解壓縮完成" -ForegroundColor Green
    } catch {
        throw "解壓縮失敗: $($_.Exception.Message)"
    }

    # 找到解壓縮的目錄
    $extractedDirs = Get-ChildItem -Path $tempDir -Directory | Where-Object { $_.Name -like "TokenMonitor-*" }
    
    if ($extractedDirs.Count -eq 0) {
        throw "找不到解壓縮的 TokenMonitor 目錄"
    }
    
    $extractedDir = $extractedDirs[0].FullName
    Write-Host "📂 找到解壓縮目錄: $($extractedDirs[0].Name)" -ForegroundColor Gray

    # 根據模式執行不同的安裝
    Write-Host "🔧 執行 $Mode 模式安裝..." -ForegroundColor Yellow

    switch ($Mode) {
        "full" {
            # 檢查部署腳本是否存在
            $deployScript = Join-Path $extractedDir "scripts\deploy-tokenmonitor.ps1"
            
            if (Test-Path $deployScript) {
                Write-Host "🚀 執行完整部署腳本..." -ForegroundColor Yellow
                & $deployScript -TargetPath $TargetPath -Mode full -Force
            } else {
                # 手動複製檔案 (備用方案)
                Write-Host "⚠️  找不到部署腳本，執行手動安裝..." -ForegroundColor Yellow
                
                $tokenMonitorDir = Join-Path $TargetPath "TokenMonitor"
                if (-not (Test-Path $tokenMonitorDir)) {
                    New-Item -ItemType Directory -Path $tokenMonitorDir -Force | Out-Null
                }
                
                # 複製主要檔案
                $filesToCopy = @(
                    @{ src = "scripts"; dest = "scripts" },
                    @{ src = "src"; dest = "src" },
                    @{ src = ".kiro"; dest = ".kiro" },
                    @{ src = "docs"; dest = "docs" },
                    @{ src = "package.json"; dest = "package.json" }
                )
                
                foreach ($file in $filesToCopy) {
                    $srcPath = Join-Path $extractedDir $file.src
                    $destPath = Join-Path $tokenMonitorDir $file.dest
                    
                    if (Test-Path $srcPath) {
                        if (Test-Path $srcPath -PathType Container) {
                            Copy-Item -Path $srcPath -Destination $destPath -Recurse -Force
                        } else {
                            Copy-Item -Path $srcPath -Destination $destPath -Force
                        }
                        Write-Host "✅ 複製: $($file.src)" -ForegroundColor Green
                    }
                }
                
                # 創建必要目錄
                $requiredDirs = @("data", "tests/data", "tests/reports")
                foreach ($dir in $requiredDirs) {
                    $dirPath = Join-Path $tokenMonitorDir $dir
                    if (-not (Test-Path $dirPath)) {
                        New-Item -ItemType Directory -Path $dirPath -Force | Out-Null
                        Write-Host "📁 創建目錄: $dir" -ForegroundColor Blue
                    }
                }
            }
        }
        
        "lite" {
            Write-Host "📦 安裝輕量版..." -ForegroundColor Yellow
            
            $liteDir = Join-Path $TargetPath "token-monitor"
            if (-not (Test-Path $liteDir)) {
                New-Item -ItemType Directory -Path $liteDir -Force | Out-Null
            }
            
            # 創建輕量版監控腳本
            $liteScriptContent = @"
// TokenMonitor Lite - 輕量版 Token 監控
const fs = require('fs');
const path = require('path');

class TokenMonitor {
    constructor(logPath = './token-usage.log') {
        this.logPath = logPath;
    }
    
    log(event, tokens, cost = 0, model = 'unknown') {
        const record = {
            timestamp: new Date().toISOString(),
            event,
            tokens: parseInt(tokens),
            cost: parseFloat(cost),
            model,
            session: 'lite-' + Date.now()
        };
        
        fs.appendFileSync(this.logPath, JSON.stringify(record) + '\n');
    }
    
    analyze(days = 7) {
        if (!fs.existsSync(this.logPath)) {
            return { total: 0, cost: 0, records: 0 };
        }
        
        const lines = fs.readFileSync(this.logPath, 'utf8').split('\n').filter(Boolean);
        const cutoff = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
        
        const records = lines
            .map(line => JSON.parse(line))
            .filter(r => new Date(r.timestamp) > cutoff);
        
        const total = records.reduce((sum, r) => sum + r.tokens, 0);
        const cost = records.reduce((sum, r) => sum + r.cost, 0);
        
        return { 
            total, 
            cost: cost.toFixed(6), 
            records: records.length,
            daily: (total / days).toFixed(0),
            dailyCost: (cost / days).toFixed(6)
        };
    }
    
    report() {
        const stats = this.analyze();
        console.log('📊 TokenMonitor Lite 報告');
        console.log('========================');
        console.log(`總 Token: ${stats.total}`);
        console.log(`總成本: $${stats.cost}`);
        console.log(`記錄數: ${stats.records}`);
        console.log(`日均 Token: ${stats.daily}`);
        console.log(`日均成本: $${stats.dailyCost}`);
    }
}

module.exports = TokenMonitor;

// CLI 使用
if (require.main === module) {
    const monitor = new TokenMonitor();
    const command = process.argv[2];
    
    switch (command) {
        case 'report':
            monitor.report();
            break;
        case 'log':
            const [, , , event, tokens, cost] = process.argv;
            monitor.log(event, tokens, cost);
            console.log(`✅ 記錄: ${event} - ${tokens} tokens`);
            break;
        default:
            console.log('用法: node token-monitor.js [report|log <event> <tokens> <cost>]');
    }
}
"@
            
            $liteScriptPath = Join-Path $liteDir "token-monitor.js"
            Set-Content -Path $liteScriptPath -Value $liteScriptContent -Encoding UTF8
            Write-Host "✅ 創建: token-monitor.js" -ForegroundColor Green
            
            # 創建 README
            $readmeContent = @"
# TokenMonitor Lite

輕量版 Token 使用監控工具

## 使用方式

\`\`\`javascript
const TokenMonitor = require('./token-monitor');
const monitor = new TokenMonitor();

// 記錄使用
monitor.log('chat_message', 150, 0.00045);

// 查看報告
monitor.report();
\`\`\`

## CLI 使用

\`\`\`bash
# 查看報告
node token-monitor.js report

# 記錄使用
node token-monitor.js log chat_message 150 0.00045
\`\`\`

## 安裝來源

此輕量版由 TokenMonitor 一鍵安裝腳本自動生成
GitHub: https://github.com/$GitHubRepo
版本: $Version
安裝時間: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')
"@
            
            $readmePath = Join-Path $liteDir "README.md"
            Set-Content -Path $readmePath -Value $readmeContent -Encoding UTF8
            Write-Host "✅ 創建: README.md" -ForegroundColor Green
        }
        
        "npm" {
            Write-Host "📦 安裝 NPM 套件模板..." -ForegroundColor Yellow
            
            $npmDir = Join-Path $TargetPath "kiro-token-monitor"
            if (-not (Test-Path $npmDir)) {
                New-Item -ItemType Directory -Path $npmDir -Force | Out-Null
            }
            
            # 複製 NPM 套件模板
            $npmTemplateDir = Join-Path $extractedDir "templates\npm-package"
            if (Test-Path $npmTemplateDir) {
                Copy-Item -Path "$npmTemplateDir\*" -Destination $npmDir -Recurse -Force
                Write-Host "✅ 複製 NPM 套件模板" -ForegroundColor Green
            } else {
                Write-Warning "找不到 NPM 套件模板，創建基本結構..."
                
                # 創建基本的 package.json
                $packageJson = @{
                    name = "kiro-token-monitor"
                    version = "1.0.0"
                    description = "AI Token usage monitoring tool"
                    main = "index.js"
                    bin = @{ tokusage = "./bin/tokusage.js" }
                    dependencies = @{
                        commander = "^9.0.0"
                        chokidar = "^3.5.3"
                    }
                } | ConvertTo-Json -Depth 10
                
                Set-Content -Path (Join-Path $npmDir "package.json") -Value $packageJson -Encoding UTF8
                Write-Host "✅ 創建基本 package.json" -ForegroundColor Green
            }
        }
    }

    Write-Host ""
    Write-Host "🎉 TokenMonitor 安裝完成！" -ForegroundColor Green
    Write-Host ""
    
    # 顯示後續步驟
    Write-Host "📋 後續步驟:" -ForegroundColor Cyan
    switch ($Mode) {
        "full" {
            Write-Host "1. cd `"$TargetPath\TokenMonitor`"" -ForegroundColor Gray
            Write-Host "2. npm install" -ForegroundColor Gray
            Write-Host "3. .\scripts\tokusage.ps1 daily" -ForegroundColor Gray
        }
        "lite" {
            Write-Host "1. cd `"$TargetPath\token-monitor`"" -ForegroundColor Gray
            Write-Host "2. node token-monitor.js report" -ForegroundColor Gray
        }
        "npm" {
            Write-Host "1. cd `"$TargetPath\kiro-token-monitor`"" -ForegroundColor Gray
            Write-Host "2. npm install" -ForegroundColor Gray
            Write-Host "3. node bin/tokusage.js --help" -ForegroundColor Gray
        }
    }
    
    Write-Host ""
    Write-Host "💡 提示:" -ForegroundColor Yellow
    Write-Host "• 查看文件了解更多功能" -ForegroundColor Gray
    Write-Host "• 定期更新以獲得最新功能" -ForegroundColor Gray
    Write-Host "• 遇到問題請查看 GitHub Issues" -ForegroundColor Gray
    Write-Host ""
    Write-Host "🌟 如果覺得有用，請給我們一個 Star！" -ForegroundColor Yellow
    Write-Host "🔗 https://github.com/$GitHubRepo" -ForegroundColor Blue

} catch {
    Write-Host ""
    Write-Error "❌ 安裝失敗: $($_.Exception.Message)"
    Write-Host ""
    Write-Host "🔧 故障排除建議:" -ForegroundColor Yellow
    Write-Host "• 檢查網路連線" -ForegroundColor Gray
    Write-Host "• 確認 GitHub 倉庫存在且可訪問" -ForegroundColor Gray
    Write-Host "• 檢查目標路徑權限" -ForegroundColor Gray
    Write-Host "• 嘗試使用不同的版本標籤" -ForegroundColor Gray
    exit 1
} finally {
    # 清理臨時檔案
    if (Test-Path $tempDir) {
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
        Write-Host "🧹 清理臨時檔案" -ForegroundColor Gray
    }
}
