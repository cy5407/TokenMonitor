/**
 * 假的測試程式碼 - 用於測試 Token 監控系統
 * 這段程式碼包含各種常見的程式設計模式和結構
 * 目的是產生足夠的內容來測試 Token 計算功能
 */

// 引入必要的模組
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');
const { EventEmitter } = require('events');

// 定義常數
const CONFIG = {
    MAX_RETRIES: 3,
    TIMEOUT: 5000,
    API_VERSION: 'v2.1',
    DEFAULT_ENCODING: 'utf8'
};

/**
 * 假的用戶管理類別
 * 這是一個模擬的用戶管理系統
 */
class FakeUserManager extends EventEmitter {
    constructor(options = {}) {
        super();
        this.users = new Map();
        this.config = { ...CONFIG, ...options };
        this.isInitialized = false;
        
        // 綁定方法到實例
        this.createUser = this.createUser.bind(this);
        this.deleteUser = this.deleteUser.bind(this);
        this.updateUser = this.updateUser.bind(this);
    }
    
    /**
     * 初始化用戶管理器
     * @param {Object} settings - 初始化設定
     * @returns {Promise<boolean>} 初始化結果
     */
    async initialize(settings = {}) {
        try {
            console.log('🚀 初始化用戶管理器...');
            
            // 模擬異步初始化過程
            await this.delay(1000);
            
            // 載入預設用戶
            const defaultUsers = [
                { id: 1, name: 'Alice', email: 'alice@example.com', role: 'admin' },
                { id: 2, name: 'Bob', email: 'bob@example.com', role: 'user' },
                { id: 3, name: 'Charlie', email: 'charlie@example.com', role: 'moderator' }
            ];
            
            for (const user of defaultUsers) {
                await this.createUser(user);
            }
            
            this.isInitialized = true;
            this.emit('initialized', { userCount: this.users.size });
            
            console.log('✅ 用戶管理器初始化完成');
            return true;
            
        } catch (error) {
            console.error('❌ 初始化失敗:', error.message);
            this.emit('error', error);
            return false;
        }
    }
    
    /**
     * 創建新用戶
     * @param {Object} userData - 用戶資料
     * @returns {Promise<Object>} 創建的用戶物件
     */
    async createUser(userData) {
        // 驗證輸入資料
        if (!userData || !userData.name || !userData.email) {
            throw new Error('用戶資料不完整');
        }
        
        // 檢查電子郵件是否已存在
        const existingUser = Array.from(this.users.values())
            .find(user => user.email === userData.email);
            
        if (existingUser) {
            throw new Error(`電子郵件 ${userData.email} 已被使用`);
        }
        
        // 生成用戶ID和密碼雜湊
        const userId = userData.id || this.generateUserId();
        const passwordHash = this.hashPassword(userData.password || 'defaultPassword123');
        
        // 創建用戶物件
        const user = {
            id: userId,
            name: userData.name,
            email: userData.email,
            role: userData.role || 'user',
            passwordHash: passwordHash,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            isActive: true,
            loginCount: 0,
            lastLoginAt: null,
            preferences: {
                theme: 'light',
                language: 'zh-TW',
                notifications: true
            }
        };
        
        // 儲存用戶
        this.users.set(userId, user);
        
        // 觸發事件
        this.emit('userCreated', { user: this.sanitizeUser(user) });
        
        console.log(`👤 用戶已創建: ${user.name} (${user.email})`);
        
        return this.sanitizeUser(user);
    }
    
    /**
     * 更新用戶資料
     * @param {number} userId - 用戶ID
     * @param {Object} updateData - 更新資料
     * @returns {Promise<Object>} 更新後的用戶物件
     */
    async updateUser(userId, updateData) {
        const user = this.users.get(userId);
        
        if (!user) {
            throw new Error(`找不到用戶 ID: ${userId}`);
        }
        
        // 驗證更新資料
        const allowedFields = ['name', 'email', 'role', 'preferences', 'isActive'];
        const updates = {};
        
        for (const [key, value] of Object.entries(updateData)) {
            if (allowedFields.includes(key)) {
                updates[key] = value;
            }
        }
        
        // 如果更新電子郵件，檢查是否重複
        if (updates.email && updates.email !== user.email) {
            const existingUser = Array.from(this.users.values())
                .find(u => u.email === updates.email && u.id !== userId);
                
            if (existingUser) {
                throw new Error(`電子郵件 ${updates.email} 已被使用`);
            }
        }
        
        // 應用更新
        const updatedUser = {
            ...user,
            ...updates,
            updatedAt: new Date().toISOString()
        };
        
        this.users.set(userId, updatedUser);
        
        // 觸發事件
        this.emit('userUpdated', { 
            userId, 
            updates, 
            user: this.sanitizeUser(updatedUser) 
        });
        
        console.log(`📝 用戶已更新: ${updatedUser.name}`);
        
        return this.sanitizeUser(updatedUser);
    }
    
