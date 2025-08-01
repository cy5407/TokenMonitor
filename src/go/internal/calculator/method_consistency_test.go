package calculator

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// TestTiktokenVsEstimationAccuracy 測試 tiktoken 和估算方法的準確性
func TestTiktokenVsEstimationAccuracy(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	if !calculator.IsTiktokenAvailable() {
		t.Skip("Tiktoken 不可用，跳過準確性比較測試")
	}

	testCases := []struct {
		name           string
		text           string
		expectedAccuracy float64 // 預期準確度閾值 (%)
		description    string
	}{
		{
			name:           "簡單英文",
			text:           "Hello world",
			expectedAccuracy: 80.0,
			description:    "簡單英文句子應該有較高準確度",
		},
		{
			name:           "簡單中文", 
			text:           "你好世界",
			expectedAccuracy: 60.0,
			description:    "簡單中文準確度可能較低",
		},
		{
			name:           "中英混合",
			text:           "Hello 世界 testing 測試",
			expectedAccuracy: 70.0,
			description:    "中英混合文本",
		},
		{
			name:           "長英文段落",
			text:           "This is a longer paragraph with multiple sentences. It contains various words and punctuation marks. The purpose is to test token calculation accuracy on longer English texts.",
			expectedAccuracy: 85.0,
			description:    "長英文段落準確度應該較高",
		},
		{
			name:           "長中文段落",
			text:           "這是一個較長的中文段落，包含多個句子。它包含各種詞語和標點符號。目的是測試在較長中文文本上的 Token 計算準確度。",
			expectedAccuracy: 65.0,
			description:    "長中文段落準確度",
		},
		{
			name:           "程式碼片段",
			text:           "function calculateTokens(text) {\n  const words = text.split(' ');\n  return words.length;\n}",
			expectedAccuracy: 75.0,
			description:    "程式碼應該有中等準確度",
		},
		{
			name:           "JSON格式",
			text:           `{"name": "test", "values": [1, 2, 3], "enabled": true, "description": "測試用的 JSON 資料"}`,
			expectedAccuracy: 70.0,
			description:    "JSON 格式文本",
		},
		{
			name:           "技術文檔混合",
			text:           "The API uses REST endpoints. 端點 /api/users 返回用戶列表。Parameters include limit and offset for pagination.",
			expectedAccuracy: 75.0,
			description:    "技術文檔風格的中英混合",
		},
	}

	var totalAccuracy float64
	var validTests int

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 計算 tiktoken 方法的結果
			tiktokenTokens, err := calculator.CalculateTokens(tc.text, "tiktoken")
			if err != nil {
				t.Errorf("Tiktoken 計算失敗: %v", err)
				return
			}

			// 計算估算方法的結果
			estimationTokens, err := calculator.CalculateTokens(tc.text, "estimation")
			if err != nil {
				t.Errorf("估算方法計算失敗: %v", err)
				return
			}

			// 計算準確度
			if tiktokenTokens == 0 {
				t.Errorf("Tiktoken 返回 0 個 token，無法計算準確度")
				return
			}

			difference := int(math.Abs(float64(tiktokenTokens - estimationTokens)))
			accuracy := 100.0 - (float64(difference)/float64(tiktokenTokens))*100.0

			// 記錄結果
			t.Logf("%s:", tc.name)
			t.Logf("  文本: %q", tc.text)
			t.Logf("  Tiktoken: %d tokens", tiktokenTokens)
			t.Logf("  估算: %d tokens", estimationTokens)
			t.Logf("  差異: %d tokens", difference)
			t.Logf("  準確度: %.1f%%", accuracy)

			// 檢查是否達到預期準確度
			if accuracy < tc.expectedAccuracy {
				t.Logf("警告: 準確度 %.1f%% 低於預期 %.1f%%", accuracy, tc.expectedAccuracy)
			}

			// 累計統計
			totalAccuracy += accuracy
			validTests++
		})
	}

	// 計算平均準確度
	if validTests > 0 {
		avgAccuracy := totalAccuracy / float64(validTests)
		t.Logf("平均準確度: %.1f%% (%d 個測試)", avgAccuracy, validTests)

		// 總體準確度應該達到一定標準
		if avgAccuracy < 70.0 {
			t.Logf("警告: 平均準確度 %.1f%% 可能需要改進", avgAccuracy)
		}
	}
}

