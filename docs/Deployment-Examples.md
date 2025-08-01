# 🚀 TokenMonitor 部署範例

## 📋 實際部署案例

### 案例一：完整部署到新專案

```powershell
# 1. 部署到目標專案
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\MyAIProject" -Mode full

# 2. 進入目標專案
cd "C:\MyAIProject\TokenMonitor"

# 3. 安裝依賴
npm install

# 4. 測試運行
.\scripts\tokusage.ps1 daily
```

**預期結果**:
```
🔍 執行每日 Token 使用分析...
📋 分析 0 筆記錄...

┌────────────────────────────────────────┐
│ Claude Code Token Usage Report - Daily │
└────────────────────────────────────────┘

📊 使用統計摘要:
   • 記錄天數: 0 天
   • 總 Token: 0
   • 總成本: $0.00
```

### 案例二：輕量部署到現有專案

```powershell
# 1. 輕量部署
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\ExistingProject" -Mode lite

# 2. 進入專案
cd "C:\ExistingProject\token-monitor"

# 3. 測試基本功能
node token-monitor.js report
```

**預期結果**:
```
📊 TokenMonitor Lite 報告
========================
總 Token: 0
總成本: $0.000000
記錄數: 0
日均 Token: 0
日均成本: $0.000000
```

### 案例三：NPM 套件部署

```powershell
# 1. 生成 NPM 套件
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\NPMPackages" -Mode npm

# 2. 進入套件目錄
cd "C:\NPMPackages\kiro-token-monitor"

# 3. 安裝依賴
npm install

# 4. 測試 CLI
node bin/tokusage.js --help
```

**預期結果**:
```
Usage: tokusage [options] [command]

AI Token usage monitoring and analysis tool

Options:
  -V, --version     output the version number
  -h, --help        display help for command

Commands:
  daily [options]   Show daily token usage report
  summary [options] Show detailed usage summary
  cleanup [options] Clean up old usage records
  log <event> <tokens> [cost] Manually log token usage
  status            Show monitoring status
  install [options] Install TokenMonitor in current project
  help [command]    display help for command
```

---

## 🔧 實際使用場景

### 場景一：React 專案整合

```bash
# 在 React 專案中
cd my-react-app

# 安裝 TokenMonitor
npm install kiro-token-monitor

# 在 package.json 中添加腳本
{
  "scripts": {
    "token:daily": "tokusage daily",
    "token:summary": "tokusage summary",
    "token:status": "tokusage status"
  }
}

# 使用
npm run token:daily
```

### 場景二：Python 專案整合

```bash
# 在 Python 專案中創建 token-monitor 目錄
mkdir token-monitor
cd token-monitor

# 複製輕量版監控腳本
# (使用 deploy-tokenmonitor.ps1 -Mode lite)

# 在 Python 中使用
python -c "
import subprocess
import json

# 記錄使用
subprocess.run(['node', 'token-monitor.js', 'log', 'python_script', '200', '0.0006'])

# 查看報告
subprocess.run(['node', 'token-monitor.js', 'report'])
"
```

### 場景三：CI/CD 整合

```yaml
# .github/workflows/token-monitor.yml
name: Daily Token Report

on:
  schedule:
    - cron: '0 9 * * *'  # 每日 9AM UTC

jobs:
  token-report:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Node.js
        uses: actions/setup-node@v2
        with:
          node-version: '16'
          
      - name: Install TokenMonitor
        run: npm install kiro-token-monitor
        
      - name: Generate Daily Report
        run: npx tokusage daily
        
      - name: Upload Report
        uses: actions/upload-artifact@v2
        with:
          name: token-usage-report
          path: data/kiro-usage.log
```

---

## 📊 部署驗證腳本

### 自動驗證腳本

