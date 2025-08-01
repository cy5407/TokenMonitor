package cost

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"
	"token-monitor/internal/errors"
	"token-monitor/internal/types"
	"gopkg.in/yaml.v3"
)

// PricingEngine 定價引擎
type PricingEngine struct {
	models           map[string]*types.PricingModel
	mutex            sync.RWMutex
	lastUpdate       time.Time
	defaultModel     string
	validationRules  map[string]ValidationRule
	errorHandler     errors.ErrorHandler
}

// ValidationRule 定價模型驗證規則
type ValidationRule struct {
	MinPrice     float64 `yaml:"min_price"`
	MaxPrice     float64 `yaml:"max_price"`
	RequireCache bool    `yaml:"require_cache"`
	RequireBatch bool    `yaml:"require_batch"`
}

// PricingConfig 定價配置結構
type PricingConfig struct {
	Pricing    map[string]PricingModelConfig `yaml:"pricing"`
	Validation map[string]ValidationRule     `yaml:"validation"`
	Default    string                        `yaml:"default"`
}

// PricingModelConfig 配置文件中的定價模型
type PricingModelConfig struct {
	Input         float64 `yaml:"input"`
	Output        float64 `yaml:"output"`
	CacheRead     float64 `yaml:"cache_read"`
	CacheWrite    float64 `yaml:"cache_write"`
	BatchDiscount float64 `yaml:"batch_discount"`
}

// NewPricingEngine 創建新的定價引擎
func NewPricingEngine() *PricingEngine {
	pe := &PricingEngine{
		models:          make(map[string]*types.PricingModel),
		validationRules: make(map[string]ValidationRule),
		defaultModel:    "claude-sonnet-4.0",
		errorHandler:    errors.NewErrorHandler(),
	}
	
	// 載入預設模型
	pe.LoadDefaultModels()
	
	return pe
}

// AddPricingModel 添加定價模型
func (pe *PricingEngine) AddPricingModel(name string, model *types.PricingModel) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	pe.models[name] = model
	pe.lastUpdate = time.Now()
}

// LoadFromConfig 從配置文件載入定價模型
func (pe *PricingEngine) LoadFromConfig(configPath string) error {
	ctx := context.Background()
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	// 檢查文件存在性
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		warnErr := errors.New(errors.ErrCodeFileNotFound, "配置文件不存在，使用預設模型")
		warnErr = warnErr.WithContext(errors.ErrorContext{
			Operation:  "load_config",
			Component:  "pricing_engine",
			Parameters: map[string]interface{}{
				"config_path": configPath,
			},
		})
		pe.errorHandler.Handle(ctx, warnErr)
		return nil
	}
	
	// 讀取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		appErr := errors.Wrap(err, errors.ErrCodeConfigLoad, "讀取配置文件失敗")
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "read_config_file",
			Component:  "pricing_engine",
			Parameters: map[string]interface{}{
				"config_path": configPath,
			},
		})
		return pe.errorHandler.Handle(ctx, appErr)
	}
	
	var config PricingConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		appErr := errors.Wrap(err, errors.ErrCodeInvalidConfigFormat, "配置文件格式無效")
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "parse_config_file",
			Component:  "pricing_engine",
			Parameters: map[string]interface{}{
				"config_path": configPath,
			},
		})
		return pe.errorHandler.Handle(ctx, appErr)
	}
	
	// 清除現有模型
	pe.models = make(map[string]*types.PricingModel)
	
	// 載入新模型
	for name, modelConfig := range config.Pricing {
		if err := pe.validateModelConfig(name, modelConfig); err != nil {
			warnErr := errors.Wrap(err, errors.ErrCodeConfigValidation, fmt.Sprintf("無效的定價模型: %s", name))
			warnErr = warnErr.WithContext(errors.ErrorContext{
				Operation:  "validate_model_config",
				Component:  "pricing_engine",
				Parameters: map[string]interface{}{
					"model_name": name,
				},
			})
			pe.errorHandler.Handle(ctx, warnErr)
			continue
		}
		
		pe.models[name] = &types.PricingModel{
			Name:          name,
			InputPrice:    modelConfig.Input,
			OutputPrice:   modelConfig.Output,
			CacheRead:     modelConfig.CacheRead,
			CacheWrite:    modelConfig.CacheWrite,
			BatchDiscount: modelConfig.BatchDiscount,
		}
	}
	
	// 載入驗證規則
	pe.validationRules = config.Validation
	
	// 設定預設模型
	if config.Default != "" && pe.models[config.Default] != nil {
		pe.defaultModel = config.Default
	}
	
	pe.lastUpdate = time.Now()
	log.Printf("Loaded %d pricing models from config", len(pe.models))
	
	return nil
}

