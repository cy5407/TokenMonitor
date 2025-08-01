package cmd

import (
	"fmt"
	"os"


	"token-monitor/internal/calculator"
	"token-monitor/internal/services"

	"github.com/spf13/cobra"
)

// compareCmd 比較命令
var compareCmd = &cobra.Command{
	Use:	"compare [text]",
	Short:	"比較不同 Token 計算方法的結果",
	Long: `比較 tiktoken 和估算方法的 Token 計算結果。

範例：
  token-monitor compare "Hello world"
  token-monitor compare "你好世界"
  echo "長文本內容" | token-monitor compare --stdin`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var text string
		var err error

		// 取得輸入文本
		useStdin, _ := cmd.Flags().GetBool("stdin")
		if useStdin {
			// 從標準輸入讀取
			var input []byte
			input, err = os.ReadFile("/dev/stdin")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
				os.Exit(1)
			}
			text = string(input)
		} else if len(args) > 0 {
			text = args[0]
		} else {
			fmt.Fprintf(os.Stderr, "Error: Please provide text as argument or use --stdin flag\n")
			os.Exit(1)
		}

		// 從服務容器取得計算器
		calc := services.GetInstance().TokenCalculator

		// 驗證文本
		if calcImpl, ok := calc.(*calculator.TokenCalculatorImpl); ok {
			if err := calcImpl.ValidateText(text); err != nil {
				fmt.Fprintf(os.Stderr, "Text validation failed: %v\n", err)
				os.Exit(1)
			}

			// 比較計算方法
			comparison := calcImpl.CompareCalculationMethods(text)

			fmt.Printf("📊 Token 計算方法比較\n")
			fmt.Printf("====================\n")
			fmt.Printf("文本長度: %d 字符\n\n", comparison["text_length"])

			// 顯示估算結果
			if estimation, ok := comparison["estimation"].(map[string]interface{}); ok {
				fmt.Printf("🔢 估算方法:\n")
				fmt.Printf("  Token 數量: %d\n", estimation["tokens"])
				fmt.Printf("  計算方法: %s\n\n", estimation["method"])
			}

			// 顯示 tiktoken 結果
			if tiktoken, ok := comparison["tiktoken"].(map[string]interface{}); ok {
				fmt.Printf("🎯 Tiktoken 方法:\n")
				fmt.Printf("  Token 數量: %d\n", tiktoken["tokens"])
				fmt.Printf("  計算方法: %s\n\n", tiktoken["method"])
			} else {
				fmt.Printf("⚠️  Tiktoken 不可用，使用估算方法\n\n")
			}

			// 顯示比較結果
			if comp, ok := comparison["comparison"].(map[string]interface{}); ok {
				fmt.Printf("📈 比較結果:\n")
				fmt.Printf("  差異: %d tokens\n", comp["difference"])
				fmt.Printf("  準確度: %.2f%%\n", comp["accuracy_percent"])
				fmt.Printf("  建議方法: %s\n", comp["preferred_method"])
			}

		} else {
			fmt.Fprintf(os.Stderr, "Error: Calculator type assertion failed\n")
			os.Exit(1)
		}
	},
}

func init() {
	// compareCmd is added to rootCmd in root.go

	// 比較相關的 flags
	compareCmd.Flags().BoolP("stdin", "i", false, "從標準輸入讀取文本")
}
