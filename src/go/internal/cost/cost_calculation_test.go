package cost

import (
	"testing"
	"token-monitor/internal/types"
)

// TestCostCalculationAccuracy 測試定價計算的準確性
func TestCostCalculationAccuracy(t *testing.T) {
	calculator := NewCostCalculator()

	testCases := []struct {
		name         string
		inputTokens  int
		outputTokens int
		model        string
		expectedCost float64
		tolerance    float64
	}{
		{
			name:         "Claude Sonnet 4.0 - 小量 tokens",
			inputTokens:  100,
			outputTokens: 200,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.0033, // (100/1M * $3) + (200/1M * $15) = $0.0003 + $0.003 = $0.0033
			tolerance:    0.000001,
		},
		{
			name:         "Claude Sonnet 4.0 - 中量 tokens",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.033, // (1000/1M * $3) + (2000/1M * $15) = $0.003 + $0.03 = $0.033
			tolerance:    0.000001,
		},
		{
			name:         "Claude Sonnet 4.0 - 大量 tokens",
			inputTokens:  10000,
			outputTokens: 20000,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.33, // (10000/1M * $3) + (20000/1M * $15) = $0.03 + $0.3 = $0.33
			tolerance:    0.000001,
		},
		{
			name:         "Claude Haiku 3.5 - 測試不同模型",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-haiku-3.5",
			expectedCost: 0.0088, // (1000/1M * $0.8) + (2000/1M * $4) = $0.0008 + $0.008 = $0.0088
			tolerance:    0.000001,
		},
		{
			name:         "Claude Opus 4.0 - 高價模型",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-opus-4.0",
			expectedCost: 0.165, // (1000/1M * $15) + (2000/1M * $75) = $0.015 + $0.15 = $0.165
			tolerance:    0.000001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculator.CalculateCost(tc.inputTokens, tc.outputTokens, tc.model)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if absFloat64(result.TotalCost-tc.expectedCost) > tc.tolerance {
				t.Errorf("Expected cost %.6f, got %.6f (difference: %.6f)",
					tc.expectedCost, result.TotalCost, absFloat64(result.TotalCost-tc.expectedCost))
			}

			// 驗證成本分解
			expectedInputCost := float64(tc.inputTokens) / 1_000_000 * result.CostDetails.InputRate
			expectedOutputCost := float64(tc.outputTokens) / 1_000_000 * result.CostDetails.OutputRate

			if absFloat64(result.InputCost-expectedInputCost) > tc.tolerance {
				t.Errorf("Expected input cost %.6f, got %.6f", expectedInputCost, result.InputCost)
			}

			if absFloat64(result.OutputCost-expectedOutputCost) > tc.tolerance {
				t.Errorf("Expected output cost %.6f, got %.6f", expectedOutputCost, result.OutputCost)
			}
		})
	}
}

// TestDifferentPricingModels 測試不同定價模型的支援
func TestDifferentPricingModels(t *testing.T) {
	calculator := NewCostCalculator()

	// 測試所有支援的模型
	supportedModels := calculator.GetSupportedModels()
	if len(supportedModels) == 0 {
		t.Error("Expected at least one supported model")
		return
	}

	inputTokens := 1000
	outputTokens := 2000

	for _, model := range supportedModels {
		t.Run(model, func(t *testing.T) {
			// 測試基本成本計算
			result, err := calculator.CalculateCost(inputTokens, outputTokens, model)
			if err != nil {
				t.Errorf("Failed to calculate cost for model %s: %v", model, err)
				return
			}

			// 驗證結果結構
			if result.PricingModel != model {
				t.Errorf("Expected pricing model %s, got %s", model, result.PricingModel)
			}

			if result.Currency != "USD" {
				t.Errorf("Expected currency USD, got %s", result.Currency)
			}

			if result.TotalCost <= 0 {
				t.Errorf("Expected positive total cost, got %.6f", result.TotalCost)
			}

			if result.InputCost <= 0 {
				t.Errorf("Expected positive input cost, got %.6f", result.InputCost)
			}

			if result.OutputCost <= 0 {
				t.Errorf("Expected positive output cost, got %.6f", result.OutputCost)
			}

			// 驗證成本分解的一致性
			expectedTotal := result.InputCost + result.OutputCost
			if absFloat64(result.TotalCost-expectedTotal) > 0.000001 {
				t.Errorf("Cost breakdown inconsistent: input(%.6f) + output(%.6f) = %.6f, but total is %.6f",
					result.InputCost, result.OutputCost, expectedTotal, result.TotalCost)
			}

			// 測試定價資訊
			pricingInfo, err := calculator.GetPricingInfo(model)
			if err != nil {
				t.Errorf("Failed to get pricing info for model %s: %v", model, err)
				return
			}

			if pricingInfo.Name != model {
				t.Errorf("Expected pricing model name %s, got %s", model, pricingInfo.Name)
			}

			if pricingInfo.InputPrice <= 0 {
				t.Errorf("Expected positive input price, got %.2f", pricingInfo.InputPrice)
			}

			if pricingInfo.OutputPrice <= 0 {
				t.Errorf("Expected positive output price, got %.2f", pricingInfo.OutputPrice)
			}
		})
	}
}

