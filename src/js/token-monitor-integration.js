/**
 * Token Monitor Kiro IDE 整合腳本
 * 這個腳本會被 Kiro IDE 的 Agent Hook 系統調用
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

// Token Monitor 執行檔路徑
const TOKEN_MONITOR_PATH = path.join(__dirname, 'token-monitor.exe');
const DATA_DIR = path.join(__dirname, 'data');
const LOG_FILE = path.join(DATA_DIR, 'kiro-usage.log');

// 確保資料目錄存在
if (!fs.existsSync(DATA_DIR)) {
    fs.mkdirSync(DATA_DIR, { recursive: true });
}

/**
 * 主要 Hook 執行函數
 * @param {Object} context - Kiro IDE 提供的上下文
 */
async function execute(context) {
    try {
        console.log('🎯 Token Monitor Hook 被觸發!');
        console.log('📋 接收到的上下文:', JSON.stringify(context, null, 2));
        
        // 當 kiro-usage.log 檔案被更新時，分析整個檔案
        if (context && context.filePath && context.filePath.includes('kiro-usage.log')) {
            console.log('📊 檢測到 kiro-usage.log 檔案更新，開始分析...');
            return await analyzeUsageLog();
        }
        
        // 處理其他類型的事件
        let event, data, timestamp;
        
        if (context && typeof context === 'object') {
            event = context.event || context.type || context.eventType || 'file.saved';
            data = context.data || context.payload || context.content || context;
            timestamp = context.timestamp || context.time;
        } else {
            console.log('⚠️ 無效的上下文格式:', JSON.stringify(context));
            return { success: false, message: 'Invalid context format' };
        }
        
        console.log(`🎯 Token Monitor: Processing ${event}`);
        
        let result = null;
        
        switch (event) {
            case 'file.saved':
            case 'fileSaved':
            case 'fileEdited':
                // 如果是 kiro-usage.log 檔案被儲存，分析整個檔案
                if (data && data.filePath && data.filePath.includes('kiro-usage.log')) {
                    result = await analyzeUsageLog();
                } else {
                    result = await handleFileSave(data, timestamp);
                }
                break;
            case 'chat.message.sent':
            case 'chatMessageSent':
            case 'message.sent':
            case 'user.message':
                result = await handleChatMessage(data, 'sent', timestamp);
                break;
            case 'chat.message.received':
            case 'chatMessageReceived':
            case 'message.received':
            case 'ai.response':
            case 'assistant.message':
                result = await handleChatMessage(data, 'received', timestamp);
                break;
            case 'terminal.command.executed':
            case 'terminalCommandExecuted':
            case 'command.executed':
                result = await handleTerminalCommand(data, timestamp);
                break;
            case 'terminal.output.received':
            case 'command.output':
                result = await handleTerminalOutput(data, timestamp);
                break;
            // 新增：Kiro 工具調用監控
            case 'tool.fsWrite':
            case 'tool.fsAppend':
            case 'tool.strReplace':
                result = await handleToolExecution(data, event, timestamp);
                break;
            case 'agent.codeGeneration':
            case 'agent.documentGeneration':
            case 'agent.taskExecution':
                result = await handleAgentTask(data, event, timestamp);
                break;
            case 'kiro.conversation.turn':
                result = await handleConversationTurn(data, timestamp);
                break;
            default:
                console.log(`⚠️ Unknown event: ${event}`);
                console.log(`📋 Available data:`, JSON.stringify(context, null, 2));
                // 如果是未知事件但包含 kiro-usage.log，仍然嘗試分析
                if (JSON.stringify(context).includes('kiro-usage.log')) {
                    result = await analyzeUsageLog();
                } else {
                    // 嘗試通用處理
                    result = await handleGenericEvent(context, timestamp);
                }
        }
        
        // 記錄結果到日誌
        if (result) {
            await logUsage(result);
            console.log(`✅ Token Monitor: 分析完成 - ${result.tokens || 0} tokens`);
        }
        
        return { 
            success: true, 
            message: `Token monitoring completed for ${event}`,
            data: result
        };
        
    } catch (error) {
        console.error('❌ Token Monitor Error:', error.message);
        return { 
            success: false, 
            error: error.message 
        };
    }
}

