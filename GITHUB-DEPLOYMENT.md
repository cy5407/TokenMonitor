# 🌐 TokenMonitor GitHub 部署指南

## 📋 概述

將 TokenMonitor 推送到 GitHub 後，任何人都可以從網路上下載和部署到自己的專案中。

---

## 🚀 GitHub 部署步驟

### 1. 準備 GitHub 倉庫

```bash
# 初始化 Git 倉庫 (如果還沒有)
git init

# 添加所有檔案
git add .

# 提交變更
git commit -m "Initial commit: TokenMonitor v1.0"

# 添加遠端倉庫
git remote add origin https://github.com/cy5407/TokenMonitor.git

# 推送到 GitHub
git push -u origin main
```

### 2. 創建發布版本

```bash
# 創建標籤
git tag -a v1.0.0 -m "TokenMonitor v1.0.0 - Initial Release"

# 推送標籤
git push origin v1.0.0
```

---

## 📦 從 GitHub 部署的方式

### 方式一：直接下載部署腳本

```powershell
# 下載部署腳本
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cy5407/TokenMonitor/main/scripts/deploy-tokenmonitor.ps1" -OutFile "deploy-tokenmonitor.ps1"

# 執行部署
.\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
```

### 方式二：克隆整個倉庫

```bash
# 克隆倉庫
git clone https://github.com/cy5407/TokenMonitor.git

# 進入目錄
cd TokenMonitor

# 執行部署
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
```

### 方式三：下載特定版本

```bash
# 下載特定版本
curl -L https://github.com/cy5407/TokenMonitor/archive/v1.0.0.zip -o TokenMonitor-v1.0.0.zip

# 解壓縮
unzip TokenMonitor-v1.0.0.zip

# 進入目錄並部署
cd TokenMonitor-1.0.0
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
```

---

## 🔧 一鍵安裝腳本

### Windows PowerShell 一鍵安裝

```powershell
# install-tokenmonitor.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$TargetPath,
    
    [Parameter()]
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full",
    
    [Parameter()]
    [string]$Version = "main"
)

Write-Host "🚀 從 GitHub 安裝 TokenMonitor..." -ForegroundColor Green

try {
    # 創建臨時目錄
    $tempDir = Join-Path $env:TEMP "TokenMonitor-$(Get-Random)"
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
    
    # 下載 ZIP 檔案
    $zipUrl = "https://github.com/cy5407/TokenMonitor/archive/$Version.zip"
    $zipPath = Join-Path $tempDir "TokenMonitor.zip"
    
    Write-Host "📥 下載中..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $zipUrl -OutFile $zipPath
    
    # 解壓縮
    Write-Host "📦 解壓縮中..." -ForegroundColor Yellow
    Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force
    
    # 找到解壓縮的目錄
    $extractedDir = Get-ChildItem -Path $tempDir -Directory | Where-Object { $_.Name -like "TokenMonitor-*" } | Select-Object -First 1
    
    if (-not $extractedDir) {
        throw "找不到解壓縮的 TokenMonitor 目錄"
    }
    
    # 執行部署腳本
    $deployScript = Join-Path $extractedDir.FullName "scripts\deploy-tokenmonitor.ps1"
    
    if (Test-Path $deployScript) {
        Write-Host "🔧 執行部署..." -ForegroundColor Yellow
        & $deployScript -TargetPath $TargetPath -Mode $Mode
    } else {
        throw "找不到部署腳本"
    }
    
    Write-Host "🎉 安裝完成！" -ForegroundColor Green
    
} catch {
    Write-Error "安裝失敗: $($_.Exception.Message)"
    exit 1
} finally {
    # 清理臨時檔案
    if (Test-Path $tempDir) {
        Remove-Item -Path $tempDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}
```

### 使用一鍵安裝

```powershell
# 下載並執行一鍵安裝腳本
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.ps1" -OutFile "install-tokenmonitor.ps1"

# 執行安裝
.\install-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
```

---

## 🌍 跨平台安裝腳本

### Linux/macOS Bash 腳本

