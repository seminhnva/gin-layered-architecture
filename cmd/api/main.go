package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seminhnva/gin-layered-architecture/internal/app"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("set up config: %v", err)
	}

	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatalf("set up application: %v", err)
	}
	if err := application.Run(ctx); err != nil {
		log.Fatalf("run application: %v", err)
	}
}
