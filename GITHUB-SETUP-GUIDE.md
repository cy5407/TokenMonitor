# 🌐 GitHub 儲存庫建立完整指南

## 📋 步驟總覽

1. **在 GitHub 建立儲存庫** (手動)
2. **準備本地檔案** (腳本協助)
3. **推送到 GitHub** (Git 命令)
4. **測試安裝** (驗證)

---

## 🏗️ 步驟一：在 GitHub 建立儲存庫

### 1. 登入 GitHub
前往 [GitHub.com](https://github.com) 並登入你的帳號

### 2. 建立新儲存庫
1. 點擊右上角的 **"+"** 按鈕
2. 選擇 **"New repository"**
3. 填寫儲存庫資訊：
   - **Repository name**: `TokenMonitor` (或你想要的名稱)
   - **Description**: `AI Token usage monitoring and cost analysis tool`
   - **Visibility**: 選擇 `Public` (讓其他人可以使用)
   - **Initialize**: ❌ **不要勾選** "Add a README file"
   - **Initialize**: ❌ **不要勾選** "Add .gitignore"
   - **Initialize**: ❌ **不要勾選** "Choose a license"

### 3. 建立儲存庫
點擊 **"Create repository"** 按鈕

### 4. 記錄儲存庫資訊
建立完成後，你會看到類似這樣的頁面：
```
https://github.com/你的用戶名/TokenMonitor
```

記下你的：
- **GitHub 用戶名**: `你的用戶名`
- **儲存庫名稱**: `TokenMonitor`

---

## 🔧 步驟二：準備本地檔案

### 1. 執行準備腳本
```powershell
# 替換成你的實際 GitHub 用戶名
.\prepare-github.ps1 -GitHubUsername "你的GitHub用戶名"
```

**範例**:
```powershell
.\prepare-github.ps1 -GitHubUsername "john123"
```

### 2. 腳本會自動完成
- ✅ 更新所有安裝腳本中的 GitHub 路徑
- ✅ 創建 `.gitignore` 檔案
- ✅ 創建專業的 `README.md`
- ✅ 創建 `LICENSE` 檔案
- ✅ 設置 GitHub Actions 工作流程

---

## 📤 步驟三：推送到 GitHub

### 1. 初始化 Git 倉庫
```bash
git init
```

### 2. 設定 Git 使用者資訊 (如果還沒設定)
```bash
git config --global user.name "你的姓名"
git config --global user.email "你的email@example.com"
```

### 3. 添加所有檔案
```bash
git add .
```

### 4. 提交變更
```bash
git commit -m "Initial commit: TokenMonitor v1.0"
```

### 5. 添加遠端儲存庫
```bash
# 替換成你的實際 GitHub 路徑
git remote add origin https://github.com/你的用戶名/TokenMonitor.git
```

**範例**:
```bash
git remote add origin https://github.com/john123/TokenMonitor.git
```

### 6. 推送到 GitHub
```bash
git branch -M main
git push -u origin main
```

### 7. 創建版本標籤 (可選)
```bash
git tag -a v1.0.0 -m "TokenMonitor v1.0.0 - Initial Release"
git push origin v1.0.0
```

---

## ✅ 步驟四：驗證部署

### 1. 檢查 GitHub 頁面
前往你的 GitHub 儲存庫頁面，確認：
- ✅ 所有檔案都已上傳
- ✅ README.md 顯示正常
- ✅ 有 LICENSE 檔案
- ✅ 有 .github/workflows/ 目錄

### 2. 測試安裝腳本
在另一台電腦或新目錄中測試：

**Windows 測試**:
```powershell
# 替換成你的實際 GitHub 路徑
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/你的用戶名/TokenMonitor/main/install-tokenmonitor.ps1" -OutFile "test-install.ps1"
.\test-install.ps1 -TargetPath "C:\TestInstall" -Mode lite
```

**Linux/macOS 測試**:
```bash
# 替換成你的實際 GitHub 路徑
curl -sSL https://raw.githubusercontent.com/你的用戶名/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./test-install --mode lite
```

### 3. 驗證安裝結果
檢查安裝是否成功：
- ✅ 檔案正確複製
- ✅ 腳本可以執行
- ✅ 沒有錯誤訊息

---

## 🚨 常見問題和解決方案

### Q1: 推送時出現 "Permission denied" 錯誤
**解決方案**:
1. 檢查 GitHub 用戶名和儲存庫名稱是否正確
2. 確認你有該儲存庫的寫入權限
3. 可能需要設定 SSH 金鑰或使用 Personal Access Token

### Q2: 安裝腳本下載失敗
**解決方案**:
1. 確認儲存庫是 Public
2. 檢查檔案路徑是否正確
3. 等待幾分鐘讓 GitHub 同步

### Q3: 腳本中的路徑沒有更新
**解決方案**:
```powershell
# 重新執行準備腳本
.\prepare-github.ps1 -GitHubUsername "你的正確用戶名"
```

### Q4: Git 推送被拒絕
**解決方案**:
```bash
# 如果遠端有檔案，先拉取
git pull origin main --allow-unrelated-histories

# 然後再推送
git push -u origin main
```

---

## 📋 完整檢查清單

### 建立儲存庫前
- [ ] 已有 GitHub 帳號
- [ ] 決定好儲存庫名稱
- [ ] 本地已完成 TokenMonitor 開發

### 建立儲存庫時
- [ ] 儲存庫設為 Public
- [ ] 沒有勾選初始化選項
- [ ] 記錄了正確的 GitHub 路徑

### 準備檔案時
- [ ] 執行了 `prepare-github.ps1`
- [ ] 確認 GitHub 用戶名正確
- [ ] 檢查生成的檔案內容

### 推送到 GitHub 時
- [ ] Git 初始化成功
- [ ] 遠端 URL 設定正確
- [ ] 所有檔案都已提交
- [ ] 推送沒有錯誤

### 驗證部署後
- [ ] GitHub 頁面顯示正常
- [ ] 安裝腳本可以下載
- [ ] 測試安裝成功

---

## 🎉 完成後的效果

完成所有步驟後，任何人都可以這樣使用你的 TokenMonitor：

### Windows 用戶
```powershell
iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/你的用戶名/TokenMonitor/main/quick-install.ps1'))
```

### Linux/macOS 用戶
```bash
curl -sSL https://raw.githubusercontent.com/你的用戶名/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full
```

### 手動安裝
```bash
git clone https://github.com/你的用戶名/TokenMonitor.git
cd TokenMonitor
npm install
```

---

## 💡 後續維護

### 更新版本
```bash
# 修改程式碼後
git add .
git commit -m "Update: 新功能描述"
git push

# 發布新版本
git tag -a v1.1.0 -m "TokenMonitor v1.1.0 - 新功能"
git push origin v1.1.0
```

### 處理 Issues
- 定期檢查 GitHub Issues
- 回應使用者問題
- 修復 bug 並發布更新

### 推廣專案
- 在社群分享
- 寫部落格文章
- 參與相關討論

---

**總結**: 你需要先手動在 GitHub 建立儲存庫，然後使用我們的腳本準備檔案，最後用 Git 命令推送。腳本會幫你準備所有必要的檔案，但不會自動建立 GitHub 儲存庫。