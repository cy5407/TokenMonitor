/**
 * Token 監控系統測試腳本
 * 用於驗證改進後的監控功能
 */

const fs = require('fs');
const path = require('path');

// 引入監控模組
const tokenMonitor = require('./token-monitor-integration.js');
const { ManualTokenCalculator } = require('./.kiro/hooks/manual-token-calc.js');

class TokenMonitoringTest {
    constructor() {
        this.testResults = [];
        this.logFile = 'data/kiro-usage.log';
    }
    
    /**
     * 運行所有測試
     */
    async runAllTests() {
        console.log('🧪 開始 Token 監控系統測試...\n');
        
        // 備份現有日誌
        await this.backupExistingLog();
        
        try {
            // 測試1: 聊天對話監控
            await this.testChatMonitoring();
            
            // 測試2: 工具調用監控
            await this.testToolExecutionMonitoring();
            
            // 測試3: 代理任務監控
            await this.testAgentTaskMonitoring();
            
            // 測試4: 檔案操作監控
            await this.testFileOperationMonitoring();
            
            // 測試5: 對話回合監控
            await this.testConversationTurnMonitoring();
            
            // 測試6: 手動分析器
            await this.testManualAnalyzer();
            
            // 顯示測試結果
            this.displayTestResults();
            
        } finally {
            // 恢復原始日誌
            await this.restoreOriginalLog();
        }
    }
    
    /**
     * 備份現有日誌
     */
    async backupExistingLog() {
        if (fs.existsSync(this.logFile)) {
            const backupFile = this.logFile + '.backup';
            fs.copyFileSync(this.logFile, backupFile);
            console.log(`📋 已備份現有日誌到 ${backupFile}`);
        }
        
        // 清空日誌檔案以進行測試
        fs.writeFileSync(this.logFile, '', 'utf8');
    }
    
    /**
     * 恢復原始日誌
     */
    async restoreOriginalLog() {
        const backupFile = this.logFile + '.backup';
        if (fs.existsSync(backupFile)) {
            fs.copyFileSync(backupFile, this.logFile);
            fs.unlinkSync(backupFile);
            console.log(`📋 已恢復原始日誌檔案`);
        }
    }
    
    /**
     * 測試聊天對話監控
     */
    async testChatMonitoring() {
        console.log('🧪 測試1: 聊天對話監控');
        
        const testCases = [
            {
                event: 'chat.message.sent',
                data: {
                    content: '請幫我寫一個 JavaScript 函數來計算費波納契數列',
                    sessionId: 'test-session-1'
                }
            },
            {
                event: 'chat.message.received',
                data: {
                    content: '好的，我來幫你寫一個計算費波納契數列的函數：\n\n```javascript\nfunction fibonacci(n) {\n    if (n <= 1) return n;\n    return fibonacci(n - 1) + fibonacci(n - 2);\n}\n```',
                    sessionId: 'test-session-1',
                    model: 'claude-sonnet-4.0'
                }
            }
        ];
        
        for (const testCase of testCases) {
            try {
                const result = await tokenMonitor.execute(testCase);
                this.recordTestResult('聊天對話監控', testCase.event, result.success, result.data);
            } catch (error) {
                this.recordTestResult('聊天對話監控', testCase.event, false, error.message);
            }
        }
        
        console.log('✅ 聊天對話監控測試完成\n');
    }
    
    /**
     * 測試工具調用監控
     */
    async testToolExecutionMonitoring() {
        console.log('🧪 測試2: 工具調用監控');
        
        const testCases = [
            {
                event: 'tool.fsWrite',
                data: {
                    path: 'fibonacci.js',
                    text: 'function fibonacci(n) {\n    if (n <= 1) return n;\n    return fibonacci(n - 1) + fibonacci(n - 2);\n}\n\nmodule.exports = fibonacci;'
                }
            },
            {
                event: 'tool.strReplace',
                data: {
                    path: 'fibonacci.js',
                    oldStr: 'function fibonacci(n) {',
                    newStr: 'function fibonacci(n) {\n    // 計算費波納契數列'
                }
            },
            {
                event: 'tool.fsAppend',
                data: {
                    path: 'fibonacci.js',
                    text: '\n// 測試函數\nconsole.log(fibonacci(10));'
                }
            }
        ];
        
        for (const testCase of testCases) {
            try {
                const result = await tokenMonitor.execute(testCase);
                this.recordTestResult('工具調用監控', testCase.event, result.success, result.data);
            } catch (error) {
                this.recordTestResult('工具調用監控', testCase.event, false, error.message);
            }
        }
        
        console.log('✅ 工具調用監控測試完成\n');
    }
    
    /**
     * 測試代理任務監控
     */
    async testAgentTaskMonitoring() {
        console.log('🧪 測試3: 代理任務監控');
        
        const testCases = [
            {
                event: 'agent.codeGeneration',
                data: {
                    description: '生成 JavaScript 函數',
                    generatedContent: 'function calculateSum(arr) {\n    return arr.reduce((sum, num) => sum + num, 0);\n}',
                    sessionId: 'agent-session-1'
                }
            },
            {
                event: 'agent.documentGeneration',
                data: {
                    description: '生成 README 文件',
                    generatedContent: '# 專案說明\n\n這是一個 JavaScript 工具函數庫。\n\n## 安裝\n\n```bash\nnpm install\n```',
                    sessionId: 'agent-session-1'
                }
            }
        ];
        
        for (const testCase of testCases) {
            try {
                const result = await tokenMonitor.execute(testCase);
                this.recordTestResult('代理任務監控', testCase.event, result.success, result.data);
            } catch (error) {
                this.recordTestResult('代理任務監控', testCase.event, false, error.message);
            }
        }
        
        console.log('✅ 代理任務監控測試完成\n');
    }
    
