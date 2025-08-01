package testutils

import (
	"time"

	"token-monitor/internal/types"
)

// CreateTestReport 建立一個用於測試的 BasicReport 實例
func CreateTestReport() *types.BasicReport {
	return &types.BasicReport{
		GeneratedAt:  time.Now(),
		TotalRecords: 5,
		Summary: types.ReportSummary{
			TotalActivities:          5,
			TotalTokens:              types.TokenUsage{TotalTokens: 1250},
			AverageTokensPerActivity: 250.0,
		},
		ByActivity: map[types.ActivityType]types.ActivityReport{
			types.ActivityCoding: {
				ActivityType:  types.ActivityCoding,
				Count:         3,
				Tokens:        types.TokenUsage{TotalTokens: 750},
				AverageTokens: 250.0,
				Percentage:    60.0,
			},
			types.ActivityDebugging: {
				ActivityType:  types.ActivityDebugging,
				Count:         2,
				Tokens:        types.TokenUsage{TotalTokens: 500},
				AverageTokens: 250.0,
				Percentage:    40.0,
			},
		},
		Statistics: types.ReportStatistics{
			TokenDistribution: types.TokenDistributionStats{
				Total:   1250,
				Average: 250.0,
				Min:     200,
				Max:     300,
				Median:  250.0,
			},
		},
	}
}
