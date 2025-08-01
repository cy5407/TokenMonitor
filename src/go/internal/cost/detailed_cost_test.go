package cost

import (
	"testing"
	"time"
	"token-monitor/internal/types"
)

// TestCalculateDetailedCost 測試詳細成本計算
func TestCalculateDetailedCost(t *testing.T) {
	calculator := NewCostCalculator()

	testCases := []struct {
		name         string
		inputTokens  int
		outputTokens int
		model        string
		options      *CostOptions
		expectError  bool
		validateFunc func(*testing.T, *types.CostBreakdown)
	}{
		{
			name:         "標準計費模式",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:         StandardBilling,
				SessionID:    "test-session-1",
				ActivityType: types.ActivityCoding,
			},
			expectError: false,
			validateFunc: func(t *testing.T, breakdown *types.CostBreakdown) {
				if breakdown.TokenCounts.Input != 1000 {
					t.Errorf("Expected input tokens 1000, got %d", breakdown.TokenCounts.Input)
				}
				if breakdown.TokenCounts.Output != 2000 {
					t.Errorf("Expected output tokens 2000, got %d", breakdown.TokenCounts.Output)
				}
				if breakdown.TokenCounts.Total != 3000 {
					t.Errorf("Expected total tokens 3000, got %d", breakdown.TokenCounts.Total)
				}
				if breakdown.CostDetails.BillingMode != "standard" {
					t.Errorf("Expected billing mode 'standard', got %s", breakdown.CostDetails.BillingMode)
				}
				if breakdown.SessionID != "test-session-1" {
					t.Errorf("Expected session ID 'test-session-1', got %s", breakdown.SessionID)
				}
				if breakdown.ActivityType != types.ActivityCoding {
					t.Errorf("Expected activity type 'coding', got %s", breakdown.ActivityType)
				}
				// 驗證成本計算：input: 1000/1M * $3 = $0.003, output: 2000/1M * $15 = $0.03
				expectedInputCost := 0.003
				expectedOutputCost := 0.03
				expectedTotalCost := expectedInputCost + expectedOutputCost

				if absFloat(breakdown.InputCost-expectedInputCost) > 0.0001 {
					t.Errorf("Expected input cost %.6f, got %.6f", expectedInputCost, breakdown.InputCost)
				}
				if absFloat(breakdown.OutputCost-expectedOutputCost) > 0.0001 {
					t.Errorf("Expected output cost %.6f, got %.6f", expectedOutputCost, breakdown.OutputCost)
				}
				if absFloat(breakdown.TotalCost-expectedTotalCost) > 0.0001 {
					t.Errorf("Expected total cost %.6f, got %.6f", expectedTotalCost, breakdown.TotalCost)
				}
			},
		},
		{
			name:         "快取計費模式",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:             CacheBilling,
				CacheReadTokens:  500,
				CacheWriteTokens: 300,
				SessionID:        "test-session-2",
				ActivityType:     types.ActivityDebugging,
			},
			expectError: false,
			validateFunc: func(t *testing.T, breakdown *types.CostBreakdown) {
				if breakdown.TokenCounts.CacheRead != 500 {
					t.Errorf("Expected cache read tokens 500, got %d", breakdown.TokenCounts.CacheRead)
				}
				if breakdown.TokenCounts.CacheWrite != 300 {
					t.Errorf("Expected cache write tokens 300, got %d", breakdown.TokenCounts.CacheWrite)
				}
				if breakdown.TokenCounts.Total != 3800 { // 1000 + 2000 + 500 + 300
					t.Errorf("Expected total tokens 3800, got %d", breakdown.TokenCounts.Total)
				}
				if breakdown.CostDetails.BillingMode != "cache" {
					t.Errorf("Expected billing mode 'cache', got %s", breakdown.CostDetails.BillingMode)
				}

				// 驗證快取成本計算
				// Cache read: 500/1M * $0.30 = $0.00015
				// Cache write: 300/1M * $3.75 = $0.001125
				expectedCacheReadCost := 0.00015
				expectedCacheWriteCost := 0.001125

				if absFloat(breakdown.CacheReadCost-expectedCacheReadCost) > 0.000001 {
					t.Errorf("Expected cache read cost %.6f, got %.6f", expectedCacheReadCost, breakdown.CacheReadCost)
				}
				if absFloat(breakdown.CacheWriteCost-expectedCacheWriteCost) > 0.000001 {
					t.Errorf("Expected cache write cost %.6f, got %.6f", expectedCacheWriteCost, breakdown.CacheWriteCost)
				}
			},
		},
		{
			name:         "批次計費模式",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:         BatchBilling,
				IsBatch:      true,
				SessionID:    "test-session-3",
				ActivityType: types.ActivityDocumentation,
			},
			expectError: false,
			validateFunc: func(t *testing.T, breakdown *types.CostBreakdown) {
				if breakdown.CostDetails.BillingMode != "batch" {
					t.Errorf("Expected billing mode 'batch', got %s", breakdown.CostDetails.BillingMode)
				}
				if breakdown.CostDetails.DiscountRate != 0.5 {
					t.Errorf("Expected discount rate 0.5, got %.2f", breakdown.CostDetails.DiscountRate)
				}

				// 驗證批次折扣：50% 折扣
				// 原始成本：input: $0.003, output: $0.03, total: $0.033
				// 折扣後：$0.033 * 0.5 = $0.0165
				expectedTotalCost := 0.0165
				expectedBatchDiscount := 0.0165 // 50% 的折扣金額

				if absFloat(breakdown.TotalCost-expectedTotalCost) > 0.0001 {
					t.Errorf("Expected total cost %.6f, got %.6f", expectedTotalCost, breakdown.TotalCost)
				}
				if absFloat(breakdown.BatchDiscount-expectedBatchDiscount) > 0.0001 {
					t.Errorf("Expected batch discount %.6f, got %.6f", expectedBatchDiscount, breakdown.BatchDiscount)
				}
			},
		},
		{
			name:         "無效輸入 - 負數 tokens",
			inputTokens:  -100,
			outputTokens: 1000,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode: StandardBilling,
			},
			expectError: true,
		},
		{
			name:         "無效模型",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "invalid-model",
			options: &CostOptions{
				Mode: StandardBilling,
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			breakdown, err := calculator.CalculateDetailedCost(tc.inputTokens, tc.outputTokens, tc.model, tc.options)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if breakdown == nil {
				t.Errorf("Expected breakdown but got nil")
				return
			}

			// 基本驗證
			if breakdown.Currency != "USD" {
				t.Errorf("Expected currency 'USD', got %s", breakdown.Currency)
			}
			if breakdown.PricingModel != tc.model {
				t.Errorf("Expected pricing model %s, got %s", tc.model, breakdown.PricingModel)
			}
			if breakdown.Timestamp.IsZero() {
				t.Errorf("Expected non-zero timestamp")
			}

			// 執行自定義驗證
			if tc.validateFunc != nil {
				tc.validateFunc(t, breakdown)
			}
		})
	}
}

