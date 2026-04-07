package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/seminhnva/gin-layered-architecture/internal/auth"
	passwordService "github.com/seminhnva/gin-layered-architecture/internal/auth/password"
	"github.com/seminhnva/gin-layered-architecture/internal/bootstrap"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
	"github.com/seminhnva/gin-layered-architecture/internal/constants"
	"github.com/seminhnva/gin-layered-architecture/internal/db"
	"github.com/seminhnva/gin-layered-architecture/internal/db/sqlc"
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
	dbpool *pgxpool.Pool
}

type ModuleDeps struct {
	db              *pgxpool.Pool
	Queries         *sqlc.Queries
	PasswordService auth.Hasher
}

func NewApplication(cfg *config.Config) (*Application, error) {
	dbpool, err := db.NewPool(context.Background(), cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	deps := &ModuleDeps{
		Queries:         sqlc.New(dbpool),
		PasswordService: passwordService.NewPasswordService(),
		db:              dbpool,
	}
	r := gin.New()
	modules := []Module{
		NewUserModule(deps),
	}
	logOpts := bootstrap.NewLoggerOptions(cfg.Logger)

	httpLogger, err := logger.InitLogger(string(constants.HttpLogFilePath),
		logOpts,
	)
	if err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("init http logger: %w", err)
	}
	recoveryLogger, err := logger.InitLogger(string(constants.RecoveryLogFilePath), logOpts)
	if err != nil {
		dbpool.Close()
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
		dbpool: dbpool,
	}, nil
}

func (a *Application) Run(ctx context.Context, appLogger *zerolog.Logger) error {
	defer func() {
		if a.dbpool != nil {
			a.dbpool.Close()
		}
	}()

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
