/**
 * 全面 Token 監控器
 * 自動監控所有 Kiro 活動的 Token 消耗
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// 引入主要的監控整合腳本
const tokenMonitor = require('../../token-monitor-integration.js');

class ComprehensiveTokenMonitor {
    constructor() {
        this.logFile = 'data/kiro-usage.log';
        this.sessionId = this.generateSessionId();
        this.isMonitoring = true;
        this.stats = {
            totalEvents: 0,
            totalTokens: 0,
            totalCost: 0,
            sessionStart: new Date().toISOString()
        };
        
        // 確保資料目錄存在
        this.ensureDataDirectory();
        
        console.log(`🚀 全面 Token 監控器啟動 - Session: ${this.sessionId}`);
    }
    
    ensureDataDirectory() {
        const dataDir = path.dirname(this.logFile);
        if (!fs.existsSync(dataDir)) {
            fs.mkdirSync(dataDir, { recursive: true });
        }
    }
    
    generateSessionId() {
        const timestamp = Date.now();
        const random = Math.random().toString(36).substring(2, 8);
        return `monitor-${timestamp}-${random}`;
    }
    
    /**
     * 主要執行函數 - 由 Kiro Hook 系統調用
     */
    async execute(context) {
        try {
            if (!this.isMonitoring) {
                return { success: true, message: 'Monitoring is disabled' };
            }
            
            console.log('🔍 全面監控器接收事件:', JSON.stringify(context, null, 2));
            
            // 增強上下文資訊
            const enhancedContext = this.enhanceContext(context);
            
            // 使用主要監控整合腳本處理事件
            const result = await tokenMonitor.execute(enhancedContext);
            
            // 更新統計資訊
            if (result.success && result.data) {
                this.updateStats(result.data);
            }
            
            // 定期分析（每5分鐘）
            if (this.shouldRunAnalysis()) {
                await this.runPeriodicAnalysis();
            }
            
            return {
                success: true,
                message: 'Comprehensive monitoring completed',
                sessionId: this.sessionId,
                stats: this.stats,
                data: result.data
            };
            
        } catch (error) {
            console.error('❌ 全面監控器錯誤:', error.message);
            return {
                success: false,
                error: error.message,
                sessionId: this.sessionId
            };
        }
    }
    
    /**
     * 增強上下文資訊
     */
    enhanceContext(context) {
        const enhanced = {
            ...context,
            sessionId: this.sessionId,
            timestamp: context.timestamp || new Date().toISOString(),
            monitorVersion: '2.0.0'
        };
        
        // 根據不同的觸發類型增強上下文
        if (context.trigger) {
            switch (context.trigger.type) {
                case 'fileChange':
                    enhanced.event = this.mapFileEvent(context);
                    break;
                case 'chatEvent':
                    enhanced.event = this.mapChatEvent(context);
                    break;
                case 'toolExecution':
                    enhanced.event = this.mapToolEvent(context);
                    break;
                case 'agentActivity':
                    enhanced.event = this.mapAgentEvent(context);
                    break;
            }
        }
        
        // 如果沒有明確的事件類型，嘗試推斷
        if (!enhanced.event) {
            enhanced.event = this.inferEventType(context);
        }
        
        return enhanced;
    }
    
    /**
     * 映射檔案事件
     */
    mapFileEvent(context) {
        const eventMap = {
            'created': 'file.created',
            'modified': 'file.modified',
            'saved': 'file.saved'
        };
        
        return eventMap[context.trigger.event] || 'file.changed';
    }
    
    /**
     * 映射聊天事件
     */
    mapChatEvent(context) {
        const eventMap = {
            'message.sent': 'chat.message.sent',
            'message.received': 'chat.message.received',
            'conversation.turn': 'kiro.conversation.turn'
        };
        
        return eventMap[context.trigger.event] || 'chat.unknown';
    }
    
    /**
     * 映射工具事件
     */
    mapToolEvent(context) {
        const eventMap = {
            'tool.fsWrite': 'tool.fsWrite',
            'tool.fsAppend': 'tool.fsAppend',
            'tool.strReplace': 'tool.strReplace',
            'tool.executePwsh': 'tool.executePwsh',
            'tool.readFile': 'tool.readFile'
        };
        
        return eventMap[context.trigger.event] || 'tool.unknown';
    }
    
    /**
     * 映射代理事件
     */
    mapAgentEvent(context) {
        const eventMap = {
            'agent.codeGeneration': 'agent.codeGeneration',
            'agent.documentGeneration': 'agent.documentGeneration',
            'agent.taskExecution': 'agent.taskExecution'
        };
        
        return eventMap[context.trigger.event] || 'agent.unknown';
    }
    
    /**
     * 推斷事件類型
     */
    inferEventType(context) {
        // 根據上下文內容推斷事件類型
        if (context.filePath || context.file) {
            return 'file.saved';
        }
        
        if (context.message || context.content) {
            return 'chat.message.received';
        }
        
        if (context.tool || context.toolName) {
            return 'tool.execution';
        }
        
        if (context.task || context.taskType) {
            return 'agent.taskExecution';
        }
        
        return 'generic.event';
    }
    
    /**
     * 更新統計資訊
     */
    updateStats(data) {
        if (!data) return;
        
        this.stats.totalEvents++;
        
        if (data.tokens) {
            this.stats.totalTokens += data.tokens;
        }
        
        if (data.cost_analysis && data.cost_analysis.cost_usd) {
            this.stats.totalCost += data.cost_analysis.cost_usd;
        } else if (data.cost_analysis && data.cost_analysis.total_cost) {
            this.stats.totalCost += data.cost_analysis.total_cost;
        }
        
        // 每100個事件顯示一次統計
        if (this.stats.totalEvents % 100 === 0) {
            this.displayStats();
        }
    }
    
    /**
     * 顯示統計資訊
     */
    displayStats() {
        const uptime = new Date() - new Date(this.stats.sessionStart);
        const uptimeMinutes = Math.floor(uptime / 60000);
        
        console.log('\n📊 ===== 監控統計 =====');
        console.log(`⏱️ 運行時間: ${uptimeMinutes} 分鐘`);
        console.log(`📈 總事件數: ${this.stats.totalEvents}`);
        console.log(`🔢 總 Token: ${this.stats.totalTokens}`);
        console.log(`💰 總成本: ${this.stats.totalCost.toFixed(6)} USD`);
        console.log(`📊 平均每分鐘: ${(this.stats.totalEvents / Math.max(uptimeMinutes, 1)).toFixed(1)} 事件`);
        console.log('========================\n');
    }
    
    /**
     * 檢查是否應該運行分析
     */
    shouldRunAnalysis() {
        if (!this.lastAnalysis) {
            this.lastAnalysis = Date.now();
            return false;
        }
        
        const timeSinceLastAnalysis = Date.now() - this.lastAnalysis;
        const analysisInterval = 300000; // 5分鐘
        
        return timeSinceLastAnalysis >= analysisInterval;
    }
    
    /**
     * 運行定期分析
     */
    async runPeriodicAnalysis() {
        try {
            console.log('🔄 運行定期 Token 使用分析...');
            
            const analysis = await tokenMonitor.analyzeUsageLog();
            
            if (analysis.success) {
                console.log('✅ 定期分析完成');
                
                // 如果成本超過閾值，發送通知
                if (analysis.summary && analysis.summary.totalCost > 1.0) {
                    console.log(`⚠️ 成本警告: 總成本已達 ${analysis.summary.totalCost.toFixed(6)} USD`);
                }
            }
            
            this.lastAnalysis = Date.now();
            
        } catch (error) {
            console.error('❌ 定期分析失敗:', error.message);
        }
    }
    
    /**
     * 停止監控
     */
    stop() {
        this.isMonitoring = false;
        console.log('🛑 全面 Token 監控器已停止');
        this.displayStats();
    }
    
    /**
     * 重新啟動監控
     */
    restart() {
        this.isMonitoring = true;
        this.sessionId = this.generateSessionId();
        this.stats = {
            totalEvents: 0,
            totalTokens: 0,
            totalCost: 0,
            sessionStart: new Date().toISOString()
        };
        console.log(`🔄 全面 Token 監控器重新啟動 - Session: ${this.sessionId}`);
    }
}

