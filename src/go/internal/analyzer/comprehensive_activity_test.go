package analyzer

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestActivityClassificationAccuracy 測試活動分類的準確性
func TestActivityClassificationAccuracy(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	testCases := []struct {
		name         string
		content      string
		expectedType types.ActivityType
		description  string
	}{
		{
			name:         "Coding_Chinese",
			content:      "請幫我實作一個計算 Token 的函數",
			expectedType: types.ActivityCoding,
			description:  "中文編程請求",
		},
		{
			name:         "Coding_English",
			content:      "I need to implement a function that calculates tokens",
			expectedType: types.ActivityCoding,
			description:  "英文編程請求",
		},
		{
			name:         "Debugging_Chinese",
			content:      "這個程式有錯誤，請幫我修復這個問題",
			expectedType: types.ActivityDebugging,
			description:  "中文除錯請求",
		},
		{
			name:         "Debugging_English",
			content:      "There's a bug in my code, can you help me fix this error?",
			expectedType: types.ActivityDebugging,
			description:  "英文除錯請求",
		},
		{
			name:         "Documentation_Chinese",
			content:      "請幫我更新 README 文件和撰寫 API 說明",
			expectedType: types.ActivityDocumentation,
			description:  "中文文件請求",
		},
		{
			name:         "Documentation_English",
			content:      "I need to write documentation for this API",
			expectedType: types.ActivityDocumentation,
			description:  "英文文件請求",
		},
		{
			name:         "SpecDevelopment_Chinese",
			content:      "我需要設計一個新的系統架構和建立需求規格",
			expectedType: types.ActivitySpecDev,
			description:  "中文規格開發請求",
		},
		{
			name:         "SpecDevelopment_English",
			content:      "I need to create a system design and requirement specification",
			expectedType: types.ActivitySpecDev,
			description:  "英文規格開發請求",
		},
		{
			name:         "Chat_Chinese",
			content:      "你好，我想聊聊天氣如何",
			expectedType: types.ActivityTypeChat,
			description:  "中文聊天請求",
		},
		{
			name:         "Chat_English",
			content:      "Hello, I have a question about programming",
			expectedType: types.ActivityTypeChat,
			description:  "英文聊天請求",
		},
		{
			name:         "Empty_Content",
			content:      "",
			expectedType: types.ActivityTypeChat,
			description:  "空內容預設為聊天",
		},
		{
			name:         "Mixed_Content",
			content:      "請幫我實作一個函數來修復這個 bug 並更新文件",
			expectedType: types.ActivityDocumentation, // 文件關鍵字權重較高
			description:  "混合內容測試",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzer.ClassifyActivity(tc.content)
			if result != tc.expectedType {
				t.Errorf("測試 %s 失敗: 期望 %v, 得到 %v (內容: %s)",
					tc.description, tc.expectedType, result, tc.content)
			}
		})
	}
}

// TestBatchActivityAnalysis 測試批次活動分析
func TestBatchActivityAnalysis(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	contents := []string{
		"實作一個新功能",
		"修復這個錯誤",
		"更新文件",
		"設計系統架構",
		"一般聊天內容",
	}

	expected := []types.ActivityType{
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
		types.ActivitySpecDev,
		types.ActivityTypeChat,
	}

	results := analyzer.AnalyzeActivityBatch(contents)

	if len(results) != len(expected) {
		t.Fatalf("結果數量不符: 期望 %d, 得到 %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("批次分析第 %d 項失敗: 期望 %v, 得到 %v", i, expected[i], result)
		}
	}
}

