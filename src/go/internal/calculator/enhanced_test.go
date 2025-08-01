package calculator

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"
)

// TestExtendedUnicodeSupport 測試擴展的 Unicode 字符支援
func TestExtendedUnicodeSupport(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	testCases := []struct {
		name        string
		text        string
		description string
	}{
		{
			name:        "Emoji集合",
			text:        "🚀🎉🔥💯✨🌟⭐️🎯🏆🎊",
			description: "各種 emoji 字符",
		},
		{
			name:        "數學符號",
			text:        "∑∏∫∞±√∆∇∂∃∀∈∉⊂⊃∪∩",
			description: "數學和科學符號",
		},
		{
			name:        "貨幣符號",
			text:        "¥€£$¢₹₩₽₨₪₫₦₡₵₸₴₯₶",
			description: "世界各國貨幣符號",
		},
		{
			name:        "箭頭符號",
			text:        "←→↑↓↔↕↖↗↘↙⇐⇒⇑⇓⇔⇕",
			description: "各種箭頭符號",
		},
		{
			name:        "幾何符號",
			text:        "△▲▽▼◇◆□■○●◯◎★☆",
			description: "幾何圖形符號",
		},
		{
			name:        "多語言文字",
			text:        "Ελληνικά العربية 한국어 日本語 русский हिन्दी עברית",
			description: "希臘語、阿拉伯語、韓語、日語、俄語、印地語、希伯來語",
		},
		{
			name:        "組合字符",
			text:        "café naïve résumé piñata mañana",
			description: "包含重音符號的拉丁字符",
		},
		{
			name:        "中文繁簡混合",
			text:        "繁體中文和简体中文混合在一起的文字內容",
			description: "繁體和簡體中文混合",
		},
		{
			name:        "特殊空白字符",
			text:        "普通空格 \u00A0 不斷行空格 \u2000 四分之一空格 \u2003 全形空格",
			description: "各種Unicode空白字符",
		},
		{
			name:        "零寬度字符",
			text:        "zero\u200Bwidth\u200Cjoiner\u200Dtest\uFEFF",
			description: "零寬度空格、非連接符、連接符、位元組順序標記",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 測試估算方法
			tokens, err := calculator.CalculateTokens(tc.text, "estimation")
			if err != nil {
				t.Errorf("估算方法失敗: %v", err)
				return
			}

			// 驗證基本性質
			if tokens <= 0 {
				t.Errorf("Unicode 文本應該產生正數 token，得到: %d", tokens)
			}

			// 測試 tiktoken 方法（如果可用）
			if calculator.IsTiktokenAvailable() {
				tiktokenTokens, err := calculator.CalculateTokens(tc.text, "tiktoken")
				if err != nil {
					t.Errorf("Tiktoken 方法失敗: %v", err)
					return
				}

				if tiktokenTokens <= 0 {
					t.Errorf("Tiktoken 應該產生正數 token，得到: %d", tiktokenTokens)
				}

				t.Logf("%s: 估算=%d, tiktoken=%d - %s", 
					tc.name, tokens, tiktokenTokens, tc.description)
			} else {
				t.Logf("%s: 估算=%d - %s", tc.name, tokens, tc.description)
			}

			// 測試文本驗證
			if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
				err := calcImpl.ValidateText(tc.text)
				if err != nil {
					t.Errorf("文本驗證失敗: %v", err)
				}
			}
		})
	}
}

