package cost

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"
	"token-monitor/internal/types"
)

// Optimizer 成本優化分析器
type Optimizer struct {
	pricingEngine    *PricingEngine
	cacheThreshold   int     // 快取閾值（token 數量）
	batchThreshold   int     // 批次處理閾值
	confidenceMin    float64 // 最小信心度
	minSaving        float64 // 最小節省金額（USD）
}

// OptimizationContext 優化分析上下文
type OptimizationContext struct {
	TotalCost       float64
	TotalTokens     int
	SessionCount    int
	ActivityStats   map[types.ActivityType]ActivityStats
	TimeRange       struct {
		Start time.Time
		End   time.Time
	}
	UsagePatterns   []UsagePattern
}

// ActivityStats 活動統計
type ActivityStats struct {
	Count       int
	TokensUsed  int
	Cost        float64
	AvgPerRound float64
}

// UsagePattern 使用模式
type UsagePattern struct {
	Type        string
	Frequency   int
	TokensUsed  int
	Cost        float64
	Potential   float64 // 優化潛力
}

// NewOptimizer 創建新的優化器
func NewOptimizer(pricingEngine *PricingEngine) *Optimizer {
	return &Optimizer{
		pricingEngine:  pricingEngine,
		cacheThreshold: 1000,   // 1K tokens 以上建議快取
		batchThreshold: 5,      // 5個以上請求建議批次處理
		confidenceMin:  0.7,    // 70% 最小信心度
		minSaving:      0.01,   // 最小節省 $0.01
	}
}

// AnalyzeAndSuggest 分析並提供優化建議
func (o *Optimizer) AnalyzeAndSuggest(records []types.UsageRecord) (*types.OptimizationSuggestions, error) {
	if len(records) == 0 {
		return &types.OptimizationSuggestions{
			Suggestions:   []types.OptimizationSuggestion{},
			TotalSavings:  0,
			CurrentCost:   0,
			OptimizedCost: 0,
		}, nil
	}
	
	// 建立分析上下文
	context, err := o.buildContext(records)
	if err != nil {
		return nil, fmt.Errorf("failed to build optimization context: %w", err)
	}
	
	var suggestions []types.OptimizationSuggestion
	totalSavings := 0.0
	
	// 分析快取機會
	cacheSuggestions, cacheSavings := o.analyzeCacheOpportunities(context)
	suggestions = append(suggestions, cacheSuggestions...)
	totalSavings += cacheSavings
	
	// 分析批次處理機會
	batchSuggestions, batchSavings := o.analyzeBatchOpportunities(context)
	suggestions = append(suggestions, batchSuggestions...)
	totalSavings += batchSavings
	
	// 分析模型選擇優化
	modelSuggestions, modelSavings := o.analyzeModelOptimization(context)
	suggestions = append(suggestions, modelSuggestions...)
	totalSavings += modelSavings
	
	// 分析工作流程優化
	workflowSuggestions, workflowSavings := o.analyzeWorkflowOptimization(context)
	suggestions = append(suggestions, workflowSuggestions...)
	totalSavings += workflowSavings
	
	// 過濾低信心度和低節省的建議
	filteredSuggestions := o.filterSuggestions(suggestions)
	
	return &types.OptimizationSuggestions{
		Suggestions:   filteredSuggestions,
		TotalSavings:  totalSavings,
		CurrentCost:   context.TotalCost,
		OptimizedCost: context.TotalCost - totalSavings,
	}, nil
}