// GetPricingModel 取得定價模型
func (pe *PricingEngine) GetPricingModel(name string) (*types.PricingModel, error) {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	// 如果名稱為空，使用預設模型
	if name == "" {
		name = pe.defaultModel
	}
	
	model, exists := pe.models[name]
	if !exists {
		return nil, errors.Newf(errors.ErrCodeInvalidPricingModel, "定價模型 '%s' 不存在", name)
	}
	return model, nil
}

// GetSupportedModels 取得支援的模型列表
func (pe *PricingEngine) GetSupportedModels() []string {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	models := make([]string, 0, len(pe.models))
	for name := range pe.models {
		models = append(models, name)
	}
	
	// 排序確保一致性
	sort.Strings(models)
	return models
}

// GetDefaultModel 取得預設模型名稱
func (pe *PricingEngine) GetDefaultModel() string {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	return pe.defaultModel
}

// SetDefaultModel 設定預設模型
func (pe *PricingEngine) SetDefaultModel(modelName string) error {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	
	if _, exists := pe.models[modelName]; !exists {
		return fmt.Errorf("model '%s' does not exist", modelName)
	}
	
	pe.defaultModel = modelName
	return nil
}

// GetLastUpdate 取得最後更新時間
func (pe *PricingEngine) GetLastUpdate() time.Time {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	return pe.lastUpdate
}

// validateModelConfig 驗證定價模型配置
func (pe *PricingEngine) validateModelConfig(name string, config PricingModelConfig) error {
	// 基本驗證
	if config.Input < 0 || config.Output < 0 {
		return fmt.Errorf("prices cannot be negative")
	}
	
	if config.CacheRead < 0 || config.CacheWrite < 0 {
		return fmt.Errorf("cache prices cannot be negative")
	}
	
	if config.BatchDiscount < 0 || config.BatchDiscount > 1 {
		return fmt.Errorf("batch discount must be between 0 and 1")
	}
	
	// 如果有特定驗證規則
	if rule, exists := pe.validationRules[name]; exists {
		if config.Input < rule.MinPrice || config.Input > rule.MaxPrice {
			return fmt.Errorf("input price out of range [%f, %f]", rule.MinPrice, rule.MaxPrice)
		}
		
		if config.Output < rule.MinPrice || config.Output > rule.MaxPrice {
			return fmt.Errorf("output price out of range [%f, %f]", rule.MinPrice, rule.MaxPrice)
		}
		
		if rule.RequireCache && (config.CacheRead == 0 || config.CacheWrite == 0) {
			return fmt.Errorf("cache pricing required but not provided")
		}
		
		if rule.RequireBatch && config.BatchDiscount == 0 {
			return fmt.Errorf("batch discount required but not provided")
		}
	}
	
	return nil
}

