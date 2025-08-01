package reporter

import (
	"os"
	"strings"
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestNewHTMLGenerator 測試 HTML 生成器建立
func TestNewHTMLGenerator(t *testing.T) {
	config := types.ReportConfig{Format: "html"}
	generator := NewHTMLGenerator(config)

	if generator == nil {
		t.Fatal("HTML 生成器建立失敗")
	}

	if generator.config.Format != "html" {
		t.Errorf("配置設定錯誤: 期望 html, 得到 %s", generator.config.Format)
	}

	if generator.template == nil {
		t.Error("HTML 模板應該被初始化")
	}
}

// TestGenerateHTML 測試 HTML 生成
func TestGenerateHTML(t *testing.T) {
	generator := NewHTMLGenerator(types.ReportConfig{})
	report := createTestHTMLReport()

	htmlData, err := generator.GenerateHTML(report)
	if err != nil {
		t.Fatalf("生成 HTML 失敗: %v", err)
	}

	htmlString := string(htmlData)

	// 驗證 HTML 結構
	if !strings.Contains(htmlString, "<!DOCTYPE html>") {
		t.Error("HTML 應包含 DOCTYPE 聲明")
	}

	if !strings.Contains(htmlString, "Token 使用分析報告") {
		t.Error("HTML 應包含報告標題")
	}

	if !strings.Contains(htmlString, "Chart.js") {
		t.Error("HTML 應包含 Chart.js 腳本")
	}

	// 驗證響應式設計
	if !strings.Contains(htmlString, "@media (max-width: 768px)") {
		t.Error("HTML 應包含響應式 CSS")
	}

	// 驗證資料綁定
	if !strings.Contains(htmlString, "coding") {
		t.Error("HTML 應包含活動類型資料")
	}
}

// TestSaveHTML 測試 HTML 檔案儲存
func TestSaveHTML(t *testing.T) {
	generator := NewHTMLGenerator(types.ReportConfig{})
	report := createTestHTMLReport()

	tempFile := t.TempDir() + "/test_report.html"

	err := generator.SaveHTML(report, tempFile)
	if err != nil {
		t.Fatalf("儲存 HTML 檔案失敗: %v", err)
	}

	// 驗證檔案是否存在
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Error("HTML 檔案應該被建立")
	}
}

// TestHTMLTemplateExecution 測試 HTML 模板執行
func TestHTMLTemplateExecution(t *testing.T) {
	generator := NewHTMLGenerator(types.ReportConfig{})

	// 測試空報告
	emptyReport := &types.BasicReport{
		GeneratedAt:  time.Now(),
		TotalRecords: 0,
		Summary: types.ReportSummary{
			TotalActivities:          0,
			TotalTokens:              types.TokenUsage{TotalTokens: 0},
			AverageTokensPerActivity: 0,
		},
		ByActivity: make(map[types.ActivityType]types.ActivityReport),
		Statistics: types.ReportStatistics{},
	}

	_, err := generator.GenerateHTML(emptyReport)
	if err != nil {
		t.Errorf("空報告應該能正常處理: %v", err)
	}
}

// TestHTMLChartData 測試圖表資料生成
func TestHTMLChartData(t *testing.T) {
	generator := NewHTMLGenerator(types.ReportConfig{})
	report := createTestHTMLReport()

	htmlData, err := generator.GenerateHTML(report)
	if err != nil {
		t.Fatalf("生成 HTML 失敗: %v", err)
	}

	htmlString := string(htmlData)

	// 驗證圓餅圖資料
	if !strings.Contains(htmlString, "activityChart") {
		t.Error("HTML 應包含活動圖表元素")
	}

	// 驗證趨勢圖資料
	if !strings.Contains(htmlString, "trendChart") {
		t.Error("HTML 應包含趨勢圖表元素")
	}

	// 驗證圖表配置
	if !strings.Contains(htmlString, "doughnut") {
		t.Error("HTML 應包含圓餅圖配置")
	}

	if !strings.Contains(htmlString, "line") {
		t.Error("HTML 應包含線圖配置")
	}
}

// createTestHTMLReport 建立測試 HTML 報告
func createTestHTMLReport() *types.BasicReport {
	return &types.BasicReport{
		GeneratedAt:  time.Now(),
		TotalRecords: 3,
		Summary: types.ReportSummary{
			TotalActivities:          3,
			TotalTokens:              types.TokenUsage{TotalTokens: 750},
			AverageTokensPerActivity: 250.0,
		},
		ByActivity: map[types.ActivityType]types.ActivityReport{
			types.ActivityCoding: {
				Count:         2,
				Tokens:        types.TokenUsage{TotalTokens: 500},
				AverageTokens: 250.0,
				Percentage:    66.7,
			},
			types.ActivityDebugging: {
				Count:         1,
				Tokens:        types.TokenUsage{TotalTokens: 250},
				AverageTokens: 250.0,
				Percentage:    33.3,
			},
		},
		Statistics: types.ReportStatistics{
			TokenDistribution: types.TokenDistributionStats{
				Total:   750,
				Average: 250.0,
				Min:     200,
				Max:     300,
				Median:  250.0,
			},
			ActivityTrends: types.ActivityTrendsStats{
				HourlyDistribution: map[int]int{
					9:  1,
					14: 2,
				},
				PeakHour:      14,
				PeakHourCount: 2,
			},
		},
	}
}