// TestActivitySummaryGeneration 測試活動摘要生成
func TestActivitySummaryGeneration(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	// 建立測試活動數據
	activities := []types.Activity{
		{
			ID:        "1",
			Type:      types.ActivityCoding,
			Content:   "實作功能",
			Tokens:    types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300},
			StartTime: time.Now().Add(-10 * time.Minute),
			EndTime:   time.Now().Add(-5 * time.Minute),
		},
		{
			ID:        "2",
			Type:      types.ActivityCoding,
			Content:   "另一個功能",
			Tokens:    types.TokenUsage{InputTokens: 150, OutputTokens: 250, TotalTokens: 400},
			StartTime: time.Now().Add(-8 * time.Minute),
			EndTime:   time.Now().Add(-3 * time.Minute),
		},
		{
			ID:        "3",
			Type:      types.ActivityDebugging,
			Content:   "修復錯誤",
			Tokens:    types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200},
			StartTime: time.Now().Add(-6 * time.Minute),
			EndTime:   time.Now().Add(-2 * time.Minute),
		},
	}

	summary := analyzer.GenerateActivitySummary(activities)

	// 驗證基本統計
	if summary.TotalActivities != 3 {
		t.Errorf("總活動數錯誤: 期望 3, 得到 %d", summary.TotalActivities)
	}

	// 驗證活動類型計數
	if summary.ActivityCounts[types.ActivityCoding] != 2 {
		t.Errorf("編程活動計數錯誤: 期望 2, 得到 %d", summary.ActivityCounts[types.ActivityCoding])
	}

	if summary.ActivityCounts[types.ActivityDebugging] != 1 {
		t.Errorf("除錯活動計數錯誤: 期望 1, 得到 %d", summary.ActivityCounts[types.ActivityDebugging])
	}

	// 驗證 Token 使用統計
	codingTokens := summary.TokenUsage[types.ActivityCoding]
	expectedCodingTokens := types.TokenUsage{InputTokens: 250, OutputTokens: 450, TotalTokens: 700}

	if codingTokens != expectedCodingTokens {
		t.Errorf("編程活動 Token 統計錯誤: 期望 %+v, 得到 %+v", expectedCodingTokens, codingTokens)
	}

	// 驗證總 Token 統計
	expectedTotalTokens := types.TokenUsage{InputTokens: 330, OutputTokens: 570, TotalTokens: 900}
	if summary.TotalTokens != expectedTotalTokens {
		t.Errorf("總 Token 統計錯誤: 期望 %+v, 得到 %+v", expectedTotalTokens, summary.TotalTokens)
	}
}

// TestEfficiencyMetricsCalculation 測試效率指標計算
func TestEfficiencyMetricsCalculation(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	// 建立測試數據
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{TotalTokens: 300},
			StartTime: time.Now().Add(-10 * time.Minute),
			EndTime:   time.Now().Add(-5 * time.Minute), // 5 分鐘
		},
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{TotalTokens: 600},
			StartTime: time.Now().Add(-8 * time.Minute),
			EndTime:   time.Now().Add(-3 * time.Minute), // 5 分鐘
		},
	}

	data := types.ActivityData{
		ActivitiesByType: map[types.ActivityType][]types.Activity{
			types.ActivityCoding: activities,
		},
	}

	metrics := analyzer.CalculateEfficiencyMetrics(data)

	// 驗證平均 Token 數
	expectedAvgTokens := 450.0 // (300 + 600) / 2
	if metrics.AverageTokensPerActivity[types.ActivityCoding] != expectedAvgTokens {
		t.Errorf("平均 Token 數錯誤: 期望 %.1f, 得到 %.1f",
			expectedAvgTokens, metrics.AverageTokensPerActivity[types.ActivityCoding])
	}

	// 驗證平均時間
	expectedAvgTime := 5 * time.Minute
	if metrics.AverageTimePerActivity[types.ActivityCoding] != expectedAvgTime {
		t.Errorf("平均時間錯誤: 期望 %v, 得到 %v",
			expectedAvgTime, metrics.AverageTimePerActivity[types.ActivityCoding])
	}

	// 驗證每分鐘 Token 數
	expectedTokensPerMin := 90.0 // 900 tokens / 10 minutes
	if metrics.TokensPerMinute[types.ActivityCoding] != expectedTokensPerMin {
		t.Errorf("每分鐘 Token 數錯誤: 期望 %.1f, 得到 %.1f",
			expectedTokensPerMin, metrics.TokensPerMinute[types.ActivityCoding])
	}
}

// TestActivityTypeDistribution 測試活動類型分佈計算
func TestActivityTypeDistribution(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	activities := []types.Activity{
		{Type: types.ActivityCoding},
		{Type: types.ActivityCoding},
		{Type: types.ActivityCoding},
		{Type: types.ActivityDebugging},
		{Type: types.ActivityTypeChat},
	}

	distribution := analyzer.GetActivityTypeDistribution(activities)

	// 驗證分佈百分比
	expectedCoding := 60.0    // 3/5 * 100
	expectedDebugging := 20.0 // 1/5 * 100
	expectedChat := 20.0      // 1/5 * 100

	if distribution[types.ActivityCoding] != expectedCoding {
		t.Errorf("編程活動分佈錯誤: 期望 %.1f%%, 得到 %.1f%%",
			expectedCoding, distribution[types.ActivityCoding])
	}

	if distribution[types.ActivityDebugging] != expectedDebugging {
		t.Errorf("除錯活動分佈錯誤: 期望 %.1f%%, 得到 %.1f%%",
			expectedDebugging, distribution[types.ActivityDebugging])
	}

	if distribution[types.ActivityTypeChat] != expectedChat {
		t.Errorf("聊天活動分佈錯誤: 期望 %.1f%%, 得到 %.1f%%",
			expectedChat, distribution[types.ActivityTypeChat])
	}
}