/**
 * 處理聊天訊息
 */
async function handleChatMessage(data, direction, timestamp) {
    // 嘗試從不同的資料結構中提取內容
    let content = '';
    
    if (data && data.content) {
        content = data.content;
    } else if (data && data.message) {
        content = data.message;
    } else if (data && data.text) {
        content = data.text;
    } else if (typeof data === 'string') {
        content = data;
    } else {
        console.log('⚠️ 無法提取聊天內容，data:', JSON.stringify(data));
        return null;
    }
    
    if (!content || content.trim().length === 0) {
        console.log('⚠️ 聊天內容為空');
        return null;
    }
    
    const tokens = await calculateTokens(content);
    const activityType = classifyActivity(content);
    const costAnalysis = calculateCost(tokens, direction);
    
    const record = {
        timestamp: timestamp || new Date().toISOString(),
        event: 'chat_message',
        direction: direction,
        content_length: content.length,
        tokens: tokens,
        activity_type: activityType,
        model: (data && data.model) || 'claude-sonnet-4.0',
        session_id: (data && data.sessionId) || generateSessionId(),
        cost_analysis: costAnalysis
    };
    
    console.log(`📊 處理聊天訊息: ${direction}, ${tokens} tokens, 活動類型: ${activityType}`);
    return record;
}

/**
 * 處理檔案儲存
 */
async function handleFileSave(data, timestamp) {
    if (!data || !data.content) {
        return null;
    }
    
    const content = data.content;
    const tokens = await calculateTokens(content);
    const activityType = classifyActivityByFile(data.filePath || '', content);
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'file_save',
        file_path: data.filePath || 'unknown',
        file_type: path.extname(data.filePath || ''),
        content_length: content.length,
        tokens: tokens,
        activity_type: activityType
    };
}

/**
 * 處理終端機命令
 */
async function handleTerminalCommand(data, timestamp) {
    if (!data || !data.command) {
        return null;
    }
    
    const command = data.command;
    const tokens = await calculateTokens(command);
    const activityType = classifyTerminalActivity(command);
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'terminal_command',
        command: command,
        content_length: command.length,
        tokens: tokens,
        activity_type: activityType,
        working_directory: data.workingDirectory || process.cwd()
    };
}

/**
 * 處理終端機輸出
 */
async function handleTerminalOutput(data, timestamp) {
    if (!data || !data.output) {
        return null;
    }
    
    const output = data.output;
    const tokens = await calculateTokens(output);
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'terminal_output',
        output_length: output.length,
        tokens: tokens,
        activity_type: 'terminal_output',
        exit_code: data.exitCode || 0
    };
}

/**
 * 處理 Kiro 工具執行（fsWrite, fsAppend, strReplace 等）
 */
async function handleToolExecution(data, toolType, timestamp) {
    if (!data) {
        return null;
    }
    
    let content = '';
    let filePath = '';
    let operation = '';
    
    // 根據不同工具類型提取內容
    switch (toolType) {
        case 'tool.fsWrite':
            content = data.text || data.content || '';
            filePath = data.path || '';
            operation = 'write';
            break;
        case 'tool.fsAppend':
            content = data.text || data.content || '';
            filePath = data.path || '';
            operation = 'append';
            break;
        case 'tool.strReplace':
            content = data.newStr || '';
            filePath = data.path || '';
            operation = 'replace';
            break;
        default:
            content = JSON.stringify(data);
            operation = 'unknown';
    }
    
    if (!content || content.trim().length === 0) {
        return null;
    }
    
    const tokens = await calculateTokens(content);
    const activityType = classifyActivityByFile(filePath, content);
    const costAnalysis = calculateCost(tokens, 'output'); // 工具輸出視為 output
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'tool_execution',
        tool_type: toolType,
        operation: operation,
        file_path: filePath,
        file_type: path.extname(filePath || ''),
        content_length: content.length,
        tokens: tokens,
        activity_type: activityType,
        model: 'claude-sonnet-4.0',
        session_id: generateSessionId(),
        cost_analysis: costAnalysis,
        content_preview: content.substring(0, 100) + (content.length > 100 ? '...' : '')
    };
}

