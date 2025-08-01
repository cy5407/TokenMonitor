package reporter

import (
	"encoding/json"
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestNewReportGenerator 測試報告生成器建立
func TestNewReportGenerator(t *testing.T) {
	config := types.ReportConfig{
		Format:     "json",
		OutputPath: "/tmp/report.json",
	}

	generator := NewReportGenerator(config)

	if generator == nil {
		t.Fatal("報告生成器建立失敗")
	}

	if generator.config.Format != "json" {
		t.Errorf("配置設定錯誤: 期望 json, 得到 %s", generator.config.Format)
	}
}

// TestGenerateBasicReport 測試基礎報告生成
func TestGenerateBasicReport(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	// 建立測試數據
	baseTime := time.Now()
	activities := []types.Activity{
		{
			ID:        "1",
			Type:      types.ActivityCoding,
			Content:   "實作功能",
			Tokens:    types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300},
			Timestamp: baseTime,
		},
		{
			ID:        "2",
			Type:      types.ActivityCoding,
			Content:   "另一個功能",
			Tokens:    types.TokenUsage{InputTokens: 150, OutputTokens: 250, TotalTokens: 400},
			Timestamp: baseTime.Add(1 * time.Hour),
		},
		{
			ID:        "3",
			Type:      types.ActivityDebugging,
			Content:   "修復錯誤",
			Tokens:    types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200},
			Timestamp: baseTime.Add(2 * time.Hour),
		},
	}

	report, err := generator.GenerateBasicReport(activities)
	if err != nil {
		t.Fatalf("生成基礎報告失敗: %v", err)
	}

	// 驗證基本資訊
	if report.TotalRecords != 3 {
		t.Errorf("總記錄數錯誤: 期望 3, 得到 %d", report.TotalRecords)
	}

	// 驗證摘要
	if report.Summary.TotalActivities != 3 {
		t.Errorf("總活動數錯誤: 期望 3, 得到 %d", report.Summary.TotalActivities)
	}

	expectedTotalTokens := types.TokenUsage{InputTokens: 330, OutputTokens: 570, TotalTokens: 900}
	if report.Summary.TotalTokens != expectedTotalTokens {
		t.Errorf("總 Token 統計錯誤: 期望 %+v, 得到 %+v", expectedTotalTokens, report.Summary.TotalTokens)
	}

	// 驗證活動計數
	if report.Summary.ActivityCounts[types.ActivityCoding] != 2 {
		t.Errorf("編程活動計數錯誤: 期望 2, 得到 %d", report.Summary.ActivityCounts[types.ActivityCoding])
	}

	if report.Summary.ActivityCounts[types.ActivityDebugging] != 1 {
		t.Errorf("除錯活動計數錯誤: 期望 1, 得到 %d", report.Summary.ActivityCounts[types.ActivityDebugging])
	}

	// 驗證平均 Token 數
	expectedAvg := 300.0 // 900 / 3
	if report.Summary.AverageTokensPerActivity != expectedAvg {
		t.Errorf("平均 Token 數錯誤: 期望 %.1f, 得到 %.1f", expectedAvg, report.Summary.AverageTokensPerActivity)
	}
}

// TestGenerateBasicReportEmptyInput 測試空輸入的基礎報告生成
func TestGenerateBasicReportEmptyInput(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	report, err := generator.GenerateBasicReport([]types.Activity{})
	if err != nil {
		t.Fatalf("生成空報告失敗: %v", err)
	}

	if report.TotalRecords != 0 {
		t.Errorf("空報告總記錄數應為 0, 得到 %d", report.TotalRecords)
	}

	if report.Summary.TotalActivities != 0 {
		t.Errorf("空報告總活動數應為 0, 得到 %d", report.Summary.TotalActivities)
	}
}

// TestGenerateJSONReport 測試 JSON 報告生成
func TestGenerateJSONReport(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	activities := []types.Activity{
		{
			ID:        "1",
			Type:      types.ActivityCoding,
			Content:   "測試功能",
			Tokens:    types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300},
			Timestamp: time.Now(),
		},
	}

	jsonData, err := generator.GenerateJSONReport(activities)
	if err != nil {
		t.Fatalf("生成 JSON 報告失敗: %v", err)
	}

	// 驗證 JSON 格式
	var report types.BasicReport
	err = json.Unmarshal(jsonData, &report)
	if err != nil {
		t.Fatalf("JSON 解析失敗: %v", err)
	}

	if report.TotalRecords != 1 {
		t.Errorf("JSON 報告總記錄數錯誤: 期望 1, 得到 %d", report.TotalRecords)
	}
}

