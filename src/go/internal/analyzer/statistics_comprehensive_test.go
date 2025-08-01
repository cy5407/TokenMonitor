package analyzer

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestTokenUsageCalculation 測試 Token 使用量計算的準確性
func TestTokenUsageCalculation(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	activities := []types.Activity{
		{
			Type:   types.ActivityCoding,
			Tokens: types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300},
		},
		{
			Type:   types.ActivityCoding,
			Tokens: types.TokenUsage{InputTokens: 150, OutputTokens: 250, TotalTokens: 400},
		},
		{
			Type:   types.ActivityDebugging,
			Tokens: types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200},
		},
	}

	usage := stats.CalculateTokenUsageByActivity(activities)

	// 驗證編程活動的 Token 使用量
	codingUsage := usage[types.ActivityCoding]
	expectedCoding := types.TokenUsage{InputTokens: 250, OutputTokens: 450, TotalTokens: 700}
	if codingUsage != expectedCoding {
		t.Errorf("編程活動 Token 使用量錯誤: 期望 %+v, 得到 %+v", expectedCoding, codingUsage)
	}

	// 驗證除錯活動的 Token 使用量
	debuggingUsage := usage[types.ActivityDebugging]
	expectedDebugging := types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200}
	if debuggingUsage != expectedDebugging {
		t.Errorf("除錯活動 Token 使用量錯誤: 期望 %+v, 得到 %+v", expectedDebugging, debuggingUsage)
	}
}

// TestActivityTotalsCalculation 測試活動總和統計功能
func TestActivityTotalsCalculation(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	baseTime := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300},
			StartTime: baseTime,
			EndTime:   baseTime.Add(10 * time.Minute),
		},
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{InputTokens: 150, OutputTokens: 250, TotalTokens: 400},
			StartTime: baseTime.Add(15 * time.Minute),
			EndTime:   baseTime.Add(25 * time.Minute),
		},
		{
			Type:      types.ActivityDebugging,
			Tokens:    types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200},
			StartTime: baseTime.Add(30 * time.Minute),
			EndTime:   baseTime.Add(35 * time.Minute),
		},
	}

	totals := stats.CalculateActivityTotals(activities)

	// 驗證總活動數
	if totals.TotalActivities != 3 {
		t.Errorf("總活動數錯誤: 期望 3, 得到 %d", totals.TotalActivities)
	}

	// 驗證總 Token 使用量
	expectedTotalTokens := types.TokenUsage{InputTokens: 330, OutputTokens: 570, TotalTokens: 900}
	if totals.TotalTokens != expectedTotalTokens {
		t.Errorf("總 Token 使用量錯誤: 期望 %+v, 得到 %+v", expectedTotalTokens, totals.TotalTokens)
	}

	// 驗證總時間
	expectedTotalTime := 25 * time.Minute // 10 + 10 + 5
	if totals.TotalTime != expectedTotalTime {
		t.Errorf("總時間錯誤: 期望 %v, 得到 %v", expectedTotalTime, totals.TotalTime)
	}

	// 驗證各類型統計
	codingTotal := totals.ByType[types.ActivityCoding]
	if codingTotal.Count != 2 {
		t.Errorf("編程活動計數錯誤: 期望 2, 得到 %d", codingTotal.Count)
	}

	expectedCodingTokens := types.TokenUsage{InputTokens: 250, OutputTokens: 450, TotalTokens: 700}
	if codingTotal.Tokens != expectedCodingTokens {
		t.Errorf("編程活動 Token 統計錯誤: 期望 %+v, 得到 %+v", expectedCodingTokens, codingTotal.Tokens)
	}

	expectedCodingTime := 20 * time.Minute
	if codingTotal.TotalTime != expectedCodingTime {
		t.Errorf("編程活動時間統計錯誤: 期望 %v, 得到 %v", expectedCodingTime, codingTotal.TotalTime)
	}
}

