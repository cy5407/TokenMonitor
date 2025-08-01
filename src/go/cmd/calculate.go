package cmd

import (
	"fmt"
	"os"

	"token-monitor/internal/calculator"
	"token-monitor/internal/services"

	"github.com/spf13/cobra"
)

// calculateCmd 計算命令
var calculateCmd = &cobra.Command{
	Use:   "calculate [text]",
	Short: "計算文本的 Token 數量",
	Long: `計算指定文本的 Token 數量。

範例：
  token-monitor calculate "Hello world"
  token-monitor calculate "你好世界"
  token-monitor calculate "Hello 世界" --method estimation
  echo "長文本內容" | token-monitor calculate --stdin`,
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

		// 取得參數
		method, _ := cmd.Flags().GetString("method")
		showDetails, _ := cmd.Flags().GetBool("details")
		showDistribution, _ := cmd.Flags().GetBool("distribution")

		// 從服務容器取得計算器
		calc := services.GetInstance().TokenCalculator

		// 驗證文本
		if calcImpl, ok := calc.(*calculator.TokenCalculatorImpl); ok {
			if err := calcImpl.ValidateText(text); err != nil {
				fmt.Fprintf(os.Stderr, "Text validation failed: %v\n", err)
				os.Exit(1)
			}
		}

		// 計算 Token
		tokens, err := calc.CalculateTokens(text, method)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Token calculation failed: %v\n", err)
			os.Exit(1)
		}

		// 顯示基本結果
		fmt.Printf("Token 數量: %d\n", tokens)

		if showDetails {
			fmt.Printf("文本長度: %d 字符\n", len(text))
			fmt.Printf("計算方法: %s\n", method)

			// 顯示支援的方法
			methods := calc.GetSupportedMethods()
			fmt.Printf("支援的方法: %v\n", methods)

			// 顯示快取統計（如果可用）
			if calcImpl, ok := calc.(*calculator.TokenCalculatorImpl); ok {
				stats := calcImpl.GetCacheStats()
				fmt.Printf("快取統計: %v\n", stats)

				// 顯示 tiktoken 資訊
				tiktokenInfo := calcImpl.GetTiktokenInfo()
				fmt.Printf("Tiktoken 資訊: %v\n", tiktokenInfo)
			}
		}

		if showDistribution {
			distribution, err := calc.AnalyzeTokenDistribution(text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Token distribution analysis failed: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("\nToken 分佈:\n")
			fmt.Printf("  英文 Token: %d\n", distribution.EnglishTokens)
			fmt.Printf("  中文 Token: %d\n", distribution.ChineseTokens)
			fmt.Printf("  總計 Token: %d\n", distribution.TotalTokens)
			fmt.Printf("  計算方法: %s\n", distribution.Method)
		}
	},
}

func init() {
	rootCmd.AddCommand(calculateCmd)

	// 計算相關的 flags
	calculateCmd.Flags().StringP("method", "m", "", "計算方法 (estimation, tiktoken, auto)")
	calculateCmd.Flags().BoolP("details", "d", false, "顯示詳細資訊")
	calculateCmd.Flags().BoolP("distribution", "t", false, "顯示 Token 分佈")
	calculateCmd.Flags().BoolP("stdin", "i", false, "從標準輸入讀取文本")
}