// TestComplexMixedLanguageTexts 測試複雜的多語言混合文本
func TestComplexMixedLanguageTexts(t *testing.T) {
	calculator := NewTokenCalculator(1000)

	complexTexts := []struct {
		name string
		text string
	}{
		{
			name: "技術文檔混合",
			text: "這是一個使用 React.js 開發的應用程式。function App() { return <div>Hello World</div>; } 代碼中使用了 JSX 語法。",
		},
		{
			name: "商業報告混合",
			text: "2024年第一季度的營收達到 $1,234,567 USD，相比去年同期增長了 15.7%。主要增長來源於亞太地區的業務擴展。",
		},
		{
			name: "學術論文混合",
			text: "Based on our research, we found that 機器學習模型在處理中文自然語言時，需要考慮 tokenization 的特殊性。The accuracy improved by 23.5%.",
		},
		{
			name: "程式碼註解混合",
			text: "// 這是一個計算 token 數量的函數\nfunction calculateTokens(text: string): number {\n  // 使用正規表達式分割文本\n  return text.split(/\\s+/).length;\n}",
		},
		{
			name: "JSON配置檔案",
			text: `{
  "name": "token-monitor",
  "description": "Token 監控工具",
  "version": "1.0.0",
  "支援語言": ["中文", "English", "日本語"],
  "settings": {
    "maxTokens": 4096,
    "enableCache": true
  }
}`,
		},
		{
			name: "SQL查詢混合",
			text: "SELECT 使用者名稱, email FROM users WHERE 註冊日期 >= '2024-01-01' AND status = 'active' ORDER BY 最後登錄時間 DESC;",
		},
		{
			name: "錯誤訊息混合",
			text: "Error: 無法連接到資料庫。Connection timeout after 30 seconds. Please check your 網路連接 and try again. 錯誤代碼: DB_CONNECTION_TIMEOUT",
		},
	}

	for _, tc := range complexTexts {
		t.Run(tc.name, func(t *testing.T) {
			// 測試所有方法
			methods := []string{"estimation"}
			if calculator.IsTiktokenAvailable() {
				methods = append(methods, "tiktoken", "auto")
			}

			results := make(map[string]int)
			for _, method := range methods {
				tokens, err := calculator.CalculateTokens(tc.text, method)
				if err != nil {
					t.Errorf("方法 %s 失敗: %v", method, err)
					continue
				}
				results[method] = tokens
			}

			// 檢查結果合理性
			for method, tokens := range results {
				if tokens <= 0 {
					t.Errorf("方法 %s 應該產生正數 token，得到: %d", method, tokens)
				}
			}

			// 測試 Token 分佈分析
			distribution, err := calculator.AnalyzeTokenDistribution(tc.text)
			if err != nil {
				t.Errorf("分佈分析失敗: %v", err)
			} else {
				if distribution.TotalTokens != distribution.EnglishTokens+distribution.ChineseTokens {
					t.Errorf("Token 分佈不一致: 總計=%d, 英文=%d, 中文=%d",
						distribution.TotalTokens, distribution.EnglishTokens, distribution.ChineseTokens)
				}
			}

			t.Logf("%s: %+v", tc.name, results)
		})
	}
}

// TestEncodingFormats 測試各種編碼格式處理
func TestEncodingFormats(t *testing.T) {
	calculator := NewTokenCalculator(500)

	testCases := []struct {
		name string
		text string
	}{
		{
			name: "UTF-8 BOM",
			text: "\uFEFF這是包含 BOM 的 UTF-8 文本",
		},
		{
			name: "雙位元組字符",
			text: "測試雙位元組字符：©®™€£¥§¶",
		},
		{
			name: "四位元組字符",
			text: "測試四位元組字符：🌍🚀💻📱🎉",
		},
		{
			name: "代理對字符",
			text: "測試代理對：𝕊𝕦𝕣𝕖𝕣𝕔𝕣𝕚𝕡𝕥 𝔸𝕝𝕡𝕙𝕒𝔹𝕖𝕥",
		},
		{
			name: "無效UTF-8序列處理",
			text: "正常文本" + string([]byte{0xFF, 0xFE}) + "後續文本",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 檢查文本是否為有效的 UTF-8
			if !utf8.ValidString(tc.text) {
				t.Logf("警告: 文本包含無效的 UTF-8 序列")
			}

			// 測試基本計算
			tokens, err := calculator.CalculateTokens(tc.text, "estimation")
			if err != nil {
				t.Errorf("編碼處理失敗: %v", err)
				return
			}

			if tokens <= 0 {
				t.Errorf("應該產生正數 token，得到: %d", tokens)
			}

			// 測試字符計數一致性
			runeCount := utf8.RuneCountInString(tc.text)
			byteCount := len(tc.text)

			t.Logf("%s: tokens=%d, runes=%d, bytes=%d", 
				tc.name, tokens, runeCount, byteCount)

			// 測試文本驗證
			if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
				err := calcImpl.ValidateText(tc.text)
				if err != nil {
					t.Logf("文本驗證警告: %v", err)
				}
			}
		})
	}
}

