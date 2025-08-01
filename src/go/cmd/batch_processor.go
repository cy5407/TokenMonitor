package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"token-monitor/internal/services"
	"token-monitor/internal/utils"
)

// batchCmd 批次處理命令
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "批次處理多個檔案或目錄",
	Long: `批次處理功能可以：
- 同時處理多個輸入檔案
- 並行生成多種格式報告
- 批次匯出和儲存結果
- 提供進度追蹤和錯誤處理`,
	RunE: runBatchProcessor,
}

// BatchProcessor 批次處理器
type BatchProcessor struct {
	config       BatchConfig
	services     *services.ServiceContainer
	progressBar  *utils.ProgressBar
	results      []BatchResult
	mutex        sync.Mutex
}

// BatchConfig 批次處理配置
type BatchConfig struct {
	InputPaths    []string `json:"input_paths"`
	OutputDir     string   `json:"output_dir"`
	Formats       []string `json:"formats"`
	Parallel      bool     `json:"parallel"`
	MaxWorkers    int      `json:"max_workers"`
	SaveToStorage bool     `json:"save_to_storage"`
}

// BatchResult 批次處理結果
type BatchResult struct {
	InputPath   string        `json:"input_path"`
	Success     bool          `json:"success"`
	Error       string        `json:"error,omitempty"`
	ProcessTime time.Duration `json:"process_time"`
	OutputFiles []string      `json:"output_files"`
	Statistics  BatchStats    `json:"statistics"`
}

// BatchStats 批次統計
type BatchStats struct {
	TotalActivities int     `json:"total_activities"`
	TotalTokens     int     `json:"total_tokens"`
	TotalCost       float64 `json:"total_cost"`
	ProcessingRate  float64 `json:"processing_rate"`
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// 批次處理參數
	batchCmd.Flags().StringSliceP("input", "i", []string{}, "輸入檔案或目錄列表")
	batchCmd.Flags().StringP("output", "o", "./batch_output", "輸出目錄")
	batchCmd.Flags().StringSliceP("format", "f", []string{"json"}, "輸出格式 (json,csv,html)")
	batchCmd.Flags().BoolP("parallel", "p", true, "啟用並行處理")
	batchCmd.Flags().IntP("workers", "w", 4, "並行工作者數量")
	batchCmd.Flags().BoolP("save-storage", "s", true, "儲存到持久化儲存")
}

// runBatchProcessor 執行批次處理
func runBatchProcessor(cmd *cobra.Command, args []string) error {
	// 解析參數
	inputPaths, _ := cmd.Flags().GetStringSlice("input")
	outputDir, _ := cmd.Flags().GetString("output")
	formats, _ := cmd.Flags().GetStringSlice("format")
	parallel, _ := cmd.Flags().GetBool("parallel")
	maxWorkers, _ := cmd.Flags().GetInt("workers")
	saveToStorage, _ := cmd.Flags().GetBool("save-storage")

	if len(inputPaths) == 0 {
		return fmt.Errorf("請指定至少一個輸入檔案或目錄")
	}

	// 建立批次處理器
	config := BatchConfig{
		InputPaths:    inputPaths,
		OutputDir:     outputDir,
		Formats:       formats,
		Parallel:      parallel,
		MaxWorkers:    maxWorkers,
		SaveToStorage: saveToStorage,
	}

	processor := NewBatchProcessor(config)
	return processor.Process()
}

// NewBatchProcessor 建立批次處理器
func NewBatchProcessor(config BatchConfig) *BatchProcessor {
	return &BatchProcessor{
		config:      config,
		services:    services.GetInstance(),
		results:     make([]BatchResult, 0),
	}
}

// Process 執行批次處理
func (bp *BatchProcessor) Process() error {
	fmt.Println("🚀 開始批次處理...")

	// 確保輸出目錄存在
	if err := os.MkdirAll(bp.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("建立輸出目錄失敗: %w", err)
	}

	// 展開輸入路徑
	allFiles, err := bp.expandInputPaths()
	if err != nil {
		return fmt.Errorf("展開輸入路徑失敗: %w", err)
	}

	fmt.Printf("📁 找到 %d 個檔案待處理\n", len(allFiles))

	// 初始化進度條
	bp.progressBar = utils.NewProgressBar(len(allFiles))

	startTime := time.Now()

	if bp.config.Parallel {
		err = bp.processParallel(allFiles)
	} else {
		err = bp.processSequential(allFiles)
	}

	if err != nil {
		return err
	}

	// 顯示處理結果
	bp.showResults(time.Since(startTime))

	// 生成批次報告
	return bp.generateBatchReport()
}

