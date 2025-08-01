#!/usr/bin/env pwsh

<#
.SYNOPSIS
    TokenMonitor GitHub 快速設置腳本

.DESCRIPTION
    這個腳本會引導你完成 GitHub 儲存庫的設置過程

.PARAMETER GitHubUsername
    你的 GitHub 使用者名稱

.PARAMETER RepoName
    儲存庫名稱 (預設: TokenMonitor)

.EXAMPLE
    .\quick-setup-github.ps1 -GitHubUsername "myusername"
#>

param(
    [Parameter(HelpMessage="GitHub 使用者名稱")]
    [string]$GitHubUsername,
    
    [Parameter(HelpMessage="儲存庫名稱")]
    [string]$RepoName = "TokenMonitor",
    
    [Parameter(HelpMessage="跳過互動式設置")]
    [switch]$NonInteractive,
    
    [Parameter(HelpMessage="顯示幫助")]
    [switch]$Help
)

if ($Help) {
    Write-Host @"
🚀 TokenMonitor GitHub 快速設置工具

這個腳本會引導你完成 GitHub 儲存庫的完整設置過程。

用法:
    quick-setup-github.ps1 [選項]

參數:
    -GitHubUsername    GitHub 使用者名稱
    -RepoName         儲存庫名稱 (預設: TokenMonitor)
    -NonInteractive   跳過互動式設置
    -Help             顯示此幫助

範例:
    .\quick-setup-github.ps1
    .\quick-setup-github.ps1 -GitHubUsername "myusername"

"@ -ForegroundColor Cyan
    exit 0
}

Write-Host @"
🚀 TokenMonitor GitHub 快速設置
================================

這個腳本會幫助你將 TokenMonitor 部署到 GitHub，
讓全世界的開發者都能使用你的工具！

"@ -ForegroundColor Green

# 互動式獲取資訊
if (-not $NonInteractive) {
    if (-not $GitHubUsername) {
        Write-Host "📝 請提供你的 GitHub 資訊:" -ForegroundColor Yellow
        Write-Host ""
        
        do {
            $GitHubUsername = Read-Host "GitHub 使用者名稱"
            if (-not $GitHubUsername) {
                Write-Host "❌ 使用者名稱不能為空，請重新輸入" -ForegroundColor Red
            }
        } while (-not $GitHubUsername)
    }
    
    Write-Host ""
    $customRepoName = Read-Host "儲存庫名稱 (預設: $RepoName，直接按 Enter 使用預設值)"
    if ($customRepoName) {
        $RepoName = $customRepoName
    }
    
    Write-Host ""
    Write-Host "📋 確認資訊:" -ForegroundColor Cyan
    Write-Host "GitHub 使用者: $GitHubUsername" -ForegroundColor Gray
    Write-Host "儲存庫名稱: $RepoName" -ForegroundColor Gray
    Write-Host ""
    
    $confirm = Read-Host "確認無誤？(y/N)"
    if ($confirm -ne 'y' -and $confirm -ne 'Y') {
        Write-Host "❌ 已取消設置" -ForegroundColor Red
        exit 0
    }
}

if (-not $GitHubUsername) {
    Write-Error "請提供 GitHub 使用者名稱"
    exit 1
}

Write-Host ""
Write-Host "🔧 開始設置過程..." -ForegroundColor Green

