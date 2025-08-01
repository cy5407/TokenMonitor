/**
 * Enhanced Token Monitor - Table Format Output
 * 專業表格格式的 Token 使用量報告工具
 */

const fs = require('fs');
const path = require('path');

class EnhancedTokenReporter {
    constructor() {
        this.logFile = 'data/kiro-usage.log';
        this.dailyStats = new Map();
        this.totalStats = {
            input: 0,
            output: 0,
            cacheCreate: 0,
            cacheRead: 0,
            totalTokens: 0,
            cost: 0
        };
    }

    async generateReport() {
        try {
            console.log('\n🚀 生成專業 Token 使用報告...\n');

            if (!fs.existsSync(this.logFile)) {
                throw new Error(`找不到使用記錄檔案: ${this.logFile}`);
            }

            const logData = fs.readFileSync(this.logFile, 'utf8');
            const lines = logData.trim().split('\n').filter(line => line.trim());

            // 處理每筆記錄
            for (const line of lines) {
                try {
                    const record = JSON.parse(line);
                    this.processRecord(record);
                } catch (parseError) {
                    console.warn(`⚠️ 無法解析記錄: ${line.substring(0, 50)}...`);
                }
            }

            this.printTableReport();
            return this.getTotalStats();

        } catch (error) {
            console.error('❌ 報告生成失敗:', error.message);
            throw error;
        }
    }

