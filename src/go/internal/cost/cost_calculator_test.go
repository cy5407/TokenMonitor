package cost

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"token-monitor/internal/types"
	"gopkg.in/yaml.v3"
)

// TestNewCostCalculator 測試創建新的成本計算器
func TestNewCostCalculator(t *testing.T) {
	calculator := NewCostCalculator()
	
	if calculator == nil {
		t.Fatal("Expected non-nil calculator")
	}
	
	if calculator.pricingEngine == nil {
		t.Fatal("Expected non-nil pricing engine")
	}
	
	if calculator.sessionCosts == nil {
		t.Fatal("Expected non-nil session costs map")
	}
	
	if calculator.dailyCosts == nil {
		t.Fatal("Expected non-nil daily costs map")
	}
}

// TestCalculateCost 測試基本成本計算
func TestCalculateCost(t *testing.T) {
	calculator := NewCostCalculator()
	
	tests := []struct {
		name         string
		inputTokens  int
		outputTokens int
		model        string
		expectError  bool
		expectedCost float64
	}{
		{
			name:         "Valid calculation",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			expectError:  false,
			expectedCost: 0.0105, // (1000/1M * 3.0) + (500/1M * 15.0)
		},
		{
			name:         "Empty model uses default",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "",
			expectError:  false,
			expectedCost: 0.0105,
		},
		{
			name:         "Negative input tokens",
			inputTokens:  -100,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			expectError:  true,
		},
		{
			name:         "Negative output tokens",
			inputTokens:  1000,
			outputTokens: -500,
			model:        "claude-sonnet-4.0",
			expectError:  true,
		},
		{
			name:         "Invalid model",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "invalid-model",
			expectError:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.CalculateCost(tt.inputTokens, tt.outputTokens, tt.model)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Fatal("Expected non-nil result")
			}
			
			if abs(result.TotalCost-tt.expectedCost) > 0.0001 {
				t.Errorf("Expected cost %.6f, got %.6f", tt.expectedCost, result.TotalCost)
			}
		})
	}
}

// TestCalculateCostWithOptions 測試帶選項的成本計算
func TestCalculateCostWithOptions(t *testing.T) {
	calculator := NewCostCalculator()
	
	tests := []struct {
		name         string
		inputTokens  int
		outputTokens int
		model        string
		options      *CostOptions
		expectError  bool
	}{
		{
			name:         "Standard billing",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:      StandardBilling,
				SessionID: "test-session",
			},
			expectError: false,
		},
		{
			name:         "Cache billing",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:             CacheBilling,
				CacheReadTokens:  200,
				CacheWriteTokens: 100,
				SessionID:        "test-session",
			},
			expectError: false,
		},
		{
			name:         "Batch billing",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:      BatchBilling,
				IsBatch:   true,
				SessionID: "test-session",
			},
			expectError: false,
		},
		{
			name:         "Invalid cache tokens",
			inputTokens:  1000,
			outputTokens: 500,
			model:        "claude-sonnet-4.0",
			options: &CostOptions{
				Mode:             CacheBilling,
				CacheReadTokens:  -100,
				CacheWriteTokens: 100,
			},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.CalculateCostWithOptions(tt.inputTokens, tt.outputTokens, tt.model, tt.options)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Fatal("Expected non-nil result")
			}
			
			// 驗證會話成本追蹤
			if tt.options.SessionID != "" {
				sessionCost := calculator.GetSessionCost(tt.options.SessionID)
				if sessionCost == 0 {
					t.Errorf("Expected session cost to be tracked")
				}
			}
		})
	}
}

// TestGetPricingInfo 測試取得定價資訊
func TestGetPricingInfo(t *testing.T) {
	calculator := NewCostCalculator()
	
	tests := []struct {
		name        string
		model       string
		expectError bool
	}{
		{
			name:        "Valid model",
			model:       "claude-sonnet-4.0",
			expectError: false,
		},
		{
			name:        "Invalid model",
			model:       "invalid-model",
			expectError: true,
		},
		{
			name:        "Empty model",
			model:       "",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.GetPricingInfo(tt.model)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Fatal("Expected non-nil result")
			}
			
			if result.Name != tt.model {
				t.Errorf("Expected model name %s, got %s", tt.model, result.Name)
			}
		})
	}
}

