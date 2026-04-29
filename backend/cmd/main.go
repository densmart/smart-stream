package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/densmart/smart-stream/internal/adapters/api/client"
	"github.com/densmart/smart-stream/internal/adapters/api/web"
	"github.com/densmart/smart-stream/internal/domain/repo"
	"github.com/densmart/smart-stream/pkg/configger"
	"github.com/densmart/smart-stream/pkg/logger"
	"github.com/spf13/viper"
)

func main() {
	configger.InitConfig(configger.DefaultCfgPath, "config", "yaml")
	logger.InitLogger()
	logger.Infof("[main] logger initialized. level: %s", logger.GetLevel())

	appCtx, cancel := context.WithCancel(context.Background())

	logger.Infof("[main] starting OLTP DB connection...")
	oltpRepo, err := repo.NewOltpRepo(appCtx, "postgres")
	if err != nil {
		logger.Fatalf("[main] error starting Repo: %s", err.Error())
	}
	logger.Infof("[main] DB connection established")

	// init web API
	boAPIServer := web.NewWebAPIServer(web.NewWebAPIRouter(&oltpRepo).InitRoutes())
	logger.Infof("[main] starting Web API http server -> 0.0.0.0:%s", viper.GetString("api.web.port"))
	go func() {
		if err = boAPIServer.Run(); err != nil {
			logger.Errorf("[main] error starting Web API server: %s", err.Error())
		}
	}()

	// init client API
	clientAPIServer := client.NewClientAPIServer(client.NewClientAPIRouter(&oltpRepo).InitRoutes())
	logger.Infof("[main] starting Client API http server -> 0.0.0.0:%s", viper.GetString("api.client.port"))
	go func() {
		if err = clientAPIServer.Run(); err != nil {
			logger.Errorf("[main] error starting Client API server: %s", err.Error())
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// wait for shutdown
	<-quit

	logger.Infof("[main] stopping app...")

	// stop web API server
	logger.Infof("...stop Web API server...")
	if err = boAPIServer.Stop(appCtx); err != nil {
		logger.Debugf("[main] error stopping Web API server: %s", err.Error())
	}
	logger.Infof("...Web API server stopped...")

	// stop client API server
	logger.Infof("...stop Client API server...")
	if err = clientAPIServer.Stop(appCtx); err != nil {
		logger.Debugf("[main] error stopping Client API server: %s", err.Error())
	}
	logger.Infof("...Client API server stopped...")

	// cancel context
	cancel()

	// close database connections
	logger.Infof("[main] closing DB connections...")
	oltpRepo.Close()
	logger.Infof("[main] app stopped")
}
