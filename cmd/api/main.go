package main

import (
	"log"

	"github.com/seminhnva/gin-layered-architecture/internal/app"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
)

func main() {
	cfg := config.NewConfig()
	application := app.NewApplication(cfg)
	if err := application.Run(); err != nil {
		log.Fatal("Couldn't run server")
	}
}
