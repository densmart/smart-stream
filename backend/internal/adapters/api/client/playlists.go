package client

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/gin-gonic/gin"
)

// getPlaylists возвращает список всех плейлистов
func (r *ClientAPIRouter) getPlaylists(c *gin.Context) {
	var data dto.SearchPlaylistsDTO

	if err := c.BindQuery(&data); err != nil {
		logger.Debugf("[client api] getPlaylists bind query error: %s", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	playlists, totals, ucErr := usecases.SearchPlaylists(*r.oltp, data)
	if ucErr != nil {
		logger.Debugf("[client api] getPlaylists error: %s", ucErr.Message)
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

	var results []dto.PlaylistDTO
	for _, item := range playlists {
		result := dto.PlaylistDTO{
			ID:        item.ID,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			Name:      item.Name,
			Type:      item.Type,
			Poster:    item.Poster,
		}
		results = append(results, result)
	}

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}

// getPlaylistMedia возвращает список медиа в конкретном плейлисте
func (r *ClientAPIRouter) getPlaylistMedia(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.SearchMediaDTO

	if err := c.BindQuery(&data); err != nil {
		logger.Debugf("[client api] getPlaylistMedia bind query error: %s", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Фильтруем по playlist_id
	data.PlaylistID = &playlistID

	media, totals, ucErr := usecases.SearchMedia(*r.oltp, data)
	if ucErr != nil {
		logger.Debugf("[client api] getPlaylistMedia error: %s", ucErr.Message)
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

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
