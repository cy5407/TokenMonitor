package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// AppError 應用程式錯誤結構
type AppError struct {
	Code      ErrorCode     `json:"code"`
	Message   string        `json:"message"`
	MessageZH string        `json:"message_zh"`
	Category  ErrorCategory `json:"category"`
	Severity  ErrorSeverity `json:"severity"`
	Context   ErrorContext  `json:"context"`
	Cause     error         `json:"cause,omitempty"`
	Stack     []StackFrame  `json:"stack,omitempty"`
	Metadata  ErrorMetadata `json:"metadata"`
	Timestamp time.Time     `json:"timestamp"`
}

// ErrorContext 錯誤上下文資訊
type ErrorContext struct {
	Operation   string                 `json:"operation"`
	Component   string                 `json:"component"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	SessionID   string                 `json:"session_id,omitempty"`
	UserID      string                 `json:"user_id,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	Environment string                 `json:"environment,omitempty"`
}

// StackFrame 堆疊框架資訊
type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// Error 實作 error 介面
func (e *AppError) Error() string {
	if e.MessageZH != "" {
		return fmt.Sprintf("[%s] %s (%s)", e.Code, e.MessageZH, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Is 實作錯誤比較
func (e *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return e.Code == appErr.Code
	}
	return false
}

// Unwrap 實作錯誤鏈
func (e *AppError) Unwrap() error {
	return e.Cause
}

// WithCause 添加原因錯誤
func (e *AppError) WithCause(cause error) *AppError {
	e.Cause = cause
	return e
}

// WithContext 添加上下文
func (e *AppError) WithContext(ctx ErrorContext) *AppError {
	e.Context = ctx
	return e
}

// WithParameter 添加參數
func (e *AppError) WithParameter(key string, value interface{}) *AppError {
	if e.Context.Parameters == nil {
		e.Context.Parameters = make(map[string]interface{})
	}
	e.Context.Parameters[key] = value
	return e
}

// WithStack 添加堆疊資訊
func (e *AppError) WithStack() *AppError {
	e.Stack = captureStack(2) // 跳過當前函數和呼叫者
	return e
}

// ToJSON 轉換為 JSON
func (e *AppError) ToJSON() ([]byte, error) {
	return json.MarshalIndent(e, "", "  ")
}

// GetLocalizedMessage 取得本地化訊息
func (e *AppError) GetLocalizedMessage(lang string) string {
	switch strings.ToLower(lang) {
	case "zh", "zh-tw", "zh-cn", "chinese":
		if e.MessageZH != "" {
			return e.MessageZH
		}
	}
	return e.Message
}

// GetSolution 取得解決方案
func (e *AppError) GetSolution(lang string) string {
	switch strings.ToLower(lang) {
	case "zh", "zh-tw", "zh-cn", "chinese":
		if e.Metadata.SolutionZH != "" {
			return e.Metadata.SolutionZH
		}
	}
	return e.Metadata.Solution
}

// IsRetryable 檢查是否可重試
func (e *AppError) IsRetryable() bool {
	return e.Metadata.Retryable
}

// GetRetryPolicy 取得重試策略
func (e *AppError) GetRetryPolicy() *RetryPolicy {
	return e.Metadata.RetryPolicy
}

// captureStack 捕獲堆疊資訊
func captureStack(skip int) []StackFrame {
	var frames []StackFrame
	for i := skip; i < skip+10; i++ { // 限制堆疊深度
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		
		// 過濾內部函數
		funcName := fn.Name()
		if strings.Contains(funcName, "runtime.") ||
		   strings.Contains(funcName, "reflect.") {
			continue
		}
		
		frames = append(frames, StackFrame{
			Function: funcName,
			File:     file,
			Line:     line,
		})
	}
	return frames
}

// New 建立新的應用程式錯誤
func New(code ErrorCode, message string) *AppError {
	metadata, _ := GetErrorMetadata(code)
	
	err := &AppError{
		Code:      code,
		Message:   message,
		MessageZH: metadata.MessageZH,
		Category:  metadata.Category,
		Severity:  metadata.Severity,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
	
	if metadata.MessageZH == "" && message != "" {
		err.MessageZH = message
	}
	
	return err.WithStack()
}

// Newf 建立格式化的應用程式錯誤
func Newf(code ErrorCode, format string, args ...interface{}) *AppError {
	return New(code, fmt.Sprintf(format, args...))
}

// Wrap 包裝現有錯誤
func Wrap(err error, code ErrorCode, message string) *AppError {
	if err == nil {
		return nil
	}
	
	// 如果已經是 AppError，更新它
	if appErr, ok := err.(*AppError); ok {
		if message != "" {
			appErr.Message = message
		}
		return appErr
	}
	
	return New(code, message).WithCause(err)
}

// Wrapf 包裝現有錯誤（格式化）
func Wrapf(err error, code ErrorCode, format string, args ...interface{}) *AppError {
	return Wrap(err, code, fmt.Sprintf(format, args...))
}

// From 從錯誤碼建立錯誤
func From(code ErrorCode) *AppError {
	metadata, exists := GetErrorMetadata(code)
	if !exists {
		return New(code, "Unknown error")
	}
	
	return New(code, metadata.Message)
}

// IsCode 檢查錯誤是否為特定代碼
func IsCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// IsCategory 檢查錯誤是否為特定類別
func IsCategory(err error, category ErrorCategory) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Category == category
	}
	return false
}

// IsSeverity 檢查錯誤是否為特定嚴重級別
func IsSeverity(err error, severity ErrorSeverity) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Severity == severity
	}
	return false
}