/**
 * 處理 Kiro 代理任務執行
 */
async function handleAgentTask(data, taskType, timestamp) {
    if (!data) {
        return null;
    }
    
    let content = '';
    let taskDescription = '';
    
    // 提取任務內容
    if (data.generatedContent) {
        content = data.generatedContent;
    } else if (data.output) {
        content = data.output;
    } else if (data.result) {
        content = typeof data.result === 'string' ? data.result : JSON.stringify(data.result);
    } else {
        content = JSON.stringify(data);
    }
    
    taskDescription = data.description || data.task || taskType;
    
    if (!content || content.trim().length === 0) {
        return null;
    }
    
    const tokens = await calculateTokens(content);
    const activityType = classifyAgentActivity(taskType, content);
    const costAnalysis = calculateCost(tokens, 'output');
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'agent_task',
        task_type: taskType,
        task_description: taskDescription,
        content_length: content.length,
        tokens: tokens,
        activity_type: activityType,
        model: 'claude-sonnet-4.0',
        session_id: data.sessionId || generateSessionId(),
        cost_analysis: costAnalysis,
        content_preview: content.substring(0, 100) + (content.length > 100 ? '...' : ''),
        execution_time: data.executionTime || null
    };
}

/**
 * 處理完整的對話回合
 */
async function handleConversationTurn(data, timestamp) {
    if (!data) {
        return null;
    }
    
    const userInput = data.userInput || data.input || '';
    const assistantOutput = data.assistantOutput || data.output || '';
    
    const inputTokens = await calculateTokens(userInput);
    const outputTokens = await calculateTokens(assistantOutput);
    const totalTokens = inputTokens + outputTokens;
    
    const inputCost = calculateCost(inputTokens, 'sent');
    const outputCost = calculateCost(outputTokens, 'received');
    const totalCost = inputCost.cost_usd + outputCost.cost_usd;
    
    const activityType = classifyActivity(userInput + ' ' + assistantOutput);
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'conversation_turn',
        user_input_length: userInput.length,
        assistant_output_length: assistantOutput.length,
        input_tokens: inputTokens,
        output_tokens: outputTokens,
        total_tokens: totalTokens,
        activity_type: activityType,
        model: data.model || 'claude-sonnet-4.0',
        session_id: data.sessionId || generateSessionId(),
        cost_analysis: {
            input_cost: inputCost.cost_usd,
            output_cost: outputCost.cost_usd,
            total_cost: totalCost,
            currency: 'USD'
        },
        tools_used: data.toolsUsed || [],
        execution_time: data.executionTime || null
    };
}

/**
 * 處理通用事件
 */
async function handleGenericEvent(context, timestamp) {
    if (!context || typeof context !== 'object') {
        return null;
    }
    
    // 嘗試從上下文中提取內容
    let content = '';
    if (context.content) content = context.content;
    else if (context.text) content = context.text;
    else if (context.message) content = context.message;
    else if (context.data && typeof context.data === 'string') content = context.data;
    else content = JSON.stringify(context);
    
    if (!content || content.trim().length === 0) {
        return null;
    }
    
    const tokens = await calculateTokens(content);
    const activityType = classifyActivity(content);
    
    return {
        timestamp: timestamp || new Date().toISOString(),
        event: 'generic_event',
        event_type: context.event || context.type || 'unknown',
        content_length: content.length,
        tokens: tokens,
        activity_type: activityType,
        raw_context: JSON.stringify(context).substring(0, 200) + '...'
    };
}

/**
 * 計算 Token 數量
 */
