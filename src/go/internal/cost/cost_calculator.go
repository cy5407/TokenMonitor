package cost

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	"token-monitor/internal/types"

	"gopkg.in/yaml.v3"
)

// CostCalculatorImpl CostCalculator 介面實作
type CostCalculatorImpl struct {
	pricingEngine *PricingEngine
	configPath    string
	mutex         sync.RWMutex

	// 成本追蹤
	sessionCosts map[string]float64
	dailyCosts   map[string]float64

	// 最後更新時間
	lastConfigUpdate time.Time
}

// BillingMode 計費模式
type BillingMode int

const (
	StandardBilling BillingMode = iota
	CacheBilling
	BatchBilling
)

// CostOptions 成本計算選項
type CostOptions struct {
	Mode             BillingMode
	CacheReadTokens  int
	CacheWriteTokens int
	IsBatch          bool
	SessionID        string
	ActivityType     types.ActivityType
}

// ConfigData 配置文件結構
type ConfigData struct {
	Pricing map[string]struct {
		Input         float64 `yaml:"input"`
		Output        float64 `yaml:"output"`
		CacheRead     float64 `yaml:"cache_read"`
		CacheWrite    float64 `yaml:"cache_write"`
		BatchDiscount float64 `yaml:"batch_discount"`
	} `yaml:"pricing"`
}

// NewCostCalculator 創建新的成本計算器
func NewCostCalculator() *CostCalculatorImpl {
	return &CostCalculatorImpl{
		pricingEngine: NewPricingEngine(),
		sessionCosts:  make(map[string]float64),
		dailyCosts:    make(map[string]float64),
	}
}

// CalculateCost 計算成本（實作 CostCalculator 介面）
func (cc *CostCalculatorImpl) CalculateCost(inputTokens, outputTokens int, model string) (*types.CostBreakdown, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	// 輸入驗證
	if inputTokens < 0 || outputTokens < 0 {
		return nil, fmt.Errorf("token counts cannot be negative: input=%d, output=%d", inputTokens, outputTokens)
	}

	if model == "" {
		model = "claude-sonnet-4.0" // 預設模型
	}

	// 獲取定價模型
	pricingModel, err := cc.pricingEngine.GetPricingModel(model)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing model %s: %w", model, err)
	}

	// 基本成本計算
	breakdown, err := cc.pricingEngine.CalculateBasicCost(inputTokens, outputTokens, model)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate cost for model %s: %w", model, err)
	}

	// 設定詳細資訊以確保測試可以訪問
	breakdown.TokenCounts = types.TokenCounts{
		Input:  inputTokens,
		Output: outputTokens,
		Total:  inputTokens + outputTokens,
	}
	breakdown.CostDetails = types.CostDetails{
		InputRate:   pricingModel.InputPrice,
		OutputRate:  pricingModel.OutputPrice,
		BillingMode: "standard",
	}
	breakdown.Timestamp = time.Now()

	return breakdown, nil
}

// CalculateDetailedCost 計算詳細成本（新增功能）
func (cc *CostCalculatorImpl) CalculateDetailedCost(inputTokens, outputTokens int, model string, options *CostOptions) (*types.CostBreakdown, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	// 輸入驗證
	if err := cc.validateInput(inputTokens, outputTokens, model, options); err != nil {
		return nil, err
	}

	// 獲取定價模型
	pricingModel, err := cc.pricingEngine.GetPricingModel(model)
	if err != nil {
		return nil, fmt.Errorf("failed to get pricing model %s: %w", model, err)
	}

	// 建立詳細的成本分解
	breakdown := &types.CostBreakdown{
		Currency:     "USD",
		PricingModel: model,
		Timestamp:    time.Now(),
		TokenCounts: types.TokenCounts{
			Input:  inputTokens,
			Output: outputTokens,
			Total:  inputTokens + outputTokens,
		},
		CostDetails: types.CostDetails{
			InputRate:   pricingModel.InputPrice,
			OutputRate:  pricingModel.OutputPrice,
			BillingMode: "standard",
		},
	}

	// 設定會話和活動資訊
	if options != nil {
		breakdown.SessionID = options.SessionID
		breakdown.ActivityType = options.ActivityType

		// 設定快取 Token 數量
		if options.CacheReadTokens > 0 || options.CacheWriteTokens > 0 {
			breakdown.TokenCounts.CacheRead = options.CacheReadTokens
			breakdown.TokenCounts.CacheWrite = options.CacheWriteTokens
			breakdown.TokenCounts.Total += options.CacheReadTokens + options.CacheWriteTokens
		}
	}

	// 根據計費模式計算成本
	switch options.Mode {
	case CacheBilling:
		err = cc.calculateCacheCost(breakdown, pricingModel, options)
	case BatchBilling:
		err = cc.calculateBatchCost(breakdown, pricingModel, options)
	default:
		err = cc.calculateStandardCost(breakdown, pricingModel)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to calculate cost: %w", err)
	}

	// 記錄成本到會話和日常追蹤
	if options != nil && options.SessionID != "" {
		cc.sessionCosts[options.SessionID] += breakdown.TotalCost
	}

	today := time.Now().Format("2006-01-02")
	cc.dailyCosts[today] += breakdown.TotalCost

	return breakdown, nil
}

