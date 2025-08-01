package types

import "time"

// ActivityType 定義活動類型
type ActivityType string

const (
	ActivityCoding        ActivityType = "coding"
	ActivityDebugging     ActivityType = "debugging"
	ActivityDocumentation ActivityType = "documentation"
	ActivitySpecDev       ActivityType = "spec-development"
	ActivityChat          ActivityType = "chat"
	ActivityTypeChat      ActivityType = "chat" // Alias for consistency
)

// TokenDistribution Token 分佈資訊
type TokenDistribution struct {
	EnglishTokens int    `json:"english_tokens"`
	ChineseTokens int    `json:"chinese_tokens"`
	TotalTokens   int    `json:"total_tokens"`
	Method        string `json:"method"`
}

// Activity 活動記錄
type Activity struct {
	ID          string       `json:"id"`
	Type        ActivityType `json:"type"`
	Description string       `json:"description"`
	Content     string       `json:"content"`
	Rounds      int          `json:"rounds"`
	Timestamp   time.Time    `json:"timestamp"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
	Tokens      TokenUsage   `json:"tokens"`
}

// TokenUsage Token 使用量
type TokenUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ActivitySummary 活動摘要
type ActivitySummary struct {
	TotalActivities int                            `json:"total_activities"`
	ActivityCounts  map[ActivityType]int           `json:"activity_counts"`
	TokenUsage      map[ActivityType]TokenUsage    `json:"token_usage"`
	TimeSpent       map[ActivityType]time.Duration `json:"time_spent"`
	TotalTokens     TokenUsage                     `json:"total_tokens"`
	GeneratedAt     time.Time                      `json:"generated_at"`
	ByType          map[ActivityType]int           `json:"by_type"`        // Keep for backward compatibility
	TokensByType    map[ActivityType]int           `json:"tokens_by_type"` // Keep for backward compatibility
	CostByType      map[ActivityType]float64       `json:"cost_by_type"`   // Keep for backward compatibility
}

// ActivityData 活動資料集合
type ActivityData struct {
	Activities       []Activity                  `json:"activities"`
	ActivitiesByType map[ActivityType][]Activity `json:"activities_by_type"`
	TimeRange        TimeRange                   `json:"time_range"`
}

// TimeRange 時間範圍
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UsageRecord 使用記錄
type UsageRecord struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"session_id"`
	Activity  Activity  `json:"activity"`
	Tokens    struct {
		Input             int    `json:"input"`
		Output            int    `json:"output"`
		Total             int    `json:"total"`
		CalculationMethod string `json:"calculation_method"`
	} `json:"tokens"`
	Cost struct {
		Input        float64 `json:"input"`
		Output       float64 `json:"output"`
		Total        float64 `json:"total"`
		Currency     string  `json:"currency"`
		PricingModel string  `json:"pricing_model"`
	} `json:"cost"`
}

// ActivityTotals 活動總和統計
type ActivityTotals struct {
	TotalActivities int                                `json:"total_activities"`
	TotalTokens     TokenUsage                         `json:"total_tokens"`
	TotalTime       time.Duration                      `json:"total_time"`
	ByType          map[ActivityType]ActivityTypeTotal `json:"by_type"`
	CalculatedAt    time.Time                          `json:"calculated_at"`
}

// ActivityTypeTotal 活動類型總計
type ActivityTypeTotal struct {
	Count     int           `json:"count"`
	Tokens    TokenUsage    `json:"tokens"`
	TotalTime time.Duration `json:"total_time"`
}

// UsagePatternAnalysis 使用模式分析
type UsagePatternAnalysis struct {
	AnalyzedAt      time.Time               `json:"analyzed_at"`
	TotalActivities int                     `json:"total_activities"`
	Patterns        map[string]UsagePattern `json:"patterns"`
	Insights        []string                `json:"insights"`
}

// UsagePattern 使用模式
type UsagePattern struct {
	ActivityType     ActivityType  `json:"activity_type"`
	Count            int           `json:"count"`
	TotalTokens      int           `json:"total_tokens"`
	TotalTime        time.Duration `json:"total_time"`
	AverageTokens    float64       `json:"average_tokens"`
	AverageTime      time.Duration `json:"average_time"`
	TokensPerMinute  float64       `json:"tokens_per_minute"`
	TokenVariability float64       `json:"token_variability"`
}

// EfficiencyMetrics 效率指標
type EfficiencyMetrics struct {
	CalculatedAt             time.Time                      `json:"calculated_at"`
	AverageTokensPerActivity map[ActivityType]float64       `json:"average_tokens_per_activity"`
	AverageTimePerActivity   map[ActivityType]time.Duration `json:"average_time_per_activity"`
	TokensPerMinute          map[ActivityType]float64       `json:"tokens_per_minute"`
}