async function calculateTokens(text) {
    if (!text || text.trim().length === 0) {
        return 0;
    }
    
    try {
        // 使用 Token Monitor 執行檔計算
        const command = `"${TOKEN_MONITOR_PATH}" calculate "${text.replace(/"/g, '\\"')}"`;
        const output = execSync(command, { 
            encoding: 'utf8', 
            timeout: 10000,
            stdio: ['pipe', 'pipe', 'ignore'] // 忽略 stderr
        });
        
        // 解析輸出
        const match = output.match(/Token 數量:\s*(\d+)/);
        if (match) {
            return parseInt(match[1]);
        }
        
        // 如果解析失敗，使用估算
        return estimateTokens(text);
        
    } catch (error) {
        console.log(`⚠️ Token calculation failed, using estimation: ${error.message}`);
        return estimateTokens(text);
    }
}

/**
 * 簡單的 Token 估算
 */
function estimateTokens(text) {
    // 英文字符約 4 字符/token，中文字符約 1.5 字符/token
    const englishChars = (text.match(/[a-zA-Z0-9\s]/g) || []).length;
    const chineseChars = text.length - englishChars;
    
    return Math.ceil(englishChars / 4 + chineseChars / 1.5);
}

/**
 * 活動類型分類
 */
function classifyActivity(content) {
    const patterns = {
        coding: /(?:function|class|implement|程式|函數|變數|程式碼|寫.*程式|實作.*功能|建立.*函數)/i,
        debugging: /(?:error|bug|fix|錯誤|修復|除錯|修復.*問題|解決.*錯誤|debug)/i,
        documentation: /(?:README|document|文件|說明|註解|更新.*文件|撰寫.*說明|文檔)/i,
        'spec-development': /(?:spec|requirement|design|需求|設計|規格|需求分析|設計文件)/i,
        chat: /(?:chat|question|help|問題|協助|請問|如何|怎麼)/i
    };

    for (const [type, pattern] of Object.entries(patterns)) {
        if (pattern.test(content)) {
            return type;
        }
    }
    
    return 'general';
}

/**
 * 根據檔案類型分類活動
 */
function classifyActivityByFile(filePath, content) {
    const ext = path.extname(filePath).toLowerCase();
    const fileName = path.basename(filePath).toLowerCase();
    
    // 程式碼檔案
    if (['.js', '.ts', '.py', '.go', '.java', '.cpp', '.c', '.cs', '.php', '.rb'].includes(ext)) {
        return 'coding';
    }
    
    // 文件檔案
    if (['.md', '.txt', '.doc', '.docx', '.rst'].includes(ext) || fileName.includes('readme')) {
        return 'documentation';
    }
    
    // 規格檔案
    if (fileName.includes('spec') || fileName.includes('requirement') || fileName.includes('design')) {
        return 'spec-development';
    }
    
    // 配置檔案
    if (['.json', '.yaml', '.yml', '.xml', '.ini', '.conf'].includes(ext)) {
        return 'configuration';
    }
    
    // 根據內容分類
    return classifyActivity(content);
}

/**
 * 根據終端機命令分類活動
 */
function classifyTerminalActivity(command) {
    const cmd = command.toLowerCase().trim();
    
    // 建置相關命令
    if (/^(npm|yarn|pnpm)\s+(build|run\s+build|compile)/.test(cmd) || 
        /^(go\s+build|make|cmake|gradle|mvn\s+compile)/.test(cmd)) {
        return 'build';
    }
    
    // 測試相關命令
    if (/^(npm|yarn|pnpm)\s+(test|run\s+test)/.test(cmd) || 
        /^(go\s+test|pytest|jest|mocha)/.test(cmd)) {
        return 'testing';
    }
    
    // 部署相關命令
    if (/^(docker|kubectl|helm|terraform)/.test(cmd) || 
        /deploy|push|publish/.test(cmd)) {
        return 'deployment';
    }
    
    // Git 相關命令
    if (/^git\s+/.test(cmd)) {
        return 'version_control';
    }
    
    // 套件管理
    if (/^(npm|yarn|pnpm|pip|go\s+mod)\s+(install|add|get)/.test(cmd)) {
        return 'package_management';
    }
    
    // 檔案操作
    if (/^(ls|dir|cat|type|cp|copy|mv|move|rm|del)/.test(cmd)) {
        return 'file_operations';
    }
    
    // 開發工具
    if (/^(code|vim|nano|emacs)/.test(cmd)) {
        return 'development_tools';
    }
    
    return 'terminal_general';
}