// TestAnalyzeCostTrends 測試成本趨勢分析
func TestAnalyzeCostTrends(t *testing.T) {
	calculator := NewCostCalculator()

	// 創建測試資料
	now := time.Now()
	records := []types.UsageRecord{
		{
			Timestamp: now.AddDate(0, 0, -3),
			Activity: types.Activity{
				Type:   types.ActivityCoding,
				Rounds: 5,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  1000,
				Output: 2000,
				Total:  3000,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				PricingModel: "claude-sonnet-4.0",
			},
		},
		{
			Timestamp: now.AddDate(0, 0, -2),
			Activity: types.Activity{
				Type:   types.ActivityDebugging,
				Rounds: 3,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  1500,
				Output: 2500,
				Total:  4000,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				PricingModel: "claude-sonnet-4.0",
			},
		},
		{
			Timestamp: now.AddDate(0, 0, -1),
			Activity: types.Activity{
				Type:   types.ActivityDocumentation,
				Rounds: 2,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  800,
				Output: 1200,
				Total:  2000,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				PricingModel: "claude-sonnet-4.0",
			},
		},
	}

	trends, err := calculator.AnalyzeCostTrends(records, "daily")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if trends == nil {
		t.Errorf("Expected trends but got nil")
		return
	}

	// 驗證基本屬性
	if trends.TimeRange != "daily" {
		t.Errorf("Expected time range 'daily', got %s", trends.TimeRange)
	}

	if len(trends.DataPoints) != 3 {
		t.Errorf("Expected 3 data points, got %d", len(trends.DataPoints))
	}

	if trends.TotalCost <= 0 {
		t.Errorf("Expected positive total cost, got %.6f", trends.TotalCost)
	}

	if trends.AverageCost <= 0 {
		t.Errorf("Expected positive average cost, got %.6f", trends.AverageCost)
	}

	// 驗證預測
	if len(trends.Predictions) == 0 {
		t.Errorf("Expected predictions but got none")
	}

	for _, prediction := range trends.Predictions {
		if prediction.Confidence <= 0 || prediction.Confidence > 1 {
			t.Errorf("Expected confidence between 0 and 1, got %.2f", prediction.Confidence)
		}
	}
}

