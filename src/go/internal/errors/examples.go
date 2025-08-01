package errors

import (
	"context"
	"fmt"
	"time"
)

// 使用範例文件 - 展示如何使用錯誤處理機制

// ExampleBasicErrorHandling 基本錯誤處理範例
func ExampleBasicErrorHandling() {
	// 建立錯誤處理器
	handler := NewErrorHandler()
	
	// 建立一個應用錯誤
	err := New(ErrCodeTokenCalculation, "Token 計算失敗")
	err = err.WithContext(ErrorContext{
		Operation:  "calculate_tokens",
		Component:  "token_calculator",
		Parameters: map[string]interface{}{
			"text_length": 1000,
			"method":      "tiktoken",
		},
	})
	
	// 處理錯誤
	ctx := context.Background()
	handledErr := handler.Handle(ctx, err)
	
	fmt.Printf("處理後的錯誤: %v\n", handledErr)
}

// ExampleRetryMechanism 重試機制範例
func ExampleRetryMechanism() {
	_ = NewErrorHandler()
	retryManager := NewRetryManager()
	
	// 建立一個會失敗幾次然後成功的函數
	attemptCount := 0
	fn := func() error {
		attemptCount++
		if attemptCount < 3 {
			return New(ErrCodeNetworkConnection, "網路連接失敗")
		}
		return nil // 第三次嘗試成功
	}
	
	// 設定重試策略
	policy := &RetryPolicy{
		MaxRetries:    5,
		InitialDelay:  time.Millisecond * 100,
		MaxDelay:      time.Second * 5,
		BackoffFactor: 2.0,
		RetryableErrors: []ErrorCode{ErrCodeNetworkConnection},
	}
	
	// 執行帶重試的操作
	ctx := context.Background()
	success, err := retryManager.Execute(ctx, fn, policy)
	
	if success {
		fmt.Printf("操作成功，嘗試次數: %d\n", attemptCount)
	} else {
		fmt.Printf("操作失敗: %v\n", err)
	}
}

// ExampleCircuitBreaker 斷路器範例
func ExampleCircuitBreaker() {
	// 建立斷路器配置
	config := CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second * 5,
		HalfOpenMaxCalls: 2,
	}
	
	breaker := NewSimpleCircuitBreaker(config)
	
	// 模擬多次失敗
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		err := breaker.Call(ctx, func() error {
			if i < 3 {
				return fmt.Errorf("模擬失敗 %d", i+1)
			}
			return nil
		})
		
		fmt.Printf("嘗試 %d: 狀態=%s, 錯誤=%v\n", i+1, breaker.GetState(), err)
		
		if breaker.IsOpen() {
			fmt.Println("斷路器已開啟，停止嘗試")
			break
		}
	}
}

// ExampleErrorSanitization 錯誤清理範例
func ExampleErrorSanitization() {
	// 建立包含敏感資訊的錯誤
	err := New(ErrCodeTokenCalculation, "計算失敗")
	err = err.WithContext(ErrorContext{
		Operation: "api_call",
		Component: "calculator",
		Parameters: map[string]interface{}{
			"api_key":  "sk-1234567890abcdef",
			"password": "secret123",
			"username": "user@example.com",
			"data":     "safe data",
		},
	})
	
	fmt.Println("原始錯誤:")
	if jsonData, jsonErr := err.ToJSON(); jsonErr == nil {
		fmt.Println(string(jsonData))
	}
	
	// 清理敏感資訊
	sanitized := Sanitize(err)
	
	fmt.Println("\n清理後的錯誤:")
	if sanitizedErr, ok := sanitized.(*AppError); ok {
		if jsonData, jsonErr := sanitizedErr.ToJSON(); jsonErr == nil {
			fmt.Println(string(jsonData))
		}
	}
}

// ExampleErrorListener 錯誤監聽器範例
func ExampleErrorListener() {
	handler := NewErrorHandler()
	logger := NewDefaultLogger()
	listener := NewDefaultErrorListener(logger)
	
	// 註冊監聽器
	handler.RegisterListener(listener)
	
	// 建立不同嚴重級別的錯誤
	errors := []*AppError{
		New(ErrCodeInvalidText, "低級錯誤"),
		New(ErrCodeTokenCalculation, "中級錯誤"),
		New(ErrCodeSystemResource, "關鍵錯誤"),
	}
	
	ctx := context.Background()
	for _, err := range errors {
		fmt.Printf("處理錯誤: %s (嚴重級別: %s)\n", err.Code, err.Severity)
		handler.Handle(ctx, err)
	}
}

