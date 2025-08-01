package calculator

import (
	"strings"
	"testing"
)

// TestComprehensiveTokenCalculation 全面的 Token 計算測試
func TestComprehensiveTokenCalculation(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	testCases := []struct {
		name          string
		text          string
		expectedRange [2]int // [min, max] expected tokens
		description   string
	}{
		{
			name:          "空字串",
			text:          "",
			expectedRange: [2]int{0, 0},
			description:   "空字串應該返回 0 個 token",
		},
		{
			name:          "單個英文字母",
			text:          "a",
			expectedRange: [2]int{1, 1},
			description:   "單個字母應該是 1 個 token",
		},
		{
			name:          "單個中文字",
			text:          "中",
			expectedRange: [2]int{1, 2},
			description:   "單個中文字可能是 1-2 個 token",
		},
		{
			name:          "簡單英文句子",
			text:          "Hello world",
			expectedRange: [2]int{2, 3},
			description:   "簡單英文句子",
		},
		{
			name:          "簡單中文句子",
			text:          "你好世界",
			expectedRange: [2]int{2, 8},
			description:   "簡單中文句子，tiktoken 通常產生更多 token",
		},
		{
			name:          "中英混合短句",
			text:          "Hello 世界",
			expectedRange: [2]int{2, 6},
			description:   "中英混合文本",
		},
		{
			name:          "中英混合長句",
			text:          "這是一個 test 測試，包含 English 和中文 content。",
			expectedRange: [2]int{10, 25},
			description:   "複雜的中英混合文本",
		},
		{
			name:          "程式碼片段",
			text:          "func main() { fmt.Println(\"Hello, World!\") }",
			expectedRange: [2]int{10, 15},
			description:   "Go 程式碼片段",
		},
		{
			name:          "JSON 格式",
			text:          `{"name": "test", "value": 123, "enabled": true}`,
			expectedRange: [2]int{12, 18},
			description:   "JSON 格式文本",
		},
		{
			name:          "包含標點符號",
			text:          "Hello, world! How are you? I'm fine, thank you.",
			expectedRange: [2]int{12, 18},
			description:   "包含各種標點符號的英文",
		},
		{
			name:          "包含中文標點",
			text:          "你好，世界！你好嗎？我很好，謝謝你。",
			expectedRange: [2]int{8, 30},
			description:   "包含中文標點符號",
		},
		{
			name:          "數字和符號",
			text:          "Price: $123.45, Quantity: 10, Total: $1,234.50",
			expectedRange: [2]int{12, 20},
			description:   "包含數字、貨幣符號和逗號",
		},
		{
			name:          "特殊字符",
			text:          "Email: user@example.com, URL: https://example.com",
			expectedRange: [2]int{10, 16},
			description:   "包含 email 和 URL",
		},
		{
			name:          "多行文本",
			text:          "Line 1\nLine 2\nLine 3",
			expectedRange: [2]int{6, 12},
			description:   "包含換行符的多行文本",
		},
		{
			name:          "重複文字",
			text:          "test test test test test",
			expectedRange: [2]int{5, 10},
			description:   "重複的單詞",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 測試預設方法
			tokens, err := calculator.CalculateTokens(tc.text, "")
			if err != nil {
				t.Errorf("計算失敗: %v", err)
				return
			}

			if tokens < tc.expectedRange[0] || tokens > tc.expectedRange[1] {
				t.Errorf("Token 數量 %d 不在預期範圍 [%d, %d] 內。文本: %q",
					tokens, tc.expectedRange[0], tc.expectedRange[1], tc.text)
			}

			t.Logf("%s: %d tokens (預期: %d-%d) - %s",
				tc.name, tokens, tc.expectedRange[0], tc.expectedRange[1], tc.description)
		})
	}
}

// TestMethodConsistency 測試不同計算方法的一致性
func TestMethodConsistency(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	testTexts := []string{
		"Hello world",
		"你好世界",
		"Hello 世界",
		"This is a test sentence.",
		"這是一個測試句子。",
		"Mixed content: 混合內容 with English and 中文.",
		"func main() { fmt.Println(\"Hello\") }",
		"",
	}

	for _, text := range testTexts {
		t.Run("Consistency_"+text[:min(20, len(text))], func(t *testing.T) {
			// 測試估算方法
			estimationTokens, err := calculator.CalculateTokens(text, "estimation")
			if err != nil {
				t.Errorf("估算方法失敗: %v", err)
				return
			}

			// 測試 tiktoken 方法（如果可用）
			if calculator.IsTiktokenAvailable() {
				tiktokenTokens, err := calculator.CalculateTokens(text, "tiktoken")
				if err != nil {
					t.Errorf("Tiktoken 方法失敗: %v", err)
					return
				}

				// 測試自動選擇方法
				autoTokens, err := calculator.CalculateTokens(text, "auto")
				if err != nil {
					t.Errorf("自動方法失敗: %v", err)
					return
				}

				// 自動方法應該選擇 tiktoken（如果可用）
				if autoTokens != tiktokenTokens {
					t.Errorf("自動方法 (%d) 應該與 tiktoken (%d) 一致", autoTokens, tiktokenTokens)
				}

				t.Logf("文本: %q", text)
				t.Logf("  估算: %d tokens", estimationTokens)
				t.Logf("  Tiktoken: %d tokens", tiktokenTokens)
				t.Logf("  自動: %d tokens", autoTokens)

				// 計算差異百分比
				if tiktokenTokens > 0 {
					diff := abs(tiktokenTokens - estimationTokens)
					accuracy := 100.0 - (float64(diff)/float64(tiktokenTokens))*100.0
					t.Logf("  準確度: %.1f%%", accuracy)
				}
			} else {
				t.Logf("Tiktoken 不可用，跳過比較測試")
			}
		})
	}
}

