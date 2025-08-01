package reporter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestNewReportService 測試報告服務建立
func TestNewReportService(t *testing.T) {
	config := types.ReportConfig{
		Format:     "json",
		OutputPath: "/tmp/report.json",
	}

	service := NewReportService(config)

	if service == nil {
		t.Fatal("報告服務建立失敗")
	}

	if service.config.Format != "json" {
		t.Errorf("配置設定錯誤: 期望 json, 得到 %s", service.config.Format)
	}
}

// TestGenerateReport 測試報告生成
func TestGenerateReport(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	options := types.ReportOptions{
		IncludeTrends: true,
	}

	result, err := service.GenerateReport(activities, options)
	if err != nil {
		t.Fatalf("生成報告失敗: %v", err)
	}

	// 驗證報告結果
	if result.BasicReport == nil {
		t.Error("基礎報告不應為空")
	}

	if len(result.JSONData) == 0 {
		t.Error("JSON 數據不應為空")
	}

	if result.TextReport == "" {
		t.Error("文字報告不應為空")
	}

	if result.BasicReport.TotalRecords != 3 {
		t.Errorf("總記錄數錯誤: 期望 3, 得到 %d", result.BasicReport.TotalRecords)
	}
}

// TestGenerateReportEmptyActivities 測試空活動列表的報告生成
func TestGenerateReportEmptyActivities(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	result, err := service.GenerateReport([]types.Activity{}, types.ReportOptions{})
	if err != nil {
		t.Fatalf("生成空報告失敗: %v", err)
	}

	if result.BasicReport.TotalRecords != 0 {
		t.Errorf("空報告總記錄數應為 0, 得到 %d", result.BasicReport.TotalRecords)
	}
}

// TestFilterActivitiesByTimeRange 測試時間範圍過濾
func TestFilterActivitiesByTimeRange(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	baseTime := time.Now()
	activities := []types.Activity{
		{ID: "1", Timestamp: baseTime},
		{ID: "2", Timestamp: baseTime.Add(1 * time.Hour)},
		{ID: "3", Timestamp: baseTime.Add(2 * time.Hour)},
		{ID: "4", Timestamp: baseTime.Add(3 * time.Hour)},
	}

	// 測試時間範圍過濾
	timeRange := types.TimeRange{
		Start: baseTime.Add(30 * time.Minute),
		End:   baseTime.Add(2*time.Hour + 30*time.Minute),
	}

	filtered := service.filterActivitiesByTimeRange(activities, timeRange)

	// 應該過濾出 ID 為 2 和 3 的活動
	if len(filtered) != 2 {
		t.Errorf("過濾結果數量錯誤: 期望 2, 得到 %d", len(filtered))
	}

	expectedIDs := []string{"2", "3"}
	for i, activity := range filtered {
		if activity.ID != expectedIDs[i] {
			t.Errorf("過濾結果順序錯誤: 期望 %s, 得到 %s", expectedIDs[i], activity.ID)
		}
	}
}

// TestFilterActivitiesNoTimeRange 測試無時間範圍過濾
func TestFilterActivitiesNoTimeRange(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	emptyTimeRange := types.TimeRange{}

	filtered := service.filterActivitiesByTimeRange(activities, emptyTimeRange)

	// 無時間範圍限制，應該返回所有活動
	if len(filtered) != len(activities) {
		t.Errorf("無時間範圍過濾應返回所有活動: 期望 %d, 得到 %d", len(activities), len(filtered))
	}
}

// TestSaveReport 測試報告儲存
func TestSaveReport(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	// 建立測試報告
	activities := createTestActivities()
	result, err := service.GenerateReport(activities, types.ReportOptions{})
	if err != nil {
		t.Fatalf("生成測試報告失敗: %v", err)
	}

	// 建立臨時目錄
	tempDir := t.TempDir()

	// 測試 JSON 格式儲存
	jsonPath := filepath.Join(tempDir, "test_report.json")
	err = service.SaveReport(result, jsonPath)
	if err != nil {
		t.Fatalf("儲存 JSON 報告失敗: %v", err)
	}

	// 驗證檔案是否存在
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Error("JSON 報告檔案未建立")
	}

	// 測試文字格式儲存
	txtPath := filepath.Join(tempDir, "test_report.txt")
	err = service.SaveReport(result, txtPath)
	if err != nil {
		t.Fatalf("儲存文字報告失敗: %v", err)
	}

	// 驗證檔案是否存在
	if _, err := os.Stat(txtPath); os.IsNotExist(err) {
		t.Error("文字報告檔案未建立")
	}
}

// TestGenerateAndSaveReport 測試生成並儲存報告
func TestGenerateAndSaveReport(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	options := types.ReportOptions{}

	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "combined_report.json")

	err := service.GenerateAndSaveReport(activities, options, outputPath)
	if err != nil {
		t.Fatalf("生成並儲存報告失敗: %v", err)
	}

	// 驗證檔案是否存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("組合報告檔案未建立")
	}
}

// TestGetAvailableTemplates 測試獲取可用模板
func TestGetAvailableTemplates(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	templates := service.GetAvailableTemplates()

	expectedTemplates := []string{"basic", "detailed", "summary"}
	for _, expected := range expectedTemplates {
		found := false
		for _, template := range templates {
			if template == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("預期模板 '%s' 未找到", expected)
		}
	}
}

