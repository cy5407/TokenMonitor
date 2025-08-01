#!/usr/bin/env pwsh

# TokenMonitor 推送到 GitHub 腳本
# 這個腳本是為 cy5407/TokenMonitor 自動生成的

Write-Host "🚀 推送 TokenMonitor 到 GitHub..." -ForegroundColor Green

try {
    # 添加遠端倉庫
    Write-Host "🔗 添加遠端倉庫..." -ForegroundColor Yellow
    git remote add origin https://github.com/cy5407/TokenMonitor.git
    
    # 設定主分支
    Write-Host "🌿 設定主分支..." -ForegroundColor Yellow
    git branch -M main
    
    # 推送到 GitHub
    Write-Host "📤 推送到 GitHub..." -ForegroundColor Yellow
    git push -u origin main
    
    Write-Host ""
    Write-Host "✅ 推送完成！" -ForegroundColor Green
    Write-Host ""
    Write-Host "🌟 你的 TokenMonitor 現在可以在這裡找到:" -ForegroundColor Yellow
    Write-Host "https://github.com/cy5407/TokenMonitor" -ForegroundColor Blue
    Write-Host ""
    Write-Host "💡 建議創建版本標籤:" -ForegroundColor Cyan
    Write-Host "git tag -a v1.0.0 -m `"TokenMonitor v1.0.0 - Initial Release`"" -ForegroundColor Gray
    Write-Host "git push origin v1.0.0" -ForegroundColor Gray
    Write-Host ""
    Write-Host "🎉 完成後，使用者可以這樣安裝:" -ForegroundColor Green
    Write-Host ""
    Write-Host "Windows 一鍵安裝:" -ForegroundColor Cyan
    Write-Host "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/quick-install.ps1'))" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Linux/macOS 一鍵安裝:" -ForegroundColor Cyan
    Write-Host "curl -sSL https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full" -ForegroundColor Blue
    
} catch {
    Write-Error "推送失敗: $($_.Exception.Message)"
    Write-Host ""
    Write-Host "🔧 可能的解決方案:" -ForegroundColor Yellow
    Write-Host "1. 確認已在 GitHub 建立儲存庫 https://github.com/cy5407/TokenMonitor" -ForegroundColor Gray
    Write-Host "2. 檢查 GitHub 使用者名稱和儲存庫名稱" -ForegroundColor Gray
    Write-Host "3. 確認有儲存庫的寫入權限" -ForegroundColor Gray
    Write-Host "4. 如果遠端已存在，先移除: git remote remove origin" -ForegroundColor Gray
}