// calculateStandardCost 計算標準成本
func (cc *CostCalculatorImpl) calculateStandardCost(breakdown *types.CostBreakdown, model *types.PricingModel) error {
	// 將 tokens 轉換為百萬 tokens 為單位
	inputMTokens := float64(breakdown.TokenCounts.Input) / 1_000_000
	outputMTokens := float64(breakdown.TokenCounts.Output) / 1_000_000

	breakdown.InputCost = inputMTokens * model.InputPrice
	breakdown.OutputCost = outputMTokens * model.OutputPrice
	breakdown.TotalCost = breakdown.InputCost + breakdown.OutputCost
	breakdown.CostDetails.BillingMode = "standard"

	return nil
}

// calculateCacheCost 計算包含快取的成本
func (cc *CostCalculatorImpl) calculateCacheCost(breakdown *types.CostBreakdown, model *types.PricingModel, options *CostOptions) error {
	// 先計算標準成本
	if err := cc.calculateStandardCost(breakdown, model); err != nil {
		return err
	}

	// 計算快取成本
	cacheReadMTokens := float64(breakdown.TokenCounts.CacheRead) / 1_000_000
	cacheWriteMTokens := float64(breakdown.TokenCounts.CacheWrite) / 1_000_000

	breakdown.CacheReadCost = cacheReadMTokens * model.CacheRead
	breakdown.CacheWriteCost = cacheWriteMTokens * model.CacheWrite
	breakdown.TotalCost += breakdown.CacheReadCost + breakdown.CacheWriteCost

	// 更新成本詳細資訊
	breakdown.CostDetails.CacheReadRate = model.CacheRead
	breakdown.CostDetails.CacheWriteRate = model.CacheWrite
	breakdown.CostDetails.BillingMode = "cache"

	return nil
}

// calculateBatchCost 計算批次成本（含折扣）
func (cc *CostCalculatorImpl) calculateBatchCost(breakdown *types.CostBreakdown, model *types.PricingModel, options *CostOptions) error {
	// 先計算標準成本
	if err := cc.calculateStandardCost(breakdown, model); err != nil {
		return err
	}

	// 應用批次折扣
	if options.IsBatch && model.BatchDiscount > 0 {
		discountMultiplier := 1.0 - model.BatchDiscount
		breakdown.BatchDiscount = breakdown.TotalCost * model.BatchDiscount
		breakdown.InputCost *= discountMultiplier
		breakdown.OutputCost *= discountMultiplier
		breakdown.TotalCost *= discountMultiplier

		breakdown.CostDetails.DiscountRate = model.BatchDiscount
		breakdown.CostDetails.BillingMode = "batch"
	}

	return nil
}

// CalculateCostWithOptions 計算成本（含選項）
func (cc *CostCalculatorImpl) CalculateCostWithOptions(inputTokens, outputTokens int, model string, options *CostOptions) (*types.CostBreakdown, error) {
	// 使用新的詳細成本計算方法
	return cc.CalculateDetailedCost(inputTokens, outputTokens, model, options)
}

