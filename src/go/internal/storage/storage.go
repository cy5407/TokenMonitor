package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"path/filepath"
	"sort"
	"strings"
	"time"

	"token-monitor/internal/types"
)

// StorageInterface 資料儲存介面
type StorageInterface interface {
	SaveUsageRecord(record types.UsageRecord) error
	LoadUsageRecords(timeRange types.TimeRange) ([]types.UsageRecord, error)
	SaveActivityData(activities []types.Activity) error
	LoadActivityData(timeRange types.TimeRange) ([]types.Activity, error)
	DeleteOldRecords(before time.Time) error
	GetStorageStats() (*StorageStats, error)
}

// JSONStorage JSON 格式資料儲存實作
type JSONStorage struct {
	dataDir     string
	maxFileSize int64
}

// StorageStats 儲存統計資訊
type StorageStats struct {
	TotalRecords    int       `json:"total_records"`
	TotalFiles      int       `json:"total_files"`
	TotalSize       int64     `json:"total_size"`
	OldestRecord    time.Time `json:"oldest_record"`
	NewestRecord    time.Time `json:"newest_record"`
	LastMaintenance time.Time `json:"last_maintenance"`
}

// NewJSONStorage 建立新的 JSON 儲存實例
func NewJSONStorage(dataDir string) *JSONStorage {
	return &JSONStorage{
		dataDir:     dataDir,
		maxFileSize: 10 * 1024 * 1024, // 10MB
	}
}

// SaveUsageRecord 儲存使用記錄
func (js *JSONStorage) SaveUsageRecord(record types.UsageRecord) error {
	if err := js.ensureDataDir(); err != nil {
		return err
	}

	filename := js.getRecordFilename(record.Timestamp)
	return js.appendToFile(filename, record)
}

// LoadUsageRecords 載入使用記錄
func (js *JSONStorage) LoadUsageRecords(timeRange types.TimeRange) ([]types.UsageRecord, error) {
	var records []types.UsageRecord

	files, err := js.getRecordFiles(timeRange)
	if err != nil {
		return nil, fmt.Errorf("獲取記錄檔案失敗: %w", err)
	}

	for _, file := range files {
		fileRecords, err := js.loadRecordsFromFile(file)
		if err != nil {
			// 記錄錯誤但繼續處理其他檔案
			continue
		}

		// 過濾時間範圍內的記錄
		for _, record := range fileRecords {
			if js.isInTimeRange(record.Timestamp, timeRange) {
				records = append(records, record)
			}
		}
	}

	// 按時間排序
	js.sortRecordsByTime(records)

	return records, nil
}

// SaveActivityData 儲存活動資料
func (js *JSONStorage) SaveActivityData(activities []types.Activity) error {
	// This is a simplified implementation. In a real scenario, you might want to merge activities.
	if err := js.ensureDataDir(); err != nil {
		return err
	}

	filename := filepath.Join(js.dataDir, "activities.json")
	data, err := json.MarshalIndent(activities, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadActivityData 載入活動資料
func (js *JSONStorage) LoadActivityData(timeRange types.TimeRange) ([]types.Activity, error) {
	filename := filepath.Join(js.dataDir, "activities.json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []types.Activity{}, nil // Return empty slice if file doesn't exist
		}
		return nil, err
	}

	var allActivities []types.Activity
	if err := json.Unmarshal(data, &allActivities); err != nil {
		return nil, err
	}

	var filteredActivities []types.Activity
	for _, activity := range allActivities {
		if js.isInTimeRange(activity.Timestamp, timeRange) {
			filteredActivities = append(filteredActivities, activity)
		}
	}

	return filteredActivities, nil
}

// DeleteOldRecords 刪除舊記錄
func (js *JSONStorage) DeleteOldRecords(before time.Time) error {
	return js.DeleteOldRecordsImpl(before)
}

// GetStorageStats 獲取儲存統計資訊
func (js *JSONStorage) GetStorageStats() (*StorageStats, error) {
	return js.GetStorageStatsImpl()
}

// 輔助方法
func (js *JSONStorage) ensureDataDir() error {
	return os.MkdirAll(js.dataDir, 0755)
}

func (js *JSONStorage) getRecordFilename(timestamp time.Time) string {
	return filepath.Join(js.dataDir, fmt.Sprintf("records_%s.json", timestamp.Format("2006-01-02")))
}

func (js *JSONStorage) appendToFile(filename string, record types.UsageRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	return err
}

// shouldRotateFile 檢查是否需要輪轉檔案
func (js *JSONStorage) shouldRotateFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false // 檔案不存在，不需要輪轉
	}
	return info.Size() > js.maxFileSize
}

