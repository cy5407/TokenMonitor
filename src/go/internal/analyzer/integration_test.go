package analyzer

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

func TestActivityAnalyzerIntegration(t *testing.T) {
	// 建立活動分析器
	analyzer := NewActivityAnalyzer()

	// 測試內容
	testContents := []string{
		"寫一個 function 來處理資料",
		"修復這個錯誤，程式一直當機",
		"更新 README 文件，添加使用說明",
		"建立系統架構設計和需求分析",
		"你好，我需要一些協助",
		"implement a new class for data processing",
		"fix this bug in the authentication module",
		"write documentation for the API endpoints",
	}

	expectedTypes := []types.ActivityType{
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
		types.ActivitySpecDev,
		types.ActivityChat,
		types.ActivityCoding,
		types.ActivityDebugging,
		types.ActivityDocumentation,
	}

	t.Run("單個活動分類", func(t *testing.T) {
		for i, content := range testContents {
			result := analyzer.ClassifyActivity(content)
			if result != expectedTypes[i] {
				t.Errorf("內容 '%s' 期望分類為 %s，實際為 %s", content, expectedTypes[i], result)
			}
		}
	})

	t.Run("批次活動分析", func(t *testing.T) {
		results := analyzer.AnalyzeActivityBatch(testContents)
		if len(results) != len(expectedTypes) {
			t.Fatalf("期望 %d 個結果，實際得到 %d 個", len(expectedTypes), len(results))
		}

		for i, result := range results {
			if result != expectedTypes[i] {
				t.Errorf("索引 %d: 期望 %s，實際 %s", i, expectedTypes[i], result)
			}
		}
	})

	t.Run("完整工作流程測試", func(t *testing.T) {
		// 建立活動記錄
		now := time.Now()
		activities := []types.Activity{
			{
				ID:        "1",
				Type:      types.ActivityCoding,
				Content:   "寫一個 function 來處理資料",
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
				Content:   "implement a new class for data processing",
				StartTime: now.Add(15 * time.Minute),
				EndTime:   now.Add(30 * time.Minute),
				Tokens: types.TokenUsage{
					InputTokens:  150,
					OutputTokens: 250,
					TotalTokens:  400,
				},
			},
			{
				ID:        "3",
				Type:      types.ActivityDebugging,
				Content:   "修復這個錯誤，程式一直當機",
				StartTime: now.Add(35 * time.Minute),
				EndTime:   now.Add(45 * time.Minute),
				Tokens: types.TokenUsage{
					InputTokens:  50,
					OutputTokens: 100,
					TotalTokens:  150,
				},
			},
			{
				ID:        "4",
				Type:      types.ActivityDocumentation,
				Content:   "更新 README 文件，添加使用說明",
				StartTime: now.Add(50 * time.Minute),
				EndTime:   now.Add(60 * time.Minute),
				Tokens: types.TokenUsage{
					InputTokens:  75,
					OutputTokens: 125,
					TotalTokens:  200,
				},
			},
		}

		// 測試活動摘要生成
		summary := analyzer.GenerateActivitySummary(activities)

		if summary.TotalActivities != 4 {
			t.Errorf("期望總活動數為 4，實際為 %d", summary.TotalActivities)
		}

		if summary.TotalTokens.TotalTokens != 1050 {
			t.Errorf("期望總 Token 數為 1050，實際為 %d", summary.TotalTokens.TotalTokens)
		}

		if summary.ActivityCounts[types.ActivityCoding] != 2 {
			t.Errorf("期望編程活動數為 2，實際為 %d", summary.ActivityCounts[types.ActivityCoding])
		}

		if summary.TokenUsage[types.ActivityCoding].TotalTokens != 700 {
			t.Errorf("期望編程活動 Token 數為 700，實際為 %d", summary.TokenUsage[types.ActivityCoding].TotalTokens)
		}

		// 測試效率指標計算
		data := types.ActivityData{
			Activities: activities,
			ActivitiesByType: map[types.ActivityType][]types.Activity{
				types.ActivityCoding:        {activities[0], activities[1]},
				types.ActivityDebugging:     {activities[2]},
				types.ActivityDocumentation: {activities[3]},
			},
		}

		metrics := analyzer.CalculateEfficiencyMetrics(data)

		// 檢查編程活動的平均 Token 數 (300 + 400) / 2 = 350
		if metrics.AverageTokensPerActivity[types.ActivityCoding] != 350.0 {
			t.Errorf("期望編程活動平均 Token 數為 350.0，實際為 %f",
				metrics.AverageTokensPerActivity[types.ActivityCoding])
		}

		// 檢查除錯活動的每分鐘 Token 數 150 / 10 = 15
		if metrics.TokensPerMinute[types.ActivityDebugging] != 15.0 {
			t.Errorf("期望除錯活動每分鐘 Token 數為 15.0，實際為 %f",
				metrics.TokensPerMinute[types.ActivityDebugging])
		}

		// 測試活動類型分佈
		distribution := analyzer.GetActivityTypeDistribution(activities)

		if distribution[types.ActivityCoding] != 50.0 {
			t.Errorf("期望編程活動佔比為 50.0%%，實際為 %f%%", distribution[types.ActivityCoding])
		}

		if distribution[types.ActivityDebugging] != 25.0 {
			t.Errorf("期望除錯活動佔比為 25.0%%，實際為 %f%%", distribution[types.ActivityDebugging])
		}

		if distribution[types.ActivityDocumentation] != 25.0 {
			t.Errorf("期望文件活動佔比為 25.0%%，實際為 %f%%", distribution[types.ActivityDocumentation])
		}
	})
}

func TestActivityAnalyzerEdgeCases(t *testing.T) {
	analyzer := NewActivityAnalyzer()

	t.Run("空輸入處理", func(t *testing.T) {
		result := analyzer.ClassifyActivity("")
		if result != types.ActivityChat {
			t.Errorf("空輸入期望分類為 chat，實際為 %s", result)
		}

		results := analyzer.AnalyzeActivityBatch([]string{})
		if len(results) != 0 {
			t.Errorf("空批次期望返回空結果，實際長度為 %d", len(results))
		}

		summary := analyzer.GenerateActivitySummary([]types.Activity{})
		if summary.TotalActivities != 0 {
			t.Errorf("空活動列表期望總數為 0，實際為 %d", summary.TotalActivities)
		}

		distribution := analyzer.GetActivityTypeDistribution([]types.Activity{})
		if len(distribution) != 0 {
			t.Errorf("空活動列表期望空分佈，實際長度為 %d", len(distribution))
		}
	})

	t.Run("混合語言內容", func(t *testing.T) {
		mixedContents := []string{
			"我需要 implement 一個新的 function",
			"這個 bug 需要修復，程式有問題",
			"Update the README file with 中文說明",
		}

		results := analyzer.AnalyzeActivityBatch(mixedContents)

		if results[0] != types.ActivityCoding {
			t.Errorf("混合編程內容期望分類為 coding，實際為 %s", results[0])
		}

		if results[1] != types.ActivityDebugging {
			t.Errorf("混合除錯內容期望分類為 debugging，實際為 %s", results[1])
		}

		if results[2] != types.ActivityDocumentation {
			t.Errorf("混合文件內容期望分類為 documentation，實際為 %s", results[2])
		}
	})
}
