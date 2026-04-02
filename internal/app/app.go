package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/routes"
	applogger "github.com/seminhnva/gin-layered-architecture/pkg/logger"
)

type logFilePath string

const (
	httpLogFilePath     logFilePath = "app.log"
	recoveryLogFilePath logFilePath = "recovery.log"
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
	httpLogger, err := initLogger(string(httpLogFilePath), cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("init http logger: %w", err)
	}
	recoveryLogger, err := initLogger(string(recoveryLogFilePath), cfg.Logger)
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

func initLogger(logFilePath string, logConfig config.LogConfig) (*zerolog.Logger, error) {
	logger, err := applogger.NewFileLogger(
		logConfig.LogFilePath+logFilePath,
		logConfig.LogLevel,
		logConfig.LogMaxSizeMB,
		logConfig.LogMaxBackups,
		logConfig.LogMaxAgeDays,
		logConfig.LogCompress,
		logConfig.LocalTime,
	)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return logger, nil
}
