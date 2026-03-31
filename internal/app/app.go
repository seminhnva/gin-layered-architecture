package app

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	server *http.Server
}

func NewApplication(cfg *config.Config) *Application {
	r := gin.New()
	modules := []Module{
		NewUserModule(),
	}
	routes.SetUpRouter(r, getModuleRoute(modules)...)
	server := &http.Server{
		Addr:              cfg.HTTPServer.ServerAddress,
		Handler:           r,
		ReadHeaderTimeout: cfg.HTTPServer.ReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      cfg.HTTPServer.WriteTimeout,
		IdleTimeout:       cfg.HTTPServer.IdleTimeout,
	}

	return &Application{
		config: cfg,
		router: r,
		module: modules,
		server: server,
	}
}

func (a *Application) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("HTTP Server listening on %s", a.config.HTTPServer.ServerAddress)
		errCh <- a.server.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.HTTPServer.ShutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server exited gracefully")
			return nil
		}
		return err
	}
}

func getModuleRoute(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}
	return routeList
}
