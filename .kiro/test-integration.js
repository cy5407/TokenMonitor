/**
 * Token Monitor Kiro IDE 整合測試腳本
 */

const fs = require('fs');
const path = require('path');

class IntegrationTester {
    constructor() {
        this.baseDir = __dirname;
        this.results = [];
    }

    /**
     * 執行所有整合測試
     */
    async runAllTests() {
        console.log('🧪 開始 Token Monitor Kiro IDE 整合測試...\n');

        await this.testHookFiles();
        await this.testCommandFiles();
        await this.testStatusBarFiles();
        await this.testSettingsFiles();
        await this.testTokenMonitorExecutable();

        this.printResults();
    }

    /**
     * 測試 Hook 檔案
     */
    async testHookFiles() {
        console.log('📋 測試方案1: Agent Hooks');
        
        const hookJson = path.join(this.baseDir, 'hooks', 'token-monitor-hook.json');
        const hookJs = path.join(this.baseDir, 'hooks', 'token-monitor-hook.js');

        this.checkFile(hookJson, '方案1: Hook 配置檔案');
        this.checkFile(hookJs, '方案1: Hook 執行腳本');

        // 測試 Hook 腳本語法
        try {
            require(hookJs);
            this.addResult('✅', '方案1: Hook 腳本語法正確');
        } catch (error) {
            this.addResult('❌', `方案1: Hook 腳本語法錯誤 - ${error.message}`);
        }

        console.log('');
    }

    /**
     * 測試命令檔案
     */
    async testCommandFiles() {
        console.log('📋 測試方案2: 命令面板');
        
        const commandsFile = path.join(this.baseDir, 'commands', 'token-monitor-commands.json');
        this.checkFile(commandsFile, '方案2: 命令配置檔案');

        // 測試 JSON 格式
        try {
            const commands = JSON.parse(fs.readFileSync(commandsFile, 'utf8'));
            this.addResult('✅', `方案2: 找到 ${commands.commands.length} 個命令`);
        } catch (error) {
            this.addResult('❌', `方案2: 命令檔案格式錯誤 - ${error.message}`);
        }

        console.log('');
    }

    /**
     * 測試狀態列檔案
     */
    async testStatusBarFiles() {
        console.log('📋 測試方案3: 狀態列顯示');
        
        const statusBarFile = path.join(this.baseDir, 'statusbar', 'token-counter.js');
        this.checkFile(statusBarFile, '方案3: 狀態列腳本');

        // 測試狀態列腳本
        try {
            const statusBar = require(statusBarFile);
            if (statusBar.initialize && statusBar.update) {
                this.addResult('✅', '方案3: 狀態列 API 正確');
            } else {
                this.addResult('❌', '方案3: 狀態列 API 不完整');
            }
        } catch (error) {
            this.addResult('❌', `方案3: 狀態列腳本錯誤 - ${error.message}`);
        }

        console.log('');
    }

    /**
     * 測試設定檔案
     */
    async testSettingsFiles() {
        console.log('📋 測試方案4: 設定檔案');
        
        const settingsFile = path.join(this.baseDir, 'settings', 'token-monitor.json');
        this.checkFile(settingsFile, '方案4: 設定檔案');

        // 測試設定檔案格式
        try {
            const settings = JSON.parse(fs.readFileSync(settingsFile, 'utf8'));
            if (settings.tokenMonitor) {
                this.addResult('✅', '方案4: 設定檔案格式正確');
                this.addResult('ℹ️', `方案4: 監控已${settings.tokenMonitor.enabled ? '啟用' : '停用'}`);
            } else {
                this.addResult('❌', '方案4: 設定檔案結構錯誤');
            }
        } catch (error) {
            this.addResult('❌', `方案4: 設定檔案格式錯誤 - ${error.message}`);
        }

        console.log('');
    }

    /**
     * 測試 Token Monitor 執行檔
     */
    async testTokenMonitorExecutable() {
        console.log('📋 測試 Token Monitor 執行檔');
        
        const exePath = path.join(this.baseDir, '..', 'token-monitor.exe');
        
        if (fs.existsSync(exePath)) {
            this.addResult('✅', 'Token Monitor 執行檔存在');
            
            // 測試執行檔是否可用
            try {
                const { execSync } = require('child_process');
                const output = execSync(`"${exePath}" calculate "test" --quiet`, { 
                    encoding: 'utf8', 
                    timeout: 5000,
                    stdio: ['pipe', 'pipe', 'ignore']
                });
                this.addResult('✅', 'Token Monitor 執行檔可正常運作');
            } catch (error) {
                this.addResult('⚠️', 'Token Monitor 執行檔可能有問題');
            }
        } else {
            this.addResult('❌', 'Token Monitor 執行檔不存在');
        }

        console.log('');
    }

    /**
     * 檢查檔案是否存在
     */
    checkFile(filePath, description) {
        if (fs.existsSync(filePath)) {
            const stats = fs.statSync(filePath);
            this.addResult('✅', `${description} (${(stats.size / 1024).toFixed(1)}KB)`);
        } else {
            this.addResult('❌', `${description} - 檔案不存在`);
        }
    }

    /**
     * 添加測試結果
     */
    addResult(status, message) {
        this.results.push({ status, message });
        console.log(`${status} ${message}`);
    }

    /**
     * 印出測試結果摘要
     */
    printResults() {
        console.log('\n' + '='.repeat(60));
        console.log('📊 整合測試結果摘要');
        console.log('='.repeat(60));

        const summary = {
            '✅': this.results.filter(r => r.status === '✅').length,
            '❌': this.results.filter(r => r.status === '❌').length,
            '⚠️': this.results.filter(r => r.status === '⚠️').length,
            'ℹ️': this.results.filter(r => r.status === 'ℹ️').length
        };

        console.log(`✅ 成功: ${summary['✅']} 項`);
        console.log(`❌ 失敗: ${summary['❌']} 項`);
        console.log(`⚠️ 警告: ${summary['⚠️']} 項`);
        console.log(`ℹ️ 資訊: ${summary['ℹ️']} 項`);

        console.log('\n📋 建議：');
        if (summary['❌'] === 0) {
            console.log('🎉 所有整合檔案都已正確部署！');
            console.log('💡 現在可以重新啟動 Kiro IDE 來啟用 Token Monitor 功能');
        } else {
            console.log('🔧 請修復上述錯誤後重新測試');
        }

        console.log('\n🚀 使用方式：');
        console.log('1. 重新啟動 Kiro IDE');
        console.log('2. 檢查狀態列是否顯示 Token 計數器');
        console.log('3. 按 Ctrl+Shift+P 搜尋 "Token Monitor" 命令');
        console.log('4. 發送聊天訊息測試自動監控功能');
    }
}

// 執行測試
if (require.main === module) {
    const tester = new IntegrationTester();
    tester.runAllTests().catch(console.error);
}

module.exports = { IntegrationTester };