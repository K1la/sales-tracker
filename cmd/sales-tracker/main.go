package main

import (
	analyticsHandler "github.com/K1la/sales-tracker/internal/api/handlers/analytics"
	itemHandler "github.com/K1la/sales-tracker/internal/api/handlers/items"
	"github.com/K1la/sales-tracker/internal/api/router"
	"github.com/K1la/sales-tracker/internal/api/server"
	"github.com/K1la/sales-tracker/internal/config"
	"github.com/K1la/sales-tracker/internal/repository"
	analyticsRepo "github.com/K1la/sales-tracker/internal/repository/analytics"
	itemRepo "github.com/K1la/sales-tracker/internal/repository/items"
	analyticsService "github.com/K1la/sales-tracker/internal/service/analytics"
	itemService "github.com/K1la/sales-tracker/internal/service/items"

	"github.com/wb-go/wbf/zlog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// инициализация глобального логгера
	zlog.InitConsole()
	// присваиваем глобальный логгер
	lg := zlog.Logger

	cfg := config.Init()

	// TODO: доделать инициализацию
	db := repository.NewDB(cfg)

	repoItem := itemRepo.New(db, lg)
	repoAnalytics := analyticsRepo.New(db, lg)

	serviceItem := itemService.New(repoItem, lg)
	serviceAnalytics := analyticsService.New(repoAnalytics, lg)

	handlerItem := itemHandler.New(serviceItem, lg)
	handlerAnalytics := analyticsHandler.New(serviceAnalytics, lg)

	r := router.New(handlerItem, handlerAnalytics)
	s := server.New(cfg.HTTPServer.Address, r)

	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting down signal %v. Shutting down...", sig)
		//cancel()
	}()

	if err := s.ListenAndServe(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to start server")
	}
	zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
}
