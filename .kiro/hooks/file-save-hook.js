// Token Monitor - 檔案保存安全檢查 Hook
// 在檔案保存時檢查機密資訊、API Keys 和 Token 洩漏

const fs = require('fs');
const path = require('path');

// 設定控制台編碼 (Windows)
if (process.platform === 'win32') {
    process.stdout.setEncoding('utf8');
    process.stderr.setEncoding('utf8');
}

// 從配置檔案載入設定
const config = JSON.parse(fs.readFileSync(path.join(__dirname, 'file-save-hook.json'), 'utf8')).config;

// 主要執行函數
async function main(context) {
    try {
        const filePath = context.filePath;
        
        // 檢查是否應該處理此檔案
        if (!shouldProcessFile(filePath)) {
            return;
        }
        
        console.log(`🔍 [Token Monitor] 檢查檔案: ${filePath}`);
        
        // 檢查檔案
        const issues = await checkFile(filePath);
        
        if (issues.length > 0) {
            // 顯示通知
            if (config.notification.enabled && config.notification.showPopup) {
                showNotification(issues);
            }
            
            // 高亮顯示問題
            if (config.visualization.enabled) {
                highlightIssues(context, issues);
            }
            
            console.log(`⚠️  [Token Monitor] 在 ${filePath} 中發現 ${issues.length} 個潛在安全問題`);
            issues.forEach(issue => {
                console.log(`   行 ${issue.line}: ${issue.message} (${issue.severity})`);
            });
        } else {
            console.log(`✅ [Token Monitor] ${filePath} 安全檢查通過`);
        }
        
    } catch (error) {
        console.error('[Token Monitor] 檔案保存安全檢查失敗:', error.message);
    }
}

// 檢查是否應該處理此檔案
function shouldProcessFile(filePath) {
    // 檢查檔案副檔名
    const ext = path.extname(filePath);
    if (!config.detection.fileExtensions.includes(ext)) {
        return false;
    }
    
    // 檢查排除模式
    const normalizedPath = filePath.replace(/\\\\/g, '/');
    return !config.detection.excludePatterns.some(pattern => 
        normalizedPath.includes(pattern)
    );
}

// 檢查單個檔案
async function checkFile(filePath) {
    const issues = [];
    
    try {
        // 檢查檔案大小
        const stats = fs.statSync(filePath);
        if (stats.size > config.detection.maxFileSize) {
            return issues; // 跳過大檔案
        }
        
        // 讀取檔案內容 (強制 UTF-8 編碼)
        const content = fs.readFileSync(filePath, { encoding: 'utf8' });
        const lines = content.split('\n');
        
        // 檢查每種模式
        for (const [type, patterns] of Object.entries(config.patterns)) {
            if (type === 'excludeValues') continue; // 跳過排除值
            
            for (const pattern of patterns) {
                const regex = new RegExp(pattern, 'gi');
                
                lines.forEach((line, index) => {
                    const matches = line.match(regex);
                    if (matches) {
                        matches.forEach(match => {
                            // 檢查是否在排除清單中
                            if (config.patterns.excludeValues && 
                                config.patterns.excludeValues.some(excludeValue => 
                                    match.toLowerCase().includes(excludeValue.toLowerCase()))) {
                                return; // 跳過排除的值
                            }
                            
                            // 檢查是否是註解或測試程式碼
                            const trimmedLine = line.trim();
                            if (trimmedLine.startsWith('//') || 
                                trimmedLine.startsWith('#') ||
                                trimmedLine.includes('TODO') ||
                                trimmedLine.includes('FIXME') ||
                                trimmedLine.includes('example') ||
                                trimmedLine.includes('test')) {
                                return; // 跳過註解和測試程式碼
                            }
                            
                            issues.push({
                                file: filePath,
                                line: index + 1,
                                column: line.indexOf(match) + 1,
                                type: type,
                                pattern: pattern,
                                match: match,
                                message: `可能的${getTypeDescription(type)}: ${match.substring(0, 50)}${match.length > 50 ? '...' : ''}`,
                                severity: getSeverity(type)
                            });
                        });
                    }
                });
            }
        }
        
    } catch (error) {
        console.error(`[Token Monitor] 檢查檔案失敗 ${filePath}:`, error.message);
    }
    
    return issues;
}

// 獲取類型描述
function getTypeDescription(type) {
    const descriptions = {
        apiKeys: 'API 金鑰',
        tokens: 'Token 或存取權杖',
        secrets: '機密資訊', 
        credentials: '認證資訊',
        urls: '敏感 URL'
    };
    return descriptions[type] || type;
}

// 獲取嚴重性等級
function getSeverity(type) {
    const severities = {
        apiKeys: 'high',
        tokens: 'high',
        secrets: 'high',
        credentials: 'medium',
        urls: 'low'
    };
    return severities[type] || 'medium';
}

// 顯示通知
function showNotification(issues) {
    const highSeverityCount = issues.filter(i => i.severity === 'high').length;
    const message = highSeverityCount > 0 
        ? `發現 ${highSeverityCount} 個高風險安全問題`
        : `發現 ${issues.length} 個潛在安全問題`;
    
    console.log(`🚨 [Token Monitor] 通知: ${message}`);
}

// 高亮顯示問題
function highlightIssues(context, issues) {
    issues.forEach(issue => {
        console.log(`🎯 [Token Monitor] 高亮: ${issue.file}:${issue.line}:${issue.column} - ${issue.message}`);
    });
}

// 執行主函數
if (require.main === module) {
    // 測試模式
    const testContext = {
        filePath: process.argv[2] || 'test.go'
    };
    main(testContext);
}

module.exports = { main, checkFile };