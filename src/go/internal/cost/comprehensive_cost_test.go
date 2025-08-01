package cost

import (
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestComprehensiveCostCalculation 綜合成本計算測試
func TestComprehensiveCostCalculation(t *testing.T) {
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
			name:         "Claude Sonnet 4.0 - 精確計算",
			inputTokens:  1000,
			outputTokens: 2000,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.033, // (1000 * 3 + 2000 * 15) / 1000000 = 0.033
			tolerance:    0.001,
		},
		{
			name:         "Claude Haiku 3.5 - 低成本模型",
			inputTokens:  5000,
			outputTokens: 3000,
			model:        "claude-haiku-3.5",
			expectedCost: 0.016, // (5000 * 0.8 + 3000 * 4) / 1000000 = 0.016
			tolerance:    0.001,
		},
		{
			name:         "Claude Opus 4.0 - 高成本模型",
			inputTokens:  500,
			outputTokens: 1000,
			model:        "claude-opus-4.0",
			expectedCost: 0.0825, // (500 * 15 + 1000 * 75) / 1000000 = 0.0825
			tolerance:    0.001,
		},
		{
			name:         "零 Token 輸入",
			inputTokens:  0,
			outputTokens: 1000,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.015, // (0 * 3 + 1000 * 15) / 1000000 = 0.015
			tolerance:    0.001,
		},
		{
			name:         "零 Token 輸出",
			inputTokens:  1000,
			outputTokens: 0,
			model:        "claude-sonnet-4.0",
			expectedCost: 0.003, // (1000 * 3 + 0 * 15) / 1000000 = 0.003
			tolerance:    0.001,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			breakdown, err := calculator.CalculateCost(tc.inputTokens, tc.outputTokens, tc.model)
			if err != nil {
				t.Fatalf("計算成本時發生錯誤: %v", err)
			}

			if breakdown == nil {
				t.Fatal("成本分解結果為 nil")
			}

			diff := breakdown.TotalCost - tc.expectedCost
			if diff < 0 {
				diff = -diff
			}

			if diff > tc.tolerance {
				t.Errorf("成本計算不準確: 期望 %f，實際 %f，差異 %f",
					tc.expectedCost, breakdown.TotalCost, diff)
			}

			// 驗證成本分解的正確性
			expectedInputCost := float64(tc.inputTokens) * breakdown.CostDetails.InputRate / 1000000
			expectedOutputCost := float64(tc.outputTokens) * breakdown.CostDetails.OutputRate / 1000000

			if breakdown.InputCost != expectedInputCost {
				t.Errorf("輸入成本不正確: 期望 %f，實際 %f",
					expectedInputCost, breakdown.InputCost)
			}

			if breakdown.OutputCost != expectedOutputCost {
				t.Errorf("輸出成本不正確: 期望 %f，實際 %f",
					expectedOutputCost, breakdown.OutputCost)
			}
		})
	}
}

// TestPricingModelSupport 測試不同定價模型的支援
func TestPricingModelSupport(t *testing.T) {
	calculator := NewCostCalculator()

	supportedModels := []string{
		"claude-sonnet-4.0",
		"claude-haiku-3.5",
		"claude-opus-4.0",
	}

	for _, model := range supportedModels {
		t.Run("模型支援測試: "+model, func(t *testing.T) {
			// 測試定價資訊獲取
			pricing, err := calculator.GetPricingInfo(model)
			if err != nil {
				t.Fatalf("獲取定價資訊失敗: %v", err)
			}

			if pricing == nil {
				t.Fatal("定價資訊為 nil")
			}

			if pricing.Name != model {
				t.Errorf("模型名稱不匹配: 期望 %s，實際 %s", model, pricing.Name)
			}

			// 驗證定價資訊的合理性
			if pricing.InputPrice <= 0 {
				t.Error("輸入價格應該大於 0")
			}

			if pricing.OutputPrice <= 0 {
				t.Error("輸出價格應該大於 0")
			}

			if pricing.OutputPrice <= pricing.InputPrice {
				t.Error("輸出價格通常應該高於輸入價格")
			}

			// 測試成本計算
			breakdown, err := calculator.CalculateCost(1000, 1000, model)
			if err != nil {
				t.Fatalf("成本計算失敗: %v", err)
			}

			if breakdown.PricingModel != model {
				t.Errorf("成本分解中的模型名稱不匹配: 期望 %s，實際 %s",
					model, breakdown.PricingModel)
			}
		})
	}

	// 測試不支援的模型
	t.Run("不支援的模型", func(t *testing.T) {
		_, err := calculator.GetPricingInfo("unsupported-model")
		if err == nil {
			t.Error("期望不支援的模型會返回錯誤")
		}

		breakdown, err := calculator.CalculateCost(1000, 1000, "unsupported-model")
		if err == nil {
			t.Error("期望不支援的模型會返回錯誤")
		}

		if breakdown != nil {
			t.Error("不支援的模型應該返回 nil breakdown")
		}
	})
}

