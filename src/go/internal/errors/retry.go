package errors

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"
)

// RetryManager 重試管理器
type RetryManager struct {
	mu          sync.RWMutex
	activeRetries map[string]*RetryState
}

// RetryState 重試狀態
type RetryState struct {
	Attempts    int           `json:"attempts"`
	LastAttempt time.Time     `json:"last_attempt"`
	NextAttempt time.Time     `json:"next_attempt"`
	Policy      *RetryPolicy  `json:"policy"`
	Delay       time.Duration `json:"delay"`
}

// RetryResult 重試結果
type RetryResult struct {
	Success      bool          `json:"success"`
	Attempts     int           `json:"attempts"`
	TotalTime    time.Duration `json:"total_time"`
	LastError    error         `json:"last_error,omitempty"`
	GaveUpReason string        `json:"gave_up_reason,omitempty"`
}

// NewRetryManager 建立新的重試管理器
func NewRetryManager() *RetryManager {
	return &RetryManager{
		activeRetries: make(map[string]*RetryState),
	}
}

// Execute 執行帶重試的操作
func (rm *RetryManager) Execute(ctx context.Context, fn func() error, policy *RetryPolicy) (bool, error) {
	if policy == nil {
		policy = &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 30,
			BackoffFactor: 2.0,
		}
	}

	_ = time.Now()
	var lastErr error
	
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		// 檢查上下文是否已取消
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		// 執行操作
		err := fn()
		if err == nil {
			return true, nil // 成功
		}

		lastErr = err

		// 檢查是否為可重試的錯誤
		if !rm.isRetryableError(err, policy) {
			return false, err
		}

		// 如果還有重試機會
		if attempt < policy.MaxRetries {
			delay := rm.calculateDelay(attempt, policy)
			
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(delay):
				// 繼續下一次重試
			}
		}
	}

	return false, lastErr
}

// ExecuteWithJitter 執行帶抖動的重試
func (rm *RetryManager) ExecuteWithJitter(ctx context.Context, fn func() error, policy *RetryPolicy, jitterPercent float64) (bool, error) {
	if jitterPercent < 0 || jitterPercent > 1 {
		jitterPercent = 0.1 // 預設 10% 抖動
	}

	_ = time.Now()
	var lastErr error
	
	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}

		err := fn()
		if err == nil {
			return true, nil
		}

		lastErr = err

		if !rm.isRetryableError(err, policy) {
			return false, err
		}

		if attempt < policy.MaxRetries {
			baseDelay := rm.calculateDelay(attempt, policy)
			jitteredDelay := rm.addJitter(baseDelay, jitterPercent)
			
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-time.After(jitteredDelay):
			}
		}
	}

	return false, lastErr
}

// ExecuteAsync 非同步執行重試
func (rm *RetryManager) ExecuteAsync(ctx context.Context, fn func() error, policy *RetryPolicy, callback func(bool, error)) {
	go func() {
		success, err := rm.Execute(ctx, fn, policy)
		if callback != nil {
			callback(success, err)
		}
	}()
}

// GetRetryState 取得重試狀態
func (rm *RetryManager) GetRetryState(key string) (*RetryState, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	state, exists := rm.activeRetries[key]
	return state, exists
}

// SetRetryState 設定重試狀態
func (rm *RetryManager) SetRetryState(key string, state *RetryState) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	rm.activeRetries[key] = state
}

// RemoveRetryState 移除重試狀態
func (rm *RetryManager) RemoveRetryState(key string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	delete(rm.activeRetries, key)
}

// isRetryableError 檢查錯誤是否可重試
func (rm *RetryManager) isRetryableError(err error, policy *RetryPolicy) bool {
	if appErr, ok := err.(*AppError); ok {
		// 檢查策略中的可重試錯誤清單
		if policy.RetryableErrors != nil {
			for _, retryableCode := range policy.RetryableErrors {
				if appErr.Code == retryableCode {
					return true
				}
			}
			return false
		}
		
		// 使用錯誤的內建可重試性
		return appErr.IsRetryable()
	}
	
	// 非 AppError 的通用重試邏輯
	errorStr := err.Error()
	retryablePatterns := []string{
		"timeout", "connection", "network", "temporary",
		"超時", "連接", "網路", "暫時",
	}
	
	for _, pattern := range retryablePatterns {
		if contains(errorStr, pattern) {
			return true
		}
	}
	
	return false
}

// calculateDelay 計算延遲時間
func (rm *RetryManager) calculateDelay(attempt int, policy *RetryPolicy) time.Duration {
	if attempt == 0 {
		return policy.InitialDelay
	}
	
	// 指數退避
	delay := float64(policy.InitialDelay) * math.Pow(policy.BackoffFactor, float64(attempt))
	
	// 限制最大延遲
	if time.Duration(delay) > policy.MaxDelay {
		delay = float64(policy.MaxDelay)
	}
	
	return time.Duration(delay)
}

