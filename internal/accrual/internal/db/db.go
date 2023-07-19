package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB struct {
	*sqlx.DB
}

func InitDB(ctx context.Context, dsn, migrationsPath string) (*DB, error) {
	l := logger.LoggerFromContext(ctx)
	l.Debug("start Run migrations")
	err := RunMigrations(ctx, dsn, migrationsPath)

	if err != nil {
		l.Error("error Run migrations", zap.String("msg", err.Error()))
		return nil, err
	}

	l.Debug("start InitDB")
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)

	if err != nil {
		l.Error("error InitDB", zap.String("msg", err.Error()))
		return nil, err
	}
	l.Debug("success InitDB")

	return &DB{db}, nil
}

func RunMigrations(ctx context.Context, dsn, migrationsPath string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dsn)
	if err != nil {
		return fmt.Errorf("failed to get new migrate instance: %w", err)
	}

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to DB: %w", err)
		}
	}
	return nil
}
