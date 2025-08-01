package cmd

import (
	"fmt"

	"token-monitor/internal/services"

	"github.com/spf13/cobra"
)

// monitorCmd 監控命令
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "開始 Token 使用量監控",
	Long: `開始監控 Kiro IDE 的 Token 使用量。

此命令會：
- 監控 IDE 活動和對話
- 計算 Token 使用量
- 記錄成本資訊
- 提供即時統計`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing services for monitoring...")
		_ = services.GetInstance() // Ensure services are initialized
		fmt.Println("開始 Token 監控...")
		// TODO: 實作監控邏輯
	},
}

func init() {
	// monitorCmd is added to rootCmd in root.go

	// 監控相關的 flags
	monitorCmd.Flags().StringP("log-path", "l", "", "指定日誌檔案路徑")
	monitorCmd.Flags().BoolP("real-time", "r", false, "即時顯示統計")
	monitorCmd.Flags().IntP("interval", "i", 60, "統計更新間隔（秒）")
}