// TestGenerateCostReport 測試成本報告生成
func TestGenerateCostReport(t *testing.T) {
	calculator := NewCostCalculator()

	// 創建測試資料
	now := time.Now()
	records := []types.UsageRecord{
		{
			Timestamp: now.AddDate(0, 0, -1),
			Activity: types.Activity{
				Type:   types.ActivityCoding,
				Rounds: 5,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  1000,
				Output: 2000,
				Total:  3000,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				PricingModel: "claude-sonnet-4.0",
			},
		},
		{
			Timestamp: now,
			Activity: types.Activity{
				Type:   types.ActivityCoding,
				Rounds: 3,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  1500,
				Output: 2500,
				Total:  4000,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				PricingModel: "claude-sonnet-4.0",
			},
		},
	}

	options := &types.ReportOptions{
		TimeRange: types.TimeRange{
			Start: now.AddDate(0, 0, -7),
			End:   now,
		},
		IncludeTrends:       true,
		IncludeOptimization: true,
		GroupBy:             "activity",
	}

	report, err := calculator.GenerateCostReport(records, options)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if report == nil {
		t.Errorf("Expected report but got nil")
		return
	}

	// 驗證基本屬性
	if report.TotalRecords != 2 {
		t.Errorf("Expected 2 total records, got %d", report.TotalRecords)
	}

	if report.Summary.TotalTokens != 7000 {
		t.Errorf("Expected 7000 total tokens, got %d", report.Summary.TotalTokens)
	}

	if report.Summary.TotalCost <= 0 {
		t.Errorf("Expected positive total cost, got %.6f", report.Summary.TotalCost)
	}

	// 驗證按活動分組
	if len(report.ByActivity) == 0 {
		t.Errorf("Expected activity breakdown but got none")
	}

	codingSummary, exists := report.ByActivity[types.ActivityCoding]
	if !exists {
		t.Errorf("Expected coding activity summary but not found")
	} else {
		if codingSummary.RecordCount != 2 {
			t.Errorf("Expected 2 coding records, got %d", codingSummary.RecordCount)
		}
		if codingSummary.TotalTokens != 7000 {
			t.Errorf("Expected 7000 coding tokens, got %d", codingSummary.TotalTokens)
		}
	}

	// 驗證按模型分組
	if len(report.ByModel) == 0 {
		t.Errorf("Expected model breakdown but got none")
	}

	modelSummary, exists := report.ByModel["claude-sonnet-4.0"]
	if !exists {
		t.Errorf("Expected claude-sonnet-4.0 model summary but not found")
	} else {
		if modelSummary.RecordCount != 2 {
			t.Errorf("Expected 2 model records, got %d", modelSummary.RecordCount)
		}
	}

	// 驗證優化建議
	if report.Optimization == nil {
		t.Errorf("Expected optimization suggestions but got nil")
	}

	// 驗證趨勢分析
	if report.Trends == nil {
		t.Errorf("Expected trends analysis but got nil")
	}
}

// absFloat 計算浮點數絕對值
func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
