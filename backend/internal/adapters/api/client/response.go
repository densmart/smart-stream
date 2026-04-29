package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Error      string          `json:"error"`
	Pagination *PaginationInfo `json:"pagination"`
	Result     interface{}     `json:"result"`
}

type PaginationInfo struct {
	Total uint64 `json:"total"`
	Pages uint   `json:"pages"`
}

func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	resp := APIResponse{
		Result: data,
	}
	c.JSON(statusCode, resp)
}

func SuccessResponseSearch[T any](c *gin.Context, data []T, pagination PaginationInfo) {
	results := data
	if results == nil {
		results = []T{}
	}
	resp := APIResponse{
		Pagination: &pagination,
		Result:     results,
	}
	c.JSON(http.StatusOK, resp)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	resp := APIResponse{
		Error: message,
	}
	c.AbortWithStatusJSON(statusCode, resp)
}