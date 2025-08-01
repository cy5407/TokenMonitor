package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"token-monitor/internal/types"
)

// ChartGenerator 圖表生成器
type ChartGenerator struct {
	config types.ReportConfig
}

// ChartData 圖表資料結構
type ChartData struct {
	Labels   []string    `json:"labels"`
	Datasets []Dataset   `json:"datasets"`
	Options  ChartOptions `json:"options"`
}

// Dataset 資料集
type Dataset struct {
	Label           string    `json:"label"`
	Data            []float64 `json:"data"`
	BackgroundColor []string  `json:"backgroundColor,omitempty"`
	BorderColor     []string  `json:"borderColor,omitempty"`
	BorderWidth     int       `json:"borderWidth,omitempty"`
	Tension         float64   `json:"tension,omitempty"`
}

// ChartOptions 圖表選項
type ChartOptions struct {
	Responsive bool                   `json:"responsive"`
	Plugins    map[string]interface{} `json:"plugins,omitempty"`
	Scales     map[string]interface{} `json:"scales,omitempty"`
}

// NewChartGenerator 建立圖表生成器
func NewChartGenerator(config types.ReportConfig) *ChartGenerator {
	return &ChartGenerator{config: config}
}

// GenerateActivityPieChart 生成活動類型圓餅圖
func (cg *ChartGenerator) GenerateActivityPieChart(report *types.BasicReport) (*ChartData, error) {
	var labels []string
	var data []float64
	colors := []string{
		"#4CAF50", "#2196F3", "#FF9800", "#F44336", "#9C27B0",
		"#00BCD4", "#8BC34A", "#FFC107", "#E91E63", "#607D8B",
	}

	colorIndex := 0
	var backgroundColors []string

	for activityType, activityReport := range report.ByActivity {
		labels = append(labels, string(activityType))
		data = append(data, activityReport.Percentage)

		if colorIndex < len(colors) {
			backgroundColors = append(backgroundColors, colors[colorIndex])
		} else {
			backgroundColors = append(backgroundColors, "#CCCCCC")
		}
		colorIndex++
	}

	chartData := &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "活動分佈",
				Data:            data,
				BackgroundColor: backgroundColors,
				BorderColor:     []string{"#ffffff"},
				BorderWidth:     2,
			},
		},
		Options: ChartOptions{
			Responsive: true,
			Plugins: map[string]interface{}{
				"legend": map[string]interface{}{
					"position": "bottom",
				},
				"title": map[string]interface{}{
					"display": true,
					"text":    "活動類型分佈",
				},
			},
		},
	}

	return chartData, nil
}

// GenerateTokenTrendChart 生成 Token 使用趨勢圖
func (cg *ChartGenerator) GenerateTokenTrendChart(report *types.BasicReport) (*ChartData, error) {
	// 生成24小時標籤
	var labels []string
	var data []float64

	for hour := 0; hour < 24; hour++ {
		labels = append(labels, fmt.Sprintf("%02d:00", hour))

		if count, exists := report.Statistics.ActivityTrends.HourlyDistribution[hour]; exists {
			data = append(data, float64(count))
		} else {
			data = append(data, 0)
		}
	}

	chartData := &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:       "活動數量",
				Data:        data,
				BorderColor: []string{"#4CAF50"},
				BackgroundColor: []string{"rgba(76, 175, 80, 0.1)"},
				BorderWidth: 2,
				Tension:     0.4,
			},
		},
		Options: ChartOptions{
			Responsive: true,
			Scales: map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			Plugins: map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "24小時活動趨勢",
				},
			},
		},
	}

	return chartData, nil
}

// GenerateTokenDistributionChart 生成 Token 分佈圖
func (cg *ChartGenerator) GenerateTokenDistributionChart(report *types.BasicReport) (*ChartData, error) {
	var labels []string
	var inputData []float64
	var outputData []float64

	for activityType, activityReport := range report.ByActivity {
		labels = append(labels, string(activityType))
		inputData = append(inputData, float64(activityReport.Tokens.InputTokens))
		outputData = append(outputData, float64(activityReport.Tokens.OutputTokens))
	}

	chartData := &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "輸入 Token",
				Data:            inputData,
				BackgroundColor: []string{"rgba(54, 162, 235, 0.8)"},
				BorderColor:     []string{"rgba(54, 162, 235, 1)"},
				BorderWidth:     1,
			},
			{
				Label:           "輸出 Token",
				Data:            outputData,
				BackgroundColor: []string{"rgba(255, 99, 132, 0.8)"},
				BorderColor:     []string{"rgba(255, 99, 132, 1)"},
				BorderWidth:     1,
			},
		},
		Options: ChartOptions{
			Responsive: true,
			Scales: map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			Plugins: map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Token 輸入輸出分佈",
				},
			},
		},
	}

	return chartData, nil
}

