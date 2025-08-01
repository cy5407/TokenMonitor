#!/usr/bin/env pwsh

<#
.SYNOPSIS
    準備 TokenMonitor 專案推送到 GitHub

.DESCRIPTION
    這個腳本會幫助你準備 TokenMonitor 專案，包括創建 .gitignore、README.md 等必要檔案

.PARAMETER GitHubUsername
    你的 GitHub 使用者名稱

.PARAMETER RepoName
    倉庫名稱 (預設: TokenMonitor)

.EXAMPLE
    .\prepare-github.ps1 -GitHubUsername "yourusername"
#>

param(
    [Parameter(Mandatory=$true, HelpMessage="GitHub 使用者名稱")]
    [string]$GitHubUsername,
    
    [Parameter(HelpMessage="倉庫名稱")]
    [string]$RepoName = "TokenMonitor",
    
    [Parameter(HelpMessage="顯示幫助")]
    [switch]$Help
)

if ($Help) {
    Write-Host @"
🚀 TokenMonitor GitHub 準備工具

用法:
    prepare-github.ps1 -GitHubUsername <使用者名稱> [選項]

參數:
    -GitHubUsername    GitHub 使用者名稱 (必要)
    -RepoName         倉庫名稱 (預設: TokenMonitor)
    -Help             顯示此幫助

範例:
    .\prepare-github.ps1 -GitHubUsername "myusername"
    .\prepare-github.ps1 -GitHubUsername "myusername" -RepoName "MyTokenMonitor"

"@ -ForegroundColor Cyan
    exit 0
}

Write-Host "🚀 準備 TokenMonitor 專案推送到 GitHub" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host "👤 GitHub 使用者: $GitHubUsername" -ForegroundColor Gray
Write-Host "📦 倉庫名稱: $RepoName" -ForegroundColor Gray
Write-Host ""

