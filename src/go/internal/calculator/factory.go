package calculator

import (
	"token-monitor/internal/interfaces"

	"github.com/spf13/viper"
)

// NewTokenCalculatorFromConfig 從配置建立 Token 計算器
func NewTokenCalculatorFromConfig() interfaces.TokenCalculator {
	// 從配置讀取參數
	cacheSize := viper.GetInt("token_calculation.cache_size")
	if cacheSize <= 0 {
		cacheSize = 1000 // 預設值
	}

	englishCharsPerToken := viper.GetFloat64("token_calculation.estimation.english_chars_per_token")
	if englishCharsPerToken <= 0 {
		englishCharsPerToken = 4.0 // 預設值
	}

	chineseCharsPerToken := viper.GetFloat64("token_calculation.estimation.chinese_chars_per_token")
	if chineseCharsPerToken <= 0 {
		chineseCharsPerToken = 1.5 // 預設值
	}

	// 建立計算器
	calculator := NewTokenCalculator(cacheSize).(*TokenCalculatorImpl)

	// 設定估算參數
	calculator.SetEstimationParameters(englishCharsPerToken, chineseCharsPerToken)

	return calculator
}

// GetPreferredMethod 取得偏好的計算方法
func GetPreferredMethod() string {
	preferredMethod := viper.GetString("token_calculation.preferred_method")
	if preferredMethod == "" {
		preferredMethod = "estimation" // 預設值
	}
	return preferredMethod
}

// GetFallbackMethod 取得備用的計算方法
func GetFallbackMethod() string {
	fallbackMethod := viper.GetString("token_calculation.fallback_method")
	if fallbackMethod == "" {
		fallbackMethod = "estimation" // 預設值
	}
	return fallbackMethod
}

// IsCacheEnabled 檢查是否啟用快取
func IsCacheEnabled() bool {
	return viper.GetBool("token_calculation.cache_enabled")
}