// getRotatedFilename 獲取輪轉後的檔案名
func (js *JSONStorage) getRotatedFilename(timestamp time.Time) string {
	base := filepath.Join(js.dataDir, fmt.Sprintf("records_%s", timestamp.Format("2006-01-02")))
	counter := 1
	for {
		filename := fmt.Sprintf("%s_%d.json", base, counter)
		if !js.shouldRotateFile(filename) {
			return filename
		}
		counter++
	}
}

// LoadUsageRecordsImpl 載入使用記錄的實作
func (js *JSONStorage) LoadUsageRecordsImpl(timeRange types.TimeRange) ([]types.UsageRecord, error) {
	var records []types.UsageRecord

	files, err := js.getRecordFiles(timeRange)
	if err != nil {
		return nil, fmt.Errorf("獲取記錄檔案失敗: %w", err)
	}

	for _, file := range files {
		fileRecords, err := js.loadRecordsFromFile(file)
		if err != nil {
			continue // 跳過損壞的檔案
		}

		// 過濾時間範圍
		for _, record := range fileRecords {
			if js.isInTimeRange(record.Timestamp, timeRange) {
				records = append(records, record)
			}
		}
	}

	return records, nil
}

// getRecordFiles 獲取指定時間範圍內的記錄檔案
func (js *JSONStorage) getRecordFiles(timeRange types.TimeRange) ([]string, error) {
	var files []string

	err := filepath.Walk(js.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), "records_") && strings.HasSuffix(info.Name(), ".json") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// loadRecordsFromFile 從檔案載入記錄
func (js *JSONStorage) loadRecordsFromFile(filename string) ([]types.UsageRecord, error) {
	var records []types.UsageRecord

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for decoder.More() {
		var record types.UsageRecord
		if err := decoder.Decode(&record); err != nil {
			continue // 跳過損壞的記錄
		}
		records = append(records, record)
	}

	return records, nil
}

// isInTimeRange 檢查時間是否在範圍內
func (js *JSONStorage) isInTimeRange(timestamp time.Time, timeRange types.TimeRange) bool {
	if !timeRange.Start.IsZero() && timestamp.Before(timeRange.Start) {
		return false
	}
	if !timeRange.End.IsZero() && timestamp.After(timeRange.End) {
		return false
	}
	return true
}

// sortRecordsByTime 按時間排序記錄
func (js *JSONStorage) sortRecordsByTime(records []types.UsageRecord) {
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.Before(records[j].Timestamp)
	})
}

// GetStorageStatsImpl 獲取儲存統計資訊的實作
func (js *JSONStorage) GetStorageStatsImpl() (*StorageStats, error) {
	stats := &StorageStats{
		LastMaintenance: time.Now(),
	}

	err := filepath.Walk(js.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".json") {
			stats.TotalFiles++
			stats.TotalSize += info.Size()

			// 更新最舊和最新記錄時間
			if records, err := js.loadRecordsFromFile(path); err == nil && len(records) > 0 {
				stats.TotalRecords += len(records)

				firstRecord := records[0].Timestamp
				lastRecord := records[len(records)-1].Timestamp

				if stats.OldestRecord.IsZero() || firstRecord.Before(stats.OldestRecord) {
					stats.OldestRecord = firstRecord
				}

				if stats.NewestRecord.IsZero() || lastRecord.After(stats.NewestRecord) {
					stats.NewestRecord = lastRecord
				}
			}
		}

		return nil
	})

	return stats, err
}

// DeleteOldRecordsImpl 刪除舊記錄的實作
func (js *JSONStorage) DeleteOldRecordsImpl(before time.Time) error {
	var deletedFiles []string

	err := filepath.Walk(js.dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".json") {
			// 檢查檔案中的記錄是否都在指定時間之前
			if js.shouldDeleteFile(path, before) {
				if err := os.Remove(path); err == nil {
					deletedFiles = append(deletedFiles, path)
				}
			}
		}

		return nil
	})

	if len(deletedFiles) > 0 {
		fmt.Printf("已刪除 %d 個舊記錄檔案\n", len(deletedFiles))
	}

	return err
}

// shouldDeleteFile 檢查是否應該刪除檔案
func (js *JSONStorage) shouldDeleteFile(filename string, before time.Time) bool {
	records, err := js.loadRecordsFromFile(filename)
	if err != nil {
		return false
	}

	// 如果檔案中所有記錄都在指定時間之前，則可以刪除
	for _, record := range records {
		if !record.Timestamp.Before(before) {
			return false
		}
	}

	return len(records) > 0
}
