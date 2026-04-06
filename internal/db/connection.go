package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seminhnva/gin-layered-architecture/internal/config"
)

const (
	defaultMaxConns        = 50
	defaultMinConns        = 5
	defaultMaxConnLifetime = 30 * time.Minute
	defaultMaxConnIdleTime = 5 * time.Minute
	defaultHealthCheck     = 1 * time.Minute
	defaultPingTimeout     = 5 * time.Second
)

func NewPool(ctx context.Context, dbConfig config.DatabaseConfig) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(DSN(dbConfig))
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	conf.MaxConns = defaultMaxConns
	conf.MinConns = defaultMinConns
	conf.MaxConnLifetime = defaultMaxConnLifetime
	conf.MaxConnIdleTime = defaultMaxConnIdleTime
	conf.HealthCheckPeriod = defaultHealthCheck

	dbpool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, defaultPingTimeout)
	defer cancel()

	if err := dbpool.Ping(pingCtx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return dbpool, nil
}

func DSN(cfg config.DatabaseConfig) string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMod)
	return dsn
}
