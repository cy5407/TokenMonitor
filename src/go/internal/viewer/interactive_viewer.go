package viewer

import (
	"fmt"
	"strings"
	"time"

	"token-monitor/internal/types"
)

// InteractiveReportViewer 互動式報告檢視器
type InteractiveReportViewer struct {
	reports []types.BasicReport
	current int
}

// NewInteractiveReportViewer 建立互動式報告檢視器
func NewInteractiveReportViewer(reports []types.BasicReport) *InteractiveReportViewer {
	return &InteractiveReportViewer{
		reports: reports,
		current: 0,
	}
}

// ShowReportMenu 顯示報告選單
func (irv *InteractiveReportViewer) ShowReportMenu() {
	for {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Printf("📊 報告檢視器 (%d/%d)\n", irv.current+1, len(irv.reports))
		fmt.Println(strings.Repeat("=", 60))

		if len(irv.reports) > 0 {
			report := irv.reports[irv.current]
			fmt.Printf("📅 生成時間: %s\n", report.GeneratedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("📈 總活動數: %d\n", report.Summary.TotalActivities)
			fmt.Printf("🎯 總Token數: %d\n", report.Summary.TotalTokens.TotalTokens)
			fmt.Printf("📊 平均Token: %.2f\n", report.Summary.AverageTokensPerActivity)
		}

		fmt.Println(strings.Repeat("-", 60))
		fmt.Println("1. 查看詳細資訊")
		fmt.Println("2. 查看統計圖表")
		fmt.Println("3. 匯出報告")
		fmt.Println("4. 上一個報告")
		fmt.Println("5. 下一個報告")
		fmt.Println("6. 返回主選單")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Print("請選擇操作 (1-6): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			irv.showDetailedInfo()
		case "2":
			irv.showStatisticsChart()
		case "3":
			irv.exportReport()
		case "4":
			irv.previousReport()
		case "5":
			irv.nextReport()
		case "6":
			return
		default:
			fmt.Println("❌ 無效選項，請重新選擇")
		}
	}
}

// showDetailedInfo 顯示詳細資訊
func (irv *InteractiveReportViewer) showDetailedInfo() {
	if len(irv.reports) == 0 {
		fmt.Println("❌ 沒有可用的報告")
		return
	}

	report := irv.reports[irv.current]
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📋 詳細資訊")
	fmt.Println(strings.Repeat("=", 50))

	for activityType, activityReport := range report.ByActivity {
		fmt.Printf("\n🔸 %s:\n", activityType)
		fmt.Printf("   數量: %d\n", activityReport.Count)
		fmt.Printf("   Token總數: %d\n", activityReport.Tokens.TotalTokens)
		fmt.Printf("   平均Token: %.2f\n", activityReport.AverageTokens)
		fmt.Printf("   佔比: %.2f%%\n", activityReport.Percentage)
	}

	fmt.Println("\n按 Enter 繼續...")
	fmt.Scanln()
}

// showStatisticsChart 顯示統計圖表
func (irv *InteractiveReportViewer) showStatisticsChart() {
	if len(irv.reports) == 0 {
		fmt.Println("❌ 沒有可用的報告")
		return
	}

	report := irv.reports[irv.current]
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 統計圖表")
	fmt.Println(strings.Repeat("=", 50))

	// ASCII 圖表顯示
	fmt.Println("\n活動類型分佈:")
	for activityType, activityReport := range report.ByActivity {
		barLength := int(activityReport.Percentage / 2) // 縮放到50個字符內
		bar := strings.Repeat("█", barLength)
		fmt.Printf("%-15s |%s %.1f%%\n", activityType, bar, activityReport.Percentage)
	}

	fmt.Printf("\n📈 Token 統計:\n")
	fmt.Printf("   總數: %d\n", report.Statistics.TokenDistribution.Total)
	fmt.Printf("   平均: %.2f\n", report.Statistics.TokenDistribution.Average)
	fmt.Printf("   最小: %d\n", report.Statistics.TokenDistribution.Min)
	fmt.Printf("   最大: %d\n", report.Statistics.TokenDistribution.Max)
	fmt.Printf("   中位數: %.2f\n", report.Statistics.TokenDistribution.Median)

	fmt.Println("\n按 Enter 繼續...")
	fmt.Scanln()
}

// exportReport 匯出報告
func (irv *InteractiveReportViewer) exportReport() {
	if len(irv.reports) == 0 {
		fmt.Println("❌ 沒有可用的報告")
		return
	}

	fmt.Println("\n📤 匯出報告")
	fmt.Println("1. CSV 格式")
	fmt.Println("2. HTML 格式")
	fmt.Println("3. JSON 格式")
	fmt.Print("請選擇格式 (1-3): ")

	var choice string
	fmt.Scanln(&choice)

	timestamp := time.Now().Format("20060102_150405")

	switch choice {
	case "1":
		filename := fmt.Sprintf("report_%s.csv", timestamp)
		fmt.Printf("💾 匯出 CSV 到: %s\n", filename)
		// TODO: 實際匯出邏輯
	case "2":
		filename := fmt.Sprintf("report_%s.html", timestamp)
		fmt.Printf("💾 匯出 HTML 到: %s\n", filename)
		// TODO: 實際匯出邏輯
	case "3":
		filename := fmt.Sprintf("report_%s.json", timestamp)
		fmt.Printf("💾 匯出 JSON 到: %s\n", filename)
		// TODO: 實際匯出邏輯
	default:
		fmt.Println("❌ 無效格式")
	}
}

// previousReport 上一個報告
func (irv *InteractiveReportViewer) previousReport() {
	if len(irv.reports) == 0 {
		return
	}

	irv.current--
	if irv.current < 0 {
		irv.current = len(irv.reports) - 1
	}
	fmt.Printf("📋 切換到報告 %d/%d\n", irv.current+1, len(irv.reports))
}

// nextReport 下一個報告
func (irv *InteractiveReportViewer) nextReport() {
	if len(irv.reports) == 0 {
		return
	}

	irv.current++
	if irv.current >= len(irv.reports) {
		irv.current = 0
	}
	fmt.Printf("📋 切換到報告 %d/%d\n", irv.current+1, len(irv.reports))
}
