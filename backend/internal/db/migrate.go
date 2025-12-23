package db

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(ctx context.Context, dsn string) error {
	migrationsPath := os.Getenv("ATODE_MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}

	sourceURL := "file://" + filepath.ToSlash(absPath)
	m, err := migrate.New(sourceURL, normalizePostgresDSN(dsn))
	if err != nil {
		return err
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}

func normalizePostgresDSN(dsn string) string {
	// migrate's postgres driver expects scheme "postgres".
	if strings.HasPrefix(dsn, "postgresql://") {
		return "postgres" + strings.TrimPrefix(dsn, "postgresql")
	}
	return dsn
}
