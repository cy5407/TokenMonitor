package utils

import (
	"fmt"
	"strings"
)

// ProgressBar is a simple progress bar
type ProgressBar struct {
	total   int
	current int
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{total: total, current: 0}
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int) {
	p.current = current
	percentage := float64(p.current) / float64(p.total) * 100
	bar := strings.Repeat("=", int(percentage/2)) + ">"
	fmt.Printf("\r[%-50s] %d/%d (%.2f%%)", bar, p.current, p.total, percentage)
}