// TestCostOptimizationSuggestionsAccuracy 驗證成本優化建議的正確性
func TestCostOptimizationSuggestionsAccuracy(t *testing.T) {
	calculator := NewCostCalculator()

	// 建立測試使用記錄
	now := time.Now()
	records := []types.UsageRecord{
		{
			Timestamp: now,
			SessionID: "session-1",
			Activity: types.Activity{
				Type: types.ActivityCoding,
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
				Input:        0.003,
				Output:       0.030,
				Total:        0.033,
				Currency:     "USD",
				PricingModel: "claude-sonnet-4.0",
			},
		},
		{
			Timestamp: now.Add(1 * time.Hour),
			SessionID: "session-1",
			Activity: types.Activity{
				Type: types.ActivityCoding,
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
				Input:        0.0045,
				Output:       0.0375,
				Total:        0.042,
				Currency:     "USD",
				PricingModel: "claude-sonnet-4.0",
			},
		},
		// 重複的內容，適合快取優化
		{
			Timestamp: now.Add(2 * time.Hour),
			SessionID: "session-2",
			Activity: types.Activity{
				Type: types.ActivityCoding,
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
				Input:        0.003,
				Output:       0.030,
				Total:        0.033,
				Currency:     "USD",
				PricingModel: "claude-sonnet-4.0",
			},
		},
	}

	suggestions, err := calculator.CalculateOptimizationSavings(records)
	if err != nil {
		t.Fatalf("計算優化建議失敗: %v", err)
	}

	if suggestions == nil {
		t.Fatal("優化建議為 nil")
	}

	// 驗證基本統計
	expectedCurrentCost := 0.033 + 0.042 + 0.033 // 總成本
	tolerance := 0.001
	diff := suggestions.CurrentCost - expectedCurrentCost
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("當前成本計算不正確: 期望 %f，實際 %f，差異 %f",
			expectedCurrentCost, suggestions.CurrentCost, diff)
	}

	// 驗證建議的合理性（可能為空，這是正常的）
	// 如果有建議，驗證其合理性
	if len(suggestions.Suggestions) == 0 {
		t.Log("沒有生成優化建議，這可能是正常的（取決於使用模式）")
	}

	for _, suggestion := range suggestions.Suggestions {
		if suggestion.PotentialSaving < 0 {
			t.Errorf("潛在節省不應該為負數: %f", suggestion.PotentialSaving)
		}

		if suggestion.Confidence < 0 || suggestion.Confidence > 1 {
			t.Errorf("信心度應該在 0-1 之間: %f", suggestion.Confidence)
		}

		if suggestion.Description == "" {
			t.Error("建議描述不應該為空")
		}

		if suggestion.Type == "" {
			t.Error("建議類型不應該為空")
		}
	}

	// 驗證總節省計算
	calculatedTotalSavings := 0.0
	for _, suggestion := range suggestions.Suggestions {
		calculatedTotalSavings += suggestion.PotentialSaving
	}

	if suggestions.TotalSavings != calculatedTotalSavings {
		t.Errorf("總節省計算不正確: 期望 %f，實際 %f",
			calculatedTotalSavings, suggestions.TotalSavings)
	}

	// 驗證優化後成本
	expectedOptimizedCost := suggestions.CurrentCost - suggestions.TotalSavings
	if suggestions.OptimizedCost != expectedOptimizedCost {
		t.Errorf("優化後成本計算不正確: 期望 %f，實際 %f",
			expectedOptimizedCost, suggestions.OptimizedCost)
	}
}

