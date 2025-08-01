package calculator

import (
	"fmt"
	"testing"
)

func TestTokenCalculatorImpl_CalculateTokens(t *testing.T) {
	calculator := NewTokenCalculator(100)

	testCases := []struct {
		name        string
		text        string
		method      string
		expectedMin int
		expectedMax int
		expectError bool
	}{
		{
			name:        "空文本",
			text:        "",
			method:      "estimation",
			expectedMin: 0,
			expectedMax: 0,
			expectError: false,
		},
		{
			name:        "簡單英文",
			text:        "Hello world",
			method:      "estimation",
			expectedMin: 2,
			expectedMax: 4,
			expectError: false,
		},
		{
			name:        "簡單中文",
			text:        "你好世界",
			method:      "estimation",
			expectedMin: 2,
			expectedMax: 4,
			expectError: false,
		},
		{
			name:        "中英混合",
			text:        "Hello 世界",
			method:      "estimation",
			expectedMin: 2,
			expectedMax: 5,
			expectError: false,
		},
		{
			name:        "程式碼片段",
			text:        "function hello() { return 'world'; }",
			method:      "estimation",
			expectedMin: 8,
			expectedMax: 12,
			expectError: false,
		},
		{
			name:        "長文本",
			text:        "這是一個很長的文本，包含中文和English混合內容，用來測試Token計算的準確性。This is a long text with mixed Chinese and English content for testing token calculation accuracy.",
			method:      "estimation",
			expectedMin: 30,
			expectedMax: 50,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := calculator.CalculateTokens(tc.text, tc.method)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tokens < tc.expectedMin || tokens > tc.expectedMax {
				t.Errorf("Token count %d not in expected range [%d, %d]", tokens, tc.expectedMin, tc.expectedMax)
			}
		})
	}
}

func TestTokenCalculatorImpl_AnalyzeTokenDistribution(t *testing.T) {
	calculator := NewTokenCalculator(100)

	testCases := []struct {
		name string
		text string
	}{
		{
			name: "純英文",
			text: "Hello world",
		},
		{
			name: "純中文",
			text: "你好世界",
		},
		{
			name: "中英混合",
			text: "Hello 世界",
		},
		{
			name: "空文本",
			text: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distribution, err := calculator.AnalyzeTokenDistribution(tc.text)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// 基本檢查：總 Token 數應該等於英文和中文 Token 的總和
			expectedTotal := distribution.EnglishTokens + distribution.ChineseTokens
			if distribution.TotalTokens != expectedTotal {
				t.Errorf("Total tokens %d doesn't match sum of English %d and Chinese %d",
					distribution.TotalTokens, distribution.EnglishTokens, distribution.ChineseTokens)
			}

			// 檢查空文本
			if tc.text == "" {
				if distribution.TotalTokens != 0 {
					t.Errorf("Expected 0 tokens for empty text, got %d", distribution.TotalTokens)
				}
			} else {
				// 非空文本應該有正數 Token
				if distribution.TotalTokens <= 0 {
					t.Errorf("Expected positive tokens for non-empty text, got %d", distribution.TotalTokens)
				}
			}

			// 檢查計算方法
			if distribution.Method != "tiktoken" && distribution.Method != "estimation" {
				t.Errorf("Unexpected calculation method: %s", distribution.Method)
			}
		})
	}
}

func TestTokenCalculatorImpl_Cache(t *testing.T) {
	calculator := NewTokenCalculator(2) // 小的快取大小用於測試

	text1 := "Hello world"
	text2 := "你好世界"
	text3 := "Another text"

	// 第一次計算
	tokens1, err := calculator.CalculateTokens(text1, "estimation")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// 第二次計算相同文本（應該從快取取得）
	tokens1_cached, err := calculator.CalculateTokens(text1, "estimation")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if tokens1 != tokens1_cached {
		t.Errorf("Cached result %d doesn't match original %d", tokens1_cached, tokens1)
	}

	// 計算更多文本以測試快取清除
	calculator.CalculateTokens(text2, "estimation")
	calculator.CalculateTokens(text3, "estimation")

	// 檢查快取統計
	if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
		stats := calcImpl.GetCacheStats()
		cacheSize := stats["cache_size"].(int)
		if cacheSize > 2 {
			t.Errorf("Cache size %d exceeds maximum 2", cacheSize)
		}
	}
}

