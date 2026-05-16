package web

import (
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	router.Use(r.cors)

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

	media := withAuth.Group("/media")
	{
		media.POST("/", r.createMedia)
		media.GET("/:id/", r.retrieveMedia)
		media.PATCH("/:id/", r.updateMedia)
		media.DELETE("/:id/", r.deleteMedia)
		media.GET("/", r.searchMedia)
		media.GET("/browse", r.browseMediaFiles)
	}

	playlists := withAuth.Group("/playlists")
	{
		playlists.POST("/", r.createPlaylist)
		playlists.GET("/search/", r.searchPlaylistsAutocomplete)
		playlists.GET("/:id/", r.retrievePlaylist)
		playlists.PATCH("/:id/", r.updatePlaylist)
		playlists.DELETE("/:id/", r.deletePlaylist)
		playlists.GET("/", r.searchPlaylists)

		// Playlist media management
		playlists.GET("/:id/media/", r.getPlaylistMedia)
		playlists.POST("/:id/media/", r.addMediaToPlaylist)
		playlists.DELETE("/:id/media/:media_id/", r.removeMediaFromPlaylist)
		playlists.PATCH("/:id/media/:media_id/", r.updateMediaOrder)

		// Batch playlist media management
		playlists.POST("/:id/media/batch/", r.batchAddMediaToPlaylist)
		playlists.DELETE("/:id/media/batch/", r.batchRemoveMediaFromPlaylist)
		playlists.PATCH("/:id/media/batch/", r.batchUpdateMediaOrder)
	}

	upload := withAuth.Group("/upload")
	{
		upload.POST("/poster/", r.uploadPoster)
	}

	// Static files для постеров (без авторизации для публичного доступа)
	router.Static("/static/posters", viper.GetString("storage.posters-dir"))

	return router
}
