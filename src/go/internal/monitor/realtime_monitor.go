package monitor

import (
	"fmt"
	"strings"
	"time"
)

// RealTimeMonitor 即時監控器
type RealTimeMonitor struct {
	isRunning bool
	interval  time.Duration
	stats     *MonitorStats
}

// MonitorStats 監控統計
type MonitorStats struct {
	StartTime       time.Time `json:"start_time"`
	TotalActivities int       `json:"total_activities"`
	TotalTokens     int       `json:"total_tokens"`
	LastActivity    time.Time `json:"last_activity"`
	AverageRate     float64   `json:"average_rate"`
}

// NewRealTimeMonitor 建立即時監控器
func NewRealTimeMonitor(interval time.Duration) *RealTimeMonitor {
	return &RealTimeMonitor{
		interval: interval,
		stats: &MonitorStats{
			StartTime: time.Now(),
		},
	}
}

// Start 開始監控
func (rtm *RealTimeMonitor) Start() {
	rtm.isRunning = true
	rtm.stats.StartTime = time.Now()

	fmt.Println("🔄 開始即時監控...")
	fmt.Println("按 'q' + Enter 停止監控")

	go rtm.monitorLoop()

	// 等待使用者輸入停止命令
	for {
		var input string
		fmt.Scanln(&input)
		if input == "q" || input == "Q" {
			rtm.Stop()
			break
		}
	}
}

// monitorLoop 監控循環
func (rtm *RealTimeMonitor) monitorLoop() {
	ticker := time.NewTicker(rtm.interval)
	defer ticker.Stop()

	for rtm.isRunning {
		select {
		case <-ticker.C:
			rtm.updateStats()
			rtm.displayStats()
		}
	}
}

// updateStats 更新統計資訊
func (rtm *RealTimeMonitor) updateStats() {
	// TODO: 從實際資料源更新統計資訊
	// 這裡使用模擬資料
	rtm.stats.TotalActivities++
	rtm.stats.TotalTokens += 100 + (rtm.stats.TotalActivities % 200)
	rtm.stats.LastActivity = time.Now()

	elapsed := time.Since(rtm.stats.StartTime).Seconds()
	if elapsed > 0 {
		rtm.stats.AverageRate = float64(rtm.stats.TotalActivities) / elapsed
	}
}

// displayStats 顯示統計資訊
func (rtm *RealTimeMonitor) displayStats() {
	// 清除螢幕 (簡單版本)
	fmt.Print("\033[2J\033[H")

	fmt.Println("🔄 即時監控 - " + time.Now().Format("15:04:05"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("⏱️  運行時間: %v\n", time.Since(rtm.stats.StartTime).Round(time.Second))
	fmt.Printf("📊 總活動數: %d\n", rtm.stats.TotalActivities)
	fmt.Printf("🎯 總Token數: %d\n", rtm.stats.TotalTokens)
	fmt.Printf("📈 平均速率: %.2f 活動/秒\n", rtm.stats.AverageRate)
	if !rtm.stats.LastActivity.IsZero() {
		fmt.Printf("🕐 最後活動: %s\n", rtm.stats.LastActivity.Format("15:04:05"))
	}
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("按 'q' + Enter 停止監控")
}

// Stop 停止監控
func (rtm *RealTimeMonitor) Stop() {
	rtm.isRunning = false
	fmt.Println("\n✅ 監控已停止")
}