// TestCalculationMethodSwitching 測試計算方法切換邏輯
func TestCalculationMethodSwitching(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	testTexts := []string{
		"Hello world",
		"你好世界",
		"Mixed content 混合內容",
		"",
	}

	for _, text := range testTexts {
		t.Run(fmt.Sprintf("Switching_%s", text), func(t *testing.T) {
			// 測試自動方法選擇
			autoTokens, err := calculator.CalculateTokens(text, "auto")
			if err != nil {
				t.Errorf("自動方法失敗: %v", err)
				return
			}

			// 測試預設方法（空字串）
			defaultTokens, err := calculator.CalculateTokens(text, "")
			if err != nil {
				t.Errorf("預設方法失敗: %v", err)
				return
			}

			// 自動和預設應該返回相同結果
			if autoTokens != defaultTokens {
				t.Errorf("自動方法 (%d) 與預設方法 (%d) 結果不一致", autoTokens, defaultTokens)
			}

			// 測試明確指定的方法
			estimationTokens, err := calculator.CalculateTokens(text, "estimation")
			if err != nil {
				t.Errorf("估算方法失敗: %v", err)
				return
			}

			// 如果 tiktoken 可用，測試其與自動方法的一致性
			if calculator.IsTiktokenAvailable() {
				tiktokenTokens, err := calculator.CalculateTokens(text, "tiktoken")
				if err != nil {
					t.Errorf("Tiktoken 方法失敗: %v", err)
					return
				}

				// 自動方法應該選擇 tiktoken
				if autoTokens != tiktokenTokens {
					t.Errorf("自動方法應該選擇 tiktoken: auto=%d, tiktoken=%d", autoTokens, tiktokenTokens)
				}

				t.Logf("文本: %q", text)
				t.Logf("  自動: %d, 預設: %d, 估算: %d, tiktoken: %d", 
					autoTokens, defaultTokens, estimationTokens, tiktokenTokens)
			} else {
				// tiktoken 不可用時，自動方法應該選擇估算
				if autoTokens != estimationTokens {
					t.Errorf("Tiktoken 不可用時，自動方法應該選擇估算: auto=%d, estimation=%d", 
						autoTokens, estimationTokens)
				}

				t.Logf("文本: %q (tiktoken 不可用)", text)
				t.Logf("  自動: %d, 預設: %d, 估算: %d", autoTokens, defaultTokens, estimationTokens)
			}
		})
	}
}

// TestCacheMechanismConsistency 測試快取機制的一致性
func TestCacheMechanismConsistency(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	testText := "這是用於測試快取一致性的文本 Cache consistency test text"

	// 第一次計算 - 應該計算並快取
	tokens1, err := calculator.CalculateTokens(testText, "estimation")
	if err != nil {
		t.Errorf("第一次計算失敗: %v", err)
		return
	}

	// 第二次計算 - 應該從快取取得
	tokens2, err := calculator.CalculateTokens(testText, "estimation")
	if err != nil {
		t.Errorf("第二次計算失敗: %v", err)
		return
	}

	// 結果應該完全一致
	if tokens1 != tokens2 {
		t.Errorf("快取結果不一致: 第一次=%d, 第二次=%d", tokens1, tokens2)
	}

	// 使用不同方法計算相同文本
	if calculator.IsTiktokenAvailable() {
		tiktokenTokens, err := calculator.CalculateTokens(testText, "tiktoken")
		if err != nil {
			t.Errorf("Tiktoken 計算失敗: %v", err)
			return
		}

		// 再次用估算方法計算，應該還是從快取取得
		tokens3, err := calculator.CalculateTokens(testText, "estimation")
		if err != nil {
			t.Errorf("第三次計算失敗: %v", err)
			return
		}

		if tokens3 != tokens1 {
			t.Errorf("快取一致性問題: 估算方法結果不一致 %d vs %d", tokens3, tokens1)
		}

		t.Logf("快取一致性測試通過: 估算=%d, tiktoken=%d", tokens1, tiktokenTokens)
	}

	// 測試快取統計
	stats := calculator.GetCacheStats()
	cacheSize := stats["cache_size"].(int)
	if cacheSize == 0 {
		t.Error("快取應該包含至少一個項目")
	}

	t.Logf("快取統計: %+v", stats)
}