// buildContext 建立優化分析上下文
func (o *Optimizer) buildContext(records []types.UsageRecord) (*OptimizationContext, error) {
	context := &OptimizationContext{
		ActivityStats: make(map[types.ActivityType]ActivityStats),
	}
	
	var minTime, maxTime time.Time
	
	for i, record := range records {
		// 計算總成本和 tokens
		context.TotalCost += record.Cost.Total
		context.TotalTokens += record.Tokens.Total
		
		// 時間範圍
		if i == 0 {
			minTime = record.Timestamp
			maxTime = record.Timestamp
		} else {
			if record.Timestamp.Before(minTime) {
				minTime = record.Timestamp
			}
			if record.Timestamp.After(maxTime) {
				maxTime = record.Timestamp
			}
		}
		
		// 活動統計
		activityType := record.Activity.Type
		stats := context.ActivityStats[activityType]
		stats.Count++
		stats.TokensUsed += record.Tokens.Total
		stats.Cost += record.Cost.Total
		if stats.Count > 0 {
			stats.AvgPerRound = stats.Cost / float64(stats.Count)
		}
		context.ActivityStats[activityType] = stats
	}
	
	context.TimeRange.Start = minTime
	context.TimeRange.End = maxTime
	context.SessionCount = len(o.extractUniqueSessions(records))
	
	// 分析使用模式
	context.UsagePatterns = o.extractUsagePatterns(records)
	
	return context, nil
}

// analyzeCacheOpportunities 分析快取機會
func (o *Optimizer) analyzeCacheOpportunities(context *OptimizationContext) ([]types.OptimizationSuggestion, float64) {
	var suggestions []types.OptimizationSuggestion
	totalSavings := 0.0
	
	// 分析重複內容的快取機會
	for activityType, stats := range context.ActivityStats {
		if stats.TokensUsed > o.cacheThreshold && stats.Count > 2 {
			// 估算快取節省
			avgTokensPerRound := float64(stats.TokensUsed) / float64(stats.Count)
			
			// 假設 30% 的內容可以快取
			cacheableTokens := int(avgTokensPerRound * 0.3)
			
			// 計算節省成本
			originalCost, _ := o.pricingEngine.CalculateBasicCost(cacheableTokens/2, cacheableTokens/2, "claude-sonnet-4.0")
			cachedCost, _ := o.pricingEngine.CalculateCostWithCache(cacheableTokens/2, cacheableTokens/2, cacheableTokens/4, cacheableTokens/4, "claude-sonnet-4.0")
			
			if originalCost != nil && cachedCost != nil {
				saving := (originalCost.TotalCost - cachedCost.TotalCost) * float64(stats.Count)
				
				if saving > o.minSaving {
					confidence := o.calculateCacheConfidence(stats, context)
					
					suggestions = append(suggestions, types.OptimizationSuggestion{
						Type:            "cache",
						Description:     fmt.Sprintf("為 %s 活動啟用提示快取，可節省重複計算成本", activityType),
						PotentialSaving: saving,
						Confidence:      confidence,
					})
					
					totalSavings += saving
				}
			}
		}
	}
	
	return suggestions, totalSavings
}

// analyzeBatchOpportunities 分析批次處理機會
func (o *Optimizer) analyzeBatchOpportunities(context *OptimizationContext) ([]types.OptimizationSuggestion, float64) {
	var suggestions []types.OptimizationSuggestion
	totalSavings := 0.0
	
	// 尋找可批次處理的模式
	for _, pattern := range context.UsagePatterns {
		if pattern.Frequency >= o.batchThreshold && pattern.Type == "repetitive" {
			// 計算批次折扣節省
			batchSaving := pattern.Cost * 0.5 // 50% 批次折扣
			
			if batchSaving > o.minSaving {
				confidence := o.calculateBatchConfidence(pattern, context)
				
				suggestions = append(suggestions, types.OptimizationSuggestion{
					Type:            "batch",
					Description:     fmt.Sprintf("將 %d 個相似請求合併為批次處理，享受50%%折扣", pattern.Frequency),
					PotentialSaving: batchSaving,
					Confidence:      confidence,
				})
				
				totalSavings += batchSaving
			}
		}
	}
	
	return suggestions, totalSavings
}

