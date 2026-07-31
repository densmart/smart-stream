package utils

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	ffprobe "gopkg.in/vansante/go-ffprobe.v2"

	"github.com/densmart/smart-stream/pkg/logger"
)

// MediaMetadata содержит информацию о медиа-файле
type MediaMetadata struct {
	Duration int    // Длительность в секундах
	Size     uint64 // Размер файла в байтах
}

// FileSystemItemType тип элемента файловой системы
type FileSystemItemType string

const (
	FileSystemItemTypeFile      FileSystemItemType = "file"
	FileSystemItemTypeDirectory FileSystemItemType = "directory"
)

// FileSystemItem представляет файл или директорию
type FileSystemItem struct {
	Name       string
	Type       FileSystemItemType
	Path       string
	Size       *uint64
	Format     *string
	ModifiedAt *time.Time
}

// Поддерживаемые форматы видео
var supportedFormats = map[string]bool{
	".mp4":  true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".webm": true,
	".flv":  true,
	".wmv":  true,
	".m4v":  true,
	".3gp":  true,
	".ts":   true,
}

// GetMediaMetadata извлекает метаданные из медиа-файла используя ffprobe
func GetMediaMetadata(filePath string) (*MediaMetadata, error) {
	// Проверяем существование файла
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Errorf("[utils] GetMediaMetadata file not found: %s", filePath)
			return nil, fmt.Errorf("media file not found: %s", filePath)
		}
		logger.Errorf("[utils] GetMediaMetadata error accessing file: %s - %v", filePath, err)
		return nil, fmt.Errorf("failed to access media file: %w", err)
	}

	// Проверяем, что это файл, а не директория
	if fileInfo.IsDir() {
		logger.Errorf("[utils] GetMediaMetadata path is a directory: %s", filePath)
		return nil, fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(filePath))
	if !supportedFormats[ext] {
		logger.Errorf("[utils] GetMediaMetadata unsupported format: %s", ext)
		return nil, fmt.Errorf("unsupported media format: %s (supported: mp4, mkv, avi, mov, webm)", ext)
	}

	// Получаем размер файла
	fileSize := uint64(fileInfo.Size())

	// Используем ffprobe для получения длительности
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		logger.Errorf("[utils] GetMediaMetadata ffprobe error: %s - %v", filePath, err)
		return nil, fmt.Errorf("failed to read media file metadata: %w", err)
	}

	// Получаем длительность из формата
	if data.Format == nil || data.Format.DurationSeconds == 0 {
		// Если не нашли длительность в формате
		logger.Errorf("[utils] GetMediaMetadata no duration found: %s", filePath)
		return nil, fmt.Errorf("failed to determine media duration")
	}

	// Округляем длительность до целых секунд
	// Duration() возвращает time.Duration, нужно преобразовать в секунды
	durationSeconds := int(math.Round(data.Format.DurationSeconds))

	logger.Infof("[utils] GetMediaMetadata success: %s - duration: %ds, size: %d bytes", filePath, durationSeconds, fileSize)

	return &MediaMetadata{
		Duration: durationSeconds,
		Size:     fileSize,
	}, nil
}

// BrowseDirectory читает содержимое директории и возвращает список файлов и папок
// basePath - базовый путь (например, /storage/media)
// relativePath - относительный путь внутри basePath (например, "alisa" или "alisa/subfolder")
func BrowseDirectory(basePath, relativePath string) ([]FileSystemItem, string, error) {
	// Нормализация путей для защиты от path traversal
	basePath = filepath.Clean(basePath)
	relativePath = filepath.Clean(relativePath)

	// Проверка на попытки выхода за пределы basePath
	if strings.Contains(relativePath, "..") {
		logger.Errorf("[utils] BrowseDirectory path traversal attempt: %s", relativePath)
		return nil, "", fmt.Errorf("invalid path: path traversal not allowed")
	}

	// Проверка на абсолютный путь
	if filepath.IsAbs(relativePath) {
		logger.Errorf("[utils] BrowseDirectory absolute path not allowed: %s", relativePath)
		return nil, "", fmt.Errorf("invalid path: absolute paths not allowed")
	}

	// Построение полного пути
	fullPath := filepath.Join(basePath, relativePath)

	// Дополнительная проверка: финальный путь должен быть внутри basePath
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		logger.Errorf("[utils] BrowseDirectory error getting absolute base path: %v", err)
		return nil, "", fmt.Errorf("internal error: failed to resolve base path")
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		logger.Errorf("[utils] BrowseDirectory error getting absolute full path: %v", err)
		return nil, "", fmt.Errorf("internal error: failed to resolve path")
	}

	if !strings.HasPrefix(absFullPath, absBasePath) {
		logger.Errorf("[utils] BrowseDirectory path outside base directory: %s", absFullPath)
		return nil, "", fmt.Errorf("invalid path: outside allowed directory")
	}

	// Проверка существования директории
	dirInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Errorf("[utils] BrowseDirectory directory not found: %s", fullPath)
			return nil, "", fmt.Errorf("directory not found: %s", relativePath)
		}
		logger.Errorf("[utils] BrowseDirectory error accessing directory: %s - %v", fullPath, err)
		return nil, "", fmt.Errorf("failed to access directory: %w", err)
	}

	// Проверка, что это действительно директория
	if !dirInfo.IsDir() {
		logger.Errorf("[utils] BrowseDirectory path is not a directory: %s", fullPath)
		return nil, "", fmt.Errorf("path is not a directory: %s", relativePath)
	}

	// Чтение содержимого директории
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		logger.Errorf("[utils] BrowseDirectory error reading directory: %s - %v", fullPath, err)
		return nil, "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Сбор информации о файлах и папках
	items := make([]FileSystemItem, 0)

	for _, entry := range entries {
		itemPath := filepath.Join(relativePath, entry.Name())
		itemFullPath := filepath.Join(fullPath, entry.Name())

		// Получаем информацию о файле/папке
		info, err := entry.Info()
		if err != nil {
			logger.Errorf("[utils] BrowseDirectory error getting info for %s: %v", itemFullPath, err)
			continue
		}

		if entry.IsDir() {
			// Это директория
			items = append(items, FileSystemItem{
				Name: entry.Name(),
				Type: FileSystemItemTypeDirectory,
				Path: itemPath,
			})
		} else {
			// Это файл - проверяем, что это поддерживаемый формат
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if supportedFormats[ext] {
				size := uint64(info.Size())
				modTime := info.ModTime()
				format := strings.TrimPrefix(ext, ".")

				items = append(items, FileSystemItem{
					Name:       entry.Name(),
					Type:       FileSystemItemTypeFile,
					Path:       itemPath,
					Size:       &size,
					Format:     &format,
					ModifiedAt: &modTime,
				})
			}
		}
	}

	// Вычисляем parent path
	parentPath := ""
	if relativePath != "" && relativePath != "." {
		parentPath = filepath.Dir(relativePath)
		if parentPath == "." {
			parentPath = ""
		}
	}

	logger.Infof("[utils] BrowseDirectory success: %s - found %d items", fullPath, len(items))

	return items, parentPath, nil
}
