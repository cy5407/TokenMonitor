package analyzer

import (
	"sort"
	"time"

	"token-monitor/internal/types"
)

// ActivityStatistics 活動統計分析器
type ActivityStatistics struct {
	analyzer *ActivityAnalyzer
}

// NewActivityStatistics 建立新的活動統計分析器
func NewActivityStatistics(analyzer *ActivityAnalyzer) *ActivityStatistics {
	return &ActivityStatistics{
		analyzer: analyzer,
	}
}

// CalculateTokenUsageByActivity 計算各活動類型的 Token 使用統計
func (as *ActivityStatistics) CalculateTokenUsageByActivity(activities []types.Activity) map[types.ActivityType]types.TokenUsage {
	usage := make(map[types.ActivityType]types.TokenUsage)

	for _, activity := range activities {
		current := usage[activity.Type]
		current.InputTokens += activity.Tokens.InputTokens
		current.OutputTokens += activity.Tokens.OutputTokens
		current.TotalTokens += activity.Tokens.TotalTokens
		usage[activity.Type] = current
	}

	return usage
}

// CalculateActivityTotals 實作活動總和統計功能
func (as *ActivityStatistics) CalculateActivityTotals(activities []types.Activity) types.ActivityTotals {
	totals := types.ActivityTotals{
		TotalActivities: len(activities),
		TotalTokens:     types.TokenUsage{},
		TotalTime:       0,
		ByType:          make(map[types.ActivityType]types.ActivityTypeTotal),
		CalculatedAt:    time.Now(),
	}

	// 計算各活動類型的統計
	for _, activity := range activities {
		// 更新總計
		totals.TotalTokens.InputTokens += activity.Tokens.InputTokens
		totals.TotalTokens.OutputTokens += activity.Tokens.OutputTokens
		totals.TotalTokens.TotalTokens += activity.Tokens.TotalTokens

		// 計算活動時間
		var duration time.Duration
		if !activity.StartTime.IsZero() && !activity.EndTime.IsZero() {
			duration = activity.EndTime.Sub(activity.StartTime)
			totals.TotalTime += duration
		}

		// 更新各類型統計
		typeTotal := totals.ByType[activity.Type]
		typeTotal.Count++
		typeTotal.Tokens.InputTokens += activity.Tokens.InputTokens
		typeTotal.Tokens.OutputTokens += activity.Tokens.OutputTokens
		typeTotal.Tokens.TotalTokens += activity.Tokens.TotalTokens
		typeTotal.TotalTime += duration
		totals.ByType[activity.Type] = typeTotal
	}

	return totals
}

// AnalyzeUsagePatterns 建立效率分析和使用模式識別
func (as *ActivityStatistics) AnalyzeUsagePatterns(activities []types.Activity) types.UsagePatternAnalysis {
	analysis := types.UsagePatternAnalysis{
		AnalyzedAt:      time.Now(),
		TotalActivities: len(activities),
		Patterns:        make(map[string]types.UsagePattern),
		Insights:        []string{},
	}

	if len(activities) == 0 {
		return analysis
	}

	// 按活動類型分組
	byType := make(map[types.ActivityType][]types.Activity)
	for _, activity := range activities {
		byType[activity.Type] = append(byType[activity.Type], activity)
	}

	// 分析各活動類型的模式
	for activityType, typeActivities := range byType {
		pattern := as.analyzeActivityTypePattern(activityType, typeActivities)
		analysis.Patterns[string(activityType)] = pattern
	}

	// 生成洞察
	analysis.Insights = as.generateInsights(analysis.Patterns)

	return analysis
}

// analyzeActivityTypePattern 分析特定活動類型的模式
func (as *ActivityStatistics) analyzeActivityTypePattern(activityType types.ActivityType, activities []types.Activity) types.UsagePattern {
	if len(activities) == 0 {
		return types.UsagePattern{}
	}

	pattern := types.UsagePattern{
		ActivityType: activityType,
		Count:        len(activities),
		TotalTokens:  0,
		TotalTime:    0,
	}

	// 收集統計數據
	tokenCounts := make([]int, len(activities))
	durations := make([]time.Duration, 0)

	for i, activity := range activities {
		pattern.TotalTokens += activity.Tokens.TotalTokens
		tokenCounts[i] = activity.Tokens.TotalTokens

		if !activity.StartTime.IsZero() && !activity.EndTime.IsZero() {
			duration := activity.EndTime.Sub(activity.StartTime)
			pattern.TotalTime += duration
			durations = append(durations, duration)
		}
	}

	// 計算平均值
	pattern.AverageTokens = float64(pattern.TotalTokens) / float64(len(activities))
	if len(durations) > 0 {
		pattern.AverageTime = pattern.TotalTime / time.Duration(len(durations))
	}

	// 計算效率指標
	if pattern.TotalTime.Minutes() > 0 {
		pattern.TokensPerMinute = float64(pattern.TotalTokens) / pattern.TotalTime.Minutes()
	}

	// 計算變異性（標準差）
	pattern.TokenVariability = as.calculateStandardDeviation(tokenCounts, pattern.AverageTokens)

	return pattern
}

