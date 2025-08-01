/**
 * Kiro Chat Token 計算器 - 手動觸發版本
 * 分析 Kiro Chat 使用記錄並提供詳細的 Token 消耗報告
 */

const fs = require('fs');
const path = require('path');

class ManualTokenCalculator {
    constructor() {
        this.logFile = 'data/kiro-usage.log';
        this.results = {
            totalRecords: 0,
            inputTokens: 0,
            outputTokens: 0,
            totalTokens: 0,
            totalCost: 0,
            activityStats: {},
            sessionStats: {},
            modelStats: {},
            eventStats: {}
        };
    }

    async execute() {
        try {
            console.log('🚀 開始執行 Token 使用量分析...');

            if (!fs.existsSync(this.logFile)) {
                throw new Error(`找不到使用記錄檔案: ${this.logFile}`);
            }

            const logData = fs.readFileSync(this.logFile, 'utf8');
            const lines = logData.trim().split('\n').filter(line => line.trim());

            console.log(`📊 分析 ${lines.length} 筆記錄...`);

            this.results.totalRecords = lines.length;

            for (const line of lines) {
                try {
                    const record = JSON.parse(line);
                    this.processRecord(record);
                } catch (parseError) {
                    console.warn(`⚠️ 無法解析記錄: ${line}`);
                }
            }

            this.generateReport();
            return this.results;

        } catch (error) {
            console.error('❌ Token 分析執行失敗:', error.message);
            throw error;
        }
    }

    processRecord(record) {
        const tokens = record.tokens || 0;

        // 統計總 Token
        this.results.totalTokens += tokens;

        // 區分輸入/輸出 Token - 修正邏輯
        if (record.direction === 'sent' || 
            (record.event === 'chat_message' && record.direction === 'sent') ||
            record.event === 'user_input') {
            // 用戶發送的消息 = 輸入 token
            this.results.inputTokens += tokens;
        } else if (record.direction === 'received' || 
                   (record.event === 'chat_message' && record.direction === 'received')) {
            // AI 回應的消息 = 輸出 token
            this.results.outputTokens += tokens;
        } else if (record.event === 'file_save' || 
                   record.event === 'file_edit' || 
                   record.event === 'terminal_command') {
            // 文件操作和終端命令通常不區分輸入輸出，計入輸入
            this.results.inputTokens += tokens;
        }

        // 處理對話回合的特殊情況
        if (record.event === 'conversation_turn') {
            this.results.inputTokens += record.input_tokens || 0;
            this.results.outputTokens += record.output_tokens || 0;
            // 避免重複計算總 tokens
            this.results.totalTokens -= tokens;
            this.results.totalTokens += (record.input_tokens || 0) + (record.output_tokens || 0);
        }

        // 成本分析（支援多種成本格式）
        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost ||
                record.cost_analysis.cost_usd ||
                record.cost_analysis.input_cost + record.cost_analysis.output_cost ||
                0;
            this.results.totalCost += cost;
        }