// GetPricingInfo 取得定價資訊（實作 CostCalculator 介面）
func (cc *CostCalculatorImpl) GetPricingInfo(model string) (*types.PricingModel, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	if model == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	pricingModel, err := cc.pricingEngine.GetPricingModel(model)
	if err != nil {
		// 避免在日誌中洩露敏感定價資訊
		log.Printf("Pricing model not found: %s", model)
		return nil, fmt.Errorf("pricing model '%s' not available", model)
	}

	return pricingModel, nil
}

// CalculateOptimizationSavings 計算優化節省（實作 CostCalculator 介面）
func (cc *CostCalculatorImpl) CalculateOptimizationSavings(records []types.UsageRecord) (*types.OptimizationSuggestions, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	if len(records) == 0 {
		return &types.OptimizationSuggestions{
			Suggestions:   []types.OptimizationSuggestion{},
			TotalSavings:  0,
			CurrentCost:   0,
			OptimizedCost: 0,
		}, nil
	}

	optimizer := NewOptimizer(cc.pricingEngine)
	return optimizer.AnalyzeAndSuggest(records)
}

// AnalyzeCostTrends 分析成本趨勢（新增功能）
func (cc *CostCalculatorImpl) AnalyzeCostTrends(records []types.UsageRecord, timeRange string) (*types.CostTrendAnalysis, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	if len(records) == 0 {
		return nil, fmt.Errorf("no usage records provided")
	}

	// 按時間分組記錄
	groupedRecords := cc.groupRecordsByTime(records, timeRange)

	// 計算趨勢
	trends := &types.CostTrendAnalysis{
		TimeRange:   timeRange,
		DataPoints:  make([]types.CostDataPoint, 0),
		TotalCost:   0,
		AverageCost: 0,
		GrowthRate:  0,
		Predictions: make([]types.CostPrediction, 0),
	}

	// 計算每個時間點的成本
	for timeKey, timeRecords := range groupedRecords {
		totalCost := 0.0
		totalTokens := 0

		for _, record := range timeRecords {
			// 計算該記錄的成本
			breakdown, err := cc.CalculateCost(record.Tokens.Input, record.Tokens.Output, record.Cost.PricingModel)
			if err != nil {
				continue
			}
			totalCost += breakdown.TotalCost
			totalTokens += record.Tokens.Total
		}

		dataPoint := types.CostDataPoint{
			Timestamp:   timeKey,
			Cost:        totalCost,
			TokenCount:  totalTokens,
			RecordCount: len(timeRecords),
		}

		trends.DataPoints = append(trends.DataPoints, dataPoint)
		trends.TotalCost += totalCost
	}

	// 計算平均成本
	if len(trends.DataPoints) > 0 {
		trends.AverageCost = trends.TotalCost / float64(len(trends.DataPoints))
	}

	// 計算成長率
	if len(trends.DataPoints) >= 2 {
		firstCost := trends.DataPoints[0].Cost
		lastCost := trends.DataPoints[len(trends.DataPoints)-1].Cost
		if firstCost > 0 {
			trends.GrowthRate = ((lastCost - firstCost) / firstCost) * 100
		}
	}

	// 生成預測
	trends.Predictions = cc.generateCostPredictions(trends.DataPoints)

	return trends, nil
}

