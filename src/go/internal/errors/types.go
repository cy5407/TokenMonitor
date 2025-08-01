package errors

import (
	"time"
)

// ErrorCode 定義錯誤代碼類型
type ErrorCode string

// ErrorSeverity 定義錯誤嚴重級別
type ErrorSeverity string

const (
	// 錯誤嚴重級別
	SeverityCritical ErrorSeverity = "critical" // 關鍵錯誤，系統無法繼續運作
	SeverityHigh     ErrorSeverity = "high"     // 高級錯誤，影響核心功能
	SeverityMedium   ErrorSeverity = "medium"   // 中級錯誤，影響部分功能
	SeverityLow      ErrorSeverity = "low"      // 低級錯誤，不影響主要功能
	SeverityWarning  ErrorSeverity = "warning"  // 警告，可能的問題
)

const (
	// Token 計算相關錯誤
	ErrCodeTokenCalculation      ErrorCode = "TOKEN_CALCULATION_FAILED"
	ErrCodeTiktokenUnavailable   ErrorCode = "TIKTOKEN_UNAVAILABLE"
	ErrCodeInvalidText           ErrorCode = "INVALID_TEXT_INPUT"
	ErrCodeTokenCountExceeded    ErrorCode = "TOKEN_COUNT_EXCEEDED"
	ErrCodeCalculationTimeout    ErrorCode = "CALCULATION_TIMEOUT"

	// 成本計算相關錯誤
	ErrCodeCostCalculation       ErrorCode = "COST_CALCULATION_FAILED"
	ErrCodeInvalidPricingModel   ErrorCode = "INVALID_PRICING_MODEL"
	ErrCodePricingDataMissing    ErrorCode = "PRICING_DATA_MISSING"
	ErrCodeInvalidTokenCount     ErrorCode = "INVALID_TOKEN_COUNT"

	// 活動分析相關錯誤
	ErrCodeActivityAnalysis      ErrorCode = "ACTIVITY_ANALYSIS_FAILED"
	ErrCodeActivityClassification ErrorCode = "ACTIVITY_CLASSIFICATION_FAILED"
	ErrCodePatternLoadFailed     ErrorCode = "PATTERN_LOAD_FAILED"
	ErrCodeInvalidActivityType   ErrorCode = "INVALID_ACTIVITY_TYPE"

	// 資料存取相關錯誤
	ErrCodeDataAccess            ErrorCode = "DATA_ACCESS_FAILED"
	ErrCodeFileNotFound          ErrorCode = "FILE_NOT_FOUND"
	ErrCodeFilePermission        ErrorCode = "FILE_PERMISSION_DENIED"
	ErrCodeDataCorruption        ErrorCode = "DATA_CORRUPTION"
	ErrCodeStorageInsufficient   ErrorCode = "STORAGE_INSUFFICIENT"

	// 配置相關錯誤
	ErrCodeConfigLoad            ErrorCode = "CONFIG_LOAD_FAILED"
	ErrCodeConfigValidation      ErrorCode = "CONFIG_VALIDATION_FAILED"
	ErrCodeConfigMissing         ErrorCode = "CONFIG_MISSING"
	ErrCodeInvalidConfigFormat   ErrorCode = "INVALID_CONFIG_FORMAT"

	// 網路相關錯誤
	ErrCodeNetworkConnection     ErrorCode = "NETWORK_CONNECTION_FAILED"
	ErrCodeAPIRateLimit          ErrorCode = "API_RATE_LIMIT_EXCEEDED"
	ErrCodeAPIAuthentication     ErrorCode = "API_AUTHENTICATION_FAILED"
	ErrCodeRequestTimeout        ErrorCode = "REQUEST_TIMEOUT"

	// 系統相關錯誤
	ErrCodeSystemResource        ErrorCode = "SYSTEM_RESOURCE_UNAVAILABLE"
	ErrCodeMemoryInsufficient    ErrorCode = "MEMORY_INSUFFICIENT"
	ErrCodeDependencyMissing     ErrorCode = "DEPENDENCY_MISSING"
	ErrCodeInitializationFailed  ErrorCode = "INITIALIZATION_FAILED"

	// 報告生成相關錯誤
	ErrCodeReportGeneration      ErrorCode = "REPORT_GENERATION_FAILED"
	ErrCodeInvalidReportFormat   ErrorCode = "INVALID_REPORT_FORMAT"
	ErrCodeReportDataMissing     ErrorCode = "REPORT_DATA_MISSING"
	ErrCodeExportFailed          ErrorCode = "EXPORT_FAILED"
)

// ErrorCategory 定義錯誤類別
type ErrorCategory string

