# TokenMonitor 快速安裝腳本 - 一行命令安裝
# 用法: iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/quick-install.ps1'))

param(
    [string]$TargetPath = (Join-Path $PWD "TokenMonitor"),
    [ValidateSet("full", "lite", "npm")]
    [string]$Mode = "full"
)

Write-Host "🚀 TokenMonitor 快速安裝" -ForegroundColor Green
Write-Host "目標路徑: $TargetPath" -ForegroundColor Gray
Write-Host "安裝模式: $Mode" -ForegroundColor Gray

try {
    # 下載完整安裝腳本
    $installScript = (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/cy5407/TokenMonitor/main/install-tokenmonitor.ps1')
    
    # 創建臨時腳本檔案
    $tempScript = Join-Path $env:TEMP "install-tokenmonitor-temp.ps1"
    Set-Content -Path $tempScript -Value $installScript
    
    # 執行安裝
    & $tempScript -TargetPath $TargetPath -Mode $Mode
    
    # 清理
    Remove-Item $tempScript -ErrorAction SilentlyContinue
    
} catch {
    Write-Error "快速安裝失敗: $($_.Exception.Message)"
    Write-Host "請嘗試手動下載安裝腳本" -ForegroundColor Yellow
}