// GetModelInfo 取得模型詳細資訊
func (pe *PricingEngine) GetModelInfo(name string) (map[string]interface{}, error) {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	model, exists := pe.models[name]
	if !exists {
		return nil, fmt.Errorf("model '%s' not found", name)
	}
	
	info := map[string]interface{}{
		"name":            model.Name,
		"input_price":     fmt.Sprintf("$%.2f/MTok", model.InputPrice),
		"output_price":    fmt.Sprintf("$%.2f/MTok", model.OutputPrice),
		"cache_read":      fmt.Sprintf("$%.2f/MTok", model.CacheRead),
		"cache_write":     fmt.Sprintf("$%.2f/MTok", model.CacheWrite),
		"batch_discount":  fmt.Sprintf("%.0f%%", model.BatchDiscount*100),
		"cost_ratio":      model.OutputPrice / model.InputPrice,
		"cache_savings":   (model.InputPrice - model.CacheRead) / model.InputPrice * 100,
	}
	
	return info, nil
}

// GetModelComparison 取得模型比較資訊
func (pe *PricingEngine) GetModelComparison() []map[string]interface{} {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	
	var comparison []map[string]interface{}
	
	for name := range pe.models {
		if info, err := pe.GetModelInfo(name); err == nil {
			comparison = append(comparison, info)
		}
	}
	
	// 按照輸入價格排序
	sort.Slice(comparison, func(i, j int) bool {
		model1, _ := pe.models[comparison[i]["name"].(string)]
		model2, _ := pe.models[comparison[j]["name"].(string)]
		return model1.InputPrice < model2.InputPrice
	})
	
	return comparison
}

// LoadDefaultModels 載入預設定價模型
func (pe *PricingEngine) LoadDefaultModels() {
	// Claude Sonnet 4.0
	pe.AddPricingModel("claude-sonnet-4.0", &types.PricingModel{
		Name:          "claude-sonnet-4.0",
		InputPrice:    3.0,   // $3/MTok
		OutputPrice:   15.0,  // $15/MTok
		CacheRead:     0.30,  // $0.30/MTok
		CacheWrite:    3.75,  // $3.75/MTok
		BatchDiscount: 0.5,   // 50% discount
	})

	// Claude Opus 4.0
	pe.AddPricingModel("claude-opus-4.0", &types.PricingModel{
		Name:          "claude-opus-4.0",
		InputPrice:    15.0,  // $15/MTok
		OutputPrice:   75.0,  // $75/MTok
		CacheRead:     1.5,   // $1.5/MTok
		CacheWrite:    18.75, // $18.75/MTok
		BatchDiscount: 0.5,   // 50% discount
	})

	// Claude Haiku 3.5
	pe.AddPricingModel("claude-haiku-3.5", &types.PricingModel{
		Name:          "claude-haiku-3.5",
		InputPrice:    0.8, // $0.8/MTok
		OutputPrice:   4.0, // $4/MTok
		CacheRead:     0.08, // $0.08/MTok
		CacheWrite:    1.0,  // $1.0/MTok
		BatchDiscount: 0.5,  // 50% discount
	})
}

// CalculateBasicCost 計算基本成本（不含快取和批次折扣）
func (pe *PricingEngine) CalculateBasicCost(inputTokens, outputTokens int, modelName string) (*types.CostBreakdown, error) {
	ctx := context.Background()
	
	// 驗證輸入參數
	if inputTokens < 0 || outputTokens < 0 {
		appErr := errors.New(errors.ErrCodeInvalidTokenCount, "Token 數量不能為負數")
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "validate_token_count",
			Component:  "pricing_engine",
			Parameters: map[string]interface{}{
				"input_tokens":  inputTokens,
				"output_tokens": outputTokens,
				"model_name":    modelName,
			},
		})
		return nil, pe.errorHandler.Handle(ctx, appErr)
	}
	
	model, err := pe.GetPricingModel(modelName)
	if err != nil {
		appErr := errors.Wrap(err, errors.ErrCodeCostCalculation, "取得定價模型失敗")
		appErr = appErr.WithContext(errors.ErrorContext{
			Operation:  "get_pricing_model",
			Component:  "pricing_engine",
			Parameters: map[string]interface{}{
				"model_name": modelName,
			},
		})
		return nil, pe.errorHandler.Handle(ctx, appErr)
	}

	// 將 tokens 轉換為百萬 tokens 為單位
	inputMTokens := float64(inputTokens) / 1_000_000
	outputMTokens := float64(outputTokens) / 1_000_000

	inputCost := inputMTokens * model.InputPrice
	outputCost := outputMTokens * model.OutputPrice
	totalCost := inputCost + outputCost

	return &types.CostBreakdown{
		InputCost:    inputCost,
		OutputCost:   outputCost,
		TotalCost:    totalCost,
		Currency:     "USD",
		PricingModel: modelName,
	}, nil
}