const (
	CategoryToken   ErrorCategory = "token"
	CategoryCost    ErrorCategory = "cost"
	CategoryActivity ErrorCategory = "activity"
	CategoryData    ErrorCategory = "data"
	CategoryConfig  ErrorCategory = "config"
	CategoryNetwork ErrorCategory = "network"
	CategorySystem  ErrorCategory = "system"
	CategoryReport  ErrorCategory = "report"
)

// RetryPolicy 定義重試策略
type RetryPolicy struct {
	MaxRetries      int           `json:"max_retries"`
	InitialDelay    time.Duration `json:"initial_delay"`
	MaxDelay        time.Duration `json:"max_delay"`
	BackoffFactor   float64       `json:"backoff_factor"`
	RetryableErrors []ErrorCode   `json:"retryable_errors"`
}

// ErrorMetadata 錯誤元數據
type ErrorMetadata struct {
	Code        ErrorCode     `json:"code"`
	Category    ErrorCategory `json:"category"`
	Severity    ErrorSeverity `json:"severity"`
	Message     string        `json:"message"`
	MessageZH   string        `json:"message_zh"`
	Description string        `json:"description"`
	Solution    string        `json:"solution"`
	SolutionZH  string        `json:"solution_zh"`
	Retryable   bool          `json:"retryable"`
	RetryPolicy *RetryPolicy  `json:"retry_policy,omitempty"`
	HelpURL     string        `json:"help_url,omitempty"`
}

// 預定義錯誤元數據
var errorMetadataMap = map[ErrorCode]ErrorMetadata{
	ErrCodeTokenCalculation: {
		Code:        ErrCodeTokenCalculation,
		Category:    CategoryToken,
		Severity:    SeverityHigh,
		Message:     "Token calculation failed",
		MessageZH:   "Token 計算失敗",
		Description: "Failed to calculate token count for the provided text",
		Solution:    "Check input text format and try again. If the problem persists, try using a different calculation method.",
		SolutionZH:  "檢查輸入文本格式並重試。如果問題持續存在，請嘗試使用不同的計算方法。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 10,
			BackoffFactor: 2.0,
		},
	},
	ErrCodeTiktokenUnavailable: {
		Code:        ErrCodeTiktokenUnavailable,
		Category:    CategoryToken,
		Severity:    SeverityMedium,
		Message:     "Tiktoken library is not available",
		MessageZH:   "Tiktoken 函式庫不可用",
		Description: "The tiktoken library is not installed or not accessible",
		Solution:    "Install tiktoken library or use alternative calculation method (estimation).",
		SolutionZH:  "安裝 tiktoken 函式庫或使用替代計算方法（估算）。",
		Retryable:   false,
	},
	ErrCodeInvalidText: {
		Code:        ErrCodeInvalidText,
		Category:    CategoryToken,
		Severity:    SeverityLow,
		Message:     "Invalid text input",
		MessageZH:   "無效的文本輸入",
		Description: "The provided text input is invalid or empty",
		Solution:    "Provide valid, non-empty text input.",
		SolutionZH:  "提供有效的非空文本輸入。",
		Retryable:   false,
	},
	ErrCodeCostCalculation: {
		Code:        ErrCodeCostCalculation,
		Category:    CategoryCost,
		Severity:    SeverityHigh,
		Message:     "Cost calculation failed",
		MessageZH:   "成本計算失敗",
		Description: "Failed to calculate cost for the given token usage",
		Solution:    "Check pricing model configuration and token counts. Ensure all required parameters are provided.",
		SolutionZH:  "檢查定價模型配置和 Token 計數。確保提供所有必需的參數。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    2,
			InitialDelay:  time.Millisecond * 500,
			MaxDelay:      time.Second * 5,
			BackoffFactor: 2.0,
		},
	},
	ErrCodeInvalidPricingModel: {
		Code:        ErrCodeInvalidPricingModel,
		Category:    CategoryCost,
		Severity:    SeverityMedium,
		Message:     "Invalid pricing model",
		MessageZH:   "無效的定價模型",
		Description: "The specified pricing model is not supported or invalid",
		Solution:    "Use a supported pricing model. Check available models with GetSupportedModels().",
		SolutionZH:  "使用支援的定價模型。使用 GetSupportedModels() 檢查可用模型。",
		Retryable:   false,
	},
	ErrCodeDataAccess: {
		Code:        ErrCodeDataAccess,
		Category:    CategoryData,
		Severity:    SeverityHigh,
		Message:     "Data access failed",
		MessageZH:   "資料存取失敗",
		Description: "Failed to access or manipulate data",
		Solution:    "Check file permissions and disk space. Ensure the data file exists and is not corrupted.",
		SolutionZH:  "檢查檔案權限和磁碟空間。確保資料檔案存在且未損壞。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 15,
			BackoffFactor: 2.0,
		},
	},
	ErrCodeFileNotFound: {
		Code:        ErrCodeFileNotFound,
		Category:    CategoryData,
		Severity:    SeverityMedium,
		Message:     "File not found",
		MessageZH:   "檔案未找到",
		Description: "The specified file does not exist",
		Solution:    "Check the file path and ensure the file exists. Create the file if necessary.",
		SolutionZH:  "檢查檔案路徑並確保檔案存在。如有必要，請建立檔案。",
		Retryable:   false,
	},
	ErrCodeConfigLoad: {
		Code:        ErrCodeConfigLoad,
		Category:    CategoryConfig,
		Severity:    SeverityHigh,
		Message:     "Configuration load failed",
		MessageZH:   "配置載入失敗",
		Description: "Failed to load configuration file",
		Solution:    "Check configuration file exists and has valid format. Restore from backup if necessary.",
		SolutionZH:  "檢查配置檔案是否存在且格式有效。如有必要，請從備份恢復。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    2,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 5,
			BackoffFactor: 1.5,
		},
	},
	ErrCodeNetworkConnection: {
		Code:        ErrCodeNetworkConnection,
		Category:    CategoryNetwork,
		Severity:    SeverityMedium,
		Message:     "Network connection failed",
		MessageZH:   "網路連接失敗",
		Description: "Failed to establish network connection",
		Solution:    "Check network connectivity and proxy settings. Verify firewall configuration.",
		SolutionZH:  "檢查網路連接和代理設定。驗證防火牆配置。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    5,
			InitialDelay:  time.Second * 2,
			MaxDelay:      time.Second * 30,
			BackoffFactor: 2.0,
		},
	},
	ErrCodeAPIRateLimit: {
		Code:        ErrCodeAPIRateLimit,
		Category:    CategoryNetwork,
		Severity:    SeverityMedium,
		Message:     "API rate limit exceeded",
		MessageZH:   "API 速率限制超出",
		Description: "The API rate limit has been exceeded",
		Solution:    "Wait before making additional requests. Consider implementing request throttling.",
		SolutionZH:  "等待後再進行額外請求。考慮實作請求節流。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    3,
			InitialDelay:  time.Second * 10,
			MaxDelay:      time.Minute * 5,
			BackoffFactor: 3.0,
		},
	},
	ErrCodeSystemResource: {
		Code:        ErrCodeSystemResource,
		Category:    CategorySystem,
		Severity:    SeverityCritical,
		Message:     "System resource unavailable",
		MessageZH:   "系統資源不可用",
		Description: "Required system resources are not available",
		Solution:    "Free up system resources (memory, disk space) and try again.",
		SolutionZH:  "釋放系統資源（記憶體、磁碟空間）並重試。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    2,
			InitialDelay:  time.Second * 5,
			MaxDelay:      time.Second * 30,
			BackoffFactor: 3.0,
		},
	},
	ErrCodeReportGeneration: {
		Code:        ErrCodeReportGeneration,
		Category:    CategoryReport,
		Severity:    SeverityMedium,
		Message:     "Report generation failed",
		MessageZH:   "報告生成失敗",
		Description: "Failed to generate the requested report",
		Solution:    "Check input data validity and output path permissions. Try a different report format.",
		SolutionZH:  "檢查輸入資料有效性和輸出路徑權限。嘗試不同的報告格式。",
		Retryable:   true,
		RetryPolicy: &RetryPolicy{
			MaxRetries:    2,
			InitialDelay:  time.Second,
			MaxDelay:      time.Second * 10,
			BackoffFactor: 2.0,
		},
	},
}

