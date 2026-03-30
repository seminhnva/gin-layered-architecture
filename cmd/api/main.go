package main

import (
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seminhnva/gin-layered-architecture/internal/app"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	applog "github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}

	if err := applog.Setup(cfg); err != nil {
		stdlog.Fatalf("setup logger: %v", err)
	}

	application := app.NewApplication(cfg)
	if err := application.Run(ctx); err != nil {
		applog.Fatalf("run application: %v", err)
	}
}
