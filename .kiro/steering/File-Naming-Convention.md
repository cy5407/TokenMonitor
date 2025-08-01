---
inclusion: always
---

# 檔案命名規範

## 基本原則

所有檔案命名必須遵循以下規則：

### 英文字母大小寫規則
- **首個英文字母大寫**
- **分隔符號後的首個英文字母大寫**
- **其餘字母小寫**

### 支援的分隔符號
- 連字號：`-`
- 底線：`_`
- 點號：`.` (僅用於副檔名前)

## 命名範例

### ✅ 正確範例
```
Token-Monitor.js
User-Profile-Manager.ts
Database-Connection_Helper.py
Quick-Setup-Guide.md
Api-Response-Handler.jsx
Test-Data_Generator.sql
```

### ❌ 錯誤範例
```
tokenMonitor.js          # 首字母應大寫
user-profile-manager.ts  # 分隔符號後應大寫
DATABASE-CONNECTION.py   # 不應全大寫
quick_setup_guide.md     # 分隔符號後應大寫
apiResponseHandler.jsx   # 應使用分隔符號
testdatagenerator.sql    # 缺少分隔符號和大寫
```

## 特殊情況

### 配置檔案
```
Package.json
Tsconfig.json
Webpack-Config.js
```

### 測試檔案
```
User-Service.test.js
Database-Helper.spec.ts
Integration-Test_Suite.py
```

### 文件檔案
```
README.md
API-Documentation.md
User-Guide.md
Installation-Instructions.md
```

### 腳本檔案
```
Deploy-Application.ps1
Build-Project.sh
Start-Server.bat
```

## 資料夾命名

資料夾名稱也遵循相同規則：
```
Src/
Components/
Utils/
Test-Data/
Config-Files/
Documentation/
```

## 實作提醒

在創建任何新檔案時，請確保：
1. 檢查檔案名稱是否符合命名規範
2. 使用有意義的描述性名稱
3. 避免使用縮寫，除非是廣泛認知的縮寫
4. 保持名稱簡潔但具描述性

這個命名規範有助於：
- 提高程式碼可讀性
- 保持專案結構一致性
- 便於檔案搜尋和管理
- 提升團隊協作效率