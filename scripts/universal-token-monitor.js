/**
 * Universal Token Monitor - 通用 Token 監控系統
 * 監控任何 IDE 或編輯器的文件變化，自動計算 Token 使用量
 */

const fs = require('fs');
const path = require('path');
const chokidar = require('chokidar');
const { execSync } = require('child_process');

class UniversalTokenMonitor {
    constructor(options = {}) {
        this.watchPaths = options.watchPaths || ['./**/*.{md,txt,js,ts,py,java,cpp,html,css,json}'];
        this.excludePaths = options.excludePaths || ['**/node_modules/**', '**/.git/**', '**/data/**'];
        this.logFile = options.logFile || 'data/kiro-usage.log';
        this.debounceTime = options.debounceTime || 2000; // 2秒防抖
        this.fileChangeQueue = new Map();
        this.isMonitoring = false;
        
        this.ensureDataDirectory();
        
        console.log('🚀 通用 Token 監控系統初始化...');
    }

    ensureDataDirectory() {
        const dataDir = path.dirname(this.logFile);
        if (!fs.existsSync(dataDir)) {
            fs.mkdirSync(dataDir, { recursive: true });
        }
    }

    start() {
        if (this.isMonitoring) {
            console.log('⚠️ 監控系統已在運行中');
            return;
        }

        console.log('👀 開始監控文件變化...');
        console.log(`📁 監控路徑: ${this.watchPaths.join(', ')}`);
        console.log(`🚫 排除路徑: ${this.excludePaths.join(', ')}`);

        this.watcher = chokidar.watch(this.watchPaths, {
            ignored: this.excludePaths,
            persistent: true,
            ignoreInitial: true,
            usePolling: false,
            interval: 1000
        });

        this.watcher
            .on('change', (filePath) => this.handleFileChange(filePath, 'modified'))
            .on('add', (filePath) => this.handleFileChange(filePath, 'created'))
            .on('ready', () => {
                this.isMonitoring = true;
                console.log('✅ 文件監控系統已啟動');
                console.log('📊 開始自動 Token 計算...\n');
            })
            .on('error', (error) => {
                console.error('❌ 監控系統錯誤:', error);
            });

        // 監控終端命令 (Windows PowerShell 歷史)
        this.startCommandMonitoring();
    }

    stop() {
        if (this.watcher) {
            this.watcher.close();
            this.isMonitoring = false;
            console.log('🛑 監控系統已停止');
        }
    }

    handleFileChange(filePath, operation) {
        // 防抖處理 - 避免短時間內重複觸發
        if (this.fileChangeQueue.has(filePath)) {
            clearTimeout(this.fileChangeQueue.get(filePath));
        }

        this.fileChangeQueue.set(filePath, setTimeout(() => {
            this.processFileChange(filePath, operation);
            this.fileChangeQueue.delete(filePath);
        }, this.debounceTime));
    }

    async processFileChange(filePath, operation) {
        try {
            // 檢查文件是否存在且可讀
            if (!fs.existsSync(filePath)) {
                return;
            }

            const stats = fs.statSync(filePath);
            if (!stats.isFile()) {
                return;
            }

            const content = fs.readFileSync(filePath, 'utf8');
            const tokens = this.calculateTokens(content);
            const fileExt = path.extname(filePath);
            const fileName = path.basename(filePath);

            if (tokens > 0) {
                const record = this.createTokenRecord({
                    event: operation === 'created' ? 'file_created' : 'file_modified',
                    filePath: filePath,
                    fileName: fileName,
                    fileType: fileExt,
                    content: content,
                    tokens: tokens,
                    operation: operation
                });

                this.logTokenUsage(record);
                this.displayFileChange(filePath, operation, tokens);
            }

        } catch (error) {
            console.error(`❌ 處理文件變化失敗 (${filePath}):`, error.message);
        }
    }

