/**
 * 改進的 AI Token 監控系統
 * 專注於捕捉真正的 AI 對話和內容生成
 */

const fs = require('fs');
const path = require('path');
const chokidar = require('chokidar');
const { execSync } = require('child_process');

class ImprovedAIMonitor {
    constructor() {
        this.logFile = 'data/kiro-usage.log';
        this.isMonitoring = false;
        this.lastLogSize = 0;
        this.sessionId = `improved-monitor-${Date.now()}`;
        
        // AI 內容檢測模式
        this.aiPatterns = [
            // 程式碼生成模式
            /function\s+\w+\s*\([^)]*\)\s*{[\s\S]*}/,
            /class\s+\w+\s*{[\s\S]*}/,
            /const\s+\w+\s*=\s*\([^)]*\)\s*=>/,
            
            // Markdown 程式碼區塊
            /```[\s\S]*?```/,
            
            // 註解模式（AI 常生成的）
            /\/\*\*[\s\S]*?\*\//,
            /\/\/\s*(Generated|Created|AI|Assistant)/i,
            
            // 文件結構模式
            /#{1,6}\s+.+/,  // Markdown 標題
            /^\s*-\s+.+/m,  // 列表項目
            
            // 對話模式
            /^(User|Assistant|AI):/m,
            /^Q:|^A:/m
        ];
        
        console.log('🚀 改進的 AI 監控系統初始化...');
    }
    
    start() {
        if (this.isMonitoring) {
            console.log('⚠️ 監控系統已在運行中');
            return;
        }
        
        console.log('👀 開始改進的 AI 監控...');
        
        // 1. 監控檔案變化（重點關注可能包含 AI 內容的檔案）
        this.startFileMonitoring();
        
        // 2. 監控 Kiro IDE 日誌檔案變化
        this.startKiroLogMonitoring();
        
        // 3. 監控 VS Code 相關檔案
        this.startVSCodeMonitoring();
        
        this.isMonitoring = true;
        console.log('✅ 改進的監控系統已啟動\n');
    }
    
    startFileMonitoring() {
        const watchPaths = [
            './**/*.{md,txt,js,ts,py,java,cpp,html,css,json}',
            '!**/node_modules/**',
            '!**/.git/**',
            '!**/data/kiro-usage.log'
        ];
        
        this.fileWatcher = chokidar.watch(watchPaths, {
            ignored: ['**/node_modules/**', '**/.git/**'],
            persistent: true,
            ignoreInitial: true
        });
        
        this.fileWatcher.on('change', (filePath) => {
            this.handleFileChange(filePath, 'modified');
        });
        
        this.fileWatcher.on('add', (filePath) => {
            this.handleFileChange(filePath, 'created');
        });
        
        console.log('📁 檔案監控已啟動');
    }
    
    startKiroLogMonitoring() {
        // 監控 Kiro IDE 的日誌檔案變化
        if (fs.existsSync(this.logFile)) {
            this.lastLogSize = fs.statSync(this.logFile).size;
            
            fs.watchFile(this.logFile, { interval: 1000 }, () => {
                this.handleKiroLogChange();
            });
            
            console.log('📋 Kiro 日誌監控已啟動');
        }
    }
    
    startVSCodeMonitoring() {
        // 監控 VS Code 相關的檔案和目錄
        const vscodePatterns = [
            '.vscode/**/*',
            '**/.vscode/**/*'
        ];
        
        try {
            this.vscodeWatcher = chokidar.watch(vscodePatterns, {
                persistent: true,
                ignoreInitial: true
            });
            
            this.vscodeWatcher.on('change', (filePath) => {
                this.handleVSCodeActivity(filePath);
            });
            
            console.log('🔧 VS Code 監控已啟動');
        } catch (error) {
            console.log('⚠️ VS Code 監控啟動失敗:', error.message);
        }
    }
    
    async handleFileChange(filePath, operation) {
        try {
            if (!fs.existsSync(filePath)) return;
            
            const content = fs.readFileSync(filePath, 'utf8');
            
            // 檢查是否為 AI 生成的內容
            if (this.isAIGeneratedContent(content)) {
                const tokens = this.calculateTokens(content);
                const activityType = this.classifyActivity(content, filePath);
                const ide = this.detectCurrentIDE();
                
                const record = {
                    timestamp: new Date().toISOString(),
                    event: 'ai_content_detected',
                    source: 'file_monitoring',
                    ide: ide,
                    file_path: filePath,
                    file_type: path.extname(filePath),
                    operation: operation,
                    content_length: content.length,
                    tokens: tokens,
                    activity_type: activityType,
                    session_id: this.sessionId,
                    cost_analysis: this.calculateCost(tokens),
                    ai_confidence: this.calculateAIConfidence(content),
                    content_preview: content.substring(0, 200) + (content.length > 200 ? '...' : '')
                };
                
                this.logAIActivity(record);
                this.displayAIActivity(record);
            }
            
        } catch (error) {
            console.error(`❌ 處理檔案變化失敗 (${filePath}):`, error.message);
        }
    }
    
    handleKiroLogChange() {
        try {
            const currentSize = fs.statSync(this.logFile).size;
            
            if (currentSize > this.lastLogSize) {
                // 讀取新增的內容
                const stream = fs.createReadStream(this.logFile, {
                    start: this.lastLogSize,
                    end: currentSize
                });
                
                let newContent = '';
                stream.on('data', (chunk) => {
                    newContent += chunk.toString();
                });
                
                stream.on('end', () => {
                    this.analyzeNewKiroActivity(newContent);
                });
                
                this.lastLogSize = currentSize;
            }
            
        } catch (error) {
            console.error('❌ 處理 Kiro 日誌變化失敗:', error.message);
        }
    }
    
    analyzeNewKiroActivity(newContent) {
        const lines = newContent.trim().split('\n').filter(line => line.trim());
        
        for (const line of lines) {
            try {
                const record = JSON.parse(line);
                
                // 檢查是否為真實的 AI 對話
                if (record.event === 'chat_message' && 
                    record.session_id && 
                    !record.session_id.includes('test')) {
                    
                    console.log(`🤖 檢測到真實 AI 對話: ${record.direction} - ${record.tokens} tokens`);
                    
                    // 可以在這裡加入額外的分析或處理
                }
                
            } catch (error) {
                // 忽略解析錯誤
            }
        }
    }
    
    handleVSCodeActivity(filePath) {
        console.log(`🔧 VS Code 活動: ${filePath}`);
        
        // 檢查是否為 AI 相關的配置變化
        if (filePath.includes('settings.json') || filePath.includes('extensions.json')) {
            try {
                const content = fs.readFileSync(filePath, 'utf8');
                
                // 檢查是否包含 AI 擴展配置
                if (content.includes('copilot') || content.includes('chatgpt') || content.includes('ai')) {
                    console.log('🤖 檢測到 VS Code AI 擴展配置變化');
                    
                    const record = {
                        timestamp: new Date().toISOString(),
                        event: 'vscode_ai_config_change',
                        source: 'vscode_monitoring',
                        ide: 'vscode',
                        file_path: filePath,
                        activity_type: 'configuration',
                        session_id: this.sessionId
                    };
                    
                    this.logAIActivity(record);
                }
                
            } catch (error) {
                console.error('❌ 處理 VS Code 活動失敗:', error.message);
            }
        }
    }
    
    isAIGeneratedContent(content) {
        // 檢查內容是否符合 AI 生成的模式
        let score = 0;
        
        for (const pattern of this.aiPatterns) {
            if (pattern.test(content)) {
                score++;
            }
        }
        
        // 額外檢查
        if (content.length > 100 && content.includes('function')) score++;
        if (content.includes('```') && content.includes('javascript')) score++;
        if (content.includes('# ') && content.length > 50) score++;
        
        return score >= 2; // 至少符合 2 個模式才認為是 AI 生成
    }
    
    calculateAIConfidence(content) {
        let confidence = 0;
        
        // 基於多個因素計算信心度
        if (content.includes('function') && content.includes('{')) confidence += 0.3;
        if (content.includes('```')) confidence += 0.2;
        if (content.match(/\/\*\*[\s\S]*?\*\//)) confidence += 0.2;
        if (content.length > 200) confidence += 0.1;
        if (content.includes('return') && content.includes(';')) confidence += 0.2;
        
        return Math.min(confidence, 1.0);
    }
    
    calculateTokens(content) {
        // 改進的 Token 計算
        const chineseChars = (content.match(/[\u4e00-\u9fff]/g) || []).length;
        const englishChars = content.length - chineseChars;
        
        // 程式碼通常有更多的 token
        let multiplier = 1;
        if (content.includes('function') || content.includes('class')) {
            multiplier = 1.2;
        }
        
        return Math.ceil((chineseChars / 1.5 + englishChars / 4) * multiplier);
    }
    
    calculateCost(tokens) {
        const pricing = {
            input: 3.0 / 1000000,   // $3 per 1M tokens
            output: 15.0 / 1000000  // $15 per 1M tokens
        };
        
        // 假設為輸出 token（AI 生成的內容）
        const cost = tokens * pricing.output;
        
        return {
            tokens: tokens,
            cost_usd: parseFloat(cost.toFixed(6)),
            cost_type: 'output',
            model: 'claude-sonnet-4.0',
            pricing_rate: pricing.output * 1000000
        };
    }
    
    classifyActivity(content, filePath) {
        const ext = path.extname(filePath).toLowerCase();
        
        if (['.js', '.ts', '.py', '.java', '.cpp'].includes(ext)) {
            return 'ai_coding';
        }
        
        if (['.md', '.txt'].includes(ext)) {
            if (content.includes('#') || content.includes('##')) {
                return 'ai_documentation';
            }
            return 'ai_writing';
        }
        
        return 'ai_general';
    }
    
    detectCurrentIDE() {
        try {
            const processes = execSync('tasklist /FO CSV', { encoding: 'utf8' });
            
            if (processes.includes('Code.exe')) return 'vscode';
            if (processes.includes('Kiro.exe')) return 'kiro';
            
            return 'unknown';
        } catch (error) {
            return 'unknown';
        }
    }
    
    logAIActivity(record) {
        try {
            const logEntry = JSON.stringify(record) + '\n';
            fs.appendFileSync(this.logFile, logEntry, 'utf8');
        } catch (error) {
            console.error('❌ 記錄 AI 活動失敗:', error.message);
        }
    }
    
    displayAIActivity(record) {
        const timestamp = new Date(record.timestamp).toLocaleTimeString('zh-TW');
        const fileName = path.basename(record.file_path);
        
        console.log(`🤖 [${timestamp}] AI 內容檢測: ${fileName}`);
        console.log(`   📊 Token: ${record.tokens} | 信心度: ${(record.ai_confidence * 100).toFixed(1)}% | IDE: ${record.ide}`);
        console.log(`   💰 成本: ${record.cost_analysis.cost_usd.toFixed(6)} USD | 類型: ${record.activity_type}`);
        console.log('');
    }
    
    stop() {
        if (this.fileWatcher) {
            this.fileWatcher.close();
        }
        
        if (this.vscodeWatcher) {
            this.vscodeWatcher.close();
        }
        
        if (fs.existsSync(this.logFile)) {
            fs.unwatchFile(this.logFile);
        }
        
        this.isMonitoring = false;
        console.log('🛑 改進的監控系統已停止');
    }
    
    getStatus() {
        return {
            isMonitoring: this.isMonitoring,
            sessionId: this.sessionId,
            logFile: this.logFile,
            lastLogSize: this.lastLogSize
        };
    }
}

// 如果直接執行此腳本
if (require.main === module) {
    const monitor = new ImprovedAIMonitor();
    
    // 處理程序退出
    process.on('SIGINT', () => {
        console.log('\n🛑 接收到退出信號...');
        monitor.stop();
        process.exit(0);
    });
    
    // 開始監控
    monitor.start();
    
    // 顯示狀態
    setInterval(() => {
        const status = monitor.getStatus();
        console.log(`📊 監控狀態: ${status.isMonitoring ? '運行中' : '已停止'} | Session: ${status.sessionId}`);
    }, 60000); // 每分鐘顯示一次狀態
}

module.exports = { ImprovedAIMonitor };