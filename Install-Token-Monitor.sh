#!/bin/bash

# TokenMonitor 一鍵安裝腳本 (Linux/macOS)
# 從 GitHub 下載並部署 TokenMonitor

set -e

# 預設值
TARGET_PATH=""
MODE="full"
VERSION="main"
GITHUB_REPO="cy5407/TokenMonitor"

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# 顯示幫助
show_help() {
    echo -e "${CYAN}🚀 TokenMonitor 一鍵安裝工具${NC}"
    echo ""
    echo "用法: $0 --target-path <路徑> [選項]"
    echo ""
    echo "參數:"
    echo "  --target-path    目標安裝路徑 (必要)"
    echo "  --mode          部署模式 (full/lite/npm)"
    echo "  --version       版本標籤 (預設: main)"
    echo "  --github-repo   GitHub 倉庫 (預設: cy5407/TokenMonitor)"
    echo "  --help          顯示此幫助"
    echo ""
    echo "部署模式:"
    echo "  full    完整安裝 - 包含所有功能和工具"
    echo "  lite    輕量安裝 - 只包含核心監控功能"
    echo "  npm     NPM 套件 - 生成可發布的 NPM 套件"
    echo ""
    echo "範例:"
    echo "  $0 --target-path /home/user/MyProject --mode full"
    echo "  $0 --target-path /home/user/MyProject --mode lite --version v1.0.0"
    echo ""
}

# 解析參數
while [[ $# -gt 0 ]]; do
    case $1 in
        --target-path)
            TARGET_PATH="$2"
            shift 2
            ;;
        --mode)
            MODE="$2"
            shift 2
            ;;
        --version)
            VERSION="$2"
            shift 2
            ;;
        --github-repo)
            GITHUB_REPO="$2"
            shift 2
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知參數: $1${NC}"
            echo "使用 --help 查看幫助"
            exit 1
            ;;
    esac
done

# 檢查必要參數
if [ -z "$TARGET_PATH" ]; then
    echo -e "${RED}❌ 錯誤: 請提供 --target-path 參數${NC}"
    echo "使用 --help 查看幫助"
    exit 1
fi

# 驗證模式
case $MODE in
    full|lite|npm)
        ;;
    *)
        echo -e "${RED}❌ 錯誤: 無效的模式 '$MODE'${NC}"
        echo "支援的模式: full, lite, npm"
        exit 1
        ;;
esac

echo -e "${GREEN}🚀 TokenMonitor 一鍵安裝工具${NC}"
echo "================================"
echo -e "${GRAY}📁 目標路徑: $TARGET_PATH${NC}"
echo -e "${GRAY}⚙️  部署模式: $MODE${NC}"
echo -e "${GRAY}🏷️  版本: $VERSION${NC}"
echo -e "${GRAY}📦 倉庫: $GITHUB_REPO${NC}"
echo ""

# 檢查依賴
check_dependencies() {
    local missing_deps=()
    
    if ! command -v curl >/dev/null 2>&1; then
        missing_deps+=("curl")
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        missing_deps+=("tar")
    fi
    
    if ! command -v node >/dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  警告: 未找到 Node.js，某些功能可能無法使用${NC}"
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少必要依賴: ${missing_deps[*]}${NC}"
        echo "請安裝缺少的依賴後重試"
        exit 1
    fi
}

# 創建目標目錄
create_target_directory() {
    if [ ! -d "$TARGET_PATH" ]; then
        echo -e "${YELLOW}📁 創建目標目錄: $TARGET_PATH${NC}"
        mkdir -p "$TARGET_PATH"
    fi
}

