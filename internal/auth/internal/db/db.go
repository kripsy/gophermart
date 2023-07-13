package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/kripsy/gophermart/internal/auth/internal/logger"
	"github.com/kripsy/gophermart/internal/auth/internal/models"
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

func (db *DB) RegisterUser(ctx context.Context, username, hashPassword string, id int) error {
	l := logger.LoggerFromContext(ctx)
	ctx, canlcel := context.WithTimeout(ctx, time.Second)
	defer canlcel()
	l.Debug("usecase RegisterUser")

	tx, err := db.DB.Begin()
	if err != nil {
		l.Error("failed to Begin Tx in RegisterUser", zap.String("msg", err.Error()))
		return err
	}

	defer tx.Rollback()

	queryBuilder := squirrel.
		Insert("users").
		Columns("id", "username", "password").
		Values(id, username, hashPassword).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := queryBuilder.ToSql()

	if err != nil {
		l.Error("failed to build sql in RegisterUser", zap.String("msg", err.Error()))
		return err
	}

	l.Debug("success build sql", zap.String("msg", sql))

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		l.Error("failed to exec sql in RegisterUser", zap.String("msg", err.Error()))
		return err
	}

	tx.Commit()
	l.Debug("success commit RegisterUser")
	return nil
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
	queryBuilder := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where(squirrel.Eq{"username": username}).
		Suffix(")").
		PlaceholderFormat(squirrel.Dollar)
	sql, args, err := queryBuilder.ToSql()

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

func (db *DB) GetNextUserID(ctx context.Context) (int, error) {
	l := logger.LoggerFromContext(ctx)
	ctx, canlcel := context.WithTimeout(ctx, time.Second)
	defer canlcel()
	l.Debug("start GetNextUserID")

	tx, err := db.DB.Begin()
	if err != nil {
		l.Error("failed to Begin Tx in getNextUserID", zap.String("msg", err.Error()))
		return 0, err
	}

	defer tx.Rollback()

	queryBuilder := squirrel.
		Select("MAX(id)+1").
		From("users")

	qbsql, _, err := queryBuilder.ToSql()

	if err != nil {
		l.Error("failed to build sql in getNextUserID", zap.String("msg", err.Error()))
		return 0, err
	}

	l.Debug("success build sql", zap.String("msg", qbsql))

	stmt, err := tx.PrepareContext(ctx, qbsql)
	if err != nil {
		l.Error("failed to PrepareContext stmt in getNextUserID", zap.String("msg", err.Error()))
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx)
	var nextID sql.NullInt32
	err = row.Scan(&nextID)
	if err != nil {
		l.Error("failed to scan getNextUserID", zap.String("msg", err.Error()))
		return 0, err
	}

	tx.Commit()
	l.Debug("success commit getNextUserID")
	if !nextID.Valid {
		return 1, nil
	}
	return int(nextID.Int32), nil
}

func (db *DB) GetUserHashPassword(ctx context.Context, username string) (int, string, error) {
	l := logger.LoggerFromContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	l.Debug("start GetUserHashPassword")

	tx, err := db.DB.Begin()
	if err != nil {
		l.Error("failed to Begin Tx in GetUserHashPassword", zap.String("msg", err.Error()))
		return 0, "", err
	}

	defer tx.Rollback()

	var userID int
	var hashPassword string
	queryBuilder := squirrel.Select("id, password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar)
	bsql, args, err := queryBuilder.ToSql()

	if err != nil {
		l.Error("failed to build sql in GetUserHashPassword", zap.String("msg", err.Error()))
		return 0, "", err
	}

	l.Debug("success build sql in GetUserHashPassword", zap.String("msg", bsql))

	row := tx.QueryRowContext(ctx, bsql, args...)

	err = row.Scan(&userID, &hashPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Debug("error compare username and pwd", zap.String("msg", username))
			return 0, "", models.NewUserLoginError(username)
		}
		l.Error("failed scan userExists", zap.String("msg", err.Error()))
		return 0, "", err
	}

	l.Debug("success get hash password ->", zap.String("msg", hashPassword))
	tx.Commit()
	return userID, hashPassword, nil
}
