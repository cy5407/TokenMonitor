package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"token-monitor/internal/types"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	configPath   string
	config       *Config
	watchers     []ConfigWatcher
	autoSave     bool
	lastModified time.Time
}

// Config 主配置結構
type Config struct {
	Version     string           `json:"version"`
	LastUpdated time.Time        `json:"last_updated"`
	General     GeneralConfig    `json:"general"`
	Calculator  CalculatorConfig `json:"calculator"`
	Analyzer    AnalyzerConfig   `json:"analyzer"`
	Cost        CostConfig       `json:"cost"`
	Reporter    ReporterConfig   `json:"reporter"`
	Storage     StorageConfig    `json:"storage"`
	CLI         CLIConfig        `json:"cli"`
}

// GeneralConfig 一般配置
type GeneralConfig struct {
	Language        string `json:"language"`
	Timezone        string `json:"timezone"`
	LogLevel        string `json:"log_level"`
	EnableDebug     bool   `json:"enable_debug"`
	MaxConcurrency  int    `json:"max_concurrency"`
	CacheEnabled    bool   `json:"cache_enabled"`
	CacheExpiration string `json:"cache_expiration"`
}

// CalculatorConfig Token 計算器配置
type CalculatorConfig struct {
	DefaultMethod   string             `json:"default_method"`
	TiktokenEnabled bool               `json:"tiktoken_enabled"`
	EstimationRules map[string]float64 `json:"estimation_rules"`
	CacheResults    bool               `json:"cache_results"`
	MaxTokenLength  int                `json:"max_token_length"`
}

// AnalyzerConfig 分析器配置
type AnalyzerConfig struct {
	EnablePatternMatching bool                           `json:"enable_pattern_matching"`
	CustomPatterns        map[string][]string            `json:"custom_patterns"`
	ActivityWeights       map[types.ActivityType]float64 `json:"activity_weights"`
	MinConfidenceScore    float64                        `json:"min_confidence_score"`
}

// CostConfig 成本配置
type CostConfig struct {
	DefaultModel   string                  `json:"default_model"`
	PricingModels  map[string]PricingModel `json:"pricing_models"`
	Currency       string                  `json:"currency"`
	EnableBatching bool                    `json:"enable_batching"`
	BatchThreshold int                     `json:"batch_threshold"`
}

// PricingModel 定價模型
type PricingModel struct {
	InputPrice    float64 `json:"input_price"`
	OutputPrice   float64 `json:"output_price"`
	CachePrice    float64 `json:"cache_price"`
	BatchDiscount float64 `json:"batch_discount"`
}

// ReporterConfig 報告器配置
type ReporterConfig struct {
	DefaultFormat    string            `json:"default_format"`
	EnabledFormats   []string          `json:"enabled_formats"`
	OutputDirectory  string            `json:"output_directory"`
	TemplateSettings map[string]string `json:"template_settings"`
	ChartEnabled     bool              `json:"chart_enabled"`
	ChartLibrary     string            `json:"chart_library"`
}

// StorageConfig 儲存配置
type StorageConfig struct {
	DataDirectory      string `json:"data_directory"`
	MaxFileSize        int64  `json:"max_file_size"`
	CompressionEnabled bool   `json:"compression_enabled"`
	BackupEnabled      bool   `json:"backup_enabled"`
	BackupInterval     string `json:"backup_interval"`
	RetentionDays      int    `json:"retention_days"`
}

// CLIConfig CLI 配置
type CLIConfig struct {
	EnableColors    bool              `json:"enable_colors"`
	ProgressBar     bool              `json:"progress_bar"`
	VerboseOutput   bool              `json:"verbose_output"`
	InteractiveMode bool              `json:"interactive_mode"`
	DefaultCommands []string          `json:"default_commands"`
	Aliases         map[string]string `json:"aliases"`
}

// ConfigWatcher 配置監聽器介面
type ConfigWatcher interface {
	OnConfigChanged(oldConfig, newConfig *Config) error
}

// NewConfigManager 建立配置管理器
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath: configPath,
		config:     getDefaultConfig(),
		watchers:   make([]ConfigWatcher, 0),
		autoSave:   true,
	}
}