// GenerateCostReport 生成成本報告（新增功能）
func (cc *CostCalculatorImpl) GenerateCostReport(records []types.UsageRecord, options *types.ReportOptions) (*types.CostReport, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	if len(records) == 0 {
		return nil, fmt.Errorf("no usage records provided")
	}

	report := &types.CostReport{
		GeneratedAt:  time.Now(),
		TimeRange:    options.TimeRange,
		TotalRecords: len(records),
		Summary:      types.CostSummary{},
		ByActivity:   make(map[types.ActivityType]types.CostSummary),
		ByModel:      make(map[string]types.CostSummary),
		Optimization: &types.OptimizationSuggestions{},
		Trends:       &types.CostTrendAnalysis{},
	}

	// 計算總體摘要
	totalCost := 0.0
	totalTokens := 0
	modelUsage := make(map[string]int)
	activityUsage := make(map[types.ActivityType]int)

	for _, record := range records {
		// 計算成本
		breakdown, err := cc.CalculateCost(record.Tokens.Input, record.Tokens.Output, record.Cost.PricingModel)
		if err != nil {
			continue
		}

		totalCost += breakdown.TotalCost
		totalTokens += record.Tokens.Total
		modelUsage[record.Cost.PricingModel]++
		activityUsage[record.Activity.Type]++

		// 按活動類型分組
		if _, exists := report.ByActivity[record.Activity.Type]; !exists {
			report.ByActivity[record.Activity.Type] = types.CostSummary{}
		}
		activitySummary := report.ByActivity[record.Activity.Type]
		activitySummary.TotalCost += breakdown.TotalCost
		activitySummary.TotalTokens += record.Tokens.Total
		activitySummary.RecordCount++
		report.ByActivity[record.Activity.Type] = activitySummary

		// 按模型分組
		if _, exists := report.ByModel[record.Cost.PricingModel]; !exists {
			report.ByModel[record.Cost.PricingModel] = types.CostSummary{}
		}
		modelSummary := report.ByModel[record.Cost.PricingModel]
		modelSummary.TotalCost += breakdown.TotalCost
		modelSummary.TotalTokens += record.Tokens.Total
		modelSummary.RecordCount++
		report.ByModel[record.Cost.PricingModel] = modelSummary
	}

	// 設定總體摘要
	report.Summary.TotalCost = totalCost
	report.Summary.TotalTokens = totalTokens
	report.Summary.RecordCount = len(records)
	if totalTokens > 0 {
		report.Summary.AverageCostPerToken = totalCost / float64(totalTokens) * 1_000_000 // 每百萬 token 的成本
	}

	// 計算平均值
	for activityType, summary := range report.ByActivity {
		if summary.RecordCount > 0 {
			summary.AverageCostPerRecord = summary.TotalCost / float64(summary.RecordCount)
			if summary.TotalTokens > 0 {
				summary.AverageCostPerToken = summary.TotalCost / float64(summary.TotalTokens) * 1_000_000
			}
			report.ByActivity[activityType] = summary
		}
	}

	for model, summary := range report.ByModel {
		if summary.RecordCount > 0 {
			summary.AverageCostPerRecord = summary.TotalCost / float64(summary.RecordCount)
			if summary.TotalTokens > 0 {
				summary.AverageCostPerToken = summary.TotalCost / float64(summary.TotalTokens) * 1_000_000
			}
			report.ByModel[model] = summary
		}
	}

	// 生成優化建議
	optimization, err := cc.CalculateOptimizationSavings(records)
	if err == nil {
		report.Optimization = optimization
	}

	// 生成趨勢分析
	trends, err := cc.AnalyzeCostTrends(records, "daily")
	if err == nil {
		report.Trends = trends
	}

	return report, nil
}

// LoadPricingModels 載入定價模型（實作 CostCalculator 介面）
func (cc *CostCalculatorImpl) LoadPricingModels(configPath string) error {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	cc.configPath = configPath

	// 檢查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, using default pricing models", configPath)
		cc.pricingEngine.LoadDefaultModels()
		return nil
	}

	// 讀取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigData
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 載入定價模型
	cc.pricingEngine.models = make(map[string]*types.PricingModel)
	for modelName, pricing := range config.Pricing {
		cc.pricingEngine.AddPricingModel(modelName, &types.PricingModel{
			Name:          modelName,
			InputPrice:    pricing.Input,
			OutputPrice:   pricing.Output,
			CacheRead:     pricing.CacheRead,
			CacheWrite:    pricing.CacheWrite,
			BatchDiscount: pricing.BatchDiscount,
		})
	}

	cc.lastConfigUpdate = time.Now()
	log.Printf("Loaded %d pricing models from config", len(config.Pricing))

	return nil
}

// GetSupportedModels 取得支援的模型（實作 CostCalculator 介面）
func (cc *CostCalculatorImpl) GetSupportedModels() []string {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.pricingEngine.GetSupportedModels()
}

// GetSessionCost 取得會話成本
func (cc *CostCalculatorImpl) GetSessionCost(sessionID string) float64 {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.sessionCosts[sessionID]
}

// GetDailyCost 取得每日成本
func (cc *CostCalculatorImpl) GetDailyCost(date string) float64 {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.dailyCosts[date]
}