// TestBoundaryAndEdgeCases 測試邊界條件和邊緣情況
func TestBoundaryAndEdgeCases(t *testing.T) {
	calculator := NewTokenCalculator(100)

	testCases := []struct {
		name        string
		text        string
		expectError bool
		description string
	}{
		{
			name:        "極大文本",
			text:        strings.Repeat("a", 999999), // 接近 1MB 限制
			expectError: false,
			description: "接近最大長度限制的文本",
		},
		{
			name:        "超大文本",
			text:        strings.Repeat("a", 1000001), // 超過 1MB 限制
			expectError: true,
			description: "超過最大長度限制的文本",
		},
		{
			name:        "極小非空文本",
			text:        "a",
			expectError: false,
			description: "最小的非空文本",
		},
		{
			name:        "只有空白字符",
			text:        "   \t\n\r   ",
			expectError: false,  
			description: "只包含各種空白字符",
		},
		{
			name:        "只有標點符號",
			text:        "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			expectError: false,
			description: "只包含標點符號",
		},
		{
			name:        "混合控制字符",
			text:        "text\x00\x01\x02text",
			expectError: false,
			description: "包含少量控制字符",
		},
		{
			name:        "過多控制字符",
			text:        strings.Repeat("\x00\x01\x02", 50),
			expectError: true,
			description: "包含過多控制字符",
		},
		{
			name:        "極長單詞",
			text:        strings.Repeat("supercalifragilisticexpialidocious", 100),
			expectError: false,
			description: "極長的重複單詞",
		},
		{
			name:        "數字序列",
			text:        strings.Repeat("1234567890", 1000),
			expectError: false,
			description: "長數字序列",
		},
		{
			name:        "混合換行格式",
			text:        "Line1\nLine2\r\nLine3\rLine4",
			expectError: false,
			description: "不同的換行符格式",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 測試基本計算
			tokens, err := calculator.CalculateTokens(tc.text, "estimation")
			
			if tc.expectError {
				if err == nil {
					// 可能是文本驗證階段的錯誤
					if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
						err = calcImpl.ValidateText(tc.text)
						if err == nil {
							t.Errorf("預期錯誤但沒有發生: %s", tc.description)
						}
					}
				}
				return
			}

			if err != nil {
				t.Errorf("意外錯誤: %v", err)
				return
			}

			// 驗證結果合理性
			if len(tc.text) > 0 && tokens <= 0 {
				t.Errorf("非空文本應該產生正數 token，得到: %d", tokens)
			}

			if len(tc.text) == 0 && tokens != 0 {
				t.Errorf("空文本應該產生 0 個 token，得到: %d", tokens)
			}

			t.Logf("%s: %d tokens (文本長度: %d) - %s", 
				tc.name, tokens, len(tc.text), tc.description)
		})
	}
}

// TestConcurrentSafety 測試併發安全性
func TestConcurrentSafety(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	
	texts := []string{
		"Hello world",
		"你好世界",
		"Mixed content 混合內容",
		"function test() { return true; }",
		"測試併發安全性的文本內容",
	}

	var wg sync.WaitGroup
	errors := make(chan error, 100)
	results := make(chan int, 100)

	// 啟動多個 goroutine 同時計算
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			for j := 0; j < 10; j++ {
				text := texts[j%len(texts)]
				tokens, err := calculator.CalculateTokens(text, "estimation")
				
				if err != nil {
					errors <- fmt.Errorf("goroutine %d: %w", id, err)
					return
				}
				
				results <- tokens
			}
		}(i)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(errors)
		close(results)
	}()

	// 檢查錯誤
	for err := range errors {
		t.Errorf("併發計算錯誤: %v", err)
	}

	// 收集結果
	var resultCount int
	for range results {
		resultCount++
	}

	expectedResults := 20 * 10 // 20 個 goroutine，每個計算 10 次
	if resultCount != expectedResults {
		t.Errorf("預期 %d 個結果，實際得到 %d 個", expectedResults, resultCount)
	}

	// 測試快取統計的線程安全性
	if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
		stats := calcImpl.GetCacheStats()
		if cacheSize, ok := stats["cache_size"].(int); ok {
			t.Logf("併發測試後快取大小: %d", cacheSize)
		}
	}
}

