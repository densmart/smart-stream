package client

import (
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/gin-gonic/gin"
)

type ClientAPIRouter struct {
	oltp *repo.OltpRepo
}

func NewClientAPIRouter(oltp *repo.OltpRepo) *ClientAPIRouter {
	return &ClientAPIRouter{
		oltp: oltp,
	}
}

func (r *ClientAPIRouter) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "PONG")
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in/", r.signIn)
	}

	withAuth := router.Group("/")
	withAuth.Use(r.checkAuth)
	{
		withAuth.GET("/media/", r.getUnassignedMedia)
		withAuth.GET("/media/:id/stream/", r.streamMedia)
		withAuth.GET("/playlists/", r.getPlaylists)
		withAuth.GET("/playlists/:id/media/", r.getPlaylistMedia)
	}

	return router
}