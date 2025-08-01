package integration

import (
	"os"
	"testing"
	"time"

	"token-monitor/internal/analyzer"
	"token-monitor/internal/calculator"
	"token-monitor/internal/cost"
	"token-monitor/internal/reporter"
	"token-monitor/internal/storage"
	"token-monitor/internal/types"
)

// TestFullWorkflow 測試完整工作流程
func TestFullWorkflow(t *testing.T) {
	// 建立臨時目錄
	tempDir := t.TempDir()

	// 初始化各個元件
	tokenCalc, _ := calculator.NewTokenCalculator(calculator.WithMethod("estimation"))
	activityAnalyzer := analyzer.NewActivityAnalyzer(analyzer.AnalyzerConfig{})
	costCalc, _ := cost.NewCostCalculatorFromConfig(map[string]interface{}{})
	jsonStorage := storage.NewJSONStorage(tempDir)
	reportGen := reporter.NewReportGenerator()

	// 模擬使用記錄
	activities := []types.Activity{
		{
			Type:        types.ActivityCoding,
			Content:     "實作新功能的程式碼",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Tokens:      types.TokenUsage{InputTokens: 150, OutputTokens: 200, TotalTokens: 350},
		},
		{
			Type:        types.ActivityDebugging,
			Content:     "修復程式錯誤",
			Timestamp:   time.Now().Add(-1 * time.Hour),
			Tokens:      types.TokenUsage{InputTokens: 100, OutputTokens: 150, TotalTokens: 250},
		},
	}

	// 步驟 1: Token 計算 (already done in mock data)

	// 步驟 2: 活動分析
	analysisResult := activityAnalyzer.Analyze(activities)
	if analysisResult.TotalActivities != len(activities) {
		t.Errorf("活動分析錯誤: 期望 %d, 得到 %d", len(activities), analysisResult.TotalActivities)
	}

	// 步驟 3: 成本計算
	for i := range activities {
		costResult, _ := costCalc.CalculateCost(activities[i].Tokens.InputTokens, activities[i].Tokens.OutputTokens, "claude-sonnet-4.0")
		activities[i].Cost = costResult.TotalCost
	}

	// 步驟 4: 資料儲存
	if err := jsonStorage.SaveActivityData(activities); err != nil {
		t.Fatalf("儲存活動資料失敗: %v", err)
	}

	// 步驟 5: 資料載入
	timeRange := types.TimeRange{
		Start: time.Now().Add(-3 * time.Hour),
		End:   time.Now(),
	}

	loadedActivities, err := jsonStorage.LoadActivityData(timeRange)
	if err != nil {
		t.Fatalf("載入活動資料失敗: %v", err)
	}

	if len(loadedActivities) != len(activities) {
		t.Errorf("載入活動數量錯誤: 期望 %d, 得到 %d", len(activities), len(loadedActivities))
	}

	// 步驟 6: 報告生成
	reportData, err := reportGen.Generate(loadedActivities, types.ReportOptions{})
	if err != nil {
		t.Fatalf("生成報告失敗: %v", err)
	}

	// 驗證報告內容
	if reportData.TotalRecords != len(activities) {
		t.Errorf("報告記錄數錯誤: 期望 %d, 得到 %d", len(activities), reportData.TotalRecords)
	}

	if len(reportData.ByActivity) == 0 {
		t.Error("報告應包含活動分析")
	}
}