    calculateTokens(content) {
        // 簡化的 Token 計算 (類似 GPT Token 計算)
        // 1 token ≈ 4 個字符 (英文) 或 1.5 個中文字符
        const chineseChars = (content.match(/[\u4e00-\u9fff]/g) || []).length;
        const otherChars = content.length - chineseChars;
        
        const estimatedTokens = Math.ceil((chineseChars * 1.5 + otherChars) / 4);
        return estimatedTokens;
    }

    determineActivityType(filePath, content) {
        const ext = path.extname(filePath).toLowerCase();
        const fileName = path.basename(filePath).toLowerCase();

        // 根據文件類型和內容判斷活動類型
        if (['.md', '.txt', '.doc', '.docx'].includes(ext)) {
            if (content.includes('# ') || content.includes('## ')) {
                return 'documentation';
            }
            return 'writing';
        }

        if (['.js', '.ts', '.py', '.java', '.cpp', '.c', '.cs', '.php'].includes(ext)) {
            return 'coding';
        }

        if (['.json', '.yaml', '.yml', '.xml', '.ini', '.conf'].includes(ext)) {
            return 'configuration';
        }

        if (['.css', '.scss', '.less', '.html', '.htm'].includes(ext)) {
            return 'frontend';
        }

        if (fileName.includes('test') || fileName.includes('spec')) {
            return 'testing';
        }

        return 'general';
    }

    createTokenRecord({ event, filePath, fileName, fileType, content, tokens, operation }) {
        const now = new Date();
        const activityType = this.determineActivityType(filePath, content);
        
        return {
            timestamp: now.toISOString(),
            event: event,
            operation: operation,
            file_path: filePath,
            file_name: fileName,
            file_type: fileType,
            content_length: content.length,
            tokens: tokens,
            activity_type: activityType,
            model: 'universal-monitor',
            session_id: `file-monitor-${Date.now()}`,
            cost_analysis: {
                tokens: tokens,
                cost_usd: tokens * 0.000003, // 假設平均成本
                cost_type: 'input',
                model: 'universal-monitor',
                pricing_rate: 3
            },
            content_preview: content.substring(0, 100) + (content.length > 100 ? '...' : ''),
            ide_detected: this.detectIDE(),
            machine_info: {
                platform: process.platform,
                user: process.env.USERNAME || process.env.USER,
                working_directory: process.cwd()
            }
        };
    }

    detectIDE() {
        try {
            // 檢測正在運行的 IDE
            const processes = execSync('tasklist /FO CSV', { encoding: 'utf8' });
            
            if (processes.includes('Code.exe')) return 'VS Code';
            if (processes.includes('devenv.exe')) return 'Visual Studio';
            if (processes.includes('idea64.exe') || processes.includes('idea.exe')) return 'IntelliJ IDEA';
            if (processes.includes('sublime_text.exe')) return 'Sublime Text';
            if (processes.includes('notepad++.exe')) return 'Notepad++';
            if (processes.includes('atom.exe')) return 'Atom';
            if (processes.includes('vim.exe') || processes.includes('nvim.exe')) return 'Vim/Neovim';
            
            return 'Unknown';
        } catch (error) {
            return 'Unknown';
        }
    }

    logTokenUsage(record) {
        try {
            const jsonRecord = JSON.stringify(record);
            fs.appendFileSync(this.logFile, jsonRecord + '\n', 'utf8');
        } catch (error) {
            console.error('❌ 記錄 Token 使用失敗:', error.message);
        }
    }

    displayFileChange(filePath, operation, tokens) {
        const timestamp = new Date().toLocaleTimeString('zh-TW');
        const relativeePath = path.relative(process.cwd(), filePath);
        const emoji = operation === 'created' ? '📄' : '✏️';
        const action = operation === 'created' ? '創建' : '修改';
        
        console.log(`${emoji} [${timestamp}] ${action}文件: ${relativeePath}`);
        console.log(`   📊 Token: ${tokens} | 活動: ${this.determineActivityType(filePath, '')}`);
        console.log(`   🔧 IDE: ${this.detectIDE()}\n`);
    }

