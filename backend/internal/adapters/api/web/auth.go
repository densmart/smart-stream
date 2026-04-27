package web

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/adapters/dto"
	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/gin-gonic/gin"
)

func (r *WebAPIRouter) signIn(c *gin.Context) {
	var data dto.SignInRequestDTO
	if err := c.BindJSON(&data); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request data")
		return
	}
	result, err := usecases.SignIn(data)
	if err != nil {
		ErrorResponse(c, err.HttpCode, err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, result)
}
