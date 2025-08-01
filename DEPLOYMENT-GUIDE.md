# 🚀 TokenMonitor 部署指南

## 📋 概述

TokenMonitor 是一個跨 IDE 的 Token 使用監控系統，可以輕鬆部署到任何專案中，提供專業的 AI Token 使用分析和成本控制。

## 🎯 部署方式

### 方式一：完整部署 (推薦)
適合需要完整功能的專案

### 方式二：輕量部署
適合只需要基本監控的專案

### 方式三：NPM 套件部署
適合 Node.js 專案的快速整合

---

## 🔧 方式一：完整部署

### 1. 複製核心檔案

```bash
# 創建 TokenMonitor 目錄
mkdir TokenMonitor
cd TokenMonitor

# 複製必要檔案結構
mkdir -p {scripts,src/js,data,docs,.kiro/hooks}
```

### 2. 複製核心檔案清單

**必要檔案**:
```
TokenMonitor/
├── scripts/
│   ├── tokusage.ps1                    # 主要 CLI 工具
│   └── universal-token-monitor.js      # 監控核心
├── src/js/
│   ├── professional-token-cli.js       # 專業報表生成器
│   └── enhanced-token-reporter.js      # 增強報表工具
├── .kiro/hooks/
│   ├── manual-token-calc.js            # 手動計算工具
│   └── manual-token-calc.json          # Hook 配置
├── data/
│   └── kiro-usage.log                  # 使用記錄 (自動生成)
├── docs/
│   └── README.md                       # 使用說明
└── package.json                        # 依賴配置
```

### 3. 安裝依賴

```bash
npm install chokidar fs path
```

### 4. 配置啟動

```powershell
# 啟動監控
.\scripts\tokusage.ps1 daily
```

---

## ⚡ 方式二：輕量部署

### 1. 最小檔案集

只需要這些檔案：
```
project/
├── token-monitor/
│   ├── monitor.js          # 簡化監控腳本
│   ├── analyzer.js         # 簡化分析工具
│   └── usage.log          # 記錄檔案
└── package.json
```

### 2. 創建簡化監控腳本

```javascript
// token-monitor/monitor.js
const fs = require('fs');
const path = require('path');

class SimpleTokenMonitor {
    constructor() {
        this.logFile = path.join(__dirname, 'usage.log');
    }
    
    log(event, tokens, cost = 0) {
        const record = {
            timestamp: new Date().toISOString(),
            event,
            tokens,
            cost,
            session: Date.now()
        };
        
        fs.appendFileSync(this.logFile, JSON.stringify(record) + '\n');
    }
    
    analyze() {
        if (!fs.existsSync(this.logFile)) return { total: 0, cost: 0 };
        
        const lines = fs.readFileSync(this.logFile, 'utf8').split('\n').filter(Boolean);
        const records = lines.map(line => JSON.parse(line));
        
        const total = records.reduce((sum, r) => sum + r.tokens, 0);
        const cost = records.reduce((sum, r) => sum + r.cost, 0);
        
        return { total, cost, records: records.length };
    }
}

module.exports = SimpleTokenMonitor;
```

### 3. 使用方式

```javascript
const TokenMonitor = require('./token-monitor/monitor');
const monitor = new TokenMonitor();

// 記錄使用
monitor.log('chat_message', 150, 0.00045);

// 分析結果
console.log(monitor.analyze());
```

---

## 📦 方式三：NPM 套件部署

### 1. 創建 NPM 套件

```bash
# 初始化套件
npm init -y

# 設定 package.json
```

```json
{
  "name": "kiro-token-monitor",
  "version": "1.0.0",
  "description": "AI Token usage monitoring for Kiro IDE",
  "main": "index.js",
  "bin": {
    "tokusage": "./bin/tokusage.js"
  },
  "scripts": {
    "install": "node install.js",
    "monitor": "node bin/tokusage.js"
  },
  "dependencies": {
    "chokidar": "^3.5.3",
    "commander": "^9.0.0"
  }
}
```

### 2. 創建安裝腳本

```javascript
// install.js
const fs = require('fs');
const path = require('path');

console.log('🚀 安裝 TokenMonitor...');

// 創建必要目錄
const dirs = ['data', '.kiro/hooks'];
dirs.forEach(dir => {
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
        console.log(`✅ 創建目錄: ${dir}`);
    }
});

// 複製配置檔案
const configs = [
    { src: 'templates/hook-config.json', dest: '.kiro/hooks/token-monitor.json' }
];

configs.forEach(({ src, dest }) => {
    if (fs.existsSync(src)) {
        fs.copyFileSync(src, dest);
        console.log(`✅ 複製配置: ${dest}`);
    }
});

console.log('🎉 TokenMonitor 安裝完成！');
console.log('使用 "npx tokusage daily" 開始監控');
```

