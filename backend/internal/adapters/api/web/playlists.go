package web

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/gin-gonic/gin"
)

func (r *WebAPIRouter) createPlaylist(c *gin.Context) {
	var data dto.CreatePlaylistDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}
	playlist, err := usecases.CreatePlaylist(*r.oltp, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.PlaylistDTO{
		ID:        playlist.ID,
		CreatedAt: playlist.CreatedAt.Format(time.RFC3339),
		Name:      playlist.Name,
		Type:      playlist.Type,
		Poster:    playlist.Poster,
	}

	SuccessResponse(c, http.StatusCreated, result)
}

func (r *WebAPIRouter) retrievePlaylist(c *gin.Context) {
	id := c.Param("id")

	playlist, err := usecases.RetrievePlaylist(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.PlaylistDTO{
		ID:        playlist.ID,
		CreatedAt: playlist.CreatedAt.Format(time.RFC3339),
		Name:      playlist.Name,
		Type:      playlist.Type,
		Poster:    playlist.Poster,
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) updatePlaylist(c *gin.Context) {
	id := c.Param("id")
	var data dto.UpdatePlaylistDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	playlist, err := usecases.UpdatePlaylist(*r.oltp, id, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.PlaylistDTO{
		ID:        playlist.ID,
		CreatedAt: playlist.CreatedAt.Format(time.RFC3339),
		Name:      playlist.Name,
		Type:      playlist.Type,
		Poster:    playlist.Poster,
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) deletePlaylist(c *gin.Context) {
	id := c.Param("id")

	err := usecases.DeletePlaylist(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusNoContent, nil)
}

func (r *WebAPIRouter) searchPlaylists(c *gin.Context) {
	var data dto.SearchPlaylistsDTO

	if err := c.BindQuery(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	playlists, totals, err := usecases.SearchPlaylists(*r.oltp, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
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

// getPlaylistMedia returns all media in a playlist sorted by order
func (r *WebAPIRouter) getPlaylistMedia(c *gin.Context) {
	playlistID := c.Param("id")
	var filter dto.SearchMediaDTO

	if err := c.BindQuery(&filter); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	media, total, err := usecases.GetPlaylistMedia(*r.oltp, playlistID, filter)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
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

// addMediaToPlaylist adds media to a playlist with order
func (r *WebAPIRouter) addMediaToPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.AddMediaToPlaylistDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := usecases.AddMediaToPlaylist(*r.oltp, playlistID, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, nil)
}

// removeMediaFromPlaylist removes media from a playlist
func (r *WebAPIRouter) removeMediaFromPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	mediaID := c.Param("media_id")

	err := usecases.RemoveMediaFromPlaylist(*r.oltp, playlistID, mediaID)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusNoContent, nil)
}

// updateMediaOrder updates the order of media in a playlist
func (r *WebAPIRouter) updateMediaOrder(c *gin.Context) {
	playlistID := c.Param("id")
	mediaID := c.Param("media_id")
	var data dto.UpdateMediaOrderDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := usecases.UpdateMediaOrder(*r.oltp, playlistID, mediaID, data.Order)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, nil)
}

// batchAddMediaToPlaylist adds multiple media to a playlist in a single request
func (r *WebAPIRouter) batchAddMediaToPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.BatchAddMediaToPlaylistDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := usecases.BatchAddMediaToPlaylist(*r.oltp, playlistID, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, nil)
}

// batchRemoveMediaFromPlaylist removes multiple media from a playlist in a single request
func (r *WebAPIRouter) batchRemoveMediaFromPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.BatchRemoveMediaFromPlaylistDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := usecases.BatchRemoveMediaFromPlaylist(*r.oltp, playlistID, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusNoContent, nil)
}

// batchUpdateMediaOrder updates the order of multiple media in a playlist in a single request
func (r *WebAPIRouter) batchUpdateMediaOrder(c *gin.Context) {
	playlistID := c.Param("id")
	var data dto.BatchUpdateMediaOrderDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	err := usecases.BatchUpdateMediaOrder(*r.oltp, playlistID, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, nil)
}