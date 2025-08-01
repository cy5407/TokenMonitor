package reporter

import (
	"fmt"
	"strings"

	"token-monitor/internal/types"
)

// ReportTemplate 報告模板
type ReportTemplate struct {
	name     string
	template string
}

// NewReportTemplate 建立新的報告模板
func NewReportTemplate(name, template string) *ReportTemplate {
	return &ReportTemplate{
		name:     name,
		template: template,
	}
}

// TemplateManager 模板管理器
type TemplateManager struct {
	templates map[string]*ReportTemplate
}

// NewTemplateManager 建立新的模板管理器
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*ReportTemplate),
	}

	// 註冊預設模板
	tm.registerDefaultTemplates()
	return tm
}

// registerDefaultTemplates 註冊預設模板
func (tm *TemplateManager) registerDefaultTemplates() {
	// 基礎報告模板
	basicTemplate := `# Token 使用報告

## 報告摘要
- 生成時間: {{.GeneratedAt}}
- 總記錄數: {{.TotalRecords}}
- 時間範圍: {{.TimeRange.Start}} 至 {{.TimeRange.End}}

## 活動統計
- 總活動數: {{.Summary.TotalActivities}}
- 總 Token 數: {{.Summary.TotalTokens.TotalTokens}}
- 平均每活動 Token 數: {{.Summary.AverageTokensPerActivity}}

## 按活動類型統計
{{range $type, $report := .ByActivity}}
### {{$type}}
- 活動數量: {{$report.Count}}
- Token 使用: {{$report.Tokens.TotalTokens}}
- 平均 Token: {{$report.AverageTokens}}
- 佔比: {{$report.Percentage}}%
{{end}}

## Token 分佈統計
- 總計: {{.Statistics.TokenDistribution.Total}}
- 平均: {{.Statistics.TokenDistribution.Average}}
- 最小: {{.Statistics.TokenDistribution.Min}}
- 最大: {{.Statistics.TokenDistribution.Max}}
- 中位數: {{.Statistics.TokenDistribution.Median}}
`

	tm.RegisterTemplate("basic", basicTemplate)

	// 詳細報告模板
	detailedTemplate := `# 詳細 Token 使用報告

## 基本資訊
**報告生成時間**: {{.GeneratedAt}}
**分析期間**: {{.TimeRange.Start}} - {{.TimeRange.End}}
**總記錄數**: {{.TotalRecords}}

## 執行摘要
本報告分析了 {{.Summary.TotalActivities}} 個活動，總共使用了 {{.Summary.TotalTokens.TotalTokens}} 個 Token。
平均每個活動使用 {{.Summary.AverageTokensPerActivity}} 個 Token。

## 活動類型分析
{{range $type, $report := .ByActivity}}
### {{$type}} 活動
- **活動數量**: {{$report.Count}} ({{$report.Percentage}}%)
- **Token 使用量**: 
  - 輸入: {{$report.Tokens.InputTokens}}
  - 輸出: {{$report.Tokens.OutputTokens}}
  - 總計: {{$report.Tokens.TotalTokens}}
- **平均每活動**: {{$report.AverageTokens}} Token

{{end}}

## 統計分析

### Token 分佈
- **總使用量**: {{.Statistics.TokenDistribution.Total}}
- **統計指標**:
  - 平均值: {{.Statistics.TokenDistribution.Average}}
  - 最小值: {{.Statistics.TokenDistribution.Min}}
  - 最大值: {{.Statistics.TokenDistribution.Max}}
  - 中位數: {{.Statistics.TokenDistribution.Median}}

### 活動趨勢
- **峰值時間**: {{.Statistics.ActivityTrends.PeakHour}}:00 ({{.Statistics.ActivityTrends.PeakHourCount}} 個活動)
- **時間分佈**: 
{{range $hour, $count := .Statistics.ActivityTrends.HourlyDistribution}}
  - {{$hour}}:00 - {{$count}} 個活動
{{end}}

### 效率指標
{{range $type, $efficiency := .Statistics.EfficiencyMetrics.TokensPerActivity}}
- **{{$type}}**: 平均 {{$efficiency}} Token/活動
{{end}}
`

	tm.RegisterTemplate("detailed", detailedTemplate)

	// 簡潔報告模板
	summaryTemplate := `Token 使用報告 ({{.GeneratedAt}})
========================================

總活動: {{.Summary.TotalActivities}}
總 Token: {{.Summary.TotalTokens.TotalTokens}}
平均 Token/活動: {{.Summary.AverageTokensPerActivity}}

活動分佈:
{{range $type, $report := .ByActivity}}
- {{$type}}: {{$report.Count}} ({{$report.Percentage}}%)
{{end}}

Token 統計:
- 最小: {{.Statistics.TokenDistribution.Min}}
- 最大: {{.Statistics.TokenDistribution.Max}}
- 平均: {{.Statistics.TokenDistribution.Average}}
- 中位數: {{.Statistics.TokenDistribution.Median}}
`

	tm.RegisterTemplate("summary", summaryTemplate)
}

