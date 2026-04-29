package client

import (
	"net/http"
	"strings"

	"github.com/densmart/smart-stream/internal/domain/utils"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	authorizationHeader = "Authorization"
)

func (r *ClientAPIRouter) checkAuth(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		logger.Debugf("[middleware] checkAuth: empty authorization header")
		ErrorResponse(c, http.StatusUnauthorized, "empty authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || (len(headerParts) == 2 && headerParts[0] != "Bearer") {
		logger.Debugf("[middleware] checkAuth: incorrect authorization header")
		ErrorResponse(c, http.StatusUnauthorized, "incorrect authorization header")
		return
	}

	jwtToken := utils.NewJwtToken(viper.GetString("app.jwt-secret"))
	jwtToken.Access = headerParts[1]

	claims, err := jwtToken.GetAccessClaims()
	if err != nil {
		logger.Debugf("[middleware] checkAuth: invalid token: %s", err.Error())
		ErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set("user_name", claims.UserName)
	c.Next()
}