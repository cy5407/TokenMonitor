package errors

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

// TestAppError_Creation 測試 AppError 建立
func TestAppError_Creation(t *testing.T) {
	tests := []struct {
		name    string
		code    ErrorCode
		message string
	}{
		{
			name:    "Token calculation error",
			code:    ErrCodeTokenCalculation,
			message: "Failed to calculate tokens",
		},
		{
			name:    "Cost calculation error", 
			code:    ErrCodeCostCalculation,
			message: "Failed to calculate cost",
		},
		{
			name:    "Data access error",
			code:    ErrCodeDataAccess,
			message: "Failed to access data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, tt.message)
			
			if err.Code != tt.code {
				t.Errorf("Expected code %s, got %s", tt.code, err.Code)
			}
			
			if err.Message != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, err.Message)
			}
			
			if err.Timestamp.IsZero() {
				t.Error("Timestamp should not be zero")
			}
			
			if len(err.Stack) == 0 {
				t.Error("Stack should not be empty")
			}
		})
	}
}

// TestAppError_WithContext 測試添加上下文
func TestAppError_WithContext(t *testing.T) {
	err := New(ErrCodeTokenCalculation, "Test error")
	
	ctx := ErrorContext{
		Operation:  "test_operation",
		Component:  "test_component",
		SessionID:  "test_session",
		Parameters: map[string]interface{}{
			"input": "test input",
			"method": "estimation",
		},
	}
	
	err = err.WithContext(ctx)
	
	if err.Context.Operation != ctx.Operation {
		t.Errorf("Expected operation %s, got %s", ctx.Operation, err.Context.Operation)
	}
	
	if err.Context.Component != ctx.Component {
		t.Errorf("Expected component %s, got %s", ctx.Component, err.Context.Component)
	}
	
	if err.Context.SessionID != ctx.SessionID {
		t.Errorf("Expected session ID %s, got %s", ctx.SessionID, err.Context.SessionID)
	}
}

// TestAppError_WithCause 測試錯誤鏈
func TestAppError_WithCause(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	appErr := New(ErrCodeTokenCalculation, "Token calculation failed").WithCause(originalErr)
	
	if appErr.Cause != originalErr {
		t.Error("Cause should be set correctly")
	}
	
	if appErr.Unwrap() != originalErr {
		t.Error("Unwrap should return the cause")
	}
}

// TestAppError_LocalizedMessage 測試本地化訊息
func TestAppError_LocalizedMessage(t *testing.T) {
	err := New(ErrCodeTokenCalculation, "Token calculation failed")
	
	// 測試中文訊息
	zhMsg := err.GetLocalizedMessage("zh")
	if zhMsg == "" {
		t.Error("Chinese message should not be empty")
	}
	
	// 測試英文訊息
	enMsg := err.GetLocalizedMessage("en")
	if enMsg != err.Message {
		t.Error("English message should match the original message")
	}
}

// TestAppError_IsRetryable 測試可重試性
func TestAppError_IsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		code      ErrorCode
		retryable bool
	}{
		{
			name:      "Token calculation - retryable",
			code:      ErrCodeTokenCalculation,
			retryable: true,
		},
		{
			name:      "Invalid text - not retryable",
			code:      ErrCodeInvalidText,
			retryable: false,
		},
		{
			name:      "Network connection - retryable",
			code:      ErrCodeNetworkConnection,
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.code, "test error")
			
			if err.IsRetryable() != tt.retryable {
				t.Errorf("Expected retryable %v, got %v", tt.retryable, err.IsRetryable())
			}
		})
	}
}