/**
 * 根據代理任務類型分類活動
 */
function classifyAgentActivity(taskType, content) {
    // 根據任務類型分類
    if (taskType.includes('codeGeneration') || taskType.includes('code')) {
        return 'ai_code_generation';
    }
    
    if (taskType.includes('documentGeneration') || taskType.includes('document')) {
        return 'ai_documentation';
    }
    
    if (taskType.includes('taskExecution') || taskType.includes('task')) {
        return 'ai_task_execution';
    }
    
    // 根據內容分類
    const patterns = {
        'ai_code_generation': /(?:function|class|const|let|var|import|export|interface|type|async|await|=>)/i,
        'ai_documentation': /(?:README|\.md|documentation|說明|文件|註解|comment)/i,
        'ai_debugging': /(?:error|bug|fix|debug|錯誤|修復|除錯)/i,
        'ai_refactoring': /(?:refactor|optimize|improve|重構|優化|改進)/i,
        'ai_testing': /(?:test|spec|jest|mocha|測試|單元測試)/i,
        'ai_configuration': /(?:config|setting|\.json|\.yaml|\.yml|配置|設定)/i
    };
    
    for (const [type, pattern] of Object.entries(patterns)) {
        if (pattern.test(content)) {
            return type;
        }
    }
    
    return 'ai_general';
}

/**
 * 記錄使用量
 */
async function logUsage(record) {
    try {
        if (!record || !record.tokens) {
            return;
        }
        
        const logEntry = JSON.stringify(record) + '\n';
        fs.appendFileSync(LOG_FILE, logEntry, 'utf8');
        
        // 顯示即時統計
        displayRealtimeStats(record);
        
        return true;
        
    } catch (error) {
        console.error('Failed to log usage:', error.message);
        return false;
    }
}

/**
 * 顯示即時統計
 */
function displayRealtimeStats(record) {
    const timestamp = new Date(record.timestamp).toLocaleTimeString('zh-TW');
    const activity = record.activity_type;
    const tokens = record.tokens;
    const cost = record.cost_analysis ? record.cost_analysis.cost_usd || record.cost_analysis.total_cost : 0;
    
    console.log(`📊 [${timestamp}] ${activity}: ${tokens} tokens (${cost.toFixed(6)} USD)`);
    
    // 根據事件類型顯示詳細資訊
    switch (record.event) {
        case 'chat_message':
            const direction = record.direction === 'sent' ? '➡️ 用戶輸入' : '⬅️ AI 回應';
            console.log(`   ${direction} (${record.model})`);
            break;
            
        case 'file_save':
        case 'file_edit':
            const fileName = path.basename(record.file_path || '');
            console.log(`   💾 檔案: ${fileName} (${record.file_type})`);
            break;
            
        case 'tool_execution':
            console.log(`   🔧 工具: ${record.tool_type} - ${record.operation}`);
            if (record.file_path) {
                console.log(`   📁 檔案: ${path.basename(record.file_path)}`);
            }
            break;
            
        case 'agent_task':
            console.log(`   🤖 任務: ${record.task_type}`);
            if (record.task_description) {
                console.log(`   📝 描述: ${record.task_description.substring(0, 50)}...`);
            }
            break;
            
        case 'conversation_turn':
            console.log(`   💬 對話回合: 輸入${record.input_tokens} + 輸出${record.output_tokens} tokens`);
            if (record.tools_used && record.tools_used.length > 0) {
                console.log(`   🛠️ 使用工具: ${record.tools_used.join(', ')}`);
            }
            break;
            
        case 'terminal_command':
            console.log(`   💻 命令: ${record.command}`);
            break;
            
        case 'terminal_output':
            console.log(`   📤 輸出長度: ${record.output_length} 字符`);
            break;
            
        default:
            if (record.content_preview) {
                console.log(`   📄 內容預覽: ${record.content_preview}`);
            }
    }
    
    // 顯示成本資訊
    if (cost > 0) {
        console.log(`   💰 成本: ${cost.toFixed(6)} USD`);
    }
}

/**
 * 計算成本分析
 */
