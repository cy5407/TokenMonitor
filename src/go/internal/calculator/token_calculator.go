package calculator

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unicode"

	"token-monitor/internal/errors"
	"token-monitor/internal/interfaces"
	"token-monitor/internal/types"

	"github.com/pkoukk/tiktoken-go"
)

// TokenCalculatorImpl Token 計算器實作
type TokenCalculatorImpl struct {
	cache           map[string]int
	cacheMutex      sync.RWMutex
	maxCacheSize    int
	tiktokenEnabled bool
	tiktokenEncoder *tiktoken.Tiktoken
	errorHandler    errors.ErrorHandler

	// 估算演算法參數
	englishCharsPerToken float64
	chineseCharsPerToken float64
}

// NewTokenCalculator 建立新的 Token 計算器
func NewTokenCalculator(maxCacheSize int) interfaces.TokenCalculator {
	calc := &TokenCalculatorImpl{
		cache:                make(map[string]int),
		maxCacheSize:         maxCacheSize,
		tiktokenEnabled:      false,
		errorHandler:         errors.NewErrorHandler(),
		englishCharsPerToken: 4.0, // 英文約 4 字符 = 1 token
		chineseCharsPerToken: 1.5, // 中文約 1.5 字符 = 1 token
	}

	// 嘗試初始化 tiktoken
	calc.initTiktoken()

	return calc
}

// CalculateTokens 計算文本的 Token 數量
func (tc *TokenCalculatorImpl) CalculateTokens(text string, method string) (int, error) {
	ctx := context.Background()
	
	if text == "" {
		return 0, nil
	}

	// 驗證文本
	if err := tc.ValidateText(text); err != nil {
		appErr := errors.New(errors.ErrCodeInvalidText, "文本驗證失敗").WithCause(err)
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "validate_text",
			Component:  "token_calculator",
			Parameters: map[string]interface{}{
				"text_length": len(text),
				"method":      method,
			},
		})
		return 0, tc.errorHandler.Handle(ctx, appErr)
	}

	// 檢查快取
	if tokens, found := tc.getCachedTokens(text); found {
		return tokens, nil
	}

	var tokens int
	var err error

	switch method {
	case "tiktoken":
		if tc.tiktokenEnabled {
			tokens, err = tc.calculateWithTiktoken(text)
		} else {
			// tiktoken 不可用，記錄警告並回退到估算方法
			warnErr := errors.New(errors.ErrCodeTiktokenUnavailable, "Tiktoken 不可用，使用估算方法")
			warnErr = warnErr.WithContext(errors.ErrorContext{
				Operation: "fallback_to_estimation",
				Component: "token_calculator",
			})
			tc.errorHandler.Handle(ctx, warnErr)
			tokens, err = tc.calculateWithEstimation(text)
		}
	case "estimation":
		tokens, err = tc.calculateWithEstimation(text)
	default:
		// 預設使用最佳可用方法
		if tc.tiktokenEnabled {
			tokens, err = tc.calculateWithTiktoken(text)
		} else {
			tokens, err = tc.calculateWithEstimation(text)
		}
	}

	if err != nil {
		appErr := errors.Wrap(err, errors.ErrCodeTokenCalculation, "Token 計算失敗")
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "calculate_tokens",
			Component:  "token_calculator",
			Parameters: map[string]interface{}{
				"text_length": len(text),
				"method":      method,
			},
		})
		return 0, tc.errorHandler.Handle(ctx, appErr)
	}

	// 儲存到快取
	tc.setCachedTokens(text, tokens)

	return tokens, nil
}

// calculateWithEstimation 使用估算演算法計算 Token
func (tc *TokenCalculatorImpl) calculateWithEstimation(text string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 檢查上下文是否已取消
	select {
	case <-ctx.Done():
		return 0, errors.New(errors.ErrCodeCalculationTimeout, "計算超時")
	default:
	}
	
	// 分離英文和中文字符
	englishChars := 0
	chineseChars := 0

	for _, r := range text {
		if r <= unicode.MaxASCII {
			// ASCII 字符（包括英文、數字、符號）
			englishChars++
		} else if unicode.Is(unicode.Han, r) {
			// 中文字符
			chineseChars++
		} else {
			// 其他 Unicode 字符，按英文處理
			englishChars++
		}
	}

	// 計算 Token 數量
	englishTokens := float64(englishChars) / tc.englishCharsPerToken
	chineseTokens := float64(chineseChars) / tc.chineseCharsPerToken

	totalTokens := int(englishTokens + chineseTokens)

	// 至少 1 個 token（如果有內容的話）
	if totalTokens == 0 && len(text) > 0 {
		totalTokens = 1
	}

	return totalTokens, nil
}

