package analyzer

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

func TestNewActivityAnalyzer(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created, got nil")
	}

	if len(analyzer.patterns) == 0 {
		t.Error("Expected patterns to be initialized")
	}

	if len(analyzer.keywords) == 0 {
		t.Error("Expected keywords to be initialized")
	}
}

func TestClassifyActivity(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	testCases := []struct {
		name     string
		content  string
		expected types.ActivityType
	}{
		{
			name:     "簡單編程活動",
			content:  "寫一個 function 來處理資料",
			expected: types.ActivityCoding,
		},
		{
			name:     "英文編程活動",
			content:  "implement a new function to handle data processing",
			expected: types.ActivityCoding,
		},
		{
			name:     "除錯活動",
			content:  "修復這個錯誤，程式一直當機",
			expected: types.ActivityDebugging,
		},
		{
			name:     "英文除錯活動",
			content:  "fix this bug, the application keeps crashing",
			expected: types.ActivityDebugging,
		},
		{
			name:     "文件撰寫",
			content:  "更新 README 文件，添加使用說明",
			expected: types.ActivityDocumentation,
		},
		{
			name:     "英文文件撰寫",
			content:  "write documentation for the new API",
			expected: types.ActivityDocumentation,
		},
		{
			name:     "規格開發",
			content:  "建立系統架構設計和需求分析",
			expected: types.ActivitySpecDev,
		},
		{
			name:     "英文規格開發",
			content:  "create system architecture design and requirement analysis",
			expected: types.ActivitySpecDev,
		},
		{
			name:     "一般聊天",
			content:  "你好，我需要一些協助",
			expected: types.ActivityChat,
		},
		{
			name:     "英文聊天",
			content:  "hello, I need some help with this",
			expected: types.ActivityChat,
		},
		{
			name:     "空內容",
			content:  "",
			expected: types.ActivityChat,
		},
		{
			name:     "混合活動 - 編程為主",
			content:  "我需要實作一個新的 class 來處理資料",
			expected: types.ActivityCoding,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzer.ClassifyActivity(tc.content)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v for content: %s", tc.expected, result, tc.content)
			}
		})
	}
}

func TestAnalyzeActivityBatch(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	contents := []string{
		"寫一個 function",
		"修復這個 bug",
		"更新文件",
		"建立規格",
		"需要協助",
	}

	expected := []types.ActivityType{
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
		types.ActivitySpecDev,
		types.ActivityChat,
	}

	results := analyzer.AnalyzeActivityBatch(contents)

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("Expected %v at index %d, got %v", expected[i], i, result)
		}
	}
}

func TestGenerateActivitySummary(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	now := time.Now()
	activities := []types.Activity{
		{
			ID:        "1",
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
			ID:        "2",
			Type:      types.ActivityCoding,
			StartTime: now,
			EndTime:   now.Add(15 * time.Minute),
			Tokens: types.TokenUsage{
				InputTokens:  150,
				OutputTokens: 250,
				TotalTokens:  400,
			},
		},
		{
			ID:        "3",
			Type:      types.ActivityDebugging,
			StartTime: now,
			EndTime:   now.Add(5 * time.Minute),
			Tokens: types.TokenUsage{
				InputTokens:  50,
				OutputTokens: 100,
				TotalTokens:  150,
			},
		},
	}

	summary := analyzer.GenerateActivitySummary(activities)

	// 檢查總活動數
	if summary.TotalActivities != 3 {
		t.Errorf("Expected 3 total activities, got %d", summary.TotalActivities)
	}

	// 檢查活動計數
	if summary.ActivityCounts[types.ActivityCoding] != 2 {
		t.Errorf("Expected 2 coding activities, got %d", summary.ActivityCounts[types.ActivityCoding])
	}

	if summary.ActivityCounts[types.ActivityDebugging] != 1 {
		t.Errorf("Expected 1 debugging activity, got %d", summary.ActivityCounts[types.ActivityDebugging])
	}

	// 檢查 Token 使用量
	codingTokens := summary.TokenUsage[types.ActivityCoding]
	if codingTokens.TotalTokens != 700 {
		t.Errorf("Expected 700 total tokens for coding, got %d", codingTokens.TotalTokens)
	}

	debuggingTokens := summary.TokenUsage[types.ActivityDebugging]
	if debuggingTokens.TotalTokens != 150 {
		t.Errorf("Expected 150 total tokens for debugging, got %d", debuggingTokens.TotalTokens)
	}

	// 檢查總 Token 數
	if summary.TotalTokens.TotalTokens != 850 {
		t.Errorf("Expected 850 total tokens, got %d", summary.TotalTokens.TotalTokens)
	}
}

