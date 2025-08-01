package cmd

import (
	"fmt"
	"strconv"
	"token-monitor/internal/services"

	"github.com/spf13/cobra"
)

var costCmd = &cobra.Command{
	Use:   "cost",
	Short: "計算 Token 使用成本",
	Long: `計算 Token 使用成本，支援多種定價模型和計費模式。

範例:
  cost 1000 2000                           # 基本成本計算
  cost 1000 2000 --breakdown               # 顯示成本分解
  cost 1000 2000 --model claude-opus-4.0   # 使用不同模型`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("使用方式: cost <input_tokens> <output_tokens> [flags]")
			return
		}

		inputTokens, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("錯誤: 無效的輸入 token 數量: %s\n", args[0])
			return
		}

		outputTokens, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("錯誤: 無效的輸出 token 數量: %s\n", args[1])
			return
		}

		model, _ := cmd.Flags().GetString("model")
		breakdown, _ := cmd.Flags().GetBool("breakdown")

		calculator := services.GetInstance().CostCalculator
		result, err := calculator.CalculateCost(inputTokens, outputTokens, model)
		if err != nil {
			fmt.Printf("錯誤: %v\n", err)
			return
		}

		fmt.Printf("=== Token 成本計算結果 ===\n")
		fmt.Printf("模型: %s\n", result.PricingModel)
		fmt.Printf("輸入 Tokens: %d\n", inputTokens)
		fmt.Printf("輸出 Tokens: %d\n", outputTokens)
		fmt.Printf("總成本: $%.6f\n", result.TotalCost)

		if breakdown {
			fmt.Printf("\n=== 成本分解 ===\n")
			fmt.Printf("輸入成本: $%.6f\n", result.InputCost)
			fmt.Printf("輸出成本: $%.6f\n", result.OutputCost)
		}
	},
}

func init() {
	// costCmd is added to rootCmd in root.go

	costCmd.Flags().StringP("model", "m", "claude-sonnet-4.0", "定價模型")
	costCmd.Flags().BoolP("breakdown", "b", false, "顯示成本分解")
}