# 下載和解壓縮
download_and_extract() {
    # 創建臨時目錄
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    echo -e "${GRAY}📂 創建臨時目錄: $TEMP_DIR${NC}"
    
    # 構建下載 URL
    local download_url="https://github.com/$GITHUB_REPO/archive/$VERSION.tar.gz"
    local tar_file="$TEMP_DIR/TokenMonitor.tar.gz"
    
    echo -e "${YELLOW}📥 從 GitHub 下載中...${NC}"
    echo -e "${GRAY}🔗 URL: $download_url${NC}"
    
    # 下載檔案
    if ! curl -L "$download_url" -o "$tar_file" --silent --show-error; then
        echo -e "${RED}❌ 下載失敗${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 下載完成${NC}"
    
    # 檢查下載的檔案
    if [ ! -f "$tar_file" ] || [ ! -s "$tar_file" ]; then
        echo -e "${RED}❌ 下載的檔案無效或為空${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}📦 解壓縮中...${NC}"
    
    # 解壓縮
    if ! tar -xzf "$tar_file" -C "$TEMP_DIR"; then
        echo -e "${RED}❌ 解壓縮失敗${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 解壓縮完成${NC}"
    
    # 找到解壓縮的目錄
    EXTRACTED_DIR=$(find "$TEMP_DIR" -name "TokenMonitor-*" -type d | head -1)
    
    if [ -z "$EXTRACTED_DIR" ]; then
        echo -e "${RED}❌ 找不到解壓縮的 TokenMonitor 目錄${NC}"
        exit 1
    fi
    
    echo -e "${GRAY}📂 找到解壓縮目錄: $(basename "$EXTRACTED_DIR")${NC}"
}

# 安裝完整版
install_full() {
    echo -e "${YELLOW}🔧 執行完整安裝...${NC}"
    
    local tokenmonitor_dir="$TARGET_PATH/TokenMonitor"
    mkdir -p "$tokenmonitor_dir"
    
    # 複製主要檔案和目錄
    local dirs_to_copy=("scripts" "src" ".kiro" "docs")
    local files_to_copy=("package.json")
    
    for dir in "${dirs_to_copy[@]}"; do
        if [ -d "$EXTRACTED_DIR/$dir" ]; then
            cp -r "$EXTRACTED_DIR/$dir" "$tokenmonitor_dir/"
            echo -e "${GREEN}✅ 複製目錄: $dir${NC}"
        fi
    done
    
    for file in "${files_to_copy[@]}"; do
        if [ -f "$EXTRACTED_DIR/$file" ]; then
            cp "$EXTRACTED_DIR/$file" "$tokenmonitor_dir/"
            echo -e "${GREEN}✅ 複製檔案: $file${NC}"
        fi
    done
    
    # 創建必要目錄
    local required_dirs=("data" "tests/data" "tests/reports")
    for dir in "${required_dirs[@]}"; do
        mkdir -p "$tokenmonitor_dir/$dir"
        echo -e "${BLUE}📁 創建目錄: $dir${NC}"
    done
    
    # 如果有 PowerShell，設置執行權限
    if [ -f "$tokenmonitor_dir/scripts/tokusage.ps1" ]; then
        chmod +x "$tokenmonitor_dir/scripts/tokusage.ps1"
    fi
}

