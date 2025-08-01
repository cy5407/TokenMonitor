# 現代程式設計指南

## JavaScript 最佳實踐

### 1. 使用現代語法

```javascript
// 使用 const 和 let 而不是 var
const API_URL = 'https://api.example.com';
let userCount = 0;

// 使用箭頭函數
const calculateTotal = (items) => {
    return items.reduce((sum, item) => sum + item.price, 0);
};

// 使用解構賦值
const { name, email, age } = user;
const [first, second, ...rest] = numbers;

// 使用模板字符串
const message = `Hello ${name}, you have ${userCount} notifications`;
```

### 2. 異步程式設計

```javascript
// 使用 async/await 而不是 Promise.then()
async function fetchUserData(userId) {
    try {
        const response = await fetch(`${API_URL}/users/${userId}`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const userData = await response.json();
        return userData;
    } catch (error) {
        console.error('Failed to fetch user data:', error);
        throw error;
    }
}

// 並行處理多個異步操作
async function fetchMultipleUsers(userIds) {
    const promises = userIds.map(id => fetchUserData(id));
    
    try {
        const users = await Promise.all(promises);
        return users;
    } catch (error) {
        console.error('Failed to fetch multiple users:', error);
        return [];
    }
}
```

### 3. 錯誤處理

```javascript
// 建立自定義錯誤類別
class ValidationError extends Error {
    constructor(message, field) {
        super(message);
        this.name = 'ValidationError';
        this.field = field;
    }
}

// 輸入驗證函數
function validateEmail(email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    
    if (!email) {
        throw new ValidationError('Email is required', 'email');
    }
    
    if (!emailRegex.test(email)) {
        throw new ValidationError('Invalid email format', 'email');
    }
    
    return true;
}

// 使用 try-catch 處理錯誤
function processUserRegistration(userData) {
    try {
        validateEmail(userData.email);
        
        // 處理註冊邏輯
        return registerUser(userData);
    } catch (error) {
        if (error instanceof ValidationError) {
            console.error(`Validation failed for ${error.field}: ${error.message}`);
        } else {
            console.error('Unexpected error:', error);
        }
        
        throw error;
    }
}
```

## TypeScript 進階技巧

### 1. 型別定義

```typescript
// 介面定義
interface User {
    id: number;
    name: string;
    email: string;
    age?: number; // 可選屬性
    readonly createdAt: Date; // 只讀屬性
}

// 泛型介面
interface ApiResponse<T> {
    success: boolean;
    data: T;
    message?: string;
}

// 聯合型別
type Status = 'pending' | 'approved' | 'rejected';

// 映射型別
type PartialUser = Partial<User>;
type RequiredUser = Required<User>;
```

### 2. 進階型別操作

```typescript
// 條件型別
type NonNullable<T> = T extends null | undefined ? never : T;

// 工具型別
type UserKeys = keyof User; // 'id' | 'name' | 'email' | 'age' | 'createdAt'
type UserName = Pick<User, 'name'>; // { name: string }
type UserWithoutId = Omit<User, 'id'>; // User 但沒有 id 屬性

// 函數重載
function processData(data: string): string;
function processData(data: number): number;
function processData(data: boolean): boolean;
function processData(data: any): any {
    if (typeof data === 'string') {
        return data.toUpperCase();
    } else if (typeof data === 'number') {
        return data * 2;
    } else if (typeof data === 'boolean') {
        return !data;
    }
}
```

## React 現代開發模式

### 1. 函數式元件與 Hooks

