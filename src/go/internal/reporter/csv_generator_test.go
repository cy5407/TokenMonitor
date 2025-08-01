package reporter

import (
	"strings"
	"testing"

	"token-monitor/internal/testutils"
	"token-monitor/internal/types"
)

// TestNewCSVGenerator 測試 CSV 生成器建立
func TestNewCSVGenerator(t *testing.T) {
	config := types.ReportConfig{Format: "csv"}
	generator := NewCSVGenerator(config)

	if generator == nil {
		t.Fatal("CSV 生成器建立失敗")
	}

	if generator.config.Format != "csv" {
		t.Errorf("配置設定錯誤: 期望 csv, 得到 %s", generator.config.Format)
	}
}

// TestGenerateCSV 測試 CSV 生成
func TestGenerateCSV(t *testing.T) {
	generator := NewCSVGenerator(types.ReportConfig{})
	report := testutils.CreateTestReport()

	csvData, err := generator.GenerateCSV(report)
	if err != nil {
		t.Fatalf("生成 CSV 失敗: %v", err)
	}

	csvString := string(csvData)

	// 驗證 CSV 內容
	if !strings.Contains(csvString, "Token 使用報告") {
		t.Error("CSV 應包含報告標題")
	}

	if !strings.Contains(csvString, "coding") {
		t.Error("CSV 應包含活動類型")
	}

	// 驗證數據行數
	lines := strings.Split(strings.TrimSpace(csvString), "\n")
	if len(lines) < 5 {
		t.Errorf("CSV 行數不足: 期望至少 5 行, 得到 %d 行", len(lines))
	}
}

