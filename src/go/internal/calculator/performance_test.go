package calculator

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestBatchProcessingPerformance 測試大批量處理效能
func TestBatchProcessingPerformance(t *testing.T) {
	calculator := NewTokenCalculator(1000).(*TokenCalculatorImpl)

	batchSizes := []int{10, 100, 500, 1000, 5000}
	
	for _, batchSize := range batchSizes {
		t.Run(fmt.Sprintf("BatchSize_%d", batchSize), func(t *testing.T) {
			// 生成測試文本
			texts := make([]string, batchSize)
			for i := 0; i < batchSize; i++ {
				texts[i] = fmt.Sprintf("批量測試文本 %d: %s", i, 
					strings.Repeat("test content ", (i%10)+1))
			}

			// 測試批量計算效能
			start := time.Now()
			results, err := calculator.CalculateTokensForMultipleTexts(texts, "estimation")
			duration := time.Since(start)

			if err != nil {
				t.Errorf("批量計算失敗: %v", err)
				return
			}

			if len(results) != batchSize {
				t.Errorf("結果數量不匹配: 預期 %d, 實際 %d", batchSize, len(results))
				return
			}

			// 計算效能指標
			avgTimePerItem := duration / time.Duration(batchSize)
			itemsPerSecond := float64(batchSize) / duration.Seconds()

			t.Logf("批量大小 %d:", batchSize)
			t.Logf("  總耗時: %v", duration)
			t.Logf("  平均每項: %v", avgTimePerItem)
			t.Logf("  處理速度: %.2f 項/秒", itemsPerSecond)

			// 效能要求檢查
			if avgTimePerItem > 10*time.Millisecond {
				t.Logf("警告: 平均處理時間較長: %v", avgTimePerItem)
			}

			if itemsPerSecond < 100 {
				t.Logf("警告: 處理速度較慢: %.2f 項/秒", itemsPerSecond)
			}

			// 驗證結果的正確性
			for i, tokens := range results {
				if tokens <= 0 && len(texts[i]) > 0 {
					t.Errorf("項目 %d 的 token 數量異常: %d", i, tokens)
				}
			}
		})
	}
}

// TestCachePerformanceScaling 測試快取效能擴展性
func TestCachePerformanceScaling(t *testing.T) {
	cacheSizes := []int{10, 100, 500, 1000, 5000}
	
	for _, cacheSize := range cacheSizes {
		t.Run(fmt.Sprintf("CacheSize_%d", cacheSize), func(t *testing.T) {
			calculator := NewTokenCalculator(cacheSize)
			
			// 生成測試文本（數量等於快取大小）
			texts := make([]string, cacheSize)
			for i := 0; i < cacheSize; i++ {
				texts[i] = fmt.Sprintf("快取測試文本 %d", i)
			}

			// 第一次計算：填充快取
			start := time.Now()
			for _, text := range texts {
				calculator.CalculateTokens(text, "estimation")
			}
			fillDuration := time.Since(start)

			// 第二次計算：測試快取命中效能
			start = time.Now()
			for _, text := range texts {
				calculator.CalculateTokens(text, "estimation")
			}
			hitDuration := time.Since(start)

			// 計算效能提升
			improvement := float64(fillDuration) / float64(hitDuration)

			t.Logf("快取大小 %d:", cacheSize)
			t.Logf("  填充時間: %v", fillDuration)
			t.Logf("  命中時間: %v", hitDuration)
			t.Logf("  效能提升: %.2fx", improvement)

			// 快取應該顯著提升效能
			if improvement < 2.0 {
				t.Logf("警告: 快取效能提升不明顯: %.2fx", improvement)
			}

			// 檢查快取統計
			if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
				stats := calcImpl.GetCacheStats()
				usage := stats["cache_usage"].(float64)
				t.Logf("  快取使用率: %.1f%%", usage*100)
			}
		})
	}
}

