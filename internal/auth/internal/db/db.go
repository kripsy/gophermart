package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/utils"
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
	fmt.Println(migrationsPath)
	fmt.Println(dsn)
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

func (db *DB) RegisterUser(ctx context.Context, username, password string) (string, time.Time, error) {
	l := logger.LoggerFromContext(ctx)
	l.Debug("usecase RegisterUser")
	ctx, canlcel := context.WithTimeout(ctx, time.Second)
	defer canlcel()
	userIsExists, err := db.IsUserExists(ctx, username)
	if err != nil {
		l.Error("error check IsUserExists in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	l.Debug("user is exist?", zap.Bool("msg", userIsExists))

	token := ""
	expTime := time.Now().Add(100500 * time.Hour)
	hash, err := utils.GetHash(ctx, password)
	if err != nil {
		l.Error("error GetHash in RegisterUser", zap.String("msg", err.Error()))
		return "", time.Time{}, err
	}
	l.Debug("got hash in RegisterUser", zap.String("msg", string(hash)))
	return token, expTime, fmt.Errorf("not implemented yet")
}

func (db *DB) IsUserExists(ctx context.Context, username string) (bool, error) {
	l := logger.LoggerFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	l.Debug("start IsUserExists")

	tx, err := db.DB.Begin()
	if err != nil {
		l.Error("failed to Begin Tx in IsUserExists", zap.String("msg", err.Error()))
		return false, err
	}

	defer tx.Rollback()

	var userExists bool
	selectBuilder := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where(squirrel.Eq{"username": username}).
		Suffix(")").
		PlaceholderFormat(squirrel.Dollar)
	sql, args, err := selectBuilder.ToSql()

	if err != nil {
		l.Error("failed to build sql in IsUserExists", zap.String("msg", err.Error()))
		return false, err
	}

	l.Debug("success build sql", zap.String("msg", sql))

	row := tx.QueryRowContext(ctx, sql, args...)

	err = row.Scan(&userExists)

	if err != nil {
		l.Error("failed scan userExists", zap.String("msg", err.Error()))
		return false, err
	}
	l.Debug("success scan userExists, value ->", zap.Bool("msg", userExists))
	tx.Commit()
	return userExists, nil
}