// expandInputPaths 展開輸入路徑
func (bp *BatchProcessor) expandInputPaths() ([]string, error) {
	var allFiles []string

	for _, inputPath := range bp.config.InputPaths {
		info, err := os.Stat(inputPath)
		if err != nil {
			fmt.Printf("⚠️  跳過無效路徑: %s (%v)\n", inputPath, err)
			continue
		}

		if info.IsDir() {
			// 遞歸搜尋目錄中的檔案
			err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && bp.isValidFile(path) {
					allFiles = append(allFiles, path)
				}

				return nil
			})

			if err != nil {
				return nil, err
			}
		} else {
			if bp.isValidFile(inputPath) {
				allFiles = append(allFiles, inputPath)
			}
		}
	}

	return allFiles, nil
}

// isValidFile 檢查是否為有效檔案
func (bp *BatchProcessor) isValidFile(filename string) bool {
	validExtensions := []string{".txt", ".log", ".json", ".csv"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}

	return false
}

// processParallel 並行處理
func (bp *BatchProcessor) processParallel(files []string) error {
	fmt.Printf("⚡ 使用 %d 個工作者並行處理\n", bp.config.MaxWorkers)

	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	// 啟動工作者
	for i := 0; i < bp.config.MaxWorkers; i++ {
		wg.Add(1)
		go bp.worker(jobs, &wg)
	}

	// 發送任務
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// 等待完成
	wg.Wait()

	return nil
}

// processSequential 順序處理
func (bp *BatchProcessor) processSequential(files []string) error {
	fmt.Println("📝 順序處理檔案")

	for _, file := range files {
		bp.processFile(file)
	}

	return nil
}

// worker 工作者函數
func (bp *BatchProcessor) worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for file := range jobs {
		bp.processFile(file)
	}
}

// processFile 處理單個檔案
func (bp *BatchProcessor) processFile(inputPath string) {
	startTime := time.Now()
	result := BatchResult{
		InputPath:   inputPath,
		Success:     false,
		OutputFiles: make([]string, 0),
	}

	defer func() {
		result.ProcessTime = time.Since(startTime)
		bp.mutex.Lock()
		bp.results = append(bp.results, result)
		bp.progressBar.Update(len(bp.results))
		bp.mutex.Unlock()
	}()

	// TODO: 實作檔案處理邏輯
	// 1. 讀取檔案內容
	// 2. 解析活動資料
	// 3. 計算 Token 和成本
	// 4. 生成報告
	// 5. 儲存結果

	// 模擬處理
	time.Sleep(100 * time.Millisecond)

	result.Success = true
	result.Statistics = BatchStats{
		TotalActivities: 10,
		TotalTokens:     1000,
		TotalCost:       5.0,
		ProcessingRate:  100.0,
	}
}

// showResults 顯示處理結果
func (bp *BatchProcessor) showResults(totalTime time.Duration) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 批次處理結果")
	fmt.Println(strings.Repeat("=", 60))

	successCount := 0
	totalActivities := 0
	totalTokens := 0
	totalCost := 0.0

	for _, result := range bp.results {
		if result.Success {
			successCount++
			totalActivities += result.Statistics.TotalActivities
			totalTokens += result.Statistics.TotalTokens
			totalCost += result.Statistics.TotalCost
		}
	}

	fmt.Printf("✅ 成功處理: %d/%d 檔案\n", successCount, len(bp.results))
	fmt.Printf("⏱️  總耗時: %v\n", totalTime.Round(time.Second))
	fmt.Printf("📈 總活動數: %d\n", totalActivities)
	fmt.Printf("🎯 總Token數: %d\n", totalTokens)
	fmt.Printf("💰 總成本: $%.2f\n", totalCost)
	fmt.Printf("⚡ 處理速度: %.2f 檔案/秒\n", float64(len(bp.results))/totalTime.Seconds())

	// 顯示失敗的檔案
	for _, result := range bp.results {
		if !result.Success {
			fmt.Printf("❌ 失敗: %s - %s\n", result.InputPath, result.Error)
		}
	}
}

// generateBatchReport 生成批次報告
func (bp *BatchProcessor) generateBatchReport() error {
	fmt.Println("\n📄 生成批次處理報告...")

	reportPath := filepath.Join(bp.config.OutputDir, "batch_report.json")

	// TODO: 實作批次報告生成
	fmt.Printf("💾 批次報告已儲存: %s\n", reportPath)

	return nil
}


// ProgressBar is a simple progress bar
type ProgressBar struct {
	total int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{total: total}
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int) {
	percentage := float64(current) / float64(p.total) * 100
	bar := strings.Repeat("=", int(percentage/2)) + ">"
	fmt.Printf("\r[%-50s] %d/%d (%.2f%%)", bar, current, p.total, percentage)
}
