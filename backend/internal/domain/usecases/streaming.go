package usecases

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/densmart/smart-stream/internal/domain/models"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/spf13/viper"

	"github.com/densmart/smart-stream/pkg/logger"
)

// MIME types для различных видео форматов
var videoMimeTypes = map[string]string{
	".mp4":  "video/mp4",
	".mkv":  "video/x-matroska",
	".avi":  "video/x-msvideo",
	".mov":  "video/quicktime",
	".webm": "video/webm",
	".flv":  "video/x-flv",
	".wmv":  "video/x-ms-wmv",
	".m4v":  "video/x-m4v",
	".3gp":  "video/3gpp",
	".ts":   "video/mp2t",
}

// MediaStreamInfo содержит информацию о медиа-файле для стриминга
type MediaStreamInfo struct {
	Media       *models.Media
	FilePath    string
	FileSize    int64
	ContentType string
}

// HTTPRange описывает диапазон байтов для чтения
type HTTPRange struct {
	Start int64
	End   int64
}

// GetMediaStreamInfo получает информацию о медиа-файле для стриминга
func GetMediaStreamInfo(oltp repo.OltpRepo, mediaID string) (*MediaStreamInfo, *UCError) {
	// Получаем информацию о медиа из БД
	media, dbErr := oltp.RetrieveMedia(mediaID)
	if dbErr != nil {
		return nil, &UCError{
			Code:     240,
			HttpCode: http.StatusNotFound,
			Message:  "Media not found",
		}
	}

	// Получаем базовую директорию для медиа-файлов
	mediaDir := viper.GetString("storage.media-dir")

	// Формируем полный путь к файлу
	// media.Path может быть как абсолютным путём, так и относительным
	var fullPath string
	if filepath.IsAbs(media.Path) {
		fullPath = media.Path
	} else {
		fullPath = filepath.Join(mediaDir, media.Path)
	}

	// Проверяем существование файла
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debugf("[usecases] GetMediaStreamInfo file not found: %s", fullPath)
			return nil, &UCError{
				Code:     241,
				HttpCode: http.StatusNotFound,
				Message:  "Media file not found on disk",
			}
		}
		logger.Errorf("[usecases] GetMediaStreamInfo error accessing file: %s", err.Error())
		return nil, &UCError{
			Code:     242,
			HttpCode: http.StatusInternalServerError,
			Message:  "Error accessing media file",
		}
	}

	// Определяем MIME-тип на основе расширения
	ext := strings.ToLower(filepath.Ext(fullPath))
	contentType := videoMimeTypes[ext]
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return &MediaStreamInfo{
		Media:       &media,
		FilePath:    fullPath,
		FileSize:    fileInfo.Size(),
		ContentType: contentType,
	}, nil
}

// ParseRange парсит HTTP Range заголовок
func ParseRange(rangeHeader string, size int64) ([]HTTPRange, error) {
	if rangeHeader == "" {
		return nil, nil
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range header format")
	}

	ranges := []HTTPRange{}
	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")

	// Разбиваем на отдельные диапазоны (разделённые запятой)
	for _, ra := range strings.Split(rangeHeader, ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}

		// Парсим диапазон "start-end"
		parts := strings.Split(ra, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}

		var start, end int64
		var err error

		// Парсим start
		if parts[0] != "" {
			start, err = strconv.ParseInt(parts[0], 10, 64)
			if err != nil || start < 0 {
				return nil, fmt.Errorf("invalid start value")
			}
		} else {
			start = 0
		}

		// Парсим end
		if parts[1] != "" {
			end, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil || end < 0 {
				return nil, fmt.Errorf("invalid end value")
			}
		} else {
			end = size - 1
		}

		// Проверяем что start <= end
		if start > end {
			return nil, fmt.Errorf("start greater than end")
		}

		ranges = append(ranges, HTTPRange{
			Start: start,
			End:   end,
		})
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("no valid ranges")
	}

	return ranges, nil
}