// TestMemoryUsageEfficiency 測試記憶體使用效率
func TestMemoryUsageEfficiency(t *testing.T) {
	testScenarios := []struct {
		name        string
		textCount   int
		textSize    int
		description string
	}{
		{"小文本大量", 10000, 50, "測試小文本的記憶體效率"},
		{"中文本適量", 1000, 500, "測試中等文本的記憶體使用"},
		{"大文本少量", 100, 5000, "測試大文本的記憶體處理"},
		{"極大文本", 10, 50000, "測試極大文本的記憶體管理"},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			calculator := NewTokenCalculator(1000)

			// 記錄初始記憶體
			runtime.GC()
			var memBefore runtime.MemStats
			runtime.ReadMemStats(&memBefore)

			// 生成測試文本
			baseText := "記憶體效率測試 Memory efficiency test "
			texts := make([]string, scenario.textCount)
			for i := 0; i < scenario.textCount; i++ {
				repeatCount := scenario.textSize / len(baseText)
				if repeatCount == 0 {
					repeatCount = 1
				}
				texts[i] = strings.Repeat(baseText, repeatCount)
			}

			// 執行計算
			start := time.Now()
			for i, text := range texts {
				_, err := calculator.CalculateTokens(text, "estimation")
				if err != nil {
					t.Errorf("計算文本 %d 失敗: %v", i, err)
					return
				}
			}
			duration := time.Since(start)

			// 記錄最終記憶體
			runtime.GC()
			var memAfter runtime.MemStats
			runtime.ReadMemStats(&memAfter)

			// 計算記憶體使用
			memoryIncrease := memAfter.Alloc - memBefore.Alloc
			avgMemoryPerText := memoryIncrease / uint64(scenario.textCount)
			totalTextSize := uint64(scenario.textCount * scenario.textSize)
			memoryEfficiency := float64(totalTextSize) / float64(memoryIncrease)

			t.Logf("%s:", scenario.name)
			t.Logf("  處理文本: %d 個，每個約 %d 字符", scenario.textCount, scenario.textSize)
			t.Logf("  總耗時: %v", duration)
			t.Logf("  記憶體增加: %d KB", memoryIncrease/1024)
			t.Logf("  平均每文本: %d bytes", avgMemoryPerText)
			t.Logf("  記憶體效率: %.2f (文本大小/記憶體使用)", memoryEfficiency)

			// 記憶體效率檢查
			if memoryEfficiency < 1.0 {
				t.Logf("警告: 記憶體使用可能過多，效率: %.2f", memoryEfficiency)
			}

			if avgMemoryPerText > 10*1024 { // 每個文本超過 10KB 記憶體使用
				t.Logf("警告: 平均記憶體使用較高: %d bytes/文本", avgMemoryPerText)
			}
		})
	}
}

// TestConcurrentPerformance 測試併發效能
func TestConcurrentPerformance(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	text := strings.Repeat("併發效能測試 Concurrent performance test ", 100)

	goroutineCounts := []int{1, 2, 4, 8, 16, 32}
	iterationsPerGoroutine := 1000

	for _, goroutineCount := range goroutineCounts {
		t.Run(fmt.Sprintf("Goroutines_%d", goroutineCount), func(t *testing.T) {
			var wg sync.WaitGroup
			errors := make(chan error, goroutineCount)
			completions := make(chan time.Duration, goroutineCount)

			start := time.Now()

			// 啟動指定數量的 goroutine
			for i := 0; i < goroutineCount; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					
					goroutineStart := time.Now()
					
					for j := 0; j < iterationsPerGoroutine; j++ {
						_, err := calculator.CalculateTokens(text, "estimation")
						if err != nil {
							errors <- fmt.Errorf("goroutine %d: %w", id, err)
							return
						}
					}
					
					completions <- time.Since(goroutineStart)
				}(i)
			}

			wg.Wait()
			totalDuration := time.Since(start)
			close(errors)
			close(completions)

			// 檢查錯誤
			errorCount := 0
			for err := range errors {
				t.Errorf("併發錯誤: %v", err)
				errorCount++
			}

			if errorCount > 0 {
				t.Errorf("發生 %d 個併發錯誤", errorCount)
				return
			}

			// 收集效能統計
			var totalGoroutineTime time.Duration
			goroutineCount := 0
			for duration := range completions {
				totalGoroutineTime += duration
				goroutineCount++
			}

			totalOperations := goroutineCount * iterationsPerGoroutine
			operationsPerSecond := float64(totalOperations) / totalDuration.Seconds()
			avgGoroutineTime := totalGoroutineTime / time.Duration(goroutineCount)
			parallelEfficiency := float64(avgGoroutineTime) / float64(totalDuration)

			t.Logf("%d 個 goroutine:", goroutineCount)
			t.Logf("  總耗時: %v", totalDuration)
			t.Logf("  平均 goroutine 時間: %v", avgGoroutineTime)
			t.Logf("  總操作數: %d", totalOperations)
			t.Logf("  操作/秒: %.2f", operationsPerSecond)
			t.Logf("  並行效率: %.2f%%", parallelEfficiency*100)

			// 效能基準檢查
			if operationsPerSecond < 1000 {
				t.Logf("警告: 併發效能較低: %.2f 操作/秒", operationsPerSecond)
			}
		})
	}
}