// calculateStandardDeviation 計算標準差
func (as *ActivityStatistics) calculateStandardDeviation(values []int, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	var sum float64
	for _, value := range values {
		diff := float64(value) - mean
		sum += diff * diff
	}

	variance := sum / float64(len(values)-1)
	return variance // 簡化版本，實際應該開平方根
}

// generateInsights 生成使用模式洞察
func (as *ActivityStatistics) generateInsights(patterns map[string]types.UsagePattern) []string {
	insights := []string{}

	// 找出最活躍的活動類型
	var mostActive types.UsagePattern
	var mostActiveType string
	for typeName, pattern := range patterns {
		if pattern.Count > mostActive.Count {
			mostActive = pattern
			mostActiveType = typeName
		}
	}

	if mostActiveType != "" {
		insights = append(insights, "最活躍的活動類型是 "+mostActiveType+"，共 "+
			string(rune(mostActive.Count))+" 次活動")
	}

	// 找出最耗費 Token 的活動類型
	var mostTokens types.UsagePattern
	var mostTokensType string
	for typeName, pattern := range patterns {
		if pattern.TotalTokens > mostTokens.TotalTokens {
			mostTokens = pattern
			mostTokensType = typeName
		}
	}

	if mostTokensType != "" {
		insights = append(insights, "最耗費 Token 的活動類型是 "+mostTokensType+"，共使用 "+
			string(rune(mostTokens.TotalTokens))+" 個 Token")
	}

	// 找出效率最高的活動類型
	var mostEfficient types.UsagePattern
	var mostEfficientType string
	for typeName, pattern := range patterns {
		if pattern.TokensPerMinute > mostEfficient.TokensPerMinute {
			mostEfficient = pattern
			mostEfficientType = typeName
		}
	}

	if mostEfficientType != "" {
		insights = append(insights, "效率最高的活動類型是 "+mostEfficientType+"，每分鐘處理 "+
			string(rune(int(mostEfficient.TokensPerMinute)))+" 個 Token")
	}

	return insights
}

// GetTopActivitiesByTokens 取得 Token 使用量最高的活動
func (as *ActivityStatistics) GetTopActivitiesByTokens(activities []types.Activity, limit int) []types.Activity {
	if len(activities) == 0 {
		return []types.Activity{}
	}

	// 複製切片以避免修改原始數據
	sorted := make([]types.Activity, len(activities))
	copy(sorted, activities)

	// 按 Token 數量排序
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Tokens.TotalTokens > sorted[j].Tokens.TotalTokens
	})

	// 返回前 N 個
	if limit > len(sorted) {
		limit = len(sorted)
	}

	return sorted[:limit]
}

// CalculateActivityFrequency 計算活動頻率分析
func (as *ActivityStatistics) CalculateActivityFrequency(activities []types.Activity, timeWindow time.Duration) types.ActivityFrequency {
	frequency := types.ActivityFrequency{
		TimeWindow:   timeWindow,
		AnalyzedAt:   time.Now(),
		ByType:       make(map[types.ActivityType]float64),
		ByHour:       make(map[int]int),
		TotalPeriods: 0,
	}

	if len(activities) == 0 {
		return frequency
	}

	// 找出時間範圍
	var earliest, latest time.Time
	for i, activity := range activities {
		if i == 0 || activity.Timestamp.Before(earliest) {
			earliest = activity.Timestamp
		}
		if i == 0 || activity.Timestamp.After(latest) {
			latest = activity.Timestamp
		}
	}

	// 計算總時間段數
	totalDuration := latest.Sub(earliest)
	frequency.TotalPeriods = int(totalDuration / timeWindow)
	if frequency.TotalPeriods == 0 {
		frequency.TotalPeriods = 1
	}

	// 統計各類型活動數量
	typeCounts := make(map[types.ActivityType]int)
	for _, activity := range activities {
		typeCounts[activity.Type]++

		// 統計按小時分佈
		hour := activity.Timestamp.Hour()
		frequency.ByHour[hour]++
	}

	// 計算頻率（每個時間窗口的平均活動數）
	for activityType, count := range typeCounts {
		frequency.ByType[activityType] = float64(count) / float64(frequency.TotalPeriods)
	}

	return frequency
}
