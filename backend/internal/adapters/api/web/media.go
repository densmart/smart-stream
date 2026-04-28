package web

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/gin-gonic/gin"
)

func (r *WebAPIRouter) createMedia(c *gin.Context) {
	var data dto.CreateMediaDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}
	media, err := usecases.CreateMedia(*r.oltp, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.MediaDTO{
		ID:         media.ID,
		CreatedAt:  media.CreatedAt.Format(time.RFC3339),
		PlaylistID: media.PlaylistID,
		Name:       media.Name,
		Poster:     media.Poster,
		Format:     media.Format,
		Path:       media.Path,
		Duration:   media.Duration,
		Size:       media.Size,
		Order:      media.Order,
	}

	SuccessResponse(c, http.StatusCreated, result)
}

func (r *WebAPIRouter) retrieveMedia(c *gin.Context) {
	id := c.Param("id")

	media, err := usecases.RetrieveMedia(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.MediaDTO{
		ID:         media.ID,
		CreatedAt:  media.CreatedAt.Format(time.RFC3339),
		PlaylistID: media.PlaylistID,
		Name:       media.Name,
		Poster:     media.Poster,
		Format:     media.Format,
		Path:       media.Path,
		Duration:   media.Duration,
		Size:       media.Size,
		Order:      media.Order,
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) updateMedia(c *gin.Context) {
	id := c.Param("id")
	var data dto.UpdateMediaDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	media, err := usecases.UpdateMedia(*r.oltp, id, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.MediaDTO{
		ID:         media.ID,
		CreatedAt:  media.CreatedAt.Format(time.RFC3339),
		PlaylistID: media.PlaylistID,
		Name:       media.Name,
		Poster:     media.Poster,
		Format:     media.Format,
		Path:       media.Path,
		Duration:   media.Duration,
		Size:       media.Size,
		Order:      media.Order,
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) deleteMedia(c *gin.Context) {
	id := c.Param("id")

	err := usecases.DeleteMedia(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusNoContent, nil)
}

func (r *WebAPIRouter) searchMedia(c *gin.Context) {
	var data dto.SearchMediaDTO

	if err := c.BindQuery(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	media, totals, err := usecases.SearchMedia(*r.oltp, data)
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

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}