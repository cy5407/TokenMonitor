/**
 * Kiro IDE Token Monitor Hook
 * 自動監控和記錄 Token 使用量
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

class TokenMonitorHook {
    constructor() {
        this.tokenMonitorPath = path.join(__dirname, '..', 'token-monitor.exe');
        this.dataDir = path.join(__dirname, '..', 'data');
        this.logFile = path.join(this.dataDir, 'token-usage.log');
        
        // 確保資料目錄存在
        if (!fs.existsSync(this.dataDir)) {
            fs.mkdirSync(this.dataDir, { recursive: true });
        }
    }

    /**
     * 主要執行函數
     */
    async execute(context) {
        try {
            const { event, data } = context;
            
            console.log(`🔍 Token Monitor Hook triggered by: ${event}`);
            
            switch (event) {
                case 'chat.message.sent':
                    await this.handleChatMessage(data, 'sent');
                    break;
                case 'chat.message.received':
                    await this.handleChatMessage(data, 'received');
                    break;
                case 'file.save':
                    await this.handleFileSave(data);
                    break;
                case 'agent.task.completed':
                    await this.handleTaskCompleted(data);
                    break;
                default:
                    console.log(`⚠️ Unknown event: ${event}`);
            }
            
            return { success: true, message: 'Token monitoring completed' };
            
        } catch (error) {
            console.error('❌ Token Monitor Hook error:', error);
            return { success: false, error: error.message };
        }
    }

    /**
     * 處理聊天訊息
     */
    async handleChatMessage(data, direction) {
        if (!data.content) return;
        
        const tokenCount = await this.calculateTokens(data.content);
        const activityType = this.classifyActivity(data.content);
        
        const record = {
            timestamp: new Date().toISOString(),
            event: 'chat_message',
            direction: direction,
            content_length: data.content.length,
            token_count: tokenCount,
            activity_type: activityType,
            session_id: data.sessionId || 'unknown',
            model: data.model || 'claude-sonnet-4.0'
        };
        
        await this.saveRecord(record);
        
        console.log(`📊 Chat ${direction}: ${tokenCount} tokens (${activityType})`);
    }

    /**
     * 處理檔案儲存
     */
    async handleFileSave(data) {
        if (!data.content || !data.filePath) return;
        
        const tokenCount = await this.calculateTokens(data.content);
        const activityType = this.classifyActivityByFile(data.filePath, data.content);
        
        const record = {
            timestamp: new Date().toISOString(),
            event: 'file_save',
            file_path: data.filePath,
            file_type: path.extname(data.filePath),
            content_length: data.content.length,
            token_count: tokenCount,
            activity_type: activityType
        };
        
        await this.saveRecord(record);
        
        console.log(`💾 File saved: ${tokenCount} tokens (${activityType})`);
    }

    /**
     * 處理任務完成
     */
    async handleTaskCompleted(data) {
        if (!data.description) return;
        
        const tokenCount = await this.calculateTokens(data.description);
        
        const record = {
            timestamp: new Date().toISOString(),
            event: 'task_completed',
            task_id: data.taskId,
            task_type: data.taskType || 'unknown',
            description_length: data.description.length,
            token_count: tokenCount,
            activity_type: 'task_management',
            duration: data.duration || 0
        };
        
        await this.saveRecord(record);
        
        console.log(`✅ Task completed: ${tokenCount} tokens`);
    }

    /**
     * 計算 Token 數量
     */
    async calculateTokens(text) {
        try {
            const command = `"${this.tokenMonitorPath}" calculate "${text.replace(/"/g, '\\"')}"`;
            const output = execSync(command, { encoding: 'utf8', timeout: 10000 });
            
            // 解析輸出中的 Token 數量
            const match = output.match(/Token 數量:\s*(\d+)/);
            return match ? parseInt(match[1]) : 0;
            
        } catch (error) {
            console.error('Token calculation error:', error.message);
            return this.estimateTokens(text);
        }
    }

    /**
     * 簡單的 Token 估算（備用方法）
     */
    estimateTokens(text) {
        // 簡單估算：英文 4 字符/token，中文 1.5 字符/token
        const englishChars = (text.match(/[a-zA-Z0-9\s]/g) || []).length;
        const chineseChars = text.length - englishChars;
        
        return Math.ceil(englishChars / 4 + chineseChars / 1.5);
    }

    /**
     * 活動類型分類
     */
    classifyActivity(content) {
        const patterns = {
            coding: /(?:function|class|implement|程式|函數|變數|程式碼|寫.*程式|實作.*功能)/i,
            debugging: /(?:error|bug|fix|錯誤|修復|除錯|修復.*問題|解決.*錯誤)/i,
            documentation: /(?:README|document|文件|說明|註解|更新.*文件|撰寫.*說明)/i,
            'spec-development': /(?:spec|requirement|design|需求|設計|規格|需求分析)/i,
            chat: /(?:chat|question|help|問題|協助|請問|如何)/i
        };

        for (const [type, pattern] of Object.entries(patterns)) {
            if (pattern.test(content)) {
                return type;
            }
        }
        
        return 'general';
    }

    /**
     * 根據檔案類型分類活動
     */
    classifyActivityByFile(filePath, content) {
        const ext = path.extname(filePath).toLowerCase();
        const fileName = path.basename(filePath).toLowerCase();
        
        if (['.js', '.ts', '.py', '.go', '.java', '.cpp', '.c', '.cs'].includes(ext)) {
            return 'coding';
        }
        
        if (['.md', '.txt', '.doc', '.docx'].includes(ext) || fileName.includes('readme')) {
            return 'documentation';
        }
        
        if (fileName.includes('spec') || fileName.includes('requirement')) {
            return 'spec-development';
        }
        
        return this.classifyActivity(content);
    }

    /**
     * 儲存記錄
     */
    async saveRecord(record) {
        try {
            const logEntry = JSON.stringify(record) + '\n';
            fs.appendFileSync(this.logFile, logEntry, 'utf8');
            
            // 同時儲存到 Token Monitor 的資料格式
            await this.saveToTokenMonitorFormat(record);
            
        } catch (error) {
            console.error('Failed to save record:', error);
        }
    }

    /**
     * 儲存為 Token Monitor 格式
     */
    async saveToTokenMonitorFormat(record) {
        try {
            const command = `"${this.tokenMonitorPath}" store-record '${JSON.stringify(record)}'`;
            execSync(command, { timeout: 5000 });
        } catch (error) {
            // 如果 Token Monitor 不支援此命令，忽略錯誤
            console.log('Note: Token Monitor store-record command not available');
        }
    }

    /**
     * 生成使用報告
     */
    async generateReport() {
        try {
            const command = `"${this.tokenMonitorPath}" report --format json --output "${path.join(this.dataDir, 'usage-report.json')}"`;
            execSync(command, { timeout: 30000 });
            console.log('📊 Usage report generated');
        } catch (error) {
            console.error('Failed to generate report:', error);
        }
    }
}

// Hook 執行入口
async function main(context) {
    const hook = new TokenMonitorHook();
    return await hook.execute(context);
}

// 如果直接執行此腳本，進行測試
if (require.main === module) {
    const testContext = {
        event: 'chat.message.sent',
        data: {
            content: '請幫我寫一個 JavaScript 函數來計算斐波那契數列',
            sessionId: 'test-session',
            model: 'claude-sonnet-4.0'
        }
    };
    
    main(testContext).then(result => {
        console.log('Test result:', result);
    }).catch(error => {
        console.error('Test error:', error);
    });
}

module.exports = { main };