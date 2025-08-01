#!/usr/bin/env pwsh

<#
.SYNOPSIS
    TokenMonitor 部署腳本 - 將 TokenMonitor 部署到其他專案

.DESCRIPTION
    這個腳本可以將 TokenMonitor 系統部署到任何專案中，支援多種部署模式

.PARAMETER TargetPath
    目標專案路徑

.PARAMETER Mode
    部署模式: full (完整), lite (輕量), npm (套件)

.PARAMETER Force
    強制覆蓋現有檔案

.EXAMPLE
    .\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
    完整部署到指定專案

.EXAMPLE
    .\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode lite
    輕量部署，只包含核心功能
#>

param(
    [Parameter(Mandatory=$false, HelpMessage="目標專案路徑")]
    [string]$TargetPath,
    
    [Parameter(HelpMessage="部署模式")]
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full",
    
    [Parameter(HelpMessage="強制覆蓋現有檔案")]
    [switch]$Force,
    
    [Parameter(HelpMessage="顯示幫助")]
    [switch]$Help
)

# 顯示幫助
if ($Help) {
    Write-Host @"
🚀 TokenMonitor 部署工具

用法:
    deploy-tokenmonitor.ps1 -TargetPath <路徑> [選項]

參數:
    -TargetPath     目標專案路徑 (必要)
    -Mode          部署模式 (full/lite/npm)
    -Force         強制覆蓋現有檔案
    -Help          顯示此幫助

部署模式:
    full    完整部署 - 包含所有功能和工具
    lite    輕量部署 - 只包含核心監控功能
    npm     NPM 套件 - 生成可發布的 NPM 套件

範例:
    .\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
    .\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode lite -Force

"@ -ForegroundColor Cyan
    exit 0
}

# 檢查必要參數
if (-not $TargetPath -and -not $Help) {
    Write-Error "請提供目標路徑參數 -TargetPath"
    Write-Host "使用 -Help 查看詳細說明" -ForegroundColor Yellow
    exit 1
}

# 檢查目標路徑
if ($TargetPath -and -not (Test-Path $TargetPath)) {
    Write-Error "目標路徑不存在: $TargetPath"
    exit 1
}

# 獲取腳本目錄
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent $ScriptDir

Write-Host "🚀 開始部署 TokenMonitor" -ForegroundColor Green
Write-Host "📁 來源路徑: $ProjectRoot" -ForegroundColor Gray
Write-Host "📁 目標路徑: $TargetPath" -ForegroundColor Gray
Write-Host "⚙️  部署模式: $Mode" -ForegroundColor Gray
Write-Host ""

