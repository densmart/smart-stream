package web

import (
	"net/http"

	"github.com/densmart/smart-stream/internal/domain/usecases"
	"github.com/gin-gonic/gin"
)

func (r *WebAPIRouter) uploadPoster(c *gin.Context) {
	// Получаем файл из запроса
	file, err := c.FormFile("file")
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, "File is required")
		return
	}

	// Используем usecase для загрузки
	uploadedFile, ucErr := usecases.UploadPoster(file)
	if ucErr != nil {
		ErrorResponse(c, ucErr.HttpCode, ucErr.Message)
		return
	}

	SuccessResponse(c, http.StatusCreated, gin.H{
		"path":     uploadedFile.Path,
		"filename": uploadedFile.Filename,
	})
}