// initTiktoken 初始化 tiktoken 編碼器
func (tc *TokenCalculatorImpl) initTiktoken() {
	ctx := context.Background()
	
	// 嘗試初始化 tiktoken 編碼器 (使用 cl100k_base，適用於 GPT-3.5/GPT-4)
	encoder, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		warnErr := errors.Wrap(err, errors.ErrCodeTiktokenUnavailable, "Tiktoken 初始化失敗")
		warnErr = warnErr.WithContext(errors.ErrorContext{
			Operation: "init_tiktoken",
			Component: "token_calculator",
			Parameters: map[string]interface{}{
				"encoding": "cl100k_base",
			},
		})
		tc.errorHandler.Handle(ctx, warnErr)
		tc.tiktokenEnabled = false
		return
	}

	tc.tiktokenEncoder = encoder
	tc.tiktokenEnabled = true
	fmt.Println("✅ Tiktoken initialized successfully")
}

// calculateWithTiktoken 使用 tiktoken 計算 Token
func (tc *TokenCalculatorImpl) calculateWithTiktoken(text string) (int, error) {
	if !tc.tiktokenEnabled || tc.tiktokenEncoder == nil {
		return tc.calculateWithEstimation(text)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// 檢查上下文是否已取消
	select {
	case <-ctx.Done():
		return 0, errors.New(errors.ErrCodeCalculationTimeout, "Tiktoken 計算超時")
	default:
	}

	// 使用 tiktoken 進行精確計算
	defer func() {
		if r := recover(); r != nil {
			panic(errors.New(errors.ErrCodeTokenCalculation, fmt.Sprintf("Tiktoken 計算發生恐慌: %v", r)))
		}
	}()
	
	tokens := tc.tiktokenEncoder.Encode(text, nil, nil)
	return len(tokens), nil
}

// AnalyzeTokenDistribution 分析 Token 分佈
func (tc *TokenCalculatorImpl) AnalyzeTokenDistribution(text string) (*types.TokenDistribution, error) {
	if text == "" {
		return &types.TokenDistribution{
			EnglishTokens: 0,
			ChineseTokens: 0,
			TotalTokens:   0,
			Method:        "estimation",
		}, nil
	}

	// 分析字符分佈
	englishChars := 0
	chineseChars := 0

	for _, r := range text {
		if r <= unicode.MaxASCII {
			englishChars++
		} else if unicode.Is(unicode.Han, r) {
			chineseChars++
		} else {
			englishChars++
		}
	}

	// 計算各類型的 Token 數量
	englishTokens := int(float64(englishChars) / tc.englishCharsPerToken)
	chineseTokens := int(float64(chineseChars) / tc.chineseCharsPerToken)
	totalTokens := englishTokens + chineseTokens

	if totalTokens == 0 && len(text) > 0 {
		totalTokens = 1
		if chineseChars > englishChars {
			chineseTokens = 1
		} else {
			englishTokens = 1
		}
	}

	method := "estimation"
	if tc.tiktokenEnabled {
		method = "tiktoken"
		// 如果使用 tiktoken，重新計算總 Token 數
		if actualTotalTokens, err := tc.calculateWithTiktoken(text); err == nil {
			// tiktoken 不區分中英文，所以我們需要估算分佈
			if len(text) > 0 {
				ratio := float64(chineseChars) / float64(len(text))
				chineseTokens = int(float64(actualTotalTokens) * ratio)
				englishTokens = actualTotalTokens - chineseTokens
				totalTokens = actualTotalTokens
			}
		}
	}

	return &types.TokenDistribution{
		EnglishTokens: englishTokens,
		ChineseTokens: chineseTokens,
		TotalTokens:   totalTokens,
		Method:        method,
	}, nil
}

// IsTiktokenAvailable 檢查 tiktoken 是否可用
func (tc *TokenCalculatorImpl) IsTiktokenAvailable() bool {
	return tc.tiktokenEnabled
}

// ClearCache 清除快取
func (tc *TokenCalculatorImpl) ClearCache() {
	tc.cacheMutex.Lock()
	defer tc.cacheMutex.Unlock()

	tc.cache = make(map[string]int)
}

// GetSupportedMethods 取得支援的計算方法
func (tc *TokenCalculatorImpl) GetSupportedMethods() []string {
	methods := []string{"estimation"}
	if tc.tiktokenEnabled {
		methods = append(methods, "tiktoken")
	}
	return methods
}

// getCachedTokens 從快取取得 Token 數量
func (tc *TokenCalculatorImpl) getCachedTokens(text string) (int, bool) {
	tc.cacheMutex.RLock()
	defer tc.cacheMutex.RUnlock()

	tokens, found := tc.cache[text]
	return tokens, found
}

// setCachedTokens 設定快取的 Token 數量
func (tc *TokenCalculatorImpl) setCachedTokens(text string, tokens int) {
	tc.cacheMutex.Lock()
	defer tc.cacheMutex.Unlock()

	// 如果快取已滿，清除一些舊的項目
	if len(tc.cache) >= tc.maxCacheSize {
		// 簡單的清除策略：清除一半的快取
		newCache := make(map[string]int)
		count := 0
		for k, v := range tc.cache {
			if count < tc.maxCacheSize/2 {
				newCache[k] = v
				count++
			}
		}
		tc.cache = newCache
	}

	tc.cache[text] = tokens
}

// SetEstimationParameters 設定估算演算法參數
func (tc *TokenCalculatorImpl) SetEstimationParameters(englishCharsPerToken, chineseCharsPerToken float64) {
	tc.englishCharsPerToken = englishCharsPerToken
	tc.chineseCharsPerToken = chineseCharsPerToken
}

// GetCacheStats 取得快取統計資訊
func (tc *TokenCalculatorImpl) GetCacheStats() map[string]interface{} {
	tc.cacheMutex.RLock()
	defer tc.cacheMutex.RUnlock()

	return map[string]interface{}{
		"cache_size":       len(tc.cache),
		"max_cache_size":   tc.maxCacheSize,
		"cache_usage":      float64(len(tc.cache)) / float64(tc.maxCacheSize),
		"tiktoken_enabled": tc.tiktokenEnabled,
	}
}

// CalculateTokensForMultipleTexts 批次計算多個文本的 Token
func (tc *TokenCalculatorImpl) CalculateTokensForMultipleTexts(texts []string, method string) ([]int, error) {
	results := make([]int, len(texts))

	for i, text := range texts {
		tokens, err := tc.CalculateTokens(text, method)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tokens for text %d: %w", i, err)
		}
		results[i] = tokens
	}

	return results, nil
}