// TestUsagePatternsAnalysis 測試使用模式分析
func TestUsagePatternsAnalysis(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	baseTime := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{TotalTokens: 300},
			StartTime: baseTime,
			EndTime:   baseTime.Add(10 * time.Minute),
		},
		{
			Type:      types.ActivityCoding,
			Tokens:    types.TokenUsage{TotalTokens: 500},
			StartTime: baseTime.Add(15 * time.Minute),
			EndTime:   baseTime.Add(20 * time.Minute),
		},
		{
			Type:      types.ActivityDebugging,
			Tokens:    types.TokenUsage{TotalTokens: 200},
			StartTime: baseTime.Add(25 * time.Minute),
			EndTime:   baseTime.Add(35 * time.Minute),
		},
	}

	analysis := stats.AnalyzeUsagePatterns(activities)

	// 驗證基本統計
	if analysis.TotalActivities != 3 {
		t.Errorf("分析的總活動數錯誤: 期望 3, 得到 %d", analysis.TotalActivities)
	}

	// 驗證編程活動模式
	codingPattern := analysis.Patterns["coding"]
	if codingPattern.Count != 2 {
		t.Errorf("編程活動模式計數錯誤: 期望 2, 得到 %d", codingPattern.Count)
	}

	if codingPattern.TotalTokens != 800 {
		t.Errorf("編程活動總 Token 錯誤: 期望 800, 得到 %d", codingPattern.TotalTokens)
	}

	expectedAvgTokens := 400.0 // (300 + 500) / 2
	if codingPattern.AverageTokens != expectedAvgTokens {
		t.Errorf("編程活動平均 Token 錯誤: 期望 %.1f, 得到 %.1f", expectedAvgTokens, codingPattern.AverageTokens)
	}

	expectedAvgTime := 7*time.Minute + 30*time.Second // (10 + 5) / 2 minutes
	if codingPattern.AverageTime != expectedAvgTime {
		t.Errorf("編程活動平均時間錯誤: 期望 %v, 得到 %v", expectedAvgTime, codingPattern.AverageTime)
	}

	// 驗證效率指標
	expectedTokensPerMin := 800.0 / 15.0 // 800 tokens / 15 minutes
	tolerance := 0.1
	if codingPattern.TokensPerMinute < expectedTokensPerMin-tolerance ||
		codingPattern.TokensPerMinute > expectedTokensPerMin+tolerance {
		t.Errorf("編程活動每分鐘 Token 數錯誤: 期望約 %.2f, 得到 %.2f",
			expectedTokensPerMin, codingPattern.TokensPerMinute)
	}

	// 驗證洞察生成
	if len(analysis.Insights) == 0 {
		t.Errorf("應該生成使用模式洞察")
	}
}

// TestTopActivitiesByTokens 測試按 Token 使用量排序功能
func TestTopActivitiesByTokens(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	activities := []types.Activity{
		{ID: "1", Tokens: types.TokenUsage{TotalTokens: 100}},
		{ID: "2", Tokens: types.TokenUsage{TotalTokens: 500}},
		{ID: "3", Tokens: types.TokenUsage{TotalTokens: 300}},
		{ID: "4", Tokens: types.TokenUsage{TotalTokens: 200}},
		{ID: "5", Tokens: types.TokenUsage{TotalTokens: 400}},
	}

	// 測試取前 3 名
	top3 := stats.GetTopActivitiesByTokens(activities, 3)

	if len(top3) != 3 {
		t.Errorf("前 3 名數量錯誤: 期望 3, 得到 %d", len(top3))
	}

	// 驗證排序正確性
	expectedOrder := []string{"2", "5", "3"} // Token 數: 500, 400, 300
	for i, activity := range top3 {
		if activity.ID != expectedOrder[i] {
			t.Errorf("排序錯誤: 位置 %d 期望 ID %s, 得到 %s", i, expectedOrder[i], activity.ID)
		}
	}

	// 測試限制超過總數的情況
	topAll := stats.GetTopActivitiesByTokens(activities, 10)
	if len(topAll) != 5 {
		t.Errorf("超過總數限制時應返回全部: 期望 5, 得到 %d", len(topAll))
	}

	// 測試空列表
	topEmpty := stats.GetTopActivitiesByTokens([]types.Activity{}, 5)
	if len(topEmpty) != 0 {
		t.Errorf("空列表應返回空結果: 期望 0, 得到 %d", len(topEmpty))
	}
}