// ActivityFrequency 活動頻率分析
type ActivityFrequency struct {
	TimeWindow   time.Duration            `json:"time_window"`
	AnalyzedAt   time.Time                `json:"analyzed_at"`
	ByType       map[ActivityType]float64 `json:"by_type"`
	ByHour       map[int]int              `json:"by_hour"`
	TotalPeriods int                      `json:"total_periods"`
}

// CostBreakdown 成本分解
type CostBreakdown struct {
	InputCost      float64      `json:"input_cost"`
	OutputCost     float64      `json:"output_cost"`
	CacheReadCost  float64      `json:"cache_read_cost,omitempty"`
	CacheWriteCost float64      `json:"cache_write_cost,omitempty"`
	BatchDiscount  float64      `json:"batch_discount,omitempty"`
	TotalCost      float64      `json:"total_cost"`
	Currency       string       `json:"currency"`
	PricingModel   string       `json:"pricing_model"`
	TokenCounts    TokenCounts  `json:"token_counts"`
	CostDetails    CostDetails  `json:"cost_details"`
	Timestamp      time.Time    `json:"timestamp"`
	SessionID      string       `json:"session_id,omitempty"`
	ActivityType   ActivityType `json:"activity_type,omitempty"`
}

// TokenCounts Token 數量詳細資訊
type TokenCounts struct {
	Input      int `json:"input"`
	Output     int `json:"output"`
	CacheRead  int `json:"cache_read,omitempty"`
	CacheWrite int `json:"cache_write,omitempty"`
	Total      int `json:"total"`
}

// CostDetails 成本詳細資訊
type CostDetails struct {
	InputRate      float64 `json:"input_rate"`                 // USD per 1M tokens
	OutputRate     float64 `json:"output_rate"`                // USD per 1M tokens
	CacheReadRate  float64 `json:"cache_read_rate,omitempty"`  // USD per 1M tokens
	CacheWriteRate float64 `json:"cache_write_rate,omitempty"` // USD per 1M tokens
	DiscountRate   float64 `json:"discount_rate,omitempty"`    // Discount percentage
	BillingMode    string  `json:"billing_mode"`               // standard, cache, batch
}

// PricingModel 定價模型
type PricingModel struct {
	Name          string  `json:"name"`
	InputPrice    float64 `json:"input_price"`    // USD per 1M tokens
	OutputPrice   float64 `json:"output_price"`   // USD per 1M tokens
	CacheRead     float64 `json:"cache_read"`     // USD per 1M tokens
	CacheWrite    float64 `json:"cache_write"`    // USD per 1M tokens
	BatchDiscount float64 `json:"batch_discount"` // Discount percentage
}

// OptimizationSuggestion 優化建議
type OptimizationSuggestion struct {
	Type            string  `json:"type"`
	Description     string  `json:"description"`
	PotentialSaving float64 `json:"potential_saving"`
	Confidence      float64 `json:"confidence"`
}

// OptimizationSuggestions 優化建議集合
type OptimizationSuggestions struct {
	Suggestions   []OptimizationSuggestion `json:"suggestions"`
	TotalSavings  float64                  `json:"total_savings"`
	CurrentCost   float64                  `json:"current_cost"`
	OptimizedCost float64                  `json:"optimized_cost"`
}

// TrendAnalysis 趨勢分析
type TrendAnalysis struct {
	DailyAverage float64 `json:"daily_average"`
	WeeklyTrend  float64 `json:"weekly_trend"`
	MonthlyTrend float64 `json:"monthly_trend"`
	GrowthRate   float64 `json:"growth_rate"`
}

// SummaryStats 摘要統計
type SummaryStats struct {
	TotalTokens       int                  `json:"total_tokens"`
	TotalCost         float64              `json:"total_cost"`
	TotalRounds       int                  `json:"total_rounds"`
	ActivityBreakdown map[ActivityType]int `json:"activity_breakdown"`
	TimeRange         struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"time_range"`
}

// ReportConfig 報告配置
type ReportConfig struct {
	Format     string `json:"format"` // json, csv, html
	OutputPath string `json:"output_path"`
	TimeRange  struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"time_range"`
	IncludeDetails bool   `json:"include_details"`
	GroupBy        string `json:"group_by"` // activity, date, session
}

// ReportOptions 報告選項
type ReportOptions struct {
	TimeRange           TimeRange `json:"time_range"`
	IncludeTrends       bool      `json:"include_trends"`
	IncludeOptimization bool      `json:"include_optimization"`
	GroupBy             string    `json:"group_by"`
}

// CostTrendAnalysis 成本趨勢分析
type CostTrendAnalysis struct {
	TimeRange   string           `json:"time_range"`
	DataPoints  []CostDataPoint  `json:"data_points"`
	TotalCost   float64          `json:"total_cost"`
	AverageCost float64          `json:"average_cost"`
	GrowthRate  float64          `json:"growth_rate"`
	Predictions []CostPrediction `json:"predictions"`
}