// TestTextSizeScalability 測試文本大小擴展性
func TestTextSizeScalability(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	
	textSizes := []int{10, 100, 1000, 10000, 100000}
	baseText := "擴展性測試 Scalability test "

	for _, size := range textSizes {
		t.Run(fmt.Sprintf("TextSize_%d", size), func(t *testing.T) {
			// 生成指定大小的文本
			repeatCount := size / len(baseText)
			if repeatCount == 0 {
				repeatCount = 1
			}
			text := strings.Repeat(baseText, repeatCount)
			actualSize := len(text)

			// 測試計算時間
			iterations := 100
			if size > 10000 {
				iterations = 10 // 大文本減少迭代次數
			}

			start := time.Now()
			for i := 0; i < iterations; i++ {
				_, err := calculator.CalculateTokens(text, "estimation")
				if err != nil {
					t.Errorf("計算失敗 (iteration %d): %v", i, err)
					return
				}
			}
			duration := time.Since(start)

			avgTimePerCalc := duration / time.Duration(iterations)
			charsPerSecond := float64(actualSize*iterations) / duration.Seconds()
			timePerChar := duration / time.Duration(actualSize*iterations)

			t.Logf("文本大小 %d 字符:", actualSize)
			t.Logf("  迭代次數: %d", iterations)
			t.Logf("  總耗時: %v", duration)
			t.Logf("  平均每次: %v", avgTimePerCalc)
			t.Logf("  處理速度: %.0f 字符/秒", charsPerSecond)
			t.Logf("  每字符時間: %v", timePerChar)

			// 擴展性檢查
			if timePerChar > 10*time.Nanosecond {
				t.Logf("警告: 每字符處理時間較長: %v", timePerChar)
			}

			// 檢查時間複雜度是否合理（應該接近線性）
			if size > 1000 && avgTimePerCalc > 100*time.Millisecond {
				t.Logf("警告: 大文本處理時間可能過長: %v", avgTimePerCalc)
			}
		})
	}
}

// TestCacheEvictionPerformance 測試快取淘汰效能
func TestCacheEvictionPerformance(t *testing.T) {
	cacheSize := 100
	calculator := NewTokenCalculator(cacheSize).(*TokenCalculatorImpl)

	// 生成測試文本（數量遠超快取大小）
	textCount := cacheSize * 5
	texts := make([]string, textCount)
	for i := 0; i < textCount; i++ {
		texts[i] = fmt.Sprintf("快取淘汰測試 %d: %s", i, strings.Repeat("content ", i%20+1))
	}

	// 測試快取淘汰的效能影響
	start := time.Now()
	for i, text := range texts {
		_, err := calculator.CalculateTokens(text, "estimation")
		if err != nil {
			t.Errorf("計算文本 %d 失敗: %v", i, err)
			return
		}

		// 定期檢查快取狀態
		if i%100 == 0 {
			stats := calculator.GetCacheStats()
			currentCacheSize := stats["cache_size"].(int)
			t.Logf("處理 %d 個文本後，快取大小: %d", i, currentCacheSize)
		}
	}
	duration := time.Since(start)

	// 檢查最終快取狀態
	stats := calculator.GetCacheStats()
	finalCacheSize := stats["cache_size"].(int)
	
	t.Logf("快取淘汰效能測試:")
	t.Logf("  處理文本數: %d", textCount)
	t.Logf("  快取大小限制: %d", cacheSize)
	t.Logf("  最終快取大小: %d", finalCacheSize)
	t.Logf("  總耗時: %v", duration)
	t.Logf("  平均每項: %v", duration/time.Duration(textCount))

	// 驗證快取大小控制
	if finalCacheSize > cacheSize {
		t.Errorf("快取大小超出限制: %d > %d", finalCacheSize, cacheSize)
	}

	// 測試快取淘汰後的一致性
	testText := texts[0]
	tokens1, _ := calculator.CalculateTokens(testText, "estimation")
	tokens2, _ := calculator.CalculateTokens(testText, "estimation")
	
	if tokens1 != tokens2 {
		t.Errorf("快取淘汰後結果不一致: %d vs %d", tokens1, tokens2)
	}
}