// TestActivityFrequencyCalculation 測試活動頻率分析
func TestActivityFrequencyCalculation(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	baseTime := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			Timestamp: baseTime,
		},
		{
			Type:      types.ActivityCoding,
			Timestamp: baseTime.Add(1 * time.Hour),
		},
		{
			Type:      types.ActivityDebugging,
			Timestamp: baseTime.Add(2 * time.Hour),
		},
		{
			Type:      types.ActivityTypeChat,
			Timestamp: baseTime.Add(3 * time.Hour),
		},
	}

	frequency := stats.CalculateActivityFrequency(activities, 1*time.Hour)

	// 驗證時間窗口
	if frequency.TimeWindow != 1*time.Hour {
		t.Errorf("時間窗口錯誤: 期望 %v, 得到 %v", 1*time.Hour, frequency.TimeWindow)
	}

	// 驗證總時間段數
	if frequency.TotalPeriods != 3 { // 3 小時 / 1 小時窗口 = 3 個時間段
		t.Errorf("總時間段數錯誤: 期望 3, 得到 %d", frequency.TotalPeriods)
	}

	// 驗證各類型頻率
	expectedCodingFreq := 2.0 / 3.0 // 2 個編程活動 / 3 個時間段
	tolerance := 0.01
	if frequency.ByType[types.ActivityCoding] < expectedCodingFreq-tolerance ||
		frequency.ByType[types.ActivityCoding] > expectedCodingFreq+tolerance {
		t.Errorf("編程活動頻率錯誤: 期望約 %.3f, 得到 %.3f",
			expectedCodingFreq, frequency.ByType[types.ActivityCoding])
	}

	// 驗證按小時分佈
	expectedHourCounts := make(map[int]int)
	for _, activity := range activities {
		hour := activity.Timestamp.Hour()
		expectedHourCounts[hour]++
	}

	for hour, expectedCount := range expectedHourCounts {
		if frequency.ByHour[hour] != expectedCount {
			t.Errorf("小時 %d 的活動計數錯誤: 期望 %d, 得到 %d",
				hour, expectedCount, frequency.ByHour[hour])
		}
	}
}

// TestStatisticsEdgeCases 測試統計功能的邊界情況
func TestStatisticsEdgeCases(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	t.Run("空活動列表", func(t *testing.T) {
		// 測試空列表的各種統計功能
		usage := stats.CalculateTokenUsageByActivity([]types.Activity{})
		if len(usage) != 0 {
			t.Errorf("空列表的 Token 使用統計應為空")
		}

		totals := stats.CalculateActivityTotals([]types.Activity{})
		if totals.TotalActivities != 0 {
			t.Errorf("空列表的總活動數應為 0")
		}

		analysis := stats.AnalyzeUsagePatterns([]types.Activity{})
		if analysis.TotalActivities != 0 {
			t.Errorf("空列表的使用模式分析總數應為 0")
		}

		frequency := stats.CalculateActivityFrequency([]types.Activity{}, 1*time.Hour)
		if len(frequency.ByType) != 0 {
			t.Errorf("空列表的頻率分析應為空")
		}
	})

	t.Run("單一活動", func(t *testing.T) {
		singleActivity := []types.Activity{
			{
				Type:      types.ActivityCoding,
				Tokens:    types.TokenUsage{TotalTokens: 100},
				StartTime: time.Now(),
				EndTime:   time.Now().Add(5 * time.Minute),
			},
		}

		totals := stats.CalculateActivityTotals(singleActivity)
		if totals.TotalActivities != 1 {
			t.Errorf("單一活動的總數應為 1")
		}

		if totals.ByType[types.ActivityCoding].Count != 1 {
			t.Errorf("單一活動的類型計數應為 1")
		}
	})

	t.Run("零 Token 活動", func(t *testing.T) {
		zeroTokenActivity := []types.Activity{
			{
				Type:   types.ActivityCoding,
				Tokens: types.TokenUsage{TotalTokens: 0},
			},
		}

		usage := stats.CalculateTokenUsageByActivity(zeroTokenActivity)
		if usage[types.ActivityCoding].TotalTokens != 0 {
			t.Errorf("零 Token 活動應正確處理")
		}
	})

	t.Run("無效時間範圍", func(t *testing.T) {
		invalidTimeActivity := []types.Activity{
			{
				Type:      types.ActivityCoding,
				Tokens:    types.TokenUsage{TotalTokens: 100},
				StartTime: time.Now(),
				EndTime:   time.Now().Add(-5 * time.Minute), // 結束時間早於開始時間
			},
		}

		totals := stats.CalculateActivityTotals(invalidTimeActivity)
		// 應該能處理，但時間統計可能為負數或零
		if totals.TotalActivities != 1 {
			t.Errorf("應該能處理無效時間範圍的活動")
		}
	})
}

