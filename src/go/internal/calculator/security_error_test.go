package calculator

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestTextValidationSecurity 測試文本驗證的安全性
func TestTextValidationSecurity(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	securityTestCases := []struct {
		name        string
		text        string
		expectError bool
		description string
	}{
		{
			name:        "正常文本",
			text:        "這是正常的文本 This is normal text",
			expectError: false,
			description: "正常文本應該通過驗證",
		},
		{
			name:        "包含密碼的文本",
			text:        "用戶密碼是 password123",
			expectError: false,
			description: "包含敏感詞但不違反驗證規則",
		},
		{
			name:        "極長文本攻擊",
			text:        strings.Repeat("a", 1000001),
			expectError: true,
			description: "超長文本應該被拒絕",
		},
		{
			name:        "控制字符注入",
			text:        strings.Repeat("\x00\x01\x02\x03", 50),
			expectError: true,
			description: "過多控制字符應該被拒絕",
		},
		{
			name:        "Unicode 轟炸攻擊",
			text:        strings.Repeat("💀", 10000),
			expectError: false,
			description: "大量 Unicode 字符（但在合理範圍內）",
		},
		{
			name:        "零寬度字符攻擊",
			text:        strings.Repeat("\u200B\u200C\u200D", 1000),
			expectError: false,
			description: "零寬度字符不應該觸發安全錯誤",
		},
		{
			name:        "格式化字符串攻擊",
			text:        "%s %d %x %p %n" + strings.Repeat("%s", 100),
			expectError: false,
			description: "格式化字符串不應該引起問題",
		},
		{
			name:        "SQL 注入模式",
			text:        "'; DROP TABLE users; --",
			expectError: false,
			description: "SQL 注入模式應該被當作普通文本處理",
		},
		{
			name:        "腳本注入模式",
			text:        "<script>alert('xss')</script>",
			expectError: false,
			description: "腳本標籤應該被當作普通文本處理",
		},
		{
			name:        "路徑遍歷模式",
			text:        "../../../etc/passwd",
			expectError: false,
			description: "路徑遍歷模式應該被當作普通文本處理",
		},
	}

	for _, tc := range securityTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := calculator.ValidateText(tc.text)

			if tc.expectError && err == nil {
				t.Errorf("預期安全驗證錯誤但沒有發生: %s", tc.description)
			}

			if !tc.expectError && err != nil {
				t.Errorf("意外的安全驗證錯誤: %v - %s", err, tc.description)
			}

			// 如果驗證通過，測試計算是否正常
			if err == nil {
				tokens, calcErr := calculator.CalculateTokens(tc.text, "estimation")
				if calcErr != nil {
					t.Errorf("計算失敗: %v", calcErr)
				} else if len(tc.text) > 0 && tokens <= 0 {
					t.Errorf("非空文本應該產生正數 token: %d", tokens)
				}
			}

			t.Logf("%s: 驗證%s - %s", 
				tc.name, 
				map[bool]string{true: "失敗", false: "通過"}[err != nil],
				tc.description)
		})
	}
}

// TestMemoryLimitsAndProtection 測試記憶體限制和保護
func TestMemoryLimitsAndProtection(t *testing.T) {
	// 記錄初始記憶體狀態
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	calculator := NewTokenCalculator(1000)

	memoryTestCases := []struct {
		name        string
		textSize    int
		iterations  int
		description string
	}{
		{
			name:        "小文本大量計算",
			textSize:    100,
			iterations:  10000,
			description: "測試小文本的記憶體效率",
		},
		{
			name:        "中等文本適量計算",
			textSize:    10000,
			iterations:  1000,
			description: "測試中等大小文本的記憶體使用",
		},
		{
			name:        "大文本少量計算",
			textSize:    100000,
			iterations:  100,
			description: "測試大文本的記憶體處理",
		},
	}

	for _, tc := range memoryTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// 生成測試文本
			baseText := "測試記憶體使用 Memory usage test "
			text := strings.Repeat(baseText, tc.textSize/len(baseText))

			runtime.GC()
			var memBefore runtime.MemStats
			runtime.ReadMemStats(&memBefore)

			// 執行計算
			start := time.Now()
			for i := 0; i < tc.iterations; i++ {
				_, err := calculator.CalculateTokens(text, "estimation")
				if err != nil {
					t.Errorf("計算 %d 失敗: %v", i, err)
					break
				}

				// 定期檢查記憶體使用
				if i%1000 == 0 && i > 0 {
					var memCurrent runtime.MemStats
					runtime.ReadMemStats(&memCurrent)
					currentUsage := memCurrent.Alloc - memBefore.Alloc
					
					// 如果記憶體使用超過 100MB，發出警告
					if currentUsage > 100*1024*1024 {
						t.Logf("警告: 記憶體使用較高: %d MB (迭代 %d)", 
							currentUsage/(1024*1024), i)
					}
				}
			}
			duration := time.Since(start)

			runtime.GC()
			var memAfter runtime.MemStats
			runtime.ReadMemStats(&memAfter)

			memoryIncrease := memAfter.Alloc - memBefore.Alloc
			
			t.Logf("%s: %d 次計算耗時 %v", tc.name, tc.iterations, duration)
			t.Logf("  記憶體增加: %d KB", memoryIncrease/1024)
			t.Logf("  平均每次: %d bytes", memoryIncrease/uint64(tc.iterations))

			// 檢查記憶體洩漏的軟性指標
			if memoryIncrease > uint64(tc.iterations)*1000 { // 每次迭代超過 1KB 可能有問題
				t.Logf("警告: 可能存在記憶體洩漏，每次迭代平均使用 %d bytes", 
					memoryIncrease/uint64(tc.iterations))
			}
		})
	}

	// 檢查總體記憶體增長
	runtime.GC()
	runtime.ReadMemStats(&m2)
	totalIncrease := m2.Alloc - m1.Alloc
	t.Logf("總記憶體增加: %d KB", totalIncrease/1024)
}

