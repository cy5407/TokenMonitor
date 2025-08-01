package errors

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// ErrorHandler 錯誤處理器介面
type ErrorHandler interface {
	Handle(ctx context.Context, err error) error
	HandleWithRecovery(ctx context.Context, err error, recoveryFn func() error) error
	RegisterListener(listener ErrorListener)
	SetCircuitBreaker(breaker CircuitBreaker)
	SetLogger(logger Logger)
}

// ErrorListener 錯誤監聽器
type ErrorListener interface {
	OnError(ctx context.Context, err *AppError)
	OnRecovery(ctx context.Context, err *AppError, recovered bool)
}

// CircuitBreaker 斷路器介面
type CircuitBreaker interface {
	Call(ctx context.Context, fn func() error) error
	IsOpen() bool
	Reset()
	GetState() CircuitState
}

// CircuitState 斷路器狀態
type CircuitState string

const (
	CircuitClosed   CircuitState = "closed"   // 正常狀態
	CircuitOpen     CircuitState = "open"     // 斷開狀態
	CircuitHalfOpen CircuitState = "half-open" // 半開狀態
)

// Logger 日誌記錄器介面
type Logger interface {
	Error(ctx context.Context, err error, fields map[string]interface{})
	Warn(ctx context.Context, message string, fields map[string]interface{})
	Info(ctx context.Context, message string, fields map[string]interface{})
	Debug(ctx context.Context, message string, fields map[string]interface{})
}

// DefaultErrorHandler 預設錯誤處理器
type DefaultErrorHandler struct {
	listeners      []ErrorListener
	circuitBreaker CircuitBreaker
	logger         Logger
	retryManager   *RetryManager
	mu             sync.RWMutex
}

// NewErrorHandler 建立新的錯誤處理器
func NewErrorHandler() *DefaultErrorHandler {
	return &DefaultErrorHandler{
		listeners:    make([]ErrorListener, 0),
		retryManager: NewRetryManager(),
		logger:       NewDefaultLogger(),
	}
}

// Handle 處理錯誤
func (h *DefaultErrorHandler) Handle(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	// 轉換為 AppError
	var appErr *AppError
	if ae, ok := err.(*AppError); ok {
		appErr = ae
	} else {
		appErr = New(ErrCodeSystemResource, err.Error()).WithCause(err)
	}

	// 清理敏感資訊
	sanitizedErr := Sanitize(appErr).(*AppError)

	// 記錄錯誤
	h.logError(ctx, sanitizedErr)

	// 通知監聽器
	h.notifyListeners(ctx, sanitizedErr)

	// 檢查是否需要斷路器處理
	if h.circuitBreaker != nil && h.shouldUseCircuitBreaker(sanitizedErr) {
		return h.circuitBreaker.Call(ctx, func() error {
			return sanitizedErr
		})
	}

	return sanitizedErr
}

// HandleWithRecovery 處理錯誤並嘗試恢復
func (h *DefaultErrorHandler) HandleWithRecovery(ctx context.Context, err error, recoveryFn func() error) error {
	if err == nil {
		return nil
	}

	// 先處理錯誤
	handledErr := h.Handle(ctx, err)
	
	appErr, ok := handledErr.(*AppError)
	if !ok {
		return handledErr
	}

	// 檢查是否可恢復
	if !appErr.IsRetryable() || recoveryFn == nil {
		return handledErr
	}

	// 嘗試恢復
	retryPolicy := appErr.GetRetryPolicy()
	if retryPolicy == nil {
		retryPolicy = &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 10,
			BackoffFactor: 2.0,
		}
	}

	recovered, finalErr := h.retryManager.Execute(ctx, recoveryFn, retryPolicy)
	
	// 通知恢復結果
	h.notifyRecovery(ctx, appErr, recovered)

	if recovered {
		h.logger.Info(ctx, "錯誤恢復成功", map[string]interface{}{
			"error_code": appErr.Code,
			"operation":  appErr.Context.Operation,
		})
		return nil
	}

	return finalErr
}

// RegisterListener 註冊錯誤監聽器
func (h *DefaultErrorHandler) RegisterListener(listener ErrorListener) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.listeners = append(h.listeners, listener)
}

// SetCircuitBreaker 設定斷路器
func (h *DefaultErrorHandler) SetCircuitBreaker(breaker CircuitBreaker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.circuitBreaker = breaker
}

// SetLogger 設定日誌記錄器
func (h *DefaultErrorHandler) SetLogger(logger Logger) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.logger = logger
}

// logError 記錄錯誤
func (h *DefaultErrorHandler) logError(ctx context.Context, err *AppError) {
	fields := map[string]interface{}{
		"error_code": err.Code,
		"category":   err.Category,
		"severity":   err.Severity,
		"operation":  err.Context.Operation,
		"component":  err.Context.Component,
		"timestamp":  err.Timestamp,
	}

	if err.Context.SessionID != "" {
		fields["session_id"] = err.Context.SessionID
	}

	if err.Context.RequestID != "" {
		fields["request_id"] = err.Context.RequestID
	}

	switch err.Severity {
	case SeverityCritical, SeverityHigh:
		h.logger.Error(ctx, err, fields)
	case SeverityMedium:
		h.logger.Warn(ctx, err.Error(), fields)
	default:
		h.logger.Info(ctx, err.Error(), fields)
	}
}

