package reporter

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"token-monitor/internal/types"
)

// HTMLGenerator HTML 格式報告生成器
type HTMLGenerator struct {
	config   types.ReportConfig
	template *template.Template
}

// NewHTMLGenerator 建立新的 HTML 生成器
func NewHTMLGenerator(config types.ReportConfig) *HTMLGenerator {
	return &HTMLGenerator{
		config:   config,
		template: template.Must(template.New("report").Parse(htmlTemplate)),
	}
}

// GenerateHTML 生成 HTML 格式報告
func (hg *HTMLGenerator) GenerateHTML(report *types.BasicReport) ([]byte, error) {
	// TODO: 實作 HTML 生成邏輯
	var output strings.Builder

	err := hg.template.Execute(&output, report)
	if err != nil {
		return nil, fmt.Errorf("HTML 模板執行失敗: %w", err)
	}

	return []byte(output.String()), nil
}

// SaveHTML 儲存 HTML 報告到檔案
func (hg *HTMLGenerator) SaveHTML(report *types.BasicReport, outputPath string) error {
	htmlData, err := hg.GenerateHTML(report)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, htmlData, 0644)
}

// 增強的 HTML 模板
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Token 使用報告</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; padding: 20px; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container { 
            max-width: 1200px; margin: 0 auto; 
            background: white; border-radius: 10px; 
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #4CAF50, #45a049); 
            color: white; padding: 30px; text-align: center;
        }
        .header h1 { margin: 0; font-size: 2.5em; font-weight: 300; }
        .header .subtitle { opacity: 0.9; margin-top: 10px; }
        
        .summary { 
            display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px; padding: 30px; background: #f8f9fa;
        }
        .summary-card { 
            background: white; padding: 20px; border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: center;
        }
        .summary-card h3 { margin: 0 0 10px 0; color: #666; font-size: 0.9em; text-transform: uppercase; }
        .summary-card .value { font-size: 2em; font-weight: bold; color: #4CAF50; }
        
        .content { padding: 30px; }
        .section { margin-bottom: 40px; }
        .section h2 { 
            color: #333; border-bottom: 3px solid #4CAF50; 
            padding-bottom: 10px; margin-bottom: 20px;
        }
        
        table { 
            width: 100%; border-collapse: collapse; 
            background: white; border-radius: 8px; overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        th { 
            background: linear-gradient(135deg, #4CAF50, #45a049); 
            color: white; padding: 15px; text-align: left; font-weight: 500;
        }
        td { padding: 12px 15px; border-bottom: 1px solid #eee; }
        tr:hover { background: #f5f5f5; }
        
        .progress-bar {
            width: 100%; height: 20px; background: #eee; border-radius: 10px; overflow: hidden;
        }
        .progress-fill {
            height: 100%; background: linear-gradient(90deg, #4CAF50, #45a049);
            transition: width 0.3s ease;
        }
        
        .chart-container { 
            margin: 20px 0; padding: 20px; 
            background: white; border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        
        @media (max-width: 768px) {
            .summary { grid-template-columns: 1fr; }
            .container { margin: 10px; }
            body { padding: 10px; }
        }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script>
        document.addEventListener('DOMContentLoaded', function() {
            // 活動類型分佈圓餅圖
            const activityData = {
                labels: [{{range $type, $report := .ByActivity}}'{{$type}}',{{end}}],
                datasets: [{
                    data: [{{range $type, $report := .ByActivity}}{{$report.Percentage}},{{end}}],
                    backgroundColor: ['#4CAF50', '#2196F3', '#FF9800', '#F44336', '#9C27B0'],
                    borderWidth: 2,
                    borderColor: '#fff'
                }]
            };
            
            const ctx = document.getElementById('activityChart');
            if (ctx) {
                new Chart(ctx, {
                    type: 'doughnut',
                    data: activityData,
                    options: {
                        responsive: true,
                        plugins: {
                            legend: { position: 'bottom' },
                            title: { display: true, text: '活動類型分佈' }
                        }
                    }
                });
            }
            
            // Token 使用趨勢圖
            const hourlyData = {
                labels: Array.from({length: 24}, (_, i) => i + ':00'),
                datasets: [{
                    label: '活動數量',
                    data: [{{range $hour := .Statistics.ActivityTrends.HourlyDistribution}}{{$hour}},{{end}}],
                    borderColor: '#4CAF50',
                    backgroundColor: 'rgba(76, 175, 80, 0.1)',
                    tension: 0.4
                }]
            };
            
            const trendCtx = document.getElementById('trendChart');
            if (trendCtx) {
                new Chart(trendCtx, {
                    type: 'line',
                    data: hourlyData,
                    options: {
                        responsive: true,
                        scales: {
                            y: { beginAtZero: true }
                        },
                        plugins: {
                            title: { display: true, text: '24小時活動趨勢' }
                        }
                    }
                });
            }
        });
    </script>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Token 使用分析報告</h1>
            <div class="subtitle">生成時間: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</div>
        </div>
        
        <div class="summary">
            <div class="summary-card">
                <h3>總活動數</h3>
                <div class="value">{{.Summary.TotalActivities}}</div>
            </div>
            <div class="summary-card">
                <h3>總 Token 數</h3>
                <div class="value">{{.Summary.TotalTokens.TotalTokens}}</div>
            </div>
            <div class="summary-card">
                <h3>平均 Token</h3>
                <div class="value">{{printf "%.1f" .Summary.AverageTokensPerActivity}}</div>
            </div>
            <div class="summary-card">
                <h3>記錄數</h3>
                <div class="value">{{.TotalRecords}}</div>
            </div>
        </div>
        
        <div class="content">
            <div class="section">
                <h2>📊 活動類型分析</h2>
                <table>
                    <thead>
                        <tr>
                            <th>活動類型</th>
                            <th>數量</th>
                            <th>Token 總數</th>
                            <th>平均 Token</th>
                            <th>佔比</th>
                            <th>分佈</th>
                        </tr>
                    </thead>
                    <tbody>
                        {{range $type, $report := .ByActivity}}
                        <tr>
                            <td><strong>{{$type}}</strong></td>
                            <td>{{$report.Count}}</td>
                            <td>{{$report.Tokens.TotalTokens}}</td>
                            <td>{{printf "%.2f" $report.AverageTokens}}</td>
                            <td>{{printf "%.1f" $report.Percentage}}%</td>
                            <td>
                                <div class="progress-bar">
                                    <div class="progress-fill" style="width: {{$report.Percentage}}%"></div>
                                </div>
                            </td>
                        </tr>
                        {{end}}
                    </tbody>
                </table>
            </div>
            
            <div class="section">
                <h2>📈 統計摘要</h2>
                <div class="chart-container">
                    <table>
                        <tr><td><strong>Token 總數</strong></td><td>{{.Statistics.TokenDistribution.Total}}</td></tr>
                        <tr><td><strong>平均值</strong></td><td>{{printf "%.2f" .Statistics.TokenDistribution.Average}}</td></tr>
                        <tr><td><strong>最小值</strong></td><td>{{.Statistics.TokenDistribution.Min}}</td></tr>
                        <tr><td><strong>最大值</strong></td><td>{{.Statistics.TokenDistribution.Max}}</td></tr>
                        <tr><td><strong>中位數</strong></td><td>{{printf "%.2f" .Statistics.TokenDistribution.Median}}</td></tr>
                        <tr><td><strong>峰值時間</strong></td><td>{{printf "%02d:00" .Statistics.ActivityTrends.PeakHour}} ({{.Statistics.ActivityTrends.PeakHourCount}} 個活動)</td></tr>
                    </table>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`