```powershell
# verify-deployment.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$DeploymentPath,
    
    [Parameter()]
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full"
)

Write-Host "🔍 驗證 TokenMonitor 部署..." -ForegroundColor Green

$errors = @()

switch ($Mode) {
    "full" {
        $requiredFiles = @(
            "TokenMonitor/scripts/tokusage.ps1",
            "TokenMonitor/src/js/professional-token-cli.js",
            "TokenMonitor/package.json"
        )
        
        foreach ($file in $requiredFiles) {
            $filePath = Join-Path $DeploymentPath $file
            if (-not (Test-Path $filePath)) {
                $errors += "缺少檔案: $file"
            } else {
                Write-Host "✅ 檔案存在: $file" -ForegroundColor Green
            }
        }
        
        # 測試 CLI 功能
        try {
            Push-Location (Join-Path $DeploymentPath "TokenMonitor")
            $output = & ".\scripts\tokusage.ps1" "daily" 2>&1
            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ CLI 功能正常" -ForegroundColor Green
            } else {
                $errors += "CLI 測試失敗: $output"
            }
        } catch {
            $errors += "CLI 執行錯誤: $($_.Exception.Message)"
        } finally {
            Pop-Location
        }
    }
    
    "lite" {
        $liteFile = Join-Path $DeploymentPath "token-monitor/token-monitor.js"
        if (Test-Path $liteFile) {
            Write-Host "✅ 輕量版檔案存在" -ForegroundColor Green
            
            # 測試輕量版功能
            try {
                Push-Location (Join-Path $DeploymentPath "token-monitor")
                $output = node "token-monitor.js" "report" 2>&1
                if ($LASTEXITCODE -eq 0) {
                    Write-Host "✅ 輕量版功能正常" -ForegroundColor Green
                } else {
                    $errors += "輕量版測試失敗: $output"
                }
            } catch {
                $errors += "輕量版執行錯誤: $($_.Exception.Message)"
            } finally {
                Pop-Location
            }
        } else {
            $errors += "缺少輕量版檔案"
        }
    }
    
    "npm" {
        $packageFile = Join-Path $DeploymentPath "kiro-token-monitor/package.json"
        if (Test-Path $packageFile) {
            Write-Host "✅ NPM 套件檔案存在" -ForegroundColor Green
        } else {
            $errors += "缺少 NPM 套件檔案"
        }
    }
}

if ($errors.Count -eq 0) {
    Write-Host "🎉 部署驗證成功！" -ForegroundColor Green
    exit 0
} else {
    Write-Host "❌ 部署驗證失敗:" -ForegroundColor Red
    $errors | ForEach-Object { Write-Host "  - $_" -ForegroundColor Red }
    exit 1
}
```

### 使用驗證腳本

```powershell
# 驗證完整部署
.\verify-deployment.ps1 -DeploymentPath "C:\MyProject" -Mode full

# 驗證輕量部署
.\verify-deployment.ps1 -DeploymentPath "C:\MyProject" -Mode lite
```

---

## 🎯 部署最佳實務

### 1. 版本管理

```json
// 在目標專案的 package.json 中
{
  "devDependencies": {
    "kiro-token-monitor": "^1.0.0"
  },
  "scripts": {
    "postinstall": "tokusage install",
    "token:check": "tokusage status",
    "token:report": "tokusage daily"
  }
}
```

### 2. 環境配置

```bash
# .env 檔案
TOKEN_LOG_LEVEL=info
TOKEN_MAX_LOG_SIZE=10MB
TOKEN_RETENTION_DAYS=30
TOKEN_AUTO_CLEANUP=true
```

### 3. 忽略檔案設定

```gitignore
# .gitignore
data/kiro-usage.log
TokenMonitor/data/
token-monitor/token-usage.log
node_modules/
```

### 4. 文件模板

```markdown
# 專案中的 TOKEN-USAGE.md
# Token 使用監控

本專案已整合 TokenMonitor 系統。

## 查看使用情況

\`\`\`bash
npm run token:report
\`\`\`

## 清理舊記錄

\`\`\`bash
npx tokusage cleanup --days 30
\`\`\`

## 配置檔案

- `.kiro/settings/token-monitor.json` - 主要配置
- `data/kiro-usage.log` - 使用記錄
```

---

這些範例展示了如何在不同場景下部署和使用 TokenMonitor，確保系統能夠適應各種專案需求！