### 3. 發布和使用

```bash
# 發布套件
npm publish

# 在其他專案中使用
npm install kiro-token-monitor
npx tokusage daily
```

---

## 🔄 自動化部署腳本

### PowerShell 部署腳本

```powershell
# deploy-tokenmonitor.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$TargetPath,
    
    [Parameter()]
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full"
)

Write-Host "🚀 部署 TokenMonitor 到: $TargetPath" -ForegroundColor Green

switch ($Mode) {
    "full" {
        # 完整部署邏輯
        $requiredFiles = @(
            "scripts/tokusage.ps1",
            "src/js/professional-token-cli.js",
            "src/js/enhanced-token-reporter.js",
            ".kiro/hooks/manual-token-calc.js",
            ".kiro/hooks/manual-token-calc.json"
        )
        
        foreach ($file in $requiredFiles) {
            $sourcePath = Join-Path $PSScriptRoot $file
            $destPath = Join-Path $TargetPath $file
            $destDir = Split-Path $destPath -Parent
            
            if (-not (Test-Path $destDir)) {
                New-Item -ItemType Directory -Path $destDir -Force | Out-Null
            }
            
            Copy-Item $sourcePath $destPath -Force
            Write-Host "✅ 複製: $file" -ForegroundColor Green
        }
    }
    
    "lite" {
        # 輕量部署邏輯
        Write-Host "📦 輕量部署模式" -ForegroundColor Yellow
        # 實作輕量部署邏輯
    }
    
    "npm" {
        # NPM 部署邏輯
        Write-Host "📦 NPM 套件部署" -ForegroundColor Yellow
        # 實作 NPM 部署邏輯
    }
}

Write-Host "🎉 部署完成！" -ForegroundColor Green
```

### 使用部署腳本

```powershell
# 完整部署到其他專案
.\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full

# 輕量部署
.\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode lite
```

---

## 📋 部署檢查清單

### 部署前檢查
- [ ] 確認目標專案支援 Node.js
- [ ] 檢查 PowerShell 執行權限
- [ ] 確認網路連線 (用於下載依賴)

### 部署後驗證
- [ ] 運行 `tokusage --help` 檢查 CLI
- [ ] 測試 `tokusage daily` 生成報告
- [ ] 檢查 `data/kiro-usage.log` 檔案生成
- [ ] 驗證 Hook 系統運作

### 故障排除
- **找不到檔案**: 檢查路徑配置
- **權限錯誤**: 確認 PowerShell 執行策略
- **依賴問題**: 運行 `npm install`

---

## 🎯 最佳實務

### 1. 版本控制
```bash
# 排除記錄檔案
echo "data/kiro-usage.log" >> .gitignore
echo "node_modules/" >> .gitignore
```

### 2. 配置管理
```javascript
// config.js
module.exports = {
    logLevel: process.env.TOKEN_LOG_LEVEL || 'info',
    maxLogSize: process.env.TOKEN_MAX_LOG_SIZE || '10MB',
    retentionDays: process.env.TOKEN_RETENTION_DAYS || 30
};
```

### 3. 監控自動化
```json
{
  "scripts": {
    "token:daily": "tokusage daily",
    "token:summary": "tokusage summary",
    "token:clean": "tokusage clean --days 30"
  }
}
```

---

## 🔧 客製化選項

### 1. 自訂報表格式
修改 `src/js/professional-token-cli.js` 中的報表模板

### 2. 整合 CI/CD
```yaml
# .github/workflows/token-monitor.yml
name: Token Usage Report
on:
  schedule:
    - cron: '0 9 * * *'  # 每日 9AM
jobs:
  token-report:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - run: npm install
      - run: npx tokusage daily
```

### 3. 警報設定
```javascript
// 成本警報
if (dailyCost > 1.0) {
    console.warn('⚠️  每日成本超過 $1.00');
    // 發送通知邏輯
}
```

---

## 📞 支援和維護

### 更新檢查
```bash
# 檢查更新
npm outdated kiro-token-monitor

# 更新到最新版本
npm update kiro-token-monitor
```

### 問題回報
如遇到問題，請提供：
1. 錯誤訊息
2. 系統環境 (OS, Node.js 版本)
3. 部署方式
4. 重現步驟

---

這個部署指南提供了多種靈活的部署方式，讓你可以根據不同專案的需求選擇最適合的部署策略！