// CostDataPoint 成本資料點
type CostDataPoint struct {
	Timestamp   time.Time `json:"timestamp"`
	Cost        float64   `json:"cost"`
	TokenCount  int       `json:"token_count"`
	RecordCount int       `json:"record_count"`
}

// CostPrediction 成本預測
type CostPrediction struct {
	Date          time.Time `json:"date"`
	PredictedCost float64   `json:"predicted_cost"`
	Confidence    float64   `json:"confidence"`
}

// CostReport 成本報告
type CostReport struct {
	GeneratedAt  time.Time                    `json:"generated_at"`
	TimeRange    TimeRange                    `json:"time_range"`
	TotalRecords int                          `json:"total_records"`
	Summary      CostSummary                  `json:"summary"`
	ByActivity   map[ActivityType]CostSummary `json:"by_activity"`
	ByModel      map[string]CostSummary       `json:"by_model"`
	Optimization *OptimizationSuggestions     `json:"optimization"`
	Trends       *CostTrendAnalysis           `json:"trends"`
}

// CostSummary 成本摘要
type CostSummary struct {
	TotalCost            float64 `json:"total_cost"`
	TotalTokens          int     `json:"total_tokens"`
	RecordCount          int     `json:"record_count"`
	AverageCostPerRecord float64 `json:"average_cost_per_record"`
	AverageCostPerToken  float64 `json:"average_cost_per_token"` // USD per 1M tokens
}

// CostEfficiencyAnalysis 成本效率分析
type CostEfficiencyAnalysis struct {
	OverallEfficiency float64                  `json:"overall_efficiency"`
	ByActivity        map[ActivityType]float64 `json:"by_activity"`
	ByModel           map[string]float64       `json:"by_model"`
	Recommendations   []string                 `json:"recommendations"`
}

// StringToActivityType 將字串轉換為 ActivityType
func StringToActivityType(s string) ActivityType {
	switch s {
	case "coding":
		return ActivityCoding
	case "debugging":
		return ActivityDebugging
	case "documentation":
		return ActivityDocumentation
	case "spec-development":
		return ActivitySpecDev
	case "chat":
		return ActivityChat
	default:
		return ActivityChat // 預設為聊天類型
	}
}

// String 返回 ActivityType 的字串表示
func (at ActivityType) String() string {
	return string(at)
}

// BasicReport 基礎報告
type BasicReport struct {
	GeneratedAt  time.Time                       `json:"generated_at"`
	TotalRecords int                             `json:"total_records"`
	TimeRange    TimeRange                       `json:"time_range"`
	Summary      ReportSummary                   `json:"summary"`
	ByActivity   map[ActivityType]ActivityReport `json:"by_activity"`
	Statistics   ReportStatistics                `json:"statistics"`
}

// ReportSummary 報告摘要
type ReportSummary struct {
	TotalActivities          int                  `json:"total_activities"`
	TotalTokens              TokenUsage           `json:"total_tokens"`
	ActivityCounts           map[ActivityType]int `json:"activity_counts"`
	AverageTokensPerActivity float64              `json:"average_tokens_per_activity"`
}

// ActivityReport 活動報告
type ActivityReport struct {
	ActivityType  ActivityType `json:"activity_type"`
	Count         int          `json:"count"`
	Tokens        TokenUsage   `json:"tokens"`
	AverageTokens float64      `json:"average_tokens"`
	Percentage    float64      `json:"percentage"`
}

// ReportStatistics 報告統計
type ReportStatistics struct {
	TokenDistribution TokenDistributionStats  `json:"token_distribution"`
	ActivityTrends    ActivityTrends          `json:"activity_trends"`
	EfficiencyMetrics ReportEfficiencyMetrics `json:"efficiency_metrics"`
}

// TokenDistributionStats Token 分佈統計
type TokenDistributionStats struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
	Min     int     `json:"min"`
	Max     int     `json:"max"`
	Median  float64 `json:"median"`
	Input   int     `json:"input"`
	Output  int     `json:"output"`
}

// ActivityTrends 活動趨勢
type ActivityTrends struct {
	HourlyDistribution map[int]int `json:"hourly_distribution"`
	PeakHour           int         `json:"peak_hour"`
	PeakHourCount      int         `json:"peak_hour_count"`
}

// ReportEfficiencyMetrics 報告效率指標
type ReportEfficiencyMetrics struct {
	TokensPerActivity map[ActivityType]float64 `json:"tokens_per_activity"`
}

// ReportResult 報告結果
type ReportResult struct {
	BasicReport *BasicReport  `json:"basic_report"`
	JSONData    []byte        `json:"json_data,omitempty"`
	TextReport  string        `json:"text_report,omitempty"`
	GeneratedAt time.Time     `json:"generated_at"`
	Options     ReportOptions `json:"options"`
}