// TestRegisterCustomTemplate 測試註冊自定義模板
func TestRegisterCustomTemplate(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	customTemplate := "Custom: {{.TotalRecords}} records"
	service.RegisterCustomTemplate("custom", customTemplate)

	templates := service.GetAvailableTemplates()

	found := false
	for _, template := range templates {
		if template == "custom" {
			found = true
			break
		}
	}

	if !found {
		t.Error("自定義模板未註冊成功")
	}
}

// TestGenerateReportWithTemplate 測試使用指定模板生成報告
func TestGenerateReportWithTemplate(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	options := types.ReportOptions{}

	// 測試基礎模板
	basicReport, err := service.GenerateReportWithTemplate(activities, "basic", options)
	if err != nil {
		t.Fatalf("使用基礎模板生成報告失敗: %v", err)
	}

	if !strings.Contains(basicReport, "Token 使用報告") {
		t.Error("基礎模板報告應包含標題")
	}

	// 測試詳細模板
	detailedReport, err := service.GenerateReportWithTemplate(activities, "detailed", options)
	if err != nil {
		t.Fatalf("使用詳細模板生成報告失敗: %v", err)
	}

	if !strings.Contains(detailedReport, "詳細 Token 使用報告") {
		t.Error("詳細模板報告應包含詳細標題")
	}
}

// TestGetReportStatistics 測試獲取報告統計
func TestGetReportStatistics(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	stats := service.GetReportStatistics(activities)

	if stats == nil {
		t.Fatal("統計資訊不應為空")
	}

	// 驗證統計資訊
	if stats.TokenDistribution.Total != 900 {
		t.Errorf("Token 總數錯誤: 期望 900, 得到 %d", stats.TokenDistribution.Total)
	}
}

// TestValidateReportOptions 測試報告選項驗證
func TestValidateReportOptions(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	// 測試有效選項
	validOptions := types.ReportOptions{
		TimeRange: types.TimeRange{
			Start: time.Now(),
			End:   time.Now().Add(1 * time.Hour),
		},
		GroupBy: "activity",
	}

	err := service.ValidateReportOptions(validOptions)
	if err != nil {
		t.Errorf("有效選項驗證失敗: %v", err)
	}

	// 測試無效時間範圍
	invalidTimeOptions := types.ReportOptions{
		TimeRange: types.TimeRange{
			Start: time.Now().Add(1 * time.Hour),
			End:   time.Now(),
		},
	}

	err = service.ValidateReportOptions(invalidTimeOptions)
	if err == nil {
		t.Error("無效時間範圍應該返回錯誤")
	}

	// 測試無效分組選項
	invalidGroupOptions := types.ReportOptions{
		GroupBy: "invalid_group",
	}

	err = service.ValidateReportOptions(invalidGroupOptions)
	if err == nil {
		t.Error("無效分組選項應該返回錯誤")
	}
}

// TestGenerateQuickSummary 測試快速摘要生成
func TestGenerateQuickSummary(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	activities := createTestActivities()
	summary := service.GenerateQuickSummary(activities)

	// 驗證摘要內容
	if !strings.Contains(summary, "3 個活動") {
		t.Error("摘要應包含活動數量")
	}

	if !strings.Contains(summary, "900 個 Token") {
		t.Error("摘要應包含 Token 總數")
	}

	if !strings.Contains(summary, "活動分佈") {
		t.Error("摘要應包含活動分佈")
	}

	// 測試空活動列表
	emptySummary := service.GenerateQuickSummary([]types.Activity{})
	if emptySummary != "無活動記錄" {
		t.Errorf("空活動摘要錯誤: 期望 '無活動記錄', 得到 '%s'", emptySummary)
	}
}

// TestGetSupportedFormats 測試獲取支援格式
func TestGetSupportedFormats(t *testing.T) {
	service := NewReportService(types.ReportConfig{})

	formats := service.GetSupportedFormats()
	expectedFormats := []string{"json", "txt", "md"}

	if len(formats) != len(expectedFormats) {
		t.Errorf("支援格式數量錯誤: 期望 %d, 得到 %d", len(expectedFormats), len(formats))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("預期格式 '%s' 未找到", expected)
		}
	}
}

// TestServiceConfigManagement 測試服務配置管理
func TestServiceConfigManagement(t *testing.T) {
	initialConfig := types.ReportConfig{
		Format:     "json",
		OutputPath: "/tmp/initial.json",
	}

	service := NewReportService(initialConfig)

	// 測試獲取配置
	config := service.GetConfig()
	if config.Format != "json" {
		t.Errorf("獲取配置錯誤: 期望 json, 得到 %s", config.Format)
	}

	// 測試設定配置
	newConfig := types.ReportConfig{
		Format:     "txt",
		OutputPath: "/tmp/new.txt",
	}

	service.SetConfig(newConfig)
	updatedConfig := service.GetConfig()

	if updatedConfig.Format != "txt" {
		t.Errorf("設定配置錯誤: 期望 txt, 得到 %s", updatedConfig.Format)
	}
}

// createTestActivities 建立測試活動數據
func createTestActivities() []types.Activity {
	baseTime := time.Now()

	return []types.Activity{
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
}