// RegisterTemplate 註冊模板
func (tm *TemplateManager) RegisterTemplate(name, template string) {
	tm.templates[name] = NewReportTemplate(name, template)
}

// GetTemplate 獲取模板
func (tm *TemplateManager) GetTemplate(name string) (*ReportTemplate, error) {
	template, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("模板 '%s' 不存在", name)
	}
	return template, nil
}

// ListTemplates 列出所有模板
func (tm *TemplateManager) ListTemplates() []string {
	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// RenderReport 渲染報告
func (tm *TemplateManager) RenderReport(templateName string, report *types.BasicReport) (string, error) {
	template, err := tm.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	return template.Render(report)
}

// Render 渲染模板
func (rt *ReportTemplate) Render(report *types.BasicReport) (string, error) {
	result := rt.template

	// 替換基本欄位
	result = strings.ReplaceAll(result, "{{.GeneratedAt}}", report.GeneratedAt.Format("2006-01-02 15:04:05"))
	result = strings.ReplaceAll(result, "{{.TotalRecords}}", fmt.Sprintf("%d", report.TotalRecords))

	// 替換時間範圍
	result = strings.ReplaceAll(result, "{{.TimeRange.Start}}", report.TimeRange.Start.Format("2006-01-02 15:04:05"))
	result = strings.ReplaceAll(result, "{{.TimeRange.End}}", report.TimeRange.End.Format("2006-01-02 15:04:05"))

	// 替換摘要資訊
	result = strings.ReplaceAll(result, "{{.Summary.TotalActivities}}", fmt.Sprintf("%d", report.Summary.TotalActivities))
	result = strings.ReplaceAll(result, "{{.Summary.TotalTokens.TotalTokens}}", fmt.Sprintf("%d", report.Summary.TotalTokens.TotalTokens))
	result = strings.ReplaceAll(result, "{{.Summary.AverageTokensPerActivity}}", fmt.Sprintf("%.2f", report.Summary.AverageTokensPerActivity))

	// 替換統計資訊
	result = strings.ReplaceAll(result, "{{.Statistics.TokenDistribution.Total}}", fmt.Sprintf("%d", report.Statistics.TokenDistribution.Total))
	result = strings.ReplaceAll(result, "{{.Statistics.TokenDistribution.Average}}", fmt.Sprintf("%.2f", report.Statistics.TokenDistribution.Average))
	result = strings.ReplaceAll(result, "{{.Statistics.TokenDistribution.Min}}", fmt.Sprintf("%d", report.Statistics.TokenDistribution.Min))
	result = strings.ReplaceAll(result, "{{.Statistics.TokenDistribution.Max}}", fmt.Sprintf("%d", report.Statistics.TokenDistribution.Max))
	result = strings.ReplaceAll(result, "{{.Statistics.TokenDistribution.Median}}", fmt.Sprintf("%.2f", report.Statistics.TokenDistribution.Median))

	// 替換活動趨勢
	result = strings.ReplaceAll(result, "{{.Statistics.ActivityTrends.PeakHour}}", fmt.Sprintf("%d", report.Statistics.ActivityTrends.PeakHour))
	result = strings.ReplaceAll(result, "{{.Statistics.ActivityTrends.PeakHourCount}}", fmt.Sprintf("%d", report.Statistics.ActivityTrends.PeakHourCount))

	// 處理循環結構
	result = rt.renderActivityLoop(result, report.ByActivity)
	result = rt.renderHourlyDistributionLoop(result, report.Statistics.ActivityTrends.HourlyDistribution)
	result = rt.renderEfficiencyLoop(result, report.Statistics.EfficiencyMetrics.TokensPerActivity)

	return result, nil
}

// renderActivityLoop 渲染活動循環
func (rt *ReportTemplate) renderActivityLoop(template string, byActivity map[types.ActivityType]types.ActivityReport) string {
	// 找到循環標記
	start := "{{range $type, $report := .ByActivity}}"
	end := "{{end}}"

	startIdx := strings.Index(template, start)
	if startIdx == -1 {
		return template
	}

	endIdx := strings.Index(template[startIdx:], end)
	if endIdx == -1 {
		return template
	}

	endIdx += startIdx

	// 提取循環內容
	loopContent := template[startIdx+len(start) : endIdx]

	// 生成循環結果
	var result strings.Builder
	for activityType, report := range byActivity {
		loopResult := loopContent
		loopResult = strings.ReplaceAll(loopResult, "{{$type}}", string(activityType))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.Count}}", fmt.Sprintf("%d", report.Count))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.Tokens.InputTokens}}", fmt.Sprintf("%d", report.Tokens.InputTokens))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.Tokens.OutputTokens}}", fmt.Sprintf("%d", report.Tokens.OutputTokens))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.Tokens.TotalTokens}}", fmt.Sprintf("%d", report.Tokens.TotalTokens))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.AverageTokens}}", fmt.Sprintf("%.2f", report.AverageTokens))
		loopResult = strings.ReplaceAll(loopResult, "{{$report.Percentage}}", fmt.Sprintf("%.2f", report.Percentage))

		result.WriteString(loopResult)
	}

	// 替換原始模板
	return template[:startIdx] + result.String() + template[endIdx+len(end):]
}