```jsx
import React, { useState, useEffect, useCallback, useMemo } from 'react';

// 自定義 Hook
function useLocalStorage(key, initialValue) {
    const [storedValue, setStoredValue] = useState(() => {
        try {
            const item = window.localStorage.getItem(key);
            return item ? JSON.parse(item) : initialValue;
        } catch (error) {
            console.error('Error reading from localStorage:', error);
            return initialValue;
        }
    });

    const setValue = useCallback((value) => {
        try {
            setStoredValue(value);
            window.localStorage.setItem(key, JSON.stringify(value));
        } catch (error) {
            console.error('Error writing to localStorage:', error);
        }
    }, [key]);

    return [storedValue, setValue];
}

// 使用自定義 Hook 的元件
function UserProfile({ userId }) {
    const [user, setUser] = useState(null);
    const [loading, setLoading] = useState(true);
    const [preferences, setPreferences] = useLocalStorage('userPreferences', {});

    // 使用 useEffect 處理副作用
    useEffect(() => {
        let isCancelled = false;

        async function loadUser() {
            try {
                setLoading(true);
                const userData = await fetchUserData(userId);
                
                if (!isCancelled) {
                    setUser(userData);
                }
            } catch (error) {
                if (!isCancelled) {
                    console.error('Failed to load user:', error);
                }
            } finally {
                if (!isCancelled) {
                    setLoading(false);
                }
            }
        }

        loadUser();

        // 清理函數
        return () => {
            isCancelled = true;
        };
    }, [userId]);

    // 使用 useMemo 優化計算
    const displayName = useMemo(() => {
        if (!user) return '';
        return `${user.name} (${user.email})`;
    }, [user]);

    // 使用 useCallback 優化函數
    const handlePreferenceChange = useCallback((key, value) => {
        setPreferences(prev => ({
            ...prev,
            [key]: value
        }));
    }, [setPreferences]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (!user) {
        return <div>User not found</div>;
    }

    return (
        <div className="user-profile">
            <h1>{displayName}</h1>
            <p>Age: {user.age || 'Not specified'}</p>
            <p>Member since: {new Date(user.createdAt).toLocaleDateString()}</p>
            
            <div className="preferences">
                <h2>Preferences</h2>
                <label>
                    <input
                        type="checkbox"
                        checked={preferences.notifications || false}
                        onChange={(e) => handlePreferenceChange('notifications', e.target.checked)}
                    />
                    Enable notifications
                </label>
            </div>
        </div>
    );
}
```

### 2. 狀態管理

```jsx
// 使用 useReducer 管理複雜狀態
import React, { useReducer, createContext, useContext } from 'react';

// 定義狀態和動作
const initialState = {
    users: [],
    loading: false,
    error: null,
    selectedUser: null
};

function userReducer(state, action) {
    switch (action.type) {
        case 'FETCH_USERS_START':
            return {
                ...state,
                loading: true,
                error: null
            };
        
        case 'FETCH_USERS_SUCCESS':
            return {
                ...state,
                loading: false,
                users: action.payload,
                error: null
            };
        
        case 'FETCH_USERS_ERROR':
            return {
                ...state,
                loading: false,
                error: action.payload
            };
        
        case 'SELECT_USER':
            return {
                ...state,
                selectedUser: action.payload
            };
        
        default:
            return state;
    }
}

// 建立 Context
const UserContext = createContext();

// Provider 元件
export function UserProvider({ children }) {
    const [state, dispatch] = useReducer(userReducer, initialState);

    const actions = {
        fetchUsers: async () => {
            dispatch({ type: 'FETCH_USERS_START' });
            
            try {
                const users = await fetchMultipleUsers([1, 2, 3, 4, 5]);
                dispatch({ type: 'FETCH_USERS_SUCCESS', payload: users });
            } catch (error) {
                dispatch({ type: 'FETCH_USERS_ERROR', payload: error.message });
            }
        },
        
        selectUser: (user) => {
            dispatch({ type: 'SELECT_USER', payload: user });
        }
    };

    return (
        <UserContext.Provider value={{ state, actions }}>
            {children}
        </UserContext.Provider>
    );
}

// 自定義 Hook 使用 Context
export function useUsers() {
    const context = useContext(UserContext);
    
    if (!context) {
        throw new Error('useUsers must be used within a UserProvider');
    }
    
    return context;
}
```

## 總結

現代程式設計強調：

1. **可讀性** - 程式碼應該易於理解和維護
2. **可測試性** - 函數應該純粹且易於測試
3. **效能** - 適當使用記憶化和優化技巧
4. **型別安全** - 使用 TypeScript 提供更好的開發體驗
5. **錯誤處理** - 優雅地處理各種錯誤情況

這些實踐將幫助你寫出更高品質、更可維護的程式碼。