// getDefaultConfig 獲取預設配置
func getDefaultConfig() *Config {
	return &Config{
		Version:     "1.0.0",
		LastUpdated: time.Now(),
		General: GeneralConfig{
			Language:        "zh-TW",
			Timezone:        "Asia/Taipei",
			LogLevel:        "info",
			EnableDebug:     false,
			MaxConcurrency:  4,
			CacheEnabled:    true,
			CacheExpiration: "1h",
		},
		Calculator: CalculatorConfig{
			DefaultMethod:   "estimation",
			TiktokenEnabled: true,
			EstimationRules: map[string]float64{
				"chinese_char": 1.5,
				"english_word": 1.0,
				"punctuation":  0.5,
			},
			CacheResults:   true,
			MaxTokenLength: 100000,
		},
		Analyzer: AnalyzerConfig{
			EnablePatternMatching: true,
			CustomPatterns: map[string][]string{
				"coding":    {"實作", "開發", "程式", "function", "class"},
				"debugging": {"修復", "錯誤", "bug", "debug", "fix"},
				"testing":   {"測試", "test", "驗證", "check"},
			},
			ActivityWeights: map[types.ActivityType]float64{
				types.ActivityCoding:    1.0,
				types.ActivityDebugging: 1.2,
				types.ActivityChat:   0.8,
			},
			MinConfidenceScore: 0.7,
		},
		Cost: CostConfig{
			DefaultModel: "claude-sonnet-4.0",
			PricingModels: map[string]PricingModel{
				"claude-sonnet-4.0": {
					InputPrice:    0.003,
					OutputPrice:   0.015,
					CachePrice:    0.0015,
					BatchDiscount: 0.1,
				},
			},
			Currency:       "USD",
			EnableBatching: true,
			BatchThreshold: 1000,
		},
		Reporter: ReporterConfig{
			DefaultFormat:   "json",
			EnabledFormats:  []string{"json", "csv", "html"},
			OutputDirectory: "./reports",
			TemplateSettings: map[string]string{
				"theme":    "default",
				"language": "zh-TW",
				"timezone": "Asia/Taipei",
			},
			ChartEnabled: true,
			ChartLibrary: "chart.js",
		},
		Storage: StorageConfig{
			DataDirectory:      "./data",
			MaxFileSize:        10 * 1024 * 1024, // 10MB
			CompressionEnabled: true,
			BackupEnabled:      true,
			BackupInterval:     "24h",
			RetentionDays:      30,
		},
		CLI: CLIConfig{
			EnableColors:    true,
			ProgressBar:     true,
			VerboseOutput:   false,
			InteractiveMode: false,
			DefaultCommands: []string{"help"},
			Aliases: map[string]string{
				"calc": "calculate",
				"gen":  "generate",
				"rpt":  "report",
			},
		},
	}
}

// LoadConfig 載入配置
func (cm *ConfigManager) LoadConfig() error {
	// 檢查配置檔案是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// 建立預設配置檔案
		return cm.SaveConfig()
	}

	// 讀取配置檔案
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("讀取配置檔案失敗: %w", err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置檔案失敗: %w", err)
	}

	// 合併預設配置（處理新增的配置項）
	cm.config = cm.mergeWithDefaults(&config)

	// 更新最後修改時間
	if info, err := os.Stat(cm.configPath); err == nil {
		cm.lastModified = info.ModTime()
	}

	return nil
}

// SaveConfig 儲存配置
func (cm *ConfigManager) SaveConfig() error {
	// 確保配置目錄存在
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("建立配置目錄失敗: %w", err)
	}

	// 更新最後修改時間
	cm.config.LastUpdated = time.Now()

	// 序列化為 JSON
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失敗: %w", err)
	}

	// 寫入檔案
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("寫入配置檔案失敗: %w", err)
	}

	return nil
}

// GetConfig 獲取配置
func (cm *ConfigManager) GetConfig() *Config {
	return cm.config
}

// UpdateConfig 更新配置
func (cm *ConfigManager) UpdateConfig(updateFunc func(*Config) error) error {
	oldConfig := *cm.config // 複製舊配置

	// 執行更新函數
	if err := updateFunc(cm.config); err != nil {
		return fmt.Errorf("更新配置失敗: %w", err)
	}

	// 通知監聽器
	for _, watcher := range cm.watchers {
		if err := watcher.OnConfigChanged(&oldConfig, cm.config); err != nil {
			return fmt.Errorf("配置變更通知失敗: %w", err)
		}
	}

	// 自動儲存
	if cm.autoSave {
		return cm.SaveConfig()
	}

	return nil
}

// AddWatcher 添加配置監聽器
func (cm *ConfigManager) AddWatcher(watcher ConfigWatcher) {
	cm.watchers = append(cm.watchers, watcher)
}

// RemoveWatcher 移除配置監聽器
func (cm *ConfigManager) RemoveWatcher(watcher ConfigWatcher) {
	for i, w := range cm.watchers {
		if w == watcher {
			cm.watchers = append(cm.watchers[:i], cm.watchers[i+1:]...)
			break
		}
	}
}

