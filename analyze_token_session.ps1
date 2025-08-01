Write-Host "🔍 分析 TokenMonitor 會話檔案..." -ForegroundColor Cyan

$sessionFiles = Get-ChildItem "C:\Users\yutachang\AppData\Roaming\kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\" | 
                Where-Object { $_.Name -like "*TE9HIOaXpeiqjOafpeipog*" }

foreach ($file in $sessionFiles) {
    Write-Host "`n📄 檔案: $($file.Name)" -ForegroundColor Yellow
    Write-Host "   大小: $(($file.Length/1KB).ToString('F2')) KB" -ForegroundColor Gray
    
    $content = Get-Content $file.FullName -Raw
    
    if ($content.Contains("token") -or $content.Contains("chat") -or $content.Contains("message")) {
        Write-Host "   ✅ 包含相關內容！" -ForegroundColor Green
        
        # 計算可能的 token 數量
        $tokenMatches = [regex]::Matches($content, "token[s]?[:\s]*(\d+)", [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
        if ($tokenMatches.Count -gt 0) {
            $totalTokens = 0
            foreach ($match in $tokenMatches) {
                $totalTokens += [int]$match.Groups[1].Value
            }
            Write-Host "   💎 發現 Token 數據: $totalTokens tokens" -ForegroundColor Magenta
        }
        
        # 檢查是否有對話內容
        if ($content.Contains("出師表") -or $content.Contains("程式設計指南")) {
            Write-Host "   🎯 發現我們的對話內容！" -ForegroundColor Green
        }
    } else {
        Write-Host "   ❌ 不包含相關內容" -ForegroundColor Red
    }
}