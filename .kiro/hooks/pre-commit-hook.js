// Token Monitor - 提交前安全檢查 Hook
// 檢查暫存區檔案中的機密資訊、API Keys 和 Token 洩漏

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// 設定控制台編碼 (Windows)
if (process.platform === 'win32') {
    process.stdout.setEncoding('utf8');
    process.stderr.setEncoding('utf8');
}

// 從配置檔案載入設定
const config = JSON.parse(fs.readFileSync(path.join(__dirname, 'pre-commit-hook.json'), 'utf8')).config;

// 主要執行函數
async function main() {
    try {
        console.log('🔍 [Token Monitor] Security Check Starting...');

        // 獲取暫存區檔案
        const stagedFiles = getStagedFiles();
        if (stagedFiles.length === 0) {
            console.log('✅ [Token Monitor] No files to check');
            return;
        }

        console.log(`📁 [Token Monitor] Checking ${stagedFiles.length} files...`);

        // 檢查每個檔案
        const issues = [];
        for (const file of stagedFiles) {
            const fileIssues = await checkFile(file);
            if (fileIssues.length > 0) {
                issues.push(...fileIssues);
            }
        }

        // 生成報告
        if (config.reporting.enabled) {
            generateReport(issues, stagedFiles);
        }

        // 處理結果
        if (issues.length > 0) {
            console.log(`❌ [Token Monitor] Found ${issues.length} security issues:`);
            
            // 按嚴重性分組顯示
            const groupedIssues = groupIssuesBySeverity(issues);
            
            if (groupedIssues.high && groupedIssues.high.length > 0) {
                console.log(`\n🚨 HIGH SEVERITY (${groupedIssues.high.length}):`);
                groupedIssues.high.forEach(issue => {
                    console.log(`   ${issue.file}:${issue.line} - ${issue.type}: ${issue.message}`);
                });
            }
            
            if (groupedIssues.medium && groupedIssues.medium.length > 0) {
                console.log(`\n⚠️  MEDIUM SEVERITY (${groupedIssues.medium.length}):`);
                groupedIssues.medium.forEach(issue => {
                    console.log(`   ${issue.file}:${issue.line} - ${issue.type}: ${issue.message}`);
                });
            }
            
            if (groupedIssues.low && groupedIssues.low.length > 0) {
                console.log(`\nℹ️  LOW SEVERITY (${groupedIssues.low.length}):`);
                groupedIssues.low.forEach(issue => {
                    console.log(`   ${issue.file}:${issue.line} - ${issue.type}: ${issue.message}`);
                });
            }
            
            console.log('\n💡 建議:');
            console.log('   - 移除或加密敏感資訊');
            console.log('   - 使用環境變數或配置檔案');
            console.log('   - 將敏感檔案加入 .gitignore');
            console.log('   - 或使用 --no-verify 跳過檢查 (不建議)');
            
            // 只有高嚴重性問題才阻止提交
            const highSeverityCount = groupedIssues.high ? groupedIssues.high.length : 0;
            if (highSeverityCount > 0) {
                console.log(`\n🛑 發現 ${highSeverityCount} 個高風險問題，阻止提交`);
                process.exit(1);
            } else {
                console.log('\n⚠️  發現中低風險問題，但允許提交');
            }
        } else {
            console.log('✅ [Token Monitor] Security check passed');
        }

    } catch (error) {
        console.error('❌ [Token Monitor] 安全檢查執行失敗:', error.message);
        process.exit(1);
    }
}

// 獲取暫存區檔案
function getStagedFiles() {
    try {
        const output = execSync('git diff --cached --name-only --diff-filter=ACM', { encoding: 'utf8' });
        return output.trim().split('\n').filter(file => {
            if (!file) return false;

            // 檢查檔案副檔名
            const ext = path.extname(file);
            if (!config.detection.fileExtensions.includes(ext)) return false;

            // 檢查排除模式
            return !config.detection.excludePatterns.some(pattern =>
                file.includes(pattern.replace('/', path.sep))
            );
        });
    } catch (error) {
        console.error('[Token Monitor] 無法獲取暫存區檔案:', error.message);
        return [];
    }
}

// 檢查單個檔案
async function checkFile(filePath) {
    const issues = [];

    try {
        // 檢查檔案大小
        const stats = fs.statSync(filePath);
        if (stats.size > config.detection.maxFileSize) {
            console.log(`⚠️  [Token Monitor] 跳過大檔案: ${filePath} (${stats.size} bytes)`);
            return issues;
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
                                filePath.includes('test') ||
                                filePath.includes('_test.go')) {
                                return; // 跳過註解和測試程式碼
                            }

                            issues.push({
                                file: filePath,
                                line: index + 1,
                                type: type,
                                pattern: pattern,
                                match: match,
                                message: `Found potential ${getTypeDescription(type)}: ${match.substring(0, 50)}${match.length > 50 ? '...' : ''}`,
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
        apiKeys: 'API Key',
        tokens: 'Token/Access Token',
        secrets: 'Secret',
        credentials: 'Credential',
        urls: 'Sensitive URL'
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

// 按嚴重性分組問題
function groupIssuesBySeverity(issues) {
    return issues.reduce((groups, issue) => {
        const severity = issue.severity;
        if (!groups[severity]) {
            groups[severity] = [];
        }
        groups[severity].push(issue);
        return groups;
    }, {});
}

// 生成報告
function generateReport(issues, checkedFiles) {
    const report = {
        timestamp: new Date().toISOString(),
        project: 'Token Monitor',
        summary: {
            totalFiles: checkedFiles.length,
            issuesFound: issues.length,
            status: issues.length > 0 ? 'ISSUES_FOUND' : 'PASSED',
            severityBreakdown: {
                high: issues.filter(i => i.severity === 'high').length,
                medium: issues.filter(i => i.severity === 'medium').length,
                low: issues.filter(i => i.severity === 'low').length
            }
        },
        files: checkedFiles,
        issues: issues.map(issue => ({
            ...issue,
            // 截斷敏感內容以避免在報告中洩漏
            match: issue.match.substring(0, 20) + '...'
        }))
    };

    try {
        // 確保報告目錄存在
        const reportDir = path.dirname(config.reporting.outputPath);
        if (!fs.existsSync(reportDir)) {
            fs.mkdirSync(reportDir, { recursive: true });
        }

        fs.writeFileSync(config.reporting.outputPath, JSON.stringify(report, null, 2), { encoding: 'utf8' });
        console.log(`📊 [Token Monitor] 報告已生成: ${config.reporting.outputPath}`);
    } catch (error) {
        console.error('[Token Monitor] 生成報告失敗:', error.message);
    }
}

// 執行主函數
if (require.main === module) {
    main();
}

module.exports = { main, checkFile };