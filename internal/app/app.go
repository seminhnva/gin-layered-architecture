package app

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	applog "github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	module []Module
}

func NewApplication(cfg *config.Config) *Application {
	router := gin.New()
	modules := []Module{
		NewUserModule(),
	}

	routes.SetUpRouter(router, cfg, getModuleRoute(modules)...)

	server := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	return &Application{
		config: cfg,
		router: router,
		server: server,
		module: modules,
	}
}

func (a *Application) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		applog.Infof("HTTP server listening on %s", a.config.ServerAddress)
		errCh <- a.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		applog.Infof("Shutting down HTTP server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.config.ShutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}

		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
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
