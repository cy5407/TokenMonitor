package cmd

import (
	"fmt"

	"token-monitor/internal/services"
	"token-monitor/internal/types"

	"github.com/spf13/cobra"
)

// analyzeCmd 分析命令
var analyzeCmd = &cobra.Command{
	Use:	"analyze",
	Short: "分析 Token 使用模式",
	Long: `分析 Token 使用模式和效率。

提供以下分析：
- 活動類型分析
- 效率指標計算
- 使用趨勢分析
- 成本優化建議`,
	Run: func(cmd *cobra.Command, args []string) {
		activityType, _ := cmd.Flags().GetString("activity-type")

		analyzer := services.GetInstance().ActivityAnalyzer
		storage := services.GetInstance().Storage

		activities, err := storage.LoadActivityData(types.TimeRange{})
		if err != nil {
			fmt.Printf("Error loading activities: %v\n", err)
			return
		}

		if activityType != "" {
			fmt.Printf("分析 %s 活動的 Token 使用情況\n", activityType)
			// TODO: Filter activities by type
		}

		summary := analyzer.GenerateActivitySummary(activities)

		fmt.Printf("分析完成，共處理 %d 個活動。\n", summary.TotalActivities)
		fmt.Printf("總 Token 用量: %d\n", summary.TotalTokens.TotalTokens)
		fmt.Printf("活動分佈:\n")
		for actType, count := range summary.ActivityCounts {
			fmt.Printf("  - %s: %d\n", actType, count)
		}
	},
}

func init() {
	// analyzeCmd is added to rootCmd in root.go

	// 分析相關的 flags
	analyzeCmd.Flags().StringP("activity-type", "a", "", "指定活動類型 (coding, debugging, documentation, spec-development, chat)")
	analyzeCmd.Flags().StringP("period", "p", "7d", "分析期間 (1d, 7d, 30d)")
	analyzeCmd.Flags().BoolP("efficiency", "e", false, "顯示效率分析")
	analyzeCmd.Flags().BoolP("trends", "t", false, "顯示趨勢分析")
}
