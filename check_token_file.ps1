Write-Host "🔍 檢查包含 Token 數據的檔案..." -ForegroundColor Cyan

$tokenFile = Get-ChildItem "C:\Users\yutachang\AppData\Roaming\kiro\User\globalStorage\kiro.kiroagent\workspace-sessions\" -Recurse -Filter "96dc1507-64fa-4896-87b9-8ffc177521de.json"

if ($tokenFile) {
    Write-Host "📄 檔案: $($tokenFile.Name)" -ForegroundColor Yellow
    Write-Host "   大小: $(($tokenFile.Length/1KB).ToString('F2')) KB" -ForegroundColor Gray
    
    $content = Get-Content $tokenFile.FullName -Raw
    
    Write-Host "`n📊 內容分析:" -ForegroundColor Cyan
    Write-Host "   總字符數: $($content.Length)" -ForegroundColor Gray
    Write-Host "   包含 'Ultrathink': $($content.Contains('Ultrathink'))" -ForegroundColor Gray
    Write-Host "   包含 '獵人': $($content.Contains('獵人'))" -ForegroundColor Gray
    Write-Host "   包含 'monitor': $($content.Contains('monitor'))" -ForegroundColor Gray
    Write-Host "   包含 '出師表': $($content.Contains('出師表'))" -ForegroundColor Gray
    Write-Host "   包含 'token': $($content.Contains('token'))" -ForegroundColor Gray
    
    # 顯示內容預覽
    Write-Host "`n📝 內容預覽 (前 800 字符):" -ForegroundColor Cyan
    $preview = $content.Substring(0, [Math]::Min(800, $content.Length))
    Write-Host $preview -ForegroundColor White
    
    # 嘗試解析 JSON 結構
    try {
        $json = $content | ConvertFrom-Json
        Write-Host "`n🏗️ JSON 結構分析:" -ForegroundColor Cyan
        
        if ($json.history) {
            Write-Host "   ✅ 包含 history 陣列，長度: $($json.history.Count)" -ForegroundColor Green
            
            # 檢查最近的幾個對話
            $recentMessages = $json.history | Select-Object -Last 3
            foreach ($msg in $recentMessages) {
                if ($msg.message) {
                    $role = $msg.message.role
                    $contentPreview = $msg.message.content.Substring(0, [Math]::Min(100, $msg.message.content.Length))
                    Write-Host "   💬 $role`: $contentPreview..." -ForegroundColor Cyan
                }
            }
        }
        
        if ($json.messages) {
            Write-Host "   ✅ 包含 messages 陣列，長度: $($json.messages.Count)" -ForegroundColor Green
        }
        
    } catch {
        Write-Host "   ❌ JSON 解析失敗: $($_.Exception.Message)" -ForegroundColor Red
    }
    
} else {
    Write-Host "❌ 找不到指定的檔案" -ForegroundColor Red
}