// renderHourlyDistributionLoop 渲染小時分佈循環
func (rt *ReportTemplate) renderHourlyDistributionLoop(template string, hourlyDist map[int]int) string {
	start := "{{range $hour, $count := .Statistics.ActivityTrends.HourlyDistribution}}"
	end := "{{end}}"

	startIdx := strings.Index(template, start)
	if startIdx == -1 {
		return template
	}

	endIdx := strings.Index(template[startIdx:], end)
	if endIdx == -1 {
		return template
	}

	endIdx += startIdx

	loopContent := template[startIdx+len(start) : endIdx]

	var result strings.Builder
	for hour := 0; hour < 24; hour++ {
		if count, exists := hourlyDist[hour]; exists && count > 0 {
			loopResult := loopContent
			loopResult = strings.ReplaceAll(loopResult, "{{$hour}}", fmt.Sprintf("%d", hour))
			loopResult = strings.ReplaceAll(loopResult, "{{$count}}", fmt.Sprintf("%d", count))
			result.WriteString(loopResult)
		}
	}

	return template[:startIdx] + result.String() + template[endIdx+len(end):]
}

// renderEfficiencyLoop 渲染效率循環
func (rt *ReportTemplate) renderEfficiencyLoop(template string, efficiency map[types.ActivityType]float64) string {
	start := "{{range $type, $efficiency := .Statistics.EfficiencyMetrics.TokensPerActivity}}"
	end := "{{end}}"

	startIdx := strings.Index(template, start)
	if startIdx == -1 {
		return template
	}

	endIdx := strings.Index(template[startIdx:], end)
	if endIdx == -1 {
		return template
	}

	endIdx += startIdx

	loopContent := template[startIdx+len(start) : endIdx]

	var result strings.Builder
	for activityType, eff := range efficiency {
		loopResult := loopContent
		loopResult = strings.ReplaceAll(loopResult, "{{$type}}", string(activityType))
		loopResult = strings.ReplaceAll(loopResult, "{{$efficiency}}", fmt.Sprintf("%.2f", eff))
		result.WriteString(loopResult)
	}

	return template[:startIdx] + result.String() + template[endIdx+len(end):]
}

// GetName 獲取模板名稱
func (rt *ReportTemplate) GetName() string {
	return rt.name
}

// GetTemplate 獲取模板內容
func (rt *ReportTemplate) GetTemplate() string {
	return rt.template
}
