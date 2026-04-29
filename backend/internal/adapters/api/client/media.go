package client

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/densmart/smart-stream/internal/domain/utils"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/gin-gonic/gin"
)

// getUnassignedMedia возвращает список медиа, которые НЕ привязаны к плейлистам
func (r *ClientAPIRouter) getUnassignedMedia(c *gin.Context) {
	var data dto.SearchMediaDTO
	if err := c.BindQuery(&data); err != nil {
		logger.Debugf("[client api] getUnassignedMedia bind query error: %s", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Устанавливаем флаг для получения только медиа без плейлистов
	data.OnlyUnassigned = utils.Ptr(true)

	media, totals, ucErr := usecases.SearchMedia(*r.oltp, data)
	if ucErr != nil {
		logger.Debugf("[client api] getUnassignedMedia error: %s", ucErr.Message)
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

	// Преобразуем модели в DTO
	var results []dto.MediaDTO
	for _, item := range media {
		result := dto.MediaDTO{
			ID:         item.ID,
			CreatedAt:  item.CreatedAt.Format(time.RFC3339),
			PlaylistID: item.PlaylistID,
			Name:       item.Name,
			Poster:     item.Poster,
			Format:     item.Format,
			Path:       item.Path,
			Duration:   item.Duration,
			Size:       item.Size,
			Order:      item.Order,
		}
		results = append(results, result)
	}

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}
