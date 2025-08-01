package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"token-monitor/internal/types"
)

// TestNewJSONStorage 測試 JSON 儲存建立
func TestNewJSONStorage(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	if storage == nil {
		t.Fatal("JSON 儲存建立失敗")
	}

	if storage.dataDir != tempDir {
		t.Errorf("資料目錄設定錯誤: 期望 %s, 得到 %s", tempDir, storage.dataDir)
	}
}

// TestSaveUsageRecord 測試使用記錄儲存
func TestSaveUsageRecord(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	record := types.UsageRecord{
		Timestamp: time.Now(),
		SessionID: "test-session",
		Activity: types.Activity{
			Type: types.ActivityCoding,
		},
	}

	err := storage.SaveUsageRecord(record)
	if err != nil {
		t.Fatalf("儲存使用記錄失敗: %v", err)
	}

	// 驗證檔案是否建立
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("讀取目錄失敗: %v", err)
	}

	if len(files) == 0 {
		t.Error("應該建立至少一個檔案")
	}
}

// TestEnsureDataDir 測試資料目錄建立
func TestEnsureDataDir(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	storage := NewJSONStorage(subDir)

	err := storage.ensureDataDir()
	if err != nil {
		t.Fatalf("建立資料目錄失敗: %v", err)
	}

	// 驗證目錄是否存在
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Error("資料目錄應該被建立")
	}
}

// TestGetRecordFilename 測試記錄檔案名生成
func TestGetRecordFilename(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewJSONStorage(tempDir)

	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	filename := storage.getRecordFilename(timestamp)

	expectedFilename := filepath.Join(tempDir, "records_2024-01-15.json")
	if filename != expectedFilename {
		t.Errorf("檔案名錯誤: 期望 %s, 得到 %s", expectedFilename, filename)
	}
}
