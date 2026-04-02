package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seminhnva/gin-layered-architecture/internal/app"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("set up config: %v", err)
	}

	appLogger, err := logger.InitLogger(string(constants.AppLogFilePath), cfg.Logger)
	if err != nil {
		log.Fatalf("init app logger: %v", err)
	}
	application, err := app.NewApplication(cfg)
	if err != nil {
		appLogger.Fatal().Err(err).Msg("[API]Fail to initialize application")
	}
	if err := application.Run(ctx, appLogger); err != nil {
		appLogger.Fatal().Err(err).Msg("[API]Fail to run application")
	}
}