// notifyListeners 通知監聽器
func (h *DefaultErrorHandler) notifyListeners(ctx context.Context, err *AppError) {
	h.mu.RLock()
	listeners := make([]ErrorListener, len(h.listeners))
	copy(listeners, h.listeners)
	h.mu.RUnlock()

	for _, listener := range listeners {
		go func(l ErrorListener) {
			defer func() {
				if r := recover(); r != nil {
					h.logger.Error(ctx, fmt.Errorf("error listener panic: %v", r), nil)
				}
			}()
			l.OnError(ctx, err)
		}(listener)
	}
}

// notifyRecovery 通知恢復結果
func (h *DefaultErrorHandler) notifyRecovery(ctx context.Context, err *AppError, recovered bool) {
	h.mu.RLock()
	listeners := make([]ErrorListener, len(h.listeners))
	copy(listeners, h.listeners)
	h.mu.RUnlock()

	for _, listener := range listeners {
		go func(l ErrorListener) {
			defer func() {
				if r := recover(); r != nil {
					h.logger.Error(ctx, fmt.Errorf("recovery listener panic: %v", r), nil)
				}
			}()
			l.OnRecovery(ctx, err, recovered)
		}(listener)
	}
}

// shouldUseCircuitBreaker 檢查是否應使用斷路器
func (h *DefaultErrorHandler) shouldUseCircuitBreaker(err *AppError) bool {
	// 對網路相關和系統資源錯誤使用斷路器
	return err.Category == CategoryNetwork || 
		   err.Category == CategorySystem ||
		   err.Severity == SeverityCritical
}

// DefaultLogger 預設日誌記錄器
type DefaultLogger struct {
	logger *log.Logger
}

// NewDefaultLogger 建立預設日誌記錄器
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stderr, "[TokenMonitor] ", log.LstdFlags|log.Lshortfile),
	}
}

// Error 記錄錯誤日誌
func (l *DefaultLogger) Error(ctx context.Context, err error, fields map[string]interface{}) {
	l.logWithLevel("ERROR", err.Error(), fields)
}

// Warn 記錄警告日誌
func (l *DefaultLogger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLevel("WARN", message, fields)
}

// Info 記錄資訊日誌
func (l *DefaultLogger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLevel("INFO", message, fields)
}

// Debug 記錄除錯日誌
func (l *DefaultLogger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	l.logWithLevel("DEBUG", message, fields)
}

// logWithLevel 按級別記錄日誌
func (l *DefaultLogger) logWithLevel(level string, message string, fields map[string]interface{}) {
	logMsg := fmt.Sprintf("[%s] %s", level, message)
	
	if fields != nil && len(fields) > 0 {
		logMsg += " | Fields: "
		for k, v := range fields {
			logMsg += fmt.Sprintf("%s=%v ", k, v)
		}
	}
	
	l.logger.Println(logMsg)
}

// DefaultErrorListener 預設錯誤監聽器
type DefaultErrorListener struct {
	logger Logger
}

// NewDefaultErrorListener 建立預設錯誤監聽器
func NewDefaultErrorListener(logger Logger) *DefaultErrorListener {
	return &DefaultErrorListener{
		logger: logger,
	}
}

// OnError 錯誤發生時的處理
func (l *DefaultErrorListener) OnError(ctx context.Context, err *AppError) {
	// 根據錯誤嚴重級別執行不同的處理
	switch err.Severity {
	case SeverityCritical:
		l.handleCriticalError(ctx, err)
	case SeverityHigh:
		l.handleHighError(ctx, err)
	case SeverityMedium:
		l.handleMediumError(ctx, err)
	default:
		l.handleLowError(ctx, err)
	}
}

// OnRecovery 錯誤恢復時的處理
func (l *DefaultErrorListener) OnRecovery(ctx context.Context, err *AppError, recovered bool) {
	status := "失敗"
	if recovered {
		status = "成功"
	}
	
	l.logger.Info(ctx, fmt.Sprintf("錯誤恢復%s", status), map[string]interface{}{
		"error_code": err.Code,
		"recovered":  recovered,
		"operation":  err.Context.Operation,
	})
}

// handleCriticalError 處理關鍵錯誤
func (l *DefaultErrorListener) handleCriticalError(ctx context.Context, err *AppError) {
	l.logger.Error(ctx, err, map[string]interface{}{
		"action": "系統需要立即關注",
		"recommendation": "檢查系統資源和配置",
	})
}

// handleHighError 處理高級錯誤
func (l *DefaultErrorListener) handleHighError(ctx context.Context, err *AppError) {
	l.logger.Warn(ctx, "高級錯誤需要關注", map[string]interface{}{
		"error_code": err.Code,
		"solution": err.GetSolution("zh"),
	})
}

// handleMediumError 處理中級錯誤
func (l *DefaultErrorListener) handleMediumError(ctx context.Context, err *AppError) {
	l.logger.Info(ctx, "中級錯誤", map[string]interface{}{
		"error_code": err.Code,
		"message": err.GetLocalizedMessage("zh"),
	})
}

// handleLowError 處理低級錯誤
func (l *DefaultErrorListener) handleLowError(ctx context.Context, err *AppError) {
	l.logger.Debug(ctx, "低級錯誤", map[string]interface{}{
		"error_code": err.Code,
		"context": err.Context.Operation,
	})
}