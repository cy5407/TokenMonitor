# TokenMonitor 使用指南

> 🎉 **專案已整理完成！** 所有檔案已重新組織到合適的資料夾中

## 📁 新的專案結構

```
TokenMonitor/
├── 📂 scripts/          # 主要腳本工具
│   ├── tokusage.ps1     # 主要 CLI 工具
│   ├── universal-monitor.ps1
│   └── legacy/          # 舊版檔案
├── 📂 src/              # 原始碼
│   ├── js/              # JavaScript 模組
│   └── go/              # Go 語言模組
├── 📂 tests/            # 測試檔案
│   ├── reports/         # 測試報告
│   └── data/            # 測試資料
├── 📂 docs/             # 文件
├── 📂 data/             # 使用記錄
└── 📂 build/            # 編譯輸出
```

## 🚀 快速開始

### 1. 使用主要 CLI 工具
```powershell
# 查看每日報告
.\scripts\tokusage.ps1 daily

# 查看詳細統計
.\scripts\tokusage.ps1 summary

# 查看幫助
.\scripts\tokusage.ps1 --help
```

## 概述

這個改進的 Token 監控系統能夠全面追蹤你與 Kiro 的所有互動，包括：

1. **聊天對話** - 你的輸入和 Kiro 的回應
2. **工具調用** - Kiro 執行任務時使用的工具（fsWrite, strReplace 等）
3. **程式碼生成** - Kiro 生成的程式碼內容
4. **文件創建** - Kiro 創建的 Markdown 文件和其他文件
5. **代理任務** - Kiro 執行的各種自動化任務

## 系統組件

### 1. 主要監控整合腳本
- **檔案**: `token-monitor-integration.js`
- **功能**: 處理所有類型的事件並計算 Token 消耗
- **支援事件**:
  - 聊天訊息 (`chat.message.sent`, `chat.message.received`)
  - 工具執行 (`tool.fsWrite`, `tool.fsAppend`, `tool.strReplace`)
  - 代理任務 (`agent.codeGeneration`, `agent.documentGeneration`)
  - 檔案操作 (`file.saved`, `file.created`, `file.modified`)
  - 對話回合 (`kiro.conversation.turn`)

### 2. 全面監控器
- **檔案**: `.kiro/hooks/comprehensive-token-monitor.js`
- **功能**: 自動監控所有 Kiro 活動
- **特點**:
  - 即時監控和記錄
  - 自動成本計算
  - 定期分析（每5分鐘）
  - 統計資訊顯示

### 3. 手動分析器
- **檔案**: `.kiro/hooks/manual-token-calc.js`
- **功能**: 手動觸發的詳細分析
- **改進**:
  - 支援更多事件類型
  - 增強的統計分類
  - 詳細的成本分析

## 使用方法

### 自動監控
系統會自動監控所有活動，無需手動操作。監控資料會記錄到 `data/kiro-usage.log`。

### 手動分析
1. 在 Kiro 的 Agent Hooks 面板中找到 "Kiro Chat Token Calculator"
2. 點擊執行按鈕
3. 查看詳細的分析報告

### 查看即時統計
監控系統會在控制台顯示即時統計：
```
📊 [14:30:25] ai_code_generation: 45 tokens (0.000675 USD)
   🔧 工具: tool.fsWrite - write
   📁 檔案: fibonacci.js
   💰 成本: 0.000675 USD
```

## 監控的活動類型

### 聊天相關
- `chat` - 一般聊天對話
- `coding` - 程式設計相關討論
- `debugging` - 除錯和問題解決
- `documentation` - 文件撰寫

### AI 生成內容
- `ai_code_generation` - AI 生成程式碼
- `ai_documentation` - AI 生成文件
- `ai_debugging` - AI 協助除錯
- `ai_refactoring` - AI 協助重構
- `ai_testing` - AI 協助測試

### 工具操作
- `tool_execution` - 工具執行
- `file_operations` - 檔案操作
- `configuration` - 配置檔案修改

## 成本計算

系統使用 Claude Sonnet 4.0 的定價：
- **輸入 Token**: $3.00 / 1M tokens
- **輸出 Token**: $15.00 / 1M tokens

## 分析報告內容

### 基本統計
- 總記錄數
- 輸入/輸出 Token 數量
- 總成本（USD）

### 活動類型統計
- 各類型活動的次數和 Token 消耗
- 成本分析

### 會話統計
- 各會話的 Token 使用情況
- 最活躍的會話

### 模型使用統計
- 不同模型的使用情況
- 各模型的成本分析

### 事件類型統計
- 各種事件的發生頻率
- Token 消耗分布

## 測試系統

運行測試腳本來驗證監控系統：
```bash
node test-token-monitoring.js
```

測試包括：
- 聊天對話監控
- 工具調用監控
- 代理任務監控
- 檔案操作監控
- 對話回合監控
- 手動分析器

## 配置選項

### 全面監控器設定 (`.kiro/hooks/comprehensive-token-monitor.json`)
```json
{
  "settings": {
    "logFile": "data/kiro-usage.log",
    "realtimeDisplay": true,
    "costTracking": true,
    "sessionTracking": true,
    "detailedLogging": true,
    "autoAnalysis": true,
    "analysisInterval": 300000
  }
}
```

### 手動分析器設定 (`.kiro/hooks/manual-token-calc.json`)
```json
{
  "settings": {
    "logFile": "data/kiro-usage.log",
    "outputFormat": "detailed",
    "includeCostAnalysis": true,
    "showActivityBreakdown": true
  }
}
```

## 故障排除

### 常見問題

1. **監控沒有記錄資料**
   - 檢查 `data` 資料夾是否存在
   - 確認 Hook 是否已啟用
   - 查看控制台錯誤訊息

2. **Token 計算不準確**
   - 確認 `token-monitor.exe` 是否存在
   - 檢查網路連線（如果使用線上 API）
   - 查看估算函數是否正常運作

3. **分析報告為空**
   - 確認 `kiro-usage.log` 檔案存在且有內容
   - 檢查檔案權限
   - 確認日誌格式正確

### 日誌檔案格式
每筆記錄都是一個 JSON 物件：
```json
{
  "timestamp": "2025-07-31T05:57:38.098Z",
  "event": "chat_message",
  "direction": "sent",
  "content_length": 35,
  "tokens": 32,
  "activity_type": "general",
  "model": "claude-sonnet-4.0",
  "session_id": "session-1753941459410-d8s3ud",
  "cost_analysis": {
    "tokens": 32,
    "cost_usd": 0.000096,
    "cost_type": "input",
    "model": "claude-sonnet-4.0",
    "pricing_rate": 3
  }
}
```

## 進階功能

### 自訂活動分類
你可以修改 `classifyActivity` 和 `classifyAgentActivity` 函數來自訂活動分類邏輯。

### 成本警告
當總成本超過設定閾值時，系統會自動發出警告。

### 定期報告
系統每5分鐘會自動運行一次分析，提供使用情況摘要。

## 總結

這個改進的監控系統現在能夠：

✅ **完整監控你的輸入和 Kiro 的輸出**
✅ **追蹤所有工具調用和程式碼生成**
✅ **提供詳細的成本分析**
✅ **支援多種事件類型**
✅ **即時顯示統計資訊**
✅ **自動化分析和報告**

你現在可以全面了解與 Kiro 互動的 Token 消耗情況，包括聊天對話和所有任務執行過程中產生的內容。