try {
    # 步驟 1: 檢查 Git 是否已安裝
    Write-Host ""
    Write-Host "1️⃣ 檢查 Git 環境..." -ForegroundColor Yellow
    
    try {
        $gitVersion = git --version
        Write-Host "✅ Git 已安裝: $gitVersion" -ForegroundColor Green
    } catch {
        Write-Host "❌ 未找到 Git，請先安裝 Git" -ForegroundColor Red
        Write-Host "下載地址: https://git-scm.com/download" -ForegroundColor Blue
        exit 1
    }
    
    # 步驟 2: 檢查 Git 配置
    Write-Host ""
    Write-Host "2️⃣ 檢查 Git 配置..." -ForegroundColor Yellow
    
    $gitUser = git config --global user.name 2>$null
    $gitEmail = git config --global user.email 2>$null
    
    if (-not $gitUser -or -not $gitEmail) {
        Write-Host "⚠️  Git 使用者資訊未設定" -ForegroundColor Yellow
        
        if (-not $NonInteractive) {
            Write-Host ""
            if (-not $gitUser) {
                $userName = Read-Host "請輸入你的姓名"
                git config --global user.name "$userName"
                Write-Host "✅ 設定 Git 使用者名稱: $userName" -ForegroundColor Green
            }
            
            if (-not $gitEmail) {
                $userEmail = Read-Host "請輸入你的 Email"
                git config --global user.email "$userEmail"
                Write-Host "✅ 設定 Git Email: $userEmail" -ForegroundColor Green
            }
        } else {
            Write-Host "❌ 請先設定 Git 使用者資訊:" -ForegroundColor Red
            Write-Host "git config --global user.name `"你的姓名`"" -ForegroundColor Gray
            Write-Host "git config --global user.email `"你的email@example.com`"" -ForegroundColor Gray
            exit 1
        }
    } else {
        Write-Host "✅ Git 使用者: $gitUser <$gitEmail>" -ForegroundColor Green
    }
    
    # 步驟 3: 準備檔案
    Write-Host ""
    Write-Host "3️⃣ 準備 GitHub 檔案..." -ForegroundColor Yellow
    
    if (Test-Path "prepare-github.ps1") {
        & ".\prepare-github.ps1" -GitHubUsername $GitHubUsername -RepoName $RepoName
        Write-Host "✅ 檔案準備完成" -ForegroundColor Green
    } else {
        Write-Host "❌ 找不到 prepare-github.ps1 腳本" -ForegroundColor Red
        exit 1
    }
    
    # 步驟 4: 初始化 Git
    Write-Host ""
    Write-Host "4️⃣ 初始化 Git 倉庫..." -ForegroundColor Yellow
    
    if (-not (Test-Path ".git")) {
        git init
        Write-Host "✅ Git 倉庫初始化完成" -ForegroundColor Green
    } else {
        Write-Host "✅ Git 倉庫已存在" -ForegroundColor Green
    }
    
    # 步驟 5: 添加檔案
    Write-Host ""
    Write-Host "5️⃣ 添加檔案到 Git..." -ForegroundColor Yellow
    
    git add .
    Write-Host "✅ 檔案添加完成" -ForegroundColor Green
    
    # 步驟 6: 提交變更
    Write-Host ""
    Write-Host "6️⃣ 提交變更..." -ForegroundColor Yellow
    
    git commit -m "Initial commit: TokenMonitor v1.0"
    Write-Host "✅ 變更提交完成" -ForegroundColor Green
    
    # 步驟 7: 顯示後續步驟
    Write-Host ""
    Write-Host "🎉 本地設置完成！" -ForegroundColor Green
    Write-Host ""
    Write-Host "📋 接下來你需要手動完成以下步驟:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "1. 在 GitHub 建立儲存庫:" -ForegroundColor Yellow
    Write-Host "   • 前往 https://github.com/new" -ForegroundColor Blue
    Write-Host "   • Repository name: $RepoName" -ForegroundColor Gray
    Write-Host "   • Description: AI Token usage monitoring and cost analysis tool" -ForegroundColor Gray
    Write-Host "   • 設為 Public" -ForegroundColor Gray
    Write-Host "   • ❌ 不要勾選任何初始化選項" -ForegroundColor Gray
    Write-Host "   • 點擊 'Create repository'" -ForegroundColor Gray
    Write-Host ""
    Write-Host "2. 建立完成後，執行以下命令推送:" -ForegroundColor Yellow
    Write-Host "   git remote add origin https://github.com/$GitHubUsername/$RepoName.git" -ForegroundColor Blue
    Write-Host "   git branch -M main" -ForegroundColor Blue
    Write-Host "   git push -u origin main" -ForegroundColor Blue
    Write-Host ""
    Write-Host "3. (可選) 創建版本標籤:" -ForegroundColor Yellow
    Write-Host "   git tag -a v1.0.0 -m `"TokenMonitor v1.0.0 - Initial Release`"" -ForegroundColor Blue
    Write-Host "   git push origin v1.0.0" -ForegroundColor Blue
    Write-Host ""
    Write-Host "4. 完成後，使用者可以這樣安裝:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "   Windows:" -ForegroundColor Cyan
    Write-Host "   iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/quick-install.ps1'))" -ForegroundColor Blue
    Write-Host ""
    Write-Host "   Linux/macOS:" -ForegroundColor Cyan
    Write-Host "   curl -sSL https://raw.githubusercontent.com/$GitHubUsername/$RepoName/main/install-tokenmonitor.sh | bash -s -- --target-path ./TokenMonitor --mode full" -ForegroundColor Blue
    Write-Host ""
    
    # 創建便利腳本
    $pushScriptContent = @"
#!/usr/bin/env pwsh

# TokenMonitor 推送到 GitHub 腳本
# 這個腳本是由 quick-setup-github.ps1 自動生成的

Write-Host "🚀 推送 TokenMonitor 到 GitHub..." -ForegroundColor Green

try {
    # 添加遠端倉庫
    git remote add origin https://github.com/$GitHubUsername/$RepoName.git
    
    # 設定主分支
    git branch -M main
    
    # 推送到 GitHub
    git push -u origin main
    
    Write-Host "✅ 推送完成！" -ForegroundColor Green
    Write-Host ""
    Write-Host "🌟 你的 TokenMonitor 現在可以在這裡找到:" -ForegroundColor Yellow
    Write-Host "https://github.com/$GitHubUsername/$RepoName" -ForegroundColor Blue
    Write-Host ""
    Write-Host "💡 建議創建版本標籤:" -ForegroundColor Cyan
    Write-Host "git tag -a v1.0.0 -m `"TokenMonitor v1.0.0 - Initial Release`"" -ForegroundColor Gray
    Write-Host "git push origin v1.0.0" -ForegroundColor Gray
    
} catch {
    Write-Error "推送失敗: `$(`$_.Exception.Message)"
    Write-Host ""
    Write-Host "🔧 可能的解決方案:" -ForegroundColor Yellow
    Write-Host "1. 確認已在 GitHub 建立儲存庫" -ForegroundColor Gray
    Write-Host "2. 檢查 GitHub 使用者名稱和儲存庫名稱" -ForegroundColor Gray
    Write-Host "3. 確認有儲存庫的寫入權限" -ForegroundColor Gray
}
"@
    
    Set-Content "push-to-github.ps1" -Value $pushScriptContent -Encoding UTF8
    Write-Host "💡 已創建便利腳本: push-to-github.ps1" -ForegroundColor Cyan
    Write-Host "   在 GitHub 建立儲存庫後，執行此腳本即可推送" -ForegroundColor Gray
    Write-Host ""
    
    Write-Host "📖 詳細說明請查看: GITHUB-SETUP-GUIDE.md" -ForegroundColor Cyan
    
} catch {
    Write-Error "設置過程中發生錯誤: $($_.Exception.Message)"
    exit 1
}

Write-Host ""
Write-Host "✨ 快速設置完成！按照上述步驟完成 GitHub 部署吧！" -ForegroundColor Green