// TestCalculateByActivity 測試按活動類型統計
func TestCalculateByActivity(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	activities := []types.Activity{
		{Type: types.ActivityCoding, Tokens: types.TokenUsage{TotalTokens: 300}},
		{Type: types.ActivityCoding, Tokens: types.TokenUsage{TotalTokens: 400}},
		{Type: types.ActivityDebugging, Tokens: types.TokenUsage{TotalTokens: 200}},
	}

	byActivity := generator.calculateByActivity(activities)

	// 驗證編程活動統計
	codingReport := byActivity[types.ActivityCoding]
	if codingReport.Count != 2 {
		t.Errorf("編程活動計數錯誤: 期望 2, 得到 %d", codingReport.Count)
	}

	if codingReport.Tokens.TotalTokens != 700 {
		t.Errorf("編程活動 Token 總數錯誤: 期望 700, 得到 %d", codingReport.Tokens.TotalTokens)
	}

	expectedAvg := 350.0 // 700 / 2
	if codingReport.AverageTokens != expectedAvg {
		t.Errorf("編程活動平均 Token 錯誤: 期望 %.1f, 得到 %.1f", expectedAvg, codingReport.AverageTokens)
	}

	expectedPercentage := 66.67 // 2/3 * 100, 四捨五入
	tolerance := 0.1
	if codingReport.Percentage < expectedPercentage-tolerance || codingReport.Percentage > expectedPercentage+tolerance {
		t.Errorf("編程活動百分比錯誤: 期望約 %.2f%%, 得到 %.2f%%", expectedPercentage, codingReport.Percentage)
	}
}

// TestCalculateTokenDistribution 測試 Token 分佈計算
func TestCalculateTokenDistribution(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	activities := []types.Activity{
		{Tokens: types.TokenUsage{InputTokens: 100, OutputTokens: 200, TotalTokens: 300}},
		{Tokens: types.TokenUsage{InputTokens: 150, OutputTokens: 250, TotalTokens: 400}},
		{Tokens: types.TokenUsage{InputTokens: 80, OutputTokens: 120, TotalTokens: 200}},
	}

	distribution := generator.calculateTokenDistribution(activities)

	// 驗證總計
	if distribution.Total != 900 {
		t.Errorf("Token 總數錯誤: 期望 900, 得到 %d", distribution.Total)
	}

	if distribution.Input != 330 {
		t.Errorf("輸入 Token 總數錯誤: 期望 330, 得到 %d", distribution.Input)
	}

	if distribution.Output != 570 {
		t.Errorf("輸出 Token 總數錯誤: 期望 570, 得到 %d", distribution.Output)
	}

	// 驗證統計值
	expectedAvg := 300.0 // 900 / 3
	if distribution.Average != expectedAvg {
		t.Errorf("平均 Token 數錯誤: 期望 %.1f, 得到 %.1f", expectedAvg, distribution.Average)
	}

	if distribution.Min != 200 {
		t.Errorf("最小 Token 數錯誤: 期望 200, 得到 %d", distribution.Min)
	}

	if distribution.Max != 400 {
		t.Errorf("最大 Token 數錯誤: 期望 400, 得到 %d", distribution.Max)
	}

	expectedMedian := 300.0 // 中位數
	if distribution.Median != expectedMedian {
		t.Errorf("中位數錯誤: 期望 %.1f, 得到 %.1f", expectedMedian, distribution.Median)
	}
}

// TestCalculateActivityTrends 測試活動趨勢計算
func TestCalculateActivityTrends(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	baseTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	activities := []types.Activity{
		{Timestamp: baseTime},                    // 10:00
		{Timestamp: baseTime.Add(1 * time.Hour)}, // 11:00
		{Timestamp: baseTime.Add(1 * time.Hour)}, // 11:00 (重複)
		{Timestamp: baseTime.Add(2 * time.Hour)}, // 12:00
	}

	trends := generator.calculateActivityTrends(activities)

	// 驗證小時分佈
	if trends.HourlyDistribution[10] != 1 {
		t.Errorf("10 點活動數錯誤: 期望 1, 得到 %d", trends.HourlyDistribution[10])
	}

	if trends.HourlyDistribution[11] != 2 {
		t.Errorf("11 點活動數錯誤: 期望 2, 得到 %d", trends.HourlyDistribution[11])
	}

	if trends.HourlyDistribution[12] != 1 {
		t.Errorf("12 點活動數錯誤: 期望 1, 得到 %d", trends.HourlyDistribution[12])
	}

	// 驗證峰值時間
	if trends.PeakHour != 11 {
		t.Errorf("峰值小時錯誤: 期望 11, 得到 %d", trends.PeakHour)
	}

	if trends.PeakHourCount != 2 {
		t.Errorf("峰值小時活動數錯誤: 期望 2, 得到 %d", trends.PeakHourCount)
	}
}

