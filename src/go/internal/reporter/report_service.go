package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"token-monitor/internal/types"
)

// ReportService 報告服務
type ReportService struct {
	generator       *ReportGenerator
	templateManager *TemplateManager
	config          types.ReportConfig
}

// NewReportService 建立新的報告服務
func NewReportService(config types.ReportConfig) *ReportService {
	return &ReportService{
		generator:       NewReportGenerator(),
		templateManager: NewTemplateManager(),
		config:          config,
	}
}

// GenerateReport 生成報告
func (rs *ReportService) GenerateReport(activities []types.Activity, options types.ReportOptions) (*types.ReportResult, error) {
	// 過濾時間範圍內的活動
	filteredActivities := rs.filterActivitiesByTimeRange(activities, options.TimeRange)

	// 生成基礎報告
	basicReport, err := rs.generator.GenerateBasicReport(filteredActivities)
	if err != nil {
		return nil, fmt.Errorf("生成基礎報告失敗: %w", err)
	}

	// 生成 JSON 報告
	jsonData, err := rs.generator.GenerateJSONReport(filteredActivities)
	if err != nil {
		return nil, fmt.Errorf("生成 JSON 報告失敗: %w", err)
	}

	// 生成文字報告
	textReport, err := rs.generateTextReport(basicReport, options)
	if err != nil {
		return nil, fmt.Errorf("生成文字報告失敗: %w", err)
	}

	result := &types.ReportResult{
		BasicReport: basicReport,
		JSONData:    jsonData,
		TextReport:  textReport,
		GeneratedAt: time.Now(),
		Options:     options,
	}

	return result, nil
}

// SaveReport 儲存報告到檔案
func (rs *ReportService) SaveReport(result *types.ReportResult, outputPath string) error {
	// 確保輸出目錄存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("建立輸出目錄失敗: %w", err)
	}

	// 根據檔案副檔名決定儲存格式
	ext := filepath.Ext(outputPath)

	switch ext {
	case ".json":
		return rs.saveJSONReport(result, outputPath)
	case ".txt", ".md":
		return rs.saveTextReport(result, outputPath)
	default:
		// 預設儲存為 JSON
		return rs.saveJSONReport(result, outputPath+".json")
	}
}

// GenerateAndSaveReport 生成並儲存報告
func (rs *ReportService) GenerateAndSaveReport(activities []types.Activity, options types.ReportOptions, outputPath string) error {
	result, err := rs.GenerateReport(activities, options)
	if err != nil {
		return fmt.Errorf("生成報告失敗: %w", err)
	}

	return rs.SaveReport(result, outputPath)
}

// GetAvailableTemplates 獲取可用的模板列表
func (rs *ReportService) GetAvailableTemplates() []string {
	return rs.templateManager.ListTemplates()
}

// RegisterCustomTemplate 註冊自定義模板
func (rs *ReportService) RegisterCustomTemplate(name, template string) {
	rs.templateManager.RegisterTemplate(name, template)
}

// GenerateReportWithTemplate 使用指定模板生成報告
func (rs *ReportService) GenerateReportWithTemplate(activities []types.Activity, templateName string, options types.ReportOptions) (string, error) {
	// 過濾活動
	filteredActivities := rs.filterActivitiesByTimeRange(activities, options.TimeRange)

	// 生成基礎報告
	basicReport, err := rs.generator.GenerateBasicReport(filteredActivities)
	if err != nil {
		return "", fmt.Errorf("生成基礎報告失敗: %w", err)
	}

	// 使用模板渲染
	return rs.templateManager.RenderReport(templateName, basicReport)
}

