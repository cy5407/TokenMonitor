package reporter

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"token-monitor/internal/types"
)

// CSVGenerator CSV 格式報告生成器
type CSVGenerator struct {
	config types.ReportConfig
}

// NewCSVGenerator 建立新的 CSV 生成器
func NewCSVGenerator(config types.ReportConfig) *CSVGenerator {
	return &CSVGenerator{config: config}
}

// GenerateCSV 生成 CSV 格式報告
func (cg *CSVGenerator) GenerateCSV(report *types.BasicReport) ([]byte, error) {
	// TODO: 實作 CSV 生成邏輯
	var output strings.Builder
	writer := csv.NewWriter(&output)

	// 基本實作框架
	headers := []string{"活動類型", "數量", "Token總數", "平均Token", "佔比(%)"}
	writer.Write(headers)

	for activityType, activityReport := range report.ByActivity {
		row := []string{
			string(activityType),
			strconv.Itoa(activityReport.Count),
			strconv.Itoa(activityReport.Tokens.TotalTokens),
			fmt.Sprintf("%.2f", activityReport.AverageTokens),
			fmt.Sprintf("%.2f", activityReport.Percentage),
		}
		writer.Write(row)
	}

	writer.Flush()
	return []byte(output.String()), writer.Error()
}

// SaveCSV 儲存 CSV 報告到檔案
func (cg *CSVGenerator) SaveCSV(report *types.BasicReport, outputPath string) error {
	csvData, err := cg.GenerateCSV(report)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, csvData, 0644)
}

// writeReportHeader 寫入報告標題和摘要
func (cg *CSVGenerator) writeReportHeader(writer *csv.Writer, report *types.BasicReport) error {
	// 報告基本資訊
	writer.Write([]string{"Token 使用報告"})
	writer.Write([]string{"生成時間", report.GeneratedAt.Format("2006-01-02 15:04:05")})
	writer.Write([]string{"總記錄數", strconv.Itoa(report.TotalRecords)})
	writer.Write([]string{"總活動數", strconv.Itoa(report.Summary.TotalActivities)})
	writer.Write([]string{"總Token數", strconv.Itoa(report.Summary.TotalTokens.TotalTokens)})
	writer.Write([]string{"平均Token/活動", fmt.Sprintf("%.2f", report.Summary.AverageTokensPerActivity)})
	writer.Write([]string{}) // 空行分隔
	return nil
}

// writeActivityDetails 寫入活動詳細資訊
func (cg *CSVGenerator) writeActivityDetails(writer *csv.Writer, report *types.BasicReport) error {
	// 活動詳細資訊標題
	headers := []string{"活動類型", "數量", "輸入Token", "輸出Token", "總Token", "平均Token", "佔比(%)"}
	writer.Write(headers)

	for activityType, activityReport := range report.ByActivity {
		row := []string{
			string(activityType),
			strconv.Itoa(activityReport.Count),
			strconv.Itoa(activityReport.Tokens.InputTokens),
			strconv.Itoa(activityReport.Tokens.OutputTokens),
			strconv.Itoa(activityReport.Tokens.TotalTokens),
			fmt.Sprintf("%.2f", activityReport.AverageTokens),
			fmt.Sprintf("%.2f", activityReport.Percentage),
		}
		writer.Write(row)
	}

	writer.Write([]string{}) // 空行分隔
	return nil
}

// writeStatistics 寫入統計資訊
func (cg *CSVGenerator) writeStatistics(writer *csv.Writer, report *types.BasicReport) error {
	// Token 分佈統計
	writer.Write([]string{"統計項目", "數值"})
	writer.Write([]string{"Token總數", strconv.Itoa(report.Statistics.TokenDistribution.Total)})
	writer.Write([]string{"Token平均數", fmt.Sprintf("%.2f", report.Statistics.TokenDistribution.Average)})
	writer.Write([]string{"Token最小值", strconv.Itoa(report.Statistics.TokenDistribution.Min)})
	writer.Write([]string{"Token最大值", strconv.Itoa(report.Statistics.TokenDistribution.Max)})
	writer.Write([]string{"Token中位數", fmt.Sprintf("%.2f", report.Statistics.TokenDistribution.Median)})

	writer.Write([]string{}) // 空行分隔

	// 活動趨勢
	writer.Write([]string{"時間分析", "小時", "活動數量"})
	for hour := 0; hour < 24; hour++ {
		if count, exists := report.Statistics.ActivityTrends.HourlyDistribution[hour]; exists && count > 0 {
			writer.Write([]string{"小時分佈", fmt.Sprintf("%02d:00", hour), strconv.Itoa(count)})
		}
	}

	return nil
}