// TestBoundaryConditions 測試邊界條件
func TestBoundaryConditions(t *testing.T) {
	calculator := NewTokenCalculator(10) // 小的快取大小

	boundaryTests := []struct {
		name        string
		text        string
		expectError bool
		description string
	}{
		{
			name:        "空字串",
			text:        "",
			expectError: false,
			description: "空字串應該正常處理",
		},
		{
			name:        "單個空格",
			text:        " ",
			expectError: false,
			description: "單個空格",
		},
		{
			name:        "多個空格",
			text:        "   ",
			expectError: false,
			description: "多個空格",
		},
		{
			name:        "只有換行符",
			text:        "\n\n\n",
			expectError: false,
			description: "只有換行符",
		},
		{
			name:        "只有標點符號",
			text:        "!@#$%^&*()",
			expectError: false,
			description: "只有標點符號",
		},
		{
			name:        "極長文本",
			text:        strings.Repeat("This is a test sentence. ", 100),
			expectError: false,
			description: "極長文本（2500+ 字符）",
		},
		{
			name:        "Unicode 字符",
			text:        "🚀 🎉 🔥 💯 ✨",
			expectError: false,
			description: "Unicode emoji 字符",
		},
		{
			name:        "混合 Unicode",
			text:        "Hello 🌍 世界 🚀 World",
			expectError: false,
			description: "混合 Unicode、中英文",
		},
		{
			name:        "控制字符",
			text:        "Hello\tWorld\r\n",
			expectError: false,
			description: "包含 tab 和回車換行",
		},
		{
			name:        "重複字符",
			text:        strings.Repeat("a", 1000),
			expectError: false,
			description: "1000 個重複字符",
		},
	}

	for _, tc := range boundaryTests {
		t.Run(tc.name, func(t *testing.T) {
			tokens, err := calculator.CalculateTokens(tc.text, "")

			if tc.expectError && err == nil {
				t.Errorf("預期錯誤但沒有發生錯誤")
				return
			}

			if !tc.expectError && err != nil {
				t.Errorf("意外錯誤: %v", err)
				return
			}

			if !tc.expectError {
				if tokens < 0 {
					t.Errorf("Token 數量不應該是負數: %d", tokens)
				}

				t.Logf("%s: %d tokens - %s", tc.name, tokens, tc.description)
			}
		})
	}
}

// TestCacheEffectiveness 測試快取效果
func TestCacheEffectiveness(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	testText := "This is a test text for cache effectiveness."

	// 第一次計算（應該計算並快取）
	tokens1, err := calculator.CalculateTokens(testText, "estimation")
	if err != nil {
		t.Errorf("第一次計算失敗: %v", err)
		return
	}

	// 第二次計算（應該從快取取得）
	tokens2, err := calculator.CalculateTokens(testText, "estimation")
	if err != nil {
		t.Errorf("第二次計算失敗: %v", err)
		return
	}

	// 結果應該一致
	if tokens1 != tokens2 {
		t.Errorf("快取結果不一致: 第一次 %d, 第二次 %d", tokens1, tokens2)
	}

	// 檢查快取統計
	stats := calculator.GetCacheStats()
	cacheSize := stats["cache_size"].(int)
	if cacheSize == 0 {
		t.Error("快取應該包含至少一個項目")
	}

	t.Logf("快取統計: %+v", stats)
}

// TestTokenDistributionAccuracy 測試 Token 分佈準確性
func TestTokenDistributionAccuracy(t *testing.T) {
	calculator := NewTokenCalculator(100)

	testCases := []struct {
		name string
		text string
	}{
		{"純英文", "Hello world this is a test"},
		{"純中文", "這是一個測試句子"},
		{"中英混合", "Hello 世界 this is 測試"},
		{"程式碼", "func main() { fmt.Println(\"test\") }"},
		{"數字符號", "Price: $123.45, Total: 100%"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distribution, err := calculator.AnalyzeTokenDistribution(tc.text)
			if err != nil {
				t.Errorf("分析失敗: %v", err)
				return
			}

			// 基本一致性檢查
			if distribution.TotalTokens != distribution.EnglishTokens+distribution.ChineseTokens {
				t.Errorf("總 Token 數 (%d) 不等於英文 (%d) + 中文 (%d)",
					distribution.TotalTokens, distribution.EnglishTokens, distribution.ChineseTokens)
			}

			// 非負數檢查
			if distribution.EnglishTokens < 0 || distribution.ChineseTokens < 0 || distribution.TotalTokens < 0 {
				t.Errorf("Token 數量不應該是負數: 英文=%d, 中文=%d, 總計=%d",
					distribution.EnglishTokens, distribution.ChineseTokens, distribution.TotalTokens)
			}

			t.Logf("%s: 總計=%d, 英文=%d, 中文=%d, 方法=%s",
				tc.name, distribution.TotalTokens, distribution.EnglishTokens,
				distribution.ChineseTokens, distribution.Method)
		})
	}
}

// TestErrorHandling 測試錯誤處理
func TestErrorHandling(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	// 測試無效方法
	_, err := calculator.CalculateTokens("test", "invalid_method")
	if err != nil {
		t.Logf("無效方法正確返回錯誤: %v", err)
	}

	// 測試文本驗證
	err = calculator.ValidateText(strings.Repeat("a", 2000000)) // 超過限制
	if err == nil {
		t.Error("超大文本應該返回驗證錯誤")
	}

	// 測試控制字符過多的文本
	controlText := strings.Repeat("\x00\x01\x02", 100)
	err = calculator.ValidateText(controlText)
	if err == nil {
		t.Error("過多控制字符應該返回驗證錯誤")
	}
}

// 輔助函數已在其他檔案中定義