// TestLoadPricingModels 測試載入定價模型
func TestLoadPricingModels(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 創建臨時配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.yaml")
	
	config := map[string]interface{}{
		"pricing": map[string]interface{}{
			"test-model": map[string]interface{}{
				"input":          2.0,
				"output":         10.0,
				"cache_read":     0.2,
				"cache_write":    2.5,
				"batch_discount": 0.4,
			},
		},
	}
	
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	
	// 測試載入
	err = calculator.LoadPricingModels(configFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// 驗證模型是否載入
	models := calculator.GetSupportedModels()
	found := false
	for _, model := range models {
		if model == "test-model" {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected test-model to be loaded")
	}
	
	// 測試不存在的文件
	err = calculator.LoadPricingModels("nonexistent.yaml")
	if err != nil {
		t.Errorf("Expected no error for nonexistent file, got: %v", err)
	}
}

// TestCalculateOptimizationSavings 測試計算優化節省
func TestCalculateOptimizationSavings(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 創建測試記錄
	records := []types.UsageRecord{
		{
			Timestamp: time.Now(),
			SessionID: "session1",
			Activity: types.Activity{
				Type:        types.ActivityCoding,
				Description: "Test coding",
				Rounds:      3,
			},
			Tokens: struct {
				Input             int    `json:"input"`
				Output            int    `json:"output"`
				Total             int    `json:"total"`
				CalculationMethod string `json:"calculation_method"`
			}{
				Input:  1000,
				Output: 500,
				Total:  1500,
			},
			Cost: struct {
				Input        float64 `json:"input"`
				Output       float64 `json:"output"`
				Total        float64 `json:"total"`
				Currency     string  `json:"currency"`
				PricingModel string  `json:"pricing_model"`
			}{
				Total: 0.0105,
			},
		},
	}
	
	result, err := calculator.CalculateOptimizationSavings(records)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	
	// 測試空記錄
	emptyResult, err := calculator.CalculateOptimizationSavings([]types.UsageRecord{})
	if err != nil {
		t.Errorf("Unexpected error for empty records: %v", err)
	}
	
	if emptyResult.TotalSavings != 0 {
		t.Errorf("Expected zero savings for empty records")
	}
}

// TestGetSupportedModels 測試取得支援的模型
func TestGetSupportedModels(t *testing.T) {
	calculator := NewCostCalculator()
	
	models := calculator.GetSupportedModels()
	
	if len(models) == 0 {
		t.Errorf("Expected at least one supported model")
	}
	
	// 檢查預設模型是否存在
	found := false
	for _, model := range models {
		if model == "claude-sonnet-4.0" {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected claude-sonnet-4.0 to be in supported models")
	}
}

// TestSessionAndDailyCosts 測試會話和每日成本追蹤
func TestSessionAndDailyCosts(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 測試會話成本
	options := &CostOptions{
		Mode:      StandardBilling,
		SessionID: "test-session",
	}
	
	_, err := calculator.CalculateCostWithOptions(1000, 500, "claude-sonnet-4.0", options)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	sessionCost := calculator.GetSessionCost("test-session")
	if sessionCost == 0 {
		t.Errorf("Expected non-zero session cost")
	}
	
	// 測試每日成本
	today := time.Now().Format("2006-01-02")
	dailyCost := calculator.GetDailyCost(today)
	if dailyCost == 0 {
		t.Errorf("Expected non-zero daily cost")
	}
	
	// 測試每日成本摘要
	summary := calculator.GetDailyCostSummary(7)
	if len(summary) != 7 {
		t.Errorf("Expected 7 days in summary, got %d", len(summary))
	}
	
	// 測試清除成本
	calculator.ClearSessionCosts()
	clearedSessionCost := calculator.GetSessionCost("test-session")
	if clearedSessionCost != 0 {
		t.Errorf("Expected zero session cost after clearing")
	}
	
	calculator.ClearDailyCosts()
	clearedDailyCost := calculator.GetDailyCost(today)
	if clearedDailyCost != 0 {
		t.Errorf("Expected zero daily cost after clearing")
	}
}

// TestEstimateMonthlyBudget 測試估算月度預算
func TestEstimateMonthlyBudget(t *testing.T) {
	calculator := NewCostCalculator()
	
	budget, err := calculator.EstimateMonthlyBudget(10000, "claude-sonnet-4.0")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if budget <= 0 {
		t.Errorf("Expected positive budget, got %f", budget)
	}
	
	// 測試無效模型
	_, err = calculator.EstimateMonthlyBudget(10000, "invalid-model")
	if err == nil {
		t.Errorf("Expected error for invalid model")
	}
}

// TestComparePricingModels 測試比較定價模型
func TestComparePricingModels(t *testing.T) {
	calculator := NewCostCalculator()
	
	comparison, err := calculator.ComparePricingModels(1000, 500)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	if len(comparison) == 0 {
		t.Errorf("Expected at least one model in comparison")
	}
	
	// 驗證包含預設模型
	if _, exists := comparison["claude-sonnet-4.0"]; !exists {
		t.Errorf("Expected claude-sonnet-4.0 in comparison")
	}
}

// TestReloadConfig 測試重新載入配置
func TestReloadConfig(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 測試未設定配置路徑
	err := calculator.ReloadConfig()
	if err == nil {
		t.Errorf("Expected error when no config path set")
	}
	
	// 設定配置路徑後測試
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	
	config := map[string]interface{}{
		"pricing": map[string]interface{}{
			"test-model": map[string]interface{}{
				"input":          1.0,
				"output":         5.0,
				"cache_read":     0.1,
				"cache_write":    1.25,
				"batch_discount": 0.3,
			},
		},
	}
	
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	
	err = calculator.LoadPricingModels(configFile)
	if err != nil {
		t.Errorf("Unexpected error loading config: %v", err)
	}
	
	err = calculator.ReloadConfig()
	if err != nil {
		t.Errorf("Unexpected error reloading config: %v", err)
	}
}

// TestGetStatistics 測試取得統計資訊
func TestGetStatistics(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 添加一些測試資料
	options := &CostOptions{
		Mode:      StandardBilling,
		SessionID: "test-session",
	}
	
	_, err := calculator.CalculateCostWithOptions(1000, 500, "claude-sonnet-4.0", options)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	stats := calculator.GetStatistics()
	
	expectedKeys := []string{
		"total_sessions",
		"total_session_cost",
		"total_days",
		"total_daily_cost",
		"supported_models",
		"last_config_update",
		"average_session_cost",
		"average_daily_cost",
	}
	
	for _, key := range expectedKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("Expected key %s in statistics", key)
		}
	}
}

// TestConcurrentAccess 測試併發存取
func TestConcurrentAccess(t *testing.T) {
	calculator := NewCostCalculator()
	
	// 併發測試
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			// 併發計算成本
			_, err := calculator.CalculateCost(1000, 500, "claude-sonnet-4.0")
			if err != nil {
				t.Errorf("Concurrent access error: %v", err)
			}
			
			// 併發取得定價資訊
			_, err = calculator.GetPricingInfo("claude-sonnet-4.0")
			if err != nil {
				t.Errorf("Concurrent pricing info error: %v", err)
			}
		}(i)
	}
	
	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}
}

// abs 計算絕對值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// BenchmarkCalculateCost 基準測試成本計算
func BenchmarkCalculateCost(b *testing.B) {
	calculator := NewCostCalculator()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := calculator.CalculateCost(1000, 500, "claude-sonnet-4.0")
		if err != nil {
			b.Errorf("Benchmark error: %v", err)
		}
	}
}

// BenchmarkCalculateCostWithOptions 基準測試帶選項的成本計算
func BenchmarkCalculateCostWithOptions(b *testing.B) {
	calculator := NewCostCalculator()
	
	options := &CostOptions{
		Mode:      StandardBilling,
		SessionID: "bench-session",
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := calculator.CalculateCostWithOptions(1000, 500, "claude-sonnet-4.0", options)
		if err != nil {
			b.Errorf("Benchmark error: %v", err)
		}
	}
}