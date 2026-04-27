package web

import (
	"net/http"
	"time"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/densmart/smart-stream/internal/domain/utils"
	"github.com/gin-gonic/gin"
)

func (r *WebAPIRouter) createClient(c *gin.Context) {
	var data dto.CreateClientDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}
	client, err := usecases.CreateClient(*r.oltp, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.ClientDTO{
		ID:        client.ID,
		CreatedAt: client.CreatedAt.Format(time.RFC3339),
		Login:     client.Login,
		Email:     client.Email,
		IsActive:  client.IsActive,
	}

	SuccessResponse(c, http.StatusCreated, result)
}

func (r *WebAPIRouter) retrieveClient(c *gin.Context) {
	id := c.Param("id")

	client, err := usecases.RetrieveClient(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.ClientDTO{
		ID:             client.ID,
		CreatedAt:      client.CreatedAt.Format(time.RFC3339),
		Login:          client.Login,
		Email:          client.Email,
		IsActive:       client.IsActive,
		CurrentVersion: client.CurrentVersion,
	}

	if client.LastLoginAt != nil {
		result.LastLoginAt = utils.Ptr(client.LastLoginAt.Format(time.RFC3339))
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) updateClient(c *gin.Context) {
	id := c.Param("id")
	var data dto.UpdateClientDTO

	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	client, err := usecases.UpdateClient(*r.oltp, id, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	result := dto.ClientDTO{
		ID:             client.ID,
		CreatedAt:      client.CreatedAt.Format(time.RFC3339),
		Login:          client.Login,
		Email:          client.Email,
		IsActive:       client.IsActive,
		CurrentVersion: client.CurrentVersion,
	}
	if client.LastLoginAt != nil {
		result.LastLoginAt = utils.Ptr(client.LastLoginAt.Format(time.RFC3339))
	}

	SuccessResponse(c, http.StatusOK, result)
}

func (r *WebAPIRouter) deleteClient(c *gin.Context) {
	id := c.Param("id")

	err := usecases.DeleteClient(*r.oltp, id)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusNoContent, nil)
}

func (r *WebAPIRouter) searchClients(c *gin.Context) {
	var data dto.SearchClientsDTO

	if err := c.BindQuery(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	clients, totals, err := usecases.SearchClients(*r.oltp, data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	var results []dto.ClientDTO
	for _, client := range clients {
		item := dto.ClientDTO{
			ID:             client.ID,
			CreatedAt:      client.CreatedAt.Format(time.RFC3339),
			Login:          client.Login,
			Email:          client.Email,
			IsActive:       client.IsActive,
			CurrentVersion: client.CurrentVersion,
		}
		if client.LastLoginAt != nil {
			item.LastLoginAt = utils.Ptr(client.LastLoginAt.Format(time.RFC3339))
		}
		results = append(results, item)
	}

	paginator := NewPaginator(data.BaseSearchRequestDTO, totals)
	SuccessResponseSearch(c, results, paginator.ToRepresentation())
}