    /**
     * 刪除用戶
     * @param {number} userId - 用戶ID
     * @returns {Promise<boolean>} 刪除結果
     */
    async deleteUser(userId) {
        const user = this.users.get(userId);
        
        if (!user) {
            throw new Error(`找不到用戶 ID: ${userId}`);
        }
        
        // 軟刪除：標記為非活躍
        const deletedUser = {
            ...user,
            isActive: false,
            deletedAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };
        
        this.users.set(userId, deletedUser);
        
        // 觸發事件
        this.emit('userDeleted', { userId, user: this.sanitizeUser(user) });
        
        console.log(`🗑️ 用戶已刪除: ${user.name}`);
        
        return true;
    }
    
    /**
     * 獲取用戶列表
     * @param {Object} filters - 篩選條件
     * @returns {Array} 用戶列表
     */
    getUsers(filters = {}) {
        let users = Array.from(this.users.values());
        
        // 應用篩選條件
        if (filters.role) {
            users = users.filter(user => user.role === filters.role);
        }
        
        if (filters.isActive !== undefined) {
            users = users.filter(user => user.isActive === filters.isActive);
        }
        
        if (filters.search) {
            const searchTerm = filters.search.toLowerCase();
            users = users.filter(user => 
                user.name.toLowerCase().includes(searchTerm) ||
                user.email.toLowerCase().includes(searchTerm)
            );
        }
        
        // 排序
        users.sort((a, b) => {
            if (filters.sortBy === 'name') {
                return a.name.localeCompare(b.name);
            } else if (filters.sortBy === 'createdAt') {
                return new Date(b.createdAt) - new Date(a.createdAt);
            }
            return a.id - b.id;
        });
        
        // 分頁
        if (filters.page && filters.limit) {
            const start = (filters.page - 1) * filters.limit;
            const end = start + filters.limit;
            users = users.slice(start, end);
        }
        
        return users.map(user => this.sanitizeUser(user));
    }
    
    /**
     * 用戶登入
     * @param {string} email - 電子郵件
     * @param {string} password - 密碼
     * @returns {Promise<Object>} 登入結果
     */
    async login(email, password) {
        const user = Array.from(this.users.values())
            .find(u => u.email === email && u.isActive);
            
        if (!user) {
            throw new Error('用戶不存在或已被停用');
        }
        
        // 驗證密碼
        const passwordHash = this.hashPassword(password);
        if (passwordHash !== user.passwordHash) {
            throw new Error('密碼錯誤');
        }
        
        // 更新登入資訊
        const updatedUser = {
            ...user,
            loginCount: user.loginCount + 1,
            lastLoginAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
        };
        
        this.users.set(user.id, updatedUser);
        
        // 生成假的 JWT Token
        const token = this.generateToken(user.id);
        
        // 觸發事件
        this.emit('userLoggedIn', { 
            userId: user.id, 
            email: user.email,
            loginCount: updatedUser.loginCount
        });
        
        console.log(`🔐 用戶登入: ${user.name}`);
        
        return {
            user: this.sanitizeUser(updatedUser),
            token: token,
            expiresIn: 3600 // 1小時
        };
    }
    
    /**
     * 生成用戶ID
     * @returns {number} 新的用戶ID
     */
    generateUserId() {
        const existingIds = Array.from(this.users.keys());
        return existingIds.length > 0 ? Math.max(...existingIds) + 1 : 1;
    }
    
    /**
     * 雜湊密碼
     * @param {string} password - 原始密碼
     * @returns {string} 雜湊後的密碼
     */
    hashPassword(password) {
        return crypto.createHash('sha256')
            .update(password + 'fake-salt-12345')
            .digest('hex');
    }
    
    /**
     * 生成假的 Token
     * @param {number} userId - 用戶ID
     * @returns {string} Token
     */
    generateToken(userId) {
        const payload = {
            userId: userId,
            iat: Math.floor(Date.now() / 1000),
            exp: Math.floor(Date.now() / 1000) + 3600
        };
        
        return Buffer.from(JSON.stringify(payload)).toString('base64');
    }
    
