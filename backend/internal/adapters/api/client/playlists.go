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
			ID:          item.ID,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			Name:        item.Name,
			Type:        item.Type,
			Poster:      item.Poster,
			ParentID:    item.ParentID,
			HasChildren: item.HasChildren,
		}
		results = append(results, result)
	}

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}

// getPlaylistChildren возвращает дочерние плейлисты для указанного плейлиста
func (r *ClientAPIRouter) getPlaylistChildren(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.SearchPlaylistsDTO

	if err := c.BindQuery(&data); err != nil {
		logger.Debugf("[client api] getPlaylistChildren bind query error: %s", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	// Устанавливаем фильтр по parent_id
	data.ParentID = &playlistID

	playlists, totals, ucErr := usecases.SearchPlaylists(*r.oltp, data)
	if ucErr != nil {
		logger.Debugf("[client api] getPlaylistChildren error: %s", ucErr.Message)
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

	var results []dto.PlaylistDTO
	for _, item := range playlists {
		result := dto.PlaylistDTO{
			ID:          item.ID,
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
			Name:        item.Name,
			Type:        item.Type,
			Poster:      item.Poster,
			ParentID:    item.ParentID,
			HasChildren: item.HasChildren,
		}
		results = append(results, result)
	}

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}

// getPlaylistMedia возвращает список медиа в конкретном плейлисте
func (r *ClientAPIRouter) getPlaylistMedia(c *gin.Context) {
	playlistID := c.Param("id")
	var filter dto.SearchMediaDTO

	if err := c.BindQuery(&filter); err != nil {
		logger.Debugf("[client api] getPlaylistMedia bind query error: %s", err.Error())
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	media, total, ucErr := usecases.GetPlaylistMedia(*r.oltp, playlistID, filter)
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

	paginator := NewPaginator(filter.BaseSearchRequestDTO, total)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}