// 全域監控器實例
let globalMonitor = null;

/**
 * Hook 執行入口
 */
async function main(context) {
    try {
        // 初始化全域監控器（如果尚未初始化）
        if (!globalMonitor) {
            globalMonitor = new ComprehensiveTokenMonitor();
        }
        
        // 執行監控
        const result = await globalMonitor.execute(context);
        
        return result;
        
    } catch (error) {
        console.error('❌ 全面監控器主函數錯誤:', error.message);
        return {
            success: false,
            error: error.message
        };
    }
}

/**
 * 停止監控
 */
function stopMonitoring() {
    if (globalMonitor) {
        globalMonitor.stop();
    }
}

/**
 * 重新啟動監控
 */
function restartMonitoring() {
    if (globalMonitor) {
        globalMonitor.restart();
    } else {
        globalMonitor = new ComprehensiveTokenMonitor();
    }
}

// 如果直接執行此檔案
if (require.main === module) {
    console.log('🧪 測試全面 Token 監控器...');
    
    // 測試不同類型的事件
    const testEvents = [
        {
            trigger: { type: 'fileChange', event: 'saved' },
            filePath: 'test.js',
            content: 'console.log("Hello World");'
        },
        {
            trigger: { type: 'chatEvent', event: 'message.sent' },
            message: '請幫我寫一個 JavaScript 函數'
        },
        {
            trigger: { type: 'toolExecution', event: 'tool.fsWrite' },
            tool: 'fsWrite',
            path: 'output.js',
            text: 'function test() { return "Hello"; }'
        }
    ];
    
    async function runTests() {
        for (const event of testEvents) {
            console.log(`\n🧪 測試事件: ${JSON.stringify(event, null, 2)}`);
            const result = await main(event);
            console.log(`✅ 結果: ${JSON.stringify(result, null, 2)}`);
        }
        
        // 停止監控
        stopMonitoring();
    }
    
    runTests().catch(error => {
        console.error('❌ 測試失敗:', error);
        process.exit(1);
    });
}

module.exports = {
    main,
    stopMonitoring,
    restartMonitoring,
    ComprehensiveTokenMonitor
};