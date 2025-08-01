package cmd

import (
	"fmt"
	"time"

	"token-monitor/internal/services"
	"token-monitor/internal/types"

	"github.com/spf13/cobra"
)

var (
	reportFormat        string
	reportTemplate      string
	reportOutput        string
	reportStartTime     string
	reportEndTime       string
	includeTrends       bool
	includeOptimization bool
)

// reportCmd 報告命令
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "生成 Token 使用報告",
	Long: `生成詳細的 Token 使用報告，支援多種格式和模板。

範例:
  token-monitor report --format json --output report.json
  token-monitor report --template detailed --output report.txt
  token-monitor report --start "2024-01-01" --end "2024-01-31"`,
	RunE: runReportCommand,
}

func init() {
	// reportCmd is added to rootCmd in root.go

	// 報告格式選項
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "json",
		"報告格式 (json, txt, md)")

	// 模板選項
	reportCmd.Flags().StringVarP(&reportTemplate, "template", "t", "basic",
		"報告模板 (basic, detailed, summary)")

	// 輸出檔案
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "",
		"輸出檔案路徑")

	// 時間範圍
	reportCmd.Flags().StringVar(&reportStartTime, "start", "",
		"開始時間 (YYYY-MM-DD 或 YYYY-MM-DD HH:MM:SS)")
	reportCmd.Flags().StringVar(&reportEndTime, "end", "",
		"結束時間 (YYYY-MM-DD 或 YYYY-MM-DD HH:MM:SS)")

	// 包含選項
	reportCmd.Flags().BoolVar(&includeTrends, "trends", false,
		"包含趨勢分析")
	reportCmd.Flags().BoolVar(&includeOptimization, "optimization", false,
		"包含優化建議")

	reportCmd.AddCommand(listTemplatesCmd)
	reportCmd.AddCommand(previewReportCmd)
}

// runReportCommand 執行報告命令
func runReportCommand(cmd *cobra.Command, args []string) error {
	storage := services.GetInstance().Storage
	reportGenerator := services.GetInstance().ReportGenerator

	// 載入活動數據
	activities, err := storage.LoadActivityData(types.TimeRange{})
	if err != nil {
		return fmt.Errorf("載入活動數據失敗: %w", err)
	}

	if len(activities) == 0 {
		fmt.Println("沒有找到活動數據")
		return nil
	}

	// 解析時間範圍
	timeRange, err := parseTimeRange(reportStartTime, reportEndTime)
	if err != nil {
		return fmt.Errorf("解析時間範圍失敗: %w", err)
	}

	// 建立報告選項
	options := types.ReportOptions{
		TimeRange:           timeRange,
		IncludeTrends:       includeTrends,
		IncludeOptimization: includeOptimization,
	}

	// 生成報告
	reportData, err := reportGenerator.Generate(activities, options)
	if err != nil {
		return fmt.Errorf("生成報告數據失敗: %w", err)
	}

	if reportOutput != "" {
		// 儲存到檔案
		err = reportGenerator.Save(reportData, reportFormat, reportOutput)
		if err != nil {
			return fmt.Errorf("儲存報告失敗: %w", err)
		}
		fmt.Printf("報告已儲存到: %s\n", reportOutput)
	} else {
		// 輸出到控制台
		output, err := reportGenerator.Render(reportData, reportTemplate)
		if err != nil {
			return fmt.Errorf("渲染報告失敗: %w", err)
		}
		fmt.Println(output)
	}

	return nil
}

// parseTimeRange 解析時間範圍
func parseTimeRange(startStr, endStr string) (types.TimeRange, error) {
	var timeRange types.TimeRange

	if startStr != "" {
		start, err := parseTimeString(startStr)
		if err != nil {
			return timeRange, fmt.Errorf("解析開始時間失敗: %w", err)
		}
		timeRange.Start = start
	}

	if endStr != "" {
		end, err := parseTimeString(endStr)
		if err != nil {
			return timeRange, fmt.Errorf("解析結束時間失敗: %w", err)
		}
		timeRange.End = end
	}

	return timeRange, nil
}

// parseTimeString 解析時間字符串
func parseTimeString(timeStr string) (time.Time, error) {
	// 支援的時間格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("無法解析時間格式: %s", timeStr)
}

// listTemplatesCmd 列出可用模板命令
var listTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "列出可用的報告模板",
	RunE: func(cmd *cobra.Command, args []string) error {
		reportGenerator := services.GetInstance().ReportGenerator
		templates := reportGenerator.GetAvailableTemplates()

		fmt.Println("可用的報告模板:")
		for _, template := range templates {
			fmt.Printf("  - %s\n", template)
		}

		return nil
	},
}

// previewReportCmd 預覽報告命令
var previewReportCmd = &cobra.Command{
	Use:   "preview [template]",
	Short: "預覽報告模板",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := "basic"
		if len(args) > 0 {
			templateName = args[0]
		}

		storage := services.GetInstance().Storage
		reportGenerator := services.GetInstance().ReportGenerator

		// 載入活動數據
		activities, err := storage.LoadActivityData(types.TimeRange{})
		if err != nil {
			return fmt.Errorf("載入活動數據失敗: %w", err)
		}

		// 生成預覽
		options := types.ReportOptions{IncludeTrends: true}
		preview, err := reportGenerator.RenderPreview(activities, templateName, options)
		if err != nil {
			return fmt.Errorf("生成預覽失敗: %w", err)
		}

		fmt.Printf("=== %s 模板預覽 ===\n", templateName)
		fmt.Println(preview)

		return nil
	},
}