// TestMethodSwitchingPerformance 測試方法切換效能
func TestMethodSwitchingPerformance(t *testing.T) {
	calculator := NewTokenCalculator(1000)
	
	if !calculator.IsTiktokenAvailable() {
		t.Skip("需要 tiktoken 進行方法切換效能測試")
	}

	testTexts := []string{
		"短文本",
		strings.Repeat("中等長度文本 ", 50),
		strings.Repeat("較長的文本內容 ", 200),
	}

	methods := []string{"estimation", "tiktoken", "auto"}

	for _, text := range testTexts {
		t.Run(fmt.Sprintf("TextLength_%d", len(text)), func(t *testing.T) {
			results := make(map[string]time.Duration)
			
			for _, method := range methods {
				iterations := 1000
				start := time.Now()
				
				for i := 0; i < iterations; i++ {
					_, err := calculator.CalculateTokens(text, method)
					if err != nil {
						t.Errorf("方法 %s 失敗: %v", method, err)
						break
					}
				}
				
				duration := time.Since(start)
				results[method] = duration
				
				avgTime := duration / time.Duration(iterations)
				t.Logf("方法 %s: 總時間 %v, 平均 %v", method, duration, avgTime)
			}

			// 比較方法效能
			estimationTime := results["estimation"]
			tiktokenTime := results["tiktoken"]
			autoTime := results["auto"]

			t.Logf("效能比較 (文本長度: %d):", len(text))
			t.Logf("  估算方法: %v", estimationTime)
			t.Logf("  Tiktoken: %v", tiktokenTime)
			t.Logf("  自動選擇: %v", autoTime)

			// 自動方法的效能開銷應該很小
			if autoTime > tiktokenTime*110/100 { // 允許 10% 的開銷
				t.Logf("警告: 自動方法開銷較大: %v vs %v", autoTime, tiktokenTime)
			}
		})
	}
}

// BenchmarkTokenCalculationSizes 不同文本大小的基準測試
func BenchmarkTokenCalculationSizes(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	
	benchmarkSizes := []struct {
		name string
		size int
	}{
		{"XSmall", 10},
		{"Small", 100},
		{"Medium", 1000},
		{"Large", 10000},
		{"XLarge", 100000},
	}

	baseText := "基準測試文本 Benchmark test text "
	
	for _, size := range benchmarkSizes {
		repeatCount := size.size / len(baseText)
		if repeatCount == 0 {
			repeatCount = 1
		}
		text := strings.Repeat(baseText, repeatCount)
		
		b.Run(size.name, func(b *testing.B) {
			b.SetBytes(int64(len(text)))
			b.ResetTimer()
			
			for i := 0; i < b.N; i++ {
				calculator.CalculateTokens(text, "estimation")
			}
		})
	}
}

// BenchmarkConcurrentAccess 併發存取基準測試
func BenchmarkConcurrentAccess(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	text := strings.Repeat("併發基準測試 ", 100)
	
	// 預熱快取
	calculator.CalculateTokens(text, "estimation")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			calculator.CalculateTokens(text, "estimation")
		}
	})
}

// BenchmarkCacheOperations 快取操作基準測試
func BenchmarkCacheOperations(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	
	// 生成測試文本
	texts := make([]string, 100)
	for i := 0; i < 100; i++ {
		texts[i] = fmt.Sprintf("快取基準測試 %d", i)
	}

	b.Run("CacheMiss", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := fmt.Sprintf("unique text %d", i)
			calculator.CalculateTokens(text, "estimation")
		}
	})

	b.Run("CacheHit", func(b *testing.B) {
		// 預熱快取
		for _, text := range texts {
			calculator.CalculateTokens(text, "estimation")
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			text := texts[i%len(texts)]
			calculator.CalculateTokens(text, "estimation")
		}
	})
}

// BenchmarkMethodComparison 方法比較基準測試
func BenchmarkMethodComparison(b *testing.B) {
	calculator := NewTokenCalculator(1000)
	text := strings.Repeat("方法比較基準測試 Method comparison benchmark ", 100)

	methods := []string{"estimation"}
	if calculator.IsTiktokenAvailable() {
		methods = append(methods, "tiktoken")
	}

	for _, method := range methods {
		b.Run(method, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				calculator.CalculateTokens(text, method)
			}
		})
	}
}