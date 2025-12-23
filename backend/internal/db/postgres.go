package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	PingTimeout     time.Duration
}

func LoadConfig() (Config, error) {
	dsn := os.Getenv("ATODE_DATABASE_URL")
	if dsn == "" {
		return Config{}, errors.New("ATODE_DATABASE_URL is required")
	}

	maxOpenConns, err := envInt("ATODE_DB_MAX_OPEN_CONNS", 10)
	if err != nil {
		return Config{}, err
	}
	maxIdleConns, err := envInt("ATODE_DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return Config{}, err
	}

	connMaxLifetime, err := envDuration("ATODE_DB_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}
	connMaxIdleTime, err := envDuration("ATODE_DB_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return Config{}, err
	}
	pingTimeout, err := envDuration("ATODE_DB_PING_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}

	return Config{
		DSN:             dsn,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		ConnMaxIdleTime: connMaxIdleTime,
		PingTimeout:     pingTimeout,
	}, nil
}

func Open(ctx context.Context) (*sql.DB, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return OpenWithConfig(ctx, cfg)
}

func OpenWithConfig(ctx context.Context, cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func HealthCheck(ctx context.Context, db *sql.DB) error {
	return db.PingContext(ctx)
}

func envInt(key string, defaultValue int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func envDuration(key string, defaultValue time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, err
	}
	return d, nil
}