// ValidateText 驗證文本是否適合 Token 計算
func (tc *TokenCalculatorImpl) ValidateText(text string) error {
	if len(text) > 1000000 { // 1MB 限制
		return errors.Newf(errors.ErrCodeInvalidText, "文本過長: %d 字符（最大限制: 1,000,000）", len(text))
	}

	// 檢查是否包含過多的控制字符
	controlChars := 0
	for _, r := range text {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			controlChars++
		}
	}

	if controlChars > len(text)/10 { // 如果控制字符超過 10%
		return errors.Newf(errors.ErrCodeInvalidText, "文本包含過多控制字符: %d（超過總長度的10%%）", controlChars)
	}

	return nil
}

// GetTiktokenInfo 取得 tiktoken 資訊
func (tc *TokenCalculatorImpl) GetTiktokenInfo() map[string]interface{} {
	info := map[string]interface{}{
		"enabled": tc.tiktokenEnabled,
	}

	if tc.tiktokenEnabled && tc.tiktokenEncoder != nil {
		info["encoding"] = "cl100k_base"
		info["model_compatibility"] = []string{"gpt-3.5-turbo", "gpt-4", "text-embedding-ada-002"}
	}

	return info
}

// EnableTiktoken 手動啟用 tiktoken（如果初始化失敗）
func (tc *TokenCalculatorImpl) EnableTiktoken() error {
	if tc.tiktokenEnabled {
		return nil // 已經啟用
	}

	encoder, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return fmt.Errorf("failed to initialize tiktoken: %w", err)
	}

	tc.tiktokenEncoder = encoder
	tc.tiktokenEnabled = true
	return nil
}

// DisableTiktoken 停用 tiktoken
func (tc *TokenCalculatorImpl) DisableTiktoken() {
	tc.tiktokenEnabled = false
	tc.tiktokenEncoder = nil
}

// CompareCalculationMethods 比較不同計算方法的結果
func (tc *TokenCalculatorImpl) CompareCalculationMethods(text string) map[string]interface{} {
	result := map[string]interface{}{
		"text_length": len(text),
	}

	// 估算方法
	if estimationTokens, err := tc.calculateWithEstimation(text); err == nil {
		result["estimation"] = map[string]interface{}{
			"tokens": estimationTokens,
			"method": "estimation",
		}
	}

	// tiktoken 方法
	if tc.tiktokenEnabled {
		if tiktokenTokens, err := tc.calculateWithTiktoken(text); err == nil {
			result["tiktoken"] = map[string]interface{}{
				"tokens": tiktokenTokens,
				"method": "tiktoken",
			}

			// 計算差異
			if estimationData, ok := result["estimation"].(map[string]interface{}); ok {
				if estimationTokens, ok := estimationData["tokens"].(int); ok {
					difference := tiktokenTokens - estimationTokens
					accuracy := 100.0 - (float64(abs(difference))/float64(tiktokenTokens))*100.0
					result["comparison"] = map[string]interface{}{
						"difference":       difference,
						"accuracy_percent": accuracy,
						"preferred_method": "tiktoken",
					}
				}
			}
		}
	}

	return result
}

// abs 計算絕對值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
