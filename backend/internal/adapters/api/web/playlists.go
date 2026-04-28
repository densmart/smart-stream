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