// GetErrorMetadata 取得錯誤元數據
func GetErrorMetadata(code ErrorCode) (ErrorMetadata, bool) {
	metadata, exists := errorMetadataMap[code]
	return metadata, exists
}

// IsRetryable 檢查錯誤是否可重試
func IsRetryable(code ErrorCode) bool {
	if metadata, exists := errorMetadataMap[code]; exists {
		return metadata.Retryable
	}
	return false
}

// GetRetryPolicy 取得重試策略
func GetRetryPolicy(code ErrorCode) *RetryPolicy {
	if metadata, exists := errorMetadataMap[code]; exists {
		return metadata.RetryPolicy
	}
	return nil
}

// GetErrorsByCategory 依類別取得錯誤
func GetErrorsByCategory(category ErrorCategory) []ErrorCode {
	var codes []ErrorCode
	for code, metadata := range errorMetadataMap {
		if metadata.Category == category {
			codes = append(codes, code)
		}
	}
	return codes
}

// GetErrorsBySeverity 依嚴重級別取得錯誤
func GetErrorsBySeverity(severity ErrorSeverity) []ErrorCode {
	var codes []ErrorCode
	for code, metadata := range errorMetadataMap {
		if metadata.Severity == severity {
			codes = append(codes, code)
		}
	}
	return codes
}