package main

import (
	"context"
	"github.com/K1la/sales-tracker/internal/api/router"
	"github.com/K1la/sales-tracker/internal/api/server"
	"github.com/K1la/sales-tracker/internal/config"
	"github.com/K1la/sales-tracker/internal/repository"

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
	dbItem := repository.NewDB(cfg)
	dbAnalytics := repository.NewDB(cfg)

	repoItem := repository.New(db)
	repoAnalytics := repository.New(db)

	itemService := service.New(repo)
	analyticsService := service.New(repo)

	itemHandler := handler.New(srvc)
	analyticsHandler := handler.New(srvc)
	r := router.New(hndlr)
	s := server.New(cfg.HTTPServer.Address, r)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting down signal %v. Shutting down...", sig)
		cancel()
	}()

	if err := s.ListenAndServe(); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to start server")
	}
	zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
}