// TestConcurrentSafetyAndRaceConditions 測試併發安全性和競爭條件
func TestConcurrentSafetyAndRaceConditions(t *testing.T) {
	calculator := NewTokenCalculator(100)
	
	// 測試併發寫入快取
	t.Run("ConcurrentCacheWrites", func(t *testing.T) {
		var wg sync.WaitGroup
		errors := make(chan error, 100)
		
		// 啟動多個 goroutine 同時寫入不同的快取項目
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < 50; j++ {
					text := fmt.Sprintf("併發測試文本 %d-%d", id, j)
					_, err := calculator.CalculateTokens(text, "estimation")
					if err != nil {
						errors <- fmt.Errorf("goroutine %d: %w", id, err)
						return
					}
				}
			}(i)
		}

		// 等待完成
		done := make(chan bool)
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// 成功完成
		case err := <-errors:
			t.Errorf("併發寫入錯誤: %v", err)
		case <-time.After(10 * time.Second):
			t.Error("併發測試超時")
		}

		close(errors)
		for err := range errors {
			t.Errorf("併發錯誤: %v", err)
		}
	})

	// 測試併發讀寫混合
	t.Run("ConcurrentReadWrites", func(t *testing.T) {
		sharedText := "共享的測試文本 Shared test text"
		
		// 預先計算一次
		calculator.CalculateTokens(sharedText, "estimation")

		var wg sync.WaitGroup
		errors := make(chan error, 100)

		// 啟動讀取 goroutine
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < 100; j++ {
					_, err := calculator.CalculateTokens(sharedText, "estimation")
					if err != nil {
						errors <- fmt.Errorf("reader %d: %w", id, err)
						return
					}
				}
			}(i)
		}

		// 啟動寫入 goroutine
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				
				for j := 0; j < 50; j++ {
					text := fmt.Sprintf("寫入測試 %d-%d", id, j)
					_, err := calculator.CalculateTokens(text, "estimation")
					if err != nil {
						errors <- fmt.Errorf("writer %d: %w", id, err)
						return
					}
				}
			}(i)
		}

		// 等待完成
		done := make(chan bool)
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// 成功完成
		case err := <-errors:
			t.Errorf("併發讀寫錯誤: %v", err)
		case <-time.After(10 * time.Second):
			t.Error("併發讀寫測試超時")
		}

		close(errors)
		for err := range errors {
			t.Errorf("併發錯誤: %v", err)
		}
	})
}

// TestSensitiveInformationProtection 測試敏感資訊保護
func TestSensitiveInformationProtection(t *testing.T) {
	calculator := NewTokenCalculator(100)

	sensitivePatterns := []struct {
		name        string
		text        string
		description string
	}{
		{
			name:        "信用卡號",
			text:        "信用卡號碼: 4111-1111-1111-1111",
			description: "信用卡號碼應該被正常處理但不記錄到日誌",
		},
		{
			name:        "身分證號",
			text:        "身分證字號: A123456789",
			description: "身分證號碼",
		},
		{
			name:        "電話號碼",
			text:        "聯絡電話: 02-2345-6789",
			description: "電話號碼",
		},
		{
			name:        "電子郵件",
			text:        "email: user@example.com",
			description: "電子郵件地址",
		},
		{
			name:        "API密鑰",
			text:        "API_KEY=sk-1234567890abcdef",
			description: "API 密鑰",
		},
		{
			name:        "密碼",
			text:        "password: secretPassword123!",
			description: "密碼",
		},
		{
			name:        "JWT Token",
			text:        "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			description: "JWT Token",
		},
	}

	for _, tc := range sensitivePatterns {
		t.Run(tc.name, func(t *testing.T) {
			// 測試基本計算功能
			tokens, err := calculator.CalculateTokens(tc.text, "estimation")
			if err != nil {
				t.Errorf("計算失敗: %v", err)
				return
			}

			if tokens <= 0 {
				t.Errorf("應該產生正數 token: %d", tokens)
			}

			// 測試文本驗證
			if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
				err := calcImpl.ValidateText(tc.text)
				if err != nil {
					t.Errorf("文本驗證失敗: %v", err)
				}
			}

			// 重要：在實際記錄中不顯示敏感內容
			t.Logf("%s: %d tokens (文本長度: %d) - %s", 
				tc.name, tokens, len(tc.text), tc.description)
			// 注意：這裡故意不記錄實際文本內容
		})
	}
}