        // 活動類型統計
        const activity = record.activity_type || 'unknown';
        if (!this.results.activityStats[activity]) {
            this.results.activityStats[activity] = { count: 0, tokens: 0, cost: 0 };
        }
        this.results.activityStats[activity].count++;
        this.results.activityStats[activity].tokens += tokens;
        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost || record.cost_analysis.cost_usd || 0;
            this.results.activityStats[activity].cost += cost;
        }

        // 會話統計
        if (record.session_id) {
            if (!this.results.sessionStats[record.session_id]) {
                this.results.sessionStats[record.session_id] = {
                    inputTokens: 0,
                    outputTokens: 0,
                    totalTokens: 0,
                    cost: 0,
                    events: 0,
                    lastActivity: record.timestamp
                };
            }

            const session = this.results.sessionStats[record.session_id];
            session.events++;
            session.lastActivity = record.timestamp;

            if (record.event === 'conversation_turn') {
                session.inputTokens += record.input_tokens || 0;
                session.outputTokens += record.output_tokens || 0;
                session.totalTokens += (record.input_tokens || 0) + (record.output_tokens || 0);
            } else {
                session.totalTokens += tokens;
                if (record.direction === 'sent') {
                    session.inputTokens += tokens;
                } else if (record.direction === 'received' ||
                    ['tool_execution', 'agent_task', 'file_save'].includes(record.event)) {
                    session.outputTokens += tokens;
                }
            }

            if (record.cost_analysis) {
                const cost = record.cost_analysis.total_cost || record.cost_analysis.cost_usd || 0;
                session.cost += cost;
            }
        }

        // 模型統計
        const model = record.model || 'unknown';
        if (!this.results.modelStats[model]) {
            this.results.modelStats[model] = {
                inputTokens: 0,
                outputTokens: 0,
                totalTokens: 0,
                cost: 0,
                events: 0
            };
        }

        const modelStat = this.results.modelStats[model];
        modelStat.events++;

        if (record.event === 'conversation_turn') {
            modelStat.inputTokens += record.input_tokens || 0;
            modelStat.outputTokens += record.output_tokens || 0;
            modelStat.totalTokens += (record.input_tokens || 0) + (record.output_tokens || 0);
        } else {
            modelStat.totalTokens += tokens;
            if (record.direction === 'sent') {
                modelStat.inputTokens += tokens;
            } else if (record.direction === 'received' ||
                ['tool_execution', 'agent_task', 'file_save'].includes(record.event)) {
                modelStat.outputTokens += tokens;
            }
        }

        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost || record.cost_analysis.cost_usd || 0;
            modelStat.cost += cost;
        }

        // 事件類型統計
        const eventType = record.event || 'unknown';
        if (!this.results.eventStats) {
            this.results.eventStats = {};
        }
        if (!this.results.eventStats[eventType]) {
            this.results.eventStats[eventType] = { count: 0, tokens: 0, cost: 0 };
        }
        this.results.eventStats[eventType].count++;
        this.results.eventStats[eventType].tokens += tokens;
        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost || record.cost_analysis.cost_usd || 0;
            this.results.eventStats[eventType].cost += cost;
        }
    }

    generateReport() {
        console.log('\n📊 ===== Kiro Chat Token 使用分析報告 =====');
        console.log(`📈 總記錄數: ${this.results.totalRecords}`);
        console.log(`🔢 輸入 Token: ${this.results.inputTokens}`);
        console.log(`🔢 輸出 Token: ${this.results.outputTokens}`);
        console.log(`🔢 總 Token: ${this.results.totalTokens}`);
        console.log(`💰 預估成本: $${this.results.totalCost.toFixed(6)} USD`);

        console.log('\n📋 活動類型統計:');
        Object.entries(this.results.activityStats).forEach(([activity, stats]) => {
            console.log(`  ${activity}: ${stats.count} 次, ${stats.tokens} tokens`);
        });

        console.log('\n💡 會話統計 (前5個):');
        const topSessions = Object.entries(this.results.sessionStats)
            .sort(([, a], [, b]) => b.totalTokens - a.totalTokens)
            .slice(0, 5);

        topSessions.forEach(([sessionId, stats]) => {
            console.log(`  ${sessionId}: ${stats.totalTokens} tokens ($${stats.cost.toFixed(6)})`);
        });

        console.log('\n🤖 模型使用統計:');
        Object.entries(this.results.modelStats).forEach(([model, stats]) => {
            console.log(`  ${model}: 輸入${stats.inputTokens} + 輸出${stats.outputTokens} = ${stats.totalTokens} tokens ($${stats.cost.toFixed(6)})`);
        });

        console.log('\n📊 事件類型統計:');
        Object.entries(this.results.eventStats).forEach(([eventType, stats]) => {
            console.log(`  ${eventType}: ${stats.count} 次, ${stats.tokens} tokens (${stats.cost.toFixed(6)} USD)`);
        });

        console.log('\n✅ Token 分析完成！');
    }
}

// Hook 執行入口
async function main() {
    const calculator = new ManualTokenCalculator();
    try {
        const results = await calculator.execute();

        // 回傳結果給 Kiro
        return {
            success: true,
            message: 'Token 分析完成',
            data: results
        };
    } catch (error) {
        return {
            success: false,
            message: `Token 分析失敗: ${error.message}`,
            error: error.stack
        };
    }
}

// 如果直接執行此檔案
if (require.main === module) {
    main().then(result => {
        console.log('\n📋 執行結果:', JSON.stringify(result, null, 2));
    }).catch(error => {
        console.error('❌ 執行錯誤:', error);
        process.exit(1);
    });
}

module.exports = { ManualTokenCalculator, main };