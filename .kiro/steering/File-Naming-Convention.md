---
inclusion: always
---

# 檔案命名規範

## 基本原則

採用 **Go 語言風格的命名習慣**，強調簡潔、清晰和一致性：

### 檔案命名規則
- **全小寫字母**
- **使用底線 `_` 作為分隔符號**
- **避免使用連字號 `-`**
- **名稱應簡潔且具描述性**

### 變數和函數命名規則
- **camelCase** 用於私有/內部使用
- **PascalCase** 用於公開/匯出使用
- **全大寫 + 底線** 用於常數

## 檔案命名範例

### ✅ 正確範例（Go 風格）
```
token_monitor.js
user_profile_manager.ts
database_connection_helper.py
quick_setup_guide.md
api_response_handler.jsx
test_data_generator.sql
```

### ❌ 錯誤範例
```
Token-Monitor.js         # 不應使用連字號和大寫
tokenMonitor.js          # 不應使用 camelCase
USER_PROFILE_MANAGER.ts  # 不應全大寫
quick-setup-guide.md     # 不應使用連字號
apiResponseHandler.jsx   # 應使用底線分隔
testdatagenerator.sql    # 缺少分隔符號
```

## 程式碼內命名範例

### JavaScript/TypeScript
```javascript
// ✅ 正確：私有函數使用 camelCase
function calculateTokenUsage(text) { }
const userConfig = { };

// ✅ 正確：匯出函數使用 PascalCase
export function CalculateTokenUsage(text) { }
export const UserConfig = { };

// ✅ 正確：常數使用全大寫 + 底線
const MAX_TOKEN_LIMIT = 1000;
const DEFAULT_CONFIG_PATH = './config.yaml';
```

### PowerShell
```powershell
# ✅ 正確：函數使用 PascalCase（PowerShell 慣例）
function Get-TokenUsage { }
function Set-UserConfiguration { }

# ✅ 正確：變數使用 camelCase
$tokenCount = 0
$configPath = "./config.yaml"

# ✅ 正確：常數使用全大寫
$MAX_RETRIES = 3
$DEFAULT_TIMEOUT = 30
```

## 特殊情況

### 配置檔案
```
package.json
tsconfig.json
webpack_config.js
```

### 測試檔案
```
user_service.test.js
database_helper.spec.ts
integration_test_suite.py
```

### 文件檔案
```
README.md              # 特殊情況：保持傳統大寫
api_documentation.md
user_guide.md
installation_instructions.md
```

### 腳本檔案
```
deploy_application.ps1
build_project.sh
start_server.bat
```

## 資料夾命名

資料夾名稱遵循相同的 Go 風格規則：
```
src/
components/
utils/
test_data/
config_files/
documentation/
```

## Go 風格的程式碼組織原則

### 1. 簡潔性原則
- 名稱應該簡短但有意義
- 避免不必要的縮寫
- 優先使用常見的英文單詞

### 2. 一致性原則
- 整個專案使用統一的命名風格
- 相似功能的檔案使用相似的命名模式
- 保持命名的可預測性

### 3. 可讀性原則
- 名稱應該能清楚表達其用途
- 避免使用模糊或容易混淆的名稱
- 使用完整的英文單詞而非縮寫

## 實作提醒

在創建任何新檔案時，請確保：
1. **檔案名稱使用全小寫 + 底線**
2. **程式碼內命名遵循 Go 風格慣例**
3. **使用有意義的描述性名稱**
4. **避免不必要的縮寫**
5. **保持整個專案的命名一致性**

## Go 風格的優勢

採用 Go 語言的命名風格帶來以下好處：
- **簡潔明瞭** - 避免過度複雜的命名
- **一致性強** - 統一的命名規則減少混淆
- **可讀性高** - 清晰的命名提升程式碼理解度
- **國際化友好** - 全小寫檔案名在各系統間相容性更好
- **搜尋友好** - 底線分隔的檔案名更容易搜尋和過濾

## 遷移指導

對於現有檔案的重新命名：
1. **優先處理新建檔案** - 新檔案直接採用新規範
2. **逐步重構舊檔案** - 在修改時順便重新命名
3. **更新相關引用** - 確保所有引用都同步更新
4. **測試驗證** - 重新命名後進行功能測試

這個基於 Go 風格的命名規範將使專案更加專業和易於維護。