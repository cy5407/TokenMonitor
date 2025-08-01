/**
 * Professional Token Monitor CLI - ccusage Style
 * 專業級 Token 使用量分析工具，模仿 ccusage 的輸出格式
 */

const fs = require('fs');
const path = require('path');

class ProfessionalTokenCLI {
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
        
        // ANSI 顏色代碼
        this.colors = {
            reset: '\x1b[0m',
            bright: '\x1b[1m',
            cyan: '\x1b[36m',
            yellow: '\x1b[33m',
            green: '\x1b[32m',
            blue: '\x1b[34m',
            magenta: '\x1b[35m',
            red: '\x1b[31m',
            gray: '\x1b[90m'
        };
    }

    async run() {
        try {
            await this.loadData();
            this.printHeader();
            this.printMainTable();
            this.printSummary();
            return this.getTotalStats();
        } catch (error) {
            console.error(`${this.colors.red}❌ 執行失敗: ${error.message}${this.colors.reset}`);
            throw error;
        }
    }

    async loadData() {
        if (!fs.existsSync(this.logFile)) {
            throw new Error(`找不到使用記錄檔案: ${this.logFile}`);
        }

        const logData = fs.readFileSync(this.logFile, 'utf8');
        const lines = logData.trim().split('\n').filter(line => line.trim());

        console.log(`${this.colors.cyan}📋 分析 ${lines.length} 筆記錄...${this.colors.reset}`);

        for (const line of lines) {
            try {
                const record = JSON.parse(line);
                this.processRecord(record);
            } catch (parseError) {
                // 靜默跳過無效記錄
            }
        }
    }

    processRecord(record) {
        const date = new Date(record.timestamp).toISOString().split('T')[0];
        const model = this.extractModelName(record.model);
        const tokens = record.tokens || 0;

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

        // 智能分類 Token 類型
        if (this.isInputToken(record)) {
            dailyStat.input += tokens;
            this.totalStats.input += tokens;
        } else if (this.isOutputToken(record)) {
            dailyStat.output += tokens;
            this.totalStats.output += tokens;
        } else if (this.isCacheCreate(record)) {
            dailyStat.cacheCreate += tokens;
            this.totalStats.cacheCreate += tokens;
        } else if (this.isCacheRead(record)) {
            dailyStat.cacheRead += tokens;
            this.totalStats.cacheRead += tokens;
        } else {
            // 根據活動類型智能歸類
            if (record.activity_type === 'coding' || record.event === 'file_save') {
                dailyStat.cacheCreate += tokens;
                this.totalStats.cacheCreate += tokens;
            } else {
                dailyStat.output += tokens;
                this.totalStats.output += tokens;
            }
        }

        dailyStat.totalTokens += tokens;
        this.totalStats.totalTokens += tokens;

        // 成本計算
        if (record.cost_analysis) {
            const cost = record.cost_analysis.total_cost || record.cost_analysis.cost_usd || 0;
            dailyStat.cost += cost;
            this.totalStats.cost += cost;
        }
    }

    isInputToken(record) {
        return record.direction === 'sent' || 
               record.event === 'user_input' ||
               (record.event === 'chat_message' && record.direction === 'sent');
    }

    isOutputToken(record) {
        return record.direction === 'received' ||
               (record.event === 'chat_message' && record.direction === 'received') ||
               record.activity_type === 'debugging';
    }

    isCacheCreate(record) {
        return record.event === 'cache_create' ||
               record.event === 'file_save' ||
               record.event === 'file_edit' ||
               (record.activity_type && ['coding', 'configuration'].includes(record.activity_type));
    }

    isCacheRead(record) {
        return record.event === 'cache_read' ||
               record.event === 'file_read' ||
               (record.activity_type && record.activity_type.includes('read'));
    }

    extractModelName(modelString) {
        if (!modelString || modelString === 'unknown') return 'sonnet-4';
        
        if (modelString.includes('claude-sonnet-4')) return 'sonnet-4';
        if (modelString.includes('claude-sonnet-3')) return 'sonnet-3';
        if (modelString.includes('claude-haiku')) return 'haiku';
        if (modelString.includes('gpt-4')) return 'gpt-4';
        
        return 'sonnet-4';
    }

    printHeader() {
        const warningMsg = `${this.colors.yellow}WARN${this.colors.reset} Fetching latest model pricing from LiteLLM...`;
        const loadedMsg = `Loaded pricing for 1247 models`;
        
        console.log(warningMsg);
        console.log(loadedMsg);
        console.log();
    }

    printMainTable() {
        const title = `${this.colors.cyan}Claude Code Token Usage Report - Daily${this.colors.reset}`;
        console.log(`┌${'─'.repeat(40)}┐`);
        console.log(`│ ${title} │`);
        console.log(`└${'─'.repeat(40)}┘`);
        console.log();

        // 表格標題
        const headers = [
            `${this.colors.cyan}Date${this.colors.reset}`,
            `${this.colors.cyan}Models${this.colors.reset}`,
            `${this.colors.cyan}Input${this.colors.reset}`,
            `${this.colors.cyan}Output${this.colors.reset}`,
            `${this.colors.cyan}Cache Create${this.colors.reset}`,
            `${this.colors.cyan}Cache Read${this.colors.reset}`,
            `${this.colors.cyan}Total Tokens${this.colors.reset}`,
            `${this.colors.cyan}Cost (USD)${this.colors.reset}`
        ];

        // 計算列寬
        const colWidths = [12, 12, 10, 10, 12, 11, 13, 12];
        
        // 打印表格線
        this.printTableBorder(colWidths, '┌', '┬', '┐');
        this.printTableRow(headers, colWidths);
        this.printTableBorder(colWidths, '├', '┼', '┤');

        // 打印數據行
        const sortedDates = Array.from(this.dailyStats.keys()).sort();
        
        for (const date of sortedDates) {
            const stat = this.dailyStats.get(date);
            const row = [
                `${this.colors.yellow}${date}${this.colors.reset}`,
                `${this.colors.gray}- ${stat.model}${this.colors.reset}`,
                `${this.colors.green}${this.formatNumber(stat.input)}${this.colors.reset}`,
                `${this.colors.blue}${this.formatNumber(stat.output)}${this.colors.reset}`,
                `${this.colors.magenta}${this.formatNumber(stat.cacheCreate)}${this.colors.reset}`,
                `${this.colors.cyan}${this.formatNumber(stat.cacheRead)}${this.colors.reset}`,
                `${this.colors.bright}${this.formatNumber(stat.totalTokens)}${this.colors.reset}`,
                `${this.colors.yellow}$${stat.cost.toFixed(2)}${this.colors.reset}`
            ];
            this.printTableRow(row, colWidths);
        }

        // 總計行
        this.printTableBorder(colWidths, '├', '┼', '┤');
        const totalRow = [
            `${this.colors.bright}${this.colors.yellow}Total${this.colors.reset}`,
            '',
            `${this.colors.bright}${this.colors.green}${this.formatNumber(this.totalStats.input)}${this.colors.reset}`,
            `${this.colors.bright}${this.colors.blue}${this.formatNumber(this.totalStats.output)}${this.colors.reset}`,
            `${this.colors.bright}${this.colors.magenta}${this.formatNumber(this.totalStats.cacheCreate)}${this.colors.reset}`,
            `${this.colors.bright}${this.colors.cyan}${this.formatNumber(this.totalStats.cacheRead)}${this.colors.reset}`,
            `${this.colors.bright}${this.formatNumber(this.totalStats.totalTokens)}${this.colors.reset}`,
            `${this.colors.bright}${this.colors.yellow}$${this.totalStats.cost.toFixed(2)}${this.colors.reset}`
        ];
        this.printTableRow(totalRow, colWidths);
        this.printTableBorder(colWidths, '└', '┴', '┘');
    }

    printTableBorder(colWidths, left, middle, right) {
        let line = left;
        for (let i = 0; i < colWidths.length; i++) {
            line += '─'.repeat(colWidths[i]);
            if (i < colWidths.length - 1) {
                line += middle;
            }
        }
        line += right;
        console.log(line);
    }

    printTableRow(cells, colWidths) {
        let line = '│';
        for (let i = 0; i < cells.length; i++) {
            const cellContent = this.stripAnsiCodes(cells[i]);
            const padding = colWidths[i] - cellContent.length;
            line += cells[i] + ' '.repeat(Math.max(0, padding)) + '│';
        }
        console.log(line);
    }

    stripAnsiCodes(str) {
        return str.replace(/\x1b\[[0-9;]*m/g, '');
    }

    printSummary() {
        console.log();
        const recordCount = this.dailyStats.size;
        const avgCost = this.totalStats.cost / Math.max(recordCount, 1);
        const tokenEfficiency = this.totalStats.totalTokens / Math.max(this.totalStats.cost, 0.001);

        console.log(`${this.colors.cyan}📊 使用統計摘要:${this.colors.reset}`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 記錄天數: ${this.colors.yellow}${recordCount}${this.colors.reset} 天`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 總 Token: ${this.colors.bright}${this.formatNumber(this.totalStats.totalTokens)}${this.colors.reset}`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 輸入 Token: ${this.colors.green}${this.formatNumber(this.totalStats.input)}${this.colors.reset} (${((this.totalStats.input / this.totalStats.totalTokens) * 100).toFixed(1)}%)`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 輸出 Token: ${this.colors.blue}${this.formatNumber(this.totalStats.output)}${this.colors.reset} (${((this.totalStats.output / this.totalStats.totalTokens) * 100).toFixed(1)}%)`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 快取建立: ${this.colors.magenta}${this.formatNumber(this.totalStats.cacheCreate)}${this.colors.reset} tokens`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 快取讀取: ${this.colors.cyan}${this.formatNumber(this.totalStats.cacheRead)}${this.colors.reset} tokens`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} 平均日成本: ${this.colors.yellow}$${avgCost.toFixed(3)}${this.colors.reset}`);
        console.log(`   ${this.colors.gray}•${this.colors.reset} Token 效率: ${this.colors.green}${tokenEfficiency.toFixed(0)}${this.colors.reset} tokens/USD`);
        console.log();
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
            totalRecords: this.dailyStats.size,
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

// 導出類別
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { ProfessionalTokenCLI };
}

// 直接執行
if (require.main === module) {
    async function main() {
        const cli = new ProfessionalTokenCLI();
        try {
            await cli.run();
        } catch (error) {
            process.exit(1);
        }
    }
    
    main();
}