// GetDailyCostSummary 取得每日成本摘要
func (cc *CostCalculatorImpl) GetDailyCostSummary(days int) map[string]float64 {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	summary := make(map[string]float64)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		if cost, exists := cc.dailyCosts[date]; exists {
			summary[date] = cost
		} else {
			summary[date] = 0
		}
	}

	return summary
}

// EstimateMonthlyBudget 估算月度預算
func (cc *CostCalculatorImpl) EstimateMonthlyBudget(dailyTokens int, model string) (float64, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.pricingEngine.EstimateMonthlyBudget(dailyTokens, model)
}

// ComparePricingModels 比較不同定價模型
func (cc *CostCalculatorImpl) ComparePricingModels(inputTokens, outputTokens int) (map[string]*types.CostBreakdown, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.pricingEngine.ComparePricingModels(inputTokens, outputTokens)
}

// ReloadConfig 重新載入配置（用於熱更新）
func (cc *CostCalculatorImpl) ReloadConfig() error {
	if cc.configPath == "" {
		return fmt.Errorf("no config path set")
	}

	return cc.LoadPricingModels(cc.configPath)
}

// GetLastConfigUpdate 取得最後配置更新時間
func (cc *CostCalculatorImpl) GetLastConfigUpdate() time.Time {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	return cc.lastConfigUpdate
}

// ClearSessionCosts 清除會話成本記錄
func (cc *CostCalculatorImpl) ClearSessionCosts() {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	cc.sessionCosts = make(map[string]float64)
}

// ClearDailyCosts 清除每日成本記錄
func (cc *CostCalculatorImpl) ClearDailyCosts() {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	cc.dailyCosts = make(map[string]float64)
}

// validateInput 驗證輸入參數
func (cc *CostCalculatorImpl) validateInput(inputTokens, outputTokens int, model string, options *CostOptions) error {
	if inputTokens < 0 || outputTokens < 0 {
		return fmt.Errorf("token counts cannot be negative: input=%d, output=%d", inputTokens, outputTokens)
	}

	if inputTokens > 10_000_000 || outputTokens > 10_000_000 {
		return fmt.Errorf("token counts exceed maximum limit (10M): input=%d, output=%d", inputTokens, outputTokens)
	}

	if model == "" {
		return fmt.Errorf("model name cannot be empty")
	}

	if options != nil {
		if options.CacheReadTokens < 0 || options.CacheWriteTokens < 0 {
			return fmt.Errorf("cache token counts cannot be negative: read=%d, write=%d",
				options.CacheReadTokens, options.CacheWriteTokens)
		}
	}

	return nil
}

// groupRecordsByTime 按時間分組記錄
func (cc *CostCalculatorImpl) groupRecordsByTime(records []types.UsageRecord, timeRange string) map[time.Time][]types.UsageRecord {
	grouped := make(map[time.Time][]types.UsageRecord)

	for _, record := range records {
		var timeKey time.Time

		switch timeRange {
		case "hourly":
			timeKey = record.Timestamp.Truncate(time.Hour)
		case "daily":
			timeKey = record.Timestamp.Truncate(24 * time.Hour)
		case "weekly":
			// 取得週的開始時間（週一）
			weekday := int(record.Timestamp.Weekday())
			if weekday == 0 {
				weekday = 7 // 將週日從0改為7
			}
			timeKey = record.Timestamp.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
		case "monthly":
			timeKey = time.Date(record.Timestamp.Year(), record.Timestamp.Month(), 1, 0, 0, 0, 0, record.Timestamp.Location())
		default:
			timeKey = record.Timestamp.Truncate(24 * time.Hour)
		}

		grouped[timeKey] = append(grouped[timeKey], record)
	}

	return grouped
}

