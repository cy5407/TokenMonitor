# Kiro IDE 整合測試

這是一個測試檔案，用來驗證 Kiro IDE 的 Token 監控是否正常工作。

## 測試步驟

1. 建立這個檔案
2. 與 Kiro AI 進行對話
3. 檢查 `data/kiro-usage.log` 是否有新記錄

測試時間: 2025-08-01

## 預期結果

如果整合正常，應該會在日誌中看到：
- `chat_message` 事件
- 正確的 token 計算
- 實際的成本分析

如果只看到 `command_executed` 事件，表示只有終端監控在工作，而不是真正的 Kiro IDE 整合。