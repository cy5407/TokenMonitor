# 專案開發偏好設定

## 溝通風格要求

### 基本溝通原則
- **平等對話** - 我們是平等的協作夥伴，使用「你」而非「您」
- **中立客觀** - 以事實和邏輯為基礎進行討論，避免迎合或奉承
- **直接有效** - 直接表達觀點和建議，不需要過度禮貌用語
- **專業討論** - 專注於技術問題的分析和解決方案
- **坦誠回饋** - 如有不同意見或更好的方案，直接提出討論

## 任務執行原則

### 最小化動作原則
- **執行任務時請採用最小化動作**
- 不要創造使用者未明確提及的功能或檔案
- 專注於解決使用者明確要求的問題
- 若認為需要補充額外功能，**必須先與使用者討論確認**
- 避免過度工程化或添加不必要的複雜性

### 功能補充流程
1. 識別使用者明確要求的核心功能
2. 評估是否需要額外支援功能
3. 若需要補充功能，先詢問使用者意見
4. 獲得確認後再進行實作

## 測試檔案管理

### 測試檔案位置
- 所有測試相關檔案統一放置於 **根目錄下的 `/Tests` 資料夾**
- 測試檔案命名遵循專案命名規範

### 測試腳本要求
測試資料夾內必須包含：
- **測試檔案本身**
- **檔案移動腳本** (如需要)
- **檔案重新命名腳本** (如需要)
- **測試執行腳本**
- **測試清理腳本**

### 測試檔案結構範例
```
/Tests
├── Test-Feature-Name/
│   ├── Test-Cases.js
│   ├── Move-Files.ps1
│   ├── Rename-Files.ps1
│   ├── Run-Tests.ps1
│   └── Cleanup-Tests.ps1
```

## 任務完成報告

### 報告格式要求
- **每個任務完成後必須產生 Markdown 報告**
- 報告檔案命名格式：`YYYYMMDDHHMM-任務名稱.md`
- 報告內容應包含：
  - 任務概述
  - 執行步驟
  - 建立的檔案清單
  - 修改的檔案清單
  - 測試結果 (如適用)
  - 注意事項或後續建議

### 報告命名範例
```
202501011430-Token-Monitor-Setup.md
202501011445-Database-Migration.md
202501011500-API-Integration-Test.md
```

## 程式碼品質要求

### 程式碼撰寫原則（Go 風格）
- **遵循 Go 語言的命名慣例和程式碼風格**
- **簡潔性優先** - 避免過度複雜的實作
- **明確性** - 程式碼應該清楚表達其意圖
- **一致性** - 整個專案使用統一的風格
- **添加適當的註解說明**
- **確保程式碼可讀性和維護性**
- **避免硬編碼，使用配置檔案或環境變數**

### Go 風格的命名規範
- **檔案名稱**：全小寫 + 底線 (例：`token_monitor.js`)
- **函數/方法**：
  - 私有：camelCase (例：`calculateUsage`)
  - 公開：PascalCase (例：`CalculateUsage`)
- **變數**：
  - 私有：camelCase (例：`tokenCount`)
  - 公開：PascalCase (例：`TokenCount`)
- **常數**：全大寫 + 底線 (例：`MAX_TOKEN_LIMIT`)

### 錯誤處理（Go 風格）
- **明確的錯誤處理** - 不忽略任何錯誤
- **提供有意義的錯誤訊息**
- **使用包裝錯誤提供上下文**
- **記錄重要的操作日誌**
- **優雅的失敗處理**

## 文件撰寫規範

### 文件內容要求
- 使用清晰簡潔的語言
- 提供實際的使用範例
- 包含必要的安裝和設定說明
- 說明相依性和系統需求

### 文件結構
- 使用適當的標題層級
- 善用清單和表格整理資訊
- 提供程式碼區塊範例
- 包含目錄索引 (長文件)

## 溝通協作原則

### 與使用者互動
- 在執行重要變更前先確認
- 清楚說明執行的步驟和原因
- 主動回報進度和遇到的問題
- 提供替代方案供使用者選擇

### 問題解決流程
1. 理解使用者需求
2. 分析問題和限制條件
3. 提出解決方案
4. 獲得使用者確認
5. 執行並回報結果
## Go 
風格程式碼撰寫指導

### 函數設計原則
- **單一職責** - 每個函數只做一件事
- **短小精悍** - 函數應該簡潔明瞭
- **明確的參數和回傳值**
- **避免深層嵌套** - 使用早期返回減少嵌套

### 程式碼組織
- **邏輯分組** - 相關功能放在一起
- **依賴注入** - 避免硬編碼依賴
- **介面導向** - 定義清晰的介面契約
- **模組化設計** - 可重用的元件設計

### 註解風格
- **公開函數必須有註解**
- **註解應該說明 "為什麼" 而不只是 "做什麼"**
- **使用簡潔清晰的語言**
- **保持註解與程式碼同步**

### 範例對照

#### JavaScript 範例
```javascript
// ✅ Go 風格
// CalculateTokenUsage 計算指定文本的 token 使用量
export function CalculateTokenUsage(text, model = 'claude-sonnet-4.0') {
    if (!text) {
        return { error: 'text cannot be empty' };
    }
    
    const tokenCount = estimateTokens(text);
    const cost = calculateCost(tokenCount, model);
    
    return { tokenCount, cost };
}

// calculateCost 計算指定 token 數量的成本
function calculateCost(tokens, model) {
    const pricing = getPricing(model);
    return tokens * pricing.input / 1000000;
}
```

#### PowerShell 範例
```powershell
# ✅ Go 風格
# Get-TokenUsage 獲取 token 使用統計
function Get-TokenUsage {
    param(
        [string]$FilePath,
        [string]$Model = "claude-sonnet-4.0"
    )
    
    if (-not (Test-Path $FilePath)) {
        Write-Error "File not found: $FilePath"
        return
    }
    
    $content = Get-Content $FilePath -Raw
    $tokenCount = Measure-TokenCount $content
    $cost = Get-TokenCost $tokenCount $Model
    
    return @{
        TokenCount = $tokenCount
        Cost = $cost
        Model = $Model
    }
}
```