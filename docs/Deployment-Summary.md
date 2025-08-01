# 🚀 TokenMonitor 部署方案總結

## 📋 部署選項概覽

TokenMonitor 提供了三種靈活的部署方式，滿足不同專案的需求：

| 部署方式 | 適用場景 | 檔案大小 | 功能完整度 | 設置複雜度 |
|---------|---------|---------|-----------|-----------|
| **完整部署** | 需要全功能的專案 | ~2MB | 100% | 中等 |
| **輕量部署** | 只需基本監控 | ~50KB | 60% | 簡單 |
| **NPM 套件** | Node.js 專案 | ~500KB | 90% | 簡單 |

---

## 🔧 快速部署指令

### 1. 完整部署 (推薦)
```powershell
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\YourProject" -Mode full
```

**包含功能**:
- ✅ 完整的 CLI 工具 (`tokusage.ps1`)
- ✅ 專業報表生成器
- ✅ Kiro IDE Hook 整合
- ✅ 多模型支援
- ✅ 成本分析
- ✅ 自動化監控

### 2. 輕量部署
```powershell
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\YourProject" -Mode lite
```

**包含功能**:
- ✅ 基本 Token 記錄
- ✅ 簡單報表生成
- ✅ 成本計算
- ❌ 進階分析功能
- ❌ Hook 整合

### 3. NPM 套件部署
```powershell
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\YourProject" -Mode npm
```

**包含功能**:
- ✅ 完整 CLI 工具
- ✅ 程式化 API
- ✅ 自動安裝腳本
- ✅ 跨平台支援
- ❌ PowerShell 腳本

---

## 📊 部署後驗證

### 完整部署驗證
```powershell
# 進入部署目錄
cd "C:\YourProject\TokenMonitor"

# 安裝依賴
npm install

# 測試 CLI
.\scripts\tokusage.ps1 daily

# 預期輸出: Token 使用報告表格
```

### 輕量部署驗證
```powershell
# 進入部署目錄
cd "C:\YourProject\token-monitor"

# 測試功能
node token-monitor.js report

# 預期輸出: 基本統計資訊
```

### NPM 套件驗證
```bash
# 進入套件目錄
cd "C:\YourProject\kiro-token-monitor"

# 安裝依賴
npm install

# 測試 CLI
node bin/tokusage.js --help

# 預期輸出: CLI 幫助資訊
```

---

## 🎯 使用場景建議

### 🏢 企業級專案
**推薦**: 完整部署
- 需要詳細的成本分析
- 多人協作開發
- 需要整合到 CI/CD

### 🚀 個人專案
**推薦**: 輕量部署
- 只需要基本監控
- 快速設置
- 最小化依賴

### 📦 開源專案
**推薦**: NPM 套件
- 易於分享和安裝
- 標準化部署
- 社群友好

### 🔬 實驗性專案
**推薦**: 輕量部署
- 快速原型驗證
- 臨時使用
- 最小化影響

---

## 📋 部署檢查清單

### 部署前準備
- [ ] 確認 Node.js 已安裝 (v14+)
- [ ] 確認 PowerShell 執行權限
- [ ] 備份現有配置檔案
- [ ] 確認目標路徑存在

### 部署執行
- [ ] 選擇合適的部署模式
- [ ] 執行部署腳本
- [ ] 檢查檔案複製完整性
- [ ] 安裝必要依賴

### 部署後驗證
- [ ] 測試 CLI 工具運行
- [ ] 檢查記錄檔案生成
- [ ] 驗證報表生成功能
- [ ] 測試成本計算準確性

### 整合測試
- [ ] 在實際專案中測試
- [ ] 驗證 Hook 系統運作
- [ ] 檢查自動監控功能
- [ ] 確認報表格式正確

---

## 🔧 客製化選項

### 配置檔案位置
```
完整部署: TokenMonitor/.kiro/settings/
輕量部署: token-monitor/config.json
NPM 套件: .kiro/settings/token-monitor.json
```

### 記錄檔案位置
```
完整部署: TokenMonitor/data/kiro-usage.log
輕量部署: token-monitor/token-usage.log
NPM 套件: data/kiro-usage.log
```

### CLI 工具位置
```
完整部署: TokenMonitor/scripts/tokusage.ps1
輕量部署: token-monitor/token-monitor.js
NPM 套件: node_modules/.bin/tokusage
```

---

## 🚨 故障排除

### 常見問題

#### 1. 權限錯誤
```powershell
# 解決方案
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### 2. 找不到 Node.js
```bash
# 檢查 Node.js 安裝
node --version
npm --version

# 如未安裝，請從 nodejs.org 下載安裝
```

#### 3. 檔案路徑錯誤
```powershell
# 檢查路徑是否存在
Test-Path "C:\YourProject"

# 使用絕對路徑
.\scripts\deploy-tokenmonitor.ps1 -TargetPath "C:\Full\Path\To\Project"
```

#### 4. 依賴安裝失敗
```bash
# 清理 npm 快取
npm cache clean --force

# 重新安裝
npm install
```

### 支援資源

- 📖 **完整文件**: `docs/README.md`
- 🏗️ **架構說明**: `docs/ARCHITECTURE.md`
- 📋 **使用指南**: `docs/USAGE-GUIDE.md`
- 🚀 **部署指南**: `DEPLOYMENT-GUIDE.md`
- 💡 **範例集合**: `DEPLOYMENT-EXAMPLES.md`

---

## 🎉 部署成功！

恭喜您成功部署了 TokenMonitor！現在您可以：

1. **監控 AI Token 使用情況**
2. **分析成本和效率**
3. **生成專業報表**
4. **自動化監控流程**

### 下一步建議

1. 設置定期報表生成
2. 配置成本警報
3. 整合到開發工作流程
4. 與團隊分享使用方法

**祝您使用愉快！** 🚀✨