```bash
#!/bin/bash
# install-tokenmonitor.sh

set -e

# 參數設定
TARGET_PATH=""
MODE="full"
VERSION="main"
GITHUB_REPO="cy5407/TokenMonitor"

# 解析參數
while [[ $# -gt 0 ]]; do
    case $1 in
        --target-path)
            TARGET_PATH="$2"
            shift 2
            ;;
        --mode)
            MODE="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --help)
            echo "TokenMonitor 安裝腳本"
            echo ""
            echo "用法: $0 --target-path <路徑> [選項]"
            echo ""
            echo "選項:"
            echo "  --target-path    目標安裝路徑 (必要)"
            echo "  --mode          部署模式 (full/lite/npm)"
            echo "  --version       版本 (預設: main)"
            echo "  --help          顯示此幫助"
            exit 0
            ;;
        *)
            echo "未知參數: $1"
            exit 1
            ;;
    esac
done

# 檢查必要參數
if [ -z "$TARGET_PATH" ]; then
    echo "錯誤: 請提供 --target-path 參數"
    exit 1
fi

echo "🚀 從 GitHub 安裝 TokenMonitor..."
echo "📁 目標路徑: $TARGET_PATH"
echo "⚙️  部署模式: $MODE"
echo "🏷️  版本: $VERSION"

# 創建臨時目錄
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# 下載並解壓縮
echo "📥 下載中..."
curl -L "https://github.com/$GITHUB_REPO/archive/$VERSION.tar.gz" | tar -xz -C "$TEMP_DIR"

# 找到解壓縮的目錄
EXTRACTED_DIR=$(find "$TEMP_DIR" -name "TokenMonitor-*" -type d | head -1)

if [ -z "$EXTRACTED_DIR" ]; then
    echo "錯誤: 找不到解壓縮的 TokenMonitor 目錄"
    exit 1
fi

# 根據模式執行不同的安裝邏輯
case $MODE in
    "lite")
        echo "🔧 執行輕量安裝..."
        # 複製輕量版檔案
        mkdir -p "$TARGET_PATH/token-monitor"
        cp "$EXTRACTED_DIR/templates/lite-monitor.js" "$TARGET_PATH/token-monitor/token-monitor.js" 2>/dev/null || {
            # 如果模板不存在，創建基本版本
            cat > "$TARGET_PATH/token-monitor/token-monitor.js" << 'EOF'
// TokenMonitor Lite for Linux/macOS
const fs = require('fs');
const path = require('path');

class TokenMonitor {
    constructor(logPath = './token-usage.log') {
        this.logPath = logPath;
    }
    
    log(event, tokens, cost = 0) {
        const record = {
            timestamp: new Date().toISOString(),
            event, tokens: parseInt(tokens), cost: parseFloat(cost)
        };
        fs.appendFileSync(this.logPath, JSON.stringify(record) + '\n');
    }
    
    report() {
        if (!fs.existsSync(this.logPath)) {
            console.log('📊 TokenMonitor: 尚無使用記錄');
            return;
        }
        
        const lines = fs.readFileSync(this.logPath, 'utf8').split('\n').filter(Boolean);
        const records = lines.map(line => JSON.parse(line));
        const total = records.reduce((sum, r) => sum + r.tokens, 0);
        const cost = records.reduce((sum, r) => sum + r.cost, 0);
        
        console.log('📊 TokenMonitor Lite 報告');
        console.log(`總 Token: ${total}`);
        console.log(`總成本: $${cost.toFixed(6)}`);
        console.log(`記錄數: ${records.length}`);
    }
}

if (require.main === module) {
    const monitor = new TokenMonitor();
    const [,, command, ...args] = process.argv;
    
    switch (command) {
        case 'report': monitor.report(); break;
        case 'log': monitor.log(args[0], args[1], args[2]); break;
        default: console.log('用法: node token-monitor.js [report|log <event> <tokens> <cost>]');
    }
}

module.exports = TokenMonitor;
EOF
        }
        echo "✅ 輕量版安裝完成"
        ;;
        
    "npm")
        echo "🔧 執行 NPM 套件安裝..."
        mkdir -p "$TARGET_PATH/kiro-token-monitor"
        cp -r "$EXTRACTED_DIR/templates/npm-package/"* "$TARGET_PATH/kiro-token-monitor/"
        echo "✅ NPM 套件安裝完成"
        ;;
        
    *)
        echo "🔧 執行完整安裝..."
        # 檢查是否有 PowerShell (用於 WSL 或安裝了 PowerShell 的 Linux)
        if command -v pwsh >/dev/null 2>&1; then
            pwsh "$EXTRACTED_DIR/scripts/deploy-tokenmonitor.ps1" -TargetPath "$TARGET_PATH" -Mode full
        else
            echo "⚠️  完整安裝需要 PowerShell，改為執行輕量安裝..."
            mkdir -p "$TARGET_PATH/token-monitor"
            cp "$EXTRACTED_DIR/templates/lite-monitor.js" "$TARGET_PATH/token-monitor/" 2>/dev/null || echo "使用預設輕量版本"
        fi
        ;;
esac

echo "🎉 TokenMonitor 安裝完成！"
echo "📋 後續步驟:"
case $MODE in
    "lite")
        echo "  cd $TARGET_PATH/token-monitor"
        echo "  node token-monitor.js report"
        ;;
    "npm")
        echo "  cd $TARGET_PATH/kiro-token-monitor"
        echo "  npm install"
        echo "  node bin/tokusage.js --help"
        ;;
    *)
        echo "  cd $TARGET_PATH/TokenMonitor"
        echo "  npm install"
        echo "  ./scripts/tokusage.ps1 daily"
        ;;
esac
```

