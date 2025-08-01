# 🚀 TokenMonitor

> AI Token 使用監控和成本分析工具

[![GitHub release](https://img.shields.io/github/v/release/cy5407/TokenMonitor)](https://github.com/cy5407/TokenMonitor/releases)
[![GitHub stars](https://img.shields.io/github/stars/cy5407/TokenMonitor)](https://github.com/cy5407/TokenMonitor/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/cy5407/TokenMonitor)](https://github.com/cy5407/TokenMonitor/network)
[![License](https://img.shields.io/github/license/cy5407/TokenMonitor)](LICENSE)

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

\\\powershell
# 完整安裝
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/quick-install.ps1'))

# 自訂安裝
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.ps1" -OutFile "install.ps1"
.\install.ps1 -TargetPath "C:\MyProject" -Mode full
\\\

### Linux/macOS 一鍵安裝

\\\ash
# 完整安裝
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full

# 輕量安裝
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./token-monitor --mode lite
\\\

### 手動安裝

\\\ash
# 克隆倉庫
git clone https://github.com/cy5407/TokenMonitor.git
cd TokenMonitor

# 安裝依賴
npm install

# 執行部署
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\MyProject" -Mode full
\\\

## 📊 使用方式

### 查看每日報告

\\\powershell
# Windows
.\scripts\tokusage.ps1 daily

# Linux/macOS (如果安裝了 PowerShell)
pwsh ./scripts/tokusage.ps1 daily
\\\

### 查看詳細統計

\\\powershell
.\scripts\tokusage.ps1 summary
\\\

### 輕量版使用

\\\ash
# 查看報告
node token-monitor.js report

# 記錄使用
node token-monitor.js log chat_message 150 0.00045
\\\

## 📋 部署模式

| 模式 | 適用場景 | 檔案大小 | 功能完整度 |
|------|----------|----------|------------|
| **完整部署** | 需要全功能的專案 | ~2MB | 100% |
| **輕量部署** | 只需基本監控 | ~50KB | 60% |
| **NPM 套件** | Node.js 專案 | ~500KB | 90% |

## 🏗️ 專案結構

\\\
TokenMonitor/
├── 📂 scripts/          # 主要腳本工具
├── 📂 src/js/           # JavaScript 原始碼
├── 📂 src/go/           # Go 語言模組
├── 📂 docs/             # 完整文件
├── 📂 tests/            # 測試檔案
├── 📂 templates/        # 部署模板
└── 📂 .kiro/            # Kiro IDE 整合
\\\

## 📖 文件

- [📋 使用指南](docs/USAGE-GUIDE.md)
- [🏗️ 架構說明](docs/ARCHITECTURE.md)
- [🚀 部署指南](DEPLOYMENT-GUIDE.md)
- [🌐 GitHub 部署](GITHUB-DEPLOYMENT.md)
- [💡 部署範例](DEPLOYMENT-EXAMPLES.md)

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

1. Fork 這個專案
2. 創建你的功能分支 (\git checkout -b feature/AmazingFeature\)
3. 提交你的變更 (\git commit -m 'Add some AmazingFeature'\)
4. 推送到分支 (\git push origin feature/AmazingFeature\)
5. 開啟一個 Pull Request

## 📄 授權

這個專案使用 MIT 授權 - 查看 [LICENSE](LICENSE) 檔案了解詳情。

## 🌟 支援

如果這個專案對你有幫助，請給我們一個 ⭐！

## 📞 聯絡

- GitHub Issues: [https://github.com/cy5407/TokenMonitor/issues](https://github.com/cy5407/TokenMonitor/issues)
- 專案連結: [https://github.com/cy5407/TokenMonitor](https://github.com/cy5407/TokenMonitor)

---

**TokenMonitor** - 讓 AI 使用成本透明化 🚀
