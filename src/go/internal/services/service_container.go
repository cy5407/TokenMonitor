package services

import (
	"sync"

	"token-monitor/internal/analyzer"
	"token-monitor/internal/calculator"
	"token-monitor/internal/config"
	"token-monitor/internal/cost"
	"token-monitor/internal/interfaces"
	"token-monitor/internal/reporter"
	"token-monitor/internal/storage"

	"github.com/spf13/viper"
)

// ServiceContainer 負責管理所有服務的單一實例
type ServiceContainer struct {
	ConfigManager    *config.ConfigManager
	TokenCalculator  interfaces.TokenCalculator
	ActivityAnalyzer *analyzer.ActivityAnalyzer
	CostCalculator   *cost.CostCalculatorImpl
	ReportGenerator  *reporter.ReportGenerator
	Storage          storage.StorageInterface
}

var (
	instance *ServiceContainer
	once     sync.Once
)

// GetInstance 返回 ServiceContainer 的單一實例
func GetInstance() *ServiceContainer {
	once.Do(func() {
		// 建立 ConfigManager
		cm := config.NewConfigManager(viper.GetString("config"))
		cm.LoadConfig()

		// 建立 TokenCalculator
		calc := calculator.NewTokenCalculator(viper.GetInt("token_calculation.cache_size"))

		// 建立 ActivityAnalyzer
		an := analyzer.NewActivityAnalyzer()

		// 建立 CostCalculator
		costCalc := cost.NewCostCalculator()
		costCalc.LoadPricingModels(viper.GetString("pricing_models_path"))

		// 建立 ReportGenerator
		rg := reporter.NewReportGenerator()

		// 建立 Storage
		st := storage.NewJSONStorage(viper.GetString("storage.path"))

		instance = &ServiceContainer{
			ConfigManager:    cm,
			TokenCalculator:  calc,
			ActivityAnalyzer: an,
			CostCalculator:   costCalc,
			ReportGenerator:  rg,
			Storage:          st,
		}
	})
	return instance
}