try {
    switch ($Mode) {
        "full" {
            Write-Host "📦 執行完整部署..." -ForegroundColor Yellow
            
            # 完整部署檔案清單
            $requiredFiles = @(
                @{ src = "scripts/tokusage.ps1"; dest = "TokenMonitor/scripts/tokusage.ps1" },
                @{ src = "scripts/universal-token-monitor.js"; dest = "TokenMonitor/scripts/universal-token-monitor.js" },
                @{ src = "src/js/professional-token-cli.js"; dest = "TokenMonitor/src/js/professional-token-cli.js" },
                @{ src = "src/js/enhanced-token-reporter.js"; dest = "TokenMonitor/src/js/enhanced-token-reporter.js" },
                @{ src = ".kiro/hooks/manual-token-calc.js"; dest = "TokenMonitor/.kiro/hooks/manual-token-calc.js" },
                @{ src = ".kiro/hooks/manual-token-calc.json"; dest = "TokenMonitor/.kiro/hooks/manual-token-calc.json" },
                @{ src = "docs/README.md"; dest = "TokenMonitor/docs/README.md" },
                @{ src = "docs/USAGE-GUIDE.md"; dest = "TokenMonitor/docs/USAGE-GUIDE.md" },
                @{ src = "package.json"; dest = "TokenMonitor/package.json" }
            )
            
            foreach ($file in $requiredFiles) {
                $sourcePath = Join-Path $ProjectRoot $file.src
                $destPath = Join-Path $TargetPath $file.dest
                $destDir = Split-Path $destPath -Parent
                
                if (-not (Test-Path $destDir)) {
                    New-Item -ItemType Directory -Path $destDir -Force | Out-Null
                }
                
                if (Test-Path $sourcePath) {
                    if ((Test-Path $destPath) -and -not $Force) {
                        Write-Warning "檔案已存在，跳過: $($file.dest)"
                    } else {
                        Copy-Item $sourcePath $destPath -Force
                        Write-Host "✅ 複製: $($file.dest)" -ForegroundColor Green
                    }
                } else {
                    Write-Warning "來源檔案不存在: $($file.src)"
                }
            }
            
            # 創建必要目錄
            $requiredDirs = @("TokenMonitor/data", "TokenMonitor/tests/data", "TokenMonitor/tests/reports")
            foreach ($dir in $requiredDirs) {
                $dirPath = Join-Path $TargetPath $dir
                if (-not (Test-Path $dirPath)) {
                    New-Item -ItemType Directory -Path $dirPath -Force | Out-Null
                    Write-Host "📁 創建目錄: $dir" -ForegroundColor Blue
                }
            }
        }
        
        "lite" {
            Write-Host "📦 執行輕量部署..." -ForegroundColor Yellow
            
            # 創建輕量版監控腳本
            $liteMonitorContent = @"
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
            
            $liteDir = Join-Path $TargetPath "token-monitor"
            if (-not (Test-Path $liteDir)) {
                New-Item -ItemType Directory -Path $liteDir -Force | Out-Null
            }
            
            $liteScriptPath = Join-Path $liteDir "token-monitor.js"
            Set-Content -Path $liteScriptPath -Value $liteMonitorContent -Encoding UTF8
            Write-Host "✅ 創建: token-monitor/token-monitor.js" -ForegroundColor Green
            
            # 創建使用說明
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
"@
            
            $readmePath = Join-Path $liteDir "README.md"
            Set-Content -Path $readmePath -Value $readmeContent -Encoding UTF8
            Write-Host "✅ 創建: token-monitor/README.md" -ForegroundColor Green
        }
        
        "npm" {
            Write-Host "📦 生成 NPM 套件..." -ForegroundColor Yellow
            
            $npmDir = Join-Path $TargetPath "kiro-token-monitor"
            if (-not (Test-Path $npmDir)) {
                New-Item -ItemType Directory -Path $npmDir -Force | Out-Null
            }
            
            # 創建 package.json
            $packageJson = @{
                name = "kiro-token-monitor"
                version = "1.0.0"
                description = "AI Token usage monitoring for Kiro IDE and other development environments"
                main = "index.js"
                bin = @{
                    tokusage = "./bin/tokusage.js"
                }
                scripts = @{
                    install = "node install.js"
                    test = "node test.js"
                }
                dependencies = @{
                    chokidar = "^3.5.3"
                    commander = "^9.0.0"
                }
                keywords = @("token", "monitoring", "ai", "kiro", "cost", "analysis")
                author = "TokenMonitor Team"
                license = "MIT"
            }
            
            $packageJsonPath = Join-Path $npmDir "package.json"
            $packageJson | ConvertTo-Json -Depth 10 | Set-Content -Path $packageJsonPath -Encoding UTF8
            Write-Host "✅ 創建: package.json" -ForegroundColor Green
            
            Write-Host "📦 NPM 套件已準備完成" -ForegroundColor Green
            Write-Host "💡 執行 'cd $npmDir && npm publish' 來發布套件" -ForegroundColor Yellow
        }
    }
    
    Write-Host ""
    Write-Host "🎉 部署完成！" -ForegroundColor Green
    
    # 顯示後續步驟
    switch ($Mode) {
        "full" {
            Write-Host "📋 後續步驟:" -ForegroundColor Cyan
            Write-Host "1. cd `"$TargetPath/TokenMonitor`"" -ForegroundColor Gray
            Write-Host "2. npm install" -ForegroundColor Gray
            Write-Host "3. .\scripts\tokusage.ps1 daily" -ForegroundColor Gray
        }
        "lite" {
            Write-Host "📋 後續步驟:" -ForegroundColor Cyan
            Write-Host "1. cd `"$TargetPath/token-monitor`"" -ForegroundColor Gray
            Write-Host "2. node token-monitor.js report" -ForegroundColor Gray
        }
        "npm" {
            Write-Host "📋 後續步驟:" -ForegroundColor Cyan
            Write-Host "1. cd `"$TargetPath/kiro-token-monitor`"" -ForegroundColor Gray
            Write-Host "2. npm publish" -ForegroundColor Gray
        }
    }
    
} catch {
    Write-Error "部署失敗: $($_.Exception.Message)"
    exit 1
} finally {
    # 清理
}

Write-Host ""
Write-Host "✨ TokenMonitor 已成功部署到您的專案！" -ForegroundColor Green