// ExtractCode 從錯誤中提取錯誤代碼
func ExtractCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ""
}

// ExtractContext 從錯誤中提取上下文
func ExtractContext(err error) ErrorContext {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Context
	}
	return ErrorContext{}
}

// Sanitize 清理錯誤訊息中的敏感資訊
func Sanitize(err error) error {
	if appErr, ok := err.(*AppError); ok {
		// 建立副本以避免修改原始錯誤
		sanitized := *appErr
		
		// 清理敏感參數
		if sanitized.Context.Parameters != nil {
			params := make(map[string]interface{})
			for k, v := range sanitized.Context.Parameters {
				if isSensitiveKey(k) {
					params[k] = "[REDACTED]"
				} else {
					params[k] = sanitizeValue(v)
				}
			}
			sanitized.Context.Parameters = params
		}
		
		// 清理堆疊中的敏感路徑
		for i, frame := range sanitized.Stack {
			sanitized.Stack[i].File = sanitizePath(frame.File)
		}
		
		return &sanitized
	}
	return err
}

// isSensitiveKey 檢查鍵是否為敏感資訊
func isSensitiveKey(key string) bool {
	sensitiveKeys := []string{
		"password", "token", "key", "secret", "credential",
		"auth", "api_key", "access_token", "private_key",
		"密碼", "金鑰", "密鑰", "令牌", "憑證",
	}
	
	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveKeys {
		if strings.Contains(keyLower, sensitive) {
			return true
		}
	}
	return false
}

// sanitizeValue 清理值中的敏感資訊
func sanitizeValue(value interface{}) interface{} {
	if str, ok := value.(string); ok {
		// 簡單的敏感字串檢測
		if len(str) > 20 && (strings.Contains(str, "password") || 
			strings.Contains(str, "token") || strings.Contains(str, "key")) {
			return "[REDACTED]"
		}
	}
	return value
}

// sanitizePath 清理檔案路徑
func sanitizePath(path string) string {
	// 移除使用者特定路徑
	if strings.Contains(path, "Users") || strings.Contains(path, "用戶") {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if part == "Users" || part == "用戶" {
				if i+1 < len(parts) {
					parts[i+1] = "[USER]"
				}
				break
			}
		}
		return strings.Join(parts, "/")
	}
	return path
}