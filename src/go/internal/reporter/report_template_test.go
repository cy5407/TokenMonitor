package reporter

import (
	"strings"
	"testing"

	"token-monitor/internal/testutils"
	"token-monitor/internal/types"
)

// TestNewTemplateManager 測試模板管理器建立
func TestNewTemplateManager(t *testing.T) {
	tm := NewTemplateManager()

	if tm == nil {
		t.Fatal("模板管理器建立失敗")
	}

	// 檢查預設模板是否已註冊
	templates := tm.ListTemplates()
	expectedTemplates := []string{"basic", "detailed", "summary"}

	for _, expected := range expectedTemplates {
		found := false
		for _, template := range expectedTemplates {
			if template == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("預設模板 '%s' 未註冊", expected)
		}
	}
}

// TestRegisterTemplate 測試模板註冊
func TestRegisterTemplate(t *testing.T) {
	tm := NewTemplateManager()

	customTemplate := "Custom template: {{.TotalRecords}}"
	tm.RegisterTemplate("custom", customTemplate)

	template, err := tm.GetTemplate("custom")
	if err != nil {
		t.Fatalf("獲取自定義模板失敗: %v", err)
	}

	if template.GetName() != "custom" {
		t.Errorf("模板名稱錯誤: 期望 custom, 得到 %s", template.GetName())
	}

	if template.GetTemplate() != customTemplate {
		t.Errorf("模板內容錯誤")
	}
}

// TestGetTemplate 測試模板獲取
func TestGetTemplate(t *testing.T) {
	tm := NewTemplateManager()

	// 測試獲取存在的模板
	template, err := tm.GetTemplate("basic")
	if err != nil {
		t.Fatalf("獲取基礎模板失敗: %v", err)
	}

	if template.GetName() != "basic" {
		t.Errorf("模板名稱錯誤: 期望 basic, 得到 %s", template.GetName())
	}

	// 測試獲取不存在的模板
	_, err = tm.GetTemplate("nonexistent")
	if err == nil {
		t.Error("獲取不存在的模板應該返回錯誤")
	}
}

// TestRenderBasicTemplate 測試基礎模板渲染
func TestRenderBasicTemplate(t *testing.T) {
	tm := NewTemplateManager()

	// 建立測試報告
	report := testutils.CreateTestReport()

	// 渲染基礎模板
	result, err := tm.RenderReport("basic", report)
	if err != nil {
		t.Fatalf("渲染基礎模板失敗: %v", err)
	}

	// 驗證渲染結果
	if !strings.Contains(result, "Token 使用報告") {
		t.Error("渲染結果應包含標題")
	}

	if !strings.Contains(result, "總記錄數: 5") {
		t.Error("渲染結果應包含總記錄數")
	}

	if !strings.Contains(result, "總活動數: 5") {
		t.Error("渲染結果應包含總活動數")
	}

	if !strings.Contains(result, "總 Token 數: 1250") {
		t.Error("渲染結果應包含總 Token 數")
	}
}

// TestRenderDetailedTemplate 測試詳細模板渲染
func TestRenderDetailedTemplate(t *testing.T) {
	tm := NewTemplateManager()

	report := testutils.CreateTestReport()

	result, err := tm.RenderReport("detailed", report)
	if err != nil {
		t.Fatalf("渲染詳細模板失敗: %v", err)
	}

	// 驗證詳細模板特有內容
	if !strings.Contains(result, "詳細 Token 使用報告") {
		t.Error("渲染結果應包含詳細報告標題")
	}

	if !strings.Contains(result, "執行摘要") {
		t.Error("渲染結果應包含執行摘要")
	}

	if !strings.Contains(result, "統計分析") {
		t.Error("渲染結果應包含統計分析")
	}
}

// TestRenderSummaryTemplate 測試簡潔模板渲染
func TestRenderSummaryTemplate(t *testing.T) {
	tm := NewTemplateManager()

	report := testutils.CreateTestReport()

	result, err := tm.RenderReport("summary", report)
	if err != nil {
		t.Fatalf("渲染簡潔模板失敗: %v", err)
	}

	// 驗證簡潔模板格式
	if !strings.Contains(result, "========================================") {
		t.Error("渲染結果應包含分隔線")
	}

	if !strings.Contains(result, "總活動: 5") {
		t.Error("渲染結果應包含總活動數")
	}
}