// TestErrorRecoveryMechanisms 測試錯誤恢復機制
func TestErrorRecoveryMechanisms(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	// 測試從無效狀態恢復
	t.Run("RecoveryFromInvalidState", func(t *testing.T) {
		// 模擬破壞快取狀態
		calculator.cache = nil
		
		// 嘗試計算，應該能夠恢復
		text := "恢復測試文本 Recovery test text"
		tokens, err := calculator.CalculateTokens(text, "estimation")
		
		if err != nil {
			t.Logf("從無效狀態恢復失敗: %v", err)
			// 嘗試重新初始化
			calculator.cache = make(map[string]int)
			tokens, err = calculator.CalculateTokens(text, "estimation")
			if err != nil {
				t.Errorf("重新初始化後仍然失敗: %v", err)
			}
		}

		if tokens <= 0 {
			t.Errorf("恢復後應該產生正數 token: %d", tokens)
		}
	})

	// 測試從 tiktoken 錯誤恢復
	t.Run("RecoveryFromTiktokenError", func(t *testing.T) {
		text := "Tiktoken 錯誤恢復測試"
		
		// 保存原始狀態
		originalEnabled := calculator.tiktokenEnabled
		originalEncoder := calculator.tiktokenEncoder
		
		// 模擬 tiktoken 錯誤狀態
		calculator.tiktokenEnabled = true
		calculator.tiktokenEncoder = nil
		
		// 嘗試使用 tiktoken 計算
		tokens, err := calculator.CalculateTokens(text, "tiktoken")
		if err != nil {
			t.Logf("Tiktoken 錯誤狀態正確處理: %v", err)
		} else {
			t.Logf("從 tiktoken 錯誤狀態恢復，使用回退方法: %d tokens", tokens)
		}
		
		// 恢復原始狀態
		calculator.tiktokenEnabled = originalEnabled
		calculator.tiktokenEncoder = originalEncoder
	})

	// 測試記憶體壓力下的錯誤恢復
	t.Run("RecoveryUnderMemoryPressure", func(t *testing.T) {
		// 創建記憶體壓力
		pressureData := make([][]byte, 0)
		
		defer func() {
			// 清理記憶體壓力
			pressureData = nil
			runtime.GC()
		}()

		// 逐漸增加記憶體壓力並測試恢復能力
		for i := 0; i < 10; i++ {
			// 添加記憶體壓力
			pressureData = append(pressureData, make([]byte, 1024*1024)) // 1MB
			
			text := fmt.Sprintf("記憶體壓力測試 %d", i)
			tokens, err := calculator.CalculateTokens(text, "estimation")
			
			if err != nil {
				t.Errorf("記憶體壓力下計算失敗 (iteration %d): %v", i, err)
			} else if tokens <= 0 {
				t.Errorf("記憶體壓力下應該產生正數 token (iteration %d): %d", i, tokens)
			}
		}
	})
}

