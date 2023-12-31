package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kripsy/gophermart/internal/accrual/internal/config"
	"github.com/kripsy/gophermart/internal/accrual/internal/logger"
	"github.com/kripsy/gophermart/internal/accrual/internal/models"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store interface {
	PutOrder(ctx context.Context, number int64) (models.Order, error)
	GetOrder(ctx context.Context, number int64) (models.Order, error)
}

type DBStorage struct{}

var s Store = &DBStorage{}

func GetStorage() Store {
	return s
}

func (s *DBStorage) PutOrder(ctx context.Context, number int64) (models.Order, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("PutOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return models.Order{}, err
	}

	var ID int64
	var Number int64
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "INSERT INTO public.accrual (number, status, accrual) VALUES ($1, $2, $3) ON CONFLICT (number) DO UPDATE SET number=EXCLUDED.number RETURNING accrual.id, accrual.number, accrual.status, accrual.accrual, accrual.uploaded_at, accrual.processed_at;", number, models.StatusProcessed, number%1000).Scan(&ID, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.Number = Number
	order.Status = Status
	order.Accrual = Accrual
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}

func (s *DBStorage) GetOrder(ctx context.Context, number int64) (models.Order, error) {

	l := logger.LoggerFromContext(ctx)
	l.Info("GetOrder")
	cfg := config.GetConfig()

	conn, err := pgx.Connect(ctx, cfg.DatabaseAddress)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			l.Error("Unable close to database: ", zap.String("msg", err.Error()))
		}
	}(conn, ctx)
	if err != nil {
		l.Error("Unable to connect to database: %v\n", zap.String("msg", err.Error()))
		return models.Order{}, err
	}

	var ID int64
	var Number int64
	var Status string
	var Accrual int
	var UploadedAt pgtype.Timestamptz
	var ProcessedAt pgtype.Timestamptz

	err = conn.QueryRow(ctx, "select * from public.accrual where number=$1;", number).Scan(&ID, &Number, &Status, &Accrual, &UploadedAt, &ProcessedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Order{}, models.ErrNoAccrual()
	}

	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{}

	order.ID = ID
	order.Number = Number
	order.Status = Status
	order.Accrual = Accrual
	order.UploadedAt = UploadedAt
	order.ProcessedAt = ProcessedAt

	return order, nil
}
