# Requirements Document

## Introduction

建立一個簡單的通用 AI Token 監控系統，能夠同時支援 VS Code 和 Kiro IDE，透過檔案監控方式捕捉 AI 生成的內容和 Token 使用量。系統設計簡潔，專注於核心功能。

## Requirements

### Requirement 1

**User Story:** 作為開發者，我希望能夠監控我在 VS Code 和 Kiro IDE 中的 AI 使用量，以便了解我的 Token 消耗。

#### Acceptance Criteria

1. WHEN 我在 VS Code 中使用 AI 助手產生內容 THEN 系統 SHALL 透過檔案變化檢測並記錄
2. WHEN 我在 Kiro IDE 中與 AI 對話產生內容 THEN 系統 SHALL 同樣記錄相關數據
3. WHEN 系統檢測到 AI 相關檔案變化 THEN 系統 SHALL 在 3 秒內記錄到日誌檔案
4. WHEN 我切換不同的 IDE THEN 監控系統 SHALL 持續運作

### Requirement 2

**User Story:** 作為開發者，我希望系統能夠識別 AI 生成的內容，並準確計算 Token 使用量。

#### Acceptance Criteria

1. WHEN 系統檢測到檔案內容包含程式碼模式 THEN 系統 SHALL 識別為 AI 生成內容
2. WHEN 系統檢測到檔案內容包含 Markdown 格式 THEN 系統 SHALL 識別為 AI 生成內容
3. WHEN 系統識別到 AI 內容 THEN 系統 SHALL 計算對應的 Token 數量
4. WHEN 內容包含中英文混合文本 THEN 系統 SHALL 根據語言特性計算 Token
5. WHEN 系統計算 Token THEN 系統 SHALL 估算對應的成本

### Requirement 3

**User Story:** 作為開發者，我希望能夠區分不同 IDE 的 AI 活動，並進行基本的活動分類。

#### Acceptance Criteria

1. WHEN 系統檢測到 VS Code 相關的檔案變化 THEN 系統 SHALL 標記為 "vscode" 來源
2. WHEN 系統檢測到 Kiro IDE 相關的檔案變化 THEN 系統 SHALL 標記為 "kiro" 來源
3. WHEN 系統分析內容類型 THEN 系統 SHALL 分類為 "coding", "documentation", 或 "chat"
4. WHEN 無法確定 IDE 類型 THEN 系統 SHALL 標記為 "unknown" 來源

### Requirement 4

**User Story:** 作為開發者，我希望能夠查看監控狀態和生成使用報告。

#### Acceptance Criteria

1. WHEN 我執行狀態查詢命令 THEN 系統 SHALL 顯示當前監控狀態
2. WHEN 我執行報告生成命令 THEN 系統 SHALL 生成基本的使用統計
3. WHEN 系統運行時 THEN 系統 SHALL 在控制台顯示檢測到的活動
4. WHEN 我需要停止監控 THEN 系統 SHALL 能夠正常停止

### Requirement 5

**User Story:** 作為開發者，我希望監控系統輕量且穩定，對系統性能影響最小。

#### Acceptance Criteria

1. WHEN 系統運行時 THEN 系統 SHALL 只監控必要的檔案類型
2. WHEN 系統檢測到錯誤 THEN 系統 SHALL 記錄錯誤但繼續運行
3. WHEN 系統運行時 THEN 系統 SHALL 避免過度的資源消耗
4. WHEN 日誌檔案過大 THEN 系統 SHALL 提供基本的日誌管理