// TestInputSanitizationAndValidation 測試輸入清理和驗證
func TestInputSanitizationAndValidation(t *testing.T) {
	calculator := NewTokenCalculator(100).(*TokenCalculatorImpl)

	maliciousInputs := []struct {
		name        string
		text        string
		expectError bool
		description string
	}{
		{
			name:        "極端長度字符串",
			text:        strings.Repeat("a", 2000000),
			expectError: true,
			description: "超長字符串應該被拒絕",
		},
		{
			name:        "二進位數據",
			text:        string([]byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}),
			expectError: false,
			description: "二進位數據應該被處理（但可能產生警告）",
		},
		{
			name:        "控制字符轟炸",
			text:        strings.Repeat("\x00", 1000),
			expectError: true,
			description: "大量控制字符應該被拒絕",
		},
		{
			name:        "Unicode 正規化攻擊",
			text:        "é" + "é", // 一個是組合字符，一個是預組合字符
			expectError: false,
			description: "Unicode 正規化不同但視覺相同的字符",
		},
		{
			name:        "巨大的 Unicode 代碼點",
			text:        string(rune(0x10FFFF)), // 最大的有效 Unicode 代碼點
			expectError: false,
			description: "最大 Unicode 代碼點",
		},
	}

	for _, tc := range maliciousInputs {
		t.Run(tc.name, func(t *testing.T) {
			// 測試驗證
			err := calculator.ValidateText(tc.text)
			
			if tc.expectError && err == nil {
				t.Errorf("預期驗證錯誤但沒有發生: %s", tc.description)
			}
			
			if !tc.expectError && err != nil {
				t.Logf("驗證警告 (預期內): %v - %s", err, tc.description)
			}

			// 如果驗證通過，測試計算
			if err == nil {
				tokens, calcErr := calculator.CalculateTokens(tc.text, "estimation")
				if calcErr != nil {
					t.Logf("計算錯誤: %v", calcErr)
				} else {
					t.Logf("%s: %d tokens - %s", tc.name, tokens, tc.description)
				}
			}
		})
	}
}

// TestSecurityAuditLog 測試安全審計日誌（模擬）
func TestSecurityAuditLog(t *testing.T) {
	calculator := NewTokenCalculator(100)

	// 模擬需要審計的操作
	auditEvents := []struct {
		operation   string
		text        string
		shouldAudit bool
	}{
		{"normal_calculation", "正常計算文本", false},
		{"large_text", strings.Repeat("a", 10000), true},
		{"sensitive_content", "password: secret123", true},
		{"control_chars", "text\x00\x01text", true},
		{"unicode_heavy", strings.Repeat("🔥", 1000), true},
	}

	for _, event := range auditEvents {
		t.Run(event.operation, func(t *testing.T) {
			start := time.Now()
			
			tokens, err := calculator.CalculateTokens(event.text, "estimation")
			
			duration := time.Since(start)

			// 模擬安全審計日誌記錄
			if event.shouldAudit {
				auditInfo := map[string]interface{}{
					"timestamp":   start.Unix(),
					"operation":   event.operation,
					"text_length": len(event.text),
					"duration_ms": duration.Milliseconds(),
					"tokens":      tokens,
					"error":       err != nil,
				}
				
				t.Logf("安全審計: %+v", auditInfo)
			}

			if err != nil {
				t.Logf("操作 %s 發生錯誤: %v", event.operation, err)
			} else {
				t.Logf("操作 %s 完成: %d tokens", event.operation, tokens)
			}
		})
	}
}

// TestDenialOfServiceProtection 測試拒絕服務攻擊防護
func TestDenialOfServiceProtection(t *testing.T) {
	calculator := NewTokenCalculator(50) // 小快取以便測試

	// 測試快取淹沒攻擊
	t.Run("CacheFloodingProtection", func(t *testing.T) {
		start := time.Now()
		
		// 嘗試用大量不同的文本淹沒快取
		for i := 0; i < 1000; i++ {
			text := fmt.Sprintf("快取淹沒測試 %d %s", i, strings.Repeat("x", i%100))
			_, err := calculator.CalculateTokens(text, "estimation")
			
			if err != nil {
				t.Errorf("計算 %d 失敗: %v", i, err)
				break
			}
			
			// 檢查是否花費過長時間
			if time.Since(start) > 30*time.Second {
				t.Logf("快取淹沒測試達到時間限制，停止在 %d 次迭代", i)
				break
			}
		}
		
		duration := time.Since(start)
		t.Logf("快取淹沒測試完成，耗時: %v", duration)
		
		// 檢查快取是否正常運作
		if calcImpl, ok := calculator.(*TokenCalculatorImpl); ok {
			stats := calcImpl.GetCacheStats()
			t.Logf("最終快取統計: %+v", stats)
		}
	})

	// 測試計算複雜度攻擊
	t.Run("ComputationalComplexityProtection", func(t *testing.T) {
		complexTexts := []string{
			strings.Repeat("複雜計算測試 ", 10000),
			strings.Repeat("🚀🎉🔥💯✨", 5000),
			strings.Repeat("測試 test テスト 테스트 ", 8000),
		}

		for i, text := range complexTexts {
			start := time.Now()
			
			tokens, err := calculator.CalculateTokens(text, "estimation")
			
			duration := time.Since(start)
			
			if err != nil {
				t.Errorf("複雜文本 %d 計算失敗: %v", i, err)
				continue
			}
			
			t.Logf("複雜文本 %d: %d tokens, 耗時: %v", i, tokens, duration)
			
			// 檢查計算時間是否合理
			if duration > 5*time.Second {
				t.Logf("警告: 複雜文本 %d 計算時間過長: %v", i, duration)
			}
		}
	})
}