// CalculateCostWithCache 計算包含快取的成本
func (pe *PricingEngine) CalculateCostWithCache(inputTokens, outputTokens, cacheReadTokens, cacheWriteTokens int, modelName string) (*types.CostBreakdown, error) {
	model, err := pe.GetPricingModel(modelName)
	if err != nil {
		return nil, err
	}

	// 將 tokens 轉換為百萬 tokens 為單位
	inputMTokens := float64(inputTokens) / 1_000_000
	outputMTokens := float64(outputTokens) / 1_000_000
	cacheReadMTokens := float64(cacheReadTokens) / 1_000_000
	cacheWriteMTokens := float64(cacheWriteTokens) / 1_000_000

	inputCost := inputMTokens * model.InputPrice
	outputCost := outputMTokens * model.OutputPrice
	cacheReadCost := cacheReadMTokens * model.CacheRead
	cacheWriteCost := cacheWriteMTokens * model.CacheWrite
	
	totalCost := inputCost + outputCost + cacheReadCost + cacheWriteCost

	return &types.CostBreakdown{
		InputCost:    inputCost + cacheReadCost + cacheWriteCost, // 將快取成本歸類為輸入成本
		OutputCost:   outputCost,
		TotalCost:    totalCost,
		Currency:     "USD",
		PricingModel: modelName,
	}, nil
}

// CalculateCostWithBatchDiscount 計算包含批次折扣的成本
func (pe *PricingEngine) CalculateCostWithBatchDiscount(inputTokens, outputTokens int, modelName string, isBatch bool) (*types.CostBreakdown, error) {
	breakdown, err := pe.CalculateBasicCost(inputTokens, outputTokens, modelName)
	if err != nil {
		return nil, err
	}

	if isBatch {
		model, err := pe.GetPricingModel(modelName)
		if err != nil {
			return nil, err
		}

		// 應用批次折扣
		discountMultiplier := 1.0 - model.BatchDiscount
		breakdown.InputCost *= discountMultiplier
		breakdown.OutputCost *= discountMultiplier
		breakdown.TotalCost *= discountMultiplier
	}

	return breakdown, nil
}

// EstimateMonthlyBudget 估算月度預算
func (pe *PricingEngine) EstimateMonthlyBudget(dailyTokens int, modelName string) (float64, error) {
	// 假設 input 和 output tokens 各占一半
	inputTokens := dailyTokens / 2
	outputTokens := dailyTokens / 2

	dailyCost, err := pe.CalculateBasicCost(inputTokens, outputTokens, modelName)
	if err != nil {
		return 0, err
	}

	// 一個月按30天計算
	monthlyBudget := dailyCost.TotalCost * 30
	return monthlyBudget, nil
}

// ComparePricingModels 比較不同定價模型的成本
func (pe *PricingEngine) ComparePricingModels(inputTokens, outputTokens int) (map[string]*types.CostBreakdown, error) {
	comparison := make(map[string]*types.CostBreakdown)

	for modelName := range pe.models {
		cost, err := pe.CalculateBasicCost(inputTokens, outputTokens, modelName)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate cost for model %s: %v", modelName, err)
		}
		comparison[modelName] = cost
	}

	return comparison, nil
}