// SetAutoSave 設定自動儲存
func (cm *ConfigManager) SetAutoSave(enabled bool) {
	cm.autoSave = enabled
}

// ValidateConfig 驗證配置
func (cm *ConfigManager) ValidateConfig() error {
	config := cm.config

	// 驗證一般配置
	if config.General.Language == "" {
		return fmt.Errorf("語言設定不能為空")
	}

	if config.General.MaxConcurrency <= 0 {
		return fmt.Errorf("最大並行數必須大於 0")
	}

	// 驗證計算器配置
	if config.Calculator.MaxTokenLength <= 0 {
		return fmt.Errorf("最大 Token 長度必須大於 0")
	}

	// 驗證成本配置
	if config.Cost.DefaultModel == "" {
		return fmt.Errorf("預設模型不能為空")
	}

	// 驗證儲存配置
	if config.Storage.DataDirectory == "" {
		return fmt.Errorf("資料目錄不能為空")
	}

	if config.Storage.MaxFileSize <= 0 {
		return fmt.Errorf("最大檔案大小必須大於 0")
	}

	return nil
}

// ResetToDefaults 重置為預設配置
func (cm *ConfigManager) ResetToDefaults() error {
	oldConfig := *cm.config
	cm.config = getDefaultConfig()

	// 通知監聽器
	for _, watcher := range cm.watchers {
		if err := watcher.OnConfigChanged(&oldConfig, cm.config); err != nil {
			return fmt.Errorf("配置重置通知失敗: %w", err)
		}
	}

	// 自動儲存
	if cm.autoSave {
		return cm.SaveConfig()
	}

	return nil
}

// mergeWithDefaults 與預設配置合併
func (cm *ConfigManager) mergeWithDefaults(config *Config) *Config {
	defaultConfig := getDefaultConfig()

	// 使用反射合併配置
	cm.mergeStructs(reflect.ValueOf(defaultConfig).Elem(), reflect.ValueOf(config).Elem())

	return config
}

// mergeStructs 合併結構體
func (cm *ConfigManager) mergeStructs(defaultVal, configVal reflect.Value) {
	for i := 0; i < defaultVal.NumField(); i++ {
		defaultField := defaultVal.Field(i)
		configField := configVal.Field(i)

		if configField.IsZero() {
			configField.Set(defaultField)
		} else if defaultField.Kind() == reflect.Struct && configField.Kind() == reflect.Struct {
			cm.mergeStructs(defaultField, configField)
		}
	}
}

// GetConfigValue 獲取配置值
func (cm *ConfigManager) GetConfigValue(path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	val := reflect.ValueOf(cm.config).Elem()

	for _, part := range parts {
		val = val.FieldByName(strings.Title(part))
		if !val.IsValid() {
			return nil, fmt.Errorf("配置路徑不存在: %s", path)
		}
	}

	return val.Interface(), nil
}

// SetConfigValue 設定配置值
func (cm *ConfigManager) SetConfigValue(path string, value interface{}) error {
	return cm.UpdateConfig(func(config *Config) error {
		parts := strings.Split(path, ".")
		val := reflect.ValueOf(config).Elem()

		for i, part := range parts {
			if i == len(parts)-1 {
				field := val.FieldByName(strings.Title(part))
				if !field.IsValid() {
					return fmt.Errorf("配置欄位不存在: %s", part)
				}

				if !field.CanSet() {
					return fmt.Errorf("配置欄位不可設定: %s", part)
				}

				field.Set(reflect.ValueOf(value))
			} else {
				val = val.FieldByName(strings.Title(part))
				if !val.IsValid() {
					return fmt.Errorf("配置路徑不存在: %s", path)
				}
			}
		}

		return nil
	})
}

// ExportConfig 匯出配置
func (cm *ConfigManager) ExportConfig(outputPath string) error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失敗: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

// ImportConfig 匯入配置
func (cm *ConfigManager) ImportConfig(inputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("讀取配置檔案失敗: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置檔案失敗: %w", err)
	}

	oldConfig := *cm.config
	cm.config = cm.mergeWithDefaults(&config)

	// 通知監聽器
	for _, watcher := range cm.watchers {
		if err := watcher.OnConfigChanged(&oldConfig, cm.config); err != nil {
			return fmt.Errorf("配置匯入通知失敗: %w", err)
		}
	}

	// 自動儲存
	if cm.autoSave {
		return cm.SaveConfig()
	}

	return nil
}
