package web

import (
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/gin-gonic/gin"
)

type WebAPIRouter struct {
	oltp *repo.OltpRepo
}

func NewWebAPIRouter(oltp *repo.OltpRepo) *WebAPIRouter {
	return &WebAPIRouter{
		oltp: oltp,
	}
}

func (r *WebAPIRouter) InitRoutes() *gin.Engine {
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

	clients := withAuth.Group("/clients")
	{
		clients.POST("/", r.createClient)
		clients.GET("/:id/", r.retrieveClient)
		clients.PATCH("/:id/", r.updateClient)
		clients.DELETE("/:id/", r.deleteClient)
		clients.GET("/", r.searchClients)
	}

	return router
}
