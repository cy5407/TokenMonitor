package analyzer

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

func TestNewActivityStatistics(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	if stats == nil {
		t.Fatal("Expected statistics analyzer to be created, got nil")
	}

	if stats.analyzer != analyzer {
		t.Error("Expected analyzer to be set correctly")
	}
}

func TestCalculateTokenUsageByActivity(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	activities := []types.Activity{
		{
			Type: types.ActivityCoding,
			Tokens: types.TokenUsage{
				InputTokens:  100,
				OutputTokens: 200,
				TotalTokens:  300,
			},
		},
		{
			Type: types.ActivityCoding,
			Tokens: types.TokenUsage{
				InputTokens:  150,
				OutputTokens: 250,
				TotalTokens:  400,
			},
		},
		{
			Type: types.ActivityDebugging,
			Tokens: types.TokenUsage{
				InputTokens:  50,
				OutputTokens: 100,
				TotalTokens:  150,
			},
		},
	}

	usage := stats.CalculateTokenUsageByActivity(activities)

	// 檢查編程活動的 Token 使用量
	codingUsage := usage[types.ActivityCoding]
	if codingUsage.InputTokens != 250 {
		t.Errorf("Expected coding input tokens 250, got %d", codingUsage.InputTokens)
	}
	if codingUsage.OutputTokens != 450 {
		t.Errorf("Expected coding output tokens 450, got %d", codingUsage.OutputTokens)
	}
	if codingUsage.TotalTokens != 700 {
		t.Errorf("Expected coding total tokens 700, got %d", codingUsage.TotalTokens)
	}

	// 檢查除錯活動的 Token 使用量
	debuggingUsage := usage[types.ActivityDebugging]
	if debuggingUsage.TotalTokens != 150 {
		t.Errorf("Expected debugging total tokens 150, got %d", debuggingUsage.TotalTokens)
	}
}

func TestCalculateActivityTotals(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	now := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			StartTime: now,
			EndTime:   now.Add(10 * time.Minute),
			Tokens: types.TokenUsage{
				InputTokens:  100,
				OutputTokens: 200,
				TotalTokens:  300,
			},
		},
		{
			Type:      types.ActivityCoding,
			StartTime: now.Add(15 * time.Minute),
			EndTime:   now.Add(30 * time.Minute),
			Tokens: types.TokenUsage{
				InputTokens:  150,
				OutputTokens: 250,
				TotalTokens:  400,
			},
		},
		{
			Type:      types.ActivityDebugging,
			StartTime: now.Add(35 * time.Minute),
			EndTime:   now.Add(45 * time.Minute),
			Tokens: types.TokenUsage{
				InputTokens:  50,
				OutputTokens: 100,
				TotalTokens:  150,
			},
		},
	}

	totals := stats.CalculateActivityTotals(activities)

	// 檢查總計
	if totals.TotalActivities != 3 {
		t.Errorf("Expected 3 total activities, got %d", totals.TotalActivities)
	}

	if totals.TotalTokens.TotalTokens != 850 {
		t.Errorf("Expected 850 total tokens, got %d", totals.TotalTokens.TotalTokens)
	}

	expectedTotalTime := 35 * time.Minute // 10 + 15 + 10 分鐘
	if totals.TotalTime != expectedTotalTime {
		t.Errorf("Expected total time %v, got %v", expectedTotalTime, totals.TotalTime)
	}

	// 檢查各類型統計
	codingTotal := totals.ByType[types.ActivityCoding]
	if codingTotal.Count != 2 {
		t.Errorf("Expected 2 coding activities, got %d", codingTotal.Count)
	}

	if codingTotal.Tokens.TotalTokens != 700 {
		t.Errorf("Expected 700 coding tokens, got %d", codingTotal.Tokens.TotalTokens)
	}

	debuggingTotal := totals.ByType[types.ActivityDebugging]
	if debuggingTotal.Count != 1 {
		t.Errorf("Expected 1 debugging activity, got %d", debuggingTotal.Count)
	}
}

