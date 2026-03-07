package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	module []Module
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.Default()
	loadEnv()
	modules := []Module{
		NewUserModule(),
	}
	routes.SetUpRouter(r, getModuleRoute(modules)...)
	return &Application{
		config: cfg,
		router: r,
		module: modules,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAdress)
}

func getModuleRoute(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}

func loadEnv() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found")
	}

}