// TestEdgeCasesAndErrorHandling 測試邊界情況和錯誤處理
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	t.Run("空活動列表", func(t *testing.T) {
		summary := analyzer.GenerateActivitySummary([]types.Activity{})
		if summary.TotalActivities != 0 {
			t.Errorf("空列表總活動數應為 0, 得到 %d", summary.TotalActivities)
		}

		distribution := analyzer.GetActivityTypeDistribution([]types.Activity{})
		if len(distribution) != 0 {
			t.Errorf("空列表分佈應為空, 得到 %d 項", len(distribution))
		}
	})

	t.Run("無效時間戳記", func(t *testing.T) {
		activities := []types.Activity{
			{
				Type:      types.ActivityCoding,
				Tokens:    types.TokenUsage{TotalTokens: 100},
				StartTime: time.Time{}, // 零值時間
				EndTime:   time.Time{}, // 零值時間
			},
		}

		summary := analyzer.GenerateActivitySummary(activities)
		// 應該能正常處理，但時間統計為 0
		if summary.TotalActivities != 1 {
			t.Errorf("應該能處理無效時間戳記的活動")
		}
	})

	t.Run("極長內容", func(t *testing.T) {
		longContent := ""
		for i := 0; i < 10000; i++ {
			longContent += "這是一個很長的內容 "
		}

		// 應該能正常分類，不會崩潰
		result := analyzer.ClassifyActivity(longContent)
		if result == "" {
			t.Errorf("極長內容應該能正常分類")
		}
	})

	t.Run("特殊字符內容", func(t *testing.T) {
		specialContent := "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"
		result := analyzer.ClassifyActivity(specialContent)
		// 應該預設為聊天類型
		if result != types.ActivityTypeChat {
			t.Errorf("特殊字符內容應預設為聊天類型, 得到 %v", result)
		}
	})

	t.Run("Unicode 內容", func(t *testing.T) {
		unicodeContent := "🚀 實作一個新功能 💻"
		result := analyzer.ClassifyActivity(unicodeContent)
		// 應該能正確識別為編程類型
		if result != types.ActivityCoding {
			t.Errorf("Unicode 內容應能正確分類, 期望 %v, 得到 %v",
				types.ActivityCoding, result)
		}
	})
}

// TestPatternMatchingRobustness 測試模式匹配的穩健性
func TestPatternMatchingRobustness(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	testCases := []struct {
		name     string
		content  string
		expected types.ActivityType
	}{
		{
			name:     "大小寫混合",
			content:  "請幫我 IMPLEMENT 一個 Function",
			expected: types.ActivityCoding,
		},
		{
			name:     "中英混合",
			content:  "I need to 實作 a new feature",
			expected: types.ActivityTypeChat, // 混合語言可能被識別為聊天
		},
		{
			name:     "多個關鍵字",
			content:  "請幫我實作一個函數來修復錯誤並更新文件",
			expected: types.ActivityDocumentation, // 文件關鍵字權重較高
		},
		{
			name:     "模糊匹配",
			content:  "我想要建立一個新的程式功能",
			expected: types.ActivityCoding,
		},
		{
			name:     "否定語句",
			content:  "我不知道如何實作這個功能",
			expected: types.ActivityCoding, // 仍然包含實作關鍵字
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzer.ClassifyActivity(tc.content)
			if result != tc.expected {
				t.Errorf("模式匹配測試 %s 失敗: 期望 %v, 得到 %v",
					tc.name, tc.expected, result)
			}
		})
	}
}

// TestPerformanceWithLargeDataset 測試大數據集的效能
func TestPerformanceWithLargeDataset(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	// 建立大量測試數據
	activityTypes := []types.ActivityType{
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
		types.ActivitySpecDev,
		types.ActivityTypeChat,
	}

	largeDataset := make([]types.Activity, 10000)
	for i := 0; i < 10000; i++ {
		largeDataset[i] = types.Activity{
			ID:        string(rune(i)),
			Type:      activityTypes[i%5], // 循環使用 5 種類型
			Content:   "測試內容",
			Tokens:    types.TokenUsage{TotalTokens: i + 1},
			Timestamp: time.Now().Add(-time.Duration(i) * time.Minute),
		}
	}

	start := time.Now()

	// 測試摘要生成效能
	summary := analyzer.GenerateActivitySummary(largeDataset)

	elapsed := time.Since(start)

	// 驗證結果正確性
	if summary.TotalActivities != 10000 {
		t.Errorf("大數據集處理錯誤: 期望 10000 活動, 得到 %d", summary.TotalActivities)
	}

	// 效能要求：處理 10000 個活動應在 1 秒內完成
	if elapsed > time.Second {
		t.Errorf("效能測試失敗: 處理 10000 個活動耗時 %v, 超過 1 秒限制", elapsed)
	}

	t.Logf("大數據集處理效能: %v (10000 個活動)", elapsed)
}
