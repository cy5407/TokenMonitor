package analyzer

import (
	"regexp"
	"strings"
	"time"

	"token-monitor/internal/types"
)

// ActivityAnalyzer 負責分析和分類不同類型的 IDE 活動
type ActivityAnalyzer struct {
	patterns map[string]*regexp.Regexp
	keywords map[string][]string
}

// NewActivityAnalyzer 建立新的活動分析器實例
func NewActivityAnalyzer() *ActivityAnalyzer {
	analyzer := &ActivityAnalyzer{
		patterns: make(map[string]*regexp.Regexp),
		keywords: make(map[string][]string),
	}

	analyzer.initializePatterns()
	analyzer.initializeKeywords()

	return analyzer
}

// initializePatterns 初始化活動識別的正規表達式模式
func (aa *ActivityAnalyzer) initializePatterns() {
	patterns := map[string]string{
		"coding":           `(寫.*程式|實作.*功能|建立.*函數|implement.*function|write.*code|create.*class)`,
		"debugging":        `(修復.*問題|解決.*錯誤|debug.*issue|fix.*bug|troubleshoot)`,
		"documentation":    `(更新.*文件|撰寫.*說明|write.*readme|update.*doc|create.*documentation)`,
		"spec-development": `(建立.*規格|設計.*架構|create.*spec|design.*architecture|requirement.*analysis)`,
		"chat":             `(詢問.*問題|尋求.*協助|ask.*question|need.*help|how.*to)`,
	}

	for activityType, pattern := range patterns {
		compiled, err := regexp.Compile(`(?i)` + pattern)
		if err == nil {
			aa.patterns[activityType] = compiled
		}
	}
}

// initializeKeywords 初始化活動識別的關鍵字
func (aa *ActivityAnalyzer) initializeKeywords() {
	aa.keywords = map[string][]string{
		"coding": {
			"function", "class", "implement", "程式", "函數", "方法", "類別",
			"variable", "array", "loop", "condition", "變數", "陣列", "迴圈", "條件",
		},
		"debugging": {
			"error", "bug", "fix", "錯誤", "修復", "除錯", "問題",
			"exception", "crash", "fail", "異常", "當機", "失敗",
		},
		"documentation": {
			"README", "document", "文件", "說明", "註解", "comment",
			"manual", "guide", "tutorial", "手冊", "指南", "教學",
		},
		"spec-development": {
			"spec", "requirement", "design", "需求", "設計", "規格",
			"architecture", "plan", "analysis", "架構", "計畫", "分析",
		},
		"chat": {
			"chat", "question", "help", "問題", "協助", "詢問",
			"how", "what", "why", "怎麼", "什麼", "為什麼",
		},
	}
}

// ClassifyActivity 分析內容並分類活動類型
func (aa *ActivityAnalyzer) ClassifyActivity(content string) types.ActivityType {
	if content == "" {
		return types.ActivityTypeChat // 預設為聊天類型
	}

	content = strings.ToLower(content)
	scores := make(map[string]int)

	// 使用正規表達式模式評分
	for activityType, pattern := range aa.patterns {
		if pattern.MatchString(content) {
			scores[activityType] += 3 // 模式匹配權重較高
		}
	}

	// 使用關鍵字評分
	for activityType, keywords := range aa.keywords {
		for _, keyword := range keywords {
			if strings.Contains(content, strings.ToLower(keyword)) {
				scores[activityType]++
			}
		}
	}

	// 找出得分最高的活動類型
	maxScore := 0
	bestActivity := "chat" // 預設活動類型

	for activityType, score := range scores {
		if score > maxScore {
			maxScore = score
			bestActivity = activityType
		}
	}

	return types.StringToActivityType(bestActivity)
}

// AnalyzeActivityBatch 批次分析多個活動
func (aa *ActivityAnalyzer) AnalyzeActivityBatch(contents []string) []types.ActivityType {
	results := make([]types.ActivityType, len(contents))

	for i, content := range contents {
		results[i] = aa.ClassifyActivity(content)
	}

	return results
}

// GenerateActivitySummary 生成活動摘要統計
func (aa *ActivityAnalyzer) GenerateActivitySummary(activities []types.Activity) types.ActivitySummary {
	summary := types.ActivitySummary{
		TotalActivities: len(activities),
		ActivityCounts:  make(map[types.ActivityType]int),
		TokenUsage:      make(map[types.ActivityType]types.TokenUsage),
		TimeSpent:       make(map[types.ActivityType]time.Duration),
		GeneratedAt:     time.Now(),
	}

	// 統計各活動類型的數量和 Token 使用量
	for _, activity := range activities {
		summary.ActivityCounts[activity.Type]++

		// 累加 Token 使用量
		usage := summary.TokenUsage[activity.Type]
		usage.InputTokens += activity.Tokens.InputTokens
		usage.OutputTokens += activity.Tokens.OutputTokens
		usage.TotalTokens += activity.Tokens.TotalTokens
		summary.TokenUsage[activity.Type] = usage

		// 累加時間
		if !activity.StartTime.IsZero() && !activity.EndTime.IsZero() {
			duration := activity.EndTime.Sub(activity.StartTime)
			summary.TimeSpent[activity.Type] += duration
		}
	}

	// 計算總 Token 使用量
	for _, usage := range summary.TokenUsage {
		summary.TotalTokens.InputTokens += usage.InputTokens
		summary.TotalTokens.OutputTokens += usage.OutputTokens
		summary.TotalTokens.TotalTokens += usage.TotalTokens
	}

	return summary
}

// CalculateEfficiencyMetrics 計算效率指標
func (aa *ActivityAnalyzer) CalculateEfficiencyMetrics(data types.ActivityData) types.EfficiencyMetrics {
	metrics := types.EfficiencyMetrics{
		CalculatedAt:             time.Now(),
		AverageTokensPerActivity: make(map[types.ActivityType]float64),
		AverageTimePerActivity:   make(map[types.ActivityType]time.Duration),
		TokensPerMinute:          make(map[types.ActivityType]float64),
	}

	// 計算各活動類型的平均指標
	for activityType, activities := range data.ActivitiesByType {
		if len(activities) == 0 {
			continue
		}

		totalTokens := 0
		totalTime := time.Duration(0)

		for _, activity := range activities {
			totalTokens += activity.Tokens.TotalTokens
			if !activity.StartTime.IsZero() && !activity.EndTime.IsZero() {
				totalTime += activity.EndTime.Sub(activity.StartTime)
			}
		}

		// 平均 Token 數
		metrics.AverageTokensPerActivity[activityType] = float64(totalTokens) / float64(len(activities))

		// 平均時間
		if len(activities) > 0 {
			metrics.AverageTimePerActivity[activityType] = totalTime / time.Duration(len(activities))
		}

		// 每分鐘 Token 數
		if totalTime.Minutes() > 0 {
			metrics.TokensPerMinute[activityType] = float64(totalTokens) / totalTime.Minutes()
		}
	}

	return metrics
}

// GetActivityTypeDistribution 獲取活動類型分佈
func (aa *ActivityAnalyzer) GetActivityTypeDistribution(activities []types.Activity) map[types.ActivityType]float64 {
	if len(activities) == 0 {
		return make(map[types.ActivityType]float64)
	}

	counts := make(map[types.ActivityType]int)
	total := len(activities)

	for _, activity := range activities {
		counts[activity.Type]++
	}

	distribution := make(map[types.ActivityType]float64)
	for activityType, count := range counts {
		distribution[activityType] = float64(count) / float64(total) * 100
	}

	return distribution
}