    /**
     * 測試檔案操作監控
     */
    async testFileOperationMonitoring() {
        console.log('🧪 測試4: 檔案操作監控');
        
        const testCases = [
            {
                event: 'file.saved',
                data: {
                    filePath: 'package.json',
                    content: '{\n  "name": "test-project",\n  "version": "1.0.0",\n  "main": "index.js"\n}'
                }
            },
            {
                event: 'file.created',
                data: {
                    filePath: 'README.md',
                    content: '# 測試專案\n\n這是一個測試專案的說明文件。'
                }
            }
        ];
        
        for (const testCase of testCases) {
            try {
                const result = await tokenMonitor.execute(testCase);
                this.recordTestResult('檔案操作監控', testCase.event, result.success, result.data);
            } catch (error) {
                this.recordTestResult('檔案操作監控', testCase.event, false, error.message);
            }
        }
        
        console.log('✅ 檔案操作監控測試完成\n');
    }
    
    /**
     * 測試對話回合監控
     */
    async testConversationTurnMonitoring() {
        console.log('🧪 測試5: 對話回合監控');
        
        const testCase = {
            event: 'kiro.conversation.turn',
            data: {
                userInput: '請解釋什麼是遞迴函數',
                assistantOutput: '遞迴函數是一種調用自身的函數。它通常包含兩個部分：\n1. 基礎情況（Base Case）：停止遞迴的條件\n2. 遞迴情況（Recursive Case）：函數調用自身的部分\n\n例如計算階乘的遞迴函數：\n```javascript\nfunction factorial(n) {\n    if (n <= 1) return 1; // 基礎情況\n    return n * factorial(n - 1); // 遞迴情況\n}\n```',
                sessionId: 'conversation-session-1',
                toolsUsed: ['explanation', 'codeExample'],
                executionTime: 2500
            }
        };
        
        try {
            const result = await tokenMonitor.execute(testCase);
            this.recordTestResult('對話回合監控', testCase.event, result.success, result.data);
        } catch (error) {
            this.recordTestResult('對話回合監控', testCase.event, false, error.message);
        }
        
        console.log('✅ 對話回合監控測試完成\n');
    }
    
    /**
     * 測試手動分析器
     */
    async testManualAnalyzer() {
        console.log('🧪 測試6: 手動分析器');
        
        try {
            const calculator = new ManualTokenCalculator();
            const results = await calculator.execute();
            
            this.recordTestResult('手動分析器', 'analyze', true, {
                totalRecords: results.totalRecords,
                totalTokens: results.totalTokens,
                totalCost: results.totalCost
            });
            
        } catch (error) {
            this.recordTestResult('手動分析器', 'analyze', false, error.message);
        }
        
        console.log('✅ 手動分析器測試完成\n');
    }
    
    /**
     * 記錄測試結果
     */
    recordTestResult(category, testName, success, data) {
        this.testResults.push({
            category,
            testName,
            success,
            data,
            timestamp: new Date().toISOString()
        });
        
        const status = success ? '✅' : '❌';
        console.log(`  ${status} ${testName}: ${success ? '成功' : '失敗'}`);
        
        if (success && data && data.tokens) {
            console.log(`     Token: ${data.tokens}, 活動類型: ${data.activity_type}`);
        }
    }
    
    /**
     * 顯示測試結果
     */
    displayTestResults() {
        console.log('\n📊 ===== 測試結果摘要 =====');
        
        const categories = {};
        let totalTests = 0;
        let passedTests = 0;
        
        this.testResults.forEach(result => {
            if (!categories[result.category]) {
                categories[result.category] = { total: 0, passed: 0 };
            }
            categories[result.category].total++;
            totalTests++;
            
            if (result.success) {
                categories[result.category].passed++;
                passedTests++;
            }
        });
        
        console.log(`📈 總測試數: ${totalTests}`);
        console.log(`✅ 通過測試: ${passedTests}`);
        console.log(`❌ 失敗測試: ${totalTests - passedTests}`);
        console.log(`📊 通過率: ${((passedTests / totalTests) * 100).toFixed(1)}%`);
        
        console.log('\n📋 各類別測試結果:');
        Object.entries(categories).forEach(([category, stats]) => {
            const rate = ((stats.passed / stats.total) * 100).toFixed(1);
            console.log(`  ${category}: ${stats.passed}/${stats.total} (${rate}%)`);
        });
        
        // 檢查日誌檔案
        if (fs.existsSync(this.logFile)) {
            const logContent = fs.readFileSync(this.logFile, 'utf8');
            const logLines = logContent.trim().split('\n').filter(line => line.trim());
            console.log(`\n📄 生成日誌記錄: ${logLines.length} 筆`);
        }
        
        console.log('\n✅ 測試完成！');
    }
}

// 如果直接執行此檔案
if (require.main === module) {
    const test = new TokenMonitoringTest();
    test.runAllTests().catch(error => {
        console.error('❌ 測試執行失敗:', error);
        process.exit(1);
    });
}

module.exports = { TokenMonitoringTest };