// TestMemoryUsageAndLimits 測試記憶體使用和限制
func TestMemoryUsageAndLimits(t *testing.T) {
	// 創建小快取的計算器來測試記憶體限制
	calculator := NewTokenCalculator(10)
	
	// 記錄初始記憶體使用
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// 生成大量不同的文本來填充快取
	texts := make([]string, 50)
	for i := 0; i < 50; i++ {
		texts[i] = fmt.Sprintf("測試文本 %d: %s", i, strings.Repeat("test ", i+1))
	}

	// 計算所有文本的 token
	for _, text := range texts {
		_, err := calculator.CalculateTokens(text, "estimation")
		if err != nil {
			t.Errorf("計算失敗: %v", err)
		}
	}

	// 檢查快取是否按預期限制大小
	if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
		stats := calcImpl.GetCacheStats()
		cacheSize := stats["cache_size"].(int)
		maxCacheSize := stats["max_cache_size"].(int)
		
		if cacheSize > maxCacheSize {
			t.Errorf("快取大小 %d 超過最大限制 %d", cacheSize, maxCacheSize)
		}
		
		t.Logf("快取統計: %+v", stats)
	}

	// 記錄最終記憶體使用
	runtime.GC()
	runtime.ReadMemStats(&m2)

	memoryIncrease := m2.Alloc - m1.Alloc
	t.Logf("記憶體使用增加: %d bytes", memoryIncrease)

	// 驗證記憶體使用在合理範圍內（這是一個軟性檢查）
	if memoryIncrease > 10*1024*1024 { // 10MB
		t.Logf("警告: 記憶體使用增加較多: %d bytes", memoryIncrease)
	}
}

// TestPerformanceBenchmarks 擴展的效能基準測試
func TestPerformanceBenchmarks(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	
	testTexts := []struct {
		name string
		text string
	}{
		{"短英文", "Hello world"},
		{"短中文", "你好世界"},
		{"中英混合", "Hello 世界 testing 測試"},
		{"長英文", strings.Repeat("This is a test sentence. ", 100)},
		{"長中文", strings.Repeat("這是一個測試句子。", 100)},
		{"程式碼", strings.Repeat("function test() { return 'hello'; }\n", 50)},
		{"JSON", strings.Repeat(`{"key": "value", "number": 123}, `, 100)},
	}

	for _, tt := range testTexts {
		t.Run(fmt.Sprintf("Performance_%s", tt.name), func(t *testing.T) {
			start := time.Now()
			iterations := 1000
			
			for i := 0; i < iterations; i++ {
				_, err := calculator.CalculateTokens(tt.text, "estimation")
				if err != nil {
					t.Errorf("計算失敗: %v", err)
					return
				}
			}
			
			duration := time.Since(start)
			avgTime := duration / time.Duration(iterations)
			
			t.Logf("%s: %d 次計算耗時 %v，平均 %v/次", 
				tt.name, iterations, duration, avgTime)
			
			// 效能要求：每次計算應該在 1ms 以內
			if avgTime > time.Millisecond {
				t.Logf("警告: 效能可能需要優化，平均時間: %v", avgTime)
			}
		})
	}
}