function calculateCost(tokens, direction) {
    const model = 'claude-sonnet-4.0';
    const pricing = {
        input: 3.0 / 1000000,   // $3 per 1M input tokens
        output: 15.0 / 1000000  // $15 per 1M output tokens
    };
    
    let cost = 0;
    let costType = '';
    
    if (direction === 'sent') {
        cost = tokens * pricing.input;
        costType = 'input';
    } else if (direction === 'received') {
        cost = tokens * pricing.output;
        costType = 'output';
    }
    
    return {
        tokens: tokens,
        cost_usd: parseFloat(cost.toFixed(6)),
        cost_type: costType,
        model: model,
        pricing_rate: direction === 'sent' ? pricing.input * 1000000 : pricing.output * 1000000
    };
}

/**
 * 分析整個使用記錄檔案
 */
async function analyzeUsageLog() {
    try {
        if (!fs.existsSync(LOG_FILE)) {
            console.log('⚠️ 使用記錄檔案不存在');
            return { success: false, message: 'Log file not found' };
        }
        
        const logContent = fs.readFileSync(LOG_FILE, 'utf8');
        const lines = logContent.trim().split('\n').filter(line => line.trim());
        
        if (lines.length === 0) {
            console.log('⚠️ 使用記錄檔案為空');
            return { success: false, message: 'Log file is empty' };
        }
        
        console.log(`📊 分析 ${lines.length} 筆記錄...`);
        
        let totalInputTokens = 0;
        let totalOutputTokens = 0;
        let totalCost = 0;
        let sessionStats = {};
        let activityStats = {};
        let modelStats = {};
        let recentRecords = [];
        
        // 分析每一筆記錄
        for (const line of lines) {
            try {
                const record = JSON.parse(line);
                
                // 統計 Token
                if (record.event === 'chat_message') {
                    if (record.direction === 'sent') {
                        totalInputTokens += record.tokens || 0;
                    } else if (record.direction === 'received') {
                        totalOutputTokens += record.tokens || 0;
                    }
                    
                    // 統計會話
                    const sessionId = record.session_id || 'unknown';
                    if (!sessionStats[sessionId]) {
                        sessionStats[sessionId] = { input: 0, output: 0, cost: 0 };
                    }
                    sessionStats[sessionId][record.direction === 'sent' ? 'input' : 'output'] += record.tokens || 0;
                    
                    // 統計模型使用
                    const model = record.model || 'unknown';
                    if (!modelStats[model]) {
                        modelStats[model] = { input: 0, output: 0, cost: 0 };
                    }
                    modelStats[model][record.direction === 'sent' ? 'input' : 'output'] += record.tokens || 0;
                }
                
                // 統計活動類型
                const activity = record.activity_type || 'unknown';
                if (!activityStats[activity]) {
                    activityStats[activity] = { count: 0, tokens: 0 };
                }
                activityStats[activity].count++;
                activityStats[activity].tokens += record.tokens || 0;
                
                // 統計成本
                if (record.cost_analysis && record.cost_analysis.cost_usd) {
                    totalCost += record.cost_analysis.cost_usd;
                    if (sessionStats[record.session_id]) {
                        sessionStats[record.session_id].cost += record.cost_analysis.cost_usd;
                    }
                    if (modelStats[record.model]) {
                        modelStats[record.model].cost += record.cost_analysis.cost_usd;
                    }
                }
                
                // 收集最近的記錄（最後10筆）
                recentRecords.push({
                    timestamp: record.timestamp,
                    event: record.event,
                    direction: record.direction,
                    tokens: record.tokens,
                    activity: record.activity_type,
                    model: record.model
                });
                
            } catch (parseError) {
                console.log(`⚠️ 無法解析記錄: ${line.substring(0, 50)}...`);
            }
        }
        
        // 只保留最近的10筆記錄
        recentRecords = recentRecords.slice(-10);
        
        // 生成分析報告
        const analysis = {
            總覽: {
                總記錄數: lines.length,
                輸入Token總數: totalInputTokens,
                輸出Token總數: totalOutputTokens,
                Token總數: totalInputTokens + totalOutputTokens,
                預估總成本: `$${totalCost.toFixed(6)} USD`
            },
            會話統計: Object.entries(sessionStats).map(([sessionId, stats]) => ({
                會話ID: sessionId,
                輸入Token: stats.input,
                輸出Token: stats.output,
                總Token: stats.input + stats.output,
                成本: `$${stats.cost.toFixed(6)} USD`
            })).slice(-5), // 只顯示最近5個會話
            活動類型統計: Object.entries(activityStats).map(([activity, stats]) => ({
                活動類型: activity,
                次數: stats.count,
                Token數: stats.tokens
            })),
            模型使用統計: Object.entries(modelStats).map(([model, stats]) => ({
                模型: model,
                輸入Token: stats.input,
                輸出Token: stats.output,
                總Token: stats.input + stats.output,
                成本: `$${stats.cost.toFixed(6)} USD`
            })),
            最近記錄: recentRecords.map(record => ({
                時間: new Date(record.timestamp).toLocaleString('zh-TW'),
                事件: record.event,
                方向: record.direction === 'sent' ? '發送' : '接收',
                Token: record.tokens,
                活動: record.activity,
                模型: record.model
            }))
        };
        
        // 輸出分析結果
        console.log('\n📊 ===== Kiro Chat Token 使用分析報告 =====');
        console.log(`📈 總記錄數: ${analysis.總覽.總記錄數}`);
        console.log(`🔢 輸入 Token: ${analysis.總覽.輸入Token總數}`);
        console.log(`🔢 輸出 Token: ${analysis.總覽.輸出Token總數}`);
        console.log(`🔢 總 Token: ${analysis.總覽.Token總數}`);
        console.log(`💰 預估成本: ${analysis.總覽.預估總成本}`);
        
        console.log('\n📋 活動類型統計:');
        analysis.活動類型統計.forEach(stat => {
            console.log(`  ${stat.活動類型}: ${stat.次數} 次, ${stat.Token數} tokens`);
        });
        
        console.log('\n🔄 最近記錄:');
        analysis.最近記錄.forEach(record => {
            console.log(`  [${record.時間}] ${record.事件} (${record.方向}) - ${record.Token} tokens - ${record.活動}`);
        });
        
        return {
            success: true,
            message: 'Token 分析完成',
            analysis: analysis,
            summary: {
                totalRecords: lines.length,
                inputTokens: totalInputTokens,
                outputTokens: totalOutputTokens,
                totalTokens: totalInputTokens + totalOutputTokens,
                totalCost: totalCost
            }
        };
        
    } catch (error) {
        console.error('❌ 分析使用記錄時發生錯誤:', error.message);
        return {
            success: false,
            error: error.message
        };
    }
}