func TestCalculateEfficiencyMetrics(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	now := time.Now()
	activities := []types.Activity{
		{
			ID:        "1",
			Type:      types.ActivityCoding,
			StartTime: now,
			EndTime:   now.Add(10 * time.Minute),
			Tokens: types.TokenUsage{
				TotalTokens: 300,
			},
		},
		{
			ID:        "2",
			Type:      types.ActivityCoding,
			StartTime: now,
			EndTime:   now.Add(20 * time.Minute),
			Tokens: types.TokenUsage{
				TotalTokens: 600,
			},
		},
	}

	data := types.ActivityData{
		Activities: activities,
		ActivitiesByType: map[types.ActivityType][]types.Activity{
			types.ActivityCoding: activities,
		},
	}

	metrics := analyzer.CalculateEfficiencyMetrics(data)

	// 檢查平均 Token 數
	avgTokens := metrics.AverageTokensPerActivity[types.ActivityCoding]
	if avgTokens != 450.0 {
		t.Errorf("Expected 450.0 average tokens, got %f", avgTokens)
	}

	// 檢查平均時間
	avgTime := metrics.AverageTimePerActivity[types.ActivityCoding]
	expectedTime := 15 * time.Minute
	if avgTime != expectedTime {
		t.Errorf("Expected %v average time, got %v", expectedTime, avgTime)
	}

	// 檢查每分鐘 Token 數
	tokensPerMin := metrics.TokensPerMinute[types.ActivityCoding]
	expectedRate := 900.0 / 30.0 // 總 Token 數 / 總分鐘數
	if tokensPerMin != expectedRate {
		t.Errorf("Expected %f tokens per minute, got %f", expectedRate, tokensPerMin)
	}
}

func TestGetActivityTypeDistribution(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	activities := []types.Activity{
		{Type: types.ActivityCoding},
		{Type: types.ActivityCoding},
		{Type: types.ActivityDebugging},
		{Type: types.ActivityChat},
	}

	distribution := analyzer.GetActivityTypeDistribution(activities)

	// 檢查分佈百分比
	if distribution[types.ActivityCoding] != 50.0 {
		t.Errorf("Expected 50%% for coding, got %f%%", distribution[types.ActivityCoding])
	}

	if distribution[types.ActivityDebugging] != 25.0 {
		t.Errorf("Expected 25%% for debugging, got %f%%", distribution[types.ActivityDebugging])
	}

	if distribution[types.ActivityChat] != 25.0 {
		t.Errorf("Expected 25%% for chat, got %f%%", distribution[types.ActivityChat])
	}

	// 檢查總和是否為 100%
	total := 0.0
	for _, percentage := range distribution {
		total += percentage
	}

	if total != 100.0 {
		t.Errorf("Expected total distribution to be 100%%, got %f%%", total)
	}
}

func TestGetActivityTypeDistribution_EmptyInput(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	activities := []types.Activity{}
	distribution := analyzer.GetActivityTypeDistribution(activities)

	if len(distribution) != 0 {
		t.Errorf("Expected empty distribution for empty input, got %v", distribution)
	}
}
