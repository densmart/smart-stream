package web

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type WebAPIServer struct {
	server *http.Server
}

func NewWebAPIServer(handler http.Handler) *WebAPIServer {
	server := &http.Server{
		Addr:         ":" + viper.GetString("api.web.port"),
		Handler:      handler,
		ReadTimeout:  time.Duration(15) * time.Second,
		WriteTimeout: time.Duration(15) * time.Second,
	}
	return &WebAPIServer{server: server}
}

func (bs *WebAPIServer) Run() error {
	return bs.server.ListenAndServe()
}

func (bs *WebAPIServer) Stop(ctx context.Context) error {
	return bs.server.Shutdown(ctx)
}