/**
 * 生成會話ID
 */
function generateSessionId() {
    const timestamp = Date.now();
    const random = Math.random().toString(36).substring(2, 8);
    return `session-${timestamp}-${random}`;
}

/**
 * 生成使用報告
 */
async function generateReport() {
    try {
        const reportPath = path.join(DATA_DIR, 'kiro-usage-report.html');
        const command = `"${TOKEN_MONITOR_PATH}" report --format html --output "${reportPath}"`;
        
        execSync(command, { timeout: 30000 });
        console.log(`📊 Usage report generated: ${reportPath}`);
        
        return reportPath;
    } catch (error) {
        console.error('Failed to generate report:', error.message);
        return null;
    }
}

// 匯出主要函數供 Kiro IDE 使用
module.exports = {
    execute,
    calculateTokens,
    generateReport,
    analyzeUsageLog
};

// 如果直接執行此腳本，進行測試
if (require.main === module) {
    console.log('🧪 Testing Token Monitor Integration...');
    
    // 測試檔案儲存事件（模擬 Hook 觸發）
    const testContext = {
        event: 'file.saved',
        filePath: 'data/kiro-usage.log',
        timestamp: new Date().toISOString()
    };
    
    execute(testContext).then(result => {
        console.log('✅ Test completed:', JSON.stringify(result, null, 2));
    }).catch(error => {
        console.error('❌ Test failed:', error);
    });
}