// ExampleLocalizedMessages 本地化訊息範例
func ExampleLocalizedMessages() {
	err := New(ErrCodeTokenCalculation, "Token calculation failed")
	
	// 取得不同語言的訊息
	enMessage := err.GetLocalizedMessage("en")
	zhMessage := err.GetLocalizedMessage("zh")
	solution := err.GetSolution("zh")
	
	fmt.Printf("英文訊息: %s\n", enMessage)
	fmt.Printf("中文訊息: %s\n", zhMessage)
	fmt.Printf("解決方案: %s\n", solution)
}

// ExampleErrorChaining 錯誤鏈範例
func ExampleErrorChaining() {
	// 建立原始錯誤
	originalErr := fmt.Errorf("原始系統錯誤")
	
	// 包裝為應用錯誤
	appErr := Wrap(originalErr, ErrCodeTokenCalculation, "Token 計算過程中發生錯誤")
	appErr = appErr.WithContext(ErrorContext{
		Operation: "token_calculation",
		Component: "calculator",
	})
	
	// 再次包裝
	wrapperErr := Wrap(appErr, ErrCodeSystemResource, "系統資源不足")
	
	fmt.Printf("最終錯誤: %s\n", wrapperErr.Error())
	fmt.Printf("原因錯誤: %s\n", wrapperErr.Unwrap().Error())
	
	// 檢查錯誤類型
	if IsCode(wrapperErr.Unwrap(), ErrCodeTokenCalculation) {
		fmt.Println("包含 Token 計算錯誤")
	}
}

// ExampleErrorRecovery 錯誤恢復範例
func ExampleErrorRecovery() {
	handler := NewErrorHandler()
	
	// 模擬一個會失敗的操作
	failingOperation := func() error {
		return New(ErrCodeNetworkConnection, "網路連接失敗")
	}
	
	// 恢復函數
	recoveryFn := func() error {
		fmt.Println("嘗試恢復...")
		// 模擬恢復成功
		return nil
	}
	
	ctx := context.Background()
	err := failingOperation()
	
	// 嘗試恢復
	finalErr := handler.HandleWithRecovery(ctx, err, recoveryFn)
	
	if finalErr == nil {
		fmt.Println("錯誤已成功恢復")
	} else {
		fmt.Printf("恢復失敗: %v\n", finalErr)
	}
}

// SecurityValidationExample 安全驗證範例
func SecurityValidationExample() {
	fmt.Println("=== 錯誤處理安全性驗證 ===")
	
	// 驗證敏感資訊過濾
	fmt.Println("\n1. 敏感資訊過濾測試:")
	testSensitiveKeys := []string{
		"password", "token", "key", "secret", "credential", 
		"auth", "api_key", "access_token", "private_key",
		"密碼", "金鑰", "密鑰", "令牌", "憑證",
	}
	
	for _, key := range testSensitiveKeys {
		if isSensitiveKey(key) {
			fmt.Printf("✓ '%s' 被正確識別為敏感鍵\n", key)
		} else {
			fmt.Printf("✗ '%s' 未被識別為敏感鍵\n", key)
		}
	}
	
	// 驗證路徑清理
	fmt.Println("\n2. 路徑清理測試:")
	testPaths := []string{
		"/Users/john/Documents/secret.txt",
		"/home/user/project/config.yaml",
		"C:\\Users\\admin\\Desktop\\file.txt",
		"/tmp/safe_file.log",
	}
	
	for _, path := range testPaths {
		sanitized := sanitizePath(path)
		fmt.Printf("原始: %s -> 清理後: %s\n", path, sanitized)
	}
	
	// 驗證錯誤等級分類
	fmt.Println("\n3. 錯誤等級測試:")
	testCodes := []ErrorCode{
		ErrCodeInvalidText,      // Low
		ErrCodeTokenCalculation, // High  
		ErrCodeSystemResource,   // Critical
		ErrCodeNetworkConnection, // Medium
	}
	
	for _, code := range testCodes {
		if metadata, exists := GetErrorMetadata(code); exists {
			fmt.Printf("錯誤 %s: 等級=%s, 可重試=%t\n", 
				code, metadata.Severity, metadata.Retryable)
		}
	}
}