try {
    # 1. 更新安裝腳本中的 GitHub 倉庫路徑
    Write-Host "🔧 更新安裝腳本中的 GitHub 路徑..." -ForegroundColor Yellow
    
    $repoPath = "$GitHubUsername/$RepoName"
    
    # 更新 PowerShell 安裝腳本
    if (Test-Path "install-tokenmonitor.ps1") {
        $content = Get-Content "install-tokenmonitor.ps1" -Raw
        $content = $content -replace 'yourusername/TokenMonitor', $repoPath
        Set-Content "install-tokenmonitor.ps1" -Value $content -Encoding UTF8
        Write-Host "✅ 更新: install-tokenmonitor.ps1" -ForegroundColor Green
    }
    
    # 更新 Bash 安裝腳本
    if (Test-Path "install-tokenmonitor.sh") {
        $content = Get-Content "install-tokenmonitor.sh" -Raw
        $content = $content -replace 'yourusername/TokenMonitor', $repoPath
        Set-Content "install-tokenmonitor.sh" -Value $content -Encoding UTF8
        Write-Host "✅ 更新: install-tokenmonitor.sh" -ForegroundColor Green
    }
    
    # 更新快速安裝腳本
    if (Test-Path "quick-install.ps1") {
        $content = Get-Content "quick-install.ps1" -Raw
        $content = $content -replace 'yourusername/TokenMonitor', $repoPath
        Set-Content "quick-install.ps1" -Value $content -Encoding UTF8
        Write-Host "✅ 更新: quick-install.ps1" -ForegroundColor Green
    }
    
    # 更新文件中的 GitHub 連結
    $docsToUpdate = @("docs/README.md", "GITHUB-DEPLOYMENT.md", "DEPLOYMENT-GUIDE.md")
    foreach ($doc in $docsToUpdate) {
        if (Test-Path $doc) {
            $content = Get-Content $doc -Raw
            $content = $content -replace 'yourusername/TokenMonitor', $repoPath
            Set-Content $doc -Value $content -Encoding UTF8
            Write-Host "✅ 更新: $doc" -ForegroundColor Green
        }
    }

    # 2. 創建 .gitignore
    Write-Host "📝 創建 .gitignore..." -ForegroundColor Yellow
    
    $gitignoreContent = @"
# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# 使用記錄檔案
data/kiro-usage.log
*/data/kiro-usage.log
token-usage.log
*/token-usage.log

# 臨時檔案
*.tmp
*.temp
.DS_Store
Thumbs.db

# IDE 檔案
.vscode/
.idea/
*.swp
*.swo

# 測試輸出
test-output/
coverage/

# 編譯輸出
build/
dist/

# 環境變數
.env
.env.local

# PowerShell 執行記錄
*.ps1.log

# 備份檔案
*.bak
*.backup

# 系統檔案
.DS_Store
desktop.ini
"@
    
    Set-Content ".gitignore" -Value $gitignoreContent -Encoding UTF8
    Write-Host "✅ 創建: .gitignore" -ForegroundColor Green

    # 3. 創建主要 README.md
    Write-Host "📝 創建主要 README.md..." -ForegroundColor Yellow
    
    $readmeContent = @"
# 🚀 TokenMonitor

> AI Token 使用監控和成本分析工具

[![GitHub release](https://img.shields.io/github/v/release/$GitHubUsername/$RepoName)](https://github.com/$GitHubUsername/$RepoName/releases)
[![GitHub stars](https://img.shields.io/github/stars/$GitHubUsername/$RepoName)](https://github.com/$GitHubUsername/$RepoName/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/$GitHubUsername/$RepoName)](https://github.com/$GitHubUsername/$RepoName/network)
[![License](https://img.shields.io/github/license/$GitHubUsername/$RepoName)](LICENSE)

TokenMonitor 是一個專業的 AI Token 使用監控系統，支援跨 IDE 監控、即時成本分析和專業統計報表。

## ✨ 功能特色

- 🔍 **跨 IDE 監控** - 支援 Kiro IDE、VS Code 等多種開發環境
- 💰 **即時成本分析** - 精確計算 AI 使用成本
- 📊 **專業報表** - 類似 ccusage 的專業統計介面
- 🚀 **多種部署模式** - 完整版、輕量版、NPM 套件
- 🌍 **跨平台支援** - Windows、Linux、macOS
- ⚡ **一鍵安裝** - 從 GitHub 直接下載部署

## 🚀 快速開始

### Windows PowerShell 一鍵安裝

\`\`\`powershell
# 完整安裝
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/quick-install.ps1'))

# 自訂安裝
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/install-tokenmonitor.ps1" -OutFile "install.ps1"
.\install.ps1 -TargetPath "C:\MyProject" -Mode full
\`\`\`

### Linux/macOS 一鍵安裝

\`\`\`bash
# 完整安裝
curl -sSL https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full

# 輕量安裝
curl -sSL https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/install-tokenmonitor.sh | bash -s -- --target-path ./token-monitor --mode lite
\`\`\`

### 手動安裝

\`\`\`bash
# 克隆倉庫
git clone https://github.com/$GitHubUsername/$RepoName.git
cd $RepoName

# 安裝依賴
npm install

# 執行部署
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
\`\`\`

## 📊 使用方式

### 查看每日報告

\`\`\`powershell
# Windows
.\scripts\tokusage.ps1 daily

# Linux/macOS (如果安裝了 PowerShell)
pwsh ./scripts/tokusage.ps1 daily
\`\`\`

### 查看詳細統計

\`\`\`powershell
.\scripts\tokusage.ps1 summary
\`\`\`

### 輕量版使用

\`\`\`bash
# 查看報告
node token-monitor.js report

# 記錄使用
node token-monitor.js log chat_message 150 0.00045
\`\`\`

## 📋 部署模式

| 模式 | 適用場景 | 檔案大小 | 功能完整度 |
|------|----------|----------|------------|
| **完整部署** | 需要全功能的專案 | ~2MB | 100% |
| **輕量部署** | 只需基本監控 | ~50KB | 60% |
| **NPM 套件** | Node.js 專案 | ~500KB | 90% |

## 🏗️ 專案結構

\`\`\`
TokenMonitor/
├── 📂 scripts/          # 主要腳本工具
├── 📂 src/js/           # JavaScript 原始碼
├── 📂 src/go/           # Go 語言模組
├── 📂 docs/             # 完整文件
├── 📂 tests/            # 測試檔案
├── 📂 templates/        # 部署模板
└── 📂 .kiro/            # Kiro IDE 整合
\`\`\`

## 📖 文件

- [📋 使用指南](docs/USAGE-GUIDE.md)
- [🏗️ 架構說明](docs/ARCHITECTURE.md)
- [🚀 部署指南](DEPLOYMENT-GUIDE.md)
- [🌐 GitHub 部署](GITHUB-DEPLOYMENT.md)
- [💡 部署範例](DEPLOYMENT-EXAMPLES.md)

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

1. Fork 這個專案
2. 創建你的功能分支 (\`git checkout -b feature/AmazingFeature\`)
3. 提交你的變更 (\`git commit -m 'Add some AmazingFeature'\`)
4. 推送到分支 (\`git push origin feature/AmazingFeature\`)
5. 開啟一個 Pull Request

## 📄 授權

這個專案使用 MIT 授權 - 查看 [LICENSE](LICENSE) 檔案了解詳情。

## 🌟 支援

如果這個專案對你有幫助，請給我們一個 ⭐！

## 📞 聯絡

- GitHub Issues: [https://github.com/$GitHubUsername/$RepoName/issues](https://github.com/$GitHubUsername/$RepoName/issues)
- 專案連結: [https://github.com/$GitHubUsername/$RepoName](https://github.com/$GitHubUsername/$RepoName)

---

**TokenMonitor** - 讓 AI 使用成本透明化 🚀
"@
    
    Set-Content "README.md" -Value $readmeContent -Encoding UTF8
    Write-Host "✅ 創建: README.md" -ForegroundColor Green

    # 4. 創建 LICENSE
    Write-Host "📝 創建 LICENSE..." -ForegroundColor Yellow
    
    $licenseContent = @"
MIT License

Copyright (c) $(Get-Date -Format yyyy) $GitHubUsername

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"@
    
    Set-Content "LICENSE" -Value $licenseContent -Encoding UTF8
    Write-Host "✅ 創建: LICENSE" -ForegroundColor Green

    # 5. 創建 GitHub Actions 工作流程
    Write-Host "📝 創建 GitHub Actions..." -ForegroundColor Yellow
    
    if (-not (Test-Path ".github/workflows")) {
        New-Item -ItemType Directory -Path ".github/workflows" -Force | Out-Null
    }
    
    # 使用 Here-String 避免 PowerShell 變數替換問題
    $workflowContent = @'
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
        run: npm install --production
        
      - name: Create Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ github.ref }}
          release_name: TokenMonitor ${{ github.ref }}
          body: |
            ## 🚀 TokenMonitor Release
            
            ### 快速安裝
            
            **Windows PowerShell:**
            ```powershell
            iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/GITHUB_USERNAME/REPO_NAME/main/quick-install.ps1'))
            ```
            
            **Linux/macOS:**
            ```bash
            curl -sSL https://raw.githubusercontent.com/GITHUB_USERNAME/REPO_NAME/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full
            ```
            
            ### 功能特色
            - ✅ 跨 IDE Token 監控
            - ✅ 即時成本分析
            - ✅ 專業統計報表
            - ✅ 多種部署模式
            
          draft: false
          prerelease: false
'@
    
    # 替換佔位符
    $workflowContent = $workflowContent -replace 'GITHUB_USERNAME', $GitHubUsername
    $workflowContent = $workflowContent -replace 'REPO_NAME', $RepoName
    
    Set-Content ".github/workflows/release.yml" -Value $workflowContent -Encoding UTF8
    Write-Host "✅ 創建: .github/workflows/release.yml" -ForegroundColor Green

    # 6. 顯示 Git 命令
    Write-Host ""
    Write-Host "🎉 GitHub 準備完成！" -ForegroundColor Green
    Write-Host ""
    Write-Host "📋 接下來的步驟:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "1. 初始化 Git 倉庫 (如果還沒有):" -ForegroundColor Yellow
    Write-Host "   git init" -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. 添加所有檔案:" -ForegroundColor Yellow
    Write-Host "   git add ." -ForegroundColor Gray
    Write-Host ""
    Write-Host "3. 提交變更:" -ForegroundColor Yellow
    Write-Host "   git commit -m `"Initial commit: TokenMonitor v1.0`"" -ForegroundColor Gray
    Write-Host ""
    Write-Host "4. 添加遠端倉庫:" -ForegroundColor Yellow
    Write-Host "   git remote add origin https://github.com/$GitHubUsername/$RepoName.git" -ForegroundColor Gray
    Write-Host ""
    Write-Host "5. 推送到 GitHub:" -ForegroundColor Yellow
    Write-Host "   git branch -M main" -ForegroundColor Gray
    Write-Host "   git push -u origin main" -ForegroundColor Gray
    Write-Host ""
    Write-Host "6. 創建第一個版本標籤:" -ForegroundColor Yellow
    Write-Host "   git tag -a v1.0.0 -m `"TokenMonitor v1.0.0 - Initial Release`"" -ForegroundColor Gray
    Write-Host "   git push origin v1.0.0" -ForegroundColor Gray
    Write-Host ""
    Write-Host "🌟 推送完成後，使用者就可以用以下命令安裝:" -ForegroundColor Green
    Write-Host ""
    Write-Host "Windows:" -ForegroundColor Cyan
    Write-Host "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/quick-install.ps1'))" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Linux/macOS:" -ForegroundColor Cyan
    Write-Host "curl -sSL https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full" -ForegroundColor Blue
    Write-Host ""

} catch {
    Write-Error "準備過程中發生錯誤: $($_.Exception.Message)"
    exit 1
}

Write-Host "✨ 準備完成！現在可以推送到 GitHub 了！" -ForegroundColor Green