// TestEdgeCasesAndErrorHandling 測試邊界情況和錯誤處理
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	calculator := NewCostCalculator()

	t.Run("負數 Token 處理", func(t *testing.T) {
		breakdown, err := calculator.CalculateCost(-100, 1000, "claude-sonnet-4.0")
		if err == nil {
			t.Error("負數輸入 Token 應該返回錯誤")
		}
		if breakdown != nil {
			t.Error("錯誤情況下應該返回 nil breakdown")
		}

		breakdown, err = calculator.CalculateCost(1000, -100, "claude-sonnet-4.0")
		if err == nil {
			t.Error("負數輸出 Token 應該返回錯誤")
		}
		if breakdown != nil {
			t.Error("錯誤情況下應該返回 nil breakdown")
		}
	})

	t.Run("極大數值處理", func(t *testing.T) {
		largeTokens := 1000000000 // 10億 tokens
		breakdown, err := calculator.CalculateCost(largeTokens, largeTokens, "claude-sonnet-4.0")
		if err != nil {
			t.Fatalf("處理大數值時發生錯誤: %v", err)
		}

		if breakdown.TotalCost <= 0 {
			t.Error("大數值計算結果應該大於 0")
		}

		// 驗證計算結果的合理性
		expectedCost := float64(largeTokens)*3/1000000 + float64(largeTokens)*15/1000000
		tolerance := expectedCost * 0.01 // 1% 容差

		diff := breakdown.TotalCost - expectedCost
		if diff < 0 {
			diff = -diff
		}

		if diff > tolerance {
			t.Errorf("大數值計算不準確: 期望 %f，實際 %f", expectedCost, breakdown.TotalCost)
		}
	})

	t.Run("空模型名稱處理", func(t *testing.T) {
		breakdown, err := calculator.CalculateCost(1000, 1000, "")
		if err != nil {
			t.Fatalf("空模型名稱處理失敗: %v", err)
		}

		// 應該使用預設模型
		if breakdown.PricingModel == "" {
			t.Error("應該設定預設模型")
		}
	})

	t.Run("空使用記錄優化建議", func(t *testing.T) {
		emptyRecords := []types.UsageRecord{}
		suggestions, err := calculator.CalculateOptimizationSavings(emptyRecords)
		if err != nil {
			t.Fatalf("處理空記錄時發生錯誤: %v", err)
		}

		if suggestions.CurrentCost != 0 {
			t.Error("空記錄的當前成本應該為 0")
		}

		if suggestions.TotalSavings != 0 {
			t.Error("空記錄的總節省應該為 0")
		}
	})
}

// TestConcurrentCostCalculation 測試並發成本計算
func TestConcurrentCostCalculation(t *testing.T) {
	calculator := NewCostCalculator()

	// 並發測試
	const numGoroutines = 10
	const numCalculations = 100

	results := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			for j := 0; j < numCalculations; j++ {
				inputTokens := (goroutineID + 1) * 100
				outputTokens := (goroutineID + 1) * 200

				breakdown, err := calculator.CalculateCost(inputTokens, outputTokens, "claude-sonnet-4.0")
				if err != nil {
					results <- err
					return
				}

				if breakdown == nil {
					results <- err
					return
				}

				// 驗證計算結果的一致性
				expectedCost := float64(inputTokens)*3/1000000 + float64(outputTokens)*15/1000000
				if breakdown.TotalCost != expectedCost {
					results <- err
					return
				}
			}
			results <- nil
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Errorf("並發計算失敗: %v", err)
		}
	}
}
