# TokenMonitor 專案文件

> 🌐 **現已支援 GitHub 部署！** 可從網路直接下載安裝

## 🚀 快速安裝 (從 GitHub)

### Windows PowerShell 一鍵安裝
```powershell
# 完整安裝
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/quick-install.ps1'))

# 或下載安裝腳本
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.ps1" -OutFile "install.ps1"
.\install.ps1 -TargetPath "C:\MyProject" -Mode full
```

### Linux/macOS 一鍵安裝
```bash
# 完整安裝
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full

# 輕量安裝
curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./token-monitor --mode lite
```

## 📁 專案結構

```
TokenMonitor/
├── 📂 docs/                    # 文件資料夾
│   ├── README.md               # 專案說明
│   ├── USAGE-GUIDE.md          # 使用指南
│   └── ARCHITECTURE.md         # 架構說明
├── 📂 scripts/                 # 腳本工具
│   ├── tokusage.ps1           # 主要 CLI 工具
│   ├── universal-monitor.ps1   # 通用監控腳本
│   └── legacy/                # 舊版腳本
├── 📂 tests/                   # 測試檔案
│   ├── reports/               # 測試報告
│   └── data/                  # 測試資料
├── 📂 src/                     # 原始碼
│   ├── js/                    # JavaScript 檔案
│   └── go/                    # Go 語言檔案
├── 📂 .kiro/                   # Kiro IDE 配置
├── 📂 data/                    # 資料檔案
└── 📂 build/                   # 編譯輸出
```

## 🚀 快速開始

1. **安裝依賴**
   ```bash
   npm install
   ```

2. **啟動監控**
   ```powershell
   .\scripts\tokusage.ps1 daily
   ```

3. **查看報告**
   ```powershell
   .\scripts\tokusage.ps1 summary
   ```

## 📊 主要功能

- ✅ 跨 IDE Token 使用監控
- ✅ 即時成本分析
- ✅ 專業統計報表
- ✅ 自動化監控系統
- ✅ 多模型支援

## 📖 詳細文件

- [使用指南](USAGE-GUIDE.md)
- [架構說明](ARCHITECTURE.md)
- [API 文件](API.md)