// TestAppError_Sanitize 測試敏感資訊清理
func TestAppError_Sanitize(t *testing.T) {
	err := New(ErrCodeTokenCalculation, "Token calculation failed")
	err = err.WithContext(ErrorContext{
		Operation: "calculate_tokens",
		Parameters: map[string]interface{}{
			"input":    "test input",
			"password": "secret123",
			"api_key":  "sk-1234567890",
			"safe_param": "safe_value",
		},
	})
	
	sanitized := Sanitize(err).(*AppError)
	
	// 檢查敏感參數是否被清理
	if sanitized.Context.Parameters["password"] != "[REDACTED]" {
		t.Error("Password should be redacted")
	}
	
	if sanitized.Context.Parameters["api_key"] != "[REDACTED]" {
		t.Error("API key should be redacted")
	}
	
	// 檢查安全參數是否保留
	if sanitized.Context.Parameters["safe_param"] != "safe_value" {
		t.Error("Safe parameter should not be redacted")
	}
}

// TestErrorHandler_Handle 測試錯誤處理器
func TestErrorHandler_Handle(t *testing.T) {
	handler := NewErrorHandler()
	
	// 建立測試錯誤
	err := New(ErrCodeTokenCalculation, "Test error")
	
	ctx := context.Background()
	handledErr := handler.Handle(ctx, err)
	
	if handledErr == nil {
		t.Error("Handler should return the error")
	}
	
	appErr, ok := handledErr.(*AppError)
	if !ok {
		t.Error("Handler should return AppError")
	}
	
	if appErr.Code != ErrCodeTokenCalculation {
		t.Error("Error code should be preserved")
	}
}

// TestRetryManager_Execute 測試重試管理器
func TestRetryManager_Execute(t *testing.T) {
	manager := NewRetryManager()
	
	// 測試成功情況
	t.Run("Successful execution", func(t *testing.T) {
		callCount := 0
		fn := func() error {
			callCount++
			return nil
		}
		
		policy := &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Millisecond * 10,
			MaxDelay:      time.Second,
			BackoffFactor: 2.0,
		}
		
		ctx := context.Background()
		success, err := manager.Execute(ctx, fn, policy)
		
		if !success {
			t.Error("Execution should be successful")
		}
		
		if err != nil {
			t.Errorf("Error should be nil, got %v", err)
		}
		
		if callCount != 1 {
			t.Errorf("Function should be called once, got %d", callCount)
		}
	})
	
	// 測試重試情況
	t.Run("Retry execution", func(t *testing.T) {
		callCount := 0
		fn := func() error {
			callCount++
			if callCount < 3 {
				return New(ErrCodeNetworkConnection, "Network error")
			}
			return nil
		}
		
		policy := &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Millisecond * 10,
			MaxDelay:      time.Second,
			BackoffFactor: 2.0,
			RetryableErrors: []ErrorCode{ErrCodeNetworkConnection},
		}
		
		ctx := context.Background()
		success, err := manager.Execute(ctx, fn, policy)
		
		if !success {
			t.Error("Execution should eventually succeed")
		}
		
		if err != nil {
			t.Errorf("Error should be nil, got %v", err)
		}
		
		if callCount != 3 {
			t.Errorf("Function should be called 3 times, got %d", callCount)
		}
	})
	
	// 測試不可重試錯誤
	t.Run("Non-retryable error", func(t *testing.T) {
		callCount := 0
		fn := func() error {
			callCount++
			return New(ErrCodeInvalidText, "Invalid text")
		}
		
		policy := &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Millisecond * 10,
			MaxDelay:      time.Second,
			BackoffFactor: 2.0,
		}
		
		ctx := context.Background()
		success, err := manager.Execute(ctx, fn, policy)
		
		if success {
			t.Error("Execution should fail")
		}
		
		if err == nil {
			t.Error("Error should not be nil")
		}
		
		if callCount != 1 {
			t.Errorf("Function should be called once, got %d", callCount)
		}
	})
}