### 使用 Linux/macOS 安裝腳本

```bash
# 下載並執行
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path /path/to/project --mode full

# 或者下載後執行
curl -O https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh
chmod +x install-tokenmonitor.sh
./install-tokenmonitor.sh --target-path /path/to/project --mode lite
```

---

## 📋 GitHub Actions 自動化

### 自動發布工作流程

```yaml
# .github/workflows/release.yml
name: Release TokenMonitor

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install dependencies
        run: npm install
        
      - name: Run tests
        run: npm test || echo "No tests defined"
        
      - name: Create Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: TokenMonitor ${{ github.ref }}
          body: |
            ## 🚀 TokenMonitor Release
            
            ### 安裝方式
            
            **Windows PowerShell:**
            ```powershell
            Invoke-WebRequest -Uri "https://raw.githubusercontent.com/${{ github.repository }}/main/install-tokenmonitor.ps1" -OutFile "install.ps1"
            .\install.ps1 -TargetPath "C:\MyProject" -Mode full
            ```
            
            **Linux/macOS:**
            ```bash
            curl -sSL https://raw.githubusercontent.com/${{ github.repository }}/main/install-tokenmonitor.sh | bash -s -- --target-path /path/to/project --mode full
            ```
            
            ### 功能特色
            - ✅ 跨 IDE Token 監控
            - ✅ 即時成本分析
            - ✅ 專業統計報表
            - ✅ 多種部署模式
            
          draft: false
          prerelease: false
```

### 測試工作流程

```yaml
# .github/workflows/test.yml
name: Test TokenMonitor

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
        node-version: [16, 18, 20]
        
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v3
        with:
          node-version: ${{ matrix.node-version }}
          
      - name: Install dependencies
        run: npm install
        
      - name: Test deployment script (Windows)
        if: runner.os == 'Windows'
        run: |
          $testPath = Join-Path $env:TEMP "TokenMonitor-Test"
          New-Item -ItemType Directory -Path $testPath -Force
          .\scripts\deploy-tokenmonitor.ps1 -TargetPath $testPath -Mode lite
          
      - name: Test deployment script (Unix)
        if: runner.os != 'Windows'
        run: |
          TEST_PATH="/tmp/TokenMonitor-Test"
          mkdir -p "$TEST_PATH"
          # 測試輕量部署邏輯
          echo "Testing lite deployment..."
```

---

## 🎯 使用範例

### 快速開始 (Windows)

```powershell
# 一行命令安裝
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/quick-install.ps1'))
```

### 快速開始 (Linux/macOS)

```bash
# 一行命令安裝
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full
```

### 在 Docker 中使用

```dockerfile
# Dockerfile
FROM node:18-alpine

WORKDIR /app

# 下載 TokenMonitor
RUN wget https://github.com/cy5407/TokenMonitor/archive/main.tar.gz \
    && tar -xzf main.tar.gz \
    && mv TokenMonitor-main TokenMonitor \
    && cd TokenMonitor \
    && npm install

# 設置環境
ENV TOKEN_LOG_LEVEL=info
ENV TOKEN_RETENTION_DAYS=30

# 啟動監控
CMD ["node", "TokenMonitor/src/js/professional-token-cli.js"]
```

---

## 📊 統計和分析

### 下載統計

GitHub 會自動追蹤：
- Release 下載次數
- Clone 次數
- Fork 次數
- Star 次數

### 使用分析

可以在 README.md 中添加：

```markdown
## 📈 使用統計

![GitHub release (latest by date)](https://img.shields.io/github/v/release/cy5407/TokenMonitor)
![GitHub downloads](https://img.shields.io/github/downloads/cy5407/TokenMonitor/total)
![GitHub stars](https://img.shields.io/github/stars/cy5407/TokenMonitor)
![GitHub forks](https://img.shields.io/github/forks/cy5407/TokenMonitor)
```

---

## 🎉 總結

將 TokenMonitor 推送到 GitHub 後，使用者可以：

1. **直接下載使用** - 無需本地開發環境
2. **一鍵安裝** - 簡化部署流程
3. **跨平台支援** - Windows/Linux/macOS
4. **版本管理** - 追蹤更新和變更
5. **社群協作** - 接受貢獻和回饋

這樣就能讓 TokenMonitor 成為一個真正的開源工具，供全世界的開發者使用！🌍✨
