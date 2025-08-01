package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd 代表基礎命令
var rootCmd = &cobra.Command{
	Use:   "token-monitor",
	Short: "Kiro IDE Token 使用量監控工具",
	Long: `Token Monitor 是一個用於監控和分析 Kiro IDE Token 使用量的工具。

它提供以下功能：
- 精確的 Token 計算（支援 tiktoken）
- 活動類型分析和統計
- Claude Sonnet 4.0 成本計算
- 多格式報告生成
- 歷史趨勢分析`,
}

// Execute 執行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全域 flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.token-monitor.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")
	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")

	// 綁定 flags 到 viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// 添加子命令
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(calculateCmd)
	rootCmd.AddCommand(compareCmd)
	rootCmd.AddCommand(costCmd)
	rootCmd.AddCommand(monitorCmd)
	rootCmd.AddCommand(reportCmd)
}

// initConfig 初始化配置
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置檔案
		viper.SetConfigFile(cfgFile)
	} else {
		// 尋找 home 目錄
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// 搜尋配置檔案的位置
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".token-monitor")
	}

	// 讀取環境變數
	viper.AutomaticEnv()

	// 讀取配置檔案
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// 設定預設值
	setDefaults()
}

// setDefaults 設定預設配置值
func setDefaults() {
	// 定價設定
	viper.SetDefault("pricing.claude-sonnet-4.0.input", 3.0)
	viper.SetDefault("pricing.claude-sonnet-4.0.output", 15.0)
	viper.SetDefault("pricing.claude-sonnet-4.0.cache_read", 0.30)
	viper.SetDefault("pricing.claude-sonnet-4.0.cache_write", 3.75)
	viper.SetDefault("pricing.claude-sonnet-4.0.batch_discount", 0.5)

	// 活動模式設定
	viper.SetDefault("activities.patterns.coding", []string{"function", "class", "implement", "程式", "函數"})
	viper.SetDefault("activities.patterns.debugging", []string{"error", "bug", "fix", "錯誤", "修復"})
	viper.SetDefault("activities.patterns.documentation", []string{"README", "document", "文件", "說明"})
	viper.SetDefault("activities.patterns.spec-development", []string{"spec", "requirement", "design", "需求", "設計"})
	viper.SetDefault("activities.patterns.chat", []string{"chat", "question", "help", "問題", "協助"})

	// 儲存設定
	viper.SetDefault("storage.path", "./data")
	viper.SetDefault("storage.backup_interval", "24h")
	viper.SetDefault("storage.retention_days", 90)

	// Token 計算設定
	viper.SetDefault("token_calculation.preferred_method", "tiktoken")
	viper.SetDefault("token_calculation.fallback_method", "estimation")
	viper.SetDefault("token_calculation.cache_enabled", true)
	viper.SetDefault("token_calculation.cache_size", 1000)
}