// TestMethodComparisonDetailedAnalysis 詳細的方法比較分析
func TestMethodComparisonDetailedAnalysis(t *testing.T) {
	calculator := NewTokenCalculator(1000).(*TokenCalculatorImpl)

	if !calculator.IsTiktokenAvailable() {
		t.Skip("需要 tiktoken 進行詳細比較分析")
	}

	testCases := []struct {
		category string
		texts    []string
	}{
		{
			category: "英文文本",
			texts: []string{
				"Hello",
				"Hello world",
				"The quick brown fox jumps over the lazy dog.",
				"This is a longer English paragraph with multiple sentences and various punctuation marks.",
			},
		},
		{
			category: "中文文本",
			texts: []string{
				"你好",
				"你好世界",
				"這是一個測試句子，包含標點符號。",
				"這是一個較長的中文段落，包含多個句子和各種標點符號，用來測試中文文本的 Token 計算準確性。",
			},
		},
		{
			category: "混合文本",
			texts: []string{
				"Hello 世界",
				"This is 測試 mixed content.",
				"使用 JavaScript 開發 web application 是很常見的。",
				"The API returns JSON data. 返回的資料格式是 {\"status\": \"success\", \"data\": []}",
			},
		},
		{
			category: "程式碼",
			texts: []string{
				"x = 1",
				"function hello() { return 'world'; }",
				"const users = await db.users.find({ active: true });",
				"// 這是註解\nfunction calculateTokens(text) {\n  return text.split(' ').length;\n}",
			},
		},
	}

	for _, category := range testCases {
		t.Run(category.category, func(t *testing.T) {
			var totalEstimation, totalTiktoken int
			var accuracySum float64

			for i, text := range category.texts {
				comparison := calculator.CompareCalculationMethods(text)
				
				estimation := comparison["estimation"].(map[string]interface{})["tokens"].(int)
				tiktoken := comparison["tiktoken"].(map[string]interface{})["tokens"].(int)
				comparisonData := comparison["comparison"].(map[string]interface{})
				
				accuracy := comparisonData["accuracy_percent"].(float64)
				difference := comparisonData["difference"].(int)

				totalEstimation += estimation
				totalTiktoken += tiktoken
				accuracySum += accuracy

				t.Logf("文本 %d: 估算=%d, tiktoken=%d, 差異=%d, 準確度=%.1f%%", 
					i+1, estimation, tiktoken, difference, accuracy)
			}

			avgAccuracy := accuracySum / float64(len(category.texts))
			
			t.Logf("%s 摘要:", category.category)
			t.Logf("  總估算 tokens: %d", totalEstimation)
			t.Logf("  總 tiktoken tokens: %d", totalTiktoken)
			t.Logf("  平均準確度: %.1f%%", avgAccuracy)

			// 根據不同類別設定不同的準確度期望
			var expectedAccuracy float64
			switch category.category {
			case "英文文本":
				expectedAccuracy = 80.0
			case "中文文本":
				expectedAccuracy = 60.0
			case "混合文本":
				expectedAccuracy = 70.0
			case "程式碼":
				expectedAccuracy = 75.0
			}

			if avgAccuracy < expectedAccuracy {
				t.Logf("警告: %s 的平均準確度 %.1f%% 低於預期 %.1f%%", 
					category.category, avgAccuracy, expectedAccuracy)
			}
		})
	}
}

// TestCacheEvictionBehavior 測試快取淘汰行為
func TestCacheEvictionBehavior(t *testing.T) {
	// 使用小的快取大小來測試淘汰行為
	calculator := NewTokenCalculator(5).(*TokenCalculatorImpl)

	// 生成測試文本
	texts := make([]string, 10)
	for i := 0; i < 10; i++ {
		texts[i] = fmt.Sprintf("測試文本 %d test text %d", i, i)
	}

	// 計算所有文本（超過快取容量）
	results := make([]int, len(texts))
	for i, text := range texts {
		tokens, err := calculator.CalculateTokens(text, "estimation")
		if err != nil {
			t.Errorf("計算文本 %d 失敗: %v", i, err)
			return
		}
		results[i] = tokens
	}

	// 檢查快取大小
	stats := calculator.GetCacheStats()
	cacheSize := stats["cache_size"].(int)
	maxCacheSize := stats["max_cache_size"].(int)

	if cacheSize > maxCacheSize {
		t.Errorf("快取大小 %d 超過最大限制 %d", cacheSize, maxCacheSize)
	}

	t.Logf("快取淘汰測試: 處理了 %d 個文本，快取大小 %d/%d", 
		len(texts), cacheSize, maxCacheSize)

	// 重新計算一些文本，檢查結果一致性
	for i := 0; i < 5; i++ {
		tokens, err := calculator.CalculateTokens(texts[i], "estimation")
		if err != nil {
			t.Errorf("重新計算文本 %d 失敗: %v", i, err)
			continue
		}

		if tokens != results[i] {
			t.Errorf("重新計算結果不一致: 文本 %d, 原始=%d, 重新計算=%d", 
				i, results[i], tokens)
		}
	}
}