// analyzeModelOptimization 分析模型選擇優化
func (o *Optimizer) analyzeModelOptimization(context *OptimizationContext) ([]types.OptimizationSuggestion, float64) {
	var suggestions []types.OptimizationSuggestion
	totalSavings := 0.0
	
	// 檢查是否使用了成本較高的模型進行簡單任務
	for activityType, stats := range context.ActivityStats {
		if activityType == types.ActivityChat || activityType == types.ActivityDocumentation {
			// 這些活動可能適合使用較便宜的模型
			avgTokensPerRound := float64(stats.TokensUsed) / float64(stats.Count)
			
			// 比較不同模型的成本
			currentCost, _ := o.pricingEngine.CalculateBasicCost(int(avgTokensPerRound/2), int(avgTokensPerRound/2), "claude-sonnet-4.0")
			cheaperCost, _ := o.pricingEngine.CalculateBasicCost(int(avgTokensPerRound/2), int(avgTokensPerRound/2), "claude-haiku-3.5")
			
			if currentCost != nil && cheaperCost != nil {
				saving := (currentCost.TotalCost - cheaperCost.TotalCost) * float64(stats.Count)
				
				if saving > o.minSaving {
					confidence := o.calculateModelSwitchConfidence(activityType, stats)
					
					suggestions = append(suggestions, types.OptimizationSuggestion{
						Type:            "model-switch",
						Description:     fmt.Sprintf("對於 %s 活動使用 Claude Haiku 3.5 替代 Sonnet 4.0", activityType),
						PotentialSaving: saving,
						Confidence:      confidence,
					})
					
					totalSavings += saving
				}
			}
		}
	}
	
	return suggestions, totalSavings
}

// analyzeWorkflowOptimization 分析工作流程優化
func (o *Optimizer) analyzeWorkflowOptimization(context *OptimizationContext) ([]types.OptimizationSuggestion, float64) {
	var suggestions []types.OptimizationSuggestion
	totalSavings := 0.0
	
	// 檢查是否有低效的使用模式
	for activityType, stats := range context.ActivityStats {
		// 如果平均每輪成本很高，可能需要優化工作流程
		if stats.AvgPerRound > 0.1 { // 每輪成本超過 $0.1
			efficiency := float64(stats.TokensUsed) / stats.Cost
			
			if efficiency < 10000 { // 每美元少於 10K tokens
				potentialSaving := stats.Cost * 0.2 // 假設可節省 20%
				
				if potentialSaving > o.minSaving {
					suggestions = append(suggestions, types.OptimizationSuggestion{
						Type:            "workflow",
						Description:     fmt.Sprintf("優化 %s 工作流程，減少不必要的往返對話", activityType),
						PotentialSaving: potentialSaving,
						Confidence:      0.6, // 工作流程優化信心度相對較低
					})
					
					totalSavings += potentialSaving
				}
			}
		}
	}
	
	return suggestions, totalSavings
}

// extractUniqueSessions 提取唯一會話
func (o *Optimizer) extractUniqueSessions(records []types.UsageRecord) []string {
	sessionMap := make(map[string]bool)
	
	for _, record := range records {
		if record.SessionID != "" {
			sessionMap[record.SessionID] = true
		}
	}
	
	sessions := make([]string, 0, len(sessionMap))
	for session := range sessionMap {
		sessions = append(sessions, session)
	}
	
	return sessions
}

// extractUsagePatterns 提取使用模式
func (o *Optimizer) extractUsagePatterns(records []types.UsageRecord) []UsagePattern {
	patterns := make(map[string]*UsagePattern)
	
	for _, record := range records {
		// 簡化的模式識別：基於活動類型和 token 範圍
		tokenRange := o.getTokenRange(record.Tokens.Total)
		patternKey := fmt.Sprintf("%s-%s", record.Activity.Type, tokenRange)
		
		if pattern, exists := patterns[patternKey]; exists {
			pattern.Frequency++
			pattern.TokensUsed += record.Tokens.Total
			pattern.Cost += record.Cost.Total
		} else {
			patterns[patternKey] = &UsagePattern{
				Type:       "repetitive",
				Frequency:  1,
				TokensUsed: record.Tokens.Total,
				Cost:       record.Cost.Total,
			}
		}
	}
	
	var result []UsagePattern
	for _, pattern := range patterns {
		result = append(result, *pattern)
	}
	
	return result
}

