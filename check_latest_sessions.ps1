Write-Host "🔍 檢查最新的 Kiro 會話檔案..." -ForegroundColor Cyan

# 檢查今天的檔案
$today = Get-Date
$todayFiles = Get-ChildItem "C:\Users\yutachang\AppData\Roaming\kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\" -Recurse -Filter "*.json" | 
              Where-Object { $_.LastWriteTime.Date -eq $today.Date } | 
              Sort-Object LastWriteTime -Descending

Write-Host "📅 今天 ($($today.ToString('yyyy-MM-dd'))) 的會話檔案數量: $($todayFiles.Count)" -ForegroundColor Green

if ($todayFiles.Count -gt 0) {
    Write-Host "`n🎉 找到今天的會話檔案！" -ForegroundColor Green
    
    foreach ($file in $todayFiles | Select-Object -First 5) {
        Write-Host "`n📄 檔案: $($file.Name)" -ForegroundColor Yellow
        Write-Host "   路徑: $($file.Directory.Name)" -ForegroundColor Gray
        Write-Host "   大小: $(($file.Length/1KB).ToString('F2')) KB" -ForegroundColor Gray
        Write-Host "   修改時間: $($file.LastWriteTime)" -ForegroundColor Gray
        
        try {
            $content = Get-Content $file.FullName -Raw -ErrorAction Stop
            
            # 檢查是否包含我們的對話內容
            $hasOurContent = $false
            if ($content.Contains("出師表")) {
                Write-Host "   ✅ 包含出師表內容！" -ForegroundColor Green
                $hasOurContent = $true
            }
            if ($content.Contains("程式設計指南")) {
                Write-Host "   ✅ 包含程式設計指南內容！" -ForegroundColor Green
                $hasOurContent = $true
            }
            if ($content.Contains("Token 獵人")) {
                Write-Host "   ✅ 包含 Token 獵人內容！" -ForegroundColor Green
                $hasOurContent = $true
            }
            if ($content.Contains("Ultrathink")) {
                Write-Host "   ✅ 包含 Ultrathink 內容！" -ForegroundColor Green
                $hasOurContent = $true
            }
            
            # 檢查 token 相關內容
            if ($content.Contains("token") -or $content.Contains("Token")) {
                Write-Host "   💎 包含 Token 相關內容！" -ForegroundColor Magenta
                
                # 嘗試計算 token 數量
                $tokenMatches = [regex]::Matches($content, "token[s]?[:\s]*(\d+)", [System.Text.RegularExpressions.RegexOptions]::IgnoreCase)
                if ($tokenMatches.Count -gt 0) {
                    $totalTokens = 0
                    foreach ($match in $tokenMatches) {
                        $totalTokens += [int]$match.Groups[1].Value
                    }
                    Write-Host "   🔢 發現 Token 數據: $totalTokens tokens" -ForegroundColor Cyan
                }
            }
            
            if (-not $hasOurContent) {
                Write-Host "   ❌ 不包含我們今天的對話內容" -ForegroundColor Red
            }
            
        } catch {
            Write-Host "   ❌ 無法讀取檔案: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
} else {
    Write-Host "`n⚠️ 沒有找到今天的會話檔案" -ForegroundColor Yellow
    
    # 檢查最近的檔案
    Write-Host "`n🔍 檢查最近修改的檔案..." -ForegroundColor Cyan
    $recentFiles = Get-ChildItem "C:\Users\yutachang\AppData\Roaming\kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\" -Recurse -Filter "*.json" | 
                   Sort-Object LastWriteTime -Descending | 
                   Select-Object -First 3
    
    foreach ($file in $recentFiles) {
        Write-Host "`n📄 最近檔案: $($file.Name)" -ForegroundColor Yellow
        Write-Host "   修改時間: $($file.LastWriteTime)" -ForegroundColor Gray
        Write-Host "   大小: $(($file.Length/1KB).ToString('F2')) KB" -ForegroundColor Gray
    }
}

Write-Host "`n🎯 總結:" -ForegroundColor Cyan
if ($todayFiles.Count -gt 0) {
    Write-Host "✅ Kiro IDE 重開後有產生新的會話檔案" -ForegroundColor Green
} else {
    Write-Host "❌ Kiro IDE 重開後仍未產生今天的會話檔案" -ForegroundColor Red
    Write-Host "💡 可能的原因：" -ForegroundColor Yellow
    Write-Host "   1. 會話還沒有被儲存（需要更多時間）" -ForegroundColor Gray
    Write-Host "   2. 儲存機制有問題" -ForegroundColor Gray
    Write-Host "   3. 儲存在其他位置" -ForegroundColor Gray
}