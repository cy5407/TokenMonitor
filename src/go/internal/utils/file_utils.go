package utils

import (
	"bufio"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo 檔案資訊結構
type FileInfo struct {
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	IsDir     bool      `json:"is_dir"`
	Extension string    `json:"extension"`
	MD5Hash   string    `json:"md5_hash,omitempty"`
	LineCount int       `json:"line_count,omitempty"`
	Encoding  string    `json:"encoding,omitempty"`
}

// FileProcessor 檔案處理器
type FileProcessor struct {
	maxFileSize   int64
	allowedExts   []string
	excludeDirs   []string
	enableHashing bool
}

// NewFileProcessor 建立檔案處理器
func NewFileProcessor() *FileProcessor {
	return &FileProcessor{
		maxFileSize:   100 * 1024 * 1024, // 100MB
		allowedExts:   []string{".txt", ".log", ".json", ".csv", ".md"},
		excludeDirs:   []string{".git", "node_modules", ".kiro"},
		enableHashing: true,
	}
}

// SetMaxFileSize 設定最大檔案大小
func (fp *FileProcessor) SetMaxFileSize(size int64) {
	fp.maxFileSize = size
}

// SetAllowedExtensions 設定允許的副檔名
func (fp *FileProcessor) SetAllowedExtensions(exts []string) {
	fp.allowedExts = exts
}

// SetExcludeDirs 設定排除的目錄
func (fp *FileProcessor) SetExcludeDirs(dirs []string) {
	fp.excludeDirs = dirs
}

// ScanDirectory 掃描目錄
func (fp *FileProcessor) ScanDirectory(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳過排除的目錄
		if info.IsDir() && fp.shouldExcludeDir(info.Name()) {
			return filepath.SkipDir
		}

		// 處理檔案
		if !info.IsDir() {
			fileInfo, err := fp.processFile(path, info)
			if err != nil {
				// 記錄錯誤但繼續處理
				fmt.Printf("⚠️  處理檔案失敗 %s: %v\n", path, err)
				return nil
			}

			if fileInfo != nil {
				files = append(files, *fileInfo)
			}
		}

		return nil
	})

	return files, err
}

// processFile 處理單個檔案
func (fp *FileProcessor) processFile(path string, info os.FileInfo) (*FileInfo, error) {
	// 檢查檔案大小
	if info.Size() > fp.maxFileSize {
		return nil, fmt.Errorf("檔案過大: %d bytes", info.Size())
	}

	// 檢查副檔名
	ext := strings.ToLower(filepath.Ext(path))
	if !fp.isAllowedExtension(ext) {
		return nil, nil // 跳過不支援的檔案類型
	}

	fileInfo := &FileInfo{
		Path:      path,
		Size:      info.Size(),
		ModTime:   info.ModTime(),
		IsDir:     info.IsDir(),
		Extension: ext,
	}

	// 計算 MD5 雜湊
	if fp.enableHashing {
		hash, err := fp.calculateMD5(path)
		if err != nil {
			return nil, fmt.Errorf("計算 MD5 失敗: %w", err)
		}
		fileInfo.MD5Hash = hash
	}

	// 計算行數
	if fp.isTextFile(ext) {
		lineCount, encoding, err := fp.analyzeTextFile(path)
		if err != nil {
			return nil, fmt.Errorf("分析文字檔案失敗: %w", err)
		}
		fileInfo.LineCount = lineCount
		fileInfo.Encoding = encoding
	}

	return fileInfo, nil
}

// shouldExcludeDir 檢查是否應排除目錄
func (fp *FileProcessor) shouldExcludeDir(dirName string) bool {
	for _, excludeDir := range fp.excludeDirs {
		if dirName == excludeDir {
			return true
		}
	}
	return false
}

// isAllowedExtension 檢查是否為允許的副檔名
func (fp *FileProcessor) isAllowedExtension(ext string) bool {
	for _, allowedExt := range fp.allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// isTextFile 檢查是否為文字檔案
func (fp *FileProcessor) isTextFile(ext string) bool {
	textExts := []string{".txt", ".log", ".md", ".json", ".csv", ".yaml", ".yml"}
	for _, textExt := range textExts {
		if ext == textExt {
			return true
		}
	}
	return false
}

// calculateMD5 計算檔案 MD5 雜湊
func (fp *FileProcessor) calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// analyzeTextFile 分析文字檔案
func (fp *FileProcessor) analyzeTextFile(filePath string) (int, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	encoding := "UTF-8" // 預設編碼

	// 讀取前幾行來檢測編碼
	sampleLines := 0
	for scanner.Scan() && sampleLines < 10 {
		line := scanner.Text()

		// 簡單的編碼檢測
		if !isValidUTF8(line) {
			encoding = "Unknown"
		}

		lineCount++
		sampleLines++
	}

	// 繼續計算剩餘行數
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return lineCount, encoding, err
	}

	return lineCount, encoding, nil
}

// isValidUTF8 檢查字串是否為有效的 UTF-8
func isValidUTF8(s string) bool {
	return strings.ToValidUTF8(s, "") == s
}

// CompressFile 壓縮檔案
func CompressFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	gzWriter := gzip.NewWriter(dstFile)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, srcFile)
	return err
}

// DecompressFile 解壓縮檔案
func DecompressFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	gzReader, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, gzReader)
	return err
}

// CopyFile 複製檔案
func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 確保目標目錄存在
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MoveFile 移動檔案
func MoveFile(srcPath, dstPath string) error {
	// 先嘗試重新命名（同一檔案系統）
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}

	// 如果重新命名失敗，則複製後刪除
	if err := CopyFile(srcPath, dstPath); err != nil {
		return err
	}

	return os.Remove(srcPath)
}

// EnsureDir 確保目錄存在
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// CleanupOldFiles 清理舊檔案
func CleanupOldFiles(dirPath string, maxAge time.Duration) error {
	cutoffTime := time.Now().Add(-maxAge)

	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.ModTime().Before(cutoffTime) {
			fmt.Printf("🗑️  刪除舊檔案: %s\n", path)
			return os.Remove(path)
		}

		return nil
	})
}

// GetDirSize 獲取目錄大小
func GetDirSize(dirPath string) (int64, error) {
	var size int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}

		return nil
	})

	return size, err
}

// FindDuplicateFiles 尋找重複檔案
func FindDuplicateFiles(dirPath string) (map[string][]string, error) {
	fp := NewFileProcessor()
	files, err := fp.ScanDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	hashMap := make(map[string][]string)

	for _, file := range files {
		if file.MD5Hash != "" {
			hashMap[file.MD5Hash] = append(hashMap[file.MD5Hash], file.Path)
		}
	}

	// 只返回有重複的檔案
	duplicates := make(map[string][]string)
	for hash, paths := range hashMap {
		if len(paths) > 1 {
			duplicates[hash] = paths
		}
	}

	return duplicates, nil
}
