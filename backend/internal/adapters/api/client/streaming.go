package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/gin-gonic/gin"
)

// streamMedia обрабатывает стриминг видео с поддержкой HTTP Range requests
func (r *ClientAPIRouter) streamMedia(c *gin.Context) {
	mediaID := c.Param("id")

	// Получаем информацию о медиа-файле через use case
	streamInfo, ucErr := usecases.GetMediaStreamInfo(*r.oltp, mediaID)
	if ucErr != nil {
		logger.Debugf("[client api] streamMedia error: %s", ucErr.Message)
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

	// Открываем файл
	file, err := os.Open(streamInfo.FilePath)
	if err != nil {
		logger.Errorf("[client api] streamMedia error opening file: %s", err.Error())
		ErrorResponse(c, http.StatusInternalServerError, "Error opening media file")
		return
	}
	defer file.Close()

	// Устанавливаем заголовки для поддержки Range requests
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Type", streamInfo.ContentType)

	// Парсим Range заголовок
	rangeHeader := c.GetHeader("Range")
	ranges, err := usecases.ParseRange(rangeHeader, streamInfo.FileSize)

	// Если нет Range заголовка - отдаём весь файл
	if rangeHeader == "" || ranges == nil {
		c.Header("Content-Length", strconv.FormatInt(streamInfo.FileSize, 10))
		c.Status(http.StatusOK)
		io.Copy(c.Writer, file)
		return
	}

	// Ошибка парсинга Range
	if err != nil {
		logger.Debugf("[client api] streamMedia invalid range header: %s", err.Error())
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", streamInfo.FileSize))
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Поддерживаем только один диапазон (multipart ranges не поддерживаем)
	if len(ranges) != 1 {
		logger.Debugf("[client api] streamMedia multiple ranges not supported")
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", streamInfo.FileSize))
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	ra := ranges[0]

	// Проверяем валидность диапазона
	if ra.Start >= streamInfo.FileSize {
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", streamInfo.FileSize))
		c.Status(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Корректируем end если он больше размера файла
	if ra.End >= streamInfo.FileSize {
		ra.End = streamInfo.FileSize - 1
	}

	// Вычисляем длину контента
	contentLength := ra.End - ra.Start + 1

	// Устанавливаем позицию чтения в файле
	if _, err := file.Seek(ra.Start, io.SeekStart); err != nil {
		logger.Errorf("[client api] streamMedia error seeking file: %s", err.Error())
		ErrorResponse(c, http.StatusInternalServerError, "Error reading media file")
		return
	}

	// Устанавливаем заголовки для partial content
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", ra.Start, ra.End, streamInfo.FileSize))
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Status(http.StatusPartialContent)

	// Отправляем указанный диапазон байт
	io.CopyN(c.Writer, file, contentLength)
}