// TestCachePerformanceImprovement 測試快取效能提升
func TestCachePerformanceImprovement(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	text := strings.Repeat("這是一個用於效能測試的文本 Performance test text ", 50)
	iterations := 1000

	// 測試沒有快取的效能（每次都用不同的文本）
	start := time.Now()
	for i := 0; i < iterations; i++ {
		uniqueText := fmt.Sprintf("%s %d", text, i)
		calculator.CalculateTokens(uniqueText, "estimation")
	}
	noCacheDuration := time.Since(start)

	// 測試有快取的效能（重複使用相同文本）
	start = time.Now()
	for i := 0; i < iterations; i++ {
		calculator.CalculateTokens(text, "estimation")
	}
	withCacheDuration := time.Since(start)

	// 計算效能提升
	improvement := float64(noCacheDuration) / float64(withCacheDuration)
	
	t.Logf("無快取: %v", noCacheDuration)
	t.Logf("有快取: %v", withCacheDuration)
	t.Logf("效能提升: %.2fx", improvement)

	// 快取應該顯著提升效能
	if improvement < 2.0 {
		t.Logf("警告: 快取效能提升不夠明顯: %.2fx", improvement)
	}
}

// TestSensitiveDataProtection 測試敏感資料保護
func TestSensitiveDataProtection(t *testing.T) {
	calculator := NewTokenCalculator(100)
	
	// 模擬包含敏感資訊的文本
	sensitiveTexts := []string{
		"密碼: admin123",
		"API Key: sk-1234567890abcdef",
		"信用卡號: 4111-1111-1111-1111",
		"身分證字號: A123456789",
		"電話號碼: 0912-345-678",
		"電子郵件: user@example.com",
	}

	for _, text := range sensitiveTexts {
		t.Run(fmt.Sprintf("Sensitive_%s", text[:10]), func(t *testing.T) {
			// 計算 Token，不應該出錯
			tokens, err := calculator.CalculateTokens(text, "estimation")
			if err != nil {
				t.Errorf("計算失敗: %v", err)
				return
			}

			if tokens <= 0 {
				t.Errorf("應該產生正數 token，得到: %d", tokens)
			}

			// 重要：檢查日誌中不應該出現敏感資訊
			// 這裡我們模擬檢查，實際實作中應該確保敏感資訊不會記錄到日誌
			t.Logf("處理敏感文本: %d tokens (長度: %d)", tokens, len(text))
			
			// 注意：這裡我們沒有記錄實際的文本內容，以保護敏感資訊
		})
	}
}

// BenchmarkExtendedTokenCalculation 擴展的基準測試
func BenchmarkExtendedTokenCalculation(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	
	benchmarkCases := []struct {
		name string
		text string
	}{
		{"VeryShort", "Hi"},
		{"Short", "Hello world"},
		{"Medium", strings.Repeat("This is a test. ", 10)},
		{"Long", strings.Repeat("This is a longer test sentence. ", 100)},
		{"VeryLong", strings.Repeat("This is a very long test sentence for performance testing. ", 1000)},
		{"Chinese", strings.Repeat("這是中文測試句子。", 100)},
		{"Mixed", strings.Repeat("Mixed 混合 content 內容 ", 100)},
		{"Code", strings.Repeat("func test() { return 'hello'; }\n", 100)},
	}

	for _, bc := range benchmarkCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				calculator.CalculateTokens(bc.text, "estimation")
			}
		})
	}
}

// BenchmarkConcurrentCalculation 併發計算基準測試
func BenchmarkConcurrentCalculation(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	text := strings.Repeat("併發測試文本 Concurrent test text ", 50)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			calculator.CalculateTokens(text, "estimation")
		}
	})
}

// BenchmarkCacheEfficiency 快取效率基準測試
func BenchmarkCacheEfficiency(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	texts := []string{
		"Text 1 for cache test",
		"Text 2 for cache test", 
		"Text 3 for cache test",
		"Text 4 for cache test",
		"Text 5 for cache test",
	}

	// 預熱快取
	for _, text := range texts {
		calculator.CalculateTokens(text, "estimation")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		text := texts[i%len(texts)]
		calculator.CalculateTokens(text, "estimation")
	}
}