func TestAnalyzeUsagePatterns(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	now := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			StartTime: now,
			EndTime:   now.Add(10 * time.Minute),
			Tokens: types.TokenUsage{
				TotalTokens: 300,
			},
		},
		{
			Type:      types.ActivityCoding,
			StartTime: now.Add(15 * time.Minute),
			EndTime:   now.Add(30 * time.Minute),
			Tokens: types.TokenUsage{
				TotalTokens: 500,
			},
		},
		{
			Type:      types.ActivityDebugging,
			StartTime: now.Add(35 * time.Minute),
			EndTime:   now.Add(45 * time.Minute),
			Tokens: types.TokenUsage{
				TotalTokens: 150,
			},
		},
	}

	analysis := stats.AnalyzeUsagePatterns(activities)

	// 檢查基本統計
	if analysis.TotalActivities != 3 {
		t.Errorf("Expected 3 total activities, got %d", analysis.TotalActivities)
	}

	// 檢查模式分析
	codingPattern := analysis.Patterns["coding"]
	if codingPattern.Count != 2 {
		t.Errorf("Expected 2 coding activities, got %d", codingPattern.Count)
	}

	if codingPattern.TotalTokens != 800 {
		t.Errorf("Expected 800 total tokens for coding, got %d", codingPattern.TotalTokens)
	}

	if codingPattern.AverageTokens != 400.0 {
		t.Errorf("Expected 400.0 average tokens for coding, got %f", codingPattern.AverageTokens)
	}

	// 檢查洞察是否生成
	if len(analysis.Insights) == 0 {
		t.Error("Expected insights to be generated")
	}
}

func TestGetTopActivitiesByTokens(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	activities := []types.Activity{
		{
			ID:   "1",
			Type: types.ActivityCoding,
			Tokens: types.TokenUsage{
				TotalTokens: 100,
			},
		},
		{
			ID:   "2",
			Type: types.ActivityDebugging,
			Tokens: types.TokenUsage{
				TotalTokens: 300,
			},
		},
		{
			ID:   "3",
			Type: types.ActivityDocumentation,
			Tokens: types.TokenUsage{
				TotalTokens: 200,
			},
		},
	}

	top2 := stats.GetTopActivitiesByTokens(activities, 2)

	if len(top2) != 2 {
		t.Fatalf("Expected 2 activities, got %d", len(top2))
	}

	// 檢查排序是否正確（按 Token 數量降序）
	if top2[0].ID != "2" {
		t.Errorf("Expected first activity ID '2', got '%s'", top2[0].ID)
	}

	if top2[1].ID != "3" {
		t.Errorf("Expected second activity ID '3', got '%s'", top2[1].ID)
	}
}

func TestCalculateActivityFrequency(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	now := time.Now()
	activities := []types.Activity{
		{
			Type:      types.ActivityCoding,
			Timestamp: now,
		},
		{
			Type:      types.ActivityCoding,
			Timestamp: now.Add(1 * time.Hour),
		},
		{
			Type:      types.ActivityDebugging,
			Timestamp: now.Add(2 * time.Hour),
		},
	}

	frequency := stats.CalculateActivityFrequency(activities, 1*time.Hour)

	// 檢查基本統計
	if frequency.TimeWindow != 1*time.Hour {
		t.Errorf("Expected time window 1 hour, got %v", frequency.TimeWindow)
	}

	// 檢查按小時分佈
	expectedHours := []int{now.Hour(), now.Add(1 * time.Hour).Hour(), now.Add(2 * time.Hour).Hour()}
	for _, hour := range expectedHours {
		if frequency.ByHour[hour] == 0 {
			t.Errorf("Expected activity at hour %d", hour)
		}
	}
}

func TestEmptyInputHandling(t *testing.T) {
	analyzer := NewActivityAnalyzer()
	stats := NewActivityStatistics(analyzer)

	emptyActivities := []types.Activity{}

	// 測試空輸入處理
	usage := stats.CalculateTokenUsageByActivity(emptyActivities)
	if len(usage) != 0 {
		t.Error("Expected empty usage map for empty input")
	}

	totals := stats.CalculateActivityTotals(emptyActivities)
	if totals.TotalActivities != 0 {
		t.Error("Expected 0 total activities for empty input")
	}

	analysis := stats.AnalyzeUsagePatterns(emptyActivities)
	if analysis.TotalActivities != 0 {
		t.Error("Expected 0 total activities in pattern analysis for empty input")
	}

	top := stats.GetTopActivitiesByTokens(emptyActivities, 5)
	if len(top) != 0 {
		t.Error("Expected empty result for empty input")
	}

	frequency := stats.CalculateActivityFrequency(emptyActivities, 1*time.Hour)
	if len(frequency.ByType) != 0 {
		t.Error("Expected empty frequency map for empty input")
	}
}