func TestTokenCalculatorImpl_ValidateText(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	testCases := []struct {
		name        string
		text        string
		expectError bool
	}{
		{
			name:        "正常文本",
			text:        "Hello world 你好世界",
			expectError: false,
		},
		{
			name:        "空文本",
			text:        "",
			expectError: false,
		},
		{
			name:        "包含換行符",
			text:        "Hello\nworld\n你好\n世界",
			expectError: false,
		},
		{
			name:        "過多控制字符",
			text:        "\x00\x01\x02\x03\x04",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := calculator.ValidateText(tc.text)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestTokenCalculatorImpl_BatchCalculation(t *testing.T) {
	calculator := NewTokenCalculator(100)

	texts := []string{
		"Hello world",
		"你好世界",
		"Hello 世界",
		"",
		"function test() { return true; }",
	}

	if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
		results, err := calcImpl.CalculateTokensForMultipleTexts(texts, "estimation")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if len(results) != len(texts) {
			t.Errorf("Results length %d doesn't match texts length %d", len(results), len(texts))
			return
		}

		// 檢查空文本的結果
		if results[3] != 0 {
			t.Errorf("Expected 0 tokens for empty text, got %d", results[3])
		}

		// 檢查所有非空文本都有正數結果
		for i, result := range results {
			if i != 3 && result <= 0 { // 跳過空文本
				t.Errorf("Expected positive tokens for text %d, got %d", i, result)
			}
		}
	} else {
		t.Errorf("Calculator is not of expected type")
	}
}

func BenchmarkTokenCalculation(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	text := "這是一個用於基準測試的文本，包含中文和English混合內容。This is a benchmark test text with mixed Chinese and English content."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculator.CalculateTokens(text, "estimation")
	}
}

func BenchmarkTokenCalculationWithCache(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	text := "這是一個用於基準測試的文本，包含中文和English混合內容。This is a benchmark test text with mixed Chinese and English content."

	// 預先計算一次以填充快取
	calculator.CalculateTokens(text, "estimation")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculator.CalculateTokens(text, "estimation")
	}
}
func TestTokenCalculatorImpl_TiktokenIntegration(t *testing.T) {
	calculator := NewTokenCalculator(100)

	// 測試 tiktoken 是否可用
	if !calculator.IsTiktokenAvailable() {
		t.Skip("Tiktoken not available, skipping tiktoken-specific tests")
	}

	testCases := []struct {
		name string
		text string
	}{
		{
			name: "英文句子",
			text: "The quick brown fox jumps over the lazy dog.",
		},
		{
			name: "中文句子",
			text: "這是一個測試句子，用來驗證中文 Token 計算。",
		},
		{
			name: "程式碼",
			text: "func main() { fmt.Println(\"Hello, World!\") }",
		},
		{
			name: "混合內容",
			text: "Hello 世界! This is a test 測試.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 測試 tiktoken 方法
			tiktokenTokens, err := calculator.CalculateTokens(tc.text, "tiktoken")
			if err != nil {
				t.Errorf("Tiktoken calculation failed: %v", err)
				return
			}

			// 測試估算方法
			estimationTokens, err := calculator.CalculateTokens(tc.text, "estimation")
			if err != nil {
				t.Errorf("Estimation calculation failed: %v", err)
				return
			}

			// 兩種方法都應該返回正數（對於非空文本）
			if len(tc.text) > 0 {
				if tiktokenTokens <= 0 {
					t.Errorf("Tiktoken returned non-positive tokens: %d", tiktokenTokens)
				}
				if estimationTokens <= 0 {
					t.Errorf("Estimation returned non-positive tokens: %d", estimationTokens)
				}
			}

			t.Logf("Text: %s", tc.text)
			t.Logf("Tiktoken: %d tokens", tiktokenTokens)
			t.Logf("Estimation: %d tokens", estimationTokens)
		})
	}
}

func TestTokenCalculatorImpl_CompareCalculationMethods(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	if !calculator.IsTiktokenAvailable() {
		t.Skip("Tiktoken not available, skipping comparison tests")
	}

	testTexts := []string{
		"Hello world",
		"你好世界",
		"Hello 世界",
		"This is a longer text with multiple words and punctuation!",
		"這是一個較長的中文文本，包含多個詞語和標點符號。",
	}

	for _, text := range testTexts {
		t.Run(fmt.Sprintf("Compare_%s", text[:min(10, len(text))]), func(t *testing.T) {
			comparison := calculator.CompareCalculationMethods(text)

			// 檢查基本結構
			if _, ok := comparison["text_length"]; !ok {
				t.Error("Missing text_length in comparison")
			}

			if _, ok := comparison["estimation"]; !ok {
				t.Error("Missing estimation results in comparison")
			}

			if _, ok := comparison["tiktoken"]; !ok {
				t.Error("Missing tiktoken results in comparison")
			}

			if _, ok := comparison["comparison"]; !ok {
				t.Error("Missing comparison results")
			}

			t.Logf("Comparison for '%s': %+v", text, comparison)
		})
	}
}

func TestTokenCalculatorImpl_TiktokenInfo(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	info := calculator.GetTiktokenInfo()

	// 檢查基本資訊
	if enabled, ok := info["enabled"]; !ok {
		t.Error("Missing 'enabled' field in tiktoken info")
	} else if enabled.(bool) {
		// 如果啟用，應該有更多資訊
		if _, ok := info["encoding"]; !ok {
			t.Error("Missing 'encoding' field when tiktoken is enabled")
		}
		if _, ok := info["model_compatibility"]; !ok {
			t.Error("Missing 'model_compatibility' field when tiktoken is enabled")
		}
	}

	t.Logf("Tiktoken info: %+v", info)
}

// min 輔助函數
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