// TestStatisticsAccuracy 測試統計計算的精確性
func TestStatisticsAccuracy(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	// 建立精確的測試數據
	activities := []types.Activity{
		{
			Type:   types.ActivityCoding,
			Tokens: types.TokenUsage{InputTokens: 123, OutputTokens: 456, TotalTokens: 579},
		},
		{
			Type:   types.ActivityCoding,
			Tokens: types.TokenUsage{InputTokens: 234, OutputTokens: 567, TotalTokens: 801},
		},
		{
			Type:   types.ActivityDebugging,
			Tokens: types.TokenUsage{InputTokens: 89, OutputTokens: 123, TotalTokens: 212},
		},
	}

	usage := stats.CalculateTokenUsageByActivity(activities)

	// 精確驗證編程活動統計
	codingUsage := usage[types.ActivityCoding]
	if codingUsage.InputTokens != 357 { // 123 + 234
		t.Errorf("編程活動輸入 Token 計算錯誤: 期望 357, 得到 %d", codingUsage.InputTokens)
	}
	if codingUsage.OutputTokens != 1023 { // 456 + 567
		t.Errorf("編程活動輸出 Token 計算錯誤: 期望 1023, 得到 %d", codingUsage.OutputTokens)
	}
	if codingUsage.TotalTokens != 1380 { // 579 + 801
		t.Errorf("編程活動總 Token 計算錯誤: 期望 1380, 得到 %d", codingUsage.TotalTokens)
	}

	// 精確驗證除錯活動統計
	debuggingUsage := usage[types.ActivityDebugging]
	if debuggingUsage.InputTokens != 89 {
		t.Errorf("除錯活動輸入 Token 計算錯誤: 期望 89, 得到 %d", debuggingUsage.InputTokens)
	}
	if debuggingUsage.OutputTokens != 123 {
		t.Errorf("除錯活動輸出 Token 計算錯誤: 期望 123, 得到 %d", debuggingUsage.OutputTokens)
	}
	if debuggingUsage.TotalTokens != 212 {
		t.Errorf("除錯活動總 Token 計算錯誤: 期望 212, 得到 %d", debuggingUsage.TotalTokens)
	}
}

// TestConcurrentStatisticsCalculation 測試並發統計計算的安全性
func TestConcurrentStatisticsCalculation(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	// 建立測試數據
	activityTypes := []types.ActivityType{
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
		types.ActivitySpecDev,
		types.ActivityTypeChat,
	}

	activities := make([]types.Activity, 1000)
	for i := 0; i < 1000; i++ {
		activities[i] = types.Activity{
			Type:   activityTypes[i%5],
			Tokens: types.TokenUsage{TotalTokens: i + 1},
		}
	}

	// 並發執行統計計算
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// 執行各種統計計算
			stats.CalculateTokenUsageByActivity(activities)
			stats.CalculateActivityTotals(activities)
			stats.AnalyzeUsagePatterns(activities)
			stats.GetTopActivitiesByTokens(activities, 10)
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 如果沒有 panic，測試通過
	t.Log("並發統計計算測試通過")
}