// addJitter 添加抖動
func (rm *RetryManager) addJitter(baseDelay time.Duration, jitterPercent float64) time.Duration {
	if jitterPercent <= 0 {
		return baseDelay
	}
	
	jitterRange := float64(baseDelay) * jitterPercent
	jitter := rand.Float64()*jitterRange*2 - jitterRange // -jitterRange 到 +jitterRange
	
	finalDelay := time.Duration(float64(baseDelay) + jitter)
	if finalDelay < 0 {
		finalDelay = baseDelay / 2 // 最小為基礎延遲的一半
	}
	
	return finalDelay
}

// SimpleCircuitBreaker 簡單斷路器實作
type SimpleCircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failureCount    int
	lastFailureTime time.Time
	config          CircuitBreakerConfig
}

// CircuitBreakerConfig 斷路器配置
type CircuitBreakerConfig struct {
	FailureThreshold int           `json:"failure_threshold"`
	RecoveryTimeout  time.Duration `json:"recovery_timeout"`
	HalfOpenMaxCalls int           `json:"half_open_max_calls"`
}

// NewSimpleCircuitBreaker 建立簡單斷路器
func NewSimpleCircuitBreaker(config CircuitBreakerConfig) *SimpleCircuitBreaker {
	if config.FailureThreshold <= 0 {
		config.FailureThreshold = 5
	}
	if config.RecoveryTimeout <= 0 {
		config.RecoveryTimeout = time.Minute
	}
	if config.HalfOpenMaxCalls <= 0 {
		config.HalfOpenMaxCalls = 3
	}
	
	return &SimpleCircuitBreaker{
		state:  CircuitClosed,
		config: config,
	}
}

// Call 呼叫斷路器保護的操作
func (cb *SimpleCircuitBreaker) Call(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	// 檢查斷路器狀態
	switch cb.state {
	case CircuitOpen:
		// 檢查是否可以嘗試恢復
		if time.Since(cb.lastFailureTime) > cb.config.RecoveryTimeout {
			cb.state = CircuitHalfOpen
		} else {
			return New(ErrCodeSystemResource, "Circuit breaker is open").WithContext(ErrorContext{
				Operation: "circuit_breaker_check",
				Parameters: map[string]interface{}{
					"state": cb.state,
					"failure_count": cb.failureCount,
				},
			})
		}
	case CircuitHalfOpen:
		// 半開狀態下限制呼叫次數
		// 這裡簡化處理，實際應該計數半開狀態下的呼叫
	}
	
	// 執行操作
	err := fn()
	
	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}
	
	return err
}

// IsOpen 檢查斷路器是否開啟
func (cb *SimpleCircuitBreaker) IsOpen() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state == CircuitOpen
}

// Reset 重置斷路器
func (cb *SimpleCircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.state = CircuitClosed
	cb.failureCount = 0
	cb.lastFailureTime = time.Time{}
}

// GetState 取得斷路器狀態
func (cb *SimpleCircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// recordFailure 記錄失敗
func (cb *SimpleCircuitBreaker) recordFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()
	
	if cb.failureCount >= cb.config.FailureThreshold {
		cb.state = CircuitOpen
	}
}

// recordSuccess 記錄成功
func (cb *SimpleCircuitBreaker) recordSuccess() {
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}
	cb.failureCount = 0
}

// contains 檢查字串是否包含子字串（不區分大小寫）
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		   (str == substr || 
		    (len(str) > len(substr) && 
		     indexOf(toLower(str), toLower(substr)) >= 0))
}

// toLower 轉換為小寫
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// indexOf 查找子字串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// CreateRetryPolicyForError 為特定錯誤建立重試策略
func CreateRetryPolicyForError(code ErrorCode) *RetryPolicy {
	switch code {
	case ErrCodeNetworkConnection:
		return &RetryPolicy{
			MaxRetries:    5,
			InitialDelay:  time.Second * 2,
			MaxDelay:      time.Second * 30,
			BackoffFactor: 2.0,
			RetryableErrors: []ErrorCode{code},
		}
	case ErrCodeAPIRateLimit:
		return &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second * 10,
			MaxDelay:      time.Minute * 5,
			BackoffFactor: 3.0,
			RetryableErrors: []ErrorCode{code},
		}
	case ErrCodeTokenCalculation, ErrCodeCostCalculation:
		return &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 10,
			BackoffFactor: 2.0,
			RetryableErrors: []ErrorCode{code},
		}
	case ErrCodeDataAccess:
		return &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 15,
			BackoffFactor: 2.0,
			RetryableErrors: []ErrorCode{code},
		}
	default:
		return &RetryPolicy{
			MaxRetries:    2,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 5,
			BackoffFactor: 2.0,
		}
	}
}