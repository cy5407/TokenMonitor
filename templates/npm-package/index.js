#!/usr/bin/env node

/**
 * Kiro Token Monitor - NPM Package Entry Point
 * 
 * 提供 Token 使用監控和分析功能
 */

const fs = require('fs');
const path = require('path');
const { EventEmitter } = require('events');

class KiroTokenMonitor extends EventEmitter {
    constructor(options = {}) {
        super();
        
        this.options = {
            logFile: options.logFile || path.join(process.cwd(), 'data', 'kiro-usage.log'),
            autoFlush: options.autoFlush !== false,
            maxLogSize: options.maxLogSize || 10 * 1024 * 1024, // 10MB
            ...options
        };
        
        this.buffer = [];
        this.ensureLogDirectory();
    }
    
    /**
     * 確保日誌目錄存在
     */
    ensureLogDirectory() {
        const logDir = path.dirname(this.options.logFile);
        if (!fs.existsSync(logDir)) {
            fs.mkdirSync(logDir, { recursive: true });
        }
    }
    
    /**
     * 記錄 Token 使用
     * @param {Object} data - 使用數據
     */
    log(data) {
        const record = {
            timestamp: new Date().toISOString(),
            ...data,
            session_id: data.session_id || `session-${Date.now()}`
        };
        
        if (this.options.autoFlush) {
            this.writeToFile(record);
        } else {
            this.buffer.push(record);
        }
        
        this.emit('logged', record);
        return this;
    }
    
    /**
     * 寫入檔案
     * @param {Object} record - 記錄
     */
    writeToFile(record) {
        try {
            fs.appendFileSync(this.options.logFile, JSON.stringify(record) + '\n');
        } catch (error) {
            this.emit('error', error);
        }
    }
    
    /**
     * 刷新緩衝區
     */
    flush() {
        if (this.buffer.length > 0) {
            const data = this.buffer.map(record => JSON.stringify(record)).join('\n') + '\n';
            fs.appendFileSync(this.options.logFile, data);
            this.buffer = [];
        }
        return this;
    }
    
    /**
     * 分析使用數據
     * @param {Object} options - 分析選項
     */
    analyze(options = {}) {
        const {
            since = null,
            until = null,
            groupBy = 'day'
        } = options;
        
        if (!fs.existsSync(this.options.logFile)) {
            return {
                totalRecords: 0,
                totalTokens: 0,
                totalCost: 0,
                records: []
            };
        }
        
        const content = fs.readFileSync(this.options.logFile, 'utf8');
        const lines = content.split('\n').filter(line => line.trim());
        
        let records = lines.map(line => {
            try {
                return JSON.parse(line);
            } catch {
                return null;
            }
        }).filter(Boolean);
        
        // 時間過濾
        if (since) {
            const sinceDate = new Date(since);
            records = records.filter(r => new Date(r.timestamp) >= sinceDate);
        }
        
        if (until) {
            const untilDate = new Date(until);
            records = records.filter(r => new Date(r.timestamp) <= untilDate);
        }
        
        // 統計計算
        const totalTokens = records.reduce((sum, r) => sum + (r.tokens || 0), 0);
        const totalCost = records.reduce((sum, r) => {
            const cost = r.cost_analysis?.cost_usd || r.cost || 0;
            return sum + parseFloat(cost);
        }, 0);
        
        // 按類型分組
        const byActivity = records.reduce((acc, r) => {
            const activity = r.activity_type || 'unknown';
            if (!acc[activity]) {
                acc[activity] = { count: 0, tokens: 0, cost: 0 };
            }
            acc[activity].count++;
            acc[activity].tokens += r.tokens || 0;
            acc[activity].cost += parseFloat(r.cost_analysis?.cost_usd || r.cost || 0);
            return acc;
        }, {});
        
        // 按模型分組
        const byModel = records.reduce((acc, r) => {
            const model = r.model || 'unknown';
            if (!acc[model]) {
                acc[model] = { count: 0, tokens: 0, cost: 0 };
            }
            acc[model].count++;
            acc[model].tokens += r.tokens || 0;
            acc[model].cost += parseFloat(r.cost_analysis?.cost_usd || r.cost || 0);
            return acc;
        }, {});
        
        return {
            totalRecords: records.length,
            totalTokens,
            totalCost: parseFloat(totalCost.toFixed(6)),
            byActivity,
            byModel,
            records: records.slice(-10) // 最近10筆記錄
        };
    }
    
    /**
     * 生成報告
     * @param {Object} options - 報告選項
     */
    generateReport(options = {}) {
        const analysis = this.analyze(options);
        
        console.log('\n📊 Kiro Token Monitor 報告');
        console.log('================================');
        console.log(`📈 總記錄數: ${analysis.totalRecords}`);
        console.log(`🔢 總 Token: ${analysis.totalTokens.toLocaleString()}`);
        console.log(`💰 總成本: $${analysis.totalCost}`);
        
        if (Object.keys(analysis.byActivity).length > 0) {
            console.log('\n📋 活動類型統計:');
            Object.entries(analysis.byActivity)
                .sort(([,a], [,b]) => b.tokens - a.tokens)
                .forEach(([activity, stats]) => {
                    console.log(`  ${activity}: ${stats.count} 次, ${stats.tokens} tokens ($${stats.cost.toFixed(6)})`);
                });
        }
        
        if (Object.keys(analysis.byModel).length > 0) {
            console.log('\n🤖 模型使用統計:');
            Object.entries(analysis.byModel)
                .sort(([,a], [,b]) => b.tokens - a.tokens)
                .forEach(([model, stats]) => {
                    console.log(`  ${model}: ${stats.count} 次, ${stats.tokens} tokens ($${stats.cost.toFixed(6)})`);
                });
        }
        
        return analysis;
    }
    
    /**
     * 清理舊記錄
     * @param {number} days - 保留天數
     */
    cleanup(days = 30) {
        if (!fs.existsSync(this.options.logFile)) return 0;
        
        const cutoff = new Date(Date.now() - days * 24 * 60 * 60 * 1000);
        const content = fs.readFileSync(this.options.logFile, 'utf8');
        const lines = content.split('\n').filter(line => line.trim());
        
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
            fs.writeFileSync(this.options.logFile, validLines.join('\n') + '\n');
            console.log(`🧹 清理了 ${removedCount} 筆舊記錄`);
        }
        
        return removedCount;
    }
}

// 便利方法
KiroTokenMonitor.create = (options) => new KiroTokenMonitor(options);

// CLI 支援
if (require.main === module) {
    const monitor = new KiroTokenMonitor();
    const command = process.argv[2];
    
    switch (command) {
        case 'report':
            monitor.generateReport();
            break;
            
        case 'cleanup':
            const days = parseInt(process.argv[3]) || 30;
            monitor.cleanup(days);
            break;
            
        case 'log':
            const [, , , event, tokens, cost] = process.argv;
            if (event && tokens) {
                monitor.log({
                    event,
                    tokens: parseInt(tokens),
                    cost: parseFloat(cost || 0),
                    activity_type: 'manual'
                });
                console.log(`✅ 記錄: ${event} - ${tokens} tokens`);
            } else {
                console.log('用法: node index.js log <event> <tokens> [cost]');
            }
            break;
            
        default:
            console.log(`
🚀 Kiro Token Monitor

用法:
  node index.js report              # 生成使用報告
  node index.js cleanup [days]      # 清理舊記錄 (預設30天)
  node index.js log <event> <tokens> [cost]  # 手動記錄

範例:
  node index.js report
  node index.js cleanup 7
  node index.js log chat_message 150 0.00045
            `);
    }
}

module.exports = KiroTokenMonitor;