$sessionPath = "C:\Users\yutachang\AppData\Roaming\kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\*"

Write-Host "🔍 檢查 Kiro 會話檔案..." -ForegroundColor Cyan

$files = Get-ChildItem -Path $sessionPath -Filter "*.json" | Sort-Object Length -Descending

Write-Host "📊 找到 $($files.Count) 個會話檔案" -ForegroundColor Green

foreach ($file in $files | Select-Object -First 5) {
    Write-Host "`n📄 檔案: $($file.Name)" -ForegroundColor Yellow
    Write-Host "   大小: $(($file.Length/1KB).ToString('F2')) KB" -ForegroundColor Gray
    
    try {
        $content = Get-Content $file.FullName -Raw -ErrorAction Stop
        
        if ($content.Contains("token") -or $content.Contains("chat") -or $content.Contains("message")) {
            Write-Host "   ✅ 包含相關內容！" -ForegroundColor Green
            
            # 嘗試解析 JSON
            try {
                $json = $content | ConvertFrom-Json
                if ($json.messages -or $json.conversation -or $json.history) {
                    Write-Host "   🎯 發現對話數據結構！" -ForegroundColor Magenta
                }
            }
            catch {
                Write-Host "   ⚠️ JSON 解析失敗" -ForegroundColor Yellow
            }
        } else {
            Write-Host "   ❌ 不包含相關內容" -ForegroundColor Red
        }
    }
    catch {
        Write-Host "   ❌ 無法讀取檔案" -ForegroundColor Red
    }
}