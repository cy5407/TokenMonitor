package interfaces

import (
	"time"
	"token-monitor/internal/types"
)

// TokenCalculator Token 計算介面
type TokenCalculator interface {
	// CalculateTokens 計算文本的 Token 數量
	CalculateTokens(text string, method string) (int, error)

	// AnalyzeTokenDistribution 分析 Token 分佈
	AnalyzeTokenDistribution(text string) (*types.TokenDistribution, error)

	// IsTiktokenAvailable 檢查 tiktoken 是否可用
	IsTiktokenAvailable() bool

	// ClearCache 清除快取
	ClearCache()

	// GetSupportedMethods 取得支援的計算方法
	GetSupportedMethods() []string
}

// ActivityAnalyzer 活動分析介面
type ActivityAnalyzer interface {
	// ClassifyActivity 分類活動類型
	ClassifyActivity(content string) types.ActivityType

	// AnalyzeActivityBatch 批次分析多個活動
	AnalyzeActivityBatch(contents []string) []types.ActivityType

	// GenerateActivitySummary 生成活動摘要
	GenerateActivitySummary(activities []types.Activity) types.ActivitySummary

	// CalculateEfficiencyMetrics 計算效率指標
	CalculateEfficiencyMetrics(data types.ActivityData) types.EfficiencyMetrics

	// GetActivityTypeDistribution 獲取活動類型分佈
	GetActivityTypeDistribution(activities []types.Activity) map[types.ActivityType]float64
}

// CostCalculator 成本計算介面
type CostCalculator interface {
	// CalculateCost 計算成本
	CalculateCost(inputTokens, outputTokens int, model string) (*types.CostBreakdown, error)

	// GetPricingInfo 取得定價資訊
	GetPricingInfo(model string) (*types.PricingModel, error)

	// CalculateOptimizationSavings 計算優化節省
	CalculateOptimizationSavings(records []types.UsageRecord) (*types.OptimizationSuggestions, error)

	// LoadPricingModels 載入定價模型
	LoadPricingModels(configPath string) error

	// GetSupportedModels 取得支援的模型
	GetSupportedModels() []string
}

// ReportGenerator 報告生成介面
type ReportGenerator interface {
	// GenerateJSONReport 生成 JSON 報告
	GenerateJSONReport(data interface{}) ([]byte, error)

	// GenerateCSVReport 生成 CSV 報告
	GenerateCSVReport(data interface{}) ([]byte, error)

	// GenerateHTMLReport 生成 HTML 報告
	GenerateHTMLReport(data interface{}) ([]byte, error)

	// GenerateSummaryStats 生成摘要統計
	GenerateSummaryStats(records []types.UsageRecord) (*types.SummaryStats, error)

	// ExportReport 匯出報告到檔案
	ExportReport(data interface{}, config *types.ReportConfig) error
}

// DataStore 資料儲存介面
type DataStore interface {
	// SaveUsageRecord 儲存使用記錄
	SaveUsageRecord(record *types.UsageRecord) error

	// GetUsageRecords 取得使用記錄
	GetUsageRecords(start, end *time.Time) ([]types.UsageRecord, error)

	// GetUsageRecordsByActivity 依活動類型取得記錄
	GetUsageRecordsByActivity(activityType types.ActivityType) ([]types.UsageRecord, error)

	// DeleteOldRecords 刪除舊記錄
	DeleteOldRecords(before time.Time) error

	// BackupData 備份資料
	BackupData(path string) error

	// RestoreData 恢復資料
	RestoreData(path string) error
}

// Monitor 監控介面
type Monitor interface {
	// Start 開始監控
	Start() error

	// Stop 停止監控
	Stop() error

	// IsRunning 檢查是否正在執行
	IsRunning() bool

	// GetStatus 取得狀態
	GetStatus() map[string]interface{}
}
