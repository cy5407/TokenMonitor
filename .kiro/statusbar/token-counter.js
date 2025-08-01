/**
 * Kiro IDE 狀態列 Token 計數器
 * 即時顯示當前文件或選中文本的 Token 數量
 */

const { execSync } = require('child_process');
const path = require('path');

class TokenCounterStatusBar {
    constructor() {
        this.tokenMonitorPath = path.join(__dirname, '..', 'token-monitor.exe');
        this.isEnabled = true;
        this.updateInterval = 1000; // 1秒更新一次
        this.lastText = '';
        this.lastTokenCount = 0;
    }

    /**
     * 初始化狀態列項目
     */
    initialize() {
        return {
            id: 'token-counter',
            name: 'Token Counter',
            text: '🎯 0 tokens',
            tooltip: 'Token 數量 (點擊查看詳細資訊)',
            priority: 100,
            alignment: 'right',
            command: 'token-monitor.show-details',
            backgroundColor: '#1e1e1e',
            color: '#ffffff'
        };
    }

    /**
     * 更新 Token 計數
     */
    async updateTokenCount(context) {
        if (!this.isEnabled) return null;

        try {
            const { activeEditor, selection } = context;
            
            if (!activeEditor) {
                return this.createStatusItem('🎯 - tokens', 'No active editor');
            }

            // 獲取要計算的文本
            const text = selection && selection.text ? 
                selection.text : 
                activeEditor.document.getText();

            // 如果文本沒有變化，返回快取結果
            if (text === this.lastText) {
                return this.createStatusItem(
                    `🎯 ${this.lastTokenCount} tokens`,
                    this.createTooltip(text, this.lastTokenCount)
                );
            }

            // 計算新的 Token 數量
            const tokenCount = await this.calculateTokens(text);
            this.lastText = text;
            this.lastTokenCount = tokenCount;

            return this.createStatusItem(
                `🎯 ${tokenCount} tokens`,
                this.createTooltip(text, tokenCount)
            );

        } catch (error) {
            console.error('Token counter update error:', error);
            return this.createStatusItem('🎯 Error', 'Token calculation failed');
        }
    }

    /**
     * 計算 Token 數量
     */
    async calculateTokens(text) {
        if (!text || text.trim().length === 0) return 0;

        try {
            const command = `"${this.tokenMonitorPath}" calculate "${text.replace(/"/g, '\\"')}" --quiet`;
            const output = execSync(command, { 
                encoding: 'utf8', 
                timeout: 5000,
                stdio: ['pipe', 'pipe', 'ignore'] // 忽略 stderr
            });
            
            const match = output.match(/Token 數量:\s*(\d+)/);
            return match ? parseInt(match[1]) : this.estimateTokens(text);
            
        } catch (error) {
            return this.estimateTokens(text);
        }
    }

    /**
     * 簡單估算（備用方法）
     */
    estimateTokens(text) {
        const englishChars = (text.match(/[a-zA-Z0-9\s]/g) || []).length;
        const chineseChars = text.length - englishChars;
        return Math.ceil(englishChars / 4 + chineseChars / 1.5);
    }

    /**
     * 建立狀態列項目
     */
    createStatusItem(text, tooltip) {
        return {
            text: text,
            tooltip: tooltip,
            color: this.getColorByTokenCount(this.lastTokenCount),
            command: 'token-monitor.show-details'
        };
    }

    /**
     * 建立工具提示
     */
    createTooltip(text, tokenCount) {
        const charCount = text.length;
        const wordCount = text.split(/\s+/).filter(word => word.length > 0).length;
        const estimatedCost = this.calculateCost(tokenCount);

        return `Token 詳細資訊:
📊 Token 數量: ${tokenCount}
📝 字符數量: ${charCount}
📄 單詞數量: ${wordCount}
💰 估算成本: $${estimatedCost.toFixed(6)}
🔄 點擊查看更多詳細資訊`;
    }

    /**
     * 根據 Token 數量獲取顏色
     */
    getColorByTokenCount(count) {
        if (count === 0) return '#888888';
        if (count < 100) return '#00ff00';
        if (count < 500) return '#ffff00';
        if (count < 1000) return '#ff8800';
        return '#ff0000';
    }

    /**
     * 計算估算成本（基於 Claude Sonnet 4.0）
     */
    calculateCost(tokenCount) {
        const inputCostPer1M = 3.0; // $3 per 1M input tokens
        return (tokenCount / 1000000) * inputCostPer1M;
    }

    /**
     * 顯示詳細資訊
     */
    async showDetails(context) {
        const { activeEditor, selection } = context;
        
        if (!activeEditor) {
            return { message: 'No active editor' };
        }

        const text = selection && selection.text ? 
            selection.text : 
            activeEditor.document.getText();

        try {
            const command = `"${this.tokenMonitorPath}" calculate "${text.replace(/"/g, '\\"')}" --details`;
            const output = execSync(command, { encoding: 'utf8', timeout: 10000 });
            
            return {
                title: 'Token 計算詳細資訊',
                content: output,
                type: 'info'
            };
            
        } catch (error) {
            return {
                title: 'Token 計算錯誤',
                content: error.message,
                type: 'error'
            };
        }
    }

    /**
     * 切換啟用狀態
     */
    toggle() {
        this.isEnabled = !this.isEnabled;
        return {
            message: `Token Counter ${this.isEnabled ? 'enabled' : 'disabled'}`,
            statusItem: this.isEnabled ? 
                this.createStatusItem('🎯 0 tokens', 'Token Counter enabled') :
                null
        };
    }
}

// 匯出給 Kiro IDE 使用
module.exports = {
    TokenCounterStatusBar,
    
    // Kiro IDE 狀態列 API
    initialize() {
        const counter = new TokenCounterStatusBar();
        return counter.initialize();
    },
    
    async update(context) {
        const counter = new TokenCounterStatusBar();
        return await counter.updateTokenCount(context);
    },
    
    async onCommand(command, context) {
        const counter = new TokenCounterStatusBar();
        
        switch (command) {
            case 'token-monitor.show-details':
                return await counter.showDetails(context);
            case 'token-monitor.toggle':
                return counter.toggle();
            default:
                return { message: 'Unknown command' };
        }
    }
};