// GenerateChartJS 生成 Chart.js 腳本
func (cg *ChartGenerator) GenerateChartJS(chartType string, chartData *ChartData, canvasId string) (string, error) {
	dataJSON, err := json.Marshal(chartData)
	if err != nil {
		return "", fmt.Errorf("序列化圖表資料失敗: %w", err)
	}

	script := fmt.Sprintf(`
		const ctx_%s = document.getElementById('%s');
		if (ctx_%s) {
			new Chart(ctx_%s, {
				type: '%s',
				data: %s.datasets ? %s : {
					labels: %s.labels,
					datasets: %s.datasets
				},
				options: %s.options || {}
			});
		}
	`, canvasId, canvasId, canvasId, canvasId, chartType,
		string(dataJSON), string(dataJSON),
		string(dataJSON), string(dataJSON), string(dataJSON))

	return script, nil
}

// GenerateAllCharts 生成所有圖表
func (cg *ChartGenerator) GenerateAllCharts(report *types.BasicReport) (map[string]string, error) {
	charts := make(map[string]string)

	// 活動圓餅圖
	pieData, err := cg.GenerateActivityPieChart(report)
	if err != nil {
		return nil, fmt.Errorf("生成圓餅圖失敗: %w", err)
	}

	pieScript, err := cg.GenerateChartJS("doughnut", pieData, "activityChart")
	if err != nil {
		return nil, fmt.Errorf("生成圓餅圖腳本失敗: %w", err)
	}
	charts["activity_pie"] = pieScript

	// 趨勢線圖
	trendData, err := cg.GenerateTokenTrendChart(report)
	if err != nil {
		return nil, fmt.Errorf("生成趨勢圖失敗: %w", err)
	}

	trendScript, err := cg.GenerateChartJS("line", trendData, "trendChart")
	if err != nil {
		return nil, fmt.Errorf("生成趨勢圖腳本失敗: %w", err)
	}
	charts["token_trend"] = trendScript

	// Token 分佈圖
	distData, err := cg.GenerateTokenDistributionChart(report)
	if err != nil {
		return nil, fmt.Errorf("生成分佈圖失敗: %w", err)
	}

	distScript, err := cg.GenerateChartJS("bar", distData, "distributionChart")
	if err != nil {
		return nil, fmt.Errorf("生成分佈圖腳本失敗: %w", err)
	}
	charts["token_distribution"] = distScript

	return charts, nil
}

// GetChartHTML 獲取圖表 HTML 元素
func (cg *ChartGenerator) GetChartHTML() string {
	return `
		<div class="chart-container">
			<div class="chart-row">
				<div class="chart-item">
					<canvas id="activityChart"></canvas>
				</div>
				<div class="chart-item">
					<canvas id="trendChart"></canvas>
				</div>
			</div>
			<div class="chart-row">
				<div class="chart-item full-width">
					<canvas id="distributionChart"></canvas>
				</div>
			</div>
		</div>
		
		<style>
		.chart-container {
			margin: 20px 0;
		}
		.chart-row {
			display: flex;
			gap: 20px;
			margin-bottom: 20px;
		}
		.chart-item {
			flex: 1;
			background: white;
			padding: 20px;
			border-radius: 8px;
			box-shadow: 0 2px 10px rgba(0,0,0,0.1);
		}
		.chart-item.full-width {
			flex: none;
			width: 100%;
		}
		@media (max-width: 768px) {
			.chart-row {
				flex-direction: column;
			}
		}
		</style>
	`
}

// ExportChartData 匯出圖表資料
func (cg *ChartGenerator) ExportChartData(report *types.BasicReport, format string) ([]byte, error) {
	charts, err := cg.GenerateAllCharts(report)
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(format) {
	case "json":
		return json.MarshalIndent(charts, "", "  ")
	case "js":
		var scripts []string
		for _, script := range charts {
			scripts = append(scripts, script)
		}
		return []byte(strings.Join(scripts, "\n\n")), nil
	default:
		return nil, fmt.Errorf("不支援的匯出格式: %s", format)
	}
}