    processRecord(record) {
        const date = new Date(record.timestamp).toISOString().split('T')[0];
        const model = this.extractModelName(record.model);
        const tokens = record.tokens || 0;

        // 初始化日期統計
        if (!this.dailyStats.has(date)) {
            this.dailyStats.set(date, {
                date,
                model,
                input: 0,
                output: 0,
                cacheCreate: 0,
                cacheRead: 0,
                totalTokens: 0,
                cost: 0
            });
        }

        const dailyStat = this.dailyStats.get(date);

        // 分類 Token 類型
        if (record.direction === 'sent' || record.event === 'user_input') {
            dailyStat.input += tokens;
            this.totalStats.input += tokens;
        } else if (record.direction === 'received' || 
                   record.event === 'chat_message' && record.direction === 'received') {
            dailyStat.output += tokens;
            this.totalStats.output += tokens;
        } else if (record.event === 'cache_create' || 
                   (record.activity_type && record.activity_type.includes('cache'))) {
            dailyStat.cacheCreate += tokens;
            this.totalStats.cacheCreate += tokens;
        } else if (record.event === 'cache_read') {
            dailyStat.cacheRead += tokens;
            this.totalStats.cacheRead += tokens;
        } else {
            // 默認歸類為輸出
            dailyStat.output += tokens;
            this.totalStats.output += tokens;
        }

        dailyStat.totalTokens += tokens;
        this.totalStats.totalTokens += tokens;

        // 成本計算
        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost || 
                        record.cost_analysis.cost_usd || 0;
            dailyStat.cost += cost;
            this.totalStats.cost += cost;
        }
    }

    extractModelName(modelString) {
        if (!modelString || modelString === 'unknown') return 'sonnet-4';
        
        // 提取模型簡稱
        if (modelString.includes('claude-sonnet-4')) return 'sonnet-4';
        if (modelString.includes('claude-sonnet-3')) return 'sonnet-3';
        if (modelString.includes('claude-haiku')) return 'haiku';
        if (modelString.includes('gpt-4')) return 'gpt-4';
        if (modelString.includes('gpt-3')) return 'gpt-3';
        
        return 'sonnet-4'; // 默認值
    }

    printTableReport() {
        // 表格標題
        const title = "Claude Code Token Usage Report - Daily";
        console.log(`\n┌${'─'.repeat(title.length + 2)}┐`);
        console.log(`│ ${title} │`);
        console.log(`└${'─'.repeat(title.length + 2)}┘\n`);

        // 表格頭部
        const headers = [
            'Date'.padEnd(12),
            'Models'.padEnd(10),
            'Input'.padEnd(8),
            'Output'.padEnd(8),
            'Cache Create'.padEnd(12),
            'Cache Read'.padEnd(11),
            'Total Tokens'.padEnd(12),
            'Cost (USD)'.padEnd(10)
        ];

        // 打印表格邊框和標題
        const totalWidth = headers.reduce((sum, header) => sum + header.length, 0) + headers.length - 1;
        
        console.log('┌' + '─'.repeat(totalWidth) + '┐');
        console.log('│' + headers.join('│') + '│');
        console.log('├' + '─'.repeat(totalWidth) + '┤');

        // 打印數據行
        const sortedDates = Array.from(this.dailyStats.keys()).sort();
        
        for (const date of sortedDates) {
            const stat = this.dailyStats.get(date);
            const row = [
                date.padEnd(12),
                `- ${stat.model}`.padEnd(10),
                this.formatNumber(stat.input).padEnd(8),
                this.formatNumber(stat.output).padEnd(8),
                this.formatNumber(stat.cacheCreate).padEnd(12),
                this.formatNumber(stat.cacheRead).padEnd(11),
                this.formatNumber(stat.totalTokens).padEnd(12),
                `$${stat.cost.toFixed(2)}`.padEnd(10)
            ];
            console.log('│' + row.join('│') + '│');
        }

        // 總計行
        console.log('├' + '─'.repeat(totalWidth) + '┤');
        const totalRow = [
            'Total'.padEnd(12),
            ''.padEnd(10),
            this.formatNumber(this.totalStats.input).padEnd(8),
            this.formatNumber(this.totalStats.output).padEnd(8),
            this.formatNumber(this.totalStats.cacheCreate).padEnd(12),
            this.formatNumber(this.totalStats.cacheRead).padEnd(11),
            this.formatNumber(this.totalStats.totalTokens).padEnd(12),
            `$${this.totalStats.cost.toFixed(2)}`.padEnd(10)
        ];
        console.log('│' + totalRow.join('│') + '│');
        console.log('└' + '─'.repeat(totalWidth) + '┘\n');

        // 附加信息
        this.printAdditionalInfo();
    }

    printAdditionalInfo() {
        const recordCount = Array.from(this.dailyStats.values())
            .reduce((sum, stat) => sum + 1, 0);

        console.log(`📊 總結:`);
        console.log(`   • 記錄天數: ${recordCount} 天`);
        console.log(`   • 總 Token: ${this.formatNumber(this.totalStats.totalTokens)}`);
        console.log(`   • 輸入比例: ${((this.totalStats.input / this.totalStats.totalTokens) * 100).toFixed(1)}%`);
        console.log(`   • 輸出比例: ${((this.totalStats.output / this.totalStats.totalTokens) * 100).toFixed(1)}%`);
        console.log(`   • 平均日成本: $${(this.totalStats.cost / Math.max(recordCount, 1)).toFixed(3)}`);
        console.log(`   • Token 效率: ${(this.totalStats.totalTokens / Math.max(this.totalStats.cost, 0.001)).toFixed(0)} tokens/USD`);
    }

    formatNumber(num) {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num.toString();
    }

    getTotalStats() {
        return {
            success: true,
            totalRecords: Array.from(this.dailyStats.keys()).length,
            totalTokens: this.totalStats.totalTokens,
            totalCost: this.totalStats.cost,
            breakdown: {
                input: this.totalStats.input,
                output: this.totalStats.output,
                cacheCreate: this.totalStats.cacheCreate,
                cacheRead: this.totalStats.cacheRead
            }
        };
    }
}

// 導出類別以供其他模組使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { EnhancedTokenReporter };
}

// 如果直接執行此腳本
if (require.main === module) {
    async function main() {
        const reporter = new EnhancedTokenReporter();
        try {
            await reporter.generateReport();
        } catch (error) {
            console.error('執行失敗:', error.message);
            process.exit(1);
        }
    }
    
    main();
}