// TestRenderActivityLoop 測試活動循環渲染
func TestRenderActivityLoop(t *testing.T) {
	template := NewReportTemplate("test", `
{{range $type, $report := .ByActivity}}
Activity: {{$type}}, Count: {{$report.Count}}, Tokens: {{$report.Tokens.TotalTokens}}
{{end}}
`)

	report := testutils.CreateTestReport()

	result, err := template.Render(report)
	if err != nil {
		t.Fatalf("渲染活動循環失敗: %v", err)
	}

	// 驗證循環內容
	if !strings.Contains(result, "Activity: coding") {
		t.Error("渲染結果應包含編程活動")
	}

	if !strings.Contains(result, "Activity: debugging") {
		t.Error("渲染結果應包含除錯活動")
	}

	if !strings.Contains(result, "Count: 3") {
		t.Error("渲染結果應包含正確的活動計數")
	}
}

// TestRenderHourlyDistributionLoop 測試小時分佈循環渲染
func TestRenderHourlyDistributionLoop(t *testing.T) {
	template := NewReportTemplate("test", `
{{range $hour, $count := .Statistics.ActivityTrends.HourlyDistribution}}
Hour: {{$hour}}, Count: {{$count}}
{{end}}
`)

	report := testutils.CreateTestReport()

	result, err := template.Render(report)
	if err != nil {
		t.Fatalf("渲染小時分佈循環失敗: %v", err)
	}

	// 驗證小時分佈內容
	if !strings.Contains(result, "Hour: 10") {
		t.Error("渲染結果應包含小時資訊")
	}
}

// TestRenderEfficiencyLoop 測試效率循環渲染
func TestRenderEfficiencyLoop(t *testing.T) {
	template := NewReportTemplate("test", `
{{range $type, $efficiency := .Statistics.EfficiencyMetrics.TokensPerActivity}}
Type: {{$type}}, Efficiency: {{$efficiency}}
{{end}}
`)

	report := testutils.CreateTestReport()

	result, err := template.Render(report)
	if err != nil {
		t.Fatalf("渲染效率循環失敗: %v", err)
	}

	// 驗證效率內容
	if !strings.Contains(result, "Type: coding") {
		t.Error("渲染結果應包含活動類型")
	}

	if !strings.Contains(result, "Efficiency:") {
		t.Error("渲染結果應包含效率資訊")
	}
}

// TestRenderEmptyReport 測試空報告渲染
func TestRenderEmptyReport(t *testing.T) {
	tm := NewTemplateManager()

	emptyReport := &types.BasicReport{
		GeneratedAt:  time.Now(),
		TotalRecords: 0,
		Summary: types.ReportSummary{
			TotalActivities: 0,
			TotalTokens:     types.TokenUsage{},
			ActivityCounts:  make(map[types.ActivityType]int),
		},
		ByActivity: make(map[types.ActivityType]types.ActivityReport),
		Statistics: types.ReportStatistics{},
	}

	result, err := tm.RenderReport("basic", emptyReport)
	if err != nil {
		t.Fatalf("渲染空報告失敗: %v", err)
	}

	// 驗證空報告處理
	if !strings.Contains(result, "總記錄數: 0") {
		t.Error("空報告應顯示 0 記錄")
	}

	if !strings.Contains(result, "總活動數: 0") {
		t.Error("空報告應顯示 0 活動")
	}
}

// TestCustomTemplateRendering 測試自定義模板渲染
func TestCustomTemplateRendering(t *testing.T) {
	tm := NewTemplateManager()

	// 註冊自定義模板
	customTemplate := `
Records: {{.TotalRecords}}
Activities: {{.Summary.TotalActivities}}
Tokens: {{.Summary.TotalTokens.TotalTokens}}
Average: {{.Summary.AverageTokensPerActivity}}
`

	tm.RegisterTemplate("custom", customTemplate)

	report := testutils.CreateTestReport()

	result, err := tm.RenderReport("custom", report)
	if err != nil {
		t.Fatalf("渲染自定義模板失敗: %v", err)
	}

	// 驗證自定義模板內容
	if !strings.Contains(result, "Records: 5") {
		t.Error("自定義模板應包含記錄數")
	}

	if !strings.Contains(result, "Activities: 5") {
		t.Error("自定義模板應包含活動數")
	}

	if !strings.Contains(result, "Tokens: 1250") {
		t.Error("自定義模板應包含 Token 數")
	}
}

