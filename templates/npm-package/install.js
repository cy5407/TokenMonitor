#!/usr/bin/env node

/**
 * TokenMonitor 安裝腳本
 * 在 npm install 後自動執行，設置必要的目錄和配置
 */

const fs = require('fs');
const path = require('path');

console.log('🚀 正在安裝 Kiro Token Monitor...');

try {
    // 創建必要目錄
    const requiredDirs = [
        'data',
        '.kiro',
        '.kiro/hooks',
        '.kiro/settings'
    ];
    
    requiredDirs.forEach(dir => {
        if (!fs.existsSync(dir)) {
            fs.mkdirSync(dir, { recursive: true });
            console.log(`✅ 創建目錄: ${dir}`);
        }
    });
    
    // 創建預設配置檔案
    const defaultConfig = {
        name: "Kiro Token Monitor",
        description: "Monitor and analyze AI token usage across different IDEs",
        version: "1.0.0",
        settings: {
            logFile: "data/kiro-usage.log",
            maxLogSize: "10MB",
            retentionDays: 30,
            autoCleanup: true
        },
        monitoring: {
            enabled: true,
            events: [
                "chat_message",
                "file_save",
                "file_edit",
                "tool_execution",
                "command_executed"
            ]
        },
        reporting: {
            defaultPeriod: "daily",
            includeDetails: true,
            showCosts: true
        }
    };
    
    const configPath = '.kiro/settings/token-monitor.json';
    if (!fs.existsSync(configPath)) {
        fs.writeFileSync(configPath, JSON.stringify(defaultConfig, null, 2));
        console.log(`✅ 創建配置檔案: ${configPath}`);
    }
    
    // 創建 Hook 配置
    const hookConfig = {
        name: "Token Usage Monitor",
        description: "Automatically monitor token usage",
        trigger: "file.saved",
        enabled: true,
        script: "node_modules/kiro-token-monitor/index.js",
        options: {
            autoLog: true,
            includeFileContent: false
        }
    };
    
    const hookPath = '.kiro/hooks/token-monitor.json';
    if (!fs.existsSync(hookPath)) {
        fs.writeFileSync(hookPath, JSON.stringify(hookConfig, null, 2));
        console.log(`✅ 創建 Hook 配置: ${hookPath}`);
    }
    
    // 創建快速開始指南
    const quickStartGuide = `# Kiro Token Monitor - 快速開始

## 安裝完成！

TokenMonitor 已成功安裝到您的專案中。

## 基本使用

### 1. 查看使用狀態
\`\`\`bash
npx tokusage status
\`\`\`

### 2. 查看每日報告
\`\`\`bash
npx tokusage daily
\`\`\`

### 3. 查看詳細統計
\`\`\`bash
npx tokusage summary
\`\`\`

### 4. 手動記錄使用
\`\`\`bash
npx tokusage log chat_message 150 0.00045
\`\`\`

### 5. 清理舊記錄
\`\`\`bash
npx tokusage cleanup --days 30
\`\`\`

## 程式化使用

\`\`\`javascript
const KiroTokenMonitor = require('kiro-token-monitor');

// 創建監控實例
const monitor = new KiroTokenMonitor({
    logFile: './data/my-usage.log'
});

// 記錄使用
monitor.log({
    event: 'chat_message',
    tokens: 150,
    cost: 0.00045,
    activity_type: 'coding',
    model: 'claude-sonnet-4.0'
});

// 生成報告
monitor.generateReport();

// 分析數據
const analysis = monitor.analyze({
    since: '2025-01-01'
});
console.log(analysis);
\`\`\`

## 配置檔案

- \`.kiro/settings/token-monitor.json\` - 主要配置
- \`.kiro/hooks/token-monitor.json\` - Hook 配置

## 支援

如有問題，請查看：
- GitHub Issues
- 文件: https://github.com/your-org/kiro-token-monitor

祝您使用愉快！ 🎉
`;
    
    const guidePath = 'TOKENMONITOR-QUICKSTART.md';
    if (!fs.existsSync(guidePath)) {
        fs.writeFileSync(guidePath, quickStartGuide);
        console.log(`✅ 創建快速指南: ${guidePath}`);
    }
    
    console.log('');
    console.log('🎉 Kiro Token Monitor 安裝完成！');
    console.log('');
    console.log('📋 後續步驟:');
    console.log('1. 執行 "npx tokusage status" 檢查狀態');
    console.log('2. 執行 "npx tokusage daily" 查看使用報告');
    console.log('3. 閱讀 "TOKENMONITOR-QUICKSTART.md" 了解更多用法');
    console.log('');
    console.log('💡 提示: TokenMonitor 會自動監控您的 AI 使用情況');
    console.log('');
    
} catch (error) {
    console.error('❌ 安裝失敗:', error.message);
    process.exit(1);
}