    /**
     * 清理用戶物件（移除敏感資訊）
     * @param {Object} user - 用戶物件
     * @returns {Object} 清理後的用戶物件
     */
    sanitizeUser(user) {
        const { passwordHash, ...sanitizedUser } = user;
        return sanitizedUser;
    }
    
    /**
     * 延遲函數
     * @param {number} ms - 延遲毫秒數
     * @returns {Promise} Promise物件
     */
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    /**
     * 獲取統計資訊
     * @returns {Object} 統計資訊
     */
    getStatistics() {
        const users = Array.from(this.users.values());
        const activeUsers = users.filter(user => user.isActive);
        const roleStats = {};
        
        users.forEach(user => {
            roleStats[user.role] = (roleStats[user.role] || 0) + 1;
        });
        
        return {
            totalUsers: users.length,
            activeUsers: activeUsers.length,
            inactiveUsers: users.length - activeUsers.length,
            roleDistribution: roleStats,
            averageLoginCount: users.reduce((sum, user) => sum + user.loginCount, 0) / users.length
        };
    }
}

/**
 * 假的資料庫連接類別
 */
class FakeDatabase {
    constructor(connectionString) {
        this.connectionString = connectionString;
        this.isConnected = false;
        this.tables = new Map();
    }
    
    async connect() {
        console.log('🔌 連接到資料庫...');
        await this.delay(500);
        this.isConnected = true;
        console.log('✅ 資料庫連接成功');
    }
    
    async disconnect() {
        console.log('🔌 斷開資料庫連接...');
        this.isConnected = false;
        console.log('✅ 資料庫連接已斷開');
    }
    