// getTokenRange 取得 token 範圍分類
func (o *Optimizer) getTokenRange(tokens int) string {
	switch {
	case tokens < 100:
		return "small"
	case tokens < 1000:
		return "medium"
	case tokens < 5000:
		return "large"
	default:
		return "xlarge"
	}
}

// calculateCacheConfidence 計算快取信心度
func (o *Optimizer) calculateCacheConfidence(stats ActivityStats, context *OptimizationContext) float64 {
	confidence := 0.5 // 基礎信心度
	
	// 重複次數越多，信心度越高
	if stats.Count > 5 {
		confidence += 0.2
	}
	
	// 平均 token 數越多，快取效果越好
	avgTokens := float64(stats.TokensUsed) / float64(stats.Count)
	if avgTokens > 1000 {
		confidence += 0.2
	}
	
	// 使用時間越長，模式越穩定
	if context.TimeRange.End.Sub(context.TimeRange.Start).Hours() > 24 {
		confidence += 0.1
	}
	
	return math.Min(confidence, 1.0)
}

// calculateBatchConfidence 計算批次處理信心度
func (o *Optimizer) calculateBatchConfidence(pattern UsagePattern, context *OptimizationContext) float64 {
	confidence := 0.6 // 批次處理基礎信心度較高
	
	// 頻率越高，信心度越高
	if pattern.Frequency > 10 {
		confidence += 0.2
	}
	
	// token 數量相似度高，適合批次處理
	avgTokens := float64(pattern.TokensUsed) / float64(pattern.Frequency)
	if avgTokens > 500 && avgTokens < 2000 {
		confidence += 0.1
	}
	
	return math.Min(confidence, 1.0)
}

// calculateModelSwitchConfidence 計算模型切換信心度
func (o *Optimizer) calculateModelSwitchConfidence(activityType types.ActivityType, stats ActivityStats) float64 {
	confidence := 0.4 // 模型切換基礎信心度較低
	
	// 特定活動類型適合切換
	switch activityType {
	case types.ActivityChat:
		confidence = 0.8
	case types.ActivityDocumentation:
		confidence = 0.7
	case types.ActivityDebugging:
		confidence = 0.3 // 除錯需要更強的模型
	case types.ActivityCoding:
		confidence = 0.2 // 編碼需要更強的模型
	}
	
	// 平均 token 數較少，適合較便宜的模型
	avgTokens := float64(stats.TokensUsed) / float64(stats.Count)
	if avgTokens < 500 {
		confidence += 0.1
	}
	
	return math.Min(confidence, 1.0)
}

// filterSuggestions 過濾建議
func (o *Optimizer) filterSuggestions(suggestions []types.OptimizationSuggestion) []types.OptimizationSuggestion {
	var filtered []types.OptimizationSuggestion
	
	for _, suggestion := range suggestions {
		if suggestion.PotentialSaving >= o.minSaving && suggestion.Confidence >= o.confidenceMin {
			filtered = append(filtered, suggestion)
		}
	}
	
	// 按潛在節省排序
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PotentialSaving > filtered[j].PotentialSaving
	})
	
	return filtered
}

// SetThresholds 設定優化閾值
func (o *Optimizer) SetThresholds(cacheThreshold, batchThreshold int, confidenceMin, minSaving float64) {
	if cacheThreshold > 0 {
		o.cacheThreshold = cacheThreshold
	}
	if batchThreshold > 0 {
		o.batchThreshold = batchThreshold
	}
	if confidenceMin > 0 && confidenceMin <= 1 {
		o.confidenceMin = confidenceMin
	}
	if minSaving > 0 {
		o.minSaving = minSaving
	}
	
	log.Printf("Updated optimization thresholds: cache=%d, batch=%d, confidence=%.2f, minSaving=%.4f", 
		o.cacheThreshold, o.batchThreshold, o.confidenceMin, o.minSaving)
}

// GetThresholds 取得當前閾值設定
func (o *Optimizer) GetThresholds() map[string]interface{} {
	return map[string]interface{}{
		"cache_threshold":   o.cacheThreshold,
		"batch_threshold":   o.batchThreshold,
		"confidence_min":    o.confidenceMin,
		"min_saving":        o.minSaving,
	}
}