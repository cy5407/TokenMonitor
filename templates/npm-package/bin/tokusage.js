#!/usr/bin/env node

/**
 * TokenUsage CLI Tool
 * 命令行介面工具，提供類似 ccusage 的體驗
 */

const { Command } = require('commander');
const KiroTokenMonitor = require('../index');
const fs = require('fs');
const path = require('path');

const program = new Command();

program
    .name('tokusage')
    .description('AI Token usage monitoring and analysis tool')
    .version('1.0.0');

// Daily report command
program
    .command('daily')
    .description('Show daily token usage report')
    .option('-s, --since <date>', 'Start date (YYYY-MM-DD)')
    .option('-m, --model <model>', 'Filter by model')
    .action((options) => {
        const monitor = new KiroTokenMonitor();
        
        console.log('🔍 執行每日 Token 使用分析...');
        
        const analysisOptions = {};
        if (options.since) {
            analysisOptions.since = options.since;
        }
        
        const analysis = monitor.analyze(analysisOptions);
        
        // 生成類似 ccusage 的表格報告
        console.log('\n┌────────────────────────────────────────┐');
        console.log('│ Token Usage Report - Daily             │');
        console.log('└────────────────────────────────────────┘');
        
        console.log('\n📊 使用統計摘要:');
        console.log(`   • 總記錄數: ${analysis.totalRecords}`);
        console.log(`   • 總 Token: ${analysis.totalTokens.toLocaleString()}`);
        console.log(`   • 總成本: $${analysis.totalCost}`);
        
        if (analysis.totalRecords > 0) {
            const avgTokens = Math.round(analysis.totalTokens / analysis.totalRecords);
            const avgCost = (analysis.totalCost / analysis.totalRecords).toFixed(6);
            console.log(`   • 平均每次: ${avgTokens} tokens ($${avgCost})`);
        }
        
        if (Object.keys(analysis.byModel).length > 0) {
            console.log('\n🤖 模型使用統計:');
            Object.entries(analysis.byModel)
                .sort(([,a], [,b]) => b.tokens - a.tokens)
                .forEach(([model, stats]) => {
                    console.log(`   • ${model}: ${stats.tokens} tokens ($${stats.cost.toFixed(6)})`);
                });
        }
        
        console.log('\n💡 提示:');
        console.log('   • 使用 \'tokusage summary\' 查看詳細統計');
        console.log('   • 使用 \'tokusage cleanup\' 清理舊記錄');
    });

// Summary command
program
    .command('summary')
    .description('Show detailed usage summary')
    .option('-d, --days <days>', 'Number of days to analyze', '7')
    .action((options) => {
        const monitor = new KiroTokenMonitor();
        const days = parseInt(options.days);
        
        console.log(`📈 執行 ${days} 天使用摘要分析...`);
        
        const since = new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString().split('T')[0];
        const analysis = monitor.analyze({ since });
        
        monitor.generateReport({ since });
        
        if (analysis.records.length > 0) {
            console.log('\n🔄 最近記錄:');
            analysis.records.slice(-5).forEach(record => {
                const time = new Date(record.timestamp).toLocaleString();
                const event = record.event || 'unknown';
                const tokens = record.tokens || 0;
                const activity = record.activity_type || 'general';
                console.log(`  [${time}] ${event} - ${tokens} tokens - ${activity}`);
            });
        }
    });

// Cleanup command
program
    .command('cleanup')
    .description('Clean up old usage records')
    .option('-d, --days <days>', 'Keep records for N days', '30')
    .action((options) => {
        const monitor = new KiroTokenMonitor();
        const days = parseInt(options.days);
        
        console.log(`🧹 清理 ${days} 天前的記錄...`);
        const removed = monitor.cleanup(days);
        
        if (removed > 0) {
            console.log(`✅ 已清理 ${removed} 筆舊記錄`);
        } else {
            console.log('✅ 沒有需要清理的記錄');
        }
    });

// Log command (for manual logging)
program
    .command('log <event> <tokens> [cost]')
    .description('Manually log token usage')
    .option('-a, --activity <type>', 'Activity type', 'manual')
    .option('-m, --model <model>', 'Model name', 'manual')
    .action((event, tokens, cost, options) => {
        const monitor = new KiroTokenMonitor();
        
        monitor.log({
            event,
            tokens: parseInt(tokens),
            cost: parseFloat(cost || 0),
            activity_type: options.activity,
            model: options.model
        });
        
        console.log(`✅ 記錄成功: ${event} - ${tokens} tokens`);
    });

// Status command
program
    .command('status')
    .description('Show monitoring status')
    .action(() => {
        const monitor = new KiroTokenMonitor();
        const logFile = monitor.options.logFile;
        
        console.log('📊 TokenMonitor 狀態');
        console.log('==================');
        
        if (fs.existsSync(logFile)) {
            const stats = fs.statSync(logFile);
            const lines = fs.readFileSync(logFile, 'utf8').split('\n').filter(Boolean);
            
            console.log(`✅ 監控狀態: 運行中`);
            console.log(`📁 記錄檔案: ${logFile}`);
            console.log(`📊 記錄數量: ${lines.length}`);
            console.log(`📏 檔案大小: ${(stats.size / 1024).toFixed(2)} KB`);
            console.log(`🕒 最後更新: ${stats.mtime.toLocaleString()}`);
            
            if (lines.length > 0) {
                try {
                    const lastRecord = JSON.parse(lines[lines.length - 1]);
                    console.log(`🔄 最後記錄: ${new Date(lastRecord.timestamp).toLocaleString()}`);
                } catch (e) {
                    console.log('⚠️  最後記錄格式錯誤');
                }
            }
        } else {
            console.log(`❌ 監控狀態: 未啟動`);
            console.log(`📁 記錄檔案: ${logFile} (不存在)`);
            console.log('💡 提示: 開始使用 Token 監控功能後會自動創建記錄檔案');
        }
    });

// Install command
program
    .command('install')
    .description('Install TokenMonitor in current project')
    .option('-f, --force', 'Force overwrite existing files')
    .action((options) => {
        console.log('🚀 安裝 TokenMonitor 到當前專案...');
        
        // 創建必要目錄
        const dirs = ['data', '.kiro/hooks'];
        dirs.forEach(dir => {
            if (!fs.existsSync(dir)) {
                fs.mkdirSync(dir, { recursive: true });
                console.log(`✅ 創建目錄: ${dir}`);
            }
        });
        
        // 創建基本配置
        const configPath = '.kiro/hooks/token-monitor.json';
        if (!fs.existsSync(configPath) || options.force) {
            const config = {
                name: "Token Monitor",
                description: "Monitor AI token usage",
                trigger: "manual",
                enabled: true
            };
            
            fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
            console.log(`✅ 創建配置: ${configPath}`);
        }
        
        console.log('🎉 TokenMonitor 安裝完成！');
        console.log('💡 使用 "tokusage status" 檢查狀態');
        console.log('💡 使用 "tokusage daily" 查看使用報告');
    });

// Parse command line arguments
program.parse();

// If no command provided, show help
if (!process.argv.slice(2).length) {
    program.outputHelp();
}