// TestConcurrentCacheAccess 測試併發快取存取
func TestConcurrentCacheAccess(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	
	text := "併發快取測試文本 Concurrent cache test text"
	iterations := 100
	goroutines := 10

	// 預先計算一次以填充快取
	expectedTokens, err := calculator.CalculateTokens(text, "estimation")
	if err != nil {
		t.Errorf("預先計算失敗: %v", err)
		return
	}

	// 啟動多個 goroutine 同時存取快取
	results := make(chan int, goroutines*iterations)
	errors := make(chan error, goroutines*iterations)

	for g := 0; g < goroutines; g++ {
		go func() {
			for i := 0; i < iterations; i++ {
				tokens, err := calculator.CalculateTokens(text, "estimation")
				if err != nil {
					errors <- err
					return
				}
				results <- tokens
			}
		}()
	}

	// 收集結果
	var resultCount int
	for resultCount < goroutines*iterations {
		select {
		case err := <-errors:
			t.Errorf("併發計算錯誤: %v", err)
			return
		case tokens := <-results:
			if tokens != expectedTokens {
				t.Errorf("併發結果不一致: 預期=%d, 實際=%d", expectedTokens, tokens)
			}
			resultCount++
		case <-time.After(5 * time.Second):
			t.Errorf("測試超時，只收到 %d/%d 個結果", resultCount, goroutines*iterations)
			return
		}
	}

	t.Logf("併發快取測試通過: %d 個 goroutine，每個 %d 次迭代", goroutines, iterations)
}

// TestMethodFallbackBehavior 測試方法回退行為
func TestMethodFallbackBehavior(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	testText := "Method fallback test 方法回退測試"

	// 測試指定不存在的方法
	tokens, err := calculator.CalculateTokens(testText, "nonexistent_method")
	if err != nil {
		t.Logf("不存在的方法正確返回錯誤: %v", err)
	} else {
		t.Logf("不存在的方法回退到預設方法，tokens: %d", tokens)
	}

	// 測試 tiktoken 不可用時的回退
	originalTiktokenEnabled := calculator.tiktokenEnabled
	originalTiktokenEncoder := calculator.tiktokenEncoder

	// 暫時停用 tiktoken
	calculator.DisableTiktoken()

	tiktokenTokens, err := calculator.CalculateTokens(testText, "tiktoken")
	if err != nil {
		t.Errorf("Tiktoken 不可用時回退失敗: %v", err)
	} else {
		t.Logf("Tiktoken 不可用，回退到估算方法: %d tokens", tiktokenTokens)
	}

	// 恢復 tiktoken 狀態
	calculator.tiktokenEnabled = originalTiktokenEnabled
	calculator.tiktokenEncoder = originalTiktokenEncoder

	// 驗證恢復後的行為
	if calculator.IsTiktokenAvailable() {
		restoredTokens, err := calculator.CalculateTokens(testText, "tiktoken")
		if err != nil {
			t.Errorf("恢復 tiktoken 後計算失敗: %v", err)
		} else {
			t.Logf("Tiktoken 恢復後: %d tokens", restoredTokens)
		}
	}
}

// TestCalculationConsistencyAcrossRuns 測試跨次執行的計算一致性
func TestCalculationConsistencyAcrossRuns(t *testing.T) {
	// 創建多個計算器實例
	calc1 := NewTokenCalculator(1000)
	calc2 := NewTokenCalculator(1000)
	calc3 := NewTokenCalculator(1000)

	testTexts := []string{
		"Hello world",
		"你好世界",
		"Mixed content 混合內容",
		"function test() { return 'hello'; }",
		"這是一個較長的測試文本，包含中文和English混合內容，用來驗證不同計算器實例之間的一致性。",
	}

	for _, text := range testTexts {
		t.Run(fmt.Sprintf("Consistency_%s", text[:20]), func(t *testing.T) {
			// 用三個不同的計算器實例計算相同文本
			tokens1, err1 := calc1.CalculateTokens(text, "estimation")
			tokens2, err2 := calc2.CalculateTokens(text, "estimation")
			tokens3, err3 := calc3.CalculateTokens(text, "estimation")

			// 檢查錯誤
			if err1 != nil || err2 != nil || err3 != nil {
				t.Errorf("計算錯誤: err1=%v, err2=%v, err3=%v", err1, err2, err3)
				return
			}

			// 檢查一致性
			if tokens1 != tokens2 || tokens2 != tokens3 {
				t.Errorf("不同實例結果不一致: %d, %d, %d", tokens1, tokens2, tokens3)
				return
			}

			t.Logf("一致性測試通過: %d tokens", tokens1)

			// 如果 tiktoken 可用，也測試其一致性
			if calc1.IsTiktokenAvailable() {
				tikTokens1, _ := calc1.CalculateTokens(text, "tiktoken")
				tikTokens2, _ := calc2.CalculateTokens(text, "tiktoken")
				tikTokens3, _ := calc3.CalculateTokens(text, "tiktoken")

				if tikTokens1 != tikTokens2 || tikTokens2 != tikTokens3 {
					t.Errorf("Tiktoken 不同實例結果不一致: %d, %d, %d", 
						tikTokens1, tikTokens2, tikTokens3)
				} else {
					t.Logf("Tiktoken 一致性測試通過: %d tokens", tikTokens1)
				}
			}
		})
	}
}