// generateCostPredictions 生成成本預測
func (cc *CostCalculatorImpl) generateCostPredictions(dataPoints []types.CostDataPoint) []types.CostPrediction {
	if len(dataPoints) < 2 {
		return []types.CostPrediction{}
	}

	predictions := make([]types.CostPrediction, 0)

	// 簡單的線性預測（可以後續改進為更複雜的算法）
	if len(dataPoints) >= 3 {
		// 計算平均成長率
		totalGrowth := 0.0
		validPeriods := 0

		for i := 1; i < len(dataPoints); i++ {
			if dataPoints[i-1].Cost > 0 {
				growth := (dataPoints[i].Cost - dataPoints[i-1].Cost) / dataPoints[i-1].Cost
				totalGrowth += growth
				validPeriods++
			}
		}

		if validPeriods > 0 {
			avgGrowthRate := totalGrowth / float64(validPeriods)
			lastDataPoint := dataPoints[len(dataPoints)-1]

			// 預測未來3個時間點
			for i := 1; i <= 3; i++ {
				predictedCost := lastDataPoint.Cost * (1 + avgGrowthRate*float64(i))
				confidence := 1.0 - (float64(i) * 0.2) // 時間越遠信心度越低

				prediction := types.CostPrediction{
					Date:          lastDataPoint.Timestamp.AddDate(0, 0, i),
					PredictedCost: predictedCost,
					Confidence:    confidence,
				}
				predictions = append(predictions, prediction)
			}
		}
	}

	return predictions
}

// CalculateCostEfficiency 計算成本效率（新增功能）
func (cc *CostCalculatorImpl) CalculateCostEfficiency(records []types.UsageRecord) (*types.CostEfficiencyAnalysis, error) {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	if len(records) == 0 {
		return nil, fmt.Errorf("no usage records provided")
	}

	analysis := &types.CostEfficiencyAnalysis{
		OverallEfficiency: 0,
		ByActivity:        make(map[types.ActivityType]float64),
		ByModel:           make(map[string]float64),
		Recommendations:   make([]string, 0),
	}

	// 按活動類型分析效率
	activityStats := make(map[types.ActivityType]struct {
		totalCost   float64
		totalTokens int
		totalRounds int
	})

	for _, record := range records {
		breakdown, err := cc.CalculateCost(record.Tokens.Input, record.Tokens.Output, record.Cost.PricingModel)
		if err != nil {
			continue
		}

		stats := activityStats[record.Activity.Type]
		stats.totalCost += breakdown.TotalCost
		stats.totalTokens += record.Tokens.Total
		stats.totalRounds += record.Activity.Rounds
		activityStats[record.Activity.Type] = stats
	}

	// 計算各活動類型的效率（tokens per dollar）
	totalEfficiency := 0.0
	validActivities := 0

	for activityType, stats := range activityStats {
		if stats.totalCost > 0 {
			efficiency := float64(stats.totalTokens) / stats.totalCost
			analysis.ByActivity[activityType] = efficiency
			totalEfficiency += efficiency
			validActivities++

			// 生成建議
			if efficiency < 10000 { // 低於 10k tokens per dollar
				analysis.Recommendations = append(analysis.Recommendations,
					fmt.Sprintf("%s 活動的成本效率較低，建議優化提示詞或使用快取", activityType))
			}
		}
	}

	if validActivities > 0 {
		analysis.OverallEfficiency = totalEfficiency / float64(validActivities)
	}

	return analysis, nil
}

// GetStatistics 取得統計資訊
func (cc *CostCalculatorImpl) GetStatistics() map[string]interface{} {
	cc.mutex.RLock()
	defer cc.mutex.RUnlock()

	stats := make(map[string]interface{})

	// 會話統計
	totalSessions := len(cc.sessionCosts)
	totalSessionCost := 0.0
	for _, cost := range cc.sessionCosts {
		totalSessionCost += cost
	}

	// 每日統計
	totalDays := len(cc.dailyCosts)
	totalDailyCost := 0.0
	for _, cost := range cc.dailyCosts {
		totalDailyCost += cost
	}

	stats["total_sessions"] = totalSessions
	stats["total_session_cost"] = totalSessionCost
	stats["total_days"] = totalDays
	stats["total_daily_cost"] = totalDailyCost
	stats["supported_models"] = len(cc.pricingEngine.models)
	stats["last_config_update"] = cc.lastConfigUpdate.Format(time.RFC3339)

	if totalSessions > 0 {
		stats["average_session_cost"] = totalSessionCost / float64(totalSessions)
	}

	if totalDays > 0 {
		stats["average_daily_cost"] = totalDailyCost / float64(totalDays)
	}

	return stats
}
