/**
 * 診斷現有監控系統的問題
 */

const fs = require('fs');
const path = require('path');

class MonitorDiagnosis {
    constructor() {
        this.logFile = 'data/kiro-usage.log';
    }
    
    async diagnose() {
        console.log('🔍 診斷現有監控系統...\n');
        
        // 1. 檢查日誌檔案
        await this.checkLogFile();
        
        // 2. 分析記錄類型
        await this.analyzeRecordTypes();
        
        // 3. 檢查 AI 對話記錄
        await this.checkAIConversations();
        
        // 4. 檢查監控覆蓋範圍
        await this.checkMonitoringCoverage();
        
        // 5. 提出改進建議
        await this.suggestImprovements();
    }
    
    async checkLogFile() {
        console.log('📄 檢查日誌檔案狀態:');
        
        if (!fs.existsSync(this.logFile)) {
            console.log('❌ 日誌檔案不存在');
            return;
        }
        
        const stats = fs.statSync(this.logFile);
        const content = fs.readFileSync(this.logFile, 'utf8');
        const lines = content.trim().split('\n').filter(line => line.trim());
        
        console.log(`✅ 檔案大小: ${(stats.size / 1024).toFixed(2)} KB`);
        console.log(`✅ 記錄總數: ${lines.length}`);
        console.log(`✅ 最後修改: ${stats.mtime.toLocaleString('zh-TW')}\n`);
    }
    
    async analyzeRecordTypes() {
        console.log('📊 分析記錄類型分布:');
        
        if (!fs.existsSync(this.logFile)) {
            return;
        }
        
        const content = fs.readFileSync(this.logFile, 'utf8');
        const lines = content.trim().split('\n').filter(line => line.trim());
        
        const eventTypes = {};
        const sources = {};
        const activities = {};
        
        for (const line of lines) {
            try {
                const record = JSON.parse(line);
                
                // 統計事件類型
                const eventType = record.event || 'unknown';
                eventTypes[eventType] = (eventTypes[eventType] || 0) + 1;
                
                // 統計來源
                const source = record.model || record.ide_detected || 'unknown';
                sources[source] = (sources[source] || 0) + 1;
                
                // 統計活動類型
                const activity = record.activity_type || 'unknown';
                activities[activity] = (activities[activity] || 0) + 1;
                
            } catch (error) {
                // 忽略解析錯誤
            }
        }
        
        console.log('事件類型分布:');
        Object.entries(eventTypes).forEach(([type, count]) => {
            console.log(`  ${type}: ${count} 次`);
        });
        
        console.log('\n來源分布:');
        Object.entries(sources).forEach(([source, count]) => {
            console.log(`  ${source}: ${count} 次`);
        });
        
        console.log('\n活動類型分布:');
        Object.entries(activities).forEach(([activity, count]) => {
            console.log(`  ${activity}: ${count} 次`);
        });
        
        console.log('');
    }
    
    async checkAIConversations() {
        console.log('💬 檢查 AI 對話記錄:');
        
        if (!fs.existsSync(this.logFile)) {
            return;
        }
        
        const content = fs.readFileSync(this.logFile, 'utf8');
        const lines = content.trim().split('\n').filter(line => line.trim());
        
        let chatMessages = 0;
        let realAIInteractions = 0;
        let terminalCommands = 0;
        
        for (const line of lines) {
            try {
                const record = JSON.parse(line);
                
                if (record.event === 'chat_message') {
                    chatMessages++;
                    
                    // 檢查是否為真實的 AI 互動
                    if (record.session_id && !record.session_id.includes('test')) {
                        realAIInteractions++;
                    }
                }
                
                if (record.event === 'command_executed') {
                    terminalCommands++;
                }
                
            } catch (error) {
                // 忽略解析錯誤
            }
        }
        
        console.log(`💬 聊天訊息總數: ${chatMessages}`);
        console.log(`🤖 真實 AI 互動: ${realAIInteractions}`);
        console.log(`💻 終端命令: ${terminalCommands}`);
        
        if (realAIInteractions === 0) {
            console.log('❌ 問題：沒有檢測到真實的 AI 對話！');
        }
        
        console.log('');
    }
    
    async checkMonitoringCoverage() {
        console.log('🔍 檢查監控覆蓋範圍:');
        
        // 檢查是否有監控程序在運行
        const { execSync } = require('child_process');
        
        try {
            const processes = execSync('tasklist /FI "IMAGENAME eq node.exe" /FO CSV', { encoding: 'utf8' });
            const nodeProcesses = processes.split('\n').filter(line => line.includes('node.exe')).length - 1;
            
            console.log(`🟢 Node.js 程序數量: ${nodeProcesses}`);
            
            if (nodeProcesses === 0) {
                console.log('❌ 問題：沒有 Node.js 監控程序在運行！');
            }
            
        } catch (error) {
            console.log('⚠️ 無法檢查程序狀態');
        }
        
        // 檢查監控的檔案類型
        console.log('📁 當前監控的檔案類型: .md, .txt, .js, .ts, .py, .java, .cpp, .html, .css, .json');
        
        // 檢查是否排除了重要目錄
        console.log('🚫 排除的目錄: node_modules, .git, data');
        
        console.log('');
    }
    
    async suggestImprovements() {
        console.log('💡 改進建議:');
        
        const suggestions = [
            '1. 啟動檔案監控系統：node scripts/Universal-Token-Monitor.js',
            '2. 檢查 Kiro IDE Hook 系統是否正常運作',
            '3. 測試在 VS Code 中使用 AI 助手時是否有檔案變化',
            '4. 改進 AI 內容檢測邏輯，更準確識別 AI 生成的內容',
            '5. 加入 VS Code 擴展監控（如 GitHub Copilot）',
            '6. 改進 Kiro IDE 整合，確保能捕捉到真實對話'
        ];
        
        suggestions.forEach(suggestion => {
            console.log(`  ${suggestion}`);
        });
        
        console.log('\n🎯 核心問題：');
        console.log('  - 監控系統主要記錄終端命令，而不是真正的 AI 對話');
        console.log('  - 需要改進檢測邏輯，識別 AI 生成的內容');
        console.log('  - 需要測試和驗證不同 IDE 中的 AI 使用情況');
    }
}

// 執行診斷
if (require.main === module) {
    const diagnosis = new MonitorDiagnosis();
    diagnosis.diagnose().catch(error => {
        console.error('診斷失敗:', error);
    });
}

module.exports = { MonitorDiagnosis };