// filterActivitiesByTimeRange 按時間範圍過濾活動
func (rs *ReportService) filterActivitiesByTimeRange(activities []types.Activity, timeRange types.TimeRange) []types.Activity {
	if timeRange.Start.IsZero() && timeRange.End.IsZero() {
		return activities
	}

	var filtered []types.Activity
	for _, activity := range activities {
		// 如果沒有設定開始時間，只檢查結束時間
		if timeRange.Start.IsZero() {
			if activity.Timestamp.Before(timeRange.End) || activity.Timestamp.Equal(timeRange.End) {
				filtered = append(filtered, activity)
			}
			continue
		}

		// 如果沒有設定結束時間，只檢查開始時間
		if timeRange.End.IsZero() {
			if activity.Timestamp.After(timeRange.Start) || activity.Timestamp.Equal(timeRange.Start) {
				filtered = append(filtered, activity)
			}
			continue
		}

		// 檢查活動是否在時間範圍內
		if (activity.Timestamp.After(timeRange.Start) || activity.Timestamp.Equal(timeRange.Start)) &&
			(activity.Timestamp.Before(timeRange.End) || activity.Timestamp.Equal(timeRange.End)) {
			filtered = append(filtered, activity)
		}
	}

	return filtered
}

// generateTextReport 生成文字報告
func (rs *ReportService) generateTextReport(basicReport *types.BasicReport, options types.ReportOptions) (string, error) {
	// 根據選項決定使用哪個模板
	templateName := "basic"
	if options.IncludeTrends {
		templateName = "detailed"
	}

	return rs.templateManager.RenderReport(templateName, basicReport)
}

// saveJSONReport 儲存 JSON 報告
func (rs *ReportService) saveJSONReport(result *types.ReportResult, outputPath string) error {
	return os.WriteFile(outputPath, result.JSONData, 0644)
}

// saveTextReport 儲存文字報告
func (rs *ReportService) saveTextReport(result *types.ReportResult, outputPath string) error {
	return os.WriteFile(outputPath, []byte(result.TextReport), 0644)
}

// GetReportStatistics 獲取報告統計資訊
func (rs *ReportService) GetReportStatistics(activities []types.Activity) *types.ReportStatistics {
	if len(activities) == 0 {
		return &types.ReportStatistics{}
	}

	stats := rs.generator.calculateStatistics(activities)
	return &stats
}

// ValidateReportOptions 驗證報告選項
func (rs *ReportService) ValidateReportOptions(options types.ReportOptions) error {
	// 檢查時間範圍
	if !options.TimeRange.Start.IsZero() && !options.TimeRange.End.IsZero() {
		if options.TimeRange.Start.After(options.TimeRange.End) {
			return fmt.Errorf("開始時間不能晚於結束時間")
		}
	}

	// 檢查分組選項
	validGroupBy := []string{"", "activity", "date", "session"}
	valid := false
	for _, validOption := range validGroupBy {
		if options.GroupBy == validOption {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("無效的分組選項: %s", options.GroupBy)
	}

	return nil
}

// SetConfig 設定報告配置
func (rs *ReportService) SetConfig(config types.ReportConfig) {
	rs.config = config
	rs.generator.SetConfig(config)
}

// GetConfig 獲取報告配置
func (rs *ReportService) GetConfig() types.ReportConfig {
	return rs.config
}

// GetSupportedFormats 獲取支援的報告格式
func (rs *ReportService) GetSupportedFormats() []string {
	return []string{"json", "txt", "md"}
}

// GenerateQuickSummary 生成快速摘要
func (rs *ReportService) GenerateQuickSummary(activities []types.Activity) string {
	if len(activities) == 0 {
		return "無活動記錄"
	}

	totalTokens := 0
	activityCounts := make(map[types.ActivityType]int)

	for _, activity := range activities {
		totalTokens += activity.Tokens.TotalTokens
		activityCounts[activity.Type]++
	}

	avgTokens := float64(totalTokens) / float64(len(activities))

	summary := fmt.Sprintf("活動摘要: %d 個活動, %d 個 Token (平均 %.1f)\n",
		len(activities), totalTokens, avgTokens)

	summary += "活動分佈:\n"
	for activityType, count := range activityCounts {
		percentage := float64(count) / float64(len(activities)) * 100
		summary += fmt.Sprintf("- %s: %d (%.1f%%)\n", activityType, count, percentage)
	}

	return summary
}