// TestCostOptimizationSuggestions 測試成本優化建議的正確性
func TestCostOptimizationSuggestions(t *testing.T) {
	calculator := NewCostCalculator()

	// 創建測試使用記錄
	records := []types.UsageRecord{
		{
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
	}

	suggestions, err := calculator.CalculateOptimizationSavings(records)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if suggestions == nil {
		t.Error("Expected optimization suggestions but got nil")
		return
	}

	// 驗證建議結構（允許為空，因為優化器可能沒有找到優化機會）
	if suggestions.Suggestions == nil {
		// 初始化空的建議數組
		suggestions.Suggestions = []types.OptimizationSuggestion{}
		t.Log("Suggestions array was nil, initialized to empty array")
	}

	if suggestions.CurrentCost < 0 {
		t.Errorf("Expected non-negative current cost, got %.6f", suggestions.CurrentCost)
	}

	if suggestions.OptimizedCost < 0 {
		t.Errorf("Expected non-negative optimized cost, got %.6f", suggestions.OptimizedCost)
	}

	if suggestions.TotalSavings < 0 {
		t.Errorf("Expected non-negative total savings, got %.6f", suggestions.TotalSavings)
	}

	// 驗證邏輯一致性
	expectedSavings := suggestions.CurrentCost - suggestions.OptimizedCost
	if absFloat64(suggestions.TotalSavings-expectedSavings) > 0.000001 {
		t.Errorf("Savings calculation inconsistent: current(%.6f) - optimized(%.6f) = %.6f, but total savings is %.6f",
			suggestions.CurrentCost, suggestions.OptimizedCost, expectedSavings, suggestions.TotalSavings)
	}
}

// TestCacheBillingAccuracy 測試快取計費的準確性
func TestCacheBillingAccuracy(t *testing.T) {
	calculator := NewCostCalculator()

	inputTokens := 1000
	outputTokens := 2000
	cacheReadTokens := 500
	cacheWriteTokens := 300
	model := "claude-sonnet-4.0"

	options := &CostOptions{
		Mode:             CacheBilling,
		CacheReadTokens:  cacheReadTokens,
		CacheWriteTokens: cacheWriteTokens,
	}

	result, err := calculator.CalculateDetailedCost(inputTokens, outputTokens, model, options)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// 驗證快取成本計算
	// Claude Sonnet 4.0: cache read $0.30/MTok, cache write $3.75/MTok
	expectedCacheReadCost := float64(cacheReadTokens) / 1_000_000 * 0.30
	expectedCacheWriteCost := float64(cacheWriteTokens) / 1_000_000 * 3.75

	if absFloat64(result.CacheReadCost-expectedCacheReadCost) > 0.000001 {
		t.Errorf("Expected cache read cost %.6f, got %.6f", expectedCacheReadCost, result.CacheReadCost)
	}

	if absFloat64(result.CacheWriteCost-expectedCacheWriteCost) > 0.000001 {
		t.Errorf("Expected cache write cost %.6f, got %.6f", expectedCacheWriteCost, result.CacheWriteCost)
	}

	// 驗證總成本包含快取成本
	expectedTotalCost := result.InputCost + result.OutputCost + result.CacheReadCost + result.CacheWriteCost
	if absFloat64(result.TotalCost-expectedTotalCost) > 0.000001 {
		t.Errorf("Expected total cost %.6f, got %.6f", expectedTotalCost, result.TotalCost)
	}

	// 驗證計費模式
	if result.CostDetails.BillingMode != "cache" {
		t.Errorf("Expected billing mode 'cache', got %s", result.CostDetails.BillingMode)
	}
}

// TestBatchBillingAccuracy 測試批次計費的準確性
func TestBatchBillingAccuracy(t *testing.T) {
	calculator := NewCostCalculator()

	inputTokens := 1000
	outputTokens := 2000
	model := "claude-sonnet-4.0"

	options := &CostOptions{
		Mode:    BatchBilling,
		IsBatch: true,
	}

	result, err := calculator.CalculateDetailedCost(inputTokens, outputTokens, model, options)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// 驗證批次折扣
	// Claude Sonnet 4.0 有 50% 的批次折扣
	expectedDiscountRate := 0.5
	if absFloat64(result.CostDetails.DiscountRate-expectedDiscountRate) > 0.000001 {
		t.Errorf("Expected discount rate %.2f, got %.2f", expectedDiscountRate, result.CostDetails.DiscountRate)
	}

	// 計算原始成本（無折扣）
	originalInputCost := float64(inputTokens) / 1_000_000 * 3.0
	originalOutputCost := float64(outputTokens) / 1_000_000 * 15.0
	originalTotalCost := originalInputCost + originalOutputCost

	// 驗證折扣金額
	expectedBatchDiscount := originalTotalCost * expectedDiscountRate
	if absFloat64(result.BatchDiscount-expectedBatchDiscount) > 0.000001 {
		t.Errorf("Expected batch discount %.6f, got %.6f", expectedBatchDiscount, result.BatchDiscount)
	}

	// 驗證折扣後的總成本
	expectedDiscountedCost := originalTotalCost * (1 - expectedDiscountRate)
	if absFloat64(result.TotalCost-expectedDiscountedCost) > 0.000001 {
		t.Errorf("Expected discounted total cost %.6f, got %.6f", expectedDiscountedCost, result.TotalCost)
	}

	// 驗證計費模式
	if result.CostDetails.BillingMode != "batch" {
		t.Errorf("Expected billing mode 'batch', got %s", result.CostDetails.BillingMode)
	}
}

// TestModelComparison 測試模型比較功能
func TestModelComparison(t *testing.T) {
	calculator := NewCostCalculator()

	inputTokens := 1000
	outputTokens := 2000

	comparison, err := calculator.ComparePricingModels(inputTokens, outputTokens)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if len(comparison) == 0 {
		t.Error("Expected at least one model in comparison")
		return
	}

	// 驗證每個模型的比較結果
	for modelName, breakdown := range comparison {
		if breakdown == nil {
			t.Errorf("Expected breakdown for model %s but got nil", modelName)
			continue
		}

		if breakdown.PricingModel != modelName {
			t.Errorf("Expected pricing model %s, got %s", modelName, breakdown.PricingModel)
		}

		if breakdown.TotalCost <= 0 {
			t.Errorf("Expected positive total cost for model %s, got %.6f", modelName, breakdown.TotalCost)
		}

		if breakdown.Currency != "USD" {
			t.Errorf("Expected currency USD for model %s, got %s", modelName, breakdown.Currency)
		}
	}

	// 驗證不同模型的成本確實不同（除非定價相同）
	costs := make([]float64, 0, len(comparison))
	for _, breakdown := range comparison {
		costs = append(costs, breakdown.TotalCost)
	}

	// 檢查是否有不同的成本（至少有一對不同）
	allSame := true
	if len(costs) > 1 {
		firstCost := costs[0]
		for _, cost := range costs[1:] {
			if absFloat64(cost-firstCost) > 0.000001 {
				allSame = false
				break
			}
		}
	}

	// 如果所有成本都相同，可能是因為所有模型定價相同（不太可能但可能）
	if allSame && len(costs) > 1 {
		t.Log("Warning: All models have the same cost, this might indicate an issue")
	}
}

// absFloat64 計算浮點數絕對值
func absFloat64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