    async query(sql, params = []) {
        if (!this.isConnected) {
            throw new Error('資料庫未連接');
        }
        
        console.log(`📊 執行查詢: ${sql}`);
        
        // 模擬查詢延遲
        await this.delay(Math.random() * 100 + 50);
        
        // 返回假的查詢結果
        return {
            rows: [
                { id: 1, name: 'Sample Data 1', value: Math.random() * 100 },
                { id: 2, name: 'Sample Data 2', value: Math.random() * 100 }
            ],
            rowCount: 2,
            executionTime: Math.random() * 50 + 10
        };
    }
    
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

/**
 * 主要的測試函數
 */
async function runFakeTest() {
    console.log('🧪 開始執行假的測試程式...');
    
    try {
        // 創建用戶管理器
        const userManager = new FakeUserManager({
            MAX_RETRIES: 5,
            TIMEOUT: 10000
        });
        
        // 設置事件監聽器
        userManager.on('initialized', (data) => {
            console.log(`📊 用戶管理器已初始化，用戶數量: ${data.userCount}`);
        });
        
        userManager.on('userCreated', (data) => {
            console.log(`👤 新用戶創建事件: ${data.user.name}`);
        });
        
        userManager.on('error', (error) => {
            console.error(`❌ 錯誤事件: ${error.message}`);
        });
        
        // 初始化用戶管理器
        await userManager.initialize();
        
        // 創建新用戶
        const newUser = await userManager.createUser({
            name: 'David Chen',
            email: 'david@example.com',
            role: 'developer',
            password: 'securePassword456'
        });
        
        console.log('新用戶創建成功:', newUser);
        
        // 更新用戶
        const updatedUser = await userManager.updateUser(newUser.id, {
            name: 'David Chen (Senior)',
            preferences: {
                theme: 'dark',
                language: 'en-US',
                notifications: false
            }
        });
        
        console.log('用戶更新成功:', updatedUser);
        
        // 用戶登入
        const loginResult = await userManager.login('david@example.com', 'securePassword456');
        console.log('登入成功:', loginResult);
        
        // 獲取用戶列表
        const users = userManager.getUsers({
            isActive: true,
            sortBy: 'name',
            page: 1,
            limit: 10
        });
        
        console.log('活躍用戶列表:', users);
        
        // 獲取統計資訊
        const stats = userManager.getStatistics();
        console.log('用戶統計:', stats);
        
        // 測試資料庫操作
        const db = new FakeDatabase('fake://localhost:5432/testdb');
        await db.connect();
        
        const queryResult = await db.query(
            'SELECT * FROM users WHERE role = ? AND created_at > ?',
            ['admin', '2024-01-01']
        );
        
        console.log('查詢結果:', queryResult);
        
        await db.disconnect();
        
        // 模擬一些複雜的業務邏輯
        const businessLogicResult = await performComplexBusinessLogic(userManager, users);
        console.log('業務邏輯執行結果:', businessLogicResult);
        
        console.log('✅ 假的測試程式執行完成');
        
        return {
            success: true,
            userCount: users.length,
            statistics: stats,
            executionTime: Date.now()
        };
        
    } catch (error) {
        console.error('❌ 測試執行失敗:', error.message);
        return {
            success: false,
            error: error.message,
            executionTime: Date.now()
        };
    }
}

/**
 * 複雜的業務邏輯函數
 */
async function performComplexBusinessLogic(userManager, users) {
    console.log('🔄 執行複雜業務邏輯...');
    
    const results = [];
    
    // 模擬批次處理
    for (let i = 0; i < users.length; i++) {
        const user = users[i];
        
        // 模擬一些計算
        const score = calculateUserScore(user);
        const category = categorizeUser(user, score);
        const recommendations = generateRecommendations(user, category);
        
        results.push({
            userId: user.id,
            userName: user.name,
            score: score,
            category: category,
            recommendations: recommendations,
            processedAt: new Date().toISOString()
        });
        
        // 模擬處理延遲
        await new Promise(resolve => setTimeout(resolve, 10));
    }
    
    // 模擬聚合計算
    const aggregatedData = {
        totalProcessed: results.length,
        averageScore: results.reduce((sum, r) => sum + r.score, 0) / results.length,
        categoryDistribution: results.reduce((acc, r) => {
            acc[r.category] = (acc[r.category] || 0) + 1;
            return acc;
        }, {}),
        topUsers: results
            .sort((a, b) => b.score - a.score)
            .slice(0, 3)
            .map(r => ({ name: r.userName, score: r.score }))
    };
    
    console.log('✅ 複雜業務邏輯執行完成');
    
    return {
        individualResults: results,
        aggregatedData: aggregatedData,
        processingTime: Date.now()
    };
}

/**
 * 計算用戶分數
 */
function calculateUserScore(user) {
    let score = 0;
    
    // 基於登入次數
    score += user.loginCount * 10;
    
    // 基於角色
    const roleScores = { admin: 100, moderator: 75, developer: 60, user: 30 };
    score += roleScores[user.role] || 0;
    
    // 基於帳戶年齡
    const accountAge = Date.now() - new Date(user.createdAt).getTime();
    const daysSinceCreation = accountAge / (1000 * 60 * 60 * 24);
    score += Math.min(daysSinceCreation * 2, 200);
    
    // 基於活躍度
    if (user.isActive) score += 50;
    if (user.lastLoginAt) {
        const daysSinceLogin = (Date.now() - new Date(user.lastLoginAt).getTime()) / (1000 * 60 * 60 * 24);
        if (daysSinceLogin < 7) score += 30;
        else if (daysSinceLogin < 30) score += 15;
    }
    
    return Math.round(score);
}

/**
 * 用戶分類
 */
function categorizeUser(user, score) {
    if (score >= 300) return 'VIP';
    if (score >= 200) return 'Premium';
    if (score >= 100) return 'Regular';
    if (score >= 50) return 'Basic';
    return 'New';
}

/**
 * 生成推薦
 */
function generateRecommendations(user, category) {
    const recommendations = [];
    
    if (category === 'VIP') {
        recommendations.push('專屬客服支援', '進階功能存取', '優先技術支援');
    } else if (category === 'Premium') {
        recommendations.push('功能升級建議', '社群活動邀請', '進階教學資源');
    } else if (category === 'Regular') {
        recommendations.push('功能探索指南', '社群參與建議', '技能提升課程');
    } else if (category === 'Basic') {
        recommendations.push('基礎功能教學', '入門指南', '社群介紹');
    } else {
        recommendations.push('歡迎教學', '快速入門', '基礎設定指導');
    }
    
    // 基於用戶偏好的個人化推薦
    if (user.preferences) {
        if (user.preferences.theme === 'dark') {
            recommendations.push('深色主題優化建議');
        }
        if (user.preferences.notifications === false) {
            recommendations.push('重要通知設定建議');
        }
    }
    
    return recommendations;
}

// 如果直接執行此檔案
if (require.main === module) {
    runFakeTest()
        .then(result => {
            console.log('\n📋 最終執行結果:', JSON.stringify(result, null, 2));
            process.exit(0);
        })
        .catch(error => {
            console.error('\n❌ 程式執行錯誤:', error);
            process.exit(1);
        });
}

// 匯出模組
module.exports = {
    FakeUserManager,
    FakeDatabase,
    runFakeTest,
    performComplexBusinessLogic,
    calculateUserScore,
    categorizeUser,
    generateRecommendations
};