// TestCircuitBreaker 測試斷路器
func TestCircuitBreaker(t *testing.T) {
	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Millisecond * 100,
		HalfOpenMaxCalls: 2,
	}
	
	breaker := NewSimpleCircuitBreaker(config)
	
	// 測試初始狀態
	if breaker.GetState() != CircuitClosed {
		t.Error("Circuit breaker should start in closed state")
	}
	
	// 測試失敗計數
	ctx := context.Background()
	for i := 0; i < config.FailureThreshold; i++ {
		err := breaker.Call(ctx, func() error {
			return fmt.Errorf("test error")
		})
		if err == nil {
			t.Error("Should return error")
		}
	}
	
	// 斷路器應該開啟
	if breaker.GetState() != CircuitOpen {
		t.Error("Circuit breaker should be open after failures")
	}
	
	// 測試斷路器開啟時的行為
	err := breaker.Call(ctx, func() error {
		return nil
	})
	
	if err == nil {
		t.Error("Should return circuit breaker error when open")
	}
	
	// 等待恢復時間
	time.Sleep(config.RecoveryTimeout + time.Millisecond*10)
	
	// 現在應該可以成功呼叫
	err = breaker.Call(ctx, func() error {
		return nil
	})
	
	if err != nil {
		t.Errorf("Should succeed after recovery timeout, got %v", err)
	}
	
	// 斷路器應該回到關閉狀態
	if breaker.GetState() != CircuitClosed {
		t.Error("Circuit breaker should be closed after successful call")
	}
}

// TestErrorMetadata 測試錯誤元數據
func TestErrorMetadata(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected bool
	}{
		{ErrCodeTokenCalculation, true},
		{ErrCodeInvalidText, true},
		{ErrCodeNetworkConnection, true},
		{ErrorCode("UNKNOWN_ERROR"), false},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			metadata, exists := GetErrorMetadata(tt.code)
			
			if exists != tt.expected {
				t.Errorf("Expected metadata exists %v, got %v", tt.expected, exists)
			}
			
			if exists {
				if metadata.Code != tt.code {
					t.Errorf("Expected code %s, got %s", tt.code, metadata.Code)
				}
				
				if metadata.Message == "" {
					t.Error("Message should not be empty")
				}
				
				if metadata.MessageZH == "" {
					t.Error("Chinese message should not be empty")
				}
			}
		})
	}
}

// TestErrorCategorization 測試錯誤分類
func TestErrorCategorization(t *testing.T) {
	err := New(ErrCodeTokenCalculation, "Test error")
	
	if !IsCategory(err, CategoryToken) {
		t.Error("Error should be in token category")
	}
	
	if IsCategory(err, CategoryNetwork) {
		t.Error("Error should not be in network category")
	}
	
	if !IsSeverity(err, SeverityHigh) {
		t.Error("Error should be high severity")
	}
}

// TestErrorJSON 測試錯誤 JSON 序列化
func TestErrorJSON(t *testing.T) {
	err := New(ErrCodeTokenCalculation, "Test error")
	err = err.WithContext(ErrorContext{
		Operation: "test_operation",
		Component: "test_component",
	})
	
	jsonData, jsonErr := err.ToJSON()
	if jsonErr != nil {
		t.Errorf("JSON serialization failed: %v", jsonErr)
	}
	
	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, string(ErrCodeTokenCalculation)) {
		t.Error("JSON should contain error code")
	}
	
	if !strings.Contains(jsonStr, "test_operation") {
		t.Error("JSON should contain operation context")
	}
}

// BenchmarkErrorCreation 錯誤建立效能測試
func BenchmarkErrorCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New(ErrCodeTokenCalculation, "Test error")
	}
}

// BenchmarkErrorWithContext 錯誤上下文效能測試
func BenchmarkErrorWithContext(b *testing.B) {
	ctx := ErrorContext{
		Operation: "benchmark_test",
		Component: "test_component",
		Parameters: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
	}
	
	for i := 0; i < b.N; i++ {
		_ = New(ErrCodeTokenCalculation, "Test error").WithContext(ctx)
	}
}

// BenchmarkErrorSanitize 錯誤清理效能測試
func BenchmarkErrorSanitize(b *testing.B) {
	err := New(ErrCodeTokenCalculation, "Test error")
	err = err.WithContext(ErrorContext{
		Operation: "benchmark_test",
		Parameters: map[string]interface{}{
			"password": "secret123",
			"api_key":  "sk-1234567890",
			"normal":   "normal_value",
		},
	})
	
	for i := 0; i < b.N; i++ {
		_ = Sanitize(err)
	}
}