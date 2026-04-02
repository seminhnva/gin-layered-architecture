package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	"github.com/seminhnva/gin-layered-architecture/pkg/logger"
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

func NewApplication(cfg *config.Config) (*Application, error) {
	r := gin.New()
	modules := []Module{
		NewUserModule(),
	}
	httpLogger, err := logger.InitLogger(string(constants.HttpLogFilePath), cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init http logger: %w", err)
	}
	recoveryLogger, err := logger.InitLogger(string(constants.RecoveryLogFilePath), cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init recovery logger: %w", err)
	}

	routes.SetUpRouter(cfg.CORSAllowedOrigins, r, httpLogger, recoveryLogger, getModuleRoute(modules)...)
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
	}, nil
}

func (a *Application) Run(ctx context.Context, appLogger *zerolog.Logger) error {
	errCh := make(chan error, 1)
	go func() {
		appLogger.Info().Msgf("HTTP Server listening on %s", a.config.HTTPServer.ServerAddress)
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
			appLogger.Info().Msg("Server exited gracefully")
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