    startCommandMonitoring() {
        // 監控 PowerShell 命令歷史 (Windows)
        if (process.platform === 'win32') {
            try {
                const historyPath = path.join(
                    process.env.APPDATA, 
                    'Microsoft', 'Windows', 'PowerShell', 'PSReadLine', 
                    'ConsoleHost_history.txt'
                );

                if (fs.existsSync(historyPath)) {
                    console.log('📋 監控 PowerShell 命令歷史...');
                    
                    fs.watchFile(historyPath, { interval: 5000 }, () => {
                        this.processCommandHistory(historyPath);
                    });
                }
            } catch (error) {
                console.warn('⚠️ 無法監控命令歷史:', error.message);
            }
        }
    }

    processCommandHistory(historyPath) {
        try {
            const content = fs.readFileSync(historyPath, 'utf8');
            const lines = content.split('\n');
            const lastCommand = lines[lines.length - 2]?.trim(); // 倒數第二行是最新命令

            if (lastCommand && this.isRelevantCommand(lastCommand)) {
                const tokens = Math.ceil(lastCommand.length / 4);
                
                const record = this.createTokenRecord({
                    event: 'command_executed',
                    filePath: 'terminal',
                    fileName: 'command',
                    fileType: '.cmd',
                    content: lastCommand,
                    tokens: tokens,
                    operation: 'executed'
                });

                this.logTokenUsage(record);
                console.log(`⚡ [${new Date().toLocaleTimeString('zh-TW')}] 命令執行: ${lastCommand.substring(0, 50)}...`);
                console.log(`   📊 Token: ${tokens} | 類型: terminal\n`);
            }
        } catch (error) {
            // 忽略錯誤，避免干擾主要功能
        }
    }

    isRelevantCommand(command) {
        const relevantCommands = [
            'node', 'npm', 'git', 'code', 'python', 'pip',
            'tokusage', 'New-Item', 'Set-Content', 'Add-Content'
        ];
        
        return relevantCommands.some(cmd => command.toLowerCase().includes(cmd.toLowerCase()));
    }

    generateReport() {
        console.log('\n📊 ===== 通用 Token 監控報告 =====');
        console.log(`📁 監控中的路徑: ${this.isMonitoring ? '✅' : '❌'}`);
        console.log(`📄 日誌文件: ${this.logFile}`);
        console.log(`🕒 運行時間: ${this.isMonitoring ? '運行中' : '已停止'}`);
        
        if (fs.existsSync(this.logFile)) {
            const content = fs.readFileSync(this.logFile, 'utf8');
            const lines = content.trim().split('\n').filter(line => line.trim());
            console.log(`📝 總記錄數: ${lines.length}`);
            
            if (lines.length > 0) {
                try {
                    const lastRecord = JSON.parse(lines[lines.length - 1]);
                    console.log(`⏰ 最後記錄: ${new Date(lastRecord.timestamp).toLocaleString('zh-TW')}`);
                } catch (error) {
                    console.log('⏰ 最後記錄: 解析失敗');
                }
            }
        } else {
            console.log('📝 尚無記錄');
        }
        console.log('====================================\n');
    }
}

// 導出類別
module.exports = { UniversalTokenMonitor };

// 如果直接執行此腳本
if (require.main === module) {
    const monitor = new UniversalTokenMonitor({
        watchPaths: [
            './**/*.{md,txt,js,ts,py,java,cpp,c,cs,php,html,css,json,yaml,yml}',
            '!**/node_modules/**',
            '!**/.git/**',
            '!**/data/kiro-usage.log'
        ]
    });

    // 處理程序退出
    process.on('SIGINT', () => {
        console.log('\n🛑 接收到退出信號...');
        monitor.stop();
        monitor.generateReport();
        process.exit(0);
    });

    process.on('SIGTERM', () => {
        monitor.stop();
        process.exit(0);
    });

    // 開始監控
    monitor.start();

    // 顯示狀態報告
    setInterval(() => {
        monitor.generateReport();
    }, 300000); // 每5分鐘顯示一次報告
}