// TestCalculateEfficiencyMetrics 測試效率指標計算
func TestCalculateEfficiencyMetrics(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	activities := []types.Activity{
		{Type: types.ActivityCoding, Tokens: types.TokenUsage{TotalTokens: 300}},
		{Type: types.ActivityCoding, Tokens: types.TokenUsage{TotalTokens: 500}},
		{Type: types.ActivityDebugging, Tokens: types.TokenUsage{TotalTokens: 200}},
	}

	metrics := generator.calculateEfficiencyMetrics(activities)

	// 驗證編程活動效率
	expectedCodingEfficiency := 400.0 // (300 + 500) / 2
	if metrics.TokensPerActivity[types.ActivityCoding] != expectedCodingEfficiency {
		t.Errorf("編程活動效率錯誤: 期望 %.1f, 得到 %.1f",
			expectedCodingEfficiency, metrics.TokensPerActivity[types.ActivityCoding])
	}

	// 驗證除錯活動效率
	expectedDebuggingEfficiency := 200.0 // 200 / 1
	if metrics.TokensPerActivity[types.ActivityDebugging] != expectedDebuggingEfficiency {
		t.Errorf("除錯活動效率錯誤: 期望 %.1f, 得到 %.1f",
			expectedDebuggingEfficiency, metrics.TokensPerActivity[types.ActivityDebugging])
	}
}

// TestFindMinMax 測試最小最大值計算
func TestFindMinMax(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	values := []int{300, 100, 500, 200, 400}
	min, max := generator.findMinMax(values)

	if min != 100 {
		t.Errorf("最小值錯誤: 期望 100, 得到 %d", min)
	}

	if max != 500 {
		t.Errorf("最大值錯誤: 期望 500, 得到 %d", max)
	}

	// 測試空陣列
	emptyMin, emptyMax := generator.findMinMax([]int{})
	if emptyMin != 0 || emptyMax != 0 {
		t.Errorf("空陣列應返回 0, 0, 得到 %d, %d", emptyMin, emptyMax)
	}
}

// TestCalculateMedian 測試中位數計算
func TestCalculateMedian(t *testing.T) {
	generator := NewReportGenerator(types.ReportConfig{})

	// 測試奇數個元素
	oddValues := []int{300, 100, 500, 200, 400}
	oddMedian := generator.calculateMedian(oddValues)
	expectedOddMedian := 300.0
	if oddMedian != expectedOddMedian {
		t.Errorf("奇數中位數錯誤: 期望 %.1f, 得到 %.1f", expectedOddMedian, oddMedian)
	}

	// 測試偶數個元素
	evenValues := []int{100, 200, 300, 400}
	evenMedian := generator.calculateMedian(evenValues)
	expectedEvenMedian := 250.0 // (200 + 300) / 2
	if evenMedian != expectedEvenMedian {
		t.Errorf("偶數中位數錯誤: 期望 %.1f, 得到 %.1f", expectedEvenMedian, evenMedian)
	}

	// 測試空陣列
	emptyMedian := generator.calculateMedian([]int{})
	if emptyMedian != 0 {
		t.Errorf("空陣列中位數應為 0, 得到 %.1f", emptyMedian)
	}
}

// TestConfigManagement 測試配置管理
func TestConfigManagement(t *testing.T) {
	initialConfig := types.ReportConfig{
		Format:     "json",
		OutputPath: "/tmp/initial.json",
	}

	generator := NewReportGenerator(initialConfig)

	// 測試獲取配置
	config := generator.GetConfig()
	if config.Format != "json" {
		t.Errorf("獲取配置錯誤: 期望 json, 得到 %s", config.Format)
	}

	// 測試設定配置
	newConfig := types.ReportConfig{
		Format:     "csv",
		OutputPath: "/tmp/new.csv",
	}

	generator.SetConfig(newConfig)
	updatedConfig := generator.GetConfig()

	if updatedConfig.Format != "csv" {
		t.Errorf("設定配置錯誤: 期望 csv, 得到 %s", updatedConfig.Format)
	}

	if updatedConfig.OutputPath != "/tmp/new.csv" {
		t.Errorf("設定輸出路徑錯誤: 期望 /tmp/new.csv, 得到 %s", updatedConfig.OutputPath)
	}
}
