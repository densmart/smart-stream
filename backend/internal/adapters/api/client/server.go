package client

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type ClientAPIServer struct {
	server *http.Server
}

func NewClientAPIServer(handler http.Handler) *ClientAPIServer {
	server := &http.Server{
		Addr:         ":" + viper.GetString("api.client.port"),
		Handler:      handler,
		ReadTimeout:  time.Duration(15) * time.Second,
		WriteTimeout: time.Duration(15) * time.Second,
	}
	return &ClientAPIServer{server: server}
}

func (cs *ClientAPIServer) Run() error {
	return cs.server.ListenAndServe()
}

func (cs *ClientAPIServer) Stop(ctx context.Context) error {
	return cs.server.Shutdown(ctx)
}