# 安裝輕量版
install_lite() {
    echo -e "${YELLOW}📦 安裝輕量版...${NC}"
    
    local lite_dir="$TARGET_PATH/token-monitor"
    mkdir -p "$lite_dir"
    
    # 創建輕量版監控腳本
    cat > "$lite_dir/token-monitor.js" << 'EOF'
// TokenMonitor Lite - 輕量版 Token 監控 (Linux/macOS)
const fs = require('fs');
const path = require('path');

class TokenMonitor {
    constructor(logPath = './token-usage.log') {
        this.logPath = logPath;
    }
    
    log(event, tokens, cost = 0, model = 'unknown') {
        const record = {
            timestamp: new Date().toISOString(),
            event,
            tokens: parseInt(tokens),
            cost: parseFloat(cost),
            model,
            session: 'lite-' + Date.now(),
            platform: process.platform
        };
        
        fs.appendFileSync(this.logPath, JSON.stringify(record) + '\n');
    }
    
    analyze(days = 7) {
        if (!fs.existsSync(this.logPath)) {
            return { total: 0, cost: 0, records: 0 };
        }
        
        const lines = fs.readFileSync(this.logPath, 'utf8').split('\n').filter(Boolean);
        const cutoff = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
        
        const records = lines
            .map(line => {
                try {
                    return JSON.parse(line);
                } catch {
                    return null;
                }
            })
            .filter(r => r && new Date(r.timestamp) > cutoff);
        
        const total = records.reduce((sum, r) => sum + r.tokens, 0);
        const cost = records.reduce((sum, r) => sum + r.cost, 0);
        
        return { 
            total, 
            cost: cost.toFixed(6), 
            records: records.length,
            daily: (total / days).toFixed(0),
            dailyCost: (cost / days).toFixed(6)
        };
    }
    
    report() {
        const stats = this.analyze();
        console.log('📊 TokenMonitor Lite 報告');
        console.log('========================');
        console.log(`總 Token: ${stats.total}`);
        console.log(`總成本: $${stats.cost}`);
        console.log(`記錄數: ${stats.records}`);
        console.log(`日均 Token: ${stats.daily}`);
        console.log(`日均成本: $${stats.dailyCost}`);
        console.log(`平台: ${process.platform}`);
    }
    
    cleanup(days = 30) {
        if (!fs.existsSync(this.logPath)) return 0;
        
        const cutoff = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
        const lines = fs.readFileSync(this.logPath, 'utf8').split('\n').filter(Boolean);
        
        const validLines = lines.filter(line => {
            try {
                const record = JSON.parse(line);
                return new Date(record.timestamp) > cutoff;
            } catch {
                return false;
            }
        });
        
        const removedCount = lines.length - validLines.length;
        
        if (removedCount > 0) {
            fs.writeFileSync(this.logPath, validLines.join('\n') + '\n');
            console.log(`🧹 清理了 ${removedCount} 筆舊記錄`);
        }
        
        return removedCount;
    }
}

module.exports = TokenMonitor;

// CLI 使用
if (require.main === module) {
    const monitor = new TokenMonitor();
    const command = process.argv[2];
    
    switch (command) {
        case 'report':
            monitor.report();
            break;
        case 'log':
            const [, , , event, tokens, cost] = process.argv;
            if (event && tokens) {
                monitor.log(event, tokens, cost);
                console.log(`✅ 記錄: ${event} - ${tokens} tokens`);
            } else {
                console.log('用法: node token-monitor.js log <event> <tokens> [cost]');
            }
            break;
        case 'cleanup':
            const days = parseInt(process.argv[3]) || 30;
            monitor.cleanup(days);
            break;
        default:
            console.log('用法: node token-monitor.js [report|log <event> <tokens> <cost>|cleanup [days]]');
    }
}
EOF
    
    echo -e "${GREEN}✅ 創建: token-monitor.js${NC}"
    
    # 創建 README
    cat > "$lite_dir/README.md" << EOF
# TokenMonitor Lite

輕量版 Token 使用監控工具 (Linux/macOS)

## 使用方式

\`\`\`javascript
const TokenMonitor = require('./token-monitor');
const monitor = new TokenMonitor();

// 記錄使用
monitor.log('chat_message', 150, 0.00045);

// 查看報告
monitor.report();

// 清理舊記錄
monitor.cleanup(30);
\`\`\`

## CLI 使用

\`\`\`bash
# 查看報告
node token-monitor.js report

# 記錄使用
node token-monitor.js log chat_message 150 0.00045

# 清理舊記錄 (保留30天)
node token-monitor.js cleanup 30
\`\`\`

## 安裝資訊

- 安裝來源: TokenMonitor 一鍵安裝腳本
- GitHub: https://github.com/$GITHUB_REPO
- 版本: $VERSION
- 安裝時間: $(date '+%Y-%m-%d %H:%M:%S')
- 平台: $(uname -s)
EOF
    
    echo -e "${GREEN}✅ 創建: README.md${NC}"
    
    # 設置執行權限
    chmod +x "$lite_dir/token-monitor.js"
}

# 安裝 NPM 套件
install_npm() {
    echo -e "${YELLOW}📦 安裝 NPM 套件模板...${NC}"
    
    local npm_dir="$TARGET_PATH/kiro-token-monitor"
    mkdir -p "$npm_dir"
    
    # 複製 NPM 套件模板 (如果存在)
    if [ -d "$EXTRACTED_DIR/templates/npm-package" ]; then
        cp -r "$EXTRACTED_DIR/templates/npm-package/"* "$npm_dir/"
        echo -e "${GREEN}✅ 複製 NPM 套件模板${NC}"
    else
        echo -e "${YELLOW}⚠️  找不到 NPM 套件模板，創建基本結構...${NC}"
        
        # 創建基本的 package.json
        cat > "$npm_dir/package.json" << EOF
{
  "name": "kiro-token-monitor",
  "version": "1.0.0",
  "description": "AI Token usage monitoring tool",
  "main": "index.js",
  "bin": {
    "tokusage": "./bin/tokusage.js"
  },
  "scripts": {
    "start": "node index.js",
    "test": "echo \"No tests specified\""
  },
  "dependencies": {
    "commander": "^9.0.0",
    "chokidar": "^3.5.3"
  },
  "keywords": ["token", "monitoring", "ai", "cost"],
  "author": "TokenMonitor",
  "license": "MIT"
}
EOF
        
        echo -e "${GREEN}✅ 創建基本 package.json${NC}"
        
        # 創建基本的 index.js
        cat > "$npm_dir/index.js" << 'EOF'
// TokenMonitor NPM Package Entry Point
console.log('🚀 TokenMonitor NPM Package');
console.log('請查看 README.md 了解使用方法');

module.exports = require('./token-monitor');
EOF
        
        echo -e "${GREEN}✅ 創建基本 index.js${NC}"
    fi
}

# 顯示後續步驟
show_next_steps() {
    echo ""
    echo -e "${GREEN}🎉 TokenMonitor 安裝完成！${NC}"
    echo ""
    echo -e "${CYAN}📋 後續步驟:${NC}"
    
    case $MODE in
        full)
            echo -e "${GRAY}1. cd \"$TARGET_PATH/TokenMonitor\"${NC}"
            echo -e "${GRAY}2. npm install${NC}"
            if command -v pwsh >/dev/null 2>&1; then
                echo -e "${GRAY}3. pwsh ./scripts/tokusage.ps1 daily${NC}"
            else
                echo -e "${GRAY}3. node src/js/professional-token-cli.js${NC}"
            fi
            ;;
        lite)
            echo -e "${GRAY}1. cd \"$TARGET_PATH/token-monitor\"${NC}"
            echo -e "${GRAY}2. node token-monitor.js report${NC}"
            ;;
        npm)
            echo -e "${GRAY}1. cd \"$TARGET_PATH/kiro-token-monitor\"${NC}"
            echo -e "${GRAY}2. npm install${NC}"
            echo -e "${GRAY}3. node index.js${NC}"
            ;;
    esac
    
    echo ""
    echo -e "${YELLOW}💡 提示:${NC}"
    echo -e "${GRAY}• 查看相關文件了解更多功能${NC}"
    echo -e "${GRAY}• 定期更新以獲得最新功能${NC}"
    echo -e "${GRAY}• 遇到問題請查看 GitHub Issues${NC}"
    echo ""
    echo -e "${YELLOW}🌟 如果覺得有用，請給我們一個 Star！${NC}"
    echo -e "${BLUE}🔗 https://github.com/$GITHUB_REPO${NC}"
}

# 主要執行流程
main() {
    check_dependencies
    create_target_directory
    download_and_extract
    
    case $MODE in
        full)
            install_full
            ;;
        lite)
            install_lite
            ;;
        npm)
            install_npm
            ;;
    esac
    
    show_next_steps
}

# 錯誤處理
handle_error() {
    echo ""
    echo -e "${RED}❌ 安裝失敗: $1${NC}"
    echo ""
    echo -e "${YELLOW}🔧 故障排除建議:${NC}"
    echo -e "${GRAY}• 檢查網路連線${NC}"
    echo -e "${GRAY}• 確認 GitHub 倉庫存在且可訪問${NC}"
    echo -e "${GRAY}• 檢查目標路徑權限${NC}"
    echo -e "${GRAY}• 嘗試使用不同的版本標籤${NC}"
    echo -e "${GRAY}• 確認已安裝必要依賴 (curl, tar, node)${NC}"
    exit 1
}

# 設置錯誤處理
trap 'handle_